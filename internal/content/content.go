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
