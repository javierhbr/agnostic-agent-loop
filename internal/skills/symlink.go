package skills

import (
	"fmt"
	"os"
	"path/filepath"
)

// EnsureSymlink ensures dst is a symlink pointing to src.
// If dst is a regular file (legacy copy), it is removed and replaced.
// If dst is a symlink with the wrong target, it is re-created.
// Creates parent directories for dst as needed.
func EnsureSymlink(src, dst string) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return fmt.Errorf("create parent dir for %s: %w", dst, err)
	}

	info, err := os.Lstat(dst)
	if err == nil {
		// dst exists — check what it is
		if info.Mode()&os.ModeSymlink != 0 {
			// Already a symlink — check target
			target, readErr := os.Readlink(dst)
			if readErr == nil && target == src {
				return nil // Already correct
			}
		}
		// Regular file or wrong symlink — remove it
		if err := os.Remove(dst); err != nil {
			return fmt.Errorf("remove existing %s: %w", dst, err)
		}
	}

	return os.Symlink(src, dst)
}
