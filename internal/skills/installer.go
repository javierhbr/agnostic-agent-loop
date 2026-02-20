package skills

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// InstallResult holds the outcome of a pack installation.
type InstallResult struct {
	PackName     string
	Tool         string
	OutputDir    string
	FilesWritten []string
}

// Installer handles installing skill packs to tool-specific directories.
type Installer struct {
	Registry     *PackRegistry
	CanonicalDir string
}

// NewInstaller creates an installer with default pack registry and canonical dir.
func NewInstaller() *Installer {
	return &Installer{
		Registry:     NewPackRegistry(),
		CanonicalDir: CanonicalSkillDir,
	}
}

// NewInstallerWithCanonicalDir creates an installer with a custom canonical directory.
func NewInstallerWithCanonicalDir(dir string) *Installer {
	return &Installer{
		Registry:     NewPackRegistry(),
		CanonicalDir: dir,
	}
}

// Install writes a skill pack's files to the canonical directory and creates
// symlinks in the tool's skill directory.
func (inst *Installer) Install(packName, tool string, global bool) (*InstallResult, error) {
	pack, err := inst.Registry.GetPack(packName)
	if err != nil {
		return nil, err
	}

	outputDir, err := resolveOutputDir(tool, global)
	if err != nil {
		return nil, err
	}

	result := &InstallResult{
		PackName:  packName,
		Tool:      tool,
		OutputDir: outputDir,
	}

	for _, f := range pack.Files {
		content, err := packsFS.ReadFile(f.SrcPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read embedded file %s: %w", f.SrcPath, err)
		}

		// Write canonical file
		canonicalPath := filepath.Join(inst.CanonicalDir, f.DstPath)
		if err := os.MkdirAll(filepath.Dir(canonicalPath), 0755); err != nil {
			return nil, fmt.Errorf("failed to create canonical dir for %s: %w", canonicalPath, err)
		}
		if err := os.WriteFile(canonicalPath, content, 0644); err != nil {
			return nil, fmt.Errorf("failed to write canonical %s: %w", canonicalPath, err)
		}

		// Convert canonical path to absolute for symlink
		absCan, err := filepath.Abs(canonicalPath)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve absolute path for %s: %w", canonicalPath, err)
		}

		// Create symlink in tool dir
		destPath := filepath.Join(outputDir, f.DstPath)
		if err := EnsureSymlink(absCan, destPath); err != nil {
			return nil, fmt.Errorf("failed to symlink %s -> %s: %w", destPath, absCan, err)
		}

		result.FilesWritten = append(result.FilesWritten, destPath)
	}

	return result, nil
}

// InstallMulti installs a skill pack to multiple tools in one call.
// Returns all results collected so far plus the first error encountered (if any).
func (inst *Installer) InstallMulti(packName string, tools []string, global bool) ([]*InstallResult, error) {
	var results []*InstallResult
	for _, tool := range tools {
		result, err := inst.Install(packName, tool, global)
		if err != nil {
			return results, fmt.Errorf("failed for tool %s: %w", tool, err)
		}
		results = append(results, result)
	}
	return results, nil
}

// IsInstalled checks whether a pack's files exist and match the embedded content.
func (inst *Installer) IsInstalled(packName, tool string) bool {
	dir, ok := ToolSkillDir[tool]
	if !ok {
		return false
	}

	pack, err := inst.Registry.GetPack(packName)
	if err != nil {
		return false
	}

	for _, f := range pack.Files {
		// Check tool dir path exists (could be symlink or real file)
		path := filepath.Join(dir, f.DstPath)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return false
		}

		// Check canonical file content matches embedded
		canonicalPath := filepath.Join(inst.CanonicalDir, f.DstPath)
		actual, err := os.ReadFile(canonicalPath)
		if err != nil {
			return false
		}
		expected, err := packsFS.ReadFile(f.SrcPath)
		if err != nil {
			return false
		}
		if string(actual) != string(expected) {
			return false // Content drift
		}
	}
	return true
}

// IsInstalledAnywhere checks whether a pack is installed for any known tool.
// Returns the first tool name where it's found, or empty string.
func (inst *Installer) IsInstalledAnywhere(packName string) string {
	for tool := range ToolSkillDir {
		if inst.IsInstalled(packName, tool) {
			return tool
		}
	}
	return ""
}

// ListPacks returns all available skill packs.
func (inst *Installer) ListPacks() []SkillPack {
	return inst.Registry.GetAll()
}

// resolveOutputDir returns the correct output directory for a tool.
func resolveOutputDir(tool string, global bool) (string, error) {
	var dirMap map[string]string
	if global {
		dirMap = ToolGlobalSkillDir
	} else {
		dirMap = ToolSkillDir
	}

	dir, ok := dirMap[tool]
	if !ok {
		return "", fmt.Errorf("unsupported tool: %s (supported: %s)", tool, strings.Join(SupportedTools(), ", "))
	}

	// Expand ~ for global paths
	if strings.HasPrefix(dir, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to resolve home directory: %w", err)
		}
		dir = filepath.Join(home, dir[2:])
	}

	return dir, nil
}
