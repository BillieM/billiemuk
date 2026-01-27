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
	ContentDir    string
	TemplatesDir  string
	StaticDir     string
	DistDir       string
	Site          templates.SiteData
	IncludeDrafts bool
	DevMode       bool
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
