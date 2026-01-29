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
			Socials: []Social{
				{Name: "GitHub", URL: "https://github.com/test"},
				{Name: "LinkedIn", URL: "https://www.linkedin.com/in/test"},
			},
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
		"/posts/2026-01-15-first-post/",
		"<header",
		"<main",
		"<article",
		"<svg",
		`aria-label="GitHub"`,
		`aria-label="LinkedIn"`,
		"theme.min.css",
	}
	for _, check := range checks {
		if !strings.Contains(html, check) {
			t.Errorf("home HTML missing %q", check)
		}
	}

	forbidden := []string{
		`class="container"`,
		`class="site-title"`,
		`class="social-icon"`,
		`class="post-card-wrap"`,
		`class="post-card"`,
		`class="post-card__title"`,
		`class="post-card__meta"`,
		`class="post-card__summary"`,
		"pico.min.css",
	}
	for _, check := range forbidden {
		if strings.Contains(html, check) {
			t.Errorf("home HTML should not include %q", check)
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
		"og:title",
		"og:type",
		"canonical",
		"<article",
		"<header",
	}
	for _, check := range checks {
		if !strings.Contains(html, check) {
			t.Errorf("post HTML missing %q", check)
		}
	}
}
