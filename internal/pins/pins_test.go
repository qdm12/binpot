package pins

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// tree writes a miniature repository and returns its root. The gopls quirk of a
// Dockerfile written with CRLF line endings, and the workflow file whose name is a
// typo, are both reproduced: they are the two shapes the code has to cope with.
func tree(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	files := map[string]string{
		filepath.Join("dockerfiles", "foo", "Dockerfile"): strings.Join([]string{
			"ARG ALPINE_VERSION=3.23",
			"",
			"ARG VERSION=v1.2.3",
			"RUN git clone --depth 1 --branch ${VERSION} https://github.com/octocat/foo.git .",
			"",
		}, "\n"),
		// A CRLF file, with a trailing space after the version.
		filepath.Join("dockerfiles", "bar", "Dockerfile"): "ARG VERSION=v9.9.9 \r\nRUN echo bar\r\n",
		// The file name does not match the program name it declares.
		filepath.Join(".github", "workflows", "fooo.yml"): strings.Join([]string{
			"name: foo",
			"on:",
			"  push:",
			"jobs:",
			"  foo:",
			"    uses: ./.github/workflows/reusable.yml",
			"    with:",
			"      name: foo",
			"      version: v1.2.3",
			"",
		}, "\n"),
		// No workflow at all for bar, as for gitstatus.
		"README.md": strings.Join([]string{
			"| Program | Last version | Image tags | Architectures |",
			"| --- | --- | --- | --- |",
			"| [`foo`](https://github.com/octocat/foo) | [`v1.2.3`](https://github.com/octocat/foo/releases/tag/v1.2.3) | [**Docker Hub**](https://hub.docker.com/r/qmcgaw/binpot/tags?name=foo) | all |",
			"| [`bar`](https://github.com/octocat/bar) | [`v9.9.9`](https://github.com/octocat/bar/releases/tag/v9.9.9) | [**Docker Hub**](https://hub.docker.com/r/qmcgaw/binpot/tags?name=bar-v9.9.9) | all |",
			"",
			"A note about v1.2.3 in prose, which must never be rewritten.",
			"",
		}, "\n"),
	}
	for path, content := range files {
		full := filepath.Join(root, path)
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

func TestRead(t *testing.T) {
	root := tree(t)

	set, err := Read(root, "foo")
	if err != nil {
		t.Fatal(err)
	}
	if got, err := set.Version(); err != nil || got != "v1.2.3" {
		t.Errorf("Version() = %q, %v, want v1.2.3, nil", got, err)
	}
	if len(set.Sources) != 3 {
		t.Fatalf("read %d sources, want 3: %+v", len(set.Sources), set.Sources)
	}
	if drift := set.Drift(); len(drift) != 0 {
		t.Errorf("Drift() = %+v, want none", drift)
	}

	// The carriage return and the trailing space of a CRLF Dockerfile are not part
	// of the version.
	bar, err := Read(root, "bar")
	if err != nil {
		t.Fatal(err)
	}
	if got, err := bar.Version(); err != nil || got != "v9.9.9" {
		t.Errorf("bar Version() = %q, %v, want v9.9.9, nil", got, err)
	}
	if len(bar.Sources) != 2 {
		t.Errorf("bar read %d sources, want 2 (Dockerfile and README): %+v",
			len(bar.Sources), bar.Sources)
	}

	if _, err := Read(root, "absent"); err == nil {
		t.Error("Read(absent) succeeded, want an error")
	}
}

func TestUpdate(t *testing.T) {
	root := tree(t)

	changed, err := Update(root, "foo", "v1.2.3", "v1.3.0")
	if err != nil {
		t.Fatal(err)
	}
	want := []string{
		filepath.Join(".github", "workflows", "fooo.yml"),
		"README.md",
		filepath.Join("dockerfiles", "foo", "Dockerfile"),
	}
	if strings.Join(changed, ",") != strings.Join(want, ",") {
		t.Errorf("changed = %v, want %v", changed, want)
	}

	dockerfile := read(t, root, filepath.Join("dockerfiles", "foo", "Dockerfile"))
	if !strings.Contains(dockerfile, "ARG VERSION=v1.3.0\n") {
		t.Errorf("Dockerfile not rewritten:\n%s", dockerfile)
	}
	workflow := read(t, root, filepath.Join(".github", "workflows", "fooo.yml"))
	if !strings.Contains(workflow, "      version: v1.3.0\n") {
		t.Errorf("workflow not rewritten:\n%s", workflow)
	}
	readme := read(t, root, "README.md")
	if !strings.Contains(readme, "[`v1.3.0`](https://github.com/octocat/foo/releases/tag/v1.3.0)") {
		t.Errorf("README row not rewritten:\n%s", readme)
	}
	if !strings.Contains(readme, "A note about v1.2.3 in prose") {
		t.Error("prose outside the row was rewritten")
	}
	if strings.Contains(readme, "v1.2.3](https://github.com/octocat/foo") {
		t.Error("README row still holds the old version")
	}

	// Nothing is left to upgrade, and the second program is untouched.
	if _, err := Update(root, "foo", "v1.2.3", "v1.4.0"); err == nil {
		t.Error("Update from a version no longer pinned succeeded, want an error")
	}
	if got := read(t, root, filepath.Join("dockerfiles", "bar", "Dockerfile")); !strings.Contains(got, "v9.9.9") {
		t.Errorf("bar was touched: %q", got)
	}
}

// TestUpdateCRFLikeGopls checks that a CRLF Dockerfile keeps its line endings, and
// that a row pinning the version in the image tag follows, as logo-ls does.
func TestUpdateCRFLikeGopls(t *testing.T) {
	root := tree(t)

	if _, err := Update(root, "bar", "v9.9.9", "v10.0.0"); err != nil {
		t.Fatal(err)
	}
	got := read(t, root, filepath.Join("dockerfiles", "bar", "Dockerfile"))
	if got != "ARG VERSION=v10.0.0 \r\nRUN echo bar\r\n" {
		t.Errorf("CRLF Dockerfile = %q, want the trailing space and carriage returns kept", got)
	}
	readme := read(t, root, "README.md")
	if !strings.Contains(readme, "?name=bar-v10.0.0)") {
		t.Errorf("pinned image tag of the row not rewritten:\n%s", readme)
	}
	if strings.Contains(readme, "bar-v9.9.9") {
		t.Error("old pinned image tag survives")
	}
}

func TestUpdateRefusesNoOp(t *testing.T) {
	root := tree(t)
	for _, tc := range []struct{ from, to string }{
		{"", "v1.0.0"}, {"v1.2.3", ""}, {"v1.2.3", "v1.2.3"},
	} {
		if _, err := Update(root, "foo", tc.from, tc.to); err == nil {
			t.Errorf("Update(%q, %q) succeeded, want an error", tc.from, tc.to)
		}
	}
}

func read(t *testing.T, root, relPath string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(root, relPath))
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}
