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
