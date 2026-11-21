// Package program holds the list of programs built by this repository, the upstream
// repository each one comes from, and how to find its newest published version.
package program

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/qdm12/binpot/internal/gh"
	"github.com/qdm12/binpot/internal/semver"
)

// Track tells how a program is pinned upstream: by release tag, or by commit.
type Track string

const (
	// TrackRelease pins a released version, the case of every program but one.
	TrackRelease Track = "release"
	// TrackCommit pins a commit hash, for upstreams that publish no release.
	TrackCommit Track = "commit"
)

// Program is a program built by this repository.
type Program struct {
	// Name is the program name, the directory name under dockerfiles/ and the
	// Docker image tag prefix.
	Name string
	// Repo is the upstream repository, in owner/name form.
	Repo string
	// TagPrefix is the prefix carried by the tags of a monorepo, "gopls/" for the
	// gopls tags of golang/tools. Empty when the repository only holds one program.
	TagPrefix string
	// Track is TrackRelease or TrackCommit.
	Track Track
	// Note documents a quirk worth showing with the program, such as a substitute
	// binary on one architecture.
	Note string
}

// programs is the registry, kept sorted by name. Adding a program to the
// repository means adding a line here, so the version checker covers it.
var programs = []Program{
	{Name: "bit", Repo: "chriswalz/bit", Track: TrackRelease},
	{Name: "buildx", Repo: "docker/buildx", Track: TrackRelease},
	{Name: "compose", Repo: "docker/compose", Track: TrackRelease},
	{Name: "dlv", Repo: "go-delve/delve", Track: TrackRelease, Note: "amd64 and arm64 only"},
	{Name: "docker", Repo: "docker/cli", Track: TrackRelease},
	{Name: "gh", Repo: "cli/cli", Track: TrackRelease},
	{Name: "gitstatus", Repo: "romkatv/gitstatus", Track: TrackRelease, Note: "no CI workflow"},
	{Name: "go-outline", Repo: "ramya-rao-a/go-outline", Track: TrackCommit, Note: "pinned to a commit, upstream has no release"},
	{Name: "gofumpt", Repo: "mvdan/gofumpt", Track: TrackRelease},
	{Name: "golangci-lint", Repo: "golangci/golangci-lint", Track: TrackRelease},
	{Name: "gomock", Repo: "uber-go/mock", Track: TrackRelease},
	{Name: "gomodifytags", Repo: "fatih/gomodifytags", Track: TrackRelease},
	{Name: "gopkgs", Repo: "uudashr/gopkgs", Track: TrackRelease},
	{Name: "goplay", Repo: "haya14busa/goplay", Track: TrackRelease},
	{Name: "gopls", Repo: "golang/tools", TagPrefix: "gopls/", Track: TrackRelease},
	{Name: "gotests", Repo: "cweill/gotests", Track: TrackRelease},
	{Name: "helm", Repo: "helm/helm", Track: TrackRelease},
	{Name: "impl", Repo: "josharian/impl", Track: TrackRelease},
	{Name: "kubectl", Repo: "kubernetes/kubernetes", Track: TrackRelease},
	{Name: "kubectx", Repo: "ahmetb/kubectx", Track: TrackRelease},
	{Name: "kubens", Repo: "ahmetb/kubectx", Track: TrackRelease},
	{Name: "logo-ls", Repo: "Yash-Handa/logo-ls", Track: TrackRelease, Note: "two README rows, the second pins the version tag"},
	{Name: "mockery", Repo: "vektra/mockery", Track: TrackRelease},
	{Name: "mockgen", Repo: "uber-go/mock", Track: TrackRelease},
	{Name: "rumdl", Repo: "rvben/rumdl", Track: TrackRelease, Note: "Rust, substitute program on s390x"},
	{Name: "stern", Repo: "stern/stern", Track: TrackRelease},
	{Name: "supervisord", Repo: "ochinchina/supervisord", Track: TrackRelease},
}

// All returns every program of the registry, sorted by name.
func All() []Program {
	out := make([]Program, len(programs))
	copy(out, programs)
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}

// Select returns the programs named by names, or every program when names is empty.
// It fails on an unknown name rather than silently skipping a program.
func Select(names []string) ([]Program, error) {
	if len(names) == 0 {
		return All(), nil
	}
	byName := make(map[string]Program, len(programs))
	for _, p := range programs {
		byName[p.Name] = p
	}
	var out []Program
	for _, name := range names {
		p, ok := byName[name]
		if !ok {
			return nil, fmt.Errorf("unknown program %q, run the checker with no name to list them", name)
		}
		out = append(out, p)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out, nil
}

// Upstream is a version published upstream.
type Upstream struct {
	// Tag is the tag name as published, prefix included, such as "gopls/v0.23.0".
	Tag string
	// Version is the tag without the monorepo prefix, the value written in the
	// Dockerfile and the workflow of the program, such as "v0.23.0".
	Version string
	// Source is "release" when the version comes from a published release, and
	// "tag" when the repository has no release and only tags to go by.
	Source string `json:"source"`
}

// Source of an Upstream, as reported by Latest.
const (
	SourceRelease = "release"
	SourceTag     = "tag"
)

// Latest returns the newest stable version published upstream, skipping pre-releases
// such as v1.2.0-rc.1.
//
// Published releases are the reference, and tags are the fallback for the
// repositories that publish none. Reading tags first would be wrong: a stray tag
// outranks the real versions, ochinchina/supervisord carries a 0.61 tag next to the
// v0.7.6 release, and kubernetes/kubernetes keeps patching old minor versions, so
// the newest release of a repository is not always the highest version.
func (p Program) Latest(ctx context.Context, client *gh.Client, pages int) (Upstream, error) {
	if p.Track == TrackCommit {
		return Upstream{}, fmt.Errorf("%s is pinned to a commit, there is no release to track", p.Name)
	}
	releases, releaseErr := client.Releases(ctx, p.Repo, pages)
	if releaseErr == nil {
		var tags []string
		for _, r := range releases {
			if r.Draft || r.Prerelease {
				continue
			}
			tags = append(tags, r.TagName)
		}
		if best, ok := semver.Highest(p.stableVersions(tags)); ok {
			return Upstream{
				Tag:     p.TagPrefix + best.String(),
				Version: best.String(),
				Source:  SourceRelease,
			}, nil
		}
	}
	tagNames, tagErr := client.Tags(ctx, p.Repo, pages)
	if tagErr != nil {
		if releaseErr != nil {
			return Upstream{}, fmt.Errorf("listing releases of %s: %w", p.Repo, releaseErr)
		}
		return Upstream{}, fmt.Errorf("listing tags of %s: %w", p.Repo, tagErr)
	}
	best, ok := semver.Highest(p.stableVersions(tagNames))
	if !ok {
		return Upstream{}, fmt.Errorf("%w: %d tags and %d releases looked at for %s",
			semver.ErrNotFound, len(tagNames), len(releases), p.Repo)
	}
	return Upstream{
		Tag:     p.TagPrefix + best.String(),
		Version: best.String(),
		Source:  SourceTag,
	}, nil
}

// stableVersions keeps the candidates that carry the tag prefix of the program and
// parse as a stable version, with the prefix stripped from the returned versions.
func (p Program) stableVersions(candidates []string) []semver.Version {
	out := make([]semver.Version, 0, len(candidates))
	for _, c := range candidates {
		if c == "" {
			continue
		}
		if p.TagPrefix != "" {
			if !strings.HasPrefix(c, p.TagPrefix) {
				continue
			}
			c = strings.TrimPrefix(c, p.TagPrefix)
		} else if strings.Contains(c, "/") {
			// A tag of another program of the same repository, such as the
			// protobuf tags of golang/tools.
			continue
		}
		if v, ok := semver.Parse(c); ok && v.Stable() {
			out = append(out, v)
		}
	}
	return out
}
