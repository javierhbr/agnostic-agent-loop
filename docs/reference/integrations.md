# Superpowers Plugin Integration Guide

Quick reference for using the Superpowers plugin alongside `agentic-agent` CLI.

**Status:** ✅ Integration complete (commit 70eb704)

---

## TL;DR

| Scenario | Use |
|----------|-----|
| 🔍 **Vague requirements** | Superpowers:brainstorming |
| 📋 **Clear spec** | `agentic-agent openspec init` |
| 🔴 **TDD (critical logic)** | superpowers:test-driven-development OR `tdd` skill |
| 🔧 **Debugging mystery** | superpowers:systematic-debugging |
| 🚀 **Iterative implementation** | superpowers:executing-plans OR `/ralph-loop` |
| ✅ **Final verification** | superpowers:verification-before-completion OR `sdd/verifier` |

---

## Installation Check

**Verify Superpowers is installed:**
```bash
# Check for Superpowers skills
ls ~/.claude/skills/ | grep "superpowers:"

# Or in Claude Code
/help → Marketplace → search "Superpowers"
```

**If missing, install:**
1. Open Claude Code
2. Run `/help` → **Marketplace**
3. Search for **Superpowers**
4. Click **Install**
5. Restart Claude Code

---

## Skills Mapping

### Superpowers Skills → In-House Equivalents

| Superpowers | CLI/In-House | When to Use |
|-------------|--------------|------------|
| `brainstorming` | `brainstorming` skill | Requirements unclear or need user discovery |
| `writing-plans` | `dev-plans` skill | Planning implementation phases |
| `test-driven-development` | `tdd` skill | Write tests first (payment, auth, critical logic) |
| `using-git-worktrees` | CLI native | Parallel isolated component work |
| `executing-plans` | `/ralph-loop` | Iterative task convergence |
| `systematic-debugging` | `systematic-debugging` skill | 4-phase root cause analysis |
| `verification-before-completion` | `sdd/verifier` skill | Hard gate requiring evidence |

**Key difference:** Superpowers is **strict methodology** (enforces gates, refuses to proceed without verification). CLI skills are **fast automation** (focuses on traceability and multi-agent support).

---

## Integrated Workflows

### Workflow 1: From Vague Idea to Shipped (with Superpowers)

```
/brainstorming
  ↓ (clarify requirements with user)
/product-wizard  (create PRD)
  ↓
/writing-plans  (create dev plan)
  ↓
agentic-agent openspec init "Feature" --from prd.md
  ↓
agentic-agent task claim TASK-1
  ↓
/executing-plans  (implement with checkpoints)
  ↓
agentic-agent task complete TASK-1
  ↓
/verification-before-completion  (hard gate)
  ↓
agentic-agent openspec complete
```

### Workflow 2: From Clear Spec to Shipped (with CLI)

```
agentic-agent openspec init "Feature" --from requirements.md
  ↓
agentic-agent task list
  ↓
agentic-agent task claim TASK-1
  ↓
/ralph-loop "prompt..." --max-iterations 10
  ↓
agentic-agent task complete TASK-1
  ↓
agentic-agent openspec complete
```

### Workflow 3: Debugging with Systematic Approach

```
Encounter bug during implementation
  ↓
/systematic-debugging
  ├─ Phase 1: Identify symptom
  ├─ Phase 2: Isolate problem
  ├─ Phase 3: Narrow root cause
  └─ Phase 4: Fix and test
  ↓
If tests pass → continue
If tests fail → restart /systematic-debugging
```

---

## Updated Skills

All of these skills now reference Superpowers as an **optional** complement:

### 1. openspec/SKILL.md
- **Phase 4 (Plan):** Can use `superpowers:writing-plans` or in-house `dev-plans`
- **Phase 6 (Execute):** Can use `superpowers:executing-plans`, TDD, or `/ralph-loop`
- **Phase 7 (Verify):** Can use `superpowers:verification-before-completion` or `sdd/verifier`

### 2. sdd/developer/SKILL.md
- **Step 3.5 (New):** Optional TDD with `superpowers:test-driven-development` or in-house `tdd`
- **Step 3.5 (New):** Optional isolation with `superpowers:using-git-worktrees`
- **Step 7 (New):** Optional debugging with `superpowers:systematic-debugging`

### 3. sdd/verifier/SKILL.md
- **Step 3.5 (New):** Optional verification gate with `superpowers:verification-before-completion`

### 4. run-with-ralph/SKILL.md
- **Step 5 (Updated):** Recommends `superpowers:executing-plans` or `/ralph-loop`
- **Debugging:** Can use `superpowers:systematic-debugging` during iteration

---

## The New superpowers-bridge Skill

Comprehensive integration guide covering:
- Installation + verification
- Comparison of both systems
- Decision tree (when to use what)
- 4 complete workflow scenarios
- Skill mapping reference
- Troubleshooting

**Access it:**
```bash
agentic-agent skills list  # see superpowers-bridge
agentic-agent skills install superpowers-bridge  # optional
# Then read superpowers-bridge/SKILL.md in Claude Code
```

---

## Detection Module

New `internal/skills/superpowers.go` provides:

```go
DetectSuperpowers() bool                    // Check if plugin installed
SuperpowersInstallInstructions() string     // Installation guidance
WarnIfSuperpowersNotFound()                 // Display warning message
```

**Future integration:** These functions will be called in `run` and `skills ensure` commands to warn users when Superpowers would be beneficial.

---

## Philosophy

**Superpowers (strict methodology):**
- ✅ Brainstorming with human approval gates
- ✅ TDD enforcement (can't pass verification without tests)
- ✅ Systematic debugging (4-phase protocol)
- ✅ Hard stops before marking done
- ✅ Best for: unclear requirements, critical logic, complex debugging

**CLI + SDD Skills (fast automation):**
- ✅ Automatic git branch + commit tracking
- ✅ Multi-agent support (Claude, Cursor, Gemini, Windsurf, etc.)
- ✅ SDD v3.0 with 5-gate validation
- ✅ Ralph loop iterative convergence
- ✅ Best for: clear specs, fast iteration, team environments

**Neither is "better"** — use the right tool for the job:
- Vague requirements + critical logic → Superpowers throughout
- Clear specs + simple features → CLI for speed
- Mixed scenarios → Use both (Superpowers for planning, CLI for execution)

---

## Next Steps (Optional)

Future enhancements in the plan:

1. **Integrate detection into CLI commands**
   - `agentic-agent run --task <ID>` warns if Superpowers unavailable
   - `agentic-agent skills ensure` suggests installing Superpowers

2. **Update CLAUDE.md template**
   - Document Superpowers integration in generated rules files
   - Show decision tree for choosing between systems

3. **Add examples to documentation**
   - Include Superpowers workflows in SMALL-PROJECTS.md
   - Real-world scenarios showing both systems in action

---

## Questions?

Refer to `superpowers-bridge/SKILL.md` for:
- Detailed troubleshooting
- Full decision trees
- 4 complete workflow examples
- Skill mapping with explanations

Or check the updated skills directly:
- `openspec/SKILL.md` — Phase 4-7 integration
- `sdd/developer/SKILL.md` — Steps 3.5, 7
- `sdd/verifier/SKILL.md` — Step 3.5
- `run-with-ralph/SKILL.md` — Step 5

---

**Commit:** 70eb704 — feat: Add Superpowers plugin integration with CLI and SDD skills
