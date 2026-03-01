# Examples

Practical workflows for everyday tmux session management with p.

---

## Create a New Project

Start a new project and enter the session in one line:

```bash
mkdir my-new-project && p !$
```

> `!$` is bash/zsh history expansion. In fish or POSIX sh, type the path explicitly: `mkdir my-project && p my-project`

### Why this works

`!$` is bash history expansion—it's the last argument from the previous command:

1. `mkdir my-new-project` creates the directory
2. `!$` expands to `my-new-project` (the last argument)
3. `p my-new-project` creates a tmux session named after the directory

This is a common bash/zsh pattern that saves typing. The workflow becomes muscle memory:

```bash
mkdir <tab> && p !$
```

---

## Work in the Current Directory

Create a session in your current working directory:

```bash
p .
```

The session is named after the directory basename. If you're in `/Users/you/projects/api`, the session is named `api`.

### Why this works

When you pass `.` as the argument, p resolves it to the current directory path and derives the session name from the basename.

---

## Resume Work Fast

The most common workflow—jump back into an existing session:

```bash
p
# type a few characters
# Enter
```

Type to fuzzy filter, navigate if needed, press Enter. No flags, no options.

### Why this works

The selector shows all sessions and filters as you type. The more you type, the narrower the match. For frequent projects, 2-3 characters is usually enough.

---

## From Anywhere

Create a session from any path without changing directories first:

```bash
p ~/projects/backend
p /Users/you/work/api
```

### Why this works

p resolves the path internally—you don't need to `cd` first. The session name comes from the directory basename, so `p ~/projects/backend` creates a session named `backend`.

---

## Inside tmux

Run p from inside an existing tmux session:

```bash
# You're already in a session called "main"
p .
# Creates new session, switches to it seamlessly
# You're now in the new session, "main" is still there
```

### Why this works

p detects when it's running inside tmux and uses `switch-client` instead of `attach-session`. Your previous session stays intact—no need to manually detach.

---

## Rapid Context Switching

Quickly cycle between projects:

```bash
p
# type "fron" → matches "frontend"
# Enter
```

Then repeat for another project. The selector remembers your position but filters to what you type.

### Why this works

Fuzzy matching means partial matches work. "fron" matches "frontend", "api" matches "api-server". It's fast because you're just typing 2-4 characters and hitting Enter.
