# Architecture Decision Records

This directory contains Architecture Decision Records (ADRs) documenting significant technical decisions.

## What is an ADR?

An ADR captures the context, decision, and consequences of an architecturally significant choice.

## Format

We use the [Michael Nygard format](https://cognitect.com/blog/2011/11/15/documenting-architecture-decisions).

## Naming Convention

- Filename: `NNN-kebab-case-title.md` (e.g., `001-use-cdpath-for-discovery.md`)
- NNN = zero-padded sequence number (001, 002, 003...)
- Title in heading must match: `# NNN. Title` (e.g., `# 001. Use CDPATH for Discovery`)

## Index

- [001. Use CDPATH as sole discovery mechanism](001-use-cdpath-for-discovery.md)
- [002. Zero-configuration tmux execution](002-zero-config-tmux-execution.md)
- [003. Collision-free session naming with hash suffix](003-collision-free-session-naming.md)
- [004. Stack-based drill-down navigation](004-stack-based-drill-down-navigation.md)
- [005. Color-agnostic status bar styling](005-color-agnostic-status-bar.md)
