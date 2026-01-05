package tmux

import (
	"crypto/sha256"
	"encoding/hex"
	"path/filepath"
	"strings"
)

// sessionNameRegistry tracks used session names to detect collisions.
var sessionNameRegistry = make(map[string]string)

// GenerateSessionName creates a deterministic, collision-free session name.
// Uses the directory's base name, with a hash suffix if collision would occur.
func GenerateSessionName(dirPath string) string {
	baseName := filepath.Base(dirPath)
	// Sanitize for tmux (replace dots and colons which tmux doesn't allow)
	baseName = sanitizeSessionName(baseName)

	// Check if this base name is already used by a different path
	if existingPath, exists := sessionNameRegistry[baseName]; exists {
		if existingPath != dirPath {
			// Collision detected: append hash suffix
			return baseName + "-" + shortHash(dirPath)
		}
	}

	// Register this name
	sessionNameRegistry[baseName] = dirPath
	return baseName
}

// sanitizeSessionName removes or replaces characters that tmux doesn't allow.
func sanitizeSessionName(name string) string {
	// tmux session names cannot contain: . : (and some others)
	name = strings.ReplaceAll(name, ".", "-")
	name = strings.ReplaceAll(name, ":", "-")
	return name
}

// shortHash returns a short hash suffix for collision resolution.
func shortHash(path string) string {
	h := sha256.Sum256([]byte(path))
	return hex.EncodeToString(h[:])[:6]
}

// ResetRegistry clears the session name registry (for testing).
func ResetRegistry() {
	sessionNameRegistry = make(map[string]string)
}
