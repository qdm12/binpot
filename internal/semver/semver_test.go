package semver

import "testing"

func TestParse(t *testing.T) {
	for _, tc := range []struct {
		in     string
		ok     bool
		stable bool
		major  int
		minor  int
		patch  int
	}{
		{in: "v1.2.3", ok: true, stable: true, major: 1, minor: 2, patch: 3},
		{in: "1.2.3", ok: true, stable: true, major: 1, minor: 2, patch: 3},
		{in: "v0.7.6", ok: true, stable: true, major: 0, minor: 7, patch: 6},
		{in: "0.61", ok: true, stable: true, major: 0, minor: 61, patch: 0},
		{in: "v1.37.0-rc.1", ok: true, stable: false, major: 1, minor: 37},
		{in: "v0.0.3-test", ok: true, stable: false, patch: 3},
		{in: "v1.2.3+insecure", ok: true, stable: false, major: 1, minor: 2, patch: 3},
		{in: "9736a4b", ok: false},
		{in: "gopls/v0.23.0", ok: false},
		{in: "", ok: false},
	} {
		v, ok := Parse(tc.in)
		if ok != tc.ok {
			t.Errorf("Parse(%q) ok = %v, want %v", tc.in, ok, tc.ok)
			continue
		}
		if !ok {
			continue
		}
		if got := v.Stable(); got != tc.stable {
			t.Errorf("Parse(%q).Stable() = %v, want %v", tc.in, got, tc.stable)
		}
		if v.Major != tc.major || v.Minor != tc.minor || v.Patch != tc.patch {
			t.Errorf("Parse(%q) = %d.%d.%d, want %d.%d.%d",
				tc.in, v.Major, v.Minor, v.Patch, tc.major, tc.minor, tc.patch)
		}
		if v.String() != tc.in {
			t.Errorf("Parse(%q).String() = %q, want the raw input back", tc.in, v.String())
		}
	}
}

func TestLess(t *testing.T) {
	for _, tc := range []struct {
		low, high string
	}{
		{"v0.7.6", "v0.7.7"},
		{"v0.7.6", "v0.8.0"},
		{"v2.9.0", "v10.0.0"},
		{"v1.9.0", "v1.10.0"},
		{"v1.2.3", "v1.3.0"},
	} {
		low, _ := Parse(tc.low)
		high, _ := Parse(tc.high)
		if !low.Less(high) {
			t.Errorf("%s.Less(%s) = false, want true", tc.low, tc.high)
		}
		if high.Less(low) {
			t.Errorf("%s.Less(%s) = true, want false", tc.high, tc.low)
		}
	}
}

func TestHighestSkipsPrereleases(t *testing.T) {
	values := Scan([]string{
		"v1.38.0-alpha.0", "v1.37.2", "v1.36.6", "v1.35.10", "not-a-version",
	})
	best, ok := Highest(values)
	if !ok {
		t.Fatal("Highest found nothing, want v1.37.2")
	}
	if best.String() != "v1.37.2" {
		t.Errorf("Highest = %s, want v1.37.2", best)
	}
	if _, ok := Highest(Scan([]string{"v1.0.0-rc.1"})); ok {
		t.Error("Highest returned a pre-release")
	}
}
