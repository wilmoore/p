# 004. Stack-based Drill-down Navigation

Date: 2026-01-07

## Status

Accepted

## Context

Users need to navigate into subdirectories before creating a tmux session. The original design only showed top-level CDPATH directories, but many users organize projects in nested structures (e.g., `clients/acme/backend`).

We needed a way to:
- Allow unlimited navigation depth into subdirectories
- Provide a way to go back up the directory tree
- Maintain context of the navigation path
- Keep the UI simple and predictable

## Decision

We implemented a navigation stack approach:

1. **Navigation Stack**: A slice of `*Directory` pointers tracks the current path
2. **Drill-down**: When user selects a directory with subdirectories, they can choose to "drill down", pushing the directory onto the stack
3. **Back Navigation**: Typing `..` pops the stack and returns to the parent level
4. **Visual Indicator**: Directories with subdirectories show a `>` suffix in the listing
5. **Context Display**: When drilled down, the current directory name is shown as a header

The UI flow prompts users when selecting a directory with subdirectories:
- `[d]` drill down into subdirectories
- `[c]` create session at current directory
- `[q]` cancel and return to selection

## Consequences

**Positive:**
- Unlimited navigation depth without UI complexity
- Familiar `..` pattern for back navigation
- Clear visual feedback on drillable directories
- Stack approach is simple to implement and maintain

**Negative:**
- Additional prompt step for directories with subdirectories
- Slightly longer interaction for deep navigation
- Sessions at top level show full sessions list; drilled-down views do not

## Alternatives Considered

1. **Recursive directory listing**: Show all nested directories at once
   - Rejected: Would create overwhelming lists and violate "no filesystem scanning" principle

2. **Breadcrumb-style path input**: Let users type paths directly
   - Rejected: Adds cognitive load, less discoverable than selection

3. **Two-column layout**: Show parent and children simultaneously
   - Rejected: Over-complicates the minimal TUI design

## Related

- Planning: `.plan/.done/feat-drill-down-and-status-bar/`
- Spec: US-008 (Interactive subdirectory navigation)
