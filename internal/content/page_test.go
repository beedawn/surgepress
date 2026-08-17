package content

import (
	"testing"
)

func TestPageStruct(t *testing.T) {
	// Test creating a Page
	page := Page{
		SourcePath: "src/index.md",
		OutputPath: "index.html",
	}

	// Test all fields
	if page.SourcePath != "src/index.md" {
		t.Errorf("Expected SourcePath 'src/index.md', got '%s'", page.SourcePath)
	}

	if page.OutputPath != "index.html" {
		t.Errorf("Expected OutputPath 'index.html', got '%s'", page.OutputPath)
	}

}
