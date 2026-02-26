package tmux

import (
	"os"
	"os/exec"
	"syscall"
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
func CreateSession(sessionName, workingDir string) error {
	// Check if we're inside tmux
	if os.Getenv("TMUX") != "" {
		// Inside tmux: create detached session first, configure it, then switch
		if err := runTmux("new-session", "-d", "-s", sessionName, "-c", workingDir); err != nil {
			return err
		}
		configureSession(sessionName)
		return execTmux("switch-client", "-t", sessionName)
	}
	// Outside tmux: create detached, configure, then attach
	if err := runTmux("new-session", "-d", "-s", sessionName, "-c", workingDir); err != nil {
		return err
	}
	configureSession(sessionName)
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

// runTmux runs a tmux command and returns after completion.
func runTmux(args ...string) error {
	fullArgs := append([]string{"-f", "/dev/null"}, args...)
	cmd := exec.Command("tmux", fullArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
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
	// Near-black background (colour232 ≈ rgba(5,5,5,.9))
	// Sage green accent (colour108), muted gray (colour240), white for emphasis
	// Reset status to single line with default format (repairs any corrupted status-format)
	runTmuxSilent("set-option", target, "status", "on")
	runTmuxSilent("set-option", target, "status-format[0]", "#[align=left range=left #{E:status-left-style}]#[push-default]#{T;=/#{status-left-length}:status-left}#[pop-default]#[norange default]#[list=on align=#{status-justify}]#[list=left-marker]<#[list=right-marker]>#[list=on]#{W:#[range=window|#{window_index} #{E:window-status-style}#{?#{&&:#{window_last_flag},#{!=:#{E:window-status-last-style},default}}, #{E:window-status-last-style},}#{?#{&&:#{window_bell_flag},#{!=:#{E:window-status-bell-style},default}}, #{E:window-status-bell-style},#{?#{&&:#{||:#{window_activity_flag},#{window_silence_flag}},#{!=:#{E:window-status-activity-style},default}}, #{E:window-status-activity-style},}}]#[push-default]#{T:window-status-format}#[pop-default]#[norange default]#{?window_end_flag,,#{window-status-separator}},#[range=window|#{window_index} list=focus #{?#{!=:#{E:window-status-current-style},default},#{E:window-status-current-style},#{E:window-status-style}}#{?#{&&:#{window_last_flag},#{!=:#{E:window-status-last-style},default}}, #{E:window-status-last-style},}#{?#{&&:#{window_bell_flag},#{!=:#{E:window-status-bell-style},default}}, #{E:window-status-bell-style},#{?#{&&:#{||:#{window_activity_flag},#{window_silence_flag}},#{!=:#{E:window-status-activity-style},default}}, #{E:window-status-activity-style},}}]#[push-default]#{T:window-status-current-format}#[pop-default]#[norange default]#{?window_end_flag,,#{window-status-separator}}}#[nolist align=right range=right #{E:status-right-style}]#[push-default]#{T;=/#{status-right-length}:status-right}#[pop-default]#[norange default]")
	runTmuxSilent("set-option", target, "status-style", "bg=colour232,fg=colour240")
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
