// Package upgrade drives the version check of the programs of this repository, and
// optionally applies the upgrade to the pinned version strings, with one commit per
// program.
package upgrade

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strings"

	"github.com/qdm12/binpot/internal/gh"
	"github.com/qdm12/binpot/internal/pins"
	"github.com/qdm12/binpot/internal/program"
	"github.com/qdm12/binpot/internal/semver"
)

// Status is the outcome of checking one program.
type Status string

const (
	// StatusUpToDate means the pinned version is the newest stable one upstream.
	StatusUpToDate Status = "up to date"
	// StatusOutOfDate means a newer stable version is available upstream.
	StatusOutOfDate Status = "out of date"
	// StatusSkipped means the program could not be checked, a program pinned to a
	// commit for instance, and needs a human decision.
	StatusSkipped Status = "skipped"
	// StatusError means the check or the upgrade failed.
	StatusError Status = "error"
	// StatusUpgraded means the pins were rewritten.
	StatusUpgraded Status = "upgraded"
)

// Result is the outcome for a single program.
type Result struct {
	Program string `json:"program"`
	Repo    string `json:"repo"`
	// Pinned is the version the Dockerfile builds.
	Pinned string `json:"pinned"`
	// Latest is the newest stable version found upstream, empty when unknown.
	Latest string `json:"latest"`
	Status Status `json:"status"`
	// Detail explains a skip, an error, or drift.
	Detail string `json:"detail,omitempty"`
	// Changed lists the repository relative paths rewritten, when upgraded.
	Changed []string `json:"changed,omitempty"`
	// Commit is the short SHA of the commit created, when committed.
	Commit string `json:"commit,omitempty"`
}

// Options configures Run.
type Options struct {
	// Root is the repository root holding dockerfiles/ and .github/.
	Root string
	// Programs is the set to check, usually every program of the registry.
	Programs []program.Program
	// Client talks to GitHub.
	Client *gh.Client
	// Pages bounds how many pages of tags and releases are read per repository.
	Pages int
	// Update rewrites the version strings of the out of date programs.
	Update bool
	// Commit creates one commit per upgraded program, as Update must be set too.
	Commit bool
	// AllowDrift applies the update even when the Dockerfile, the workflow and the
	// README disagree about the current version.
	AllowDrift bool
	// DryRun reports what an update would change, without touching a file.
	DryRun bool
	// Log receives one line per program, and may be io.Discard.
	Log io.Writer
}

// Run checks every program, upgrades the out of date ones when Options.Update is
// set, and returns one Result per program in the order they were checked. An error
// is only returned for a failure of the run itself; a per program failure shows up
// as a Result with StatusError.
func Run(ctx context.Context, opts Options) ([]Result, error) {
	if opts.Root == "" {
		return nil, errors.New("no repository root")
	}
	if len(opts.Programs) == 0 {
		return nil, errors.New("no program to check")
	}
	if opts.Client == nil {
		opts.Client = gh.New("")
	}
	if opts.Pages < 1 {
		opts.Pages = 2
	}
	if opts.Log == nil {
		opts.Log = io.Discard
	}
	var results []Result
	for _, p := range opts.Programs {
		r := checkOne(ctx, opts, p)
		results = append(results, r)
		fmt.Fprintf(opts.Log, "%s\t%s\t%s\t%s\n", r.Program, r.Pinned, r.Latest, r.Status)
	}
	return results, nil
}

// checkOne handles a single program: read the pins, ask upstream, upgrade and commit
// when asked to.
func checkOne(ctx context.Context, opts Options, p program.Program) Result {
	r := Result{Program: p.Name, Repo: p.Repo}

	set, err := pins.Read(opts.Root, p.Name)
	if err != nil {
		r.Status, r.Detail = StatusError, err.Error()
		return r
	}
	pinned, err := set.Version()
	if err != nil {
		r.Status, r.Detail = StatusError, err.Error()
		return r
	}
	r.Pinned = pinned

	upstream, err := p.Latest(ctx, opts.Client, opts.Pages)
	if err != nil {
		if p.Track == program.TrackCommit || errors.Is(err, semver.ErrNotFound) {
			r.Status, r.Detail = StatusSkipped, noteOf(p, err)
			return r
		}
		r.Status, r.Detail = StatusError, err.Error()
		return r
	}
	r.Latest = upstream.Version

	if drift := set.Drift(); len(drift) > 0 && !opts.AllowDrift {
		r.Status, r.Detail = StatusSkipped, describeDrift(p.Name, drift)
		return r
	}

	parsedPinned, ok := semver.Parse(pinned)
	if !ok {
		r.Status, r.Detail = StatusSkipped,
			fmt.Sprintf("pinned value %q is not a version, upstream offers %s", pinned, upstream.Tag)
		return r
	}
	parsedLatest, ok := semver.Parse(upstream.Version)
	if !ok || !parsedPinned.Less(parsedLatest) {
		r.Status = StatusUpToDate
		return r
	}
	if !opts.Update {
		r.Status = StatusOutOfDate
		return r
	}

	target := matchStyle(upstream.Version, pinned)
	if opts.DryRun {
		r.Status, r.Detail = StatusOutOfDate, "would update to "+target
		return r
	}
	changed, err := pins.Update(opts.Root, p.Name, pinned, target)
	if err != nil {
		r.Status, r.Detail = StatusError, "updating: "+err.Error()
		return r
	}
	r.Changed = changed
	r.Status = StatusUpgraded
	if !opts.Commit {
		return r
	}
	sha, err := commit(opts.Root, p.Name, target, changed)
	if err != nil {
		r.Status, r.Detail = StatusError, "committed: "+err.Error()
		return r
	}
	r.Commit = sha
	return r
}

// commit creates the single commit of one program upgrade. The path form of git
// commit is used so nothing else that happens to be staged joins the commit.
func commit(root, name, version string, paths []string) (string, error) {
	if len(paths) == 0 {
		return "", errors.New("nothing staged to commit")
	}
	args := append([]string{"-C", root, "commit", "-q", "-m", commitMessage(name, version), "--"}, paths...)
	if out, err := exec.Command("git", args...).CombinedOutput(); err != nil {
		return "", fmt.Errorf("git commit: %w: %s", err, strings.TrimSpace(string(out)))
	}
	out, err := exec.Command("git", "-C", root, "rev-parse", "--short", "HEAD").CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("reading HEAD: %w: %s", err, strings.TrimSpace(string(out)))
	}
	return strings.TrimSpace(string(out)), nil
}

// commitMessage matches the history of the repository: `rumdl` v0.3.5.
func commitMessage(name, version string) string {
	return fmt.Sprintf("`%s` %s", name, version)
}

// matchStyle keeps the prefix style of the pinned value: a program pinned as v1.2.3
// stays pinned as v1.3.0 even when the upstream tag reads 1.3.0.
func matchStyle(latest, pinned string) string {
	if strings.HasPrefix(pinned, "v") && !strings.HasPrefix(latest, "v") {
		return "v" + latest
	}
	return latest
}

func describeDrift(name string, drift []pins.Source) string {
	parts := make([]string, 0, len(drift))
	for _, s := range drift {
		parts = append(parts, fmt.Sprintf("%s line %d has %s", s.Path, s.Line, s.Value))
	}
	return "pins disagree for " + name + ": " + strings.Join(parts, "; ") +
		" (fix them, or pass -allow-drift)"
}

func noteOf(p program.Program, err error) string {
	if p.Note != "" {
		return p.Note
	}
	if err != nil {
		return err.Error()
	}
	return "nothing to compare"
}
