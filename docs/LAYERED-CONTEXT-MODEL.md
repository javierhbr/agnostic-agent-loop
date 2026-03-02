# The 3-Tier Layered Context Model

## Overview

The 3-tier layered context model is an architecture for organizing agent guidance that **minimizes context bloat while maximizing precision**. The core principle is simple:

> **An agent that loads everything knows nothing useful.**

This document explains the model, shows how it's implemented in this project, and provides maintenance guidelines for the long term.

---

## The Core Problem It Solves

When you give an agent a large file, it doesn't get smarter — it gets noisier. It starts averaging across everything it loaded instead of focusing on what the task actually needs. The output becomes:

- Generic (covers too many cases)
- Inconsistent (picks arbitrary path through the noise)
- Slow (has to reason across excess context)
- Imprecise (misses edge cases and requirements)

**Solution:** A layered architecture where agents only load what they currently need.

---

## The 3-Tier Architecture

```
┌─────────────────────────────────────────┐
│ TIER 1: ROUTERS                         │
│ Always loaded at session start          │
│ - SKILLS.md (53 lines)                  │
│ - sdd/SKILLS.md (59 lines)              │
│ - AGENTS.md (if applicable)             │
│ Purpose: ROUTING ONLY                   │
│ Growth: NEVER (capped at ~60 lines)     │
└─────────────────────────────────────────┘
         ↓ Loads one file based on trigger
┌─────────────────────────────────────────┐
│ TIER 2: SLIM SKILL FILES                │
│ Loaded when trigger matches             │
│ - openspec/SKILL.md (59 lines)          │
│ - api-docs/SKILL.md (58 lines)          │
│ - sdd/component-spec/SKILL.md (54 lines)│
│ Purpose: ESSENTIAL ONLY                 │
│ Structure: trigger + steps + output     │
│ Growth: Capped at 100-130 lines         │
└─────────────────────────────────────────┘
         ↓ Links to resources on demand
┌─────────────────────────────────────────┐
│ TIER 3: RESOURCE FILES                  │
│ Loaded only when detail is needed       │
│ - resources/phases.md (497 lines)       │
│ - resources/hotfix-templates.md (340)   │
│ - resources/platform-spec-template.md   │
│ Purpose: DETAIL, EXAMPLES, TEMPLATES    │
│ Growth: Unbounded (as needed)           │
└─────────────────────────────────────────┘
```

### How It Works

```
User triggers: "openspec"
     ↓
SKILLS.md (53 lines) loaded
  Sees: "openspec" trigger
  Routes to: ./openspec/SKILL.md
     ↓
openspec/SKILL.md loaded (59 lines)
  Shows: 7 phase names
  Links to: "See resources/phases.md for detail"
     ↓
Agent asks: "Need Phase 1 questions?"
User: "Yes"
     ↓
resources/phases.md loaded (497 lines)
  Shows: Full questionnaire, templates, examples
     ↓
Agent completes task with:
  TOTAL CONTEXT = 53 + 59 + 497 = 609 lines
  (Instead of loading 609 lines upfront)
```

---

## Why This Works: Token Math

### Old approach (Pre-refactor)
```
User triggers: "openspec"
     ↓
Agent loads entire openspec/SKILL.md (518 lines)
+ Agent loads entire api-docs/SKILL.md (484 lines)
+ Agent loads some related playbooks (~300 lines)
     ↓
TOTAL UPFRONT LOAD = 1,302 lines
Average task uses: 200-300 lines
Waste: 80% of context is noise
```

### New approach (3-Tier Model)
```
User triggers: "openspec"
     ↓
Agent loads SKILLS.md (53 lines)
Agent loads openspec/SKILL.md (59 lines)
Agent loads resources/phases.md ON DEMAND (497 lines)
     ↓
TOTAL AT START = 112 lines
TOTAL WHEN NEEDED = 609 lines
Baseline reduction: 78% (112 vs 1,302)
On-demand reduction: 53% (609 vs 1,302)
```

---

## This Project's Implementation

### Tier 1: Routers

**`/internal/skills/packs/SKILLS.md` (53 lines)**
- Maps 13 flat-pack skills
- Maps `sdd:*` to sub-router
- Contains hard rules (one trigger, one skill, no chaining)

**`/internal/skills/packs/sdd/SKILLS.md` (59 lines)**
- Maps 15 SDD v3.0 role-specific skills
- Contains workflow routing (Low/Medium/High risk → Quick/Standard/Full)
- Contains hard rules for gates and ADRs

**Why these are tiny:**
- They only say "given trigger X, load file Y"
- They contain zero implementation detail
- They never grow (if skills are added, we add one router line, never refactor)

### Tier 2: Slim Skill Files

**All skill files follow this template:**

```markdown
---
name: [skill-name]
description: [one sentence]
---

# skill:[name]

## Does exactly this
[One sentence — the output.]

## Use this skill when
- [Trigger condition 1]
- [Trigger condition 2]

## Steps — in order, no skipping
1. [Step in 1-2 lines]
2. [Step in 1-2 lines]
3. [Step in 1-2 lines]

## Output
[Exact artifact produced]

## Done when
- [Completion condition]
- [Completion condition]

## If you need more detail
→ `resources/<file>.md`
```

**Size constraint: 54-137 lines** (Target: 60-70 lines)

Examples from this project:
- `openspec/SKILL.md` — 59 lines
- `api-docs/SKILL.md` — 58 lines
- `sdd/component-spec/SKILL.md` — 54 lines
- `sdd/hotfix/SKILL.md` — 118 lines

**Why each file is exactly this size:**
- Contains only what the agent needs to START the task
- All steps are action-oriented (not explanations)
- All detail is behind the "If you need more detail" arrow

### Tier 3: Resource Files

**Resource files are where detail lives.** They're loaded on demand and organized by section.

**Examples from this project:**

| Skill | Resource File | Lines | Content |
|-------|---------------|-------|---------|
| openspec | `resources/phases.md` | 497 | 7 full Phase workflows with questionnaires, task formats, templates |
| api-docs | `resources/api-docs-examples.md` | 728 | REST, GraphQL, Auth examples; best practices; pitfall tables |
| hotfix | `resources/hotfix-templates.md` | 340 | Bug Spec, Hotfix Spec, Follow-up Spec templates; incident examples |
| platform-spec | `resources/platform-spec-template.md` | 331 | Full file structure, domain model, UX flows, NFRs, contracts |
| platform-constitution | `resources/constitution-template.md` | 402 | 6-section governance template, gate enforcement, amendment process |

**Resource files are NEVER loaded unless:**
- The skill explicitly points to them ("If you need more detail →")
- The user explicitly asks for more information
- A step requires them to work

---

## Key Principles for Long-Term Maintenance

### 1. **Routers Never Grow**

Once you create `SKILLS.md` or `sdd/SKILLS.md`, they should never exceed 60 lines.

**When you add a new skill:**
- Add ONE line to the router (the mapping)
- Put the skill in its own file (Tier 2)
- Put the detail in resources/ (Tier 3)
- Do NOT refactor the router

**Bad example:**
```markdown
# SKILLS.md (now 150 lines)
[Contains detailed instructions for 10 different workflows]
[Embeds examples]
[Has inline templates]
```

**Good example:**
```markdown
# SKILLS.md (59 lines)
[Maps triggers to files, one line each]
[Points to resources for detail]
[No growth, no refactoring]
```

### 2. **Skill Files Have a Hard Size Limit**

Each skill file must fit on one page (60-80 lines, max 130).

**If your skill file is growing:**
1. Identify what's making it grow (examples? templates? explanations?)
2. Extract to `resources/`
3. Link from the skill file
4. Verify the skill file is back under 80 lines

**Violating this rule:**
- Makes the skill file harder to parse
- Loads unnecessary context for every task
- Breaks the "load nothing you don't need" principle

### 3. **Resources Live Outside the Skill File**

Never embed in the skill file:
- Examples (unless 1 sentence)
- Templates (if more than a few lines)
- Detailed explanations (if more than a sentence per step)
- Tables (if more than 3 rows)

**Instead:**
- Create `resources/` subdirectory in the skill folder
- Put examples there
- Link from the skill file with `→ resources/<file>.md`

**Good structure:**
```
/skills/openspec/
  SKILL.md (59 lines — triggers, steps, links)
  resources/
    phases.md (497 lines — full Phase 0-7 workflows)
```

**Bad structure:**
```
/skills/openspec/
  SKILL.md (600 lines — everything inline)
  (no resources/)
```

### 4. **Resource Files Are Organized by Section**

When you create a resource file, use clear anchors so agents can link directly to sections.

**Good example (hotfix-templates.md):**
```markdown
# Hotfix Templates

...

## Bug Spec Template (Non-Urgent)

File name: `BUG-<number>-...`

...

## Hotfix Spec Template (Production, Urgent)

File name: `HOTFIX-<number>-...`

...

## Follow-up Spec Template

File name: `FOLLOWUP-<n>-...`

...
```

**From the skill file, you can link:**
```markdown
→ See resources/hotfix-templates.md#bug-spec-template
→ See resources/hotfix-templates.md#hotfix-spec-template
```

### 5. **Link From Tier 2 to Tier 3, Never Embed**

**Bad (embeds resource):**
```markdown
# skill:write-spec

## Steps
1. Create a spec file

## Input required

[500 lines of spec template here]
```

**Good (links to resource):**
```markdown
# skill:write-spec

## Steps
1. Create a spec file
   → See resources/spec-templates.md for full template

## If you need more detail
→ resources/spec-templates.md
```

### 6. **Maintain the Tier 2 → Tier 3 Mapping**

Keep a mental map of what resources each skill points to.

**Example from this project:**

| Skill | Resources |
|-------|-----------|
| `openspec` | `resources/phases.md` |
| `api-docs` | `resources/api-docs-examples.md` |
| `sdd/hotfix` | `resources/hotfix-templates.md` |
| `sdd/platform-spec` | `resources/platform-spec-template.md` |
| `sdd/platform-constitution` | `resources/constitution-template.md` |

When you change a resource file, verify all skills that link to it still work.

---

## Long-Term Maintenance Checklist

### Monthly Review

- [ ] **Router size check**: `wc -l SKILLS.md sdd/SKILLS.md` — both should be < 70 lines
- [ ] **Skill file sizes**: `wc -l **/SKILL.md` — all should be 54-130 lines
- [ ] **Resource file count**: `find . -path "*/resources/*.md" | wc -l` — should be 15-20 for the project
- [ ] **Broken links**: Scan each skill file for `→ resources/` links; verify those files exist

### Quarterly Refactor Audit

- [ ] **Skill drift**: Are any skill files growing toward 200 lines? If so, extract and split
- [ ] **Resource bloat**: Are resource files growing past 500 lines? Consider splitting by section
- [ ] **Duplication**: Are the same examples/templates appearing in multiple resources? Consolidate
- [ ] **Staleness**: Are any resources out of date? Tag with `<!-- Last updated: 2026-03-01 -->` and update

### When Adding a New Skill

1. Create the skill folder: `/skills/[name]/`
2. Create `/skills/[name]/SKILL.md` (use template, target 60-70 lines)
3. If detail needed, create `/skills/[name]/resources/` subdirectory
4. Create resource files with sections and anchors
5. Update the router:
   - Add one line to `SKILLS.md` (alphabetical order)
   - Verify router is still < 70 lines
6. Test: Can an agent load just the skill file? Can they find resources when needed?
7. Document in project MEMORY.md

### When a Skill File Exceeds 130 Lines

**Symptoms:**
- File is getting hard to read
- Contains multiple examples
- Has explanations that are more than 1 sentence per step
- Readers say "this is overwhelming"

**How to fix:**

1. Read the entire skill file
2. Identify what's making it large:
   - Examples? → Extract to `resources/examples.md`
   - Templates? → Extract to `resources/templates.md`
   - Explanations? → Condense to 1-2 lines, link to resource
   - Edge cases? → Move to troubleshooting section in resource
3. Create the resource file with sections and anchors
4. Update the skill file to link to the resource
5. Verify skill file is back under 80 lines
6. Test with an agent

---

## Real Example: How the Model Prevents Bloat

### Scenario: Adding a new feature

**Old approach (no layered model):**
```
Product Manager: "Add new feature X to the CLI"
Developer: Adds new skill file
[Skill file grows to 300 lines with examples, templates, rules]
[Every time someone triggers this skill, they load 300 lines]
[After 10 skills added this way, SKILLS.md is 3,000 lines]
[Agents are drowning in context]
```

**New approach (with layered model):**
```
Product Manager: "Add new feature X to the CLI"
Developer:
  1. Creates /skills/feature-x/SKILL.md (65 lines)
     - Trigger, steps, output
     - Links to resources
  2. Creates /skills/feature-x/resources/templates.md (250 lines)
     - All templates, examples, edge cases
  3. Adds 1 line to SKILLS.md router
     feature-x → ./feature-x/SKILL.md
  4. SKILLS.md is still 60 lines
  5. Agent loads:
     - At trigger: 60 + 65 = 125 lines
     - When detail needed: 125 + 250 = 375 lines
     - Baseline reduction: 48% (125 vs 3,000 old bloated approach)
```

---

## Maintenance Dashboard

Create a file at `docs/TIER-MAINTENANCE.md` (updated quarterly) to track the model's health:

```markdown
# Tier Maintenance Dashboard

**Last updated:** [date]
**Maintained by:** [your name]

## Tier 1: Routers
- SKILLS.md: 53 lines ✅
- sdd/SKILLS.md: 59 lines ✅
- Hard rules: 4 rules ✅

## Tier 2: Skill Files
- Total skills: 19
- Average size: 85 lines
- Largest: sdd/adr/SKILL.md (79 lines) ✅
- Size range: 54-118 lines ✅

## Tier 3: Resource Files
- Total resources: 16
- Average size: 374 lines
- Largest: constitution-template.md (402 lines) ✅
- Smallest: component-spec-guide.md (232 lines) ✅

## Drift Detection
- Skill files over 130 lines: 0 ✅
- Routers over 70 lines: 0 ✅
- Resource files over 500 lines: 0 ✅
- Orphaned resources (linked but missing): 0 ✅
- Broken links (link but files missing): 0 ✅

## Recent changes
- Phase 3 refactor complete (2025-02-28)
- All 19 skills refactored, 16 resources created
- Next audit: 2026-05-28
```

---

## Troubleshooting

### Problem: "An agent loaded too much context and got confused"

**Root cause:** Tier 2 or 3 is loading more than needed

**Fix:**
1. Check which file the agent loaded
2. See if it's too large (> 130 lines for Tier 2, > 500 lines for Tier 3)
3. Extract the unused sections to a separate resource file
4. Add anchors so the agent can link directly to what they need

### Problem: "Routers are growing and becoming hard to maintain"

**Root cause:** Adding too much detail to SKILLS.md or sdd/SKILLS.md

**Fix:**
1. Audit the router — it should only contain trigger → file mappings
2. Any explanations or rules should move to the skill files themselves
3. If router exceeds 100 lines, split into sub-routers (like we did with sdd/SKILLS.md)

### Problem: "Resource files are becoming unmanageable (> 1000 lines)"

**Root cause:** Resource file is trying to be a catch-all for multiple skills

**Fix:**
1. Identify logically distinct sections in the resource file
2. Split into separate resource files (e.g., `hotfix-templates.md` vs `hotfix-incident-examples.md`)
3. Update all skill files that link to the old file
4. Add a `README.md` in the resources/ folder that explains which file has what

### Problem: "New team members don't know how the system is organized"

**Root cause:** The mental model isn't documented or communicated

**Fix:**
1. Create `docs/LAYERED-CONTEXT-MODEL.md` (this document)
2. Add it to onboarding checklist
3. Show new members:
   - How to find what they need (routers first, then skills, then resources)
   - How to add a new skill (template + resource)
   - Why the system is organized this way (token efficiency + precision)

---

## Summary

The 3-tier layered context model is a maintenance discipline, not a one-time refactor.

| Tier | Size | Growth | Loaded | Maintenance |
|------|------|--------|--------|-------------|
| **Routers** | ~60 lines | Never | Every session | Monthly size check |
| **Skill files** | 54-130 lines | Capped | On trigger | Extract if > 130 lines |
| **Resources** | 150-500 lines | Unbounded | On demand | Organize by section, audit quarterly |

**The golden rule:** When in doubt, extract detail to a resource file and link from the skill. It's always better to have too many resource files than a bloated skill file.

Long-term success = keeping Tier 1 and Tier 2 lean, while Tier 3 grows only as needed.
