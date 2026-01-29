# Personal Site — Design Document

## Overview

A lightweight, fast-loading personal site built as a static site generator in Go. Markdown-driven blog posts, styled with PicoCSS and a Catppuccin Lavender theme, deployed to Cloudflare Pages or GitHub Pages.

## Project Structure

```
billiemuk/
├── cmd/
│   └── main.go              # CLI entry point (build / serve / new)
├── internal/
│   ├── builder/
│   │   └── builder.go       # Orchestrates the full build pipeline
│   ├── content/
│   │   └── content.go       # Parses markdown + frontmatter into structs
│   ├── templates/
│   │   └── templates.go     # Loads and renders Go html/template files
│   └── server/
│       └── server.go        # Dev server with file-watching + live-reload
├── content/
│   └── posts/
│       └── my-first-post.md # Markdown posts live here
├── templates/
│   ├── base.html            # Shared layout (head, header, footer)
│   ├── home.html            # Homepage with post list
│   └── post.html            # Individual post page
├── static/
│   └── css/
│       └── theme.css        # Catppuccin overrides for PicoCSS
├── dist/                    # Build output (gitignored)
└── go.mod
```

## Content Model

Each post is a markdown file in `content/posts/` with YAML frontmatter:

```markdown
---
title: "My First Post"
date: 2026-01-27
summary: "A short description shown on the homepage list."
draft: true
---

Post content here.
```

### Frontmatter Fields

- `title` — Required. Used in the post page and homepage list.
- `date` — Required. Determines sort order. Displayed on the homepage and post page.
- `summary` — Optional. Shown on the homepage list. If omitted, just title and date.
- `draft` — Optional, defaults to `false`. Visible in `serve` mode, excluded from `build`.

### Content Parsing

- Goldmark parses markdown to HTML.
- `go.abhg.dev/goldmark/frontmatter` handles YAML frontmatter, decoded into Go structs.
- Slug derived from filename: `2026-01-27-my-post-title.md` -> `/posts/2026-01-27-my-post-title/`
- Posts sorted by date descending.

## Pages

### Homepage

- Site name and social links (GitHub, LinkedIn) in the header/nav.
- Simple chronological list of posts: title, date, and optional summary.
- No thumbnails, no pagination (revisit when post count warrants it).
- Mobile-first responsive design via PicoCSS.

### Post Page

- Full rendered markdown content.
- Title, date, and summary in the header.
- Shared base layout with homepage.

## Theming

### PicoCSS + Catppuccin Lavender

PicoCSS loaded from CDN. Custom `theme.css` overrides PicoCSS CSS variables with Catppuccin palette values:

- **Light mode** -> Catppuccin Latte palette, Lavender as primary/accent.
- **Dark mode** -> Catppuccin Mocha palette, Lavender as primary/accent.
- Automatic switching via `prefers-color-scheme` — no JavaScript, no manual toggle.

```css
:root {
  --pico-primary: #7287fd; /* Latte Lavender */
}

@media (prefers-color-scheme: dark) {
  :root {
    --pico-primary: #b4befe; /* Mocha Lavender */
  }
}
```

## Build Pipeline

Sequential steps executed by the `build` command:

1. **Parse content** — Read all markdown, extract frontmatter, convert to HTML.
2. **Render templates** — Feed content into `html/template`, output HTML to `dist/`.
3. **Process images** — Scan rendered HTML for `<img>` tags, compress source images to WebP, output to `dist/images/`, rewrite `src` attributes to point to optimized versions.
4. **Process static assets** — Copy `static/` to `dist/`, minify CSS via `tdewolff/minify`.
5. **Generate SEO files** — `sitemap.xml`, `robots.txt`, `feed.xml`.

### Image Handling

- Raw images stored in `content/images/`.
- Markdown references raw paths: `![alt](images/screenshot.png)`.
- Build step compresses and converts to WebP, outputs to `dist/images/`.
- HTML `src` attributes rewritten automatically to `/images/filename.webp`.

### Build Output

```
dist/
├── index.html
├── feed.xml
├── sitemap.xml
├── robots.txt
├── images/
│   └── screenshot.webp
├── posts/
│   └── 2026-01-27-my-post-title/
│       └── index.html
└── static/
    └── css/
        └── theme.min.css
```

## SEO

- `<title>` — Post title + site name.
- `<meta name="description">` — From `summary` frontmatter field.
- `<link rel="canonical">` — Full URL.
- Open Graph tags: `og:title`, `og:description`, `og:type`, `og:url`.
- `<meta name="viewport">` — Mobile-first, handled by PicoCSS.
- `sitemap.xml` — Auto-generated with all published post URLs.
- `robots.txt` — Allow-all, pointer to sitemap.
- `feed.xml` — RSS feed, auto-discoverable via `<link rel="alternate">`.
- Semantic HTML: `<article>`, `<header>`, `<main>`, `<nav>`.
- Clean URLs: `/posts/post-slug/`.

## CLI Commands

### `go run . build`

Full build pipeline. Outputs to `dist/`. Excludes drafts.

### `go run . serve`

1. Runs initial build.
2. Starts HTTP server on `localhost:8080` serving `dist/`.
3. Watches `content/`, `templates/`, `static/` via `fsnotify`.
4. On change, re-runs build and triggers browser reload via SSE.
5. Injects a small `<script>` tag in dev mode only for SSE reload.
6. Includes draft posts.

### `go run . new "Post Title"`

Creates `content/posts/YYYY-MM-DD-post-title.md` with pre-filled frontmatter (`title`, `date`, `draft: true`). Prints dev preview URL.

## Dependencies

- `github.com/yuin/goldmark` — Markdown parser.
- `go.abhg.dev/goldmark/frontmatter` — YAML frontmatter extension for Goldmark.
- `github.com/tdewolff/minify` — CSS minification.
- `github.com/fsnotify/fsnotify` — File system watching for dev server.
- Go standard library: `html/template`, `net/http`, `image` packages.

## Deployment

Static output in `dist/` deployed to Cloudflare Pages or GitHub Pages. CI runs `go run . build` on push, deploys the output directory.

## Not In Scope (For Now)

- Post thumbnails on homepage.
- Pagination.
- Categories, tags, or taxonomy.
- JSON-LD / structured data.
- Analytics.
- Manual light/dark toggle.
- Dynamic features (can be added on a subdomain later).
