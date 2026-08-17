package content

import (
	"testing"
)

func TestRenderedPageStruct(t *testing.T) {
	// Test creating a RenderedPage
	page := RenderedPage{
		OutputPath: "index.html",
		HTML:       []byte("<h1>Hello World</h1>"),
	}

	// Test all fields
	if page.OutputPath != "index.html" {
		t.Errorf("Expected OutputPath 'index.html', got '%s'", page.OutputPath)
	}

	expectedHTML := []byte("<h1>Hello World</h1>")
	if string(page.HTML) != string(expectedHTML) {
		t.Errorf("Expected HTML '<h1>Hello World</h1>', got '%s'", string(page.HTML))
	}
}
