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

// execTmux replaces the current process with tmux.
func execTmux(args ...string) error {
	tmuxPath, err := exec.LookPath("tmux")
	if err != nil {
		return err
	}

	fullArgs := append([]string{"tmux", "-f", "/dev/null"}, args...)
	return syscall.Exec(tmuxPath, fullArgs, os.Environ())
}
