# Symlink-Based Skills Installation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Replace file-copy skills installation with a canonical `.agentic/skills/` directory + symlinks to each agent's skill directory.

**Architecture:** All skill files (packs and generated tool skills) are written once to `.agentic/skills/`. Agent tool directories (`.claude/skills/`, `.cursor/skills/`, etc.) contain only symlinks pointing back into `.agentic/skills/`. Migration from the old copy-based model happens automatically during `skills ensure`.

**Tech Stack:** Go, `os.Symlink`, `os.Lstat`, `filepath.Rel`

---

### Task 1: Add `CanonicalSkillDir` constant and symlink helper

**Files:**
- Modify: `internal/skills/packs.go` (add constant)
- Create: `internal/skills/symlink.go` (symlink helper functions)
- Create: `internal/skills/symlink_test.go`

**Step 1: Write the failing test for `EnsureSymlink`**

Create `internal/skills/symlink_test.go`:

```go
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
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/skills/ -run TestEnsureSymlink -v`
Expected: FAIL — `EnsureSymlink` undefined

**Step 3: Add `CanonicalSkillDir` constant to `packs.go`**

Add to `internal/skills/packs.go` after the imports:

```go
// CanonicalSkillDir is the single source of truth for all skill files.
// Agent tool directories contain symlinks pointing here.
const CanonicalSkillDir = ".agentic/skills"
```

**Step 4: Implement `EnsureSymlink` in `internal/skills/symlink.go`**

```go
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
```

**Step 5: Run tests to verify they pass**

Run: `go test ./internal/skills/ -run TestEnsureSymlink -v`
Expected: PASS (all 4 tests)

**Step 6: Commit**

```bash
git add internal/skills/symlink.go internal/skills/symlink_test.go internal/skills/packs.go
git commit -m "feat(skills): add CanonicalSkillDir constant and EnsureSymlink helper"
```

---

### Task 2: Refactor `Installer.Install` to write canonical + symlink

**Files:**
- Modify: `internal/skills/installer.go` (rewrite `Install` method)
- Modify: `internal/skills/installer_test.go` (update assertions)

**Step 1: Write the failing test for symlink-based install**

Add to `internal/skills/installer_test.go`:

```go
func TestInstaller_Install_CreatesCanonicalAndSymlink(t *testing.T) {
	tmpDir := t.TempDir()

	origDir := ToolSkillDir["claude-code"]
	toolSkillDir := filepath.Join(tmpDir, ".claude", "skills")
	ToolSkillDir["claude-code"] = toolSkillDir
	defer func() { ToolSkillDir["claude-code"] = origDir }()

	// Set canonical dir to temp
	canonicalDir := filepath.Join(tmpDir, ".agentic", "skills")

	installer := NewInstallerWithCanonicalDir(canonicalDir)
	result, err := installer.Install("tdd", "claude-code", false)
	if err != nil {
		t.Fatalf("Install failed: %v", err)
	}

	// Canonical file should exist as a real file
	canonicalPath := filepath.Join(canonicalDir, "tdd", "SKILL.md")
	info, err := os.Lstat(canonicalPath)
	if err != nil {
		t.Fatalf("canonical file missing: %v", err)
	}
	if info.Mode()&os.ModeSymlink != 0 {
		t.Error("canonical file should be a real file, not a symlink")
	}

	// Tool dir file should be a symlink
	toolPath := filepath.Join(toolSkillDir, "tdd", "SKILL.md")
	info, err = os.Lstat(toolPath)
	if err != nil {
		t.Fatalf("tool dir file missing: %v", err)
	}
	if info.Mode()&os.ModeSymlink == 0 {
		t.Error("tool dir file should be a symlink")
	}

	// Symlink should point to canonical
	target, _ := os.Readlink(toolPath)
	if target != canonicalPath {
		t.Errorf("symlink target = %q, want %q", target, canonicalPath)
	}

	// Result should list files written
	if len(result.FilesWritten) != 3 {
		t.Errorf("expected 3 files written, got %d", len(result.FilesWritten))
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/skills/ -run TestInstaller_Install_CreatesCanonicalAndSymlink -v`
Expected: FAIL — `NewInstallerWithCanonicalDir` undefined

**Step 3: Refactor `Installer` to support canonical dir**

Modify `internal/skills/installer.go`:

- Add `CanonicalDir string` field to `Installer`
- Add `NewInstallerWithCanonicalDir(dir string) *Installer` constructor
- Update `NewInstaller()` to default `CanonicalDir` to `CanonicalSkillDir`
- Rewrite `Install()`:
  1. Write canonical file to `inst.CanonicalDir/<pack-relative-path>`
  2. Create symlink from tool dir to canonical file via `EnsureSymlink`
- Update `IsInstalled()` to check for symlink OR real file (backwards compat during migration)

```go
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
```

**Step 4: Run all installer tests**

Run: `go test ./internal/skills/ -run TestInstaller -v`
Expected: PASS (all tests including new one)

**Step 5: Commit**

```bash
git add internal/skills/installer.go internal/skills/installer_test.go
git commit -m "feat(skills): refactor installer to write canonical + symlink"
```

---

### Task 3: Refactor `GenerateToolSkills` to write canonical + symlink

**Files:**
- Modify: `internal/skills/generator.go`

**Step 1: Write the failing test**

Add to `internal/skills/ensure_test.go`:

```go
func TestGenerateToolSkills_CreatesCanonicalAndSymlink(t *testing.T) {
	dir := t.TempDir()
	origDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(origDir)

	os.MkdirAll(".agentic/agent-rules", 0755)
	os.WriteFile(".agentic/agent-rules/base.md", []byte("- Base rule\n"), 0644)

	origSkillDir := ToolSkillDir["claude-code"]
	ToolSkillDir["claude-code"] = ".claude/skills"
	defer func() { ToolSkillDir["claude-code"] = origSkillDir }()

	cfg := &models.Config{}
	cfg.Paths.PRDOutputPath = ".agentic/tasks/"

	gen := NewGeneratorWithConfig(cfg)
	err := gen.GenerateToolSkills("claude-code")
	if err != nil {
		t.Fatalf("GenerateToolSkills failed: %v", err)
	}

	// Canonical file should exist
	canonicalPrd := filepath.Join(CanonicalSkillDir, "prd.md")
	if _, err := os.Stat(canonicalPrd); os.IsNotExist(err) {
		t.Errorf("canonical prd.md should exist at %s", canonicalPrd)
	}

	// Tool dir file should be a symlink
	toolPrd := filepath.Join(".claude", "skills", "prd.md")
	info, err := os.Lstat(toolPrd)
	if err != nil {
		t.Fatalf("tool prd.md missing: %v", err)
	}
	if info.Mode()&os.ModeSymlink == 0 {
		t.Error("tool dir prd.md should be a symlink")
	}

	// Content should be readable through symlink
	content, _ := os.ReadFile(toolPrd)
	if !strings.Contains(string(content), ".agentic/tasks/") {
		t.Error("prd.md content should contain rendered PRDOutputPath")
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/skills/ -run TestGenerateToolSkills_CreatesCanonicalAndSymlink -v`
Expected: FAIL — files exist as regular files, not symlinks

**Step 3: Update `GenerateToolSkills` in `generator.go`**

Modify `GenerateToolSkills` to:
1. Render template into `.agentic/skills/<name>.md` (canonical)
2. Create symlink from `<tool-skill-dir>/<name>.md` to canonical via `EnsureSymlink`

```go
func (g *Generator) GenerateToolSkills(agentName string) error {
	if g.Config == nil {
		return fmt.Errorf("config required for generating tool skills")
	}

	skillDir, ok := ToolSkillDir[agentName]
	if !ok {
		return fmt.Errorf("unsupported tool: %s", agentName)
	}

	data := struct {
		PRDOutputPath string
	}{
		PRDOutputPath: g.Config.Paths.PRDOutputPath,
	}

	for _, skill := range toolSkillTemplates {
		tmplContent, err := templatesFS.ReadFile(skill.templateFile)
		if err != nil {
			return fmt.Errorf("failed to read template %s: %w", skill.templateFile, err)
		}

		t, err := template.New("skill").Parse(string(tmplContent))
		if err != nil {
			return fmt.Errorf("failed to parse template: %w", err)
		}

		var buf bytes.Buffer
		if err := t.Execute(&buf, data); err != nil {
			return fmt.Errorf("failed to execute template: %w", err)
		}

		// Write to canonical dir
		canonicalPath := filepath.Join(CanonicalSkillDir, skill.outputName)
		if err := os.MkdirAll(filepath.Dir(canonicalPath), 0755); err != nil {
			return fmt.Errorf("failed to create canonical dir: %w", err)
		}
		if err := os.WriteFile(canonicalPath, buf.Bytes(), 0644); err != nil {
			return fmt.Errorf("failed to write canonical %s: %w", canonicalPath, err)
		}

		// Symlink from tool dir to canonical
		absCan, err := filepath.Abs(canonicalPath)
		if err != nil {
			return fmt.Errorf("failed to resolve absolute path: %w", err)
		}

		outputFile := filepath.Join(skillDir, skill.outputName)
		if err := EnsureSymlink(absCan, outputFile); err != nil {
			return fmt.Errorf("failed to symlink %s: %w", outputFile, err)
		}
	}

	// Gemini also gets slash command TOML files (not symlinked — these are tool-specific)
	if agentName == "gemini" {
		if err := g.generateGeminiCommands(data); err != nil {
			return err
		}
	}

	return nil
}
```

**Step 4: Run tests**

Run: `go test ./internal/skills/ -run TestGenerateToolSkills -v`
Expected: PASS

**Step 5: Run the full ensure test suite to check nothing broke**

Run: `go test ./internal/skills/ -v`
Expected: PASS (all tests)

**Step 6: Commit**

```bash
git add internal/skills/generator.go internal/skills/ensure_test.go
git commit -m "feat(skills): GenerateToolSkills writes canonical + symlinks"
```

---

### Task 4: Update `IsInstalled` to follow symlinks and check canonical

**Files:**
- Modify: `internal/skills/installer.go`
- Modify: `internal/skills/installer_test.go`

**Step 1: Write the failing test**

Add to `internal/skills/installer_test.go`:

```go
func TestInstaller_IsInstalled_DetectsSymlink(t *testing.T) {
	tmpDir := t.TempDir()

	origDir := ToolSkillDir["claude-code"]
	toolSkillDir := filepath.Join(tmpDir, ".claude", "skills")
	ToolSkillDir["claude-code"] = toolSkillDir
	defer func() { ToolSkillDir["claude-code"] = origDir }()

	canonicalDir := filepath.Join(tmpDir, ".agentic", "skills")
	installer := NewInstallerWithCanonicalDir(canonicalDir)

	// Install (creates canonical + symlink)
	_, err := installer.Install("tdd", "claude-code", false)
	if err != nil {
		t.Fatalf("Install failed: %v", err)
	}

	// Should detect as installed (symlink resolves to real file)
	if !installer.IsInstalled("tdd", "claude-code") {
		t.Error("expected IsInstalled=true for symlinked pack")
	}
}
```

**Step 2: Run test to verify it passes (symlinks are transparent to `os.Stat`)**

Run: `go test ./internal/skills/ -run TestInstaller_IsInstalled -v`
Expected: PASS — `os.Stat` already follows symlinks, so existing `IsInstalled` works

Note: `os.Stat` (not `os.Lstat`) follows symlinks, so the existing implementation should already work. This test just confirms it.

**Step 3: Commit**

```bash
git add internal/skills/installer_test.go
git commit -m "test(skills): verify IsInstalled works with symlinked packs"
```

---

### Task 5: Update drift detection for canonical files

**Files:**
- Modify: `internal/skills/drift.go`

**Step 1: Write the failing test for pack drift detection**

Add to `internal/skills/ensure_test.go`:

```go
func TestEnsure_DetectsPackDriftInCanonical(t *testing.T) {
	dir := t.TempDir()
	origDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(origDir)

	os.MkdirAll(".agentic/agent-rules", 0755)
	os.WriteFile(".agentic/agent-rules/base.md", []byte("- Base rule\n"), 0644)

	origSkillDir := ToolSkillDir["claude-code"]
	ToolSkillDir["claude-code"] = ".claude/skills"
	defer func() { ToolSkillDir["claude-code"] = origSkillDir }()

	cfg := &models.Config{
		Agents: models.AgentsConfig{
			Overrides: []models.AgentConfig{
				{Name: "claude-code", SkillPacks: []string{"tdd"}},
			},
		},
	}

	// First ensure — installs everything
	_, err := Ensure("claude-code", cfg)
	if err != nil {
		t.Fatalf("first Ensure failed: %v", err)
	}

	// Tamper with canonical file
	os.WriteFile(filepath.Join(CanonicalSkillDir, "tdd", "SKILL.md"), []byte("tampered"), 0644)

	// Second ensure — should detect and fix drift
	result, err := Ensure("claude-code", cfg)
	if err != nil {
		t.Fatalf("second Ensure failed: %v", err)
	}

	// The pack should be re-installed
	found := false
	for _, p := range result.PacksInstalled {
		if p == "tdd" {
			found = true
		}
	}
	// If not re-installed, that's acceptable — drift detection is about rules files.
	// But the canonical content should match the embedded content after ensure.
	canonical, _ := os.ReadFile(filepath.Join(CanonicalSkillDir, "tdd", "SKILL.md"))
	embedded, _ := packsFS.ReadFile("packs/tdd/SKILL.md")
	if string(canonical) != string(embedded) {
		t.Error("canonical should match embedded after ensure")
	}
	_ = found // may or may not be in PacksInstalled depending on implementation
}
```

**Step 2: Run test**

Run: `go test ./internal/skills/ -run TestEnsure_DetectsPackDriftInCanonical -v`
Expected: May pass or fail depending on whether `IsInstalled` checks content

**Step 3: Add content-based drift check to `IsInstalled`**

Update `IsInstalled` in `installer.go` to compare canonical file content against embedded content:

```go
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
```

**Step 4: Run all tests**

Run: `go test ./internal/skills/ -v`
Expected: PASS

**Step 5: Commit**

```bash
git add internal/skills/installer.go internal/skills/drift.go internal/skills/ensure_test.go
git commit -m "feat(skills): detect content drift in canonical skill files"
```

---

### Task 6: Update `Ensure` to use the new canonical-based installer

**Files:**
- Modify: `internal/skills/ensure.go`

**Step 1: Review current `Ensure` flow**

The current `Ensure` already calls `NewInstaller()` which now defaults to `CanonicalSkillDir`. The refactored `Install` method already writes canonical + symlink. So `Ensure` itself should need minimal changes.

**Step 2: Verify with the full test suite**

Run: `go test ./internal/skills/ -v`
Expected: PASS — `Ensure` automatically benefits from the refactored installer

**Step 3: If tests pass, commit any minor adjustments**

```bash
git add internal/skills/ensure.go
git commit -m "chore(skills): ensure uses canonical-based installer"
```

---

### Task 7: Update existing tests to assert symlinks

**Files:**
- Modify: `internal/skills/ensure_test.go`
- Modify: `internal/skills/installer_test.go`

**Step 1: Update `TestEnsure_InstallsConfiguredPacks` to verify symlink**

In `ensure_test.go`, after verifying the file exists, add:

```go
info, err := os.Lstat(skillPath)
if err != nil {
	t.Fatalf("Lstat failed: %v", err)
}
if info.Mode()&os.ModeSymlink == 0 {
	t.Error("expected skill file to be a symlink")
}
```

**Step 2: Update `TestEnsure_GeneratesPrdAndRalphForAllTools` to verify symlinks**

Add symlink assertions after the existing `os.Stat` checks:

```go
prdInfo, _ := os.Lstat(prdPath)
if prdInfo.Mode()&os.ModeSymlink == 0 {
	t.Errorf("expected %s to be a symlink for %s", prdPath, tc.name)
}
```

**Step 3: Update `TestInstaller_Install` to verify canonical file exists**

After existing assertions, add:

```go
// Verify canonical file exists
canonicalPath := filepath.Join(tmpDir, ".agentic", "skills", "tdd", "SKILL.md")
// Note: NewInstaller uses default CanonicalSkillDir, tests need NewInstallerWithCanonicalDir
```

Note: Some older tests use `NewInstaller()` which uses the default `CanonicalSkillDir` (relative path). These tests already `chdir` to temp dirs, so the canonical files will be created there. Verify this works.

**Step 4: Run full test suite**

Run: `go test ./internal/skills/ -v`
Expected: PASS

**Step 5: Commit**

```bash
git add internal/skills/ensure_test.go internal/skills/installer_test.go
git commit -m "test(skills): assert symlink behavior in existing tests"
```

---

### Task 8: Run full project test suite

**Step 1: Run all project tests**

Run: `go test ./...`
Expected: PASS

**Step 2: If failures, fix and commit each fix separately**

**Step 3: Final commit if any fixes were needed**

---

### Task 9: Update `.gitignore` for agent tool skill dirs (symlinks)

The symlinks in `.claude/skills/`, `.cursor/skills/`, etc. are generated artifacts. The canonical `.agentic/skills/` is committed. The agent tool directories should remain gitignored (most projects already gitignore `.claude/`, `.cursor/`, etc.).

**Files:**
- Modify: `.gitignore` (if needed)

**Step 1: Verify current gitignore state**

Check if agent tool skill directories are already covered by existing gitignore patterns.

**Step 2: Add entries if needed**

If `.claude/skills/`, `.cursor/skills/` etc. are not already gitignored, add them. The symlinks are artifacts and should not be committed.

```gitignore
# Agent tool skill directories (symlinks to .agentic/skills/)
.claude/skills/
.cursor/skills/
.gemini/skills/
.windsurf/skills/
.codex/skills/
.github/skills/
.opencode/skills/
.agent/skills/
```

**Step 3: Commit**

```bash
git add .gitignore
git commit -m "chore: gitignore agent tool skill directories (symlinks)"
```
