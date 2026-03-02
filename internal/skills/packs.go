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

// ToolAgentDir maps agent tool names to their project-level agent directory.
// Only Claude Code and OpenCode have native agent support (subagents / AGENTS.md).
// Other tools will use the skill directory as fallback.
var ToolAgentDir = map[string]string{
	"claude-code": ".claude/agents",
	"opencode":    ".agents",
}

// ToolGlobalAgentDir maps agent tool names to their global/user-level agent directory.
var ToolGlobalAgentDir = map[string]string{
	"claude-code": "~/.claude/agents",
	"opencode":    "~/.agents",
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
	IsAgent bool   // if true, install to ToolAgentDir instead of ToolSkillDir (Claude Code agents only)
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
	"agentic-helper",
	"atdd",
	"code-simplification",
	"context-manager",
	"dev-plans",
	"openspec",
	"product-wizard",
	"run-with-ralph",
	"tier-enforcer",
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
		Name:        "agentic-helper",
		Description: "Agentic Agent CLI guide and executor. Teaches workflows, runs commands, automates context generation",
		Files: []SkillPackFile{
			{SrcPath: "packs/agentic-helper/AGENT.md", DstPath: "agentic-helper.md", IsAgent: true},
			{SrcPath: "packs/agentic-helper/SKILL.md", DstPath: "agentic-helper/SKILL.md"},
		},
	})

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
		Name:        "context-manager",
		Description: "Enforce reading context.md before edits and generating context.md for new directories",
		Files: []SkillPackFile{
			{SrcPath: "packs/context-manager/SKILL.md", DstPath: "context-manager/SKILL.md"},
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

	r.Register(SkillPack{
		Name:        "superpowers-bridge",
		Description: "Integration guide for using Superpowers plugin alongside CLI and SDD skills",
		Files: []SkillPackFile{
			{SrcPath: "packs/superpowers-bridge/SKILL.md", DstPath: "superpowers-bridge/SKILL.md"},
		},
	})

	r.Register(SkillPack{
		Name:        "tier-enforcer",
		Description: "Audit, create, or fix skill files for 3-tier layered context model compliance",
		Files: []SkillPackFile{
			{SrcPath: "packs/tier-enforcer/SKILL.md", DstPath: "tier-enforcer/SKILL.md"},
			{SrcPath: "packs/tier-enforcer/resources/tier-rules.md", DstPath: "tier-enforcer/resources/tier-rules.md"},
		},
	})

	r.Register(SkillPack{
		Name:        "sdd",
		Description: "Spec-Driven Development v3.0 with four roles, five gates, and complete workflow automation",
		Files: []SkillPackFile{
			// Original SDD skills
			{SrcPath: "packs/sdd/platform-spec/SKILL.md", DstPath: "sdd/platform-spec/SKILL.md"},
			{SrcPath: "packs/sdd/component-spec/SKILL.md", DstPath: "sdd/component-spec/SKILL.md"},
			{SrcPath: "packs/sdd/gate-check/SKILL.md", DstPath: "sdd/gate-check/SKILL.md"},
			{SrcPath: "packs/sdd/adr/SKILL.md", DstPath: "sdd/adr/SKILL.md"},
			{SrcPath: "packs/sdd/hotfix/SKILL.md", DstPath: "sdd/hotfix/SKILL.md"},
			// Four role-based skills
			{SrcPath: "packs/sdd/analyst/SKILL.md", DstPath: "sdd/analyst/SKILL.md"},
			{SrcPath: "packs/sdd/architect/SKILL.md", DstPath: "sdd/architect/SKILL.md"},
			{SrcPath: "packs/sdd/developer/SKILL.md", DstPath: "sdd/developer/SKILL.md"},
			{SrcPath: "packs/sdd/verifier/SKILL.md", DstPath: "sdd/verifier/SKILL.md"},
			// Product/PM and platform skills
			{SrcPath: "packs/sdd/workflow-router/SKILL.md", DstPath: "sdd/workflow-router/SKILL.md"},
			{SrcPath: "packs/sdd/initiative-definition/SKILL.md", DstPath: "sdd/initiative-definition/SKILL.md"},
			{SrcPath: "packs/sdd/risk-assessment/SKILL.md", DstPath: "sdd/risk-assessment/SKILL.md"},
			{SrcPath: "packs/sdd/stakeholder-communication/SKILL.md", DstPath: "sdd/stakeholder-communication/SKILL.md"},
			{SrcPath: "packs/sdd/platform-constitution/SKILL.md", DstPath: "sdd/platform-constitution/SKILL.md"},
			// Complete workflow guide (new)
			{SrcPath: "packs/sdd/process-guide/SKILL.md", DstPath: "sdd/process-guide/SKILL.md"},
		},
	})

	r.Register(SkillPack{
		Name:        "openclaw",
		Description: "OpenClaw autonomous agent factory pattern — orchestrator, worker, researcher, reviewer roles with agnostic-agent CLI",
		Files: []SkillPackFile{
			{SrcPath: "packs/openclaw/SKILL.md", DstPath: "openclaw/SKILL.md"},
			{SrcPath: "packs/openclaw/resources/orchestrator.md", DstPath: "openclaw/resources/orchestrator.md"},
			{SrcPath: "packs/openclaw/resources/worker.md", DstPath: "openclaw/resources/worker.md"},
			{SrcPath: "packs/openclaw/resources/coordination.md", DstPath: "openclaw/resources/coordination.md"},
			{SrcPath: "packs/openclaw/resources/researcher.md", DstPath: "openclaw/resources/researcher.md"},
			{SrcPath: "packs/openclaw/resources/reviewer.md", DstPath: "openclaw/resources/reviewer.md"},
			{SrcPath: "packs/openclaw/AGENT.md", DstPath: "openclaw-orchestrator.md", IsAgent: true},
			{SrcPath: "packs/openclaw/agents/orchestrator.md", DstPath: "openclaw-orchestrator.md", IsAgent: true},
			{SrcPath: "packs/openclaw/agents/worker.md", DstPath: "openclaw-worker.md", IsAgent: true},
			{SrcPath: "packs/openclaw/agents/researcher.md", DstPath: "openclaw-researcher.md", IsAgent: true},
			{SrcPath: "packs/openclaw/agents/reviewer.md", DstPath: "openclaw-reviewer.md", IsAgent: true},
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
