package builder

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"billiemuk/internal/templates"
)

func TestBuildCreatesDistFiles(t *testing.T) {
	root := t.TempDir()

	// Create project structure
	postsDir := filepath.Join(root, "content", "posts")
	templatesDir := filepath.Join(root, "templates")
	staticDir := filepath.Join(root, "static", "css")
	distDir := filepath.Join(root, "dist")

	for _, d := range []string{postsDir, templatesDir, staticDir} {
		if err := os.MkdirAll(d, 0755); err != nil {
			t.Fatal(err)
		}
	}

	// Create a test post
	post := `---
title: "Hello World"
date: 2026-01-15
summary: "My first post."
---

This is my **first** post.
`
	if err := os.WriteFile(filepath.Join(postsDir, "2026-01-15-hello-world.md"), []byte(post), 0644); err != nil {
		t.Fatal(err)
	}

	// Copy real templates from project root
	for _, name := range []string{"base.html", "home.html", "post.html"} {
		src, err := os.ReadFile(filepath.Join("..", "..", "templates", name))
		if err != nil {
			t.Fatalf("read template %s: %v", name, err)
		}
		if err := os.WriteFile(filepath.Join(templatesDir, name), src, 0644); err != nil {
			t.Fatal(err)
		}
	}

	// Create a minimal theme.css
	if err := os.WriteFile(filepath.Join(staticDir, "theme.css"), []byte(":root { color: red; }"), 0644); err != nil {
		t.Fatal(err)
	}

	cfg := Config{
		ContentDir:   filepath.Join(root, "content"),
		TemplatesDir: templatesDir,
		StaticDir:    filepath.Join(root, "static"),
		DistDir:      distDir,
		Site: templates.SiteData{
			Title:   "Test Blog",
			BaseURL: "https://example.com",
			Author:  "Test",
			Year:    2026,
			Socials: []templates.Social{{Name: "GitHub", URL: "https://github.com/test"}},
		},
		IncludeDrafts: false,
	}

	if err := Build(cfg); err != nil {
		t.Fatal(err)
	}

	// Verify index.html exists and contains post link
	indexHTML, err := os.ReadFile(filepath.Join(distDir, "index.html"))
	if err != nil {
		t.Fatal("dist/index.html not created")
	}
	if !strings.Contains(string(indexHTML), "Hello World") {
		t.Error("index.html missing post title")
	}

	// Verify post page exists
	postHTML, err := os.ReadFile(filepath.Join(distDir, "posts", "2026-01-15-hello-world", "index.html"))
	if err != nil {
		t.Fatal("post index.html not created")
	}
	if !strings.Contains(string(postHTML), "first") {
		t.Error("post HTML missing content")
	}

	// Verify minified CSS exists
	if _, err := os.Stat(filepath.Join(distDir, "static", "css", "theme.min.css")); err != nil {
		t.Error("minified CSS not created")
	}
}
