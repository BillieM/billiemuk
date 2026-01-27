package content

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestParsePost(t *testing.T) {
	dir := t.TempDir()
	md := `---
title: "Test Post"
date: 2026-01-15
summary: "A test summary."
draft: false
---

Hello **world**.
`
	path := filepath.Join(dir, "2026-01-15-test-post.md")
	if err := os.WriteFile(path, []byte(md), 0644); err != nil {
		t.Fatal(err)
	}

	post, err := ParsePost(path)
	if err != nil {
		t.Fatal(err)
	}

	if post.Title != "Test Post" {
		t.Errorf("title = %q, want %q", post.Title, "Test Post")
	}
	if post.Summary != "A test summary." {
		t.Errorf("summary = %q, want %q", post.Summary, "A test summary.")
	}
	if post.Draft {
		t.Error("draft = true, want false")
	}
	expectedDate := time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)
	if !post.Date.Equal(expectedDate) {
		t.Errorf("date = %v, want %v", post.Date, expectedDate)
	}
	if post.Slug != "2026-01-15-test-post" {
		t.Errorf("slug = %q, want %q", post.Slug, "2026-01-15-test-post")
	}
	if post.HTML == "" {
		t.Error("HTML is empty")
	}
	if post.HTML == md {
		t.Error("HTML was not rendered from markdown")
	}
}

func TestParsePostDraftDefault(t *testing.T) {
	dir := t.TempDir()
	md := `---
title: "No Draft Field"
date: 2026-01-15
---

Content.
`
	path := filepath.Join(dir, "2026-01-15-no-draft.md")
	if err := os.WriteFile(path, []byte(md), 0644); err != nil {
		t.Fatal(err)
	}

	post, err := ParsePost(path)
	if err != nil {
		t.Fatal(err)
	}

	if post.Draft {
		t.Error("draft should default to false")
	}
}

func TestParseAllPosts(t *testing.T) {
	dir := t.TempDir()

	posts := []struct {
		filename string
		content  string
	}{
		{"2026-01-20-second.md", "---\ntitle: \"Second\"\ndate: 2026-01-20\n---\nSecond post."},
		{"2026-01-10-first.md", "---\ntitle: \"First\"\ndate: 2026-01-10\n---\nFirst post."},
		{"2026-01-25-draft.md", "---\ntitle: \"Draft\"\ndate: 2026-01-25\ndraft: true\n---\nDraft post."},
	}
	for _, p := range posts {
		if err := os.WriteFile(filepath.Join(dir, p.filename), []byte(p.content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	// Without drafts
	result, err := ParseAllPosts(dir, false)
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 2 {
		t.Fatalf("got %d posts, want 2", len(result))
	}
	if result[0].Title != "Second" {
		t.Errorf("first post title = %q, want %q (sorted by date desc)", result[0].Title, "Second")
	}

	// With drafts
	result, err = ParseAllPosts(dir, true)
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 3 {
		t.Fatalf("got %d posts, want 3", len(result))
	}
}
