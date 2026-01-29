# Post Reading Card Design

**Goal:** Make individual post pages visually consistent with the app-like cards while improving brand hierarchy (site title stands out without overpowering content).

## Summary
We will introduce a “reading card” layout for individual posts: a single, wide panel with gentle glow, generous padding, and improved typography hierarchy. We’ll also rebalance header hierarchy by slightly reducing post title size/weight and slightly boosting the site title size/weight/contrast.

## Requirements
- Post pages use a single reading card (not list item cards).
- Maintain accessibility: clear focus-visible styles and readable contrast.
- Preserve reduced-motion support.
- Rebalance hierarchy: site title slightly stronger; post title slightly calmer.

## Non-goals
- No new content metadata (read time, tags).
- No layout overhaul across the site.
- No changes to navigation structure.

## Structure and Semantics
- Keep the post content within an `<article>` for semantics.
- Wrap the content in a `div` or `section` with a `post-reading-card` class to apply the visual panel while avoiding Pico’s default article styling.
- Keep headings and content structure unchanged for markdown rendering.

## Styling Approach
- `post-reading-card` should be wider than list cards, with larger padding and a subtler, calmer glow.
- Background slightly lighter than page, with a low-opacity outer glow (box-shadow) and no inner border.
- Reduce post title size and weight slightly (e.g., 2.1–2.3rem, weight 600).
- Increase site title size a touch (e.g., +0.1–0.15rem) and improve contrast in gradient.
- Keep underline animation for inline links in body content.

## Accessibility
- No nested links or interactive conflicts.
- Focus-visible styles preserved for cards and links.

## Testing
- Extend template tests to check for the new `post-reading-card` wrapper in post templates.
- Extend CSS tests to verify new selectors for post reading card and adjusted title styles.

