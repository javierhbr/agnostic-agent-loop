package encoding

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/javierbenavides/agentic-agent/internal/context"
	"github.com/javierbenavides/agentic-agent/internal/skills"
	"github.com/javierbenavides/agentic-agent/internal/specs"
	"github.com/javierbenavides/agentic-agent/internal/tasks"
	"github.com/javierbenavides/agentic-agent/pkg/models"
)

type ContextBundle struct {
	Task        *models.Task              `yaml:"task" json:"task"`
	Global      *models.GlobalContext     `yaml:"global" json:"global"`
	Rolling     string                    `yaml:"rolling" json:"rolling"`
	TechStack   string                    `yaml:"tech_stack,omitempty" json:"tech_stack,omitempty"`
	Workflow    string                    `yaml:"workflow,omitempty" json:"workflow,omitempty"`
	Directories []*models.DirectoryContext `yaml:"directories" json:"directories"`
	Specs             []*specs.ResolvedSpec     `yaml:"specs,omitempty" json:"specs,omitempty"`
	SkillInstructions string                    `yaml:"skill_instructions,omitempty" json:"skill_instructions,omitempty"`
	BuiltAt           time.Time                 `yaml:"built_at" json:"built_at"`
}

func CreateContextBundle(taskID string, format string, cfg *models.Config) ([]byte, error) {
	// 1. Load Task
	tm := tasks.NewTaskManager(".agentic/tasks")
	// Search in all lists
	// Simplified: just check in-progress for building context usually
	list, err := tm.LoadTasks("in-progress")
	if err != nil {
		return nil, err
	}

	var task *models.Task
	for _, t := range list.Tasks {
		if t.ID == taskID {
			task = &t
			break
		}
	}

	if task == nil {
		// Try backlog
		list, _ = tm.LoadTasks("backlog")
		for _, t := range list.Tasks {
			if t.ID == taskID {
				task = &t
				break
			}
		}
	}

	if task == nil {
		return nil, fmt.Errorf("task %s not found", taskID)
	}

	// 2. Load Global Context
	gcm := context.NewGlobalContextManager(".agentic/context")
	global, err := gcm.LoadGlobal()
	if err != nil {
		return nil, fmt.Errorf("failed to load global context: %w", err)
	}

	// 3. Load Rolling Summary
	rcm := context.NewRollingContextManager(".agentic/context")
	rolling, err := rcm.LoadRolling()
	if err != nil {
		// rolling might be optional or empty
		rolling = ""
	}

	// 4. Load Directory Contexts (simplified: load all tracked dirs?)
	// Or just from task scope?
	// For MVP, lets looking at task scope if we had it, or just root.
	// We'll scan root for now.
	dcm := context.NewDirectoryContextManager(".")
	dirs, _ := dcm.FindContextDirs(".")
	var dirContexts []*models.DirectoryContext
	for _, d := range dirs {
		ctx, err := dcm.LoadContext(d)
		if err == nil {
			dirContexts = append(dirContexts, ctx)
		}
	}

	// Load supplementary context files (optional, non-blocking)
	techStack, _ := os.ReadFile(".agentic/context/tech-stack.md")
	workflow, _ := os.ReadFile(".agentic/context/workflow-preferences.md")

	bundle := &ContextBundle{
		Task:        task,
		Global:      global,
		Rolling:     rolling,
		TechStack:   string(techStack),
		Workflow:    string(workflow),
		Directories: dirContexts,
		BuiltAt:     time.Now(),
	}

	// 5. Resolve specs if task has SpecRefs
	if len(task.SpecRefs) > 0 {
		resolver := specs.NewResolver(cfg)
		bundle.Specs = resolver.ResolveAll(task.SpecRefs)

		// Warn on stderr for unresolved specs (non-blocking)
		for _, s := range bundle.Specs {
			if !s.Found {
				fmt.Fprintf(os.Stderr, "Warning: spec %q could not be resolved: %s\n", s.Ref, s.Error)
			}
		}
	}

	// 5.5 Load skill instructions for the active agent
	if cfg.ActiveAgent != "" {
		bundle.SkillInstructions = loadSkillInstructions(cfg.ActiveAgent)
	}

	// 6. Encode
	encoder := NewToonEncoder()
	return encoder.Encode(bundle)
}

// loadSkillInstructions gathers agent rules and installed skill pack content.
func loadSkillInstructions(agent string) string {
	var parts []string

	// 1. Load the agent's generated rules file
	registry := skills.NewSkillRegistry()
	if skill, err := registry.GetSkill(agent); err == nil {
		if content, err := os.ReadFile(skill.OutputFile); err == nil {
			parts = append(parts, string(content))
		}
	}

	// 2. Walk the agent's skill directory, concatenate SKILL.md files
	if skillDir, ok := skills.ToolSkillDir[agent]; ok {
		filepath.WalkDir(skillDir, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return nil
			}
			if d.Name() == "SKILL.md" {
				if content, readErr := os.ReadFile(path); readErr == nil {
					parts = append(parts, string(content))
				}
			}
			return nil
		})
	}

	return strings.Join(parts, "\n\n---\n\n")
}
