package skills

import (
	"embed"
	"fmt"
)

//go:embed packs/*
var packsFS embed.FS

// ToolSkillDir maps agent tool names to their project-level skill directory.
var ToolSkillDir = map[string]string{
	"claude-code": ".claude/skills",
	"cursor":      ".cursor/skills",
	"gemini":      ".gemini/skills",
	"windsurf":    ".windsurf/skills",
	"antigravity": ".agent/skills",
	"codex":       ".codex/skills",
}

// ToolGlobalSkillDir maps agent tool names to their global/user-level skill directory.
var ToolGlobalSkillDir = map[string]string{
	"claude-code": "~/.claude/skills",
	"cursor":      "~/.cursor/skills",
	"gemini":      "~/.gemini/skills",
	"windsurf":    "~/.codeium/windsurf/skills",
	"antigravity": "~/.gemini/antigravity/skills",
	"codex":       "~/.codex/skills",
}

// SupportedTools returns the list of supported tool names.
func SupportedTools() []string {
	tools := make([]string, 0, len(ToolSkillDir))
	for t := range ToolSkillDir {
		tools = append(tools, t)
	}
	return tools
}

// SkillPackFile represents a single file within a skill pack.
type SkillPackFile struct {
	SrcPath string // path within embedded FS (e.g., "packs/tdd/SKILL.md")
	DstPath string // relative to tool skill dir (e.g., "tdd/SKILL.md")
}

// SkillPack is a named bundle of related skill files, tool-agnostic.
type SkillPack struct {
	Name        string
	Description string
	Files       []SkillPackFile
}

// PackRegistry maintains a map of available skill packs.
type PackRegistry struct {
	packs map[string]SkillPack
}

// NewPackRegistry creates a registry with built-in skill packs.
func NewPackRegistry() *PackRegistry {
	r := &PackRegistry{
		packs: make(map[string]SkillPack),
	}

	r.Register(SkillPack{
		Name:        "tdd",
		Description: "Test-Driven Development with red-green-refactor workflow",
		Files: []SkillPackFile{
			{SrcPath: "packs/tdd/SKILL.md", DstPath: "tdd/SKILL.md"},
			{SrcPath: "packs/tdd/resources/red-green-refactor.md", DstPath: "tdd/resources/red-green-refactor.md"},
			{SrcPath: "packs/tdd/resources/implementation-playbook.md", DstPath: "tdd/resources/implementation-playbook.md"},
		},
	})

	r.Register(SkillPack{
		Name:        "api-docs",
		Description: "Generate comprehensive API documentation from code",
		Files: []SkillPackFile{
			{SrcPath: "packs/api-docs/SKILL.md", DstPath: "api-docs/SKILL.md"},
		},
	})

	r.Register(SkillPack{
		Name:        "code-simplification",
		Description: "Review and refactor code for simplicity and maintainability",
		Files: []SkillPackFile{
			{SrcPath: "packs/code-simplification/SKILL.md", DstPath: "code-simplification/SKILL.md"},
		},
	})

	r.Register(SkillPack{
		Name:        "dev-plans",
		Description: "Create structured development plans with phased task breakdowns",
		Files: []SkillPackFile{
			{SrcPath: "packs/dev-plans/SKILL.md", DstPath: "dev-plans/SKILL.md"},
		},
	})

	r.Register(SkillPack{
		Name:        "diataxis",
		Description: "Write documentation using the Diataxis framework (tutorials, how-to, reference, explanation)",
		Files: []SkillPackFile{
			{SrcPath: "packs/diataxis/SKILL.md", DstPath: "diataxis/SKILL.md"},
			{SrcPath: "packs/diataxis/resources/principles.md", DstPath: "diataxis/resources/principles.md"},
			{SrcPath: "packs/diataxis/resources/reference.md", DstPath: "diataxis/resources/reference.md"},
		},
	})

	r.Register(SkillPack{
		Name:        "extract-wisdom",
		Description: "Extract insights and actionable takeaways from text sources",
		Files: []SkillPackFile{
			{SrcPath: "packs/extract-wisdom/SKILL.md", DstPath: "extract-wisdom/SKILL.md"},
		},
	})

	return r
}

// Register adds a skill pack to the registry.
func (r *PackRegistry) Register(pack SkillPack) {
	r.packs[pack.Name] = pack
}

// GetPack returns a skill pack by name.
func (r *PackRegistry) GetPack(name string) (SkillPack, error) {
	p, ok := r.packs[name]
	if !ok {
		return SkillPack{}, fmt.Errorf("unknown skill pack: %s", name)
	}
	return p, nil
}

// GetAll returns all registered skill packs.
func (r *PackRegistry) GetAll() []SkillPack {
	var list []SkillPack
	for _, p := range r.packs {
		list = append(list, p)
	}
	return list
}
