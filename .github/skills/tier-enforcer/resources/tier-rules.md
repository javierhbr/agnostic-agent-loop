# Tier Rules Reference

Complete reference for the 3-tier layered context model: routers, skill files, and resource files.

---

## ## #required-sections

### Required SKILL.md Sections

| Section | Required | Max Lines in SKILL.md | Detail Lives In |
|---|---|---|---|
| Frontmatter (name, description) | Yes | 3 | SKILL.md |
| `# skill:<name>` | Yes | 1 | SKILL.md |
| `## Does exactly this` | **Yes** | 2-3 | SKILL.md (overview only) |
| `## When to use` | Yes | 3-5 | SKILL.md (bullets only) |
| `## Do not use this skill when` | No | 2-3 | SKILL.md (optional) |
| Numbered steps or phases | Yes | 5-15 | SKILL.md + `→ resources/` links |
| `## Done when` | Yes | 3-5 | SKILL.md (bullets only) |
| `## If you need more detail` | **Yes** | 1-2 | SKILL.md (arrow link only) |

**Mandatory sections:** `Does exactly this`, `If you need more detail`

All other sections may be condensed or omitted if not applicable. Examples, tables, templates, and detailed explanations belong in Tier 3 resources, not inline in SKILL.md.

---

## ## #thresholds

### Size Threshold Reference

| File Type | Location | WARN at | FAIL at | Reasoning |
|---|---|---|---|---|
| Tier 1 Router | `SKILLS.md`, `sdd/SKILLS.md` | > 70 lines | > 100 lines | Routers map triggers to files only. Should never exceed 1 line per trigger + formatting (~60 lines). Anything over 70 signals prose creep. |
| Tier 2 Skill | `packs/<name>/SKILL.md` | > 130 lines | > 200 lines | Skill files must fit on one screen (~60 lines target). Under 130 is one page. Over 200 means the detail should be extracted. |
| Tier 3 Resource | `packs/<name>/resources/*.md` | > 500 lines | — (warn only) | Resources can be large (examples, templates, phases). Over 500 lines suggests splitting by topic. No hard FAIL; agents can handle 500-line documents. |

**Token economy impact:** A 130-line SKILL.md loaded unnecessarily is ~260 tokens wasted. A 500-line resource loaded only on demand costs 0 tokens until needed. Thresholds enforce the "load nothing you don't need" philosophy.

---

## ## #router-row-format

### How to Add a Trigger Row

**Format:**
```
| `<trigger>` | `<relative-path>` | <purpose (≤10 words)> |
```

**Rules:**
1. Trigger name: kebab-case (e.g., `tier-enforcer`, `code-review-ai`, `run-with-ralph`)
2. Relative path: starts with `./` and ends with `/SKILL.md` (e.g., `./tier-enforcer/SKILL.md`)
3. Purpose: 1-2 sentences, ≤ 10 words, plain English (no markdown, no bullets)
4. Alphabetical order within the trigger section (if multiple rows)
5. Do NOT include explanations, examples, or conditions — routers are trigger maps only

**Examples:**

✅ **Good:**
```
| `tier-enforcer` | `./tier-enforcer/SKILL.md` | Audit, create, or fix skill files for 3-tier compliance |
| `tier-check`    | `./tier-enforcer/SKILL.md` | Audit, create, or fix skill files for 3-tier compliance |
```

❌ **Bad:**
```
| `tier-enforcer` | `./tier-enforcer/SKILL.md` | Enforces the 3-tier model by checking that routers stay under 70 lines, skill files stay under 130, and resources are extracted properly. Use this when you want to audit an existing skill or create a new one. |
```

---

## ## #violation-examples

### Before/After: SKILL.md with Inline Table Bloat

**Before (432 lines) — VIOLATION:**
```markdown
# skill:api-docs

## Steps

1. Analyze the API

   | Endpoint | Method | Auth | Rate Limit |
   |---|---|---|---|
   | /users | GET | Bearer | 100/min |
   | /users | POST | Bearer | 50/min |
   ... (30 more rows of endpoint tables)

2. Generate documentation

   Table of REST patterns (80 lines)
   Table of GraphQL patterns (60 lines)
   Table of Auth patterns (40 lines)

3. Create examples

   Example 1: REST GET request with auth (25 lines)
   Example 2: GraphQL query (35 lines)
   ... (8 more examples)

## Done when

- Documentation generated
- Examples provided
- All endpoints covered
```

**Problems:**
- Tables inline → token bloat
- Examples inline → makes file unreadable
- Every trigger loads 432 lines → 92% waste

**After (58 lines) — FIXED:**
```markdown
# skill:api-docs

## Does exactly this

Generate comprehensive API documentation from code with examples and best practices.

## Steps

1. Analyze API structure
2. Create endpoint reference
3. Add integration examples
   → See `resources/api-examples.md` for REST, GraphQL, and Auth patterns
4. Document edge cases and error handling

## Done when

- All endpoints documented
- Examples provided for each pattern
- Errors and edge cases explained

## If you need more detail

→ `resources/api-examples.md` — Endpoint tables, REST/GraphQL/Auth patterns with 12 examples, best practices
```

**Benefits:**
- SKILL.md: 58 lines (under 130 ✅)
- Resource: 374 lines with all tables and examples ✅
- Baseline load: 58 lines only (92% reduction)
- Detail load: 58 + 374 = 432 lines (same total, but loaded on-demand)

---

### Before/After: Router with Embedded Prose

**Before (178 lines) — VIOLATION:**
```markdown
# SKILLS.md

## Overview

This router maps skill triggers to their implementation files. The system has been designed
to enforce the 3-tier layered context model where routers stay slim, skill files focus on
essentials, and resources hold all detail.

## How to use this router

When you see a trigger like "openspec" or "tier-enforcer", look it up in the table below
and load the corresponding SKILL.md file. Never load multiple files at once unless...

## Important rules

1. Routers must stay under 70 lines
2. If a skill file grows...
3. Always use resources/ for...

## Trigger table

| trigger | file | purpose |
...
```

**Problems:**
- Contains explanations (paragraphs 1-3)
- Contains rules (should be in CLAUDE.md or skill instructions)
- Loaded every session → prolixity bloat

**After (53 lines) — FIXED:**
```markdown
# SKILLS.md

This router maps triggers to skills for spec-driven development and related practices.

| Trigger | Skill File | Purpose |
|---|---|---|
| `openspec` | `./openspec/SKILL.md` | Spec-driven development from requirements |
| `api-docs` | `./api-docs/SKILL.md` | Generate API documentation |
| `tier-enforcer` | `./tier-enforcer/SKILL.md` | Audit, create, fix skill files |
| `sdd:*` | `./sdd/SKILLS.md` | SDD v3.0 role-specific skills (sub-router) |

## Hard rules

1. Always use `agentic-agent` CLI for task operations
2. Never write directly to `.agentic/` YAML files
3. Always claim a task before working on it
```

**Benefits:**
- Router: 53 lines (under 70 ✅)
- No prose, only mappings and critical rules
- Clear, scannable, maintainable

---

### Before/After: Resource over 500 Lines Split into Two

**Before (823 lines) — VIOLATION:**
```markdown
# Templates and Examples (resources/hotfix-all.md)

## Bug Spec Template
... (120 lines)

## Hotfix Spec Template
... (130 lines)

## Follow-up Spec Template
... (80 lines)

## Example: Production Outage
... (180 lines)

## Example: Data Corruption
... (90 lines)

## Example: Security Breach
... (123 lines)
```

**Problems:**
- Single file at 823 lines
- Agents must search through all examples to find template
- Unclear organization

**After (two files) — FIXED:**

`resources/hotfix-templates.md` (340 lines):
```markdown
# Hotfix Templates

## Bug Spec Template (Non-Urgent)
... (120 lines)

## Hotfix Spec Template (Production, Urgent)
... (130 lines)

## Follow-up Spec Template
... (80 lines)
```

`resources/hotfix-examples.md` (393 lines):
```markdown
# Hotfix Examples

## Example: Production Outage
... (180 lines)

## Example: Data Corruption
... (90 lines)

## Example: Security Breach
... (123 lines)
```

**From skill file, link to both:**
```markdown
Step 2: Create a spec file (bug, hotfix, or follow-up)
→ See `resources/hotfix-templates.md` for templates

Step 3: Review similar incidents for context
→ See `resources/hotfix-examples.md` for real-world examples
```

**Benefits:**
- Each resource file under 500 lines ✅
- Clear separation: templates vs examples
- Agents can load just the template if needed

---

## ## #fix-workflows

### Workflow A: Shrink a SKILL.md File (Over 130 Lines)

**Symptoms:**
- `wc -l packs/<skill>/SKILL.md` returns > 130
- File contains tables, code examples, or multi-paragraph explanations
- Reader feedback: "overwhelming" or "too much detail"

**Steps:**

1. **Read the entire file** and identify what's making it large:
   - Tables (endpoint lists, comparison matrices, requirement tables)
   - Code examples (REST requests, GraphQL queries, configuration blocks)
   - Detailed explanations (more than 1-2 sentences per step)
   - Templates (full file structures, questionnaires)

2. **Create resources/ subdirectory** if it doesn't exist:
   ```bash
   mkdir -p packs/<skill>/resources
   ```

3. **Extract each identified section** to a resource file:
   - **Tables** → `resources/<topic>-reference.md` (e.g., `endpoints-reference.md`)
   - **Examples** → `resources/<topic>-examples.md` (e.g., `rest-examples.md`)
   - **Templates** → `resources/<topic>-template.md` (e.g., `spec-template.md`)
   - **Detailed guides** → `resources/<topic>-guide.md` (e.g., `phase-1-guide.md`)

4. **Add anchors** to each resource file for direct linking:
   ```markdown
   ## #endpoints-reference

   [content]

   ## #common-patterns

   [content]
   ```

5. **Replace inline content** in SKILL.md with links:
   ```markdown
   # Old (inline)
   1. Choose an endpoint
      | GET /users | List all users |
      | POST /users | Create user |
      (30 more rows)

   # New (with link)
   1. Choose an endpoint
      → See `resources/endpoints-reference.md#rest-endpoints` for complete list
   ```

6. **Add "If you need more detail" section** if missing:
   ```markdown
   ## If you need more detail

   → `resources/<topic>.md` — Full reference tables, examples, templates
   ```

7. **Count lines** and verify SKILL.md is now < 130:
   ```bash
   wc -l packs/<skill>/SKILL.md
   ```

8. **Test** by reading the skill file and verifying you can understand the workflow without the resource file.

---

### Workflow B: Shrink a Router (Over 70 Lines)

**Symptoms:**
- `wc -l packs/SKILLS.md` or `packs/sdd/SKILLS.md` returns > 70
- File contains explanations, examples, or detailed rules
- Hard to scan the trigger list

**Steps:**

1. **Read the entire router file** and identify excess content:
   - Introductory paragraphs explaining the system
   - Multi-paragraph usage instructions
   - Detailed rules that should be in CLAUDE.md
   - Examples of how to use triggers

2. **Keep ONLY these elements:**
   - Brief description (1 sentence max)
   - Trigger-to-file mapping table
   - 2-3 critical "hard rules" (non-negotiable constraints, e.g., "Never write to .agentic/ directly")

3. **Move excess content** to appropriate locations:
   - **Usage instructions** → Move to project README or wiki
   - **Detailed rules** → Move to CLAUDE.md (project instructions)
   - **Examples** → Move to individual SKILL.md files
   - **System explanations** → Move to `docs/` folder

4. **Keep the trigger table clean:**
   ```markdown
   | Trigger | File | Purpose |
   |---|---|---|
   | `openspec` | `./openspec/SKILL.md` | Spec-driven development from requirements |
   | `api-docs` | `./api-docs/SKILL.md` | Generate API documentation |
   ```

5. **Verify the hard rules are truly critical:**
   - "Always use CLI for task operations" ✅ (critical, non-negotiable)
   - "Routers should be under 70 lines" ❌ (explanatory, belongs in LAYERED-CONTEXT-MODEL.md)
   - "Never write directly to .agentic/" ✅ (critical, prevents data corruption)

6. **Count lines** and verify router is < 70:
   ```bash
   wc -l packs/SKILLS.md
   ```

---

### Workflow C: Split a Resource File (Over 500 Lines)

**Symptoms:**
- `wc -l packs/<skill>/resources/<file>.md` returns > 500
- File contains logically distinct topics (e.g., templates + examples, API + CLI)
- Agents search through irrelevant sections to find what they need

**Steps:**

1. **Read the entire resource file** and identify distinct topics:
   - Examples: `## Example: Scenario A`, `## Example: Scenario B`, etc.
   - Templates: `## Template: Thing 1`, `## Template: Thing 2`, etc.
   - Sections: `## Phase 1`, `## Phase 2`, etc.

2. **Group related sections** into logical files:
   - If file has "Template" sections and "Example" sections → split into `<topic>-templates.md` and `<topic>-examples.md`
   - If file has Phase 0-3 and Phase 4-7 → split into `phases-0-3.md` and `phases-4-7.md`
   - If file has REST, GraphQL, gRPC → split into `rest-examples.md`, `graphql-examples.md`, `grpc-examples.md`

3. **Create new files** with clear anchors:
   ```bash
   # Original
   packs/<skill>/resources/<topic>.md (823 lines)

   # Split into
   packs/<skill>/resources/<topic>-templates.md (340 lines)
   packs/<skill>/resources/<topic>-examples.md (380 lines)
   packs/<skill>/resources/README.md (50 lines, explains which file has what)
   ```

4. **Add anchors** to each file for linking:
   ```markdown
   # File 1: templates.md

   ## #bug-spec-template
   ...

   ## #hotfix-spec-template
   ...

   # File 2: examples.md

   ## #production-outage-example
   ...
   ```

5. **Update SKILL.md** to link to the new files:
   ```markdown
   # Old
   Step 2: Create a spec
   → See resources/hotfix-templates.md for all templates

   # New
   Step 2: Create a spec (bug, hotfix, or follow-up)
   → See `resources/hotfix-templates.md#bug-spec-template` for template

   Step 3: Review examples for context
   → See `resources/hotfix-examples.md#production-outage-example` for real incident
   ```

6. **Create README.md** in resources/ directory to explain organization:
   ```markdown
   # Hotfix Resources

   - `hotfix-templates.md` — Spec templates (Bug, Hotfix, Follow-up)
   - `hotfix-examples.md` — Real-world incident examples
   - `README.md` — This file
   ```

7. **Verify** each file is under 500 lines:
   ```bash
   wc -l packs/<skill>/resources/*.md | tail -1
   ```

---

## ## #packs-registration

### Registering a SkillPack in packs.go

**Location:** `internal/skills/packs.go`, inside `NewPackRegistry()` function

**Pattern:**
```go
r.Register(SkillPack{
    Name:        "skill-name",
    Description: "One-sentence description of what this skill does",
    Files: []SkillPackFile{
        {SrcPath: "packs/skill-name/SKILL.md", DstPath: "skill-name/SKILL.md"},
        {SrcPath: "packs/skill-name/resources/topic.md", DstPath: "skill-name/resources/topic.md"},
        // Add more files as needed (resources, guides, templates)
    },
})
```

**Rules:**

1. **Name**: kebab-case, exactly matches directory name in `packs/`
2. **Description**: One sentence, clearly states what the skill does
3. **SrcPath**: Relative to `internal/skills/` (note the `//go:embed packs/*` at line 9)
   - Format: `packs/<skill-name>/<relative-path>`
   - Examples:
     - `packs/tier-enforcer/SKILL.md`
     - `packs/tier-enforcer/resources/tier-rules.md`
4. **DstPath**: Relative to the tool skill directory (e.g., `.claude/skills/`)
   - Format: `<skill-name>/<relative-path>`
   - Examples:
     - `tier-enforcer/SKILL.md`
     - `tier-enforcer/resources/tier-rules.md`
5. **Files array**: Always include SKILL.md; add resource files as needed
6. **IsAgent field**: Set to `true` for agent markdown files only (rare); omit for skills

**Example with resources:**
```go
r.Register(SkillPack{
    Name:        "openspec",
    Description: "Spec-driven development from requirements files through task execution and verification",
    Files: []SkillPackFile{
        {SrcPath: "packs/openspec/SKILL.md", DstPath: "openspec/SKILL.md"},
        {SrcPath: "packs/openspec/resources/phases.md", DstPath: "openspec/resources/phases.md"},
    },
})
```

**Example with agent file:**
```go
r.Register(SkillPack{
    Name:        "openclaw",
    Description: "OpenClaw autonomous agent factory pattern",
    Files: []SkillPackFile{
        {SrcPath: "packs/openclaw/SKILL.md", DstPath: "openclaw/SKILL.md"},
        {SrcPath: "packs/openclaw/AGENT.md", DstPath: "openclaw-orchestrator.md", IsAgent: true},
        {SrcPath: "packs/openclaw/agents/worker.md", DstPath: "openclaw-worker.md", IsAgent: true},
    },
})
```

---

## Summary

**3-Tier Compliance Checklist:**

- [ ] **Tier 1 (Routers):** SKILLS.md < 70 lines, only triggers + hard rules
- [ ] **Tier 2 (Skills):** SKILL.md < 130 lines, has all required sections
- [ ] **Tier 3 (Resources):** Each file < 500 lines, organized by topic with anchors
- [ ] **Links:** All `→ resources/` paths resolve to real files
- [ ] **Sections:** All SKILL.md files have "Does exactly this" and "If you need more detail"
- [ ] **Registration:** SkillPack registered in packs.go with correct SrcPath/DstPath
- [ ] **Validation:** `agentic-agent validate` reports no `skill-tier` FAILs
