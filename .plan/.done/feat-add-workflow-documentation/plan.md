# Implementation Plan: Workflow Documentation

## Summary

Add workflow documentation to help users understand idiomatic p usage patterns.

## Changes

### 1. README.md - Add Workflows Section

Add a "Workflows" section in the Usage area with:
- `mkdir something-new && p !$` - Create new project workflow
- `p .` - Create session in current directory
- Quick filtering technique

### 2. doc/examples.md - New File

Create new file with detailed examples:
- New project workflow (mkdir + p !$)
- Current directory workflow (p .)
- Resume work (p + fuzzy filter)
- Rapid context switching
- Running p inside tmux
- Include "Why this works" explanations

## Files Modified

- README.md
- doc/examples.md (new)

## Backlog

Item 27: Add workflow documentation
