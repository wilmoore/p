package tmux

// InjectedConfig returns the minimal tmux configuration to inject at runtime.
// These are the only overrides to stock tmux behavior.
func InjectedConfig() []string {
	return []string{
		// Enable vi-style copy mode
		"set-option -g mode-keys vi",
		// Bind v to begin-selection in copy mode
		"bind-key -T copy-mode-vi v send-keys -X begin-selection",
		// Bind y to copy selection and cancel
		"bind-key -T copy-mode-vi y send-keys -X copy-selection-and-cancel",
	}
}

// InjectConfig sends the minimal config commands to the tmux server.
func InjectConfig() error {
	configs := [][]string{
		// Vi-style copy mode
		{"set-option", "-g", "mode-keys", "vi"},
		{"bind-key", "-T", "copy-mode-vi", "v", "send-keys", "-X", "begin-selection"},
		{"bind-key", "-T", "copy-mode-vi", "y", "send-keys", "-X", "copy-selection-and-cancel"},

		// Status bar styling - minimal and color-scheme agnostic
		// Use default colors to inherit terminal theme
		{"set-option", "-g", "status-style", "bg=default,fg=default"},

		// Left side: session name
		{"set-option", "-g", "status-left", " #S "},
		{"set-option", "-g", "status-left-style", "fg=colour240"},
		{"set-option", "-g", "status-left-length", "20"},

		// Right side: git branch (if in git repo)
		{"set-option", "-g", "status-right", " #(git -C #{pane_current_path} rev-parse --abbrev-ref HEAD 2>/dev/null) "},
		{"set-option", "-g", "status-right-style", "fg=colour240"},
		{"set-option", "-g", "status-right-length", "40"},

		// Window status - current window is bold
		{"set-option", "-g", "window-status-format", " #I:#W "},
		{"set-option", "-g", "window-status-current-format", " #I:#W "},
		{"set-option", "-g", "window-status-style", "fg=colour240"},
		{"set-option", "-g", "window-status-current-style", "fg=default,bold"},

		// Refresh status bar every 5 seconds (for git branch updates)
		{"set-option", "-g", "status-interval", "5"},
	}

	for _, args := range configs {
		if err := runTmux(args...); err != nil {
			// Ignore errors - these are nice-to-have ergonomics
			continue
		}
	}
	return nil
}
