package tmux

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNormalizePathResolvesSymlinks(t *testing.T) {
	root := t.TempDir()
	target := filepath.Join(root, "target")
	if err := os.Mkdir(target, 0o755); err != nil {
		t.Fatalf("mkdir target: %v", err)
	}
	link := filepath.Join(root, "link")
	if err := os.Symlink(target, link); err != nil {
		t.Fatalf("symlink: %v", err)
	}
	targetFP, err := fingerprintPath(target)
	if err != nil {
		t.Fatalf("fingerprint target: %v", err)
	}
	linkFP, err := fingerprintPath(link)
	if err != nil {
		t.Fatalf("fingerprint link: %v", err)
	}
	if !targetFP.canonicalOK || !linkFP.canonicalOK {
		t.Fatalf("expected canonical paths to resolve")
	}
	if targetFP.canonical != linkFP.canonical {
		t.Fatalf("paths differ: %q vs %q", targetFP.canonical, linkFP.canonical)
	}
}
