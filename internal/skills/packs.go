package skills

import (
	"embed"
	"fmt"
	"sort"
)

//go:embed packs/*
var packsFS embed.FS

// CanonicalSkillDir is the single source of truth for all skill files.
// Agent tool directories contain symlinks pointing here.
const CanonicalSkillDir = ".agentic/skills"

// ToolSkillDir maps agent tool names to their project-level skill directory.
var ToolSkillDir = map[string]string{
	"claude-code": ".claude/skills",
	"cursor":      ".cursor/skills",
	"gemini":      ".gemini/skills",
	"windsurf":    ".windsurf/skills",
	"antigravity": ".agent/skills",
	"codex":       ".codex/skills",
	"copilot":     ".github/skills",
	"opencode":    ".opencode/skills",
}

// ToolGlobalSkillDir maps agent tool names to their global/user-level skill directory.
var ToolGlobalSkillDir = map[string]string{
	"claude-code": "~/.claude/skills",
	"cursor":      "~/.cursor/skills",
	"gemini":      "~/.gemini/skills",
	"windsurf":    "~/.codeium/windsurf/skills",
	"antigravity": "~/.gemini/antigravity/skills",
	"codex":       "~/.codex/skills",
	"copilot":     "~/.config/github-copilot/skills",
	"opencode":    "~/.config/opencode/skills",
}

// SupportedTools returns the list of supported tool names in sorted order.
func SupportedTools() []string {
	tools := make([]string, 0, len(ToolSkillDir))
	for t := range ToolSkillDir {
		tools = append(tools, t)
	}
	sort.Strings(tools)
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

// MandatoryPacks lists skill packs that must be installed for every agent.
// These are auto-installed during `skills ensure` and validated on startup.
var MandatoryPacks = []string{
	"atdd",
	"code-simplification",
	"dev-plans",
	"openspec",
	"product-wizard",
	"run-with-ralph",
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
		Name:        "atdd",
		Description: "Acceptance Test-Driven Development from openspec task criteria",
		Files: []SkillPackFile{
			{SrcPath: "packs/atdd/SKILL.md", DstPath: "atdd/SKILL.md"},
		},
	})

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

	r.Register(SkillPack{
		Name:        "openspec",
		Description: "Spec-driven development from requirements files using the openspec change lifecycle",
		Files: []SkillPackFile{
			{SrcPath: "packs/openspec/SKILL.md", DstPath: "openspec/SKILL.md"},
		},
	})

	r.Register(SkillPack{
		Name:        "product-wizard",
		Description: "Generate robust, production-grade Product Requirements Documents (PRDs)",
		Files: []SkillPackFile{
			{SrcPath: "packs/product-wizard/SKILL.md", DstPath: "product-wizard/SKILL.md"},
			{SrcPath: "packs/product-wizard/references/prd_template.md", DstPath: "product-wizard/references/prd_template.md"},
			{SrcPath: "packs/product-wizard/references/user_story_examples.md", DstPath: "product-wizard/references/user_story_examples.md"},
			{SrcPath: "packs/product-wizard/references/metrics_frameworks.md", DstPath: "product-wizard/references/metrics_frameworks.md"},
			{SrcPath: "packs/product-wizard/scripts/validate_prd.sh", DstPath: "product-wizard/scripts/validate_prd.sh"},
		},
	})

	r.Register(SkillPack{
		Name:        "run-with-ralph",
		Description: "Execute openspec tasks using Ralph Wiggum iterative loops",
		Files: []SkillPackFile{
			{SrcPath: "packs/run-with-ralph/SKILL.md", DstPath: "run-with-ralph/SKILL.md"},
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

// GetAll returns all registered skill packs in sorted order.
func (r *PackRegistry) GetAll() []SkillPack {
	var list []SkillPack
	for _, p := range r.packs {
		list = append(list, p)
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].Name < list[j].Name
	})
	return list
}
