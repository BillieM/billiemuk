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
		"prefers-reduced-motion: reduce",
		"--pico-background-color: #eff1f5",
		"--pico-background-color: #1e1e2e",
		"--pico-card-sectioning-background-color: #dce0e8",
		"--pico-card-sectioning-background-color: #11111b",
	}

	for _, check := range checks {
		if !strings.Contains(contents, check) {
			t.Errorf("theme.css missing %q", check)
		}
	}
}
