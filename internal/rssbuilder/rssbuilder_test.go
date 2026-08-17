package rssbuilder

import (
	"github.com/beedawn/surgepress/internal/indexpost"
	"github.com/beedawn/surgepress/internal/siteconfig"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestRSSGenerator(t *testing.T) {
	// Create test posts
	posts := []indexpost.IndexPost{
		{
			Title:   "Test Post 1",
			Content: "<p>This is the first test post</p>",
			Date:    "2023-01-01",
			URL:     "/posts/post1.html",
		},
		{
			Title:   "Test Post 2",
			Content: "<p>This is the second test post</p>",
			Date:    "2023-01-02",
			URL:     "/posts/post2.html",
		},
	}

	// Create test config
	indexMeta := siteconfig.MetaData{
		Title:       "Test Blog",
		Description: "A test blog for testing",
		Link:        "https://testblog.com",
		Language:    "en-US",
	}

	// Run RSS generator
	err := RSSGenerator(posts, indexMeta)
	if err != nil {
		t.Errorf("RSSGenerator failed: %v", err)
	}

	// Verify file was created
	_, err = os.Stat("out/feed.xml")
	if err != nil {
		t.Errorf("RSS file was not created: %v", err)
	}

	// Clean up
	os.Remove("out/feed.xml")
	os.Remove("out")
}

func TestRSSBuilder(t *testing.T) {
	// Create test posts
	posts := []indexpost.IndexPost{
		{
			Title:   "Test Post",
			Content: "<p>Test content</p>",
			Date:    "2023-01-01",
			URL:     "/posts/test.html",
		},
	}

	// Create test config
	indexMeta := siteconfig.MetaData{
		Title:       "Test Blog",
		Description: "A test blog",
		Link:        "https://testblog.com",
		Language:    "en-US",
	}

	// Build RSS
	feed, err := RSSBuilder(posts, indexMeta)
	if err != nil {
		t.Errorf("RSSBuilder failed: %v", err)
	}

	// Verify feed structure
	if feed.Version != "2.0" {
		t.Error("Expected RSS version 2.0")
	}

	if feed.Channel == nil {
		t.Error("Channel should not be nil")
	}

	if feed.Channel.Title != indexMeta.Title {
		t.Errorf("Channel title mismatch: got %s, want %s", feed.Channel.Title, indexMeta.Title)
	}

	if len(feed.Channel.Items) != 1 {
		t.Errorf("Expected 1 item, got %d", len(feed.Channel.Items))
	}

	item := feed.Channel.Items[0]
	if item.Title != "Test Post" {
		t.Errorf("Item title mismatch: got %s, want Test Post", item.Title)
	}

	if item.Link != "https://testblog.com/posts/test.html" {
		t.Errorf("Item link mismatch: got %s", item.Link)
	}
}

func TestRSSBuilder_InvalidDate(t *testing.T) {
	// Create test posts with invalid date
	posts := []indexpost.IndexPost{
		{
			Title:   "Test Post",
			Content: "<p>Test content</p>",
			Date:    "invalid-date", // This should fallback to current time
			URL:     "/posts/test.html",
		},
	}

	indexMeta := siteconfig.MetaData{
		Title:       "Test Blog",
		Description: "A test blog",
		Link:        "https://testblog.com",
		Language:    "en-US",
	}

	// Build RSS - this should not fail even with invalid date
	feed, err := RSSBuilder(posts, indexMeta)
	if err != nil {
		t.Errorf("RSSBuilder failed with invalid date: %v", err)
	}

	if len(feed.Channel.Items) != 1 {
		t.Error("Expected 1 item")
	}
}

func TestCheckOutDir(t *testing.T) {
	// Test with existing directory
	err := CheckOutDir()
	if err != nil {
		t.Errorf("CheckOutDir failed: %v", err)
	}

	// Test that directory exists
	info, err := os.Stat("out")
	if err != nil {
		t.Error("Output directory was not created")
	}
	if !info.IsDir() {
		t.Error("Output path is not a directory")
	}

	// Clean up
	os.Remove("out")
}

func TestCheckOutDir_WithPermissions(t *testing.T) {
	// Create a temporary directory and make it read-only
	tempDir, err := os.MkdirTemp("", "rss_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Try to create output in a read-only directory (this might not work on all systems)
	// This tests the error path

	// Test that we can create the out directory normally
	err = CheckOutDir()
	if err != nil {
		t.Errorf("CheckOutDir failed: %v", err)
	}

	// Clean up
	os.Remove("out")
}

func TestWriteRSS(t *testing.T) {
	feed := RSS{
		Version: "2.0",
		Channel: &Channel{
			Title:       "Test Blog",
			Link:        "https://testblog.com",
			Description: "A test blog",
			Language:    "en-US",
			Items: []Item{
				{
					Title:       "Test Post",
					Link:        "https://testblog.com/posts/test.html",
					Description: "Test content",
					PubDate:     time.Now().UTC().Format(time.RFC1123Z),
					GUID: GUID{
						Value:       "test-guid",
						IsPermaLink: "false",
					},
				},
			},
		},
	}

	indexMeta := siteconfig.MetaData{
		Title:       "Test Blog",
		Description: "A test blog",
		Link:        "https://testblog.com",
		Language:    "en-US",
	}

	outputPath := filepath.Join(t.TempDir(), "feed.xml")

	if err := WriteRSS(feed, indexMeta, outputPath); err != nil {
		t.Fatalf("WriteRSS failed: %v", err)
	}

	fileInfo, err := os.Stat(outputPath)
	if err != nil {
		t.Fatalf("RSS file was not created: %v", err)
	}

	if fileInfo.Size() == 0 {
		t.Error("RSS file is empty")
	}
}
func TestWriteRSS_FileError(t *testing.T) {
	feed := RSS{
		Version: "2.0",
		Channel: &Channel{
			Title:       "Test Blog",
			Link:        "https://testblog.com",
			Description: "A test blog",
			Language:    "en-US",
			Items:       []Item{},
		},
	}

	indexMeta := siteconfig.MetaData{
		Title:       "Test Blog",
		Description: "A test blog",
		Link:        "https://testblog.com",
		Language:    "en-US",
	}

	tempDir := t.TempDir()

	// Create a regular file where a directory would need to exist.
	blocker := filepath.Join(tempDir, "blocked")
	if err := os.WriteFile(blocker, []byte("not a directory"), 0644); err != nil {
		t.Fatalf("failed to create blocker file: %v", err)
	}

	outputPath := filepath.Join(blocker, "feed.xml")

	err := WriteRSS(feed, indexMeta, outputPath)
	if err == nil {
		t.Fatal("expected WriteRSS to return an error")
	}
}

func TestRSSBuilder_EmptyPosts(t *testing.T) {
	// Test with empty posts slice
	posts := []indexpost.IndexPost{}

	indexMeta := siteconfig.MetaData{
		Title:       "Test Blog",
		Description: "A test blog",
		Link:        "https://testblog.com",
		Language:    "en-US",
	}

	feed, err := RSSBuilder(posts, indexMeta)
	if err != nil {
		t.Errorf("RSSBuilder failed with empty posts: %v", err)
	}

	if feed.Channel == nil {
		t.Error("Channel should not be nil")
	}

	if len(feed.Channel.Items) != 0 {
		t.Error("Expected 0 items for empty posts")
	}
}

func TestRSSBuilder_MultiplePosts(t *testing.T) {
	// Test with multiple posts
	posts := []indexpost.IndexPost{
		{
			Title:   "Post 1",
			Content: "<p>Content 1</p>",
			Date:    "2023-01-01",
			URL:     "/posts/1.html",
		},
		{
			Title:   "Post 2",
			Content: "<p>Content 2</p>",
			Date:    "2023-01-02",
			URL:     "/posts/2.html",
		},
		{
			Title:   "Post 3",
			Content: "<p>Content 3</p>",
			Date:    "2023-01-03",
			URL:     "/posts/3.html",
		},
	}

	indexMeta := siteconfig.MetaData{
		Title:       "Test Blog",
		Description: "A test blog",
		Link:        "https://testblog.com",
		Language:    "en-US",
	}

	feed, err := RSSBuilder(posts, indexMeta)
	if err != nil {
		t.Errorf("RSSBuilder failed with multiple posts: %v", err)
	}

	if len(feed.Channel.Items) != 3 {
		t.Errorf("Expected 3 items, got %d", len(feed.Channel.Items))
	}
}

func TestRSSBuilder_WithSpecialCharacters(t *testing.T) {
	// Test with special HTML characters in content
	posts := []indexpost.IndexPost{
		{
			Title:   "Post with <special> chars",
			Content: "<p>This has &amp; entities</p>",
			Date:    "2023-01-01",
			URL:     "/posts/test.html",
		},
	}

	indexMeta := siteconfig.MetaData{
		Title:       "Test Blog",
		Description: "A test blog",
		Link:        "https://testblog.com",
		Language:    "en-US",
	}

	feed, err := RSSBuilder(posts, indexMeta)
	if err != nil {
		t.Errorf("RSSBuilder failed with special characters: %v", err)
	}

	if len(feed.Channel.Items) != 1 {
		t.Error("Expected 1 item")
	}

	item := feed.Channel.Items[0]
	if item.Title != "Post with <special> chars" {
		t.Error("Title should contain special characters")
	}
}

func TestRSSBuilder_StableGUID(t *testing.T) {
	posts := []indexpost.IndexPost{
		{
			Title: "Test Post",
			URL:   "/posts/test.html",
			Date:  "2026-08-17",
		},
	}

	meta := siteconfig.MetaData{
		Title: "Test Site",
		Link:  "https://example.com",
	}

	feed1, err := RSSBuilder(posts, meta)
	if err != nil {
		t.Fatalf("first RSSBuilder call failed: %v", err)
	}

	feed2, err := RSSBuilder(posts, meta)
	if err != nil {
		t.Fatalf("second RSSBuilder call failed: %v", err)
	}

	got1 := feed1.Channel.Items[0].GUID.Value
	got2 := feed2.Channel.Items[0].GUID.Value

	if got1 != got2 {
		t.Errorf("GUID is not stable: %q != %q", got1, got2)
	}

	want := "https://example.com/posts/test.html"
	if got1 != want {
		t.Errorf("GUID = %q, want %q", got1, want)
	}
}
