# Screenshots

This folder contains UI screenshots for documentation.

## Required Screenshots

1. **dashboard.png** - Main dashboard view
2. **tasks.png** - Task management interface
3. **projects.png** - Project overview
4. **team.png** - Team members view
5. **activity.png** - Activity timeline
6. **login.png** - Login page
7. **register.png** - Registration page

## Taking Screenshots

### Desktop (1920x1080)
```bash
# Use browser dev tools to capture
# Or use screenshot tool:
gnome-screenshot -d 3 -f dashboard.png
```

### Mobile (375x667)
```bash
# Use Chrome DevTools
# Device: iPhone SE
```

## Image Guidelines

- **Format:** PNG (for quality)
- **Resolution:** 1920x1080 (desktop), 375x667 (mobile)
- **File size:** Optimize with `optipng` or `pngquant`
- **Content:** Use realistic demo data

## Optimization

```bash
# Optimize PNG files
optipng *.png

# Or compress
pngquant --quality=65-80 *.png
```

---

After adding screenshots, update README.md image references.
