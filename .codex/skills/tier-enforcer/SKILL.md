---
name: tier-enforcer
description: "Audit, create, or fix skills to comply with the 3-tier layered context model. Use when creating new skills, checking an existing skill for violations, or trimming a skill that has grown too large."
---

# skill:tier-enforcer

## Does exactly this

Enforces the 3-tier layered context model: Tier 1 routers stay slim, Tier 2 SKILL.md files stay focused, and Tier 3 resources hold all detail. Guides creation, auditing, and repair of skill files.

---

## When to use

- Creating a new skill from scratch
- An existing SKILL.md has grown beyond 130 lines
- A router (SKILLS.md) has grown beyond 70 lines
- A resource file has grown beyond 500 lines
- `agentic-agent validate` reports a `skill-tier` failure or warning
- Auditing the full packs directory for tier compliance

---

## Do not use this skill when

- You are simply reading an existing skill to understand it
- The task is unrelated to skill file structure or size

---

## Tier Thresholds — quick reference

| File type | WARN at | FAIL at |
|---|---|---|
| Tier 1 Router (SKILLS.md) | > 70 lines | > 100 lines |
| Tier 2 Skill (SKILL.md) | > 130 lines | > 200 lines |
| Tier 3 Resource (resources/*.md) | > 500 lines | — (warn only) |

---

## Step 1 — Identify the mode

Ask: **"Am I creating, auditing, or fixing?"**

- **Creating** → go to Step 2
- **Auditing** → run `agentic-agent validate` first, then go to Step 3
- **Fixing a violation** → go to Step 4

---

## Step 2 — Create a new skill (Tier 2)

1. Create `packs/<name>/SKILL.md` using the required sections:
   frontmatter (`name`, `description`), `# skill:<name>`, `## Does exactly this`,
   `## When to use`, numbered steps with `→` anchors to resources, `## Done when`,
   `## If you need more detail`.
   → See `resources/tier-rules.md#required-sections` for section checklist.
2. Extract all examples, tables, and long detail into `packs/<name>/resources/<topic>.md`.
3. Count lines: SKILL.md target is 60-130 lines. Resource target is 150-500 lines.
4. Add a router row to `packs/SKILLS.md` and register the pack in `internal/skills/packs.go`.

---

## Step 3 — Audit an existing skill

1. Count lines: `wc -l packs/<name>/SKILL.md` and all files in its `resources/`.
2. Check required sections are present (see `resources/tier-rules.md#required-sections`).
3. Verify all `→ resources/<file>.md` links resolve to real files.
4. If any threshold is exceeded, proceed to Step 4.

---

## Step 4 — Fix a violation

| Violation | Fix |
|---|---|
| SKILL.md > 130 lines | Extract tables or examples to resource file; replace with `→ resources/<file>.md#anchor` |
| Router > 70 lines | Remove prose; keep only trigger table and HARD RULES |
| Resource > 500 lines | Split by topic into two files; update anchors in SKILL.md |
| Missing required section | Add section; keep it to 1-3 lines |
| Broken resource link | Create the missing file or fix the path |

→ See `resources/tier-rules.md#fix-workflows` for step-by-step repair examples.

---

## Done when

- `agentic-agent validate` reports `skill-tier` as PASS with no warnings
- All SKILL.md files have all required sections
- All `→ resources/` links resolve
- Line counts are within thresholds

---

## If you need more detail

→ `resources/tier-rules.md` — Full reference tables, required section checklist, violation examples with before/after, fix workflows, router row format, and packs.go registration pattern
