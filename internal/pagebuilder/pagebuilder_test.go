package pagebuilder

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"surgepress/internal/content"
	"surgepress/internal/pathutil"
	"surgepress/internal/siteconfig"
	"testing"
)

var testDir = filepath.Join("..", "..", "internal", "template")
var testDataDir = filepath.Join("..", "..", "test_data")
var configFilePath = filepath.Join(testDataDir, "configdata.json")

//check config file path is valid filetype

func genConfigData() siteconfig.MetaData {
	if !pathutil.ValidateFile(configFilePath) {
		panic("Invalid config file path")
	}
	configFile, _ := os.Open(configFilePath)
	defer configFile.Close()

	// 2. Initialize the target variable
	var configData siteconfig.MetaData
	//indexMeta := siteconfig.IndexMeta{
	//	Title:       "My Blog",
	//	Description: "My brand new exciting blog",
	//	Link:        "https://blog.com",
	//	Language:    "en-US",
	//}
	// 3. Create a decoder and decode the stream
	decoder := json.NewDecoder(configFile)
	if err := decoder.Decode(&configData); err != nil {
		log.Fatalf("failed to decode JSON: %v", err)
	}
	return configData
}

func TestPageBuilder_Valid(T *testing.T) {

	testPages := []content.Page{}
	testPage := content.Page{
		SourcePath: filepath.Join(testDataDir, "test.md"),
		OutputPath: "out/test.html",
	}
	testPages = append(testPages, testPage)
	configData := genConfigData()
	err := BuildPages(testPages, testDir, configData)
	//need to use these variabels and check output
	if err != nil {
		T.Error(err)
	}
	//	if len(renderedPages) != len(testPages) {
	//		T.Errorf("Expected %d pages, got %d", len(testPages), len(renderedPages))
	//	}
	//
	//	if renderedPages[0].OutputPath != "out/test.html" {
	//		T.Errorf("Output path should be out/test.html, got %v", renderedPages[0].OutputPath)
	//	}
	//	if string(renderedPages[0].HTML) != `<!DOCTYPE html>
	//<html lang="en">
	//<head>
	//<meta charset="utf-8">
	//<title>
	//test title</title>
	//		</head>
	//			<body>hello world</body>
	//</html>` {
	//		T.Errorf("HTML output does not match expected")
	//
	//	}

}

func TestPageBuilder_EmptyContent(T *testing.T) {

	testPages := []content.Page{}
	testPage := content.Page{
		SourcePath: filepath.Join(testDataDir, "blank.md"),
		OutputPath: "out/test.html",
	}
	configData := genConfigData()
	testPages = append(testPages, testPage)
	err := BuildPages(testPages, testDir, configData)
	//need to use these variabels and check output
	if err != nil {
		T.Error(err)
	}
	//	if len(renderedPages) != len(testPages) {
	//		T.Errorf("Expected %d pages, got %d", len(testPages), len(renderedPages))
	//	}
	//
	//	if renderedPages[0].OutputPath != "out/test.html" {
	//		T.Errorf("Output path should be out/test.html, got %v", renderedPages[0].OutputPath)
	//	}
	//	if string(renderedPages[0].HTML) != `<!DOCTYPE html>
	//<html lang="en">
	//<head>
	//<meta charset="utf-8">
	//<title>
	//test title</title>
	//		</head>
	//			<body></body>
	//</html>` {
	//		T.Errorf("HTML output does not match expected")
	//
	//	}

}

func TestPageBuilder_TwoPagest(T *testing.T) {

	testPages := []content.Page{}
	testPage := content.Page{
		SourcePath: filepath.Join(testDataDir, "test.md"),
		OutputPath: "out/test.html",
	}
	testPage2 := content.Page{
		SourcePath: filepath.Join(testDataDir, "test.md"),
		OutputPath: "out/test.html",
	}
	configData := genConfigData()
	testPages = append(testPages, testPage)
	testPages = append(testPages, testPage2)
	err := BuildPages(testPages, testDir, configData)
	//need to use these variabels and check output
	if err != nil {
		T.Error(err)
	}
	//if len(renderedPages) != len(testPages) {
	//	T.Errorf("Expected %d pages, got %d", len(testPages), len(renderedPages))
	//}

}

func TestPageBuilder_InvalidPath(t *testing.T) {
	testPages := []content.Page{
		{
			SourcePath: "/nonexistent/path.md",
			OutputPath: "out/test.html",
		},
	}
	configData := genConfigData()
	// This should handle the error gracefully
	err := BuildPages(testPages, testDir, configData)
	if err != nil {
		// Depending on your implementation, this might or might not error
		t.Logf("Expected error occurred: %v", err)
	}

	// Even if it errors, make sure it doesn't panic
	//if len(renderedPages) >= 0 { // Just to make sure it doesn't crash
	//	t.Logf("Function completed without panic")
	//}
}

func TestBuildPagesWithNoTemplate(t *testing.T) {
	// Create a file that might cause issues during parsing
	testMD := `---
title: "Test Page"
date: "2023-01-01"
---

This is a valid markdown document.

[link](http://example.com)

` // Missing closing parenthesis - this might cause issues

	tmpFile := "test_malformed.md"
	err := os.WriteFile(tmpFile, []byte(testMD), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	defer os.Remove(tmpFile)

	configData := siteconfig.MetaData{
		Global_Stylesheets: []string{"/css/main.css"},
	}

	testPages := []content.Page{
		{
			SourcePath: tmpFile,
			OutputPath: "out/test.html",
		},
	}

	err = BuildPages(testPages, testDir, configData)
	if err != nil {
		t.Logf("Got error (might be from md.Convert): %v", err)
	} else {
		t.Log("No error occurred - this is fine for valid markdown")
	}
}
func TestIndexBuilder_Valid(t *testing.T) {
	testPages := []content.Page{
		{
			SourcePath: filepath.Join(testDataDir, "blank.md"),
			OutputPath: "out/test.html",
		},
	}

	configData := genConfigData()
	err := BuildIndex(testPages, testDir, configData)
	if err != nil {
		t.Error(err)
	}
}

func TestIndexBuilder_BadPath(t *testing.T) {
	testPages := []content.Page{
		{
			SourcePath: "blank.md",
			OutputPath: "out/test.html",
		},
	}

	configData := genConfigData()
	err := BuildIndex(testPages, testDir, configData)
	if err == nil {
		t.Error("Supposed to throw file error!")
	}
}

func TestIndexBuilder_EmptyMarkdown(t *testing.T) {

	testMD := `
` // Missing closing parenthesis - this might cause issues

	tmpFile := "test_malformed.md"
	err := os.WriteFile(tmpFile, []byte(testMD), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	defer os.Remove(tmpFile)
	testPages := []content.Page{
		{
			SourcePath: tmpFile,
			OutputPath: "out/test.html",
		},
	}

	configData := genConfigData()
	err = BuildIndex(testPages, testDir, configData)
	if err == nil {
		t.Error("Should throw error about empty markdown, missing title, date")
	}
}

func TestIndexBuilder_EmptyTitle(t *testing.T) {

	testMD := `---
Title: ""
Date: "2023-01-01"
---

` // Missing closing parenthesis - this might cause issues

	tmpFile := "test_malformed.md"
	err := os.WriteFile(tmpFile, []byte(testMD), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	defer os.Remove(tmpFile)
	testPages := []content.Page{
		{
			SourcePath: tmpFile,
			OutputPath: "out/test.html",
		},
	}

	configData := genConfigData()
	err = BuildIndex(testPages, testDir, configData)
	if err == nil {
		t.Error("Should throw error about empty title")
	}
}

func TestIndexBuilder_MissingDate(t *testing.T) {

	testMD := `---
Title: "Test Page"
---

` // Missing closing parenthesis - this might cause issues

	tmpFile := "test_malformed.md"
	err := os.WriteFile(tmpFile, []byte(testMD), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	defer os.Remove(tmpFile)
	testPages := []content.Page{
		{
			SourcePath: tmpFile,
			OutputPath: "out/test.html",
		},
	}

	configData := genConfigData()
	err = BuildIndex(testPages, testDir, configData)
	if err == nil {
		t.Error("Should throw error about empty markdown, missing date")
	}
}

func TestIndexBuilder_EmptyDate(t *testing.T) {

	testMD := `---
Title: "Test Page"
Date: ""
---

` // Missing closing parenthesis - this might cause issues

	tmpFile := "test_malformed.md"
	err := os.WriteFile(tmpFile, []byte(testMD), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	defer os.Remove(tmpFile)
	testPages := []content.Page{
		{
			SourcePath: tmpFile,
			OutputPath: "out/test.html",
		},
	}

	configData := genConfigData()
	err = BuildIndex(testPages, testDir, configData)
	if err == nil {
		t.Error("Should throw error about empty markdown, empty date")
	}
}

func TestBuildPages_MissingTemplateReturnsError(t *testing.T) {
	tmpDir := t.TempDir()

	sourcePath := filepath.Join(tmpDir, "page.md")
	markdown := `---
Title: "Test Page"
Template: "does-not-exist.html"
---

# Hello
`

	if err := os.WriteFile(sourcePath, []byte(markdown), 0644); err != nil {
		t.Fatalf("failed to create test markdown: %v", err)
	}

	pages := []content.Page{
		{
			SourcePath: sourcePath,
			OutputPath: filepath.Join(tmpDir, "out", "page.html"),
		},
	}

	err := BuildPages(pages, tmpDir, siteconfig.MetaData{})
	if err == nil {
		t.Fatal("expected error for missing template, got nil")
	}
}

func TestBuildPages_WriteFailureReturnsError(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a valid template.
	templatePath := filepath.Join(tmpDir, "default.html")
	templateData := `<!DOCTYPE html>
<html>
<body>{{.Content}}</body>
</html>`

	if err := os.WriteFile(templatePath, []byte(templateData), 0644); err != nil {
		t.Fatalf("failed to create template: %v", err)
	}

	// Create valid Markdown.
	sourcePath := filepath.Join(tmpDir, "page.md")
	if err := os.WriteFile(sourcePath, []byte("# Hello"), 0644); err != nil {
		t.Fatalf("failed to create markdown: %v", err)
	}

	// Put a regular file where BuildPages expects a directory.
	// filepath.Join(blocker, "page.html") therefore cannot be created.
	blocker := filepath.Join(tmpDir, "blocked")
	if err := os.WriteFile(blocker, []byte("not a directory"), 0644); err != nil {
		t.Fatalf("failed to create blocker file: %v", err)
	}

	pages := []content.Page{
		{
			SourcePath: sourcePath,
			OutputPath: filepath.Join(blocker, "page.html"),
		},
	}

	err := BuildPages(pages, tmpDir, siteconfig.MetaData{})
	if err == nil {
		t.Fatal("expected file write error, got nil")
	}
}
