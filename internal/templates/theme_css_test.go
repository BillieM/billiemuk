package templates

import (
	"os"
	"strings"
	"testing"
)

func TestThemeCSSIncludesRefreshStyles(t *testing.T) {
	css, err := os.ReadFile("../../static/css/theme.css")
	if err != nil {
		t.Fatalf("read theme.css: %v", err)
	}

	contents := string(css)
	checks := []string{
		"body > header",
		"body > main",
		"header nav",
		"main article",
		"main article header",
		"main article a",
		"body > header nav ul:first-child a",
		"0 0 0 1px",
		"prefers-reduced-motion: reduce",
	}

	for _, check := range checks {
		if !strings.Contains(contents, check) {
			t.Errorf("theme.css missing %q", check)
		}
	}
}
