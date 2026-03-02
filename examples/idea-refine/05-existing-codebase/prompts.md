# Scenario 5: Prompts for Existing Codebase

## Overview

This document contains copy-paste prompts for two situations:
1. **Adding a new feature** to existing code (CSV export example)
2. **Modifying existing behavior** (filtering logic change example)

Use `scenario.md` to understand which approach applies to your situation.

---

## Part A: Adding a New Feature (CSV Export)

### Phase 1: Brainstorm

```
I want to add CSV export to our project management platform. The idea is to let
project managers export their projects to CSV for reporting and analysis.

Can you help me brainstorm this feature? I want to understand:
- Who needs this and why
- What data should be included
- Where it fits in the UI
- What edge cases might come up
- Whether we're missing any technical constraints

Please ask me clarifying questions.
```

---

### Phase 2: Product-Wizard

After you answer the brainstorm questions, use this prompt:

```
Now use the product-wizard skill to create a PRD for CSV export based on our
brainstorm conversation.

The PRD should include:
- Problem statement
- Target users and personas
- Success metrics
- Acceptance criteria
- Technical constraints
- Scope: what's IN and OUT for MVP

Please make it comprehensive and production-ready.
```

---

### Phase 3: OpenSpec Init

After the PRD is created:

```
Use openspec to initialize a change from this PRD. Call it "csv-export".

agentic-agent openspec init "CSV Export Feature" --from [path-to-prd-file]

Then create:
1. proposal.md — fill in problem, approach, scope, acceptance criteria
2. tasks.md — break into 3-4 implementation tasks with clear dependencies
```

---

### Phase 4: Implement Each Task

For each task returned by `agentic-agent task list`:

```
Claim the task:
agentic-agent task claim TASK-001

Then tell me (Claude):
"Implement this task. Use the existing [feature/pattern] as a reference if needed."

When done:
agentic-agent task complete TASK-001
```

---

### Phase 5: Complete the Change

```bash
# Verify all tasks are done
agentic-agent task list

# Close the change
agentic-agent openspec complete csv-export

# Archive it
agentic-agent openspec archive csv-export
```

---

## Part B: Modifying Existing Behavior (Filtering)

### Situation: Changing from AND logic to OR logic

This is a **breaking change** — it requires migration planning.

---

### Phase 1: Brainstorm the Migration

```
We need to change our task filtering from AND logic to OR logic. This is a
breaking change because:
- Existing saved filters will behave differently
- Users' queries will return different results
- We need a migration strategy

Can you help me brainstorm:
1. How to migrate existing saved filters (AND → OR)
2. How to roll this out safely (feature flag strategy)
3. How to verify it doesn't break existing queries
4. How to communicate this to users
5. Rollback plan if something goes wrong
```

---

### Phase 2: Product-Wizard for the Complete Picture

```
Use the product-wizard skill to create a PRD for "Filtering Logic Migration"
that covers:

1. Why we're changing (problem statement)
2. What breaks (impact analysis)
3. Migration plan (how we convert data)
4. Rollout strategy (feature flag phases)
5. Testing plan (how we verify nothing breaks)
6. Rollback plan (how we undo if needed)
7. Timeline (how long this takes)
8. User communication (when/how we tell users)

Make it a complete specification that covers the complexity of this change.
```

---

### Phase 3: OpenSpec with Decomposed Tasks

```
Use openspec to initialize a change from this PRD. Call it "filtering-or-logic".

agentic-agent openspec init "Filtering Logic Migration" --from [prd-file]

Then create tasks.md with these phases:
1. Write and test data migration script
2. Update query builder to support OR logic
3. Add feature flag to toggle old/new behavior
4. Update API contracts and docs
5. Update UI to reflect new behavior
6. Write comprehensive test suite
7. Deploy with flag OFF (safety first)
8. Monitor and gather metrics
9. Flip flag in production (staged rollout)
```

---

### Phase 4: Execute with Safety Checks

For each task:

```bash
agentic-agent task claim TASK-001
```

Tell me:
```
"Implement this task. This is a migration, so:
1. Verify migrations work on test data
2. Check queries return expected results
3. Ensure rollback is possible
4. Document any manual steps"
```

After each task:
```bash
agentic-agent validate          # Run all checks
agentic-agent task complete TASK-001
```

---

### Phase 5: Close with Confidence

```bash
# Verify all tasks done
agentic-agent task list

# Close the change
agentic-agent openspec complete filtering-or-logic

# Archive
agentic-agent openspec archive filtering-or-logic
```

---

## Part C: Quick Prompts by Situation

### "I have a small bug to fix"

```
Change [X] in [file/component] to [fix description].

Use the existing [pattern/approach] as reference.
```

**No openspec needed.** I'll:
1. Read the file
2. Edit it
3. Verify it works

---

### "I want to add a simple feature (1–3 files)"

```
Add [feature] to [component/module].

It should:
- [Behavior 1]
- [Behavior 2]
- [Acceptance criterion 1]
- [Acceptance criterion 2]

Use the existing [related feature] as a reference for patterns.
```

**No openspec needed.** I'll read code, implement, validate.

---

### "I want to add a significant feature (4+ files)"

```
I want to add [feature] to this project. Can we brainstorm this first?

Context:
- Who needs this: [description]
- Why they need it: [problem]
- Current state: [how they do it now]
- Success would be: [outcome]

Can you use the product-wizard skill to create a PRD?
```

**Then use openspec** to structure and track implementation.

---

### "I need to change something with side effects"

```
I need to change how [feature] works. This affects [list of systems].

Current behavior: [describe]
Desired behavior: [describe]
Breaking changes: [list what breaks]
Migration needed: [yes/no, describe]

Use the product-wizard skill to create a complete PRD covering:
1. Why we're changing
2. What breaks
3. How we migrate/roll out
4. How we verify it works
5. How we roll back if needed
```

**Use openspec** to ensure all side effects are planned.

---

### "I want to resume work on an in-progress change"

```
Continue implementing change [change-name].

Show me:
1. What's done
2. What's next
3. Any blockers

Then:
agentic-agent task claim [next-unclaimed-task]
```

---

## Part D: Essential CLI Commands

### Check Project Health

```bash
agentic-agent status
# Output: Progress, backlog size, in-progress count

agentic-agent task list
# Output: All tasks, status, spec references
```

### Start a Change (Spec-Driven)

```bash
agentic-agent openspec init "Feature Name" --from path/to/spec.md
# Creates: proposal.md, tasks.md in .agentic/openspec/changes/

agentic-agent openspec status feature-name
# Shows: progress, tasks, blockers
```

### Implement Task by Task

```bash
agentic-agent task claim TASK-001
# Records: branch, timestamp, locks task

agentic-agent context generate <directory>
# Generates context bundle for the directory before editing

agentic-agent validate
# Runs: linters, tests, validators before completing

agentic-agent task complete TASK-001
# Records: commits since claim, marks done
```

### Finish the Change

```bash
agentic-agent openspec complete feature-name
# Verifies: all tasks done, runs final validation

agentic-agent openspec archive feature-name
# Archives: moves to historical record
```

### Browse Workflows

```bash
agentic-agent prompts list
# Shows: all built-in prompts and recipes

agentic-agent prompts show recipe-idea-to-code
# Full: brainstorm → PRD → openspec → implement workflow

agentic-agent prompts show cli-openspec-lifecycle
# Full: init → import → execute → complete → archive
```

---

## Part E: Decision Tree

```
START: I want to change something in the codebase

  1. How many files?
     ├─ 1 file → DIRECT EDIT
     ├─ 2-3 files → DIRECT EDIT
     └─ 4+ files → Go to 2

  2. Is this a breaking change?
     ├─ No → OPENSPEC (for traceability)
     └─ Yes → Go to 3

  3. Does it affect multiple layers/systems?
     ├─ No → DIRECT EDIT (small breaking change)
     └─ Yes → OPENSPEC (mandatory for safety)

  DIRECT EDIT:
  → Tell me what to change
  → I read code, edit, validate
  → Done in 5-30 min

  OPENSPEC:
  → brainstorming skill
  → product-wizard skill
  → agentic-agent openspec init
  → task claim → implement → complete
  → Done in 1-8 hours
```

---

## Troubleshooting

### "I'm not sure if I should use openspec"

**Quick rule:** If you'd explain this change in more than 2 sentences, use openspec.

Examples:
- ✅ "Change the color of the button" → Direct edit (1 sentence)
- ✅ "Add validation to the form" → Direct edit (1-2 sentences)
- ❌ "Change how filtering works across the platform" → OpenSpec (requires migration, feature flag, communication)
- ❌ "Add multi-currency support" → OpenSpec (4+ files, multiple systems, significant change)

---

### "I started, now what?"

If you've run `agentic-agent openspec init` and created tasks:

```bash
# See what you created
agentic-agent openspec status <name>

# Start working
agentic-agent task claim <first-task>
agentic-agent context generate <directory>

# Tell me: "Implement this task"
# After done:
agentic-agent task complete <first-task>

# Repeat for each task
```

---

### "How do I know I'm done?"

```bash
# All tasks should show "done"
agentic-agent task list | grep <your-change-name>

# Then:
agentic-agent openspec complete <name>
agentic-agent openspec archive <name>
```

If validation fails:
```bash
agentic-agent validate
# Fix the issues
agentic-agent task complete <task-id>
```

---

## Real Example: Conversation Flow

```
You: "I want to add CSV export to the platform."

Me: "Let's brainstorm this. Who needs it? Why?"

You: [answer brainstorm questions]

Me: "Use product-wizard to create a PRD..."

You: "Create PRD"

Me: [PRD is generated]

You: "Now create openspec from this PRD for csv-export"

Me: "agentic-agent openspec init 'CSV Export' --from prd.md"

Me: [creates proposal.md and tasks.md]

You: "agentic-agent task claim TASK-001"

You: "Implement the service layer"

Me: [I implement the service]

You: "agentic-agent task complete TASK-001"

[Repeat for each task...]

You: "agentic-agent openspec complete csv-export"

Me: "Change archived."
```

**Total time:** 2-3 hours for a complete, tracked feature.

---

See `scenario.md` for detailed business examples.

See `learnings.md` for key lessons.
