// Package pins reads and rewrites the places a program version is written in this
// repository: the ARG VERSION of its Dockerfile, the version input of its GitHub
// workflow, and its row in the README table.
package pins

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// Kind is the sort of file a version is pinned in.
type Kind string

const (
	// KindDockerfile is dockerfiles/<program>/Dockerfile.
	KindDockerfile Kind = "Dockerfile"
	// KindWorkflow is .github/workflows/<program>.yml.
	KindWorkflow Kind = "workflow"
	// KindReadme is a program row of the README table.
	KindReadme Kind = "README"
)

const (
	dockerfileDir = "dockerfiles"
	workflowDir   = ".github/workflows"
	readmePath    = "README.md"
	// The two workflows that are not tied to a single program.
	workflowReusable = "reusable.yml"
	workflowHubDesc  = "dockerhub-description.yml"
)

// Source is one place where a version is written.
type Source struct {
	Kind Kind
	// Path is relative to the repository root.
	Path string
	// Value is the version found there.
	Value string
	// Line is the 1-indexed line holding the version, 0 for a README row where the
	// whole line is rewritten.
	Line int
}

// Set groups every pin of a program.
type Set struct {
	Program string
	Sources []Source
}

// ErrNoPin is returned when a program has no Dockerfile pin, which means this
// repository does not hold it.
var ErrNoPin = errors.New("no version pin found")

// Version returns the version pinned by the Dockerfile, the reference value: the
// build reads it, the workflow and the README are expected to agree with it.
func (s Set) Version() (string, error) {
	for _, src := range s.Sources {
		if src.Kind == KindDockerfile {
			return src.Value, nil
		}
	}
	return "", fmt.Errorf("%s: %w: expected dockerfiles/%s/Dockerfile with an ARG VERSION line",
		s.Program, ErrNoPin, s.Program)
}

// Drift lists the sources whose version differs from the Dockerfile one. It is worth
// reporting before anything gets rewritten.
func (s Set) Drift() []Source {
	want, err := s.Version()
	if err != nil {
		return nil
	}
	var out []Source
	for _, src := range s.Sources {
		if src.Value != want {
			out = append(out, src)
		}
	}
	return out
}

// versionLine matches a line holding a version, and says which capturing group
// holds it. Trailing blanks and a carriage return are allowed outside the version:
// the gopls Dockerfile is written with CRLF line endings.
type versionLine struct {
	re    *regexp.Regexp
	group int
}

func patternFor(kind Kind) versionLine {
	switch kind {
	case KindDockerfile:
		return versionLine{regexp.MustCompile(`^ARG VERSION=([^\s]+)[ \t]*\r?$`), 1}
	case KindWorkflow:
		return versionLine{regexp.MustCompile(`^(\s*version:[ \t]*)([^\s]+)[ \t]*\r?$`), 2}
	default:
		return versionLine{}
	}
}

// readmeRow matches a program row of the README table: the program name in the first
// cell and the version in the second, both wrapped in backticks and links.
var readmeRow = regexp.MustCompile("^\\|\\s*\\[`([^`]+)`\\]\\([^)]*\\)\\s*\\|\\s*\\[`([^`]+)`\\]\\(([^)]*)\\)")

// Read collects the pins of a program: its Dockerfile, its workflow when it has one,
// and the rows of the README table.
func Read(root, program string) (Set, error) {
	set := Set{Program: program}

	dockerfile := filepath.Join(dockerfileDir, program, "Dockerfile")
	value, line, ok := scan(readFile(root, dockerfile), patternFor(KindDockerfile))
	if !ok {
		return set, fmt.Errorf("%s: %w", program, ErrNoPin)
	}
	set.Sources = append(set.Sources, Source{
		Kind: KindDockerfile, Path: dockerfile, Value: value, Line: line,
	})

	if path, content := workflowFile(root, program); path != "" {
		if value, line, ok := scan(content, patternFor(KindWorkflow)); ok {
			set.Sources = append(set.Sources, Source{
				Kind: KindWorkflow, Path: path, Value: value, Line: line,
			})
		}
	}

	for _, line := range strings.Split(readFile(root, readmePath), "\n") {
		if m := readmeRow.FindStringSubmatch(line); m != nil && m[1] == program {
			set.Sources = append(set.Sources, Source{
				Kind: KindReadme, Path: readmePath, Value: m[2],
			})
		}
	}
	return set, nil
}

// Update rewrites the pins of a program from one version to another and returns the
// repository relative paths it changed. The Dockerfile and workflow lines are
// replaced, and every README row of the program has its occurrences of the version
// replaced, which covers the version cell, the release URL, and the image tag of a
// row such as logo-ls-v1.3.7.
func Update(root, program, from, to string) ([]string, error) {
	if program == "" {
		return nil, errors.New("no program given")
	}
	if from == "" || to == "" || from == to {
		return nil, fmt.Errorf("refusing to rewrite the pins of %s from %q to %q", program, from, to)
	}
	var changed []string

	if err := rewrite(root, filepath.Join(dockerfileDir, program, "Dockerfile"),
		patternFor(KindDockerfile), from, to, &changed); err != nil {
		return changed, err
	}
	if path, _ := workflowFile(root, program); path != "" {
		if err := rewrite(root, path, patternFor(KindWorkflow), from, to, &changed); err != nil {
			return changed, err
		}
	}
	if err := rewriteReadme(root, program, from, to, &changed); err != nil {
		return changed, err
	}
	if len(changed) == 0 {
		return nil, fmt.Errorf("%s: nothing to rewrite from %q to %q", program, from, to)
	}
	sort.Strings(changed)
	return changed, nil
}

// scan returns the first version of content matched by the pattern, with its
// 1-indexed line.
func scan(content string, vp versionLine) (string, int, bool) {
	if vp.re == nil {
		return "", 0, false
	}
	for i, line := range strings.Split(content, "\n") {
		m := vp.re.FindStringSubmatchIndex(line)
		if m == nil {
			continue
		}
		return line[m[2*vp.group]:m[2*vp.group+1]], i + 1, true
	}
	return "", 0, false
}

// rewrite replaces, in every line the pattern matches, the version by to when it
// equals from. Only the version bytes are touched, so the indentation of a workflow,
// the trailing blanks and the carriage returns of a CRLF file survive untouched.
func rewrite(root, relPath string, vp versionLine, from, to string, changed *[]string) error {
	if vp.re == nil {
		return fmt.Errorf("no pattern for %s files", KindDockerfile)
	}
	full := filepath.Join(root, relPath)
	before := readOrEmpty(full)
	if before == "" {
		return nil
	}
	lines := strings.Split(before, "\n")
	touched := false
	for i, line := range lines {
		m := vp.re.FindStringSubmatchIndex(line)
		if m == nil {
			continue
		}
		start, end := m[2*vp.group], m[2*vp.group+1]
		if line[start:end] != from {
			continue
		}
		lines[i] = line[:start] + to + line[end:]
		touched = true
	}
	if !touched {
		return nil
	}
	return writeIfNeeded(full, relPath, before, strings.Join(lines, "\n"), changed)
}

// rewriteReadme rewrites the version inside the README rows of a program, and only
// there, so a version written in prose elsewhere in the file is left alone.
func rewriteReadme(root, program, from, to string, changed *[]string) error {
	full := filepath.Join(root, readmePath)
	before := readOrEmpty(full)
	if before == "" {
		return nil
	}
	lines := strings.Split(before, "\n")
	touched := false
	for i, line := range lines {
		m := readmeRow.FindStringSubmatch(line)
		if m == nil || m[1] != program || !strings.Contains(line, from) {
			continue
		}
		lines[i] = strings.ReplaceAll(line, from, to)
		touched = true
	}
	if !touched {
		return nil
	}
	return writeIfNeeded(full, readmePath, before, strings.Join(lines, "\n"), changed)
}

// workflowFile returns the workflow of a program, by file name when it follows the
// convention and otherwise by the workflow declaring that program name: the file of
// gomodifytags is spelled gomoditytags.yml, a typo kept for history. It returns the
// repository relative path and the content, or empty strings when there is none, as
// for gitstatus.
func workflowFile(root, program string) (string, string) {
	preferred := filepath.ToSlash(filepath.Join(workflowDir, program+".yml"))
	if content := readOrEmpty(filepath.Join(root, workflowDir, program+".yml")); content != "" {
		return preferred, content
	}
	names, err := filepath.Glob(filepath.Join(root, workflowDir, "*.yml"))
	if err != nil {
		return "", ""
	}
	sort.Strings(names)
	declared := regexp.MustCompile(`(?m)^name:[ \t]*` + regexp.QuoteMeta(program) + `[ \t]*$`)
	for _, path := range names {
		base := filepath.Base(path)
		if base == workflowReusable || base == workflowHubDesc {
			continue
		}
		if content := readOrEmpty(path); declared.MatchString(content) {
			return workflowDir + "/" + base, content
		}
	}
	return "", ""
}

func writeIfNeeded(full, relPath, before, after string, changed *[]string) error {
	if before == after {
		return nil
	}
	info, err := os.Stat(full)
	if err != nil {
		return fmt.Errorf("stat %s: %w", relPath, err)
	}
	if err := os.WriteFile(full, []byte(after), info.Mode().Perm()); err != nil {
		return fmt.Errorf("writing %s: %w", relPath, err)
	}
	*changed = append(*changed, filepath.ToSlash(relPath))
	return nil
}

func readFile(root, relPath string) string {
	return readOrEmpty(filepath.Join(root, relPath))
}

func readOrEmpty(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return string(data)
}
