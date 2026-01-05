package tmux

import (
	"bytes"
	"os/exec"
	"strings"
)

// Session represents a tmux session.
type Session struct {
	Name string
}

// ListSessions returns all existing tmux sessions.
// Returns empty slice if no server is running.
func ListSessions() ([]Session, error) {
	cmd := exec.Command("tmux", "-f", "/dev/null", "list-sessions", "-F", "#{session_name}")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		// Check if this is a "no server running" error
		if IsNoServerError(err) {
			return nil, err
		}
		return nil, err
	}

	output := strings.TrimSpace(stdout.String())
	if output == "" {
		return nil, nil
	}

	lines := strings.Split(output, "\n")
	sessions := make([]Session, 0, len(lines))
	for _, line := range lines {
		name := strings.TrimSpace(line)
		if name != "" {
			sessions = append(sessions, Session{Name: name})
		}
	}

	return sessions, nil
}

// IsNoServerError checks if the error indicates no tmux server is running.
func IsNoServerError(err error) bool {
	if err == nil {
		return false
	}
	// tmux returns exit code 1 when no server is running
	if exitErr, ok := err.(*exec.ExitError); ok {
		return exitErr.ExitCode() == 1
	}
	return false
}
