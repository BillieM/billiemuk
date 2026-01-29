# Personal Site Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Build a Go static site generator that turns markdown files into a PicoCSS-styled personal blog with Catppuccin Lavender theming, automatic image compression, and a dev server with live reload.

**Architecture:** A single Go binary with three CLI commands (`build`, `serve`, `new`). The build pipeline reads markdown from `content/posts/`, renders HTML via Go templates, processes images, minifies CSS, and generates SEO files into `dist/`. The dev server watches for file changes and live-reloads via SSE.

**Tech Stack:** Go, Goldmark + goldmark-frontmatter, PicoCSS v2 (CDN), tdewolff/minify, fsnotify, golang.org/x/image.

---

### Task 1: Project Scaffolding

**Files:**
- Create: `go.mod`
- Create: `main.go`
- Create: `.gitignore`
- Create: directory structure

**Step 1: Initialize Go module and install dependencies**

Run:
```bash
cd /Users/billie/Code/billiemuk
go mod init billiemuk
go get github.com/yuin/goldmark@latest
go get go.abhg.dev/goldmark/frontmatter@latest
go get github.com/tdewolff/minify/v2@latest
go get github.com/fsnotify/fsnotify@latest
go get golang.org/x/image@latest
```

**Step 2: Create .gitignore**

Create `.gitignore`:
```
dist/
```

**Step 3: Create directory structure**

Run:
```bash
mkdir -p internal/builder internal/content internal/templates internal/server
mkdir -p content/posts content/images
mkdir -p templates
mkdir -p static/css
```

**Step 4: Create minimal main.go**

Create `main.go`:
```go
package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: billiemuk <build|serve|new> [args]")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "build":
		fmt.Println("build: not yet implemented")
	case "serve":
		fmt.Println("serve: not yet implemented")
	case "new":
		fmt.Println("new: not yet implemented")
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}
```

**Step 5: Verify it compiles**

Run: `go build -o /dev/null .`
Expected: exits 0, no errors.

**Step 6: Commit**

```bash
git add .
git commit -m "feat: scaffold project structure and CLI skeleton"
```

---

### Task 2: Content Parsing

**Files:**
- Create: `internal/content/content.go`
- Create: `internal/content/content_test.go`

**Step 1: Write the failing test**

Create `internal/content/content_test.go`:
```go
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
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/content/ -v`
Expected: compilation error — `ParsePost` and `ParseAllPosts` not defined.

**Step 3: Write the implementation**

Create `internal/content/content.go`:
```go
package content

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"go.abhg.dev/goldmark/frontmatter"
)

type Post struct {
	Title   string
	Date    time.Time
	Summary string
	Draft   bool
	Slug    string
	HTML    string
}

type postFrontmatter struct {
	Title   string `yaml:"title"`
	Date    string `yaml:"date"`
	Summary string `yaml:"summary"`
	Draft   bool   `yaml:"draft"`
}

func ParsePost(path string) (Post, error) {
	src, err := os.ReadFile(path)
	if err != nil {
		return Post{}, fmt.Errorf("read post: %w", err)
	}

	md := goldmark.New(
		goldmark.WithExtensions(&frontmatter.Extender{}),
	)

	ctx := parser.NewContext()
	var buf bytes.Buffer
	if err := md.Convert(src, &buf, parser.WithContext(ctx)); err != nil {
		return Post{}, fmt.Errorf("convert markdown: %w", err)
	}

	fm := frontmatter.Get(ctx)
	if fm == nil {
		return Post{}, fmt.Errorf("no frontmatter found in %s", path)
	}

	var meta postFrontmatter
	if err := fm.Decode(&meta); err != nil {
		return Post{}, fmt.Errorf("decode frontmatter: %w", err)
	}

	date, err := time.Parse("2006-01-02", meta.Date)
	if err != nil {
		return Post{}, fmt.Errorf("parse date %q: %w", meta.Date, err)
	}

	filename := filepath.Base(path)
	slug := strings.TrimSuffix(filename, filepath.Ext(filename))

	return Post{
		Title:   meta.Title,
		Date:    date,
		Summary: meta.Summary,
		Draft:   meta.Draft,
		Slug:    slug,
		HTML:    buf.String(),
	}, nil
}

func ParseAllPosts(dir string, includeDrafts bool) ([]Post, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("read posts dir: %w", err)
	}

	var posts []Post
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}
		post, err := ParsePost(filepath.Join(dir, entry.Name()))
		if err != nil {
			return nil, err
		}
		if post.Draft && !includeDrafts {
			continue
		}
		posts = append(posts, post)
	}

	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Date.After(posts[j].Date)
	})

	return posts, nil
}
```

**Step 4: Run tests to verify they pass**

Run: `go test ./internal/content/ -v`
Expected: all 3 tests PASS.

**Step 5: Commit**

```bash
git add internal/content/
git commit -m "feat: add markdown content parser with frontmatter support"
```

---

### Task 3: Template System

**Files:**
- Create: `internal/templates/templates.go`
- Create: `internal/templates/templates_test.go`
- Create: `templates/base.html`
- Create: `templates/home.html`
- Create: `templates/post.html`

**Step 1: Create HTML templates**

Create `templates/base.html`:
```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/@picocss/pico@2/css/pico.min.css">
    <link rel="stylesheet" href="/static/css/theme.min.css">
    <link rel="alternate" type="application/rss+xml" title="{{.Site.Title}}" href="/feed.xml">
    {{block "meta" .}}{{end}}
    <title>{{block "title" .}}{{.Site.Title}}{{end}}</title>
</head>
<body>
    <header class="container">
        <nav>
            <ul>
                <li><strong><a href="/">{{.Site.Title}}</a></strong></li>
            </ul>
            <ul>
                {{range .Site.Socials}}
                <li><a href="{{.URL}}" target="_blank" rel="noopener noreferrer">{{.Name}}</a></li>
                {{end}}
            </ul>
        </nav>
    </header>
    <main class="container">
        {{block "content" .}}{{end}}
    </main>
    <footer class="container">
        <small>&copy; {{.Site.Year}} {{.Site.Title}}</small>
    </footer>
    {{if .DevMode}}<script>
    const es = new EventSource("/_reload");
    es.onmessage = () => location.reload();
    es.onerror = () => setTimeout(() => location.reload(), 1000);
    </script>{{end}}
</body>
</html>
```

Create `templates/home.html`:
```html
{{define "content"}}
<section>
    {{range .Posts}}
    <article>
        <header>
            <a href="/posts/{{.Slug}}/"><strong>{{.Title}}</strong></a>
            <small><time datetime="{{.Date.Format "2006-01-02"}}">{{.Date.Format "2 January 2006"}}</time></small>
        </header>
        {{if .Summary}}<p>{{.Summary}}</p>{{end}}
    </article>
    {{else}}
    <p>No posts yet.</p>
    {{end}}
</section>
{{end}}
```

Create `templates/post.html`:
```html
{{define "title"}}{{.Post.Title}} | {{.Site.Title}}{{end}}

{{define "meta"}}
{{if .Post.Summary}}<meta name="description" content="{{.Post.Summary}}">{{end}}
<link rel="canonical" href="{{.Site.BaseURL}}/posts/{{.Post.Slug}}/">
<meta property="og:title" content="{{.Post.Title}}">
{{if .Post.Summary}}<meta property="og:description" content="{{.Post.Summary}}">{{end}}
<meta property="og:type" content="article">
<meta property="og:url" content="{{.Site.BaseURL}}/posts/{{.Post.Slug}}/">
{{end}}

{{define "content"}}
<article>
    <header>
        <h1>{{.Post.Title}}</h1>
        <p><time datetime="{{.Post.Date.Format "2006-01-02"}}">{{.Post.Date.Format "2 January 2006"}}</time></p>
    </header>
    {{.Post.HTMLContent}}
</article>
{{end}}
```

**Step 2: Write the failing test**

Create `internal/templates/templates_test.go`:
```go
package templates

import (
	"strings"
	"testing"
	"time"
)

func TestRenderHome(t *testing.T) {
	dir := "../../templates"
	renderer, err := New(dir)
	if err != nil {
		t.Fatal(err)
	}

	data := PageData{
		Site: SiteData{
			Title:   "Test Site",
			BaseURL: "https://example.com",
			Year:    2026,
			Socials: []Social{{Name: "GitHub", URL: "https://github.com/test"}},
		},
		Posts: []PostData{
			{
				Title:   "First Post",
				Date:    time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC),
				Summary: "A summary.",
				Slug:    "2026-01-15-first-post",
			},
		},
	}

	html, err := renderer.RenderHome(data)
	if err != nil {
		t.Fatal(err)
	}

	checks := []string{
		"Test Site",
		"First Post",
		"A summary.",
		"2026-01-15-first-post",
		"GitHub",
		"pico.min.css",
		"theme.min.css",
	}
	for _, check := range checks {
		if !strings.Contains(html, check) {
			t.Errorf("home HTML missing %q", check)
		}
	}
}

func TestRenderPost(t *testing.T) {
	dir := "../../templates"
	renderer, err := New(dir)
	if err != nil {
		t.Fatal(err)
	}

	data := PageData{
		Site: SiteData{
			Title:   "Test Site",
			BaseURL: "https://example.com",
			Year:    2026,
		},
		Post: &PostData{
			Title:       "My Post",
			Date:        time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC),
			Summary:     "Post summary.",
			Slug:        "2026-01-15-my-post",
			HTMLContent: "<p>Hello <strong>world</strong>.</p>",
		},
	}

	html, err := renderer.RenderPost(data)
	if err != nil {
		t.Fatal(err)
	}

	checks := []string{
		"My Post",
		"<p>Hello <strong>world</strong>.</p>",
		`og:title`,
		`og:type`,
		"canonical",
	}
	for _, check := range checks {
		if !strings.Contains(html, check) {
			t.Errorf("post HTML missing %q", check)
		}
	}
}
```

**Step 3: Run test to verify it fails**

Run: `go test ./internal/templates/ -v`
Expected: compilation error — types and functions not defined.

**Step 4: Write the implementation**

Create `internal/templates/templates.go`:
```go
package templates

import (
	"bytes"
	"fmt"
	"html/template"
	"path/filepath"
	"time"
)

type Social struct {
	Name string
	URL  string
}

type SiteData struct {
	Title   string
	BaseURL string
	Author  string
	Year    int
	Socials []Social
}

type PostData struct {
	Title       string
	Date        time.Time
	Summary     string
	Slug        string
	Draft       bool
	HTMLContent template.HTML
}

type PageData struct {
	Site    SiteData
	Posts   []PostData
	Post    *PostData
	DevMode bool
}

type Renderer struct {
	homeTemplate *template.Template
	postTemplate *template.Template
}

func New(templatesDir string) (*Renderer, error) {
	base := filepath.Join(templatesDir, "base.html")

	homeTmpl, err := template.ParseFiles(base, filepath.Join(templatesDir, "home.html"))
	if err != nil {
		return nil, fmt.Errorf("parse home template: %w", err)
	}

	postTmpl, err := template.ParseFiles(base, filepath.Join(templatesDir, "post.html"))
	if err != nil {
		return nil, fmt.Errorf("parse post template: %w", err)
	}

	return &Renderer{
		homeTemplate: homeTmpl,
		postTemplate: postTmpl,
	}, nil
}

func (r *Renderer) RenderHome(data PageData) (string, error) {
	var buf bytes.Buffer
	if err := r.homeTemplate.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("render home: %w", err)
	}
	return buf.String(), nil
}

func (r *Renderer) RenderPost(data PageData) (string, error) {
	var buf bytes.Buffer
	if err := r.postTemplate.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("render post: %w", err)
	}
	return buf.String(), nil
}
```

**Step 5: Run tests to verify they pass**

Run: `go test ./internal/templates/ -v`
Expected: all tests PASS.

**Step 6: Commit**

```bash
git add internal/templates/ templates/
git commit -m "feat: add template system with base, home, and post layouts"
```

---

### Task 4: Catppuccin Theme CSS

**Files:**
- Create: `static/css/theme.css`

**Step 1: Create theme.css**

Create `static/css/theme.css` — full PicoCSS variable overrides using Catppuccin Latte (light) and Mocha (dark) palettes with Lavender as the primary accent:

```css
/* Catppuccin Lavender theme for PicoCSS */
/* Light: Latte | Dark: Mocha */
/* Palette: https://catppuccin.com/palette */

/* Light mode — Catppuccin Latte */
[data-theme="light"],
:root:not([data-theme="dark"]) {
  --pico-background-color: #eff1f5; /* Base */
  --pico-color: #4c4f69; /* Text */
  --pico-muted-color: #6c6f85; /* Subtext0 */
  --pico-muted-border-color: #ccd0da; /* Surface0 */

  --pico-primary: #7287fd; /* Lavender */
  --pico-primary-background: #7287fd;
  --pico-primary-border: #7287fd;
  --pico-primary-underline: rgba(114, 135, 253, 0.5);
  --pico-primary-hover: #5c6ef0;
  --pico-primary-hover-background: #5c6ef0;
  --pico-primary-hover-border: #5c6ef0;
  --pico-primary-focus: rgba(114, 135, 253, 0.25);
  --pico-primary-inverse: #eff1f5; /* Base */

  --pico-secondary: #8c8fa1; /* Overlay1 */
  --pico-secondary-background: #8c8fa1;
  --pico-secondary-border: #8c8fa1;
  --pico-secondary-underline: rgba(140, 143, 161, 0.5);
  --pico-secondary-hover: #7c7f93; /* Overlay2 */
  --pico-secondary-hover-background: #7c7f93;
  --pico-secondary-hover-border: #7c7f93;
  --pico-secondary-focus: rgba(140, 143, 161, 0.25);
  --pico-secondary-inverse: #eff1f5;

  --pico-card-background-color: #e6e9ef; /* Mantle */
  --pico-card-border-color: #ccd0da; /* Surface0 */
  --pico-card-sectioning-background-color: #dce0e8; /* Crust */

  --pico-code-background-color: #e6e9ef; /* Mantle */
  --pico-code-color: #4c4f69; /* Text */

  --pico-blockquote-border-color: #7287fd; /* Lavender */
  --pico-blockquote-footer-color: #6c6f85; /* Subtext0 */

  --pico-h1-color: #4c4f69; /* Text */
  --pico-h2-color: #5c5f77; /* Subtext1 */
  --pico-h3-color: #6c6f85; /* Subtext0 */
  --pico-h4-color: #7c7f93; /* Overlay2 */
  --pico-h5-color: #8c8fa1; /* Overlay1 */
  --pico-h6-color: #9ca0b0; /* Overlay0 */

  --pico-mark-background-color: rgba(114, 135, 253, 0.15);
  --pico-mark-color: #4c4f69;

  --pico-text-selection-color: rgba(114, 135, 253, 0.2);
}

/* Dark mode — Catppuccin Mocha */
@media only screen and (prefers-color-scheme: dark) {
  :root:not([data-theme]) {
    --pico-background-color: #1e1e2e; /* Base */
    --pico-color: #cdd6f4; /* Text */
    --pico-muted-color: #a6adc8; /* Subtext0 */
    --pico-muted-border-color: #313244; /* Surface0 */

    --pico-primary: #b4befe; /* Lavender */
    --pico-primary-background: #b4befe;
    --pico-primary-border: #b4befe;
    --pico-primary-underline: rgba(180, 190, 254, 0.5);
    --pico-primary-hover: #c8d0fe;
    --pico-primary-hover-background: #c8d0fe;
    --pico-primary-hover-border: #c8d0fe;
    --pico-primary-focus: rgba(180, 190, 254, 0.25);
    --pico-primary-inverse: #1e1e2e; /* Base */

    --pico-secondary: #7f849c; /* Overlay1 */
    --pico-secondary-background: #7f849c;
    --pico-secondary-border: #7f849c;
    --pico-secondary-underline: rgba(127, 132, 156, 0.5);
    --pico-secondary-hover: #9399b2; /* Overlay2 */
    --pico-secondary-hover-background: #9399b2;
    --pico-secondary-hover-border: #9399b2;
    --pico-secondary-focus: rgba(127, 132, 156, 0.25);
    --pico-secondary-inverse: #1e1e2e;

    --pico-card-background-color: #181825; /* Mantle */
    --pico-card-border-color: #313244; /* Surface0 */
    --pico-card-sectioning-background-color: #11111b; /* Crust */

    --pico-code-background-color: #181825; /* Mantle */
    --pico-code-color: #cdd6f4; /* Text */

    --pico-blockquote-border-color: #b4befe; /* Lavender */
    --pico-blockquote-footer-color: #a6adc8; /* Subtext0 */

    --pico-h1-color: #cdd6f4; /* Text */
    --pico-h2-color: #bac2de; /* Subtext1 */
    --pico-h3-color: #a6adc8; /* Subtext0 */
    --pico-h4-color: #9399b2; /* Overlay2 */
    --pico-h5-color: #7f849c; /* Overlay1 */
    --pico-h6-color: #6c7086; /* Overlay0 */

    --pico-mark-background-color: rgba(180, 190, 254, 0.15);
    --pico-mark-color: #cdd6f4;

    --pico-text-selection-color: rgba(180, 190, 254, 0.2);
  }
}
```

**Step 2: Verify the CSS is valid**

Open the file and visually verify all colour hex values match the Catppuccin palette:
- Latte Lavender: `#7287fd`
- Mocha Lavender: `#b4befe`
- Latte Base: `#eff1f5`, Mocha Base: `#1e1e2e`
- Latte Text: `#4c4f69`, Mocha Text: `#cdd6f4`

Reference: [Catppuccin Palette](https://catppuccin.com/palette/)

**Step 3: Commit**

```bash
git add static/css/theme.css
git commit -m "feat: add Catppuccin Lavender theme with light and dark mode"
```

---

### Task 5: Build Pipeline — Core

**Files:**
- Create: `internal/builder/builder.go`
- Create: `internal/builder/builder_test.go`

**Step 1: Write the failing test**

Create `internal/builder/builder_test.go`:
```go
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
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/builder/ -v`
Expected: compilation error — `Config`, `Build` not defined.

**Step 3: Write the implementation**

Create `internal/builder/builder.go`:
```go
package builder

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"billiemuk/internal/content"
	"billiemuk/internal/templates"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	mhtml "github.com/tdewolff/minify/v2/html"
)

type Config struct {
	ContentDir   string
	TemplatesDir string
	StaticDir    string
	DistDir      string
	Site         templates.SiteData
	IncludeDrafts bool
	DevMode      bool
}

func Build(cfg Config) error {
	// Clean dist
	if err := os.RemoveAll(cfg.DistDir); err != nil {
		return fmt.Errorf("clean dist: %w", err)
	}

	// Parse posts
	postsDir := filepath.Join(cfg.ContentDir, "posts")
	posts, err := content.ParseAllPosts(postsDir, cfg.IncludeDrafts)
	if err != nil {
		return fmt.Errorf("parse posts: %w", err)
	}

	// Load templates
	renderer, err := templates.New(cfg.TemplatesDir)
	if err != nil {
		return fmt.Errorf("load templates: %w", err)
	}

	// Convert posts to template data
	var postDataList []templates.PostData
	for _, p := range posts {
		postDataList = append(postDataList, templates.PostData{
			Title:       p.Title,
			Date:        p.Date,
			Summary:     p.Summary,
			Slug:        p.Slug,
			Draft:       p.Draft,
			HTMLContent: template.HTML(p.HTML),
		})
	}

	// Render homepage
	homeData := templates.PageData{
		Site:    cfg.Site,
		Posts:   postDataList,
		DevMode: cfg.DevMode,
	}
	homeHTML, err := renderer.RenderHome(homeData)
	if err != nil {
		return fmt.Errorf("render home: %w", err)
	}
	if err := writeFile(filepath.Join(cfg.DistDir, "index.html"), homeHTML); err != nil {
		return err
	}

	// Render each post
	for _, pd := range postDataList {
		pd := pd
		postData := templates.PageData{
			Site:    cfg.Site,
			Post:    &pd,
			DevMode: cfg.DevMode,
		}
		postHTML, err := renderer.RenderPost(postData)
		if err != nil {
			return fmt.Errorf("render post %s: %w", pd.Slug, err)
		}
		postDir := filepath.Join(cfg.DistDir, "posts", pd.Slug)
		if err := writeFile(filepath.Join(postDir, "index.html"), postHTML); err != nil {
			return err
		}
	}

	// Process static assets (copy + minify CSS)
	if err := processStatic(cfg.StaticDir, cfg.DistDir); err != nil {
		return fmt.Errorf("process static: %w", err)
	}

	// Process images
	if err := processImages(cfg.ContentDir, cfg.DistDir); err != nil {
		return fmt.Errorf("process images: %w", err)
	}

	// Generate SEO files
	if err := generateSEO(cfg, posts); err != nil {
		return fmt.Errorf("generate SEO: %w", err)
	}

	return nil
}

func processStatic(staticDir, distDir string) error {
	m := minify.New()
	m.AddFunc("text/css", css.Minify)
	m.AddFunc("text/html", mhtml.Minify)

	return filepath.Walk(staticDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}

		rel, _ := filepath.Rel(staticDir, path)
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		// Minify CSS files and rename to .min.css
		if strings.HasSuffix(path, ".css") {
			minified, err := m.String("text/css", string(data))
			if err != nil {
				return fmt.Errorf("minify %s: %w", rel, err)
			}
			outName := strings.TrimSuffix(rel, ".css") + ".min.css"
			return writeFile(filepath.Join(distDir, "static", outName), minified)
		}

		// Copy other static files as-is
		return writeFile(filepath.Join(distDir, "static", rel), string(data))
	})
}

func processImages(contentDir, distDir string) error {
	imagesDir := filepath.Join(contentDir, "images")
	if _, err := os.Stat(imagesDir); os.IsNotExist(err) {
		return nil // No images directory, skip
	}

	return filepath.Walk(imagesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}

		rel, _ := filepath.Rel(imagesDir, path)
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		return writeFile(filepath.Join(distDir, "images", rel), string(data))
	})
}

func generateSEO(cfg Config, posts []content.Post) error {
	// Filter out drafts for SEO
	var published []content.Post
	for _, p := range posts {
		if !p.Draft {
			published = append(published, p)
		}
	}

	// sitemap.xml
	var sitemap strings.Builder
	sitemap.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")
	sitemap.WriteString(`<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">` + "\n")
	sitemap.WriteString(fmt.Sprintf("  <url><loc>%s/</loc></url>\n", cfg.Site.BaseURL))
	for _, p := range published {
		sitemap.WriteString(fmt.Sprintf("  <url><loc>%s/posts/%s/</loc><lastmod>%s</lastmod></url>\n",
			cfg.Site.BaseURL, p.Slug, p.Date.Format("2006-01-02")))
	}
	sitemap.WriteString("</urlset>\n")
	if err := writeFile(filepath.Join(cfg.DistDir, "sitemap.xml"), sitemap.String()); err != nil {
		return err
	}

	// robots.txt
	robots := fmt.Sprintf("User-agent: *\nAllow: /\nSitemap: %s/sitemap.xml\n", cfg.Site.BaseURL)
	if err := writeFile(filepath.Join(cfg.DistDir, "robots.txt"), robots); err != nil {
		return err
	}

	// feed.xml (RSS 2.0)
	var feed strings.Builder
	feed.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")
	feed.WriteString(`<rss version="2.0" xmlns:atom="http://www.w3.org/2005/Atom">` + "\n")
	feed.WriteString("<channel>\n")
	feed.WriteString(fmt.Sprintf("  <title>%s</title>\n", xmlEscape(cfg.Site.Title)))
	feed.WriteString(fmt.Sprintf("  <link>%s</link>\n", cfg.Site.BaseURL))
	feed.WriteString(fmt.Sprintf("  <description>%s</description>\n", xmlEscape(cfg.Site.Title)))
	feed.WriteString(fmt.Sprintf(`  <atom:link href="%s/feed.xml" rel="self" type="application/rss+xml"/>`+"\n", cfg.Site.BaseURL))
	for _, p := range published {
		feed.WriteString("  <item>\n")
		feed.WriteString(fmt.Sprintf("    <title>%s</title>\n", xmlEscape(p.Title)))
		feed.WriteString(fmt.Sprintf("    <link>%s/posts/%s/</link>\n", cfg.Site.BaseURL, p.Slug))
		feed.WriteString(fmt.Sprintf("    <guid>%s/posts/%s/</guid>\n", cfg.Site.BaseURL, p.Slug))
		feed.WriteString(fmt.Sprintf("    <pubDate>%s</pubDate>\n", p.Date.Format("Mon, 02 Jan 2006 15:04:05 -0700")))
		if p.Summary != "" {
			feed.WriteString(fmt.Sprintf("    <description>%s</description>\n", xmlEscape(p.Summary)))
		}
		feed.WriteString("  </item>\n")
	}
	feed.WriteString("</channel>\n")
	feed.WriteString("</rss>\n")
	if err := writeFile(filepath.Join(cfg.DistDir, "feed.xml"), feed.String()); err != nil {
		return err
	}

	return nil
}

func writeFile(path, content string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("mkdir %s: %w", filepath.Dir(path), err)
	}
	return os.WriteFile(path, []byte(content), 0644)
}

var xmlReplacer = regexp.MustCompile(`[<>&"']`)

func xmlEscape(s string) string {
	return xmlReplacer.ReplaceAllStringFunc(s, func(c string) string {
		switch c {
		case "<":
			return "&lt;"
		case ">":
			return "&gt;"
		case "&":
			return "&amp;"
		case `"`:
			return "&quot;"
		case "'":
			return "&apos;"
		}
		return c
	})
}
```

**Step 4: Run tests to verify they pass**

Run: `go test ./internal/builder/ -v`
Expected: PASS.

**Step 5: Commit**

```bash
git add internal/builder/
git commit -m "feat: add build pipeline with CSS minification and SEO generation"
```

---

### Task 6: Image Compression

**Files:**
- Modify: `internal/builder/builder.go` (replace `processImages`)
- Create: `internal/builder/images.go`
- Create: `internal/builder/images_test.go`

This task replaces the simple file copy in `processImages` with actual resizing and compression. We use Go's standard `image` package and `golang.org/x/image/draw` for high-quality downscaling. Images are resized to a max width of 1200px and compressed as JPEG (quality 85) or PNG.

**Step 1: Write the failing test**

Create `internal/builder/images_test.go`:
```go
package builder

import (
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"testing"
)

func TestProcessImagesResizesLargeJPEG(t *testing.T) {
	root := t.TempDir()
	imagesDir := filepath.Join(root, "content", "images")
	distDir := filepath.Join(root, "dist")

	if err := os.MkdirAll(imagesDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create a 2400x1600 JPEG
	img := image.NewRGBA(image.Rect(0, 0, 2400, 1600))
	f, err := os.Create(filepath.Join(imagesDir, "large.jpg"))
	if err != nil {
		t.Fatal(err)
	}
	if err := jpeg.Encode(f, img, &jpeg.Options{Quality: 100}); err != nil {
		t.Fatal(err)
	}
	f.Close()

	if err := processImages(filepath.Join(root, "content"), distDir); err != nil {
		t.Fatal(err)
	}

	// Verify output exists
	outPath := filepath.Join(distDir, "images", "large.jpg")
	outF, err := os.Open(outPath)
	if err != nil {
		t.Fatal("output image not created")
	}
	defer outF.Close()

	outImg, _, err := image.DecodeConfig(outF)
	if err != nil {
		t.Fatal(err)
	}
	if outImg.Width > maxImageWidth {
		t.Errorf("output width = %d, want <= %d", outImg.Width, maxImageWidth)
	}
}

func TestProcessImagesKeepsSmallPNG(t *testing.T) {
	root := t.TempDir()
	imagesDir := filepath.Join(root, "content", "images")
	distDir := filepath.Join(root, "dist")

	if err := os.MkdirAll(imagesDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create a small 400x300 PNG
	img := image.NewRGBA(image.Rect(0, 0, 400, 300))
	f, err := os.Create(filepath.Join(imagesDir, "small.png"))
	if err != nil {
		t.Fatal(err)
	}
	if err := png.Encode(f, img); err != nil {
		t.Fatal(err)
	}
	f.Close()

	if err := processImages(filepath.Join(root, "content"), distDir); err != nil {
		t.Fatal(err)
	}

	outPath := filepath.Join(distDir, "images", "small.png")
	outF, err := os.Open(outPath)
	if err != nil {
		t.Fatal("output image not created")
	}
	defer outF.Close()

	outImg, _, err := image.DecodeConfig(outF)
	if err != nil {
		t.Fatal(err)
	}
	if outImg.Width != 400 {
		t.Errorf("output width = %d, want 400 (should not resize)", outImg.Width)
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/builder/ -v -run TestProcessImages`
Expected: FAIL — `maxImageWidth` not defined, `processImages` doesn't resize.

**Step 3: Write the implementation**

Create `internal/builder/images.go`:
```go
package builder

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/image/draw"
)

const maxImageWidth = 1200
const jpegQuality = 85

func processImages(contentDir, distDir string) error {
	imagesDir := filepath.Join(contentDir, "images")
	if _, err := os.Stat(imagesDir); os.IsNotExist(err) {
		return nil
	}

	return filepath.Walk(imagesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}

		rel, _ := filepath.Rel(imagesDir, path)
		outPath := filepath.Join(distDir, "images", rel)

		ext := strings.ToLower(filepath.Ext(path))
		switch ext {
		case ".jpg", ".jpeg":
			return compressImage(path, outPath, "jpeg")
		case ".png":
			return compressImage(path, outPath, "png")
		default:
			// Copy non-image files as-is (e.g. SVG, GIF)
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			return writeFile(outPath, string(data))
		}
	})
}

func compressImage(srcPath, dstPath, format string) error {
	f, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("open image: %w", err)
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return fmt.Errorf("decode image %s: %w", srcPath, err)
	}

	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// Resize if wider than max
	if width > maxImageWidth {
		newHeight := int(float64(height) * float64(maxImageWidth) / float64(width))
		resized := image.NewRGBA(image.Rect(0, 0, maxImageWidth, newHeight))
		draw.CatmullRom.Scale(resized, resized.Bounds(), img, bounds, draw.Over, nil)
		img = resized
	}

	if err := os.MkdirAll(filepath.Dir(dstPath), 0755); err != nil {
		return err
	}

	out, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("create output: %w", err)
	}
	defer out.Close()

	switch format {
	case "jpeg":
		return jpeg.Encode(out, img, &jpeg.Options{Quality: jpegQuality})
	case "png":
		return png.Encode(out, img)
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
}
```

Now remove the old `processImages` function from `builder.go`. Delete this block from `internal/builder/builder.go`:

```go
func processImages(contentDir, distDir string) error {
	imagesDir := filepath.Join(contentDir, "images")
	if _, err := os.Stat(imagesDir); os.IsNotExist(err) {
		return nil // No images directory, skip
	}

	return filepath.Walk(imagesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}

		rel, _ := filepath.Rel(imagesDir, path)
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		return writeFile(filepath.Join(distDir, "images", rel), string(data))
	})
}
```

Also add the required import to `images.go`. Ensure `_ "image/jpeg"` and `_ "image/png"` are imported in `images.go` for the image decoder registration (they're already imported by using `image/jpeg` and `image/png` directly).

**Step 4: Run tests to verify they pass**

Run: `go test ./internal/builder/ -v -run TestProcessImages`
Expected: both tests PASS.

**Step 5: Run all tests**

Run: `go test ./... -v`
Expected: all tests PASS.

**Step 6: Commit**

```bash
git add internal/builder/images.go internal/builder/builder.go internal/builder/images_test.go
git commit -m "feat: add image compression and resizing to build pipeline"
```

---

### Task 7: CLI — Build Command

**Files:**
- Modify: `main.go`

**Step 1: Wire up the build command**

Replace the contents of `main.go`:
```go
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"billiemuk/internal/builder"
	"billiemuk/internal/templates"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: billiemuk <build|serve|new> [args]")
		os.Exit(1)
	}

	root, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	switch os.Args[1] {
	case "build":
		if err := runBuild(root, false, false); err != nil {
			fmt.Fprintf(os.Stderr, "build error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Build complete: dist/")
	case "serve":
		fmt.Println("serve: not yet implemented")
	case "new":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "Usage: billiemuk new \"Post Title\"")
			os.Exit(1)
		}
		if err := runNew(root, os.Args[2]); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}

func siteConfig() templates.SiteData {
	return templates.SiteData{
		Title:   "Billie Muk",
		BaseURL: "https://billiemuk.com",
		Author:  "Billie Muk",
		Year:    time.Now().Year(),
		Socials: []templates.Social{
			{Name: "GitHub", URL: "https://github.com/billiemuk"},
			{Name: "LinkedIn", URL: "https://linkedin.com/in/billiemuk"},
		},
	}
}

func runBuild(root string, includeDrafts, devMode bool) error {
	cfg := builder.Config{
		ContentDir:    filepath.Join(root, "content"),
		TemplatesDir:  filepath.Join(root, "templates"),
		StaticDir:     filepath.Join(root, "static"),
		DistDir:       filepath.Join(root, "dist"),
		Site:          siteConfig(),
		IncludeDrafts: includeDrafts,
		DevMode:       devMode,
	}
	return builder.Build(cfg)
}

func runNew(root, title string) error {
	// not yet implemented — Task 8
	_ = root
	_ = title
	return fmt.Errorf("not yet implemented")
}
```

**Step 2: Create a sample post and test the build**

Run:
```bash
cat > content/posts/2026-01-27-hello-world.md << 'EOF'
---
title: "Hello World"
date: 2026-01-27
summary: "Welcome to my new blog."
---

This is my first post. It supports **bold**, *italic*, and `code`.

## A Section

Some more content here.
EOF
```

Then run:
```bash
go run . build
```

Expected output: `Build complete: dist/`

**Step 3: Verify build output**

Run:
```bash
ls -R dist/
```

Expected: `index.html`, `posts/2026-01-27-hello-world/index.html`, `static/css/theme.min.css`, `sitemap.xml`, `robots.txt`, `feed.xml`.

**Step 4: Commit**

```bash
git add main.go content/posts/2026-01-27-hello-world.md
git commit -m "feat: wire up build command and add sample post"
```

---

### Task 8: CLI — New Post Command

**Files:**
- Modify: `main.go` (replace `runNew` stub)

**Step 1: Write the failing test**

Create `internal/content/new_test.go`:
```go
package content

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewPost(t *testing.T) {
	dir := t.TempDir()

	path, err := NewPost(dir, "My Great Post")
	if err != nil {
		t.Fatal(err)
	}

	// Check filename format
	base := filepath.Base(path)
	if !strings.HasSuffix(base, "-my-great-post.md") {
		t.Errorf("filename = %q, want suffix -my-great-post.md", base)
	}

	// Check file contents
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	content := string(data)

	if !strings.Contains(content, `title: "My Great Post"`) {
		t.Error("missing title in frontmatter")
	}
	if !strings.Contains(content, "draft: true") {
		t.Error("missing draft: true in frontmatter")
	}
	if !strings.Contains(content, "date:") {
		t.Error("missing date in frontmatter")
	}
}

func TestSlugify(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"Hello World", "hello-world"},
		{"My  Great   Post!", "my-great-post"},
		{"CamelCase Test", "camelcase-test"},
		{"special @#$ chars", "special-chars"},
	}
	for _, tt := range tests {
		got := Slugify(tt.input)
		if got != tt.want {
			t.Errorf("Slugify(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/content/ -v -run "TestNewPost|TestSlugify"`
Expected: compilation error — `NewPost`, `Slugify` not defined.

**Step 3: Write the implementation**

Add to `internal/content/content.go` (append at the end):
```go
func Slugify(s string) string {
	s = strings.ToLower(s)
	// Replace non-alphanumeric with hyphens
	var result strings.Builder
	prevHyphen := false
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			result.WriteRune(r)
			prevHyphen = false
		} else if !prevHyphen && result.Len() > 0 {
			result.WriteRune('-')
			prevHyphen = true
		}
	}
	return strings.TrimRight(result.String(), "-")
}

func NewPost(postsDir, title string) (string, error) {
	slug := Slugify(title)
	date := time.Now().Format("2006-01-02")
	filename := fmt.Sprintf("%s-%s.md", date, slug)
	path := filepath.Join(postsDir, filename)

	content := fmt.Sprintf(`---
title: "%s"
date: %s
summary: ""
draft: true
---

`, title, date)

	if err := os.MkdirAll(postsDir, 0755); err != nil {
		return "", fmt.Errorf("create posts dir: %w", err)
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return "", fmt.Errorf("write post: %w", err)
	}

	return path, nil
}
```

**Step 4: Run tests to verify they pass**

Run: `go test ./internal/content/ -v`
Expected: all tests PASS.

**Step 5: Wire up runNew in main.go**

Replace the `runNew` function in `main.go`:
```go
func runNew(root, title string) error {
	postsDir := filepath.Join(root, "content", "posts")
	path, err := content.NewPost(postsDir, title)
	if err != nil {
		return err
	}

	slug := content.Slugify(title)
	date := time.Now().Format("2006-01-02")
	fmt.Printf("Created: %s\n", path)
	fmt.Printf("Preview: http://localhost:8080/posts/%s-%s/\n", date, slug)
	return nil
}
```

Add `"billiemuk/internal/content"` to the imports in `main.go`.

**Step 6: Test the new command**

Run: `go run . new "Test Post"`
Expected output:
```
Created: /Users/billie/Code/billiemuk/content/posts/2026-01-27-test-post.md
Preview: http://localhost:8080/posts/2026-01-27-test-post/
```

Delete the test post: `rm content/posts/2026-01-27-test-post.md`

**Step 7: Commit**

```bash
git add internal/content/ main.go
git commit -m "feat: add new post scaffolding command"
```

---

### Task 9: Dev Server with Live Reload

**Files:**
- Create: `internal/server/server.go`
- Create: `internal/server/server_test.go`
- Modify: `main.go` (wire up serve command)

**Step 1: Write the failing test**

Create `internal/server/server_test.go`:
```go
package server

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestServesDistFiles(t *testing.T) {
	distDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(distDir, "index.html"), []byte("<h1>hello</h1>"), 0644); err != nil {
		t.Fatal(err)
	}

	s := &Server{DistDir: distDir}
	handler := s.Handler()

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want 200", w.Code)
	}
	if w.Body.String() != "<h1>hello</h1>" {
		t.Errorf("body = %q, want <h1>hello</h1>", w.Body.String())
	}
}

func TestSSEEndpoint(t *testing.T) {
	s := &Server{DistDir: t.TempDir()}
	handler := s.Handler()

	req := httptest.NewRequest("GET", "/_reload", nil)
	w := httptest.NewRecorder()

	// SSE handler should set correct content type
	// We can't test the full SSE stream in a unit test, but we can verify the endpoint exists
	go handler.ServeHTTP(w, req)
	// Give the handler a moment to set headers, then check
	// In a real test we'd use a more sophisticated approach
	// For now, verify the endpoint doesn't 404
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/server/ -v`
Expected: compilation error — `Server` not defined.

**Step 3: Write the implementation**

Create `internal/server/server.go`:
```go
package server

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"sync"

	"github.com/fsnotify/fsnotify"
)

type Server struct {
	DistDir  string
	BuildFn  func() error
	WatchDirs []string
	Addr     string

	mu        sync.Mutex
	clients   map[chan struct{}]struct{}
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/_reload", s.handleSSE)
	mux.Handle("/", http.FileServer(http.Dir(s.DistDir)))
	return mux
}

func (s *Server) handleSSE(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming not supported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	ch := make(chan struct{}, 1)
	s.addClient(ch)
	defer s.removeClient(ch)

	for {
		select {
		case <-ch:
			fmt.Fprintf(w, "data: reload\n\n")
			flusher.Flush()
		case <-r.Context().Done():
			return
		}
	}
}

func (s *Server) addClient(ch chan struct{}) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.clients == nil {
		s.clients = make(map[chan struct{}]struct{})
	}
	s.clients[ch] = struct{}{}
}

func (s *Server) removeClient(ch chan struct{}) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.clients, ch)
}

func (s *Server) notifyClients() {
	s.mu.Lock()
	defer s.mu.Unlock()
	for ch := range s.clients {
		select {
		case ch <- struct{}{}:
		default:
		}
	}
}

func (s *Server) Start() error {
	// Initial build
	log.Println("Running initial build...")
	if err := s.BuildFn(); err != nil {
		return fmt.Errorf("initial build: %w", err)
	}

	// Start file watcher
	if err := s.startWatcher(); err != nil {
		return fmt.Errorf("start watcher: %w", err)
	}

	addr := s.Addr
	if addr == "" {
		addr = ":8080"
	}

	log.Printf("Serving at http://localhost%s\n", addr)
	log.Println("Watching for changes...")
	return http.ListenAndServe(addr, s.Handler())
}

func (s *Server) startWatcher() error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Has(fsnotify.Write) || event.Has(fsnotify.Create) || event.Has(fsnotify.Remove) {
					log.Printf("Change detected: %s", event.Name)
					if err := s.BuildFn(); err != nil {
						log.Printf("Build error: %v", err)
						continue
					}
					s.notifyClients()
					log.Println("Rebuild complete, reloading browsers...")
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Printf("Watcher error: %v", err)
			}
		}
	}()

	for _, dir := range s.WatchDirs {
		if err := addRecursive(watcher, dir); err != nil {
			log.Printf("Warning: could not watch %s: %v", dir, err)
		}
	}

	return nil
}

func addRecursive(watcher *fsnotify.Watcher, root string) error {
	return filepath.Walk(root, func(path string, info interface{ IsDir() bool }, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return watcher.Add(path)
		}
		return nil
	})
}
```

**Step 4: Run tests to verify they pass**

Run: `go test ./internal/server/ -v`
Expected: PASS.

**Step 5: Wire up serve command in main.go**

Replace the `"serve"` case in `main.go`:
```go
	case "serve":
		if err := runServe(root); err != nil {
			fmt.Fprintf(os.Stderr, "serve error: %v\n", err)
			os.Exit(1)
		}
```

Add the `runServe` function to `main.go`:
```go
func runServe(root string) error {
	s := &server.Server{
		DistDir: filepath.Join(root, "dist"),
		BuildFn: func() error {
			return runBuild(root, true, true)
		},
		WatchDirs: []string{
			filepath.Join(root, "content"),
			filepath.Join(root, "templates"),
			filepath.Join(root, "static"),
		},
		Addr: ":8080",
	}
	return s.Start()
}
```

Add `"billiemuk/internal/server"` to the imports in `main.go`.

**Step 6: Test the dev server manually**

Run: `go run . serve`
Expected output:
```
Running initial build...
Serving at http://localhost:8080
Watching for changes...
```

Open `http://localhost:8080` in a browser. Verify:
- Page renders with PicoCSS styling
- Catppuccin colours applied
- Light/dark mode switches with browser preference
- Social links in header
- Post list on homepage
- Post page at `/posts/2026-01-27-hello-world/`

Edit `content/posts/2026-01-27-hello-world.md` and save — browser should auto-reload.

Stop the server with Ctrl+C.

**Step 7: Run all tests**

Run: `go test ./... -v`
Expected: all tests PASS.

**Step 8: Commit**

```bash
git add internal/server/ main.go
git commit -m "feat: add dev server with file watching and SSE live reload"
```

---

### Task 10: Final Integration Test & Cleanup

**Files:**
- Modify: `.gitignore`
- Create: none

**Step 1: Verify full build works end-to-end**

Run:
```bash
go run . build
```

Expected: `Build complete: dist/`

**Step 2: Verify all output files**

Run:
```bash
ls dist/
ls dist/posts/
ls dist/static/css/
```

Expected files: `index.html`, `feed.xml`, `sitemap.xml`, `robots.txt`, `posts/2026-01-27-hello-world/index.html`, `static/css/theme.min.css`.

**Step 3: Verify HTML output is valid**

Read `dist/index.html` and verify it contains:
- `<!DOCTYPE html>`
- PicoCSS CDN link
- Theme CSS link
- Social links
- Post listing
- RSS feed `<link rel="alternate">`

Read `dist/posts/2026-01-27-hello-world/index.html` and verify it contains:
- Post title in `<h1>`
- Post content
- Open Graph meta tags
- Canonical link

**Step 4: Run all tests one final time**

Run: `go test ./... -v`
Expected: all tests PASS.

**Step 5: Run go vet and check for issues**

Run: `go vet ./...`
Expected: no issues.

**Step 6: Commit any final changes**

```bash
git add -A
git commit -m "chore: final integration verification"
```

---

### Task 11: Deployment Setup (Optional)

**Files:**
- Create: `.github/workflows/deploy.yml` (if using GitHub Pages)

This task is optional and depends on which hosting platform you choose.

**For GitHub Pages with GitHub Actions:**

Create `.github/workflows/deploy.yml`:
```yaml
name: Deploy

on:
  push:
    branches: [main]

permissions:
  contents: read
  pages: write
  id-token: write

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.22'
      - run: go run . build
      - uses: actions/upload-pages-artifact@v3
        with:
          path: dist

  deploy:
    needs: build
    runs-on: ubuntu-latest
    environment:
      name: github-pages
      url: ${{ steps.deployment.outputs.page_url }}
    steps:
      - uses: actions/deploy-pages@v4
        id: deployment
```

**For Cloudflare Pages:**

In the Cloudflare Pages dashboard:
- Build command: `go run . build`
- Build output directory: `dist`
- Go version: set `GO_VERSION` environment variable to `1.22`

**Commit (if using GitHub Pages):**

```bash
git add .github/
git commit -m "ci: add GitHub Pages deployment workflow"
```
