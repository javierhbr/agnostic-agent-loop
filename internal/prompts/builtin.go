package prompts

var builtinPrompts = []Prompt{
	// ── Agent Prompts ───────────────────────────────────────────────────

	{
		Slug:        "claim-and-implement",
		Title:       "Claim and Implement a Task",
		Category:    CategoryAgentPrompt,
		Description: "Claim a task and implement it with TDD",
		Tags:        []string{"task", "tdd"},
		Content: `Claim task TASK-001 and implement it following TDD.

Steps:
1. Run: agentic-agent task claim TASK-001
2. Read the task details and acceptance criteria
3. Write failing tests first (RED phase)
4. Implement the minimum code to pass tests (GREEN phase)
5. Refactor for clarity (REFACTOR phase)
6. Run: agentic-agent task complete TASK-001`,
	},

	{
		Slug:        "continue-task",
		Title:       "Continue an In-Progress Task",
		Category:    CategoryAgentPrompt,
		Description: "Resume work on a previously claimed task",
		Tags:        []string{"task"},
		Content: `Continue working on my in-progress task.

Run: agentic-agent task continue
Read the task context and pick up where you left off.
Check the acceptance criteria and complete any remaining items.
When all criteria are met, run: agentic-agent task complete <TASK_ID>`,
	},

	{
		Slug:        "implement-from-spec",
		Title:       "Implement from Specification",
		Category:    CategoryAgentPrompt,
		Description: "Build a feature by reading its spec file first",
		Tags:        []string{"spec", "task"},
		Content: `Read the spec at .agentic/spec/<spec-name>.md and implement the feature.

Steps:
1. Run: agentic-agent spec resolve <spec-name>.md
2. Read and understand the specification
3. Run: agentic-agent context build --task <TASK_ID>
4. Implement according to the spec's requirements
5. Verify all acceptance criteria from the spec are met`,
	},

	{
		Slug:        "decompose-task",
		Title:       "Decompose a Large Task",
		Category:    CategoryAgentPrompt,
		Description: "Break a task into smaller subtasks",
		Tags:        []string{"task", "planning"},
		Content: `Task TASK-001 is too large. Decompose it into smaller subtasks.

Read the task with: agentic-agent task show TASK-001
Break it into 3-5 subtasks, each small enough to complete in one session.
Run: agentic-agent task decompose TASK-001 "subtask 1" "subtask 2" "subtask 3"`,
	},

	{
		Slug:        "review-and-simplify",
		Title:       "Review and Simplify Code",
		Category:    CategoryAgentPrompt,
		Description: "Run code simplification review on recent changes",
		Tags:        []string{"review", "simplify"},
		Content: `Review the code I just wrote for simplification opportunities.

Run: agentic-agent simplify .
Read the review output and apply any recommendations that improve clarity
without changing behavior. Focus on reducing complexity and improving readability.`,
	},

	{
		Slug:        "brainstorm-feature",
		Title:       "Brainstorm a New Feature",
		Category:    CategoryAgentPrompt,
		Description: "Explore and refine an idea before building it",
		Tags:        []string{"brainstorming", "planning"},
		Content: `I have an idea for a new feature: <describe your idea>.

Brainstorm this with me. Ask me questions one at a time to understand:
- What problem it solves
- Who the users are
- What the constraints are

Then propose 2-3 approaches with trade-offs and write a design doc.`,
	},

	{
		Slug:        "write-prd",
		Title:       "Write a Product Requirements Document",
		Category:    CategoryAgentPrompt,
		Description: "Generate a PRD using the product-wizard skill",
		Tags:        []string{"prd", "product", "planning"},
		Content: `Write a PRD for: <describe your feature or product>.

Use the product-wizard skill. Ask discovery questions first,
then produce a complete PRD at .agentic/spec/prd-<feature-name>.md
including: user stories, acceptance criteria, success metrics, and risks.`,
	},

	{
		Slug:        "openspec-plan",
		Title:       "Plan from Requirements with OpenSpec",
		Category:    CategoryAgentPrompt,
		Description: "Create a proposal, dev plan, and tasks from a spec file",
		Tags:        []string{"openspec", "planning", "spec"},
		Content: `Plan the implementation for .agentic/spec/<spec-name>.md using openspec.

Steps:
1. Read and understand the requirements
2. Run: agentic-agent openspec init "<change-name>" --from .agentic/spec/<spec-name>.md
3. Define the tech stack and create a development plan
4. Generate detailed task breakdown with acceptance criteria
5. Present the plan for my approval before execution`,
	},

	{
		Slug:        "openspec-execute",
		Title:       "Execute OpenSpec Tasks",
		Category:    CategoryAgentPrompt,
		Description: "Implement tasks from an approved openspec change",
		Tags:        []string{"openspec", "execute"},
		Content: `Execute the approved openspec change: <change-name>.

Work through the tasks sequentially:
1. Run: agentic-agent task list
2. Claim the next pending task
3. Implement following the dev plan
4. Run tests and verify acceptance criteria
5. Complete the task and move to the next one`,
	},

	{
		Slug:        "write-acceptance-tests",
		Title:       "Write Acceptance Tests from Criteria",
		Category:    CategoryAgentPrompt,
		Description: "Use ATDD to write executable tests before implementation",
		Tags:        []string{"atdd", "testing"},
		Content: `For task TASK-001, follow ATDD (Acceptance Test-Driven Development):

1. Read the acceptance criteria from the task
2. Write one executable acceptance test per criterion — all should FAIL
3. Implement the minimum code to make each test pass, one at a time
4. Refactor after all tests pass
5. Run the full test suite to confirm nothing broke`,
	},

	{
		Slug:        "ralph-iterate",
		Title:       "Implement Tasks with Ralph Loops",
		Category:    CategoryAgentPrompt,
		Description: "Use Ralph Wiggum iterative loops for convergent implementation",
		Tags:        []string{"ralph", "iterate"},
		Content: `Ralph my tasks. Implement the openspec tasks using iterative Ralph Wiggum loops.

For each task:
1. Claim the task
2. Build a focused prompt from the task context
3. Run a Ralph loop: implement → test → fix → converge
4. Complete the task when all tests pass
5. Move to the next task`,
	},

	{
		Slug:        "convert-prd-to-tasks",
		Title:       "Convert PRD to Task YAML",
		Category:    CategoryAgentPrompt,
		Description: "Turn a PRD into backlog tasks in YAML format",
		Tags:        []string{"prd", "tasks"},
		Content: `Convert the PRD at .agentic/spec/prd-<feature-name>.md into tasks.

Parse each user story and requirement into individual tasks with:
- Clear title and description
- Acceptance criteria from the PRD
- Dependency ordering
- One-iteration story size (break large stories down)

Write them to .agentic/tasks/backlog.yaml`,
	},

	{
		Slug:        "create-dev-plan",
		Title:       "Create a Development Plan",
		Category:    CategoryAgentPrompt,
		Description: "Generate a phased dev plan with tasks and QA checklists",
		Tags:        []string{"planning", "dev-plan"},
		Content: `Create a development plan for: <describe the feature or project>.

Include:
- Phased task breakdown with dependencies
- Requirements for each phase
- QA checklist per phase
- Risk assessment

Write the plan to DEVELOPMENT_PLAN.md`,
	},

	// ── CLI Examples ────────────────────────────────────────────────────

	{
		Slug:        "cli-task-lifecycle",
		Title:       "Task Lifecycle Commands",
		Category:    CategoryCLIExample,
		Description: "Create, claim, work, and complete a task",
		Tags:        []string{"task"},
		Content: `# Task Lifecycle
agentic-agent task create --title "Add user auth" --description "JWT-based auth"
agentic-agent task list
agentic-agent task claim TASK-001
agentic-agent task show TASK-001
agentic-agent task complete TASK-001`,
	},

	{
		Slug:        "cli-track-workflow",
		Title:       "Track Workflow Commands",
		Category:    CategoryCLIExample,
		Description: "Init a track, refine spec, activate with plan and tasks",
		Tags:        []string{"track", "plan"},
		Content: `# Track Workflow: Idea to Implementation
agentic-agent track init "User Dashboard" --type feature
agentic-agent track refine user-dashboard
agentic-agent track activate user-dashboard --decompose
agentic-agent plan show --track user-dashboard
agentic-agent plan next --track user-dashboard
agentic-agent plan mark plan.md 12 done`,
	},

	{
		Slug:        "cli-context-commands",
		Title:       "Context Management Commands",
		Category:    CategoryCLIExample,
		Description: "Generate context, scan for missing files, build bundles",
		Tags:        []string{"context"},
		Content: `# Context Management
agentic-agent context generate internal/auth    # Generate context.md
agentic-agent context scan                      # Find dirs missing context
agentic-agent context build --task TASK-001     # Full context bundle`,
	},

	{
		Slug:        "cli-skills-setup",
		Title:       "Skills Setup Commands",
		Category:    CategoryCLIExample,
		Description: "Install skill packs, ensure setup, check drift",
		Tags:        []string{"skills"},
		Content: `# Skills Setup
agentic-agent skills list                             # Available packs
agentic-agent skills install tdd --tool claude-code   # Install for one tool
agentic-agent skills install tdd --tool claude-code,cursor,gemini  # Multi-tool
agentic-agent skills ensure                           # Auto-detect and ensure
agentic-agent skills check                            # Check for drift`,
	},

	{
		Slug:        "cli-openspec-lifecycle",
		Title:       "OpenSpec Change Lifecycle",
		Category:    CategoryCLIExample,
		Description: "Init, import, execute, complete, and archive a change",
		Tags:        []string{"openspec"},
		Content: `# OpenSpec Lifecycle
agentic-agent openspec init "Auth Feature" --from .agentic/spec/auth.md
# Edit proposal.md and tasks.md...
agentic-agent openspec import auth-feature
agentic-agent openspec status auth-feature
agentic-agent openspec complete auth-feature
agentic-agent openspec archive auth-feature`,
	},

	{
		Slug:        "cli-project-status",
		Title:       "Project Status and Validation",
		Category:    CategoryCLIExample,
		Description: "Check project health, run validators, view dashboard",
		Tags:        []string{"status", "validate"},
		Content: `# Project Health
agentic-agent status                    # Dashboard with progress bar
agentic-agent status --format json      # Machine-readable
agentic-agent validate                  # Run all validation rules`,
	},

	// ── Workflow Recipes ────────────────────────────────────────────────

	{
		Slug:        "recipe-tdd-workflow",
		Title:       "TDD Workflow (Red-Green-Refactor)",
		Category:    CategoryWorkflowRecipe,
		Description: "Full TDD workflow with skill pack and task decomposition",
		Tags:        []string{"tdd", "skills"},
		Content: `# TDD Workflow Recipe
# Prerequisites: TDD skill pack installed

# 1. Install TDD skill pack
agentic-agent skills install tdd --tool claude-code

# 2. Create and claim a task
agentic-agent task create --title "Add input validation"
agentic-agent task claim TASK-001

# 3. Send to agent:
# "Implement TASK-001 following strict TDD.
#  RED: Write failing tests for all acceptance criteria.
#  GREEN: Write minimum code to pass each test.
#  REFACTOR: Clean up without changing behavior.
#  Run tests after each phase."

# 4. Complete
agentic-agent task complete TASK-001`,
	},

	{
		Slug:        "recipe-spec-to-tasks",
		Title:       "Spec to Tasks Pipeline",
		Category:    CategoryWorkflowRecipe,
		Description: "Turn a specification into tracked, implementable tasks",
		Tags:        []string{"spec", "track", "plan"},
		Content: `# Spec to Tasks Pipeline

# 1. Create a track from your idea
agentic-agent track init "Payment Processing" --type feature

# 2. Send to agent:
# "Read .agentic/tracks/payment-processing/spec.md
#  Flesh out the specification with:
#  - Detailed requirements
#  - API contracts
#  - Edge cases
#  - Acceptance criteria for each requirement"

# 3. Activate track to generate plan and tasks
agentic-agent track activate payment-processing --decompose

# 4. Work through tasks sequentially
agentic-agent plan next --track payment-processing
agentic-agent task claim TASK-001

# 5. Send to agent:
# "Continue task TASK-001. Read the spec and plan.
#  Implement according to the spec. Run tests."

# 6. Mark progress
agentic-agent task complete TASK-001
agentic-agent plan mark plan.md 5 done`,
	},

	{
		Slug:        "recipe-multi-agent",
		Title:       "Multi-Agent Collaboration",
		Category:    CategoryWorkflowRecipe,
		Description: "Switch between AI tools on the same project",
		Tags:        []string{"multi-agent", "skills"},
		Content: `# Multi-Agent Collaboration Recipe

# 1. Ensure skills for all your tools
agentic-agent skills ensure --all

# 2. Start work with Claude Code
agentic-agent task claim TASK-001
# Agent prompt: "Implement the API endpoints for TASK-001"

# 3. Switch to Cursor for UI work
agentic-agent task continue TASK-001
# Agent prompt: "Continue TASK-001. The API is done. Build the UI components."

# 4. Switch to Gemini for documentation
agentic-agent task continue TASK-001
# Agent prompt: "Continue TASK-001. Write API docs and update README."

# 5. Complete from any tool
agentic-agent task complete TASK-001`,
	},

	{
		Slug:        "recipe-new-project",
		Title:       "New Project Bootstrap",
		Category:    CategoryWorkflowRecipe,
		Description: "Set up a new project from scratch with full workflow",
		Tags:        []string{"init", "task"},
		Content: `# New Project Bootstrap

# 1. Initialize
agentic-agent init --name "my-project"
agentic-agent skills ensure

# 2. Create your first spec
# Agent prompt: "Create a specification at .agentic/spec/mvp.md
#  for the core MVP features of this project."

# 3. Create tasks from the spec
agentic-agent task create --title "Set up project structure" --spec-refs mvp.md
agentic-agent task create --title "Implement core data models" --spec-refs mvp.md
agentic-agent task create --title "Add API endpoints" --spec-refs mvp.md
agentic-agent task create --title "Write tests" --spec-refs mvp.md

# 4. Start working
agentic-agent task claim TASK-001
# Agent prompt: "Claim task TASK-001 and implement it."

# 5. Check progress
agentic-agent status`,
	},

	{
		Slug:        "recipe-atdd",
		Title:       "Acceptance Test-Driven Development",
		Category:    CategoryWorkflowRecipe,
		Description: "Write acceptance tests from task criteria before implementation",
		Tags:        []string{"atdd", "testing"},
		Content: `# ATDD Recipe

# 1. Install ATDD skill pack
agentic-agent skills install atdd --tool claude-code

# 2. View task acceptance criteria
agentic-agent task show TASK-001

# 3. Send to agent:
# "For TASK-001, follow ATDD:
#  1. Read the acceptance criteria from the task
#  2. Write one acceptance test per criterion (all should FAIL)
#  3. Implement code to make each test pass, one at a time
#  4. Refactor after all tests pass
#  5. Run the full test suite to confirm nothing broke"

# 4. Complete
agentic-agent task complete TASK-001`,
	},

	{
		Slug:        "recipe-idea-to-code",
		Title:       "Idea to Code Pipeline",
		Category:    CategoryWorkflowRecipe,
		Description: "Full pipeline: brainstorm → PRD → openspec → implement",
		Tags:        []string{"brainstorming", "prd", "openspec", "pipeline"},
		Content: `# Idea to Code Pipeline
# Complete workflow from rough idea to working code

# 1. Brainstorm the idea
# Agent prompt: "I have an idea: <your idea>.
#  Brainstorm this with me. Ask questions to understand the problem,
#  users, and constraints. Propose 2-3 approaches with trade-offs."

# 2. Write a PRD
# Agent prompt: "Write a PRD for the approach we picked.
#  Include user stories, acceptance criteria, and success metrics."

# 3. Create openspec from PRD
agentic-agent openspec init "my-feature" --from .agentic/spec/prd-my-feature.md

# 4. Review the generated proposal and tasks
agentic-agent openspec status my-feature

# 5. Execute the tasks
# Agent prompt: "Execute the approved openspec change: my-feature.
#  Work through tasks sequentially with tests."

# 6. Verify
agentic-agent validate
agentic-agent status`,
	},

	{
		Slug:        "recipe-ralph-execution",
		Title:       "Ralph Wiggum Iterative Execution",
		Category:    CategoryWorkflowRecipe,
		Description: "Implement openspec tasks using Ralph loops for convergence",
		Tags:        []string{"ralph", "openspec", "iterate"},
		Content: `# Ralph Wiggum Iterative Execution
# Use Ralph loops for reliable implementation with convergence guarantees

# 1. Ensure openspec change is approved
agentic-agent openspec status my-feature
agentic-agent task list

# 2. Send to agent:
# "Ralph my tasks. For each task in the backlog:
#  1. Claim the task
#  2. Build a focused implementation prompt
#  3. Run a Ralph loop: implement → test → fix → converge
#  4. Complete when all tests pass
#  5. Move to the next task"

# 3. Monitor progress
agentic-agent status
agentic-agent plan show my-feature`,
	},

	{
		Slug:        "recipe-prd-to-tasks",
		Title:       "PRD to Backlog Pipeline",
		Category:    CategoryWorkflowRecipe,
		Description: "Convert a PRD into structured backlog tasks ready for execution",
		Tags:        []string{"prd", "tasks", "pipeline"},
		Content: `# PRD to Backlog Pipeline

# 1. Generate a PRD (or use existing)
# Agent prompt: "Write a PRD for: <your feature>.
#  Save to .agentic/spec/prd-<feature>.md"

# 2. Convert PRD to openspec change
agentic-agent openspec init "<feature>" --from .agentic/spec/prd-<feature>.md

# 3. Review the generated tasks
agentic-agent task list
agentic-agent openspec status <feature>

# 4. Or convert directly to task YAML
# Agent prompt: "Convert the PRD at .agentic/spec/prd-<feature>.md
#  into tasks. Each task should be one-iteration size with
#  acceptance criteria and dependency ordering."

# 5. Start working
agentic-agent task claim TASK-001`,
	},
}
