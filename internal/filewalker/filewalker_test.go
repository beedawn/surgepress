package filewalker

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

var testData = filepath.Join("..", "..", "test_data")

func TestGetProjectPath_Valid(t *testing.T) {
	pages, err := WalkFiles(testData)
	if err != nil {
		t.Fatalf("WalkFiles failed: %v", err)
	}
	if len(pages) != 3 {
		t.Errorf("Expected 3 page, got %d", len(pages))
	}

	page := pages[0]

	type TestPage struct {
		SourcePath string
		OutputPath string
	}

	want := TestPage{
		SourcePath: filepath.Join(testData, "blank.md"),
		OutputPath: filepath.Join("out", "blank.html"),
	}

	if page.SourcePath != want.SourcePath {
		t.Errorf("SourcePath mismatch")
	}

	if page.OutputPath != want.OutputPath {
		t.Errorf("OutputPath mismatch")
	}

}

func TestGetProjectPath_Blank(t *testing.T) {
	pages, err := WalkFiles("")
	if err == nil {
		t.Fatalf("WalkFiles with blank path should have failed, but got pages %d", len(pages))
	}
	if len(pages) != 0 {
		t.Errorf("Expected 0 pages when path is blank, got %d", len(pages))
	}

	if err.Error() != "path is blank" {
		t.Errorf("Error mismatch, expected error 'path is blank', got %v", err)
	}

}

func TestGetProjectPath_BadPath(t *testing.T) {
	pages, err := WalkFiles("/gibberish/132874891237")
	if err == nil {
		t.Fatalf("WalkFiles with nonexistent path should have failed")
	}
	if len(pages) != 0 {
		t.Errorf("Expected 0 pages, got %d", len(pages))
	}
	if !os.IsNotExist(err) {
		t.Errorf("Expected file-not-found error, got %v", err)
	}
}

func TestGetProjectPath_EmptyDir(t *testing.T) {
	emptyDir := t.TempDir()
	pages, err := WalkFiles(emptyDir)
	if err != nil {
		t.Fatalf("Empty directory should not fail: %v", err)
	}
	if len(pages) != 0 {
		t.Errorf("Expected 0 pages when directory is empty, got %d", len(pages))
	}
}

func TestGetProjectPath_FileIsNotDirectory(t *testing.T) {
	testFile := filepath.Join(testData, "test.md")
	_, err := WalkFiles(testFile)
	if err == nil {
		t.Error("Expected error when passing file path to WalkFiles")
	}
}

func TestWalkFiles_PermissionError2(t *testing.T) {
	// Create a temporary directory
	tmpDir, err := ioutil.TempDir("", "walkfiles_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a test file
	testFile := filepath.Join(tmpDir, "test.md")
	err = ioutil.WriteFile(testFile, []byte("# Test\n\nContent"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Remove read permissions from the directory to force an error
	err = os.Chmod(tmpDir, 0200) // Write-only permissions
	if err != nil {
		t.Fatal(err)
	}

	// This should trigger the error handling
	_, err = WalkFiles(tmpDir)
	if err == nil {
		t.Error("Expected error when directory has no read permissions")
	}

	// Restore permissions
	os.Chmod(tmpDir, 0755)
}
