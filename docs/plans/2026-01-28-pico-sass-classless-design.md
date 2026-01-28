# Pico Sass Classless Design

**Goal:** Convert the site to Pico’s Sass workflow in classless mode, remove all custom classes, and preserve the existing Catppuccin-inspired styling choices.

## Summary
We will compile Pico via Sass with classless semantics enabled, move our theme customization into Sass (via Pico variables and CSS custom properties), and remove all custom classes from templates. Styling for cards, header, and icons will be expressed using semantic HTML structure and element selectors (no custom classes), to keep consistency with Pico’s philosophy.

## Requirements
- Use Pico Sass compilation with `$enable-classes: false` and `$enable-semantic-container: true`.
- Remove all custom classes from templates (site title, social icons, post cards, reading card, containers).
- Preserve the Catppuccin color scheme.
- Preserve “app-like” card behavior on the home list and reading card on post pages using semantic selectors (not custom classes).
- Keep typography defaults from Pico (no custom font sizes/weights).

## Non-goals
- No new content features (read time, tags, etc.).
- No JavaScript changes.
- No non-Pico typography overrides.

## Architecture
- Add `static/scss/theme.scss` that imports Pico with Sass settings and applies theme overrides.
- Compile Sass into `static/css/theme.css`; the build continues to minify to `theme.min.css`.
- Update templates to classless markup (header/main as semantic containers).
- Replace class-based styling with structural selectors in Sass (e.g., `body > header nav ul:first-child a` for the brand link).

## Testing
- Update template tests to remove class expectations and assert the new semantic structure.
- Update CSS tests to look for the new structural selectors instead of class selectors.
