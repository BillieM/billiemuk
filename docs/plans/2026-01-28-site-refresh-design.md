# Site Refresh Design

## Overview

Refresh the billiem.uk personal site to feel more modern and alive while maintaining simplicity. The approach: one bold gradient element (the name), subtle interactive feedback everywhere else, and clean foundations.

## Changes

### Identity & Branding

- **Site title**: Change from "Billie Muk" to "billiem"
- **Name styling**: Gradient text using Catppuccin palette (lavender → mauve range)
- **Typography**: Slightly larger/bolder name to create hierarchy against navigation

### Social Links

- **GitHub**: Update to github.com/billiem
- **LinkedIn**: Update to linkedin.com/in/billie-merz-53054418b/
- **Display**: SVG icons instead of text labels
- **Position**: Right side of header (unchanged)
- **Hover effect**: Scale up (1.05x) + soft solid glow in accent color

### Footer

- Remove entirely (no copyright, no content)

### Link Styling

- **Default**: Accent color (lavender)
- **Hover**: Underline slides in + color brightens slightly
- **Transition**: 0.2s ease on all interactive elements

### Post Cards

- **Hover**: Subtle lift (translate up 2-4px) + soft shadow appears
- **Transition**: 0.2s ease

### Global Polish

- All interactive elements get `transition: 0.2s ease`
- Smooth color transitions on theme toggle (light/dark)
- Gradient on name works in both light and dark modes

## Accessibility Requirements

1. **Icon links**: Add `aria-label` to social icon links ("GitHub", "LinkedIn")
2. **Reduced motion**: Respect `prefers-reduced-motion` - disable transitions/animations for users who prefer reduced motion
3. **Gradient contrast**: Ensure gradient text colors meet WCAG contrast ratios against backgrounds
4. **Focus states**: Verify focus indicators remain visible on all interactive elements

## Technical Notes

### Gradient Implementation

Use CSS `background: linear-gradient()` with `background-clip: text` and `color: transparent` for the name gradient. Define gradient colors for both light and dark modes.

### Icon Implementation

Use inline SVG for GitHub and LinkedIn icons. This allows:
- Styling with CSS (color, hover effects)
- No external requests
- Crisp rendering at any size

### Reduced Motion

```css
@media (prefers-reduced-motion: reduce) {
  * {
    transition: none !important;
    animation: none !important;
  }
}
```

## Files to Modify

- `main.go` - Update SiteData (title, social URLs)
- `templates/base.html` - Remove footer, update social links to use icons
- `static/css/theme.css` - Add gradient styles, transitions, hover effects
- `internal/templates/types.go` - May need to add icon field to Social struct (or handle in template)

## Out of Scope

- Skip-to-content link (low priority for simple blog)
- Blinking cursor effect (decided against)
- Gradient backgrounds (keeping backgrounds solid)
- Gradient hover effects (keeping hovers solid to let name be the hero)
