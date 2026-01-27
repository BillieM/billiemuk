package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"billiemuk/internal/builder"
	"billiemuk/internal/content"
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
