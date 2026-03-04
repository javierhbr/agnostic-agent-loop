# Learnings: Working in an Existing Codebase

Key insights from adding features and making changes in established projects.

---

## Core Principles

### 1. Always Read Existing Code First

Before proposing changes, I need to understand:
- How similar features are currently built
- What patterns are established
- Where code lives and how it's organized
- What dependencies already exist

**Why:** Writing code that matches the existing style prevents friction and makes merging easier.

**How:** When you tell me to make a change, I'll ask to read the relevant files first. Or better yet, run:
```bash
agentic-agent context generate <directory>
```

This gives me a bundle with the codebase structure, which I use as reference while implementing.

---

### 2. Use OpenSpec for Anything Multi-Layer or Breaking

| Situation | Approach |
|-----------|----------|
| Typo fix | Direct edit |
| Single file change | Direct edit |
| 1-3 files, same layer | Direct edit |
| 4+ files, multiple layers | **OpenSpec** |
| Behavior change affecting other systems | **OpenSpec** |
| Breaking change (API, data model, logic) | **OpenSpec** |
| Uncertain scope | **OpenSpec** (better safe) |

**Why:** OpenSpec forces you to plan, think through side effects, and document decisions. It looks like overhead at first, but it prevents months of pain later.

**Example:** "Add CSV export" sounds like 1 file. It's actually 4: service, endpoint, UI, tests. If you discover side effects mid-implementation, you've already done the planning work.

---

### 3. Decision Rule: Small Changes vs Full Pipeline

**Small change (< 30 min, single concern):**
```
You: "Fix the button color in the dashboard"
Me: I read Dashboard.tsx, change it, verify it looks right
You: Done
```

**Medium change (1-2 hours, 2-3 files):**
```
You: "Add a notification icon to the header"
Me: I read Header.tsx, Notifications service, maybe styles
   implement icon, integrate with service, test
You: Done
```

**Large change (2-8 hours, 4+ files, multi-layer):**
```
You: "Add CSV export"
Me: 1. Brainstorm what it means
    2. PRD with requirements
    3. OpenSpec with 4 tasks
    4. Implement service → API → UI → tests
    5. Validate everything works together
You: Done with full traceability
```

The line is fuzzy — when in doubt, **use OpenSpec**. It's better to have structure you don't need than to discover mid-implementation you need it.

---

### 4. Always Run Context Before Editing a Directory

Before I start editing files in a directory, you should run:
```bash
agentic-agent context generate <directory>
```

This does two things:
1. Generates a context bundle showing me the structure, patterns, and dependencies
2. Helps me write code that fits naturally into the existing codebase

**Example:**
```bash
agentic-agent context generate ./services/auth
# Now I understand: where auth code lives, how it's organized,
# what patterns it uses, what tests exist
```

---

### 5. Validate Before Completing

Before marking a task complete, always run:
```bash
agentic-agent validate
```

This runs quality checks: linting, tests, type checking, validators.

**Why:** Catches issues early before code reaches main branch.

**What to expect:**
```bash
$ agentic-agent validate

[✓] Lint: No issues
[✓] Types: All checks pass
[✓] Tests: 23 passing
[✓] Custom rules: All pass

Ready to complete.
```

If validation fails:
```
[✗] Tests: 2 failing
  - export.test.ts: exports with special characters
  - filters.test.ts: OR logic migration test

Fix these issues, then run validate again.
```

---

## Common Patterns

### Pattern 1: The Breaking Change

**Recognize it:** Your change breaks existing API/data contracts.

**Examples:**
- Changing database schema
- Renaming API endpoints
- Changing parameter types
- Changing enum values
- Modifying filter logic (AND → OR)

**What to do:**
1. Use **openspec** (mandatory)
2. Include migration plan (convert old data to new format)
3. Use feature flag (roll out safely)
4. Write comprehensive tests (verify nothing breaks)
5. Plan rollback (how to undo if needed)

**Never:**
- Assume "I'll migrate old data later" (you won't)
- Skip tests for breaking changes
- Push breaking changes without a migration script
- Deploy without a rollback plan

---

### Pattern 2: The Cross-Layer Change

**Recognize it:** Your change touches database → service → API → UI.

**Examples:**
- Adding a new permission system
- Changing how data is fetched/cached
- Adding a new entity type (users, projects, tasks)
- Changing authentication flow

**What to do:**
1. Use **openspec** (mandatory — this is complex)
2. Break into 5-7 focused tasks (database, service, API, UI, integration, tests, docs)
3. Implement in order: bottom layer up
4. Test each layer separately, then together
5. Use feature flags for phased rollout

**Why:** These changes have hidden dependencies. Planning upfront reveals them.

---

### Pattern 3: The Performance Change

**Recognize it:** You're optimizing queries, caching, or reducing work.

**Examples:**
- Adding database indexes
- Caching expensive computations
- Changing pagination strategy
- Optimizing a hot path

**What to do:**
1. Measure before and after (real numbers)
2. Document the bottleneck (where is time going?)
3. Test edge cases (does cache invalidate correctly?)
4. Verify it's actually faster (don't optimize by feel)

**Prompt:**
```
"I want to optimize [what]. Currently [measurement]. Goal is [target measurement].
Can you review the current implementation and suggest where to optimize?"
```

---

### Pattern 4: The Refactoring

**Recognize it:** You're reorganizing code without changing behavior.

**Examples:**
- Breaking a large file into smaller ones
- Renaming things for clarity
- Extracting shared logic
- Reorganizing imports

**Rule:** No openspec needed if behavior doesn't change.

**But:** If you're unsure whether behavior will change, **use openspec** (safer).

**Prompt:**
```
"Can you refactor [description] while preserving behavior?
Here's the current code: [file]

Goals:
- [clarity/maintainability goal]
- [organization goal]

Constraints:
- Keep behavior identical
- No API changes
- All tests must pass
"
```

---

## Anti-Patterns (Don't Do These)

### ❌ "I'll skip planning and just code"

**Result:** Discover halfway through you need to change the API, so you rewrite half the code.

**Fix:** Always brainstorm → PRD → plan for anything non-trivial.

---

### ❌ "Let me hack in the feature, clean up later"

**Result:** The "cleanup" never happens. Technical debt accumulates.

**Fix:** Get it right the first time. Use openspec to ensure proper structure.

---

### ❌ "This breaking change won't hurt anyone"

**Result:** 10 customers complain. You rush a fix. The fix has bugs.

**Fix:** Treat breaking changes as critical. Always:
1. Document what breaks
2. Provide migration path
3. Use feature flag
4. Communicate to users

---

### ❌ "The tests pass locally, good enough"

**Result:** Works on your machine, breaks in production (different environment, data, concurrency).

**Fix:** Run `agentic-agent validate` which tests against the real environment.

---

### ❌ "I'll write tests after"

**Result:** Tests never get written. Code behavior becomes undocumented.

**Fix:** Write tests as part of the task. OpenSpec enforces this by making tests a required task.

---

## Workflow Tips

### Tip 1: Use Feature Flags for Any Risky Change

```typescript
if (config.features.csvExport) {
  // New export logic
} else {
  // Old export logic
}
```

Allows you to:
- Deploy code without enabling it
- Test in production without affecting users
- Roll back instantly (flip flag off)
- Measure impact before full rollout

---

### Tip 2: Read the Tests

When understanding existing code, **read the tests first**.

Tests show:
- What the code is supposed to do
- What edge cases exist
- How other code uses this code
- What's considered correct behavior

---

### Tip 3: Ask Questions Before Coding

If requirements seem incomplete, **ask** before I start:

```
"Before I start, I want to confirm:
1. [Clarification 1]
2. [Clarification 2]
3. [Edge case question]

Should I handle [scenario]?"
```

Better to spend 5 min clarifying than 2 hours implementing the wrong thing.

---

### Tip 4: Use Git to Track Progress

Each task should produce commits:
```bash
agentic-agent task claim TASK-001
# [I implement]
agentic-agent task complete TASK-001
# Records all commits since claim time
```

This creates an audit trail and makes it easy to review what changed.

---

## Success Checklist

Before marking a task complete:

- [ ] Code is written
- [ ] All tests pass (`agentic-agent validate`)
- [ ] Code matches existing patterns (checked with context)
- [ ] No console errors or warnings
- [ ] Changes documented (if user-facing)
- [ ] Breaking changes documented (if applicable)
- [ ] Commit message is clear ("what" + "why")

Before closing a change:

- [ ] All tasks completed
- [ ] All tests pass
- [ ] Code reviewed (for larger changes, request review)
- [ ] Documentation updated
- [ ] Changelog updated
- [ ] Ready to merge to main branch

---

## Getting Unstuck

### "I started implementing but the scope grew"

**Don't:** Keep going and hope it works out.

**Do:** Stop, reassess. Tell me:
```
"I started on [task], but discovered [new work].
Original scope: [description]
New scope: [description]
Should I: a) Expand this task? b) Create separate tasks? c) Descope some items?"
```

I'll help you re-scope and stay on track.

---

### "Tests are failing and I don't know why"

**Don't:** Disable the tests or work around them.

**Do:** Let me debug. Tell me:
```
"Test [name] is failing. Expected [X], got [Y].
Here's the test: [show test code]
Here's my implementation: [show code]
"
```

---

### "I think I broke something but I'm not sure"

**Do:** Run:
```bash
agentic-agent validate
```

This tells you exactly what's broken. Show me the output:
```
"Validation failed. Here's the error: [output]"
```

---

## More Resources

- **scenario.md** — Real business examples and decision points
- **prompts.md** — Copy-paste prompts for each situation
- Parent **QUICK-REFERENCE.md** — 1-page cheatsheet
- Parent **SKILLS-GUIDE.md** — Deep dive into skills

---

**Remember:** The best code changes are the ones that are well-planned, well-tested, and well-documented. Take the time upfront to avoid pain later.
