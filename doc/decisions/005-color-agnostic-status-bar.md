# 005. Color-agnostic Status Bar Styling

Date: 2026-01-07

## Status

Accepted

## Context

The status bar is the primary visual element in tmux. We needed to style it to look polished while adhering to the "zero configuration" principle. The challenge: users have wildly different terminal color schemes (light themes, dark themes, custom palettes).

Requirements:
- Look good on any terminal theme
- No user configuration needed
- Keep it minimal and unobtrusive
- Show useful information (git branch)

## Decision

We chose a "color-agnostic" approach using:

1. **Default Background**: `bg=default` inherits the terminal's background color
2. **Muted Gray Foreground**: `colour240` (medium gray) for secondary elements
3. **Default Foreground**: `fg=default` with `bold` for the current window indicator
4. **Minimal Content**:
   - Left: Session name (muted)
   - Center: Window list (current window bold)
   - Right: Git branch from `pane_current_path` (muted)

The specific styling:
```
status-style: bg=default,fg=default
status-left-style: fg=colour240
status-right-style: fg=colour240
window-status-style: fg=colour240
window-status-current-style: fg=default,bold
```

Git branch is fetched via: `git -C #{pane_current_path} rev-parse --abbrev-ref HEAD 2>/dev/null`

## Consequences

**Positive:**
- Works universally across terminal themes
- Gray text is subtle but readable on both light and dark backgrounds
- Bold current window stands out without jarring colors
- No color clashes with user's terminal palette

**Negative:**
- Less visually distinctive than colored status bars
- `colour240` may have slightly different appearance across terminals
- Git branch command adds minor overhead (mitigated by 5s refresh interval)

## Alternatives Considered

1. **Bright accent colors**: Use green/blue for git branch
   - Rejected: Clashes with many terminal themes, especially light themes

2. **Terminal color names**: Use `red`, `green`, `blue`
   - Rejected: These map to user's palette and could clash unpredictably

3. **No styling at all**: Leave stock tmux defaults
   - Rejected: Stock green bar is visually jarring and provides poor UX

4. **Detect light/dark and switch**: Runtime theme detection
   - Rejected: Adds complexity, unreliable detection, violates simplicity principle

## Related

- Planning: `.plan/.done/feat-drill-down-and-status-bar/`
- Backlog item #14 (Default window name and status bar styling)
