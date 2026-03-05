package tmux

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

// LaunchAction describes the outcome of a CreateSession call.
type LaunchAction string

const (
	LaunchActionCreate         LaunchAction = "create"
	LaunchActionAttachExisting LaunchAction = "attach-existing"
)

// AttachToSession attaches to or switches to an existing session.
// Applies p's configuration to ensure consistent styling regardless of how the session was created.
func AttachToSession(sessionName string) error {
	// Configure the session before attaching (ensures consistent styling)
	configureSession(sessionName)

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
func CreateSession(sessionName, workingDir string) (LaunchAction, error) {
	insideTmux := os.Getenv("TMUX") != ""
	if err := newDetachedSession(sessionName, workingDir); err != nil {
		var dupErr *duplicateSessionError
		if errors.As(err, &dupErr) {
			matching, matchErr := sessionMatchesDirectory(sessionName, workingDir)
			if matchErr != nil {
				return "", matchErr
			}
			if !matching {
				return "", err
			}
			if err := AttachToSession(sessionName); err != nil {
				return "", err
			}
			return LaunchActionAttachExisting, nil
		}
		return "", err
	}

	configureSession(sessionName)
	createDefaultWindows(sessionName, workingDir)
	if insideTmux {
		return LaunchActionCreate, execTmux("switch-client", "-t", sessionName)
	}
	return LaunchActionCreate, execTmux("attach-session", "-t", sessionName)
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

func newDetachedSession(sessionName, workingDir string) error {
	args := []string{"-f", "/dev/null", "new-session", "-d", "-s", sessionName, "-c", workingDir}
	cmd := exec.Command("tmux", args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		if strings.Contains(stderr.String(), "duplicate session") {
			return &duplicateSessionError{sessionName: sessionName}
		}
		msg := strings.TrimSpace(stderr.String())
		if msg == "" {
			msg = err.Error()
		}
		return fmt.Errorf("failed to create session: %s", msg)
	}
	return nil
}

type duplicateSessionError struct {
	sessionName string
}

func (e *duplicateSessionError) Error() string {
	return fmt.Sprintf("duplicate session: %s", e.sessionName)
}

func sessionMatchesDirectory(sessionName, dir string) (bool, error) {
	existing, err := GetSessionPath(sessionName)
	if err != nil {
		return false, err
	}

	left, err := fingerprintPath(existing)
	if err != nil {
		return false, err
	}
	right, err := fingerprintPath(dir)
	if err != nil {
		return false, err
	}

	if left.canonicalOK && right.canonicalOK {
		return left.canonical == right.canonical, nil
	}
	if left.canonicalErr != nil {
		fmt.Fprintf(os.Stderr, "warning: unable to fully resolve session %s path: %v\n", sessionName, left.canonicalErr)
	}
	if right.canonicalErr != nil {
		fmt.Fprintf(os.Stderr, "warning: unable to fully resolve requested path %s: %v\n", dir, right.canonicalErr)
	}
	if left.abs == right.abs {
		return true, nil
	}
	return false, nil
}

type pathFingerprint struct {
	abs          string
	canonical    string
	canonicalOK  bool
	canonicalErr error
}

func fingerprintPath(p string) (pathFingerprint, error) {
	abs, err := filepath.Abs(p)
	if err != nil {
		return pathFingerprint{}, err
	}
	fp := pathFingerprint{abs: abs}
	resolved, err := filepath.EvalSymlinks(abs)
	if err != nil {
		if os.IsNotExist(err) {
			return fp, nil
		}
		fp.canonicalErr = err
		return fp, nil
	}
	fp.canonical = resolved
	fp.canonicalOK = true
	return fp, nil
}

// createDefaultWindows creates default windows for a new session based on P_WINDOWS env var.
// Default: home,cmd
func createDefaultWindows(sessionName, workingDir string) {
	// Get window names from env or use default
	windowsEnv := os.Getenv("P_WINDOWS")
	if windowsEnv == "" {
		windowsEnv = "home,cmd"
	}

	windows := strings.Split(windowsEnv, ",")
	if len(windows) == 0 {
		return
	}

	target := "-t" + sessionName

	// Rename first window (index 0)
	runTmuxSilent("rename-window", target+":0", strings.TrimSpace(windows[0]))

	// Create additional windows
	for i := 1; i < len(windows); i++ {
		name := strings.TrimSpace(windows[i])
		runTmuxSilent("new-window", target, "-c", workingDir, "-n", name)
	}

	// Select window 0
	runTmuxSilent("select-window", target+":0")
}

// configureSession injects minimal ergonomic defaults into a tmux session.
// Per ADR-002: vi-style copy mode bindings
// Per ADR-005: color-agnostic status bar styling
func configureSession(sessionName string) {
	target := "-t" + sessionName

	// Vi-style copy mode (ADR-002)
	runTmuxSilent("set-option", target, "mode-keys", "vi")
	runTmuxSilent("bind-key", "-T", "copy-mode-vi", "v", "send-keys", "-X", "begin-selection")
	runTmuxSilent("bind-key", "-T", "copy-mode-vi", "y", "send-keys", "-X", "copy-selection-and-cancel")

	// Status bar styling (ADR-005) - Savvy AI aesthetic
	// Dark gray background (colour235 ≈ #4e4e4e - matches Savvy AI navbar)
	// Light gray foreground (colour240), sage green accent (colour108)
	// Reset status to single line with default format (repairs any corrupted status-format)
	runTmuxSilent("set-option", target, "status", "on")
	runTmuxSilent("set-option", target, "status-format[0]", "#[align=left range=left #{E:status-left-style}]#[push-default]#{T;=/#{status-left-length}:status-left}#[pop-default]#[norange default]#[list=on align=#{status-justify}]#[list=left-marker]<#[list=right-marker]>#[list=on]#{W:#[range=window|#{window_index} #{E:window-status-style}#{?#{&&:#{window_last_flag},#{!=:#{E:window-status-last-style},default}}, #{E:window-status-last-style},}#{?#{&&:#{window_bell_flag},#{!=:#{E:window-status-bell-style},default}}, #{E:window-status-bell-style},#{?#{&&:#{||:#{window_activity_flag},#{window_silence_flag}},#{!=:#{E:window-status-activity-style},default}}, #{E:window-status-activity-style},}}]#[push-default]#{T:window-status-format}#[pop-default]#[norange default]#{?window_end_flag,,#{window-status-separator}},#[range=window|#{window_index} list=focus #{?#{!=:#{E:window-status-current-style},default},#{E:window-status-current-style},#{E:window-status-style}}#{?#{&&:#{window_last_flag},#{!=:#{E:window-status-last-style},default}}, #{E:window-status-last-style},}#{?#{&&:#{window_bell_flag},#{!=:#{E:window-status-bell-style},default}}, #{E:window-status-bell-style},#{?#{&&:#{||:#{window_activity_flag},#{window_silence_flag}},#{!=:#{E:window-status-activity-style},default}}, #{E:window-status-activity-style},}}]#[push-default]#{T:window-status-current-format}#[pop-default]#[norange default]#{?window_end_flag,,#{window-status-separator}}}#[nolist align=right range=right #{E:status-right-style}]#[push-default]#{T;=/#{status-right-length}:status-right}#[pop-default]#[norange default]")
	runTmuxSilent("set-option", target, "status-style", "bg=colour235,fg=colour240")

	// Pure black pane background (main interface)
	// Use global options (-g) so all windows inherit the black background
	runTmuxSilent("set-option", "-g", "window-style", "bg=colour16")
	runTmuxSilent("set-option", "-g", "window-active-style", "bg=colour16")

	// Remove pane borders entirely
	runTmuxSilent("set-option", target, "pane-border-status", "off")
	runTmuxSilent("set-option", target, "status-left-length", "40")
	runTmuxSilent("set-option", target, "status-right-length", "40")
	runTmuxSilent("set-option", target, "status-left", "#[fg=colour108][#S] ")
	runTmuxSilent("set-option", target, "status-right", "#[fg=colour108]#(git -C #{pane_current_path} rev-parse --abbrev-ref HEAD 2>/dev/null) ")
	runTmuxSilent("set-option", target, "status-interval", "5")

	// Window status styling with spacing
	runTmuxSilent("set-option", target, "window-status-separator", "  ")
	runTmuxSilent("set-window-option", target, "window-status-format", "#[fg=colour240] #I:#W ")
	runTmuxSilent("set-window-option", target, "window-status-current-format", "#[fg=white,bold] #I:#W ")

}

// runTmuxSilent runs a tmux command silently, ignoring errors.
// Used for configuration where failure is non-fatal.
func runTmuxSilent(args ...string) {
	fullArgs := append([]string{"-f", "/dev/null"}, args...)
	cmd := exec.Command("tmux", fullArgs...)
	_ = cmd.Run()
}
