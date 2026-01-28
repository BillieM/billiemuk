# Site Refresh Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Refresh billiem.uk with gradient name styling, icon social links, hover effects, and accessibility improvements.

**Architecture:** Update Go config for correct identity/URLs, modify base template to remove footer and use SVG icons, extend theme.css with gradient text, transitions, and hover effects.

**Tech Stack:** Go templates, CSS custom properties, inline SVG icons

---

## Task 1: Update Site Identity

**Files:**
- Modify: `main.go:54-65`

**Step 1: Update siteConfig function**

Change the site title and social URLs:

```go
func siteConfig() templates.SiteData {
	return templates.SiteData{
		Title:   "billiem",
		BaseURL: "https://billiem.uk",
		Author:  "billiem",
		Year:    time.Now().Year(),
		Socials: []templates.Social{
			{Name: "GitHub", URL: "https://github.com/billiem"},
			{Name: "LinkedIn", URL: "https://www.linkedin.com/in/billie-merz-53054418b/"},
		},
	}
}
```

**Step 2: Run tests to verify nothing breaks**

Run: `cd /Users/billie/Code/billiemuk && go test ./...`
Expected: All tests pass

**Step 3: Commit**

```bash
git add main.go
git commit -m "fix: update site title and social links"
```

---

## Task 2: Remove Footer from Template

**Files:**
- Modify: `templates/base.html:28-30`

**Step 1: Remove footer element**

Delete these lines from base.html:

```html
    <footer class="container">
        <small>&copy; {{.Site.Year}} {{.Site.Title}}</small>
    </footer>
```

The closing `</main>` tag should now be followed directly by the DevMode script block.

**Step 2: Run tests to verify template still renders**

Run: `cd /Users/billie/Code/billiemuk && go test ./internal/templates/...`
Expected: All tests pass

**Step 3: Verify visually**

Run: `cd /Users/billie/Code/billiemuk && go run . serve`
Check: http://localhost:8080 - footer should be gone

**Step 4: Commit**

```bash
git add templates/base.html
git commit -m "chore: remove footer"
```

---

## Task 3: Replace Social Text with SVG Icons

**Files:**
- Modify: `templates/base.html:19-21`

**Step 1: Replace social links with SVG icons**

Replace the social links loop:

```html
                {{range .Site.Socials}}
                <li><a href="{{.URL}}" target="_blank" rel="noopener noreferrer">{{.Name}}</a></li>
                {{end}}
```

With icon-specific rendering:

```html
                {{range .Site.Socials}}
                <li><a href="{{.URL}}" target="_blank" rel="noopener noreferrer" aria-label="{{.Name}}" class="social-icon">
                    {{if eq .Name "GitHub"}}<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="currentColor"><path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/></svg>{{end}}
                    {{if eq .Name "LinkedIn"}}<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="currentColor"><path d="M19 0h-14c-2.761 0-5 2.239-5 5v14c0 2.761 2.239 5 5 5h14c2.762 0 5-2.239 5-5v-14c0-2.761-2.238-5-5-5zm-11 19h-3v-11h3v11zm-1.5-12.268c-.966 0-1.75-.79-1.75-1.764s.784-1.764 1.75-1.764 1.75.79 1.75 1.764-.783 1.764-1.75 1.764zm13.5 12.268h-3v-5.604c0-3.368-4-3.113-4 0v5.604h-3v-11h3v1.765c1.396-2.586 7-2.777 7 2.476v6.759z"/></svg>{{end}}
                </a></li>
                {{end}}
```

**Step 2: Run tests**

Run: `cd /Users/billie/Code/billiemuk && go test ./internal/templates/...`
Expected: Tests pass (they check for "GitHub" which is now in aria-label)

**Step 3: Verify visually**

Run: `cd /Users/billie/Code/billiemuk && go run . serve`
Check: http://localhost:8080 - should see GitHub and LinkedIn icons

**Step 4: Commit**

```bash
git add templates/base.html
git commit -m "feat: replace social text links with SVG icons"
```

---

## Task 4: Add Gradient Name Styling

**Files:**
- Modify: `templates/base.html:16`
- Modify: `static/css/theme.css`

**Step 1: Add class to site title link**

In base.html, change line 16 from:

```html
                <li><strong><a href="/">{{.Site.Title}}</a></strong></li>
```

To:

```html
                <li><strong><a href="/" class="site-title">{{.Site.Title}}</a></strong></li>
```

**Step 2: Add gradient styles to theme.css**

Add at the end of theme.css:

```css
/* Site title gradient */
.site-title {
  background: linear-gradient(135deg, var(--pico-primary) 0%, #cba6f7 50%, #f5c2e7 100%);
  -webkit-background-clip: text;
  background-clip: text;
  color: transparent;
  font-size: 1.25rem;
  text-decoration: none;
}

.site-title:hover {
  text-decoration: none;
}
```

**Step 3: Add dark mode gradient colors**

Inside the dark mode media query (after line 105), add:

```css
  .site-title {
    background: linear-gradient(135deg, var(--pico-primary) 0%, #cba6f7 50%, #f5c2e7 100%);
    -webkit-background-clip: text;
    background-clip: text;
    color: transparent;
  }
```

Note: The Catppuccin Mocha mauve (#cba6f7) and pink (#f5c2e7) work in both modes.

**Step 4: Verify visually**

Run: `cd /Users/billie/Code/billiemuk && go run . serve`
Check: http://localhost:8080 - "billiem" should have gradient text in both light and dark modes

**Step 5: Commit**

```bash
git add templates/base.html static/css/theme.css
git commit -m "feat: add gradient styling to site title"
```

---

## Task 5: Add Global Transitions

**Files:**
- Modify: `static/css/theme.css`

**Step 1: Add transition styles**

Add at the end of theme.css:

```css
/* Global transitions */
a,
button,
.social-icon svg,
article {
  transition: all 0.2s ease;
}

/* Respect reduced motion preference */
@media (prefers-reduced-motion: reduce) {
  *,
  *::before,
  *::after {
    transition: none !important;
    animation: none !important;
  }
}
```

**Step 2: Verify visually**

Run: `cd /Users/billie/Code/billiemuk && go run . serve`
Check: Hover over links - transitions should be smooth

**Step 3: Commit**

```bash
git add static/css/theme.css
git commit -m "feat: add global transitions with reduced-motion support"
```

---

## Task 6: Add Social Icon Hover Effects

**Files:**
- Modify: `static/css/theme.css`

**Step 1: Add social icon hover styles**

Add at the end of theme.css:

```css
/* Social icon hover effects */
.social-icon {
  display: inline-flex;
  align-items: center;
  color: var(--pico-muted-color);
}

.social-icon:hover {
  color: var(--pico-primary);
  transform: scale(1.1);
  filter: drop-shadow(0 0 8px var(--pico-primary-focus));
}

.social-icon svg {
  display: block;
}
```

**Step 2: Verify visually**

Run: `cd /Users/billie/Code/billiemuk && go run . serve`
Check: Hover over social icons - should scale up and glow

**Step 3: Commit**

```bash
git add static/css/theme.css
git commit -m "feat: add social icon hover effects"
```

---

## Task 7: Add Link Hover Underline Effect

**Files:**
- Modify: `static/css/theme.css`

**Step 1: Add link underline styles**

Add at the end of theme.css:

```css
/* Link hover underline effect */
main a:not(.social-icon) {
  text-decoration: none;
  position: relative;
}

main a:not(.social-icon)::after {
  content: '';
  position: absolute;
  bottom: 0;
  left: 0;
  width: 0;
  height: 1px;
  background-color: var(--pico-primary);
  transition: width 0.2s ease;
}

main a:not(.social-icon):hover::after {
  width: 100%;
}

main a:not(.social-icon):hover {
  color: var(--pico-primary-hover);
}
```

**Step 2: Verify visually**

Run: `cd /Users/billie/Code/billiemuk && go run . serve`
Check: Hover over links in main content - underline should slide in

**Step 3: Commit**

```bash
git add static/css/theme.css
git commit -m "feat: add link hover underline animation"
```

---

## Task 8: Add Post Card Hover Effects

**Files:**
- Modify: `static/css/theme.css`
- Modify: `templates/home.html`

**Step 1: Add card class to article in home.html**

Change the article tag (line 4):

```html
    <article>
```

To:

```html
    <article class="post-card">
```

**Step 2: Add post card hover styles to theme.css**

Add at the end of theme.css:

```css
/* Post card hover effects */
.post-card {
  padding: 1rem;
  border-radius: 0.5rem;
  margin-bottom: 1rem;
}

.post-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px var(--pico-primary-focus);
}
```

**Step 3: Verify visually**

Run: `cd /Users/billie/Code/billiemuk && go run . serve`
Check: If posts exist, hover over them - should lift and show shadow

**Step 4: Commit**

```bash
git add static/css/theme.css templates/home.html
git commit -m "feat: add post card hover effects"
```

---

## Task 9: Final Verification

**Step 1: Run all tests**

Run: `cd /Users/billie/Code/billiemuk && go test ./...`
Expected: All tests pass

**Step 2: Build production site**

Run: `cd /Users/billie/Code/billiemuk && go run . build`
Expected: Build complete without errors

**Step 3: Visual verification checklist**

Run: `cd /Users/billie/Code/billiemuk && go run . serve`

Check at http://localhost:8080:
- [ ] Site title shows "billiem" with gradient
- [ ] Social icons display (no text)
- [ ] Social icons have hover glow/scale
- [ ] No footer present
- [ ] Link hovers show sliding underline
- [ ] Transitions are smooth (0.2s)
- [ ] Dark mode: gradient still visible
- [ ] Dark mode: all hover effects work

**Step 4: Accessibility check**

- [ ] Social icons have aria-labels (inspect element)
- [ ] Focus states visible on keyboard navigation (tab through page)
- [ ] Test with `prefers-reduced-motion` (browser dev tools)

**Step 5: Final commit if any fixes needed**

```bash
git add -A
git commit -m "fix: final polish and fixes"
```

---

## Summary

| Task | Description |
|------|-------------|
| 1 | Update site identity (title, social URLs) |
| 2 | Remove footer |
| 3 | Replace social text with SVG icons |
| 4 | Add gradient name styling |
| 5 | Add global transitions |
| 6 | Add social icon hover effects |
| 7 | Add link hover underline effect |
| 8 | Add post card hover effects |
| 9 | Final verification |
