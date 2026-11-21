// Command versions checks the programs built by this repository against the releases
// published upstream, and can bump their pinned version strings.
//
// Usage:
//
//	versions                 # check every program
//	versions stern helm      # check two programs
//	versions -update stern   # rewrite the pins of stern, no commit
//	versions -update -commit # bump every program out of date, one commit each
//
// The commit message of an upgrade is `programname` vX.Y.Z, matching the history of
// the repository. Nothing is ever pushed.
//
// A token raises the GitHub rate limit from 60 to 5000 requests an hour: export
// GH_TOKEN or GITHUB_TOKEN before a full run.
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/qdm12/binpot/internal/gh"
	"github.com/qdm12/binpot/internal/program"
	"github.com/qdm12/binpot/internal/upgrade"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "versions: error:", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	fs := flag.NewFlagSet("versions", flag.ContinueOnError)
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage: versions [-update] [-commit] [-json] [program...]")
		fs.PrintDefaults()
	}
	dir := fs.String("dir", "", "repository root, detected from the working directory by default")
	token := fs.String("token", "", "GitHub API token, read from $GH_TOKEN or $GITHUB_TOKEN by default")
	pages := fs.Int("pages", 2, "maximum pages of tags and releases read per repository")
	update := fs.Bool("update", false, "rewrite the version strings of the out of date programs")
	doCommit := fs.Bool("commit", false, "with -update: create one commit per upgraded program")
	allowDrift := fs.Bool("allow-drift", false,
		"update even when the Dockerfile, the workflow and the README disagree")
	dryRun := fs.Bool("dry-run", false, "report what -update would change, writing nothing")
	jsonOut := fs.Bool("json", false, "print the results as JSON instead of a table")
	quiet := fs.Bool("quiet", false, "do not echo the progress of each program")
	if err := fs.Parse(args); err != nil {
		return err
	}

	root, err := findRoot(*dir)
	if err != nil {
		return err
	}
	programs, err := program.Select(fs.Args())
	if err != nil {
		return err
	}
	client := gh.New(firstNonEmpty(*token, os.Getenv("GH_TOKEN"), os.Getenv("GITHUB_TOKEN")))

	var log io.Writer = os.Stdout
	if *quiet {
		log = io.Discard
	}
	results, err := upgrade.Run(context.Background(), upgrade.Options{
		Root:       root,
		Programs:   programs,
		Client:     client,
		Pages:      *pages,
		Update:     *update,
		Commit:     *doCommit,
		AllowDrift: *allowDrift,
		DryRun:     *dryRun,
		Log:        log,
	})
	if err != nil {
		return err
	}

	if *jsonOut {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(results)
	}
	printTable(results)
	changed := 0
	for _, r := range results {
		if r.Status == upgrade.StatusUpgraded {
			changed++
		}
	}
	if *update && !*dryRun {
		fmt.Printf("\n%d program(s) upgraded%s.\n", changed,
			map[bool]string{true: ", committed", false: ""}[*doCommit])
	}
	var failed []string
	for _, r := range results {
		if r.Status == upgrade.StatusError {
			failed = append(failed, r.Program+": "+r.Detail)
		}
	}
	if len(failed) > 0 {
		return fmt.Errorf("%d program(s) failed: %s", len(failed), strings.Join(failed, "; "))
	}
	return nil
}

// printTable writes an aligned report of the run.
func printTable(results []upgrade.Result) {
	rows := [][]string{{"PROGRAM", "PINNED", "LATEST", "STATUS", "DETAIL"}}
	for _, r := range results {
		rows = append(rows, []string{r.Program, r.Pinned, orDash(r.Latest), string(r.Status), r.Detail})
	}
	widths := make([]int, len(rows[0]))
	for _, row := range rows {
		for i, cell := range row {
			if n := len(cell); n > widths[i] {
				widths[i] = n
			}
		}
	}
	for _, row := range rows {
		fmt.Println(strings.Join(pad(row, widths), "  "))
	}
}

// pad pads every cell of a row to its column width, the last cell left as is.
func pad(row []string, widths []int) []string {
	out := make([]string, len(row))
	for i, cell := range row {
		out[i] = cell
		if i < len(row)-1 {
			out[i] += strings.Repeat(" ", widths[i]-len(cell))
		}
	}
	return out
}

// findRoot walks up from the working directory until it finds a directory holding
// both dockerfiles/ and .github/, so the tool works from any subdirectory.
func findRoot(dir string) (string, error) {
	if dir != "" {
		if ok, err := looksLikeRoot(dir); err != nil {
			return "", err
		} else if ok {
			return filepath.Abs(dir)
		}
		return "", fmt.Errorf("%s is not a binpot repository: no dockerfiles/ directory", dir)
	}
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("getting the working directory: %w", err)
	}
	for path := cwd; ; {
		ok, err := looksLikeRoot(path)
		if err != nil {
			return "", err
		}
		if ok {
			return path, nil
		}
		parent := filepath.Dir(path)
		if parent == path {
			return "", fmt.Errorf("no binpot repository found above %s, pass -dir", cwd)
		}
		path = parent
	}
}

func looksLikeRoot(dir string) (bool, error) {
	info, err := os.Stat(filepath.Join(dir, "dockerfiles"))
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("checking %s: %w", dir, err)
	}
	return info.IsDir(), nil
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}

func orDash(s string) string {
	if s == "" {
		return "-"
	}
	return s
}
