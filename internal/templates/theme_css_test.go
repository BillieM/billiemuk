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
		".site-title",
		".social-icon",
		".post-card",
		".post-card__title",
		".post-card__meta",
		".post-card__summary",
		".post-card:focus-visible",
		"main a.post-card",
		".post-card-wrap",
		"0 0 0 1px",
		"prefers-reduced-motion: reduce",
		"main a:not(.social-icon)::after",
	}

	for _, check := range checks {
		if !strings.Contains(contents, check) {
			t.Errorf("theme.css missing %q", check)
		}
	}
}
