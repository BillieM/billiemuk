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
