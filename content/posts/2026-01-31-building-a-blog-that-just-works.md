---
title: "Building a Blog That Just Works"
date: 2026-01-31
summary: "A technical overview of the tiny static blog stack I built in Go."
draft: false
---

I decided to start a blog. I wanted something simple: a static site I control that loads fast and is easy to write for.

The requirements were straightforward. Fast, simple, pleasant workflow. No JavaScript frameworks, no heavyweight CMS. Just markdown files that become HTML pages.

This is what I built.

---

## What I Wanted

The goals were simple:

- **Fast.** No JavaScript bloat, no client-side rendering, no framework overhead. Static HTML that loads instantly.
- **Simple.** Write markdown, get HTML. No database, no admin panel, no plugin ecosystem to maintain.
- **Pleasant workflow.** Draft a post, see it live with hot reload, deploy with a git push.
- **Good defaults.** SEO files (sitemap, robots.txt), RSS feed, and responsive design should just work.

That's it. No analytics tracking, no comment systems, no newsletter signups. Just fast pages with words on them.

---

## The Stack

I built a custom static site generator in Go. The stack:

- **[Go](https://go.dev/)** - Fast builds, single binary, stdlib has everything I need
- **[Goldmark](https://github.com/yuin/goldmark)** - Markdown parser with frontmatter support
- **[Pico CSS](https://picocss.com/)** - Classless CSS framework (semantic HTML, zero classes)
- **[GitHub Pages](https://pages.github.com/)** - Free hosting with auto-deploy via GitHub Actions

The build tool has three commands: `build`, `serve`, and `new`.

---

## How It Works

The build process is straightforward. Parse markdown files with frontmatter, convert to HTML, optimize assets, write everything to a `dist/` folder.

**Posts are markdown files with frontmatter:**

```markdown
---
title: "My Post Title"
date: 2026-01-31
summary: "A brief description"
draft: false
---

Post content here...
```

**The build pipeline:**

1. **Parse posts** - [Goldmark](https://github.com/yuin/goldmark) converts markdown to HTML, extracts frontmatter
2. **Optimize images** - Resize to max 1200px width, compress JPEGs at 85% quality
3. **Minify assets** - CSS and HTML get minified for production
4. **Generate SEO files** - Auto-create `sitemap.xml`, `robots.txt`, and `feed.xml` (RSS)
5. **Render templates** - Go templates for the homepage and post pages

[Pico CSS](https://picocss.com/) handles all the styling with zero custom classes. It's a classless framework, so semantic HTML just works. An `<article>` becomes a card, a `<header>` inside it gets a darker background. No class names needed.

---

## The Workflow

Writing a post:

```bash
# Create a new post
go run . new "Post Title"

# Edit the generated markdown file in content/posts/

# Start dev server with live reload
go run . serve

# Build for production
go run . build
```

The dev server watches for file changes and auto-reloads the browser. Edit markdown, save, see the result instantly.

Deployment is just `git push`. [GitHub Actions](https://github.com/features/actions) runs the build and publishes to [GitHub Pages](https://pages.github.com/). No manual steps, no deployment configuration. Write, commit, push, done.

---

## The Results

The site performs well.

> **[Lighthouse](https://developer.chrome.com/docs/lighthouse/) scores: 100/100** across Performance, Accessibility, Best Practices, and SEO.

No JavaScript means no bundle to download, parse, or execute. Pages load quickly, even on slow connections.

[Pico CSS](https://picocss.com/) adds minimal overhead, and the semantic HTML means good accessibility by default. The auto-generated sitemap and RSS feed handle SEO and discoverability.

It's exactly what I wanted: fast, simple, and it just works.

---

## How I Built It

I used AI-assisted development for this project. Started with [Gemini](https://gemini.google.com/) to devise the original plan (tech stack, requirements, big picture architecture). Then [Claude](https://claude.ai/) with [Superpowers](https://github.com/superpowers-marketplace/superpowers) for detailed design and planning. [Codex](https://github.com/cognition-labs/codex) handled the implementation.

The workflow: brainstorm with Gemini, plan with Claude, implement with Codex.

---

## Closing

That's the site. No frameworks, no databases, no complexity. Just markdown files that become fast HTML pages.

It's exactly what a dev blog should be: simple, fast, and mine.
