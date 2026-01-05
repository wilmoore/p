# 002. Zero-configuration tmux execution

Date: 2026-01-04

## Status

Accepted

## Context

tmux supports extensive configuration via `~/.tmux.conf`. However, user configurations can vary wildly and cause unexpected behavior. The goal is to provide consistent, predictable behavior across all environments.

## Decision

Always invoke tmux with `-f /dev/null` to disable reading of any configuration file. Inject only minimal ergonomic defaults at runtime (vi-style copy mode bindings).

## Consequences

**Positive:**
- Completely predictable behavior across all machines and users
- No config drift - what works on one machine works everywhere
- Clean slate - no need to debug user config interactions
- Fast startup - no config file parsing

**Negative:**
- Users cannot customize tmux behavior via their dotfiles when using `p`
- Some users may miss their preferred keybindings or themes
- Requires runtime injection of any desired defaults

## Alternatives Considered

1. **Read user config** - Rejected because it introduces config drift and unpredictable behavior.

2. **Provide a `p`-specific config file** - Rejected because it still requires maintaining a config file.

3. **No config injection at all** - Rejected because vi-style copy mode is a proven ergonomic necessity.

## Related

- Planning: `.plan/.done/mvp-p/`
- PRD Section 2: Core Design Principles (point 1 and 2)
