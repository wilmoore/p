Below is a complete Product Requirements Document generated in **Autopilot mode** from your Idea Pack.
No questions were required. All defaults and inferred decisions are logged explicitly.

---

# Product Requirements Document (PRD)

## Product: `p`

## Version: v1

## Document owner: Product

## Mode: Autopilot

## Confidence floor met for all required decisions

---

## 1. Product Thesis

`p` is a zero-configuration tmux session launcher that treats directories as projects, using CDPATH as the sole discovery mechanism. It removes the need to think about tmux commands, configuration, or session management. The product optimizes for daily personal use, muscle memory, and deterministic behavior.

The core value is frictionless project switching with no configuration drift, no filesystem scanning, and no tmux plugins.

---

## 2. Core Design Principles

1. Zero global configuration
   tmux is always launched with `-f /dev/null`. No user tmux config is read.

2. Minimal injected defaults
   Only ergonomics proven painful without intervention are injected. All other behavior remains stock.

3. Context-free invocation
   The command can be run from any directory, inside or outside tmux.

4. Directories equal projects
   A project is exactly one directory. No metadata, no manifests.

5. Explicit over clever
   No automatic scanning, no inference, no hidden behavior.

---

## 3. Personas

### P-001 Daily Terminal Power User

* Uses tmux daily across many projects
* Frequently switches contexts
* Values speed, predictability, and muscle memory
* Dislikes config drift and plugin ecosystems

### P-002 Minimalist Developer

* Avoids heavy dotfiles and plugins
* Wants identical behavior across machines
* Prefers tools that disappear once learned

---

## 4. Input Scenarios

* User runs `p` from a shell outside tmux
* User runs `p` from inside an existing tmux session
* CDPATH is set with one or more base directories
* CDPATH is unset or empty
* User selects an existing tmux session
* User selects a directory discovered via CDPATH
* User drills into subdirectories before session creation

---

## 5. User Journeys

### J-001 Switch to Existing Session

User wants to attach to or switch to an already running tmux session.

### J-002 Create Session From Directory

User wants to create a new tmux session derived from a project directory.

### J-003 Drill Down Into Directory

User selects a container directory and navigates deeper before creating a session.

### J-004 Failure on Missing CDPATH

User runs `p` without CDPATH configured and receives a clear failure message.

---

## 6. UX Surface Inventory

| Screen ID | Name                | Purpose                                        |
| --------- | ------------------- | ---------------------------------------------- |
| S-001     | Selector            | List sessions and directories, allow selection |
| S-002     | Drill-down Selector | Navigate subdirectories                        |
| S-003     | Failure Message     | Display fatal configuration error              |

---

## 7. Behavior and Editing Model

* Read-only interaction model
* No persistent state
* No user editing of configuration
* All behavior derived from runtime environment, tmux state, and filesystem

---

## 8. Constraints and Anti-Features

### Constraints

* Must be a single static Go binary
* Must depend only on tmux at runtime
* Must not read or write any config files
* Must not scan the filesystem beyond CDPATH resolution

### Anti-Features

* No tmux plugin support
* No layout management
* No dashboards or status views
* No recursive project discovery
* No fallback when CDPATH is unset

---

## 9. Success and Failure Criteria

### Success

* User switches projects without recalling tmux commands
* Session naming is deterministic and collision-free
* Behavior is identical across machines

### Failure

* Any implicit filesystem scanning
* Reading user tmux config
* Non-deterministic session naming
* Silent fallback when CDPATH is missing

---

## 10. North Star Metric

Number of project switches per day performed without invoking a tmux command directly.

---

## 11. Epics

* E-001 [MUST] Session Discovery and Selection
* E-002 [MUST] Session Creation From Directory
* E-003 [MUST] tmux Execution and Injection Model
* E-004 [MUST] Failure Handling and Messaging
* E-005 [SHOULD] Drill-down Navigation

---

## 12. User Stories with Acceptance Criteria

### E-001 Session Discovery and Selection

* US-001 [MUST] As a user, I can see all existing tmux sessions when running `p`.

**Acceptance Criteria**

* Given tmux has active sessions
* When `p` is invoked
* Then all sessions are listed by name
* And no directories are required to be present

---

* US-002 [MUST] As a user, I can attach to or switch to a selected session.

**Acceptance Criteria**

* Given a session is selected
* When the selection is confirmed
* Then tmux attaches to that session
* And if inside tmux, the current client is detached

---

### E-002 Session Creation From Directory

* US-003 [MUST] As a user, I can create a new tmux session from a directory discovered via CDPATH.

**Acceptance Criteria**

* Given CDPATH is set
* And a directory is selected
* When the selection is confirmed
* Then a tmux session is created
* And the session name is derived deterministically

---

* US-004 [MUST] As a user, I never see a naming collision.

**Acceptance Criteria**

* Given two directories normalize to the same base name
* When sessions are created
* Then a stable hash suffix is appended
* And no prompt is shown

---

### E-003 tmux Execution and Injection Model

* US-005 [MUST] As a user, tmux always runs without global config.

**Acceptance Criteria**

* Given any invocation path
* When tmux is launched
* Then `-f /dev/null` is always used

---

* US-006 [MUST] As a user, copy mode behaves with vi-style ergonomics.

**Acceptance Criteria**

* Given copy mode is entered
* When `v` is pressed
* Then selection begins
* And `y` copies and exits copy mode

---

### E-004 Failure Handling and Messaging

* US-007 [MUST] As a user, I receive a clear error when CDPATH is unset.

**Acceptance Criteria**

* Given CDPATH is unset or empty
* When `p` is invoked
* Then no directories are listed
* And execution exits with an error message
* And the message explains how to fix CDPATH

---

### E-005 Drill-down Navigation

* US-008 [SHOULD] As a user, I can navigate into subdirectories interactively.

**Acceptance Criteria**

* Given a selected directory contains subdirectories
* When the user chooses to drill down
* Then subdirectories are listed
* And navigation depth is user-controlled

---

## 13. Traceability Map

| Story  | Epic  | Journey | Screen | Priority |
| ------ | ----- | ------- | ------ | -------- |
| US-001 | E-001 | J-001   | S-001  | MUST     |
| US-002 | E-001 | J-001   | S-001  | MUST     |
| US-003 | E-002 | J-002   | S-001  | MUST     |
| US-004 | E-002 | J-002   | S-001  | MUST     |
| US-005 | E-003 | J-002   | S-001  | MUST     |
| US-006 | E-003 | J-002   | S-001  | MUST     |
| US-007 | E-004 | J-004   | S-003  | MUST     |
| US-008 | E-005 | J-003   | S-002  | SHOULD   |

---

## 14. Lo-fi UI Mockups (ASCII)

### S-001 Selector

Purpose: Select session or directory
Primary actions: Select, confirm
States: empty, loading, error, success

```
p

Sessions:
  [1] work
  [2] notes
  [3] sandbox

Projects:
  [a] dev/api
  [b] clients/acme
  [c] sandbox

> _
```

---

### S-002 Drill-down Selector

Purpose: Navigate into subdirectories

```
clients/acme/

  [1] backend
  [2] frontend
  [3] infra

  [..] up

> _
```

---

### S-003 Failure Message

Purpose: Explicit fatal error

```
CDPATH is not set.

p uses CDPATH to discover project directories.
Set CDPATH to enable directory-based projects.

Exiting.
```

---

## 15. Decision Log

### D-001 Use CDPATH as sole discovery mechanism

* Options: CDPATH only, filesystem scan, config file
* Evidence: Explicit in Idea Pack
* Winner: CDPATH only
* Confidence: 0.95

### D-002 No fallback when CDPATH is missing

* Options: fallback scan, prompt, hard failure
* Evidence: Explicit failure mode defined
* Winner: hard failure
* Confidence: 0.90

### D-003 Inject minimal tmux config programmatically

* Options: tmux.conf, plugin, runtime injection
* Evidence: Explicit execution model
* Winner: runtime injection
* Confidence: 0.92

### D-004 Go as implementation language

* Options: Go, Rust, Bash
* Evidence: Explicit in Idea Pack
* Winner: Go
* Confidence: 0.88

---

## 16. Assumptions

* Target platform is Unix-like systems with tmux installed
* MVP timebox is 2 to 4 weeks
* CLI interaction uses a simple TUI selector
* Hash length is short but collision-resistant
* Non-interactive shortcuts are deferred beyond v1

---

> **This PRD is complete.**
> Copy this Markdown into Word, Google Docs, Notion, or directly into a coding model.
