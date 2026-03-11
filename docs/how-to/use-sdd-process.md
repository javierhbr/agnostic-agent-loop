# Using the SDD Process Guide Skill

The **SDD Process Guide** skill is your step-by-step companion for executing the complete SDD v3.0 methodology.

## Quick Start

When you're ready to use SDD for a new feature:

```bash
# 1. Ensure SDD skills are installed
agentic-agent skills ensure

# 2. Open the process guide
cat .claude/skills/sdd/process-guide/SKILL.md
# OR use Claude Code to reference it while working
```

## What the Process Guide Covers

The skill walks through **4 complete phases**:

### Phase 0: Initiative Definition (Product Manager)
- Define clear problem statement with metrics
- Identify affected components
- Run risk assessment interview (5 questions)
- Create initiative with success criteria

**CLI:** `agentic-agent specifyify start "Name" --risk [level]`

### Phase 1: Architecture Design (Solution Architect)
- Read initiative and Platform Constitution
- Create feature-spec.md and component-spec.md files
- Run all 5 gate checks
- Create parallel development tasks

**CLI:** `agentic-agent specifyify gate-check SPEC-[ID]`

### Phase 2: Implementation (Developers, Parallel)
- Claim component task
- Implement code with logging/metrics/tracing
- Create implementation spec
- Gate check implementation

**CLI:** `agentic-agent task claim [ID]` → `agentic-agent task complete [ID]`

### Phase 3: Verification (Verifier)
- Verify every acceptance criterion
- Run full test suite
- Create verification report
- Sync spec graph
- Merge to main

**CLI:** `agentic-agent validate` → `agentic-agent specifyify sync-graph`

### Phase 4: Deployment & Success Metrics
- Deploy with feature flag OFF (safe)
- Progressive rollout (10% → 25% → 50% → 100%)
- Monitor metrics at each step
- Measure final success at day 30

**CLI:** `agentic-agent deploy` → `agentic-agent flags set`

---

## How to Use It in Your Workflow

### For Product Managers

```bash
# 1. Read Phase 0 section
#    Location: process-guide/SKILL.md → Phase 0

# 2. Follow steps 0.1 through 0.5
#    - Define problem statement
#    - Run risk assessment (5 questions in skill)
#    - Create initiative

agentic-agent specifyify start "Your Feature" --risk medium

# 3. Refer to Phase 4 when feature is ready to deploy
```

### For Solution Architects

```bash
# 1. Read Phase 1 section
#    Location: process-guide/SKILL.md → Phase 1

# 2. Follow steps 1.1 through 1.6
#    - Read initiative
#    - Review Platform Constitution
#    - Create feature-spec.md and component-spec.md
#    - Gate check design
#    - Create dev tasks

agentic-agent specifyify gate-check SPEC-[ID]
```

### For Developers

```bash
# 1. Read Phase 2 section
#    Location: process-guide/SKILL.md → Phase 2

# 2. Follow steps 2.1 through 2.6
#    - Claim task
#    - Read component-spec
#    - Implement code
#    - Create impl-spec
#    - Gate check implementation

agentic-agent task claim [TASK-ID]
agentic-agent specifyify gate-check SPEC-[SERVICE]-IMPL
agentic-agent task complete [TASK-ID]
```

### For Verifiers

```bash
# 1. Read Phase 3 section
#    Location: process-guide/SKILL.md → Phase 3

# 2. Follow steps 3.1 through 3.7
#    - Collect impl-specs from all teams
#    - Verify every AC
#    - Verify observability
#    - Run tests
#    - Create verify report
#    - Update spec graph
#    - Merge to main

agentic-agent validate
agentic-agent specifyify gate-check SPEC-[ID]
agentic-agent specifyify sync-graph
```

---

## Key Features of the Process Guide

✅ **Phase-by-phase structure** — Clear ownership and handoffs
✅ **Step numbers** — Know exactly where you are (0.1, 0.2, etc.)
✅ **CLI commands** — Every step shows the exact command to run
✅ **Troubleshooting** — Section at end for common issues
✅ **Cheat sheet** — Quick command reference by phase
✅ **Success indicators** — Know when you're doing it right

---

## Example: Following the Process for Guest Checkout

```bash
# Phase 0: PM defines initiative
agentic-agent specifyify start "Enable Guest Checkout" --risk medium

# Phase 1: Architect designs
# (Follow process-guide steps 1.1-1.6)
agentic-agent specifyify gate-check SPEC-CHECKOUT-GUEST

# Phase 2: 4 Developers implement in parallel
# Developer 1: Checkout Service
agentic-agent task claim TASK-001
# (Follow process-guide steps 2.1-2.6)
agentic-agent task complete TASK-001

# Phase 3: Verifier proves it works
# (Follow process-guide steps 3.1-3.7)
agentic-agent validate
agentic-agent specifyify sync-graph

# Phase 4: Deploy safely
agentic-agent deploy --environment production --feature-flags-all-off
agentic-agent flags set GuestCheckoutEnabled=10pct
agentic-agent flags set GuestCheckoutEnabled=100pct
```

---

## Location of the Skill

Once you install SDD skills, the process guide is available at:

**Project-level:**
```
.claude/skills/sdd/process-guide/SKILL.md
```

**Or embedded in CLI:**
```
internal/skills/packs/sdd/process-guide/SKILL.md
```

**Or in your IDE:**
Open Claude Code and search for `sdd-process-guide`

---

## Tips for Success

1. **Read phase introduction first** — Get the big picture before diving into steps
2. **Follow step numbers exactly** — 0.1, 0.2, 0.3, etc. They're sequential
3. **Use the cheat sheet** — When you need the command without reading the whole section
4. **Check the troubleshooting** — If something isn't working
5. **Keep the example workflow nearby** — Reference docs/sdd-example-workflow.md for a real walkthrough

---

## What Comes Next?

Once you finish Phase 4 (deployment), you can:

- ✅ Mark initiative as complete
- ✅ Document lessons learned
- ✅ Start the next feature with Phase 0 again

```bash
# See all completed initiatives
agentic-agent specifyify workflow show [initiative-id]

# Start next feature
agentic-agent specifyify start "Next Feature Name" --risk [level]
```

---

## Questions?

- **"Where's my initiative?"** → `.agentic/sdd/initiatives/[name].yaml`
- **"Which gate failed?"** → `agentic-agent specifyify gate-check SPEC-ID --format json | jq .`
- **"What should I implement?"** → Read the component-spec.md from your task
- **"How do I know it works?"** → Process guide Phase 3 (verification steps)

The process guide has you covered at every step! 🚀
