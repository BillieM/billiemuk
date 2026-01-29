# Pico Native Cleanup Design

**Goal:** Remove all custom CSS overrides and use Pico's built-in classless features for maximum consistency and proper Pico patterns.

## Summary

We will remove all custom article/card styling and let Pico's built-in classless semantics handle everything. The only custom CSS will be for elements Pico doesn't style: gradient site title and SVG icon sizing. We'll restructure the homepage to use Pico-native link patterns (linked headings instead of full-card links).

## Problems Being Solved

1. **Custom article styling breaking Pico's card system** - Custom padding overrides prevent Pico's header negative margins from working, causing misaligned sectioning backgrounds
2. **Inconsistency between homepage and post pages** - Different custom styling rules for list cards vs reading view
3. **Too much custom CSS** - Fighting against Pico's defaults instead of embracing them
4. **Non-standard structure** - Wrapping entire articles in links isn't the Pico way

## Requirements

- Remove all custom article/card styling
- Use Pico's default `<article>` card styling for both homepage and post pages
- Restructure homepage to use linked headings (Pico-native pattern)
- Keep only minimal custom CSS for branding elements Pico doesn't handle
- Maintain Catppuccin color theme
- Preserve semantic HTML structure

## Non-goals

- No new features or content
- No changes to build system
- No JavaScript changes
- No typography overrides

## Template Changes

### Homepage (`templates/home.html`)

**Before:**
```html
<section>
  {{range .Posts}}
  <a href="/posts/{{.Slug}}/">
    <article>
      <header>
        <h2>{{.Title}}</h2>
        <p><time>...</time></p>
      </header>
      {{if .Summary}}<p>{{.Summary}}</p>{{end}}
    </article>
  </a>
  {{end}}
</section>
```

**After:**
```html
<section>
  {{range .Posts}}
  <article>
    <header>
      <h2><a href="/posts/{{.Slug}}/">{{.Title}}</a></h2>
      <p><time datetime="{{.Date.Format "2006-01-02"}}">{{.Date.Format "2 January 2006"}}</time></p>
    </header>
    {{if .Summary}}<p>{{.Summary}}</p>{{end}}
  </article>
  {{end}}
</section>
```

**Changes:**
- Remove `<a>` wrapper around entire article
- Move link to wrap only the `<h2>` heading
- Pico automatically styles the article as a card
- Pico automatically styles the header with darker sectioning background

### Post Page (`templates/post.html`)

**No changes needed** - already uses pure semantic HTML with `<article>` and `<header>`. Pico will style it identically to homepage cards for perfect consistency.

## SCSS Changes

### Final `static/scss/theme.scss`

```scss
@use "@picocss/pico/scss/pico" with (
  $enable-classes: false,
  $enable-semantic-container: true
);

/* Catppuccin overrides (light) */
:root:not([data-theme="dark"]) {
  --pico-background-color: #eff1f5;
  --pico-color: #4c4f69;
  --pico-muted-color: #6c6f85;
  --pico-muted-border-color: #ccd0da;
  --pico-primary: #7287fd;
  --pico-primary-background: #7287fd;
  --pico-primary-border: #7287fd;
  --pico-primary-hover: #5c6ef0;
  --pico-primary-hover-background: #5c6ef0;
  --pico-primary-hover-border: #5c6ef0;
  --pico-primary-focus: rgba(114, 135, 253, 0.25);
  --pico-primary-inverse: #eff1f5;
  --pico-card-background-color: #e6e9ef;
  --pico-card-border-color: #ccd0da;
  --pico-card-sectioning-background-color: #dce0e8;
  --pico-code-background-color: #e6e9ef;
  --pico-code-color: #4c4f69;
  --pico-blockquote-border-color: #7287fd;
  --pico-blockquote-footer-color: #6c6f85;
  --pico-text-selection-color: rgba(114, 135, 253, 0.2);
}

/* Catppuccin overrides (dark) */
@media (prefers-color-scheme: dark) {
  :root:not([data-theme]) {
    --pico-background-color: #1e1e2e;
    --pico-color: #cdd6f4;
    --pico-muted-color: #a6adc8;
    --pico-muted-border-color: #313244;
    --pico-primary: #b4befe;
    --pico-primary-background: #b4befe;
    --pico-primary-border: #b4befe;
    --pico-primary-hover: #c8d0fe;
    --pico-primary-hover-background: #c8d0fe;
    --pico-primary-hover-border: #c8d0fe;
    --pico-primary-focus: rgba(180, 190, 254, 0.25);
    --pico-primary-inverse: #1e1e2e;
    --pico-card-background-color: #181825;
    --pico-card-border-color: #313244;
    --pico-card-sectioning-background-color: #11111b;
    --pico-code-background-color: #181825;
    --pico-code-color: #cdd6f4;
    --pico-blockquote-border-color: #b4befe;
    --pico-blockquote-footer-color: #a6adc8;
    --pico-text-selection-color: rgba(180, 190, 254, 0.2);
  }
}

/* Minimal custom styling (only what Pico doesn't handle) */
body > header nav ul:first-child a {
  background: linear-gradient(135deg, var(--pico-primary) 0%, #cba6f7 50%, #f5c2e7 100%);
  -webkit-background-clip: text;
  background-clip: text;
  color: transparent;
}

body > header nav ul:last-child a svg {
  width: 20px;
  height: 20px;
}
```

### Removed CSS

All of the following custom selectors are removed:

- `main > section > a` - No longer wrapping articles in links
- `main article` - Pico handles article/card styling
- `main article header` - Pico handles header sectioning
- `main article a` - Not needed
- `main > section > a:hover` - No longer wrapping articles in links
- `main > article` - Pico handles post page articles identically

## What Pico Provides Automatically

In classless mode, Pico automatically handles:

**For `<article>` elements:**
- Padding via `--pico-block-spacing-vertical` and `--pico-block-spacing-horizontal`
- Background via `--pico-card-background-color`
- Border radius via `--pico-border-radius`
- Box shadow via `--pico-card-box-shadow`
- Bottom margin for spacing between cards

**For `<header>` inside `<article>`:**
- Darker background via `--pico-card-sectioning-background-color`
- Negative margins to extend edge-to-edge within the padded card
- Proper padding to align content with article body
- Negative margin-top to pull header to top of card

**For `<a>` elements:**
- Default link styling with proper color, hover, and focus states
- Underline decoration on hover
- Focus-visible outlines for accessibility

## Benefits

1. **Fixes alignment issues** - Pico's header negative margins work correctly without custom padding interference
2. **Perfect consistency** - Homepage cards and post pages use identical styling
3. **Minimal custom CSS** - Only 2 custom rules for truly custom elements
4. **True Pico-native** - Embracing Pico's design system instead of fighting it
5. **Easier maintenance** - Less custom code to maintain
6. **Better accessibility** - Pico's built-in focus states and link patterns

## Testing

- Verify homepage cards display correctly with Pico's default styling
- Verify post pages display identically to homepage cards
- Verify header sectioning background extends edge-to-edge correctly
- Verify link styles work correctly on card headings
- Verify gradient site title displays correctly
- Verify social icons size correctly
- Verify dark mode switches correctly
- Test keyboard navigation and focus states
