package tmux

import (
	"os"
	"os/exec"
	"syscall"
)

// Choice represents the user's selection from the UI.
type Choice struct {
	IsSession bool
	Session   *Session
	Directory *Directory
}

// Execute performs the tmux action based on the user's choice.
func Execute(choice Choice) error {
	if choice.IsSession {
		return attachToSession(choice.Session.Name)
	}
	return createSession(choice.Directory)
}

// attachToSession attaches to or switches to an existing session.
func attachToSession(sessionName string) error {
	// Check if we're inside tmux
	if os.Getenv("TMUX") != "" {
		// Inside tmux: switch client
		return execTmux("switch-client", "-t", sessionName)
	}
	// Outside tmux: attach
	return execTmux("attach-session", "-t", sessionName)
}

// createSession creates a new tmux session from a directory.
func createSession(dir *Directory) error {
	sessionName := GenerateSessionName(dir.Path)

	// Check if we're inside tmux
	if os.Getenv("TMUX") != "" {
		// Inside tmux: create detached, then switch
		if err := runTmux("new-session", "-d", "-s", sessionName, "-c", dir.Path); err != nil {
			return err
		}
		// Inject minimal config
		_ = InjectConfig()
		return execTmux("switch-client", "-t", sessionName)
	}

	// Outside tmux: we need to create session and inject config
	// Create detached first so we can inject config
	if err := runTmux("new-session", "-d", "-s", sessionName, "-c", dir.Path); err != nil {
		return err
	}
	// Inject minimal config
	_ = InjectConfig()
	// Now attach
	return execTmux("attach-session", "-t", sessionName)
}

// runTmux runs a tmux command and returns.
func runTmux(args ...string) error {
	fullArgs := append([]string{"-f", "/dev/null"}, args...)
	cmd := exec.Command("tmux", fullArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// execTmux replaces the current process with tmux.
func execTmux(args ...string) error {
	tmuxPath, err := exec.LookPath("tmux")
	if err != nil {
		return err
	}

	fullArgs := append([]string{"tmux", "-f", "/dev/null"}, args...)
	return syscall.Exec(tmuxPath, fullArgs, os.Environ())
}
