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
		{"set-option", "-g", "mode-keys", "vi"},
		{"bind-key", "-T", "copy-mode-vi", "v", "send-keys", "-X", "begin-selection"},
		{"bind-key", "-T", "copy-mode-vi", "y", "send-keys", "-X", "copy-selection-and-cancel"},
	}

	for _, args := range configs {
		if err := runTmux(args...); err != nil {
			// Ignore errors - these are nice-to-have ergonomics
			continue
		}
	}
	return nil
}
