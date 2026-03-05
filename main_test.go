package main

import "testing"

func TestParseArgs(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		expectKind  commandKind
		path        string
		sessionName string
		wantErr     bool
	}{
		{"no args", nil, commandSelector, "", "", false},
		{"version", []string{"--version"}, commandVersion, "", "", false},
		{"help", []string{"--help"}, commandHelp, "", "", false},
		{"history", []string{"--log"}, commandHistory, "", "", false},
		{"create", []string{"proj"}, commandCreate, "proj", "", false},
		{"create name", []string{"proj", "--name", "custom"}, commandCreate, "proj", "custom", false},
		{"create name equals", []string{"proj", "--name=nice"}, commandCreate, "proj", "nice", false},
		{"history extra", []string{"--log", "foo"}, 0, "", "", true},
		{"unknown flag", []string{"proj", "--foo"}, 0, "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd, err := parseArgs(tt.args)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if cmd.kind != tt.expectKind {
				t.Fatalf("kind mismatch: got %v want %v", cmd.kind, tt.expectKind)
			}
			if cmd.path != tt.path {
				t.Fatalf("path mismatch: got %q want %q", cmd.path, tt.path)
			}
			if cmd.sessionName != tt.sessionName {
				t.Fatalf("session name mismatch: got %q want %q", cmd.sessionName, tt.sessionName)
			}
		})
	}
}
