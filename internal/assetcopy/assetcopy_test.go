package assetcopy

import (
	"os"
	"path/filepath"
	"testing"
)


func writeTestFile(t *testing.T, path, contents string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(contents), 0644); err != nil {
		t.Fatal(err)
	}
}

func readTestFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("expected %s to exist: %v", path, err)
	}
	return string(data)
}


func TestCopyDir_CopiesNestedFiles(t *testing.T) {
	src := t.TempDir()
	dest := filepath.Join(t.TempDir(), "out", "css")

	writeTestFile(t, filepath.Join(src, "post.css"), "h1 { color: red; }")
	writeTestFile(t, filepath.Join(src, "themes", "dark.css"), "body { background: #000; }")

	if err := CopyDir(src, dest); err != nil {
		t.Fatal(err)
	}

	if got := readTestFile(t, filepath.Join(dest, "post.css")); got != "h1 { color: red; }" {
		t.Errorf("post.css = %q", got)
	}
	if got := readTestFile(t, filepath.Join(dest, "themes", "dark.css")); got != "body { background: #000; }" {
		t.Errorf("themes/dark.css = %q", got)
	}
}

func TestCOpyDir_OverwritesExistingFiles(t *testing.T) {
	src := t.TempDir()
	dest := t.TempDir()

	writeTestFile(t, filepath.Join(dest, "post.css"), "old")
	writeTestFile(t, filepath.Join(src, "post.css"), "new")

	if err := CopyDir(src, dest); err != nil {
		t.Fatal(err)
	}
	if got := readTestFile(t, filepath.Join(dest, "post.css")); got != "new" {
		t.Errorf("expected file to be overwritten, got %q", got)
	}
}

func TestCopyDir_MissingSourceIsNotAnError(t *testing.T) {
	dest := t.TempDir()
	if err := CopyDir(filepath.Join(t.TempDir(), "does-not-exist"), dest); err != nil {
		t.Errorf("expected nil error for missing source, got %v", err)
	}
}

func TestCopyDir_SourceIsFile(t *testing.T) {
	file := filepath.Join(t.TempDir(), "not a dir.css")

	writeTestFile(t, file, "")
	if err := CopyDir(file, t.TempDir()); err == nil {
		t.Error("expected an error when source is a file")
	}
}