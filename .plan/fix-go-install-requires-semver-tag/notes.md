# Bug Fix: go install @latest fails - no semver tag

## Bug Details

**Reported:** 2026-01-10
**Branch:** fix/go-install-requires-semver-tag
**Severity:** Critical (blocks new user installation)

### Steps to Reproduce

```bash
go install github.com/wilmoore/p@latest
which p
# Output: p not found
```

### Expected Behavior

The `p` binary should be installed to `$GOPATH/bin` (or `$GOBIN`) and be available in PATH.

### Actual Behavior

- `go install` appears to succeed silently
- `which p` returns "p not found"
- pkg.go.dev shows "Oops! We couldn't find github.com/wilmoore/p"

### Environment

- Go 1.22.5 (via asdf)
- macOS Darwin 25.1.0
- Repository: github.com/wilmoore/p (public)
- GOPATH: /Users/wilmooreiii/.asdf/installs/golang/1.22.5/packages

## Root Cause Analysis

### Investigation

1. **Checked git tags** - `git tag -l` returned empty (no tags exist)
2. **Verified repo is public** - `gh repo view` confirmed `isPrivate: false`
3. **Checked if binary exists** - Found at `/Users/wilmooreiii/.asdf/installs/golang/1.22.5/packages/bin/p`
4. **Tested binary directly** - Works correctly when invoked with full path

### Root Cause

**Primary Issue:** No semantic version tags exist on the repository.

Go's module system requires at least one `vX.Y.Z` tag for `@latest` resolution. Without tags:
- Go cannot resolve `@latest` to a specific version
- The module proxy (proxy.golang.org) doesn't index the module
- pkg.go.dev cannot display documentation

**Secondary Issue (unrelated):** The GOPATH/bin directory is not in the user's shell PATH, but this is a user environment issue, not a project issue.

## Fix Implementation

### Solution

Create and push a semantic version tag `v0.1.0` to enable:
1. `go install github.com/wilmoore/p@latest` to work
2. Module indexing on proxy.golang.org
3. Documentation on pkg.go.dev

### Commands

```bash
# Create annotated tag on main branch
git checkout main
git tag -a v0.1.0 -m "Initial release"

# Push tag to remote
git push origin v0.1.0
```

### Verification

After pushing the tag:
1. Wait ~1-2 minutes for proxy.golang.org to index
2. Run: `go install github.com/wilmoore/p@latest`
3. Verify: `which p` or `$GOPATH/bin/p --help`
4. Check: https://pkg.go.dev/github.com/wilmoore/p

## Related Issues

- Users must have `$GOPATH/bin` or `$GOBIN` in their PATH (documented in README)
