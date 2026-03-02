# 010. Apply Global Window Style for Consistent Backgrounds

Date: 2026-03-02

## Status

Accepted

## Context

Per ADR-005, `p` enforces a color-agnostic status bar with a black pane background. In practice only the first window inherited `window-style` because we targeted the session (`set-option -t session`) even though `window-style` and `window-active-style` are window-level options. Newly created windows therefore reverted to tmux defaults (grey), producing inconsistent visuals once additional windows were added manually or by `P_WINDOWS`.

## Decision

Set `window-style` and `window-active-style` using global (`-g`) options immediately after the session is created. Global options become the template for all future windows in the server, guaranteeing every pane inherits the black background regardless of when or how the window was created.

## Consequences

- Positive: The background color now matches ADR-005 across every window, ensuring a cohesive UI.
- Positive: Default windows created later inherit styles automatically without per-window overrides.
- Negative: Global options apply to every subsequent session launched from the same server while `p` is running, even if launched manually. This is acceptable because `p` already runs tmux with `-f /dev/null` to isolate its environment.

## Alternatives Considered

- **Continue targeting the session** — rejected because it requires reapplying styles for every new window and failed to solve the observed bug.
- **Explicitly set styles on each new window** — rejected for added complexity and fragility; it would require hooks or repeated commands every time `new-window` runs.
- **Switch to pane borders with color fills** — rejected because it diverges from ADR-005 and does not solve the underlying inheritance issue.

## Related

- Planning: `.plan/.done/feat-default-windows-and-consistent-background/`
