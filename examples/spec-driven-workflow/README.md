# Spec-Driven Workflow: Spec Kit, OpenSpec & Autopilot

This example demonstrates how to use specification files from Spec Kit, OpenSpec, or native agentic specs — all resolved automatically through multi-directory configuration. It also shows autopilot mode for sequential task processing.

---

## How It Works

The key insight: **you don't import specs — you reference them.** Configure `specDirs` in `agnostic-agent.yaml` to point at your spec directories. The resolver searches them in order and uses the first match.

```yaml
# agnostic-agent.yaml
paths:
  specDirs:
    - .specify/specs       # Spec Kit
    - openspec/specs        # OpenSpec
    - .agentic/spec         # Agentic native (fallback)
```

When a task references `auth-requirements.md`, the resolver searches:
1. `.specify/specs/auth-requirements.md`
2. `openspec/specs/auth-requirements.md`
3. `.agentic/spec/auth-requirements.md` — found here

---

## Part A: Spec Kit Integration

### 1. Plan with Spec Kit

Spec Kit generates specs in `.specify/specs/`. After planning:

```
.specify/
└── specs/
    └── 001-auth/
        └── spec.md       # Feature spec with requirements and scenarios
```

### 2. Create tasks referencing Spec Kit specs

```bash
agentic-agent task create \
  --title "Implement JWT token service" \
  --spec-refs "001-auth/spec.md" \
  --outputs "internal/auth/jwt.go" \
  --acceptance "Token generation works,Validation rejects expired tokens"
```

The spec ref `001-auth/spec.md` resolves to `.specify/specs/001-auth/spec.md` because `.specify/specs` is first in `specDirs`.

### 3. Verify spec resolution

```bash
# List all specs across configured directories
agentic-agent specify list
```

```
auth-requirements.md  /path/to/.agentic/spec/auth-requirements.md
api-design.md         /path/to/.agentic/spec/api-design.md
001-auth/spec.md      /path/to/.specify/specs/001-auth/spec.md
auth/spec.md          /path/to/openspec/specs/auth/spec.md
```

```bash
# Resolve and read a specific spec
agentic-agent specify resolve "001-auth/spec.md"
```

Outputs the full content of the Spec Kit spec.

### 4. Claim task with readiness checks

```bash
agentic-agent task claim SPEC-001
```

```
Task SPEC-001: READY
  [+] spec-resolvable: spec "auth-requirements.md" resolved at .agentic/spec/auth-requirements.md
  [+] spec-resolvable: spec "api-design.md" resolved at .agentic/spec/api-design.md
Claimed task SPEC-001
```

Readiness checks verify that all referenced specs exist before claiming. If a spec is missing, the claim still proceeds but you see a warning.

### 5. Build context bundle with specs

```bash
agentic-agent context build --task SPEC-001
```

The output includes a `specs:` section with the full content of each resolved spec. The agent receives everything it needs — task definition, project context, and specification content — in a single bundle.

### 6. Implement and complete

```bash
# Work on the task...
agentic-agent validate
agentic-agent task complete SPEC-001
```

---

## Part B: OpenSpec Integration

### 1. Plan with OpenSpec

OpenSpec generates specs in `openspec/specs/`. After running `/opsx:new add-auth` and `/opsx:ff`:

```
openspec/
└── specs/
    └── auth/
        └── spec.md       # OpenSpec spec with proposal, design, tasks
```

### 2. Add OpenSpec directory to config

Already configured in `agnostic-agent.yaml`:

```yaml
paths:
  specDirs:
    - .specify/specs       # Spec Kit
    - openspec/specs        # OpenSpec ← specs found here
    - .agentic/spec         # Fallback
```

### 3. Create tasks referencing OpenSpec specs

```bash
agentic-agent task create \
  --title "Create User model and repository" \
  --spec-refs "auth/spec.md" \
  --outputs "internal/models/user.go,internal/repository/user_repo.go"
```

The ref `auth/spec.md` resolves to `openspec/specs/auth/spec.md`.

### 4. Execute the same way

```bash
# Verify the spec resolves
agentic-agent specify resolve "auth/spec.md"

# Claim with readiness checks
agentic-agent task claim TASK-001

# Build context (includes OpenSpec content)
agentic-agent context build --task TASK-001

# Work, validate, complete
agentic-agent validate
agentic-agent task complete TASK-001
```

### 5. Verify with OpenSpec after completion

```
/opsx:verify    # OpenSpec validates against its specs
/opsx:archive   # Archive the completed change
```

For the full spec-driven development guide, see [docs/SPEC_DRIVEN_DEVELOPMENT.md](../../docs/SPEC_DRIVEN_DEVELOPMENT.md).

---

## Part C: OpenSpec CLI (End-to-End)

The `openspec` command group handles the full change lifecycle — from a requirements file to archived implementation. Tell your agent:

> "Start a project from .agentic/spec/auth-requirements.md following openspec"

The agent uses the `openspec` skill to run these commands automatically:

### 1. Initialize a change from requirements

```bash
agentic-agent openspec init "Auth Feature" --from .agentic/spec/auth-requirements.md
```

```
Created change: auth-feature
  → Proposal:  .agentic/openspec/changes/auth-feature/proposal.md
  → Tasks:     .agentic/openspec/changes/auth-feature/tasks.md
  → Specs:     .agentic/openspec/changes/auth-feature/specs/
  → Status:    draft

→ Fill in proposal.md, then write tasks in tasks.md.
→ Then run: agentic-agent openspec import auth-feature
```

The agent reads the generated proposal template, fills in Problem/Approach/Scope/Acceptance, then writes `tasks.md` with a numbered implementation plan.

### 2. Import tasks into the backlog

```bash
agentic-agent openspec import auth-feature
```

```
Imported 4 tasks from auth-feature
  • TASK-1739000001: [auth-feature] Create User model with bcrypt
  • TASK-1739000002: [auth-feature] Add JWT token service
  • TASK-1739000003: [auth-feature] Implement login/register endpoints
  • TASK-1739000004: [auth-feature] Write integration tests

→ Run: agentic-agent task claim TASK-1739000001
```

### 3. Execute tasks sequentially

```bash
# Check progress at any time
agentic-agent openspec status auth-feature

# Work through each task
agentic-agent task claim TASK-1739000001
# ... implement ...
agentic-agent task complete TASK-1739000001

agentic-agent task claim TASK-1739000002
# ... implement ...
agentic-agent task complete TASK-1739000002
# ... repeat for all tasks ...
```

### 4. Complete and archive

```bash
# Validates all tasks are done, writes IMPLEMENTED marker
agentic-agent openspec complete auth-feature

# Moves to archive
agentic-agent openspec archive auth-feature
```

### Example prompts for your agent

| Prompt | What happens |
| ------ | ------------ |
| "Start a project from requirements.md following openspec" | Full lifecycle: init → fill proposal → write tasks → import → execute → complete |
| "Implement the features in docs/payment-spec.md using openspec" | Same flow, different source file |
| "openspec status auth-feature" | Shows task progress for an existing change |
| "Continue implementing change auth-feature" | Resumes at the next unclaimed task |

---

## Part D: Autopilot Mode

Autopilot processes backlog tasks sequentially: readiness check, claim, generate context, build bundle.

### 1. Preview with dry run

```bash
agentic-agent autopilot start --dry-run
```

```
--- Iteration 1/10 ---
Next task: [SPEC-001] Create JWT token service
Task SPEC-001: READY
  [+] spec-resolvable: spec "auth-requirements.md" resolved at .agentic/spec/auth-requirements.md
  [+] spec-resolvable: spec "api-design.md" resolved at .agentic/spec/api-design.md
[DRY RUN] Would claim task SPEC-001 and generate context

--- Iteration 2/10 ---
Next task: [SPEC-002] Implement auth middleware
Task SPEC-002: READY
  [+] spec-resolvable: spec "auth-requirements.md" resolved at .agentic/spec/auth-requirements.md
[DRY RUN] Would claim task SPEC-002 and generate context
```

### 2. Run autopilot

```bash
# Process up to 3 tasks
agentic-agent autopilot start --max-iterations 3
```

Per iteration, autopilot:
1. Finds the next claimable task (prefers tasks where all readiness checks pass)
2. Prints readiness check results
3. Claims the task
4. Generates context for each scope directory
5. Builds a context bundle with resolved specs
6. Reports the task as ready for agent execution

### 3. Stop conditions

Autopilot stops when:
- All backlog tasks are processed
- `--max-iterations` limit is reached
- You press Ctrl+C

---

## Directory Structure

```
spec-driven-workflow/
├── README.md                          # This file
├── agnostic-agent.yaml                # Multi-dir spec config
├── .agentic/
│   ├── spec/
│   │   ├── auth-requirements.md       # Native spec (auth requirements)
│   │   └── api-design.md             # Native spec (API design)
│   ├── tasks/
│   │   ├── backlog.yaml              # 2 tasks referencing specs
│   │   ├── in-progress.yaml
│   │   └── done.yaml
│   └── context/
│       ├── global-context.md
│       └── rolling-summary.md
├── .specify/
│   └── specs/
│       └── 001-auth/
│           └── spec.md               # Spec Kit sample spec
└── openspec/
    └── specs/
        └── auth/
            └── spec.md               # OpenSpec sample spec
```

## Quick Reference

| Command | Purpose |
|---------|---------|
| `spec list` | List all specs across all configured directories |
| `spec resolve <ref>` | Resolve a spec ref and print its content |
| `task create --spec-refs "..."` | Create a task that references specifications |
| `task claim <id>` | Claim task (runs readiness checks first) |
| `context build --task <id>` | Build context bundle including resolved specs |
| `autopilot start --dry-run` | Preview autopilot without making changes |
| `autopilot start` | Run autopilot to process backlog tasks |

---

## Part E: Adding Features to an Existing Codebase

This section covers the workflow for developers working in an **established project** who need to add new features or modify existing ones. Unlike Parts A-D which assume you have specs/requirements in hand, this focuses on the decision-making and planning process when you're starting from an idea.

---

### Decision: When to Use OpenSpec vs Direct Edit

**Quick Rule:** If you'd need more than 2 sentences to describe the change, use OpenSpec.

| Situation | Approach | Time | Example |
| --------- | -------- | ---- | ------- |
| Bug fix, typo, single file | Direct edit | 5-30 min | Fix button color, typo in error message |
| Small feature (1–3 files, no breaking changes) | Direct edit | 15-60 min | Add notification icon, new button |
| New feature (4+ files, cross-layer) | **OpenSpec** | 2-4 hours | CSV export, new payment method |
| Modify existing behavior (breaking change) | **OpenSpec** | 3-8 hours | Change filtering logic, new permission system |
| Have PRD/spec ready | **OpenSpec init** | 1-3 hours | Start from requirements file |

---

### The Full Pipeline: From Idea to Shipped Code

#### Phase 1: Brainstorm (15 min)

Tell the AI:

```
I want to add [feature] to this project. Let's brainstorm first.

Context:
- Who needs this: [user type]
- Why they need it: [problem]
- Current state: [how they handle it today]
- Success would be: [desired outcome]

Can you ask me clarifying questions?
```

The AI asks about scope, edge cases, constraints, and integration points. You answer with specifics.

#### Phase 2: Write a PRD (20 min)

After brainstorming, tell the AI:

```
Use the product-wizard skill to create a PRD based on our brainstorm.

The PRD should include:
- Clear problem statement
- Target users and why they need this
- Success metrics (measurable)
- Acceptance criteria
- Technical constraints (time, budget, team size)
- In-scope vs out-of-scope for MVP
```

The AI generates a production-ready PRD.

#### Phase 3: Structure with OpenSpec (15 min)

Once you have a PRD:

```bash
agentic-agent openspec init "Feature Name" --from path/to/prd.md
```

This creates:

- `proposal.md` — filled in from PRD
- `tasks.md` — decomposed into 3-4 focused tasks

You can edit these files to adjust task boundaries, dependencies, or order.

#### Phase 4: Implement Task by Task (2-4 hours)

For each task returned by `agentic-agent task list`:

```bash
# Check what's next
agentic-agent task list | grep "in-progress\|pending"

# Start a task
agentic-agent task claim TASK-001

# Before editing a directory
agentic-agent context generate <directory>

# Tell Claude to implement
# "Implement this task. Use the existing [pattern] as reference."

# Verify work
agentic-agent validate

# Mark complete
agentic-agent task complete TASK-001
```

Repeat for all tasks.

#### Phase 5: Close the Change (5 min)

```bash
# Verify all tasks are done
agentic-agent openspec status feature-name

# Close
agentic-agent openspec complete feature-name

# Archive
agentic-agent openspec archive feature-name
```

---

### Essential Prompts (Copy-Paste)

**Small change (direct edit):**

```
Change [X] in [file/component] to do [Y].
Use the existing [pattern] as a reference for style/approach.
```

**New feature (brainstorm first):**

```
I want to add [feature] to this project. Let me brainstorm:

Who needs this: [description]
Why: [problem it solves]
Current state: [how it works now]
Success looks like: [outcome]

Ask me clarifying questions to understand this better.
```

**After PRD → Ready to implement:**

```
Use the product-wizard skill to create a PRD for this feature.

Include:
- Problem statement
- Target users
- Success metrics
- Acceptance criteria
- Constraints (time, budget, team)
- MVP scope (what's in/out)
```

**Initialize change from PRD:**

```bash
agentic-agent openspec init "CSV Export" --from .agentic/spec/prd-csv-export.md
```

**Resume in-progress work:**

```
Continue implementing change [feature-name].

Show me: What's done? What's next? Any blockers?
```

---

### Built-In Recipes and Prompts

Browse all available examples:

```bash
# List all prompts and recipes
agentic-agent prompts list

# Show the full idea-to-code pipeline
agentic-agent prompts show recipe-idea-to-code

# Show openspec lifecycle
agentic-agent prompts show cli-openspec-lifecycle

# Show specific workflow examples
agentic-agent prompts show openspec-execute
agentic-agent prompts show claim-and-implement
```

---

### Real Example: CSV Export Feature

**Scenario:** You're on a React/Node project. You want to let users export projects as CSV.

**Step 1: Brainstorm (you tell me)**

```
I want to add CSV export to the project dashboard.
Who: Project managers need to export for reporting
Why: They currently copy-paste into Excel
Current: Manual export doesn't exist
Success: 60% of power users use it in first month
```

**Step 2: I ask clarifying questions**

- What data to include? (All fields? Custom columns?)
- Export from where? (Dashboard? Project view?)
- Format? (Column order, date format?)
- Bulk export? (Multiple projects?)

**Step 3: You answer with specifics**

```
- Data: Task name, status, assignee, due date
- From: Dashboard and project view
- Format: Maintain dashboard column order
- Bulk: Not in MVP, phase 2
```

**Step 4: I create PRD with product-wizard**

```text
✅ Problem: Project managers export data for reporting; no built-in export
✅ Users: Project managers, team leads, analysts
✅ Success: 60% adoption, 70% of exports complete in <5 sec
✅ Acceptance criteria:
   - Export single project as CSV
   - Include 4 standard columns
   - Works Chrome/Firefox/Safari
   - File downloads to computer
✅ Scope: MVP doesn't include filters, custom columns, scheduling
```

**Step 5: I initialize OpenSpec**

```bash
agentic-agent openspec init "CSV Export" --from prd.md
```

Tasks created:

1. Create export service (backend logic)
2. Add API endpoint for export
3. Add UI button and download flow
4. Write integration tests

**Step 6: You implement each task**

```bash
agentic-agent task claim TASK-001
# You tell me: "Implement the export service"
# I write the code
agentic-agent task complete TASK-001

agentic-agent task claim TASK-002
# etc.
```

**Step 7: All done**

```bash
agentic-agent openspec complete csv-export
agentic-agent openspec archive csv-export
```

**Result:** Feature is implemented, tested, tracked, ready to merge.

**Total time:** ~3 hours from idea to shipped.

---

### Key Lessons from Real Projects

#### ✅ Do This

1. **Always read existing code before proposing changes**
   - Understand current patterns
   - Check how similar features are built
   - Identify integration points

2. **Use `agentic-agent context generate <dir>` before editing a directory**
   - Helps me understand structure
   - Ensures I write idiomatic code

3. **Run `agentic-agent validate` before completing**
   - Catches issues early
   - Verifies tests pass

4. **Use OpenSpec even for "small" features**
   - 4 files becomes 6 once you start
   - Better to have structure than regret later

5. **Use feature flags for breaking changes**
   - Safe rollout
   - Easy rollback
   - Measure impact before full deployment

#### ❌ Don't Do This

- Skip brainstorm/PRD for "obvious" features (you'll miss edge cases)
- Underestimate scope (always bigger than it looks)
- Change behavior without migration plan
- Bypass `validate` to "save time" (costs more time later)
- Commit without running validators
- Mix breaking changes with new features in same task

---

### Getting Unstuck

#### "I'm not sure if I need OpenSpec for this"

Answer these questions:

- How many files? (4+ → OpenSpec)
- Is it a breaking change? (Yes → OpenSpec)
- Does it affect other systems? (Yes → OpenSpec)
- Would I explain this in >2 sentences? (Yes → OpenSpec)

If any answer is yes, use OpenSpec.

#### "I started but the scope grew"

Stop and reassess:

```text
I started on [task], but discovered [new work].
Original: [description]
New: [description]
Should I: a) expand task? b) separate tasks? c) descope?
```

#### "Tests are failing"

Run diagnostics:

```bash
agentic-agent validate
# Shows exactly what's wrong
```

Tell me the output and I'll help fix it.

---

### Related Resources

- **[idea-refine/05-existing-codebase/](../../idea-refine/05-existing-codebase/)** — Complete business examples and decision trees
- **[idea-refine/QUICK-REFERENCE.md](../../idea-refine/QUICK-REFERENCE.md)** — 1-page cheatsheet for existing-codebase workflows
- **[SPEC_DRIVEN_DEVELOPMENT.md](../../docs/SPEC_DRIVEN_DEVELOPMENT.md)** — Deep dive into spec-driven methodology

---

## Related Documentation

- [Spec-Driven Development Guide](../../docs/SPEC_DRIVEN_DEVELOPMENT.md)
- [Idea Refinement Examples](../../idea-refine/)
