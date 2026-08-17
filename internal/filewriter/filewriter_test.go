package filewriter

import (
	"github.com/beedawn/surgepress/internal/content"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestWriteFiles_Basic(t *testing.T) {
	// Create temporary directory for testing
	tmpDir, err := ioutil.TempDir("", "writefiles_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test pages
	pages := []content.RenderedPage{
		{
			OutputPath: filepath.Join(tmpDir, "test1.html"),
			HTML:       []byte("<h1>Test 1</h1>"),
		},
		{
			OutputPath: filepath.Join(tmpDir, "test2.html"),
			HTML:       []byte("<h1>Test 2</h1>"),
		},
	}
	for _, page := range pages {
		// Call the function
		WriteFile(&page)
	}
	// Verify files were created
	for _, page := range pages {
		_, err := os.Stat(page.OutputPath)
		if err != nil {
			t.Errorf("File %s was not created: %v", page.OutputPath, err)
		}

		// Verify content
		content, err := ioutil.ReadFile(page.OutputPath)
		if err != nil {
			t.Errorf("Failed to read %s: %v", page.OutputPath, err)
		}
		if string(content) != string(page.HTML) {
			t.Errorf("Content mismatch for %s", page.OutputPath)
		}
	}
}

func TestWriteFiles_NestedDirectories(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "writefiles_nested_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	pages := []content.RenderedPage{
		{
			OutputPath: filepath.Join(tmpDir, "nested", "deep", "test.html"),
			HTML:       []byte("<h1>Nested Test</h1>"),
		},
	}
	for _, page := range pages {
		WriteFile(&page)
	}
	// Verify nested directory structure was created
	_, err = os.Stat(filepath.Join(tmpDir, "nested", "deep", "test.html"))
	if err != nil {
		t.Errorf("Nested file was not created: %v", err)
	}
}

func TestWriteFiles_PermissionError(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "writefiles_perm_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a directory and make it read-only
	testDir := filepath.Join(tmpDir, "readonly")
	err = os.MkdirAll(testDir, 0755)
	if err != nil {
		t.Fatal(err)
	}

	// Make directory read-only
	err = os.Chmod(testDir, 0444)
	if err != nil {
		t.Fatal(err)
	}

	// This should cause an error when trying to write
	pages := []content.RenderedPage{
		{
			OutputPath: filepath.Join(testDir, "test.html"),
			HTML:       []byte("<h1>Test</h1>"),
		},
	}
	for _, page := range pages {
		err = WriteFile(&page)
		if err == nil {
			t.Error("Expected error for file permission denied")
		}

	}

	//should have error

	// Restore permissions
	os.Chmod(testDir, 0755)
}

func TestWriteFiles_PanicOnCreateDirError(t *testing.T) {
	// Test with a path that will cause CreateDir to fail
	// For example, a path with invalid characters or no write permissions

	pages := []content.RenderedPage{
		{
			OutputPath: "/invalid/path/with/forbidden:characters.html", // This might cause issues
			HTML:       []byte("<h1>Test</h1>"),
		},
	}

	// Test that it panics
	for _, page := range pages {
		err := WriteFile(&page)
		if err == nil {
			t.Error("Expected error for invalid path")
		}
	}
}
