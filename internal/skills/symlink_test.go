package skills

import (
	"os"
	"path/filepath"
	"testing"
)

func TestEnsureSymlink_CreatesNewSymlink(t *testing.T) {
	tmp := t.TempDir()
	src := filepath.Join(tmp, "canonical", "tdd", "SKILL.md")
	dst := filepath.Join(tmp, ".claude", "skills", "tdd", "SKILL.md")

	os.MkdirAll(filepath.Dir(src), 0755)
	os.WriteFile(src, []byte("# TDD Skill"), 0644)

	err := EnsureSymlink(src, dst)
	if err != nil {
		t.Fatalf("EnsureSymlink failed: %v", err)
	}

	// dst should be a symlink
	info, err := os.Lstat(dst)
	if err != nil {
		t.Fatalf("dst does not exist: %v", err)
	}
	if info.Mode()&os.ModeSymlink == 0 {
		t.Error("expected dst to be a symlink")
	}

	// Should point to src
	target, _ := os.Readlink(dst)
	if target != src {
		t.Errorf("symlink target = %q, want %q", target, src)
	}
}

func TestEnsureSymlink_ReplacesExistingFile(t *testing.T) {
	tmp := t.TempDir()
	src := filepath.Join(tmp, "canonical", "tdd", "SKILL.md")
	dst := filepath.Join(tmp, ".claude", "skills", "tdd", "SKILL.md")

	os.MkdirAll(filepath.Dir(src), 0755)
	os.WriteFile(src, []byte("# TDD Skill"), 0644)

	// Pre-existing regular file (old copy-based install)
	os.MkdirAll(filepath.Dir(dst), 0755)
	os.WriteFile(dst, []byte("old copy"), 0644)

	err := EnsureSymlink(src, dst)
	if err != nil {
		t.Fatalf("EnsureSymlink failed: %v", err)
	}

	info, _ := os.Lstat(dst)
	if info.Mode()&os.ModeSymlink == 0 {
		t.Error("expected dst to be a symlink after migration")
	}
}

func TestEnsureSymlink_IdempotentWhenCorrect(t *testing.T) {
	tmp := t.TempDir()
	src := filepath.Join(tmp, "canonical", "tdd", "SKILL.md")
	dst := filepath.Join(tmp, ".claude", "skills", "tdd", "SKILL.md")

	os.MkdirAll(filepath.Dir(src), 0755)
	os.WriteFile(src, []byte("# TDD Skill"), 0644)

	// First call
	EnsureSymlink(src, dst)
	// Second call — should be no-op
	err := EnsureSymlink(src, dst)
	if err != nil {
		t.Fatalf("idempotent EnsureSymlink failed: %v", err)
	}

	info, _ := os.Lstat(dst)
	if info.Mode()&os.ModeSymlink == 0 {
		t.Error("expected dst to still be a symlink")
	}
}

func TestEnsureSymlink_FixesWrongTarget(t *testing.T) {
	tmp := t.TempDir()
	src := filepath.Join(tmp, "canonical", "tdd", "SKILL.md")
	wrongSrc := filepath.Join(tmp, "wrong", "SKILL.md")
	dst := filepath.Join(tmp, ".claude", "skills", "tdd", "SKILL.md")

	os.MkdirAll(filepath.Dir(src), 0755)
	os.WriteFile(src, []byte("# TDD Skill"), 0644)
	os.MkdirAll(filepath.Dir(wrongSrc), 0755)
	os.WriteFile(wrongSrc, []byte("wrong"), 0644)

	// Create symlink pointing to wrong target
	os.MkdirAll(filepath.Dir(dst), 0755)
	os.Symlink(wrongSrc, dst)

	err := EnsureSymlink(src, dst)
	if err != nil {
		t.Fatalf("EnsureSymlink failed: %v", err)
	}

	target, _ := os.Readlink(dst)
	if target != src {
		t.Errorf("symlink target = %q, want %q", target, src)
	}
}
