// Package semver parses and compares the subset of semantic versions used by the
// programs of this repository: optional "v" prefix, and two or three dot separated
// numbers, optionally followed by a pre-release or build suffix.
package semver

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Version is a parsed version. The zero value is invalid.
type Version struct {
	Major, Minor, Patch int
	// Suffix is anything following the numbers, such as "-rc.1" or "+insecure".
	// It is empty for a final release.
	Suffix string
	// Raw is the string the version was parsed from, prefix included.
	Raw string
}

var pattern = regexp.MustCompile(`^v?(\d+)\.(\d+)(?:\.(\d+))?([-.+].*)?$`)

// Parse reports whether s is a version, with the numeric core of two or three parts.
func Parse(s string) (Version, bool) {
	m := pattern.FindStringSubmatch(strings.TrimSpace(s))
	if m == nil {
		return Version{}, false
	}
	patch := m[3]
	if patch == "" {
		patch = "0"
	}
	major, err := strconv.Atoi(m[1])
	if err != nil {
		return Version{}, false
	}
	minor, err := strconv.Atoi(m[2])
	if err != nil {
		return Version{}, false
	}
	npatch, err := strconv.Atoi(patch)
	if err != nil {
		return Version{}, false
	}
	return Version{
		Major:  major,
		Minor:  minor,
		Patch:  npatch,
		Suffix: m[4],
		Raw:    strings.TrimSpace(s),
	}, true
}

// Stable reports whether the version is a final release: no pre-release or build
// suffix. Release candidates tagged "-rc.1" and the like are not stable.
func (v Version) Stable() bool {
	return v.Suffix == ""
}

// Less reports whether v sorts before other, numerically on each component.
func (v Version) Less(other Version) bool {
	if v.Major != other.Major {
		return v.Major < other.Major
	}
	if v.Minor != other.Minor {
		return v.Minor < other.Minor
	}
	return v.Patch < other.Patch
}

// String returns the raw string the version was parsed from.
func (v Version) String() string {
	return v.Raw
}

// Highest returns the greatest version of stable, and reports whether the slice
// held at least one. Versions that cannot be parsed or are not stable are ignored,
// which keeps the callers free of filtering.
func Highest(stable []Version) (Version, bool) {
	var best Version
	found := false
	for _, v := range stable {
		if !v.Stable() {
			continue
		}
		if !found || best.Less(v) {
			best, found = v, true
		}
	}
	return best, found
}

// Scan parses every string of values and returns the stable versions found.
func Scan(values []string) []Version {
	out := make([]Version, 0, len(values))
	for _, s := range values {
		if v, ok := Parse(s); ok && v.Stable() {
			out = append(out, v)
		}
	}
	return out
}

// ErrNotFound is returned when no stable version could be found upstream.
var ErrNotFound = fmt.Errorf("no stable version found upstream")
