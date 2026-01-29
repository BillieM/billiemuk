# Site Refresh App-like Cards Design

**Goal:** Make post cards feel more app-like and clickable while preserving calm readability and strong accessibility.

## Summary
We will turn each post card into a single, full-card link with clear focus-visible styling, subtle lift/glow hover, and improved typography hierarchy (title > date > summary). We will keep the overall dark, developer-chic aesthetic and avoid noisy hover underlines in cards.

## Requirements
- Entire card is one link (title, date, summary).
- No nested links inside cards.
- Clear focus-visible styling for keyboard and Vimium users.
- Calm hover: slight lift, soft glow, minor background lightening.
- Improve hierarchy: title stronger, date smaller and muted.
- Keep reduced-motion support.

## Non-goals
- No major layout or typography system overhaul.
- No new data fields (e.g., read time).
- No additional pages or components.

## Structure and Semantics
- Keep `<article>` for semantics.
- Move clickable area to a single `<a class="post-card">` inside the article.
- Use sub-elements for hierarchy:
  - `.post-card__title`
  - `.post-card__meta` (date)
  - `.post-card__summary`

## Styling Approach
- `.post-card` becomes the primary interactive surface.
- Default state: slightly lighter background than page, subtle 1px border.
- Hover/focus: translateY(-2px), stronger shadow, border glow, title color shift.
- Remove underline hover for titles in cards.
- Keep inline link underline animation for main content outside cards.

## Accessibility
- Single focusable link per card.
- Clear `:focus-visible` style (outline or box-shadow).
- Maintain reduced-motion media query to disable transforms.

## Testing
- Update template tests to assert new `.post-card` link structure and class names.
- Update CSS tests to check new selectors for card and focus-visible styles.

