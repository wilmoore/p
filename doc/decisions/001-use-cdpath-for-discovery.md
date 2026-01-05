# 001. Use CDPATH as sole discovery mechanism

Date: 2026-01-04

## Status

Accepted

## Context

`p` needs a way to discover project directories. The goal is zero-configuration operation where behavior is deterministic and identical across machines.

Common approaches for project discovery include:
- Filesystem scanning (find all git repos, search home directory)
- Configuration files (list projects in a config file)
- Environment variables (use existing shell conventions)

## Decision

Use CDPATH as the sole mechanism for discovering project directories. Only immediate children of CDPATH directories are considered projects. No recursive scanning.

## Consequences

**Positive:**
- Zero configuration required - CDPATH is already set by many developers
- Explicit over clever - user controls exactly which directories are discoverable
- Identical behavior across machines (given same CDPATH)
- No surprise discoveries or slow filesystem scans
- Respects existing Unix conventions

**Negative:**
- Requires user to have CDPATH configured
- Only works with the immediate children model (no nested discovery)
- Users unfamiliar with CDPATH need to learn this shell feature

## Alternatives Considered

1. **Filesystem scanning** - Rejected because it's slow, non-deterministic, and discovers things users didn't intend to expose.

2. **Configuration file** - Rejected because it introduces config drift and requires manual maintenance.

3. **Git repository detection** - Rejected because not all projects are git repos, and it requires scanning.

## Related

- Planning: `.plan/.done/mvp-p/`
- PRD Section 2: Core Design Principles
