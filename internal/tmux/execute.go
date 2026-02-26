package tmux

import (
	"os"
	"os/exec"
	"syscall"
)

// AttachToSession attaches to or switches to an existing session.
func AttachToSession(sessionName string) error {
	// Check if we're inside tmux
	if os.Getenv("TMUX") != "" {
		// Inside tmux: switch client
		return execTmux("switch-client", "-t", sessionName)
	}
	// Outside tmux: attach
	return execTmux("attach-session", "-t", sessionName)
}

// CreateSession creates a new tmux session and attaches to it.
// If already inside tmux, creates the session detached and then switches to it.
func CreateSession(sessionName, workingDir string) error {
	// Check if we're inside tmux
	if os.Getenv("TMUX") != "" {
		// Inside tmux: create detached session first, then switch
		if err := runTmux("new-session", "-d", "-s", sessionName, "-c", workingDir); err != nil {
			return err
		}
		return execTmux("switch-client", "-t", sessionName)
	}
	// Outside tmux: create and attach in one step
	return execTmux("new-session", "-s", sessionName, "-c", workingDir)
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

// runTmux runs a tmux command and returns after completion.
func runTmux(args ...string) error {
	fullArgs := append([]string{"-f", "/dev/null"}, args...)
	cmd := exec.Command("tmux", fullArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
