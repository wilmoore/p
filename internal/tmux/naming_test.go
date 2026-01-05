package tmux

import (
	"testing"
)

func TestGenerateSessionName(t *testing.T) {
	ResetRegistry()

	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{
			name:     "simple directory name",
			path:     "/home/user/projects/myapp",
			expected: "myapp",
		},
		{
			name:     "directory with dots",
			path:     "/home/user/projects/my.app",
			expected: "my-app",
		},
		{
			name:     "directory with colons",
			path:     "/home/user/projects/my:app",
			expected: "my-app",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ResetRegistry()
			got := GenerateSessionName(tt.path)
			if got != tt.expected {
				t.Errorf("GenerateSessionName(%q) = %q, want %q", tt.path, got, tt.expected)
			}
		})
	}
}

func TestGenerateSessionNameCollision(t *testing.T) {
	ResetRegistry()

	// First path gets the clean name
	first := GenerateSessionName("/home/user/work/api")
	if first != "api" {
		t.Errorf("First path should get clean name 'api', got %q", first)
	}

	// Second path with same base name gets hash suffix
	second := GenerateSessionName("/home/user/personal/api")
	if second == "api" {
		t.Error("Second path should have hash suffix")
	}
	if len(second) <= 4 {
		t.Errorf("Second path should be longer than 'api', got %q", second)
	}

	// Same path should get same name
	firstAgain := GenerateSessionName("/home/user/work/api")
	if firstAgain != first {
		t.Errorf("Same path should get same name, got %q want %q", firstAgain, first)
	}
}

func TestSanitizeSessionName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"simple", "simple"},
		{"with.dot", "with-dot"},
		{"with:colon", "with-colon"},
		{"with.multiple.dots", "with-multiple-dots"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := sanitizeSessionName(tt.input)
			if got != tt.expected {
				t.Errorf("sanitizeSessionName(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}
