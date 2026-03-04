# 3-Tier Layered Context Model — Visual Reference

Quick visual guide to understanding and maintaining the 3-tier model.

---

## The Model at a Glance

```
                    SESSION START
                        ↓
        ┌───────────────────────────────┐
        │   TIER 1: ROUTERS             │
        │   (53-59 lines)               │
        │                               │
        │   SKILLS.md                   │
        │   ├─ openspec → openspec/...  │
        │   ├─ api-docs → api-docs/...  │
        │   ├─ tdd → tdd/...            │
        │   └─ sdd:* → sdd/SKILLS.md    │
        │                               │
        │   sdd/SKILLS.md               │
        │   ├─ analyst → analyst/...    │
        │   ├─ architect → architect/.. │
        │   └─ ...                      │
        └───────────────────────────────┘
                        ↓
                   USER TRIGGER
                   "openspec"
                        ↓
        ┌───────────────────────────────┐
        │   TIER 2: SKILL FILE          │
        │   (59 lines)                  │
        │                               │
        │   openspec/SKILL.md           │
        │   ├─ Trigger                  │
        │   ├─ Does exactly this        │
        │   ├─ When to use              │
        │   ├─ Steps 1-7                │
        │   ├─ Output                   │
        │   ├─ Done when                │
        │   └─ If detail needed →       │
        │      resources/phases.md      │
        └───────────────────────────────┘
                        ↓
            USER REQUESTS DETAIL
            "Show me Phase 1 questions"
                        ↓
        ┌───────────────────────────────┐
        │   TIER 3: RESOURCE FILE       │
        │   (497 lines)                 │
        │                               │
        │   resources/phases.md         │
        │   ├─ Phase 0 (full)           │
        │   ├─ Phase 1 (full)           │
        │   ├─ Phase 2 (full)           │
        │   ├─ ... Phase 7 (full)       │
        │   └─ Examples, templates      │
        └───────────────────────────────┘
```

---

## Load Patterns

### Pattern 1: Simple Skill (No Detail Needed)

```
User: "api-docs for my REST endpoints"
  ↓
Agent loads: SKILLS.md (53 lines)
  ↓
Agent loads: api-docs/SKILL.md (58 lines)
  ↓
Agent: "Here are the 5 steps..."
Agent completes task WITHOUT loading resources/
  ↓
TOTAL LOAD: 111 lines
```

### Pattern 2: Complex Skill (Detail Needed)

```
User: "Create a Platform Spec for guest checkout"
  ↓
Agent loads: SKILLS.md (53 lines)
  ↓
Agent loads: sdd/SKILLS.md (59 lines)
  ↓
Agent loads: sdd/platform-spec/SKILL.md (128 lines)
  ↓
Agent: "I need to see the template..."
  ↓
Agent loads: resources/platform-spec-template.md (331 lines)
  ↓
Agent: "Here's the Domain Model section..."
  ↓
TOTAL LOAD: 571 lines (but only when needed)
BASELINE: 240 lines (53+59+128 without detail)
```

### Pattern 3: Multi-Step Skill (Resources Loaded Mid-Task)

```
User: "Use openspec to implement a feature"
  ↓
Agent loads: SKILLS.md + openspec/SKILL.md = 112 lines
Agent completes Step 1-3 without resources
  ↓
Step 4: "Need to create task format..."
Agent loads: resources/phases.md → Phase 4 section
  ↓
Agent completes Step 4-7 with detail as needed
  ↓
TOTAL LOAD: ~400 lines
EFFICIENCY: Loaded detail only when needed
```

---

## File Organization

### Correct Structure (3-Tier Model)

```
internal/skills/packs/
├─ SKILLS.md (53 lines) — ROUTER
├─ openspec/
│  ├─ SKILL.md (59 lines) — SLIM FILE
│  └─ resources/
│     └─ phases.md (497 lines) — DETAIL
├─ api-docs/
│  ├─ SKILL.md (58 lines) — SLIM FILE
│  └─ resources/
│     └─ api-docs-examples.md (728 lines) — DETAIL
├─ sdd/
│  ├─ SKILLS.md (59 lines) — SUB-ROUTER
│  ├─ analyst/
│  │  ├─ SKILL.md (70 lines) — SLIM FILE
│  │  └─ (no resources if not needed)
│  ├─ hotfix/
│  │  ├─ SKILL.md (118 lines) — SLIM FILE
│  │  └─ resources/
│  │     └─ hotfix-templates.md (340 lines) — DETAIL
│  ├─ platform-spec/
│  │  ├─ SKILL.md (128 lines) — SLIM FILE
│  │  └─ resources/
│  │     └─ platform-spec-template.md (331 lines) — DETAIL
│  └─ ...
```

**Key properties:**
- Routers at the top (SKILLS.md, sdd/SKILLS.md)
- Slim files in skill folders
- Detail extracted to resources/ subdirectories
- One resource file per major section/topic

### Incorrect Structure (Bloated Model)

```
internal/skills/packs/
├─ SKILLS.md (200 lines) ❌ TOO BIG
│  ├─ Contains examples
│  ├─ Contains rules
│  ├─ Contains all detail inline
├─ openspec/
│  ├─ SKILL.md (450 lines) ❌ TOO BIG
│  │  ├─ All 7 phases inline
│  │  ├─ All templates inline
│  │  ├─ All examples inline
│  │  └─ (no resources/ directory)
└─ ...
```

**Problems:**
- Routers bloated (200+ lines)
- Skill files bloated (450+ lines)
- No extraction to resources/
- Every trigger loads unnecessary detail
- 92% token waste

---

## Size Targets

### Tier 1: Routers

```
        Line Count          Ideal Size        Max Size
SKILLS.md      |===|           53 lines         70 lines
sdd/SKILLS.md  |==|            59 lines         70 lines

Never grows beyond one line per skill + formatting
```

**If a router exceeds 70 lines:**
- Find what was added inline
- Extract it to the appropriate Tier 2 or Tier 3
- Router should only contain trigger → file mappings

### Tier 2: Skill Files

```
        Line Count          Ideal Size        Max Size
openspec/      |===|           59 lines        130 lines
api-docs/      |===|           58 lines        130 lines
hotfix/        |=====|         118 lines       130 lines
architect/     |====|          65 lines        130 lines

Most files cluster around 60-80 lines
Anything over 130 needs extraction
```

**If a skill file exceeds 130 lines:**
1. Read it to identify what's making it large
2. Extract to `resources/[topic].md`
3. Link from the skill file
4. Verify skill file is back under 100 lines

### Tier 3: Resource Files

```
        Line Count          Ideal Size        Max Size
phases/md      |================|  497 lines    500 lines
hotfix-tmpl/   |===========|      340 lines    500 lines
api-examples/  |==================| 728 lines   (large OK)
constitution/  |============|      402 lines   500 lines

Can grow large (150-500 lines)
Split if exceeds 500 lines
```

**If a resource file exceeds 500 lines:**
1. Identify logically distinct sections
2. Create separate resource files (e.g., `templates.md`, `examples.md`)
3. Update skill files to link to multiple resources
4. Create a `README.md` in resources/ explaining the layout

---

## Decision Tree: Where Does This Content Belong?

```
                    New Content
                        ↓
        Does it fit in 1-2 sentences?
                    ↙           ↘
                  YES            NO
                   ↓              ↓
        ┌─────────────────┐  Need more space
        │ Keep in Tier 2  │
        │ (SKILL.md)      │
        └─────────────────┘      ↓
                          Is it related to
                          the current step?
                             ↙       ↘
                           YES        NO
                            ↓          ↓
            ┌───────────────────┐  Keep separate
            │ Inline in         │  resource file
            │ Tier 2 with link  │
            │ to Tier 3         │
            └───────────────────┘

Example flows:

"What are the 3 phases?"
  ↓ (1 sentence answer)
  → Tier 2: "Phase 1: Requirements (3 steps)"
  → Tier 3: "See resources/phases.md"

"Show me Phase 1 questions"
  ↓ (needs detail)
  → Tier 3: Full questionnaire template

"Here's an example API response"
  ↓ (is an example)
  → Tier 3: "resources/api-examples.md"

"What's the first step?"
  ↓ (1-2 lines)
  → Tier 2: "1. Read requirements file"
```

---

## Token Efficiency Comparison

### Before Refactor (Bloated Model)

```
When user says "openspec"

Total lines loaded UPFRONT:
  SKILL.md (518 lines) — everything inline
  + overhead files
  = ~700+ lines EVERY TIME

Example: 10 triggers × 700 lines = 7,000 lines overhead
Result: Expensive, noisy, imprecise
```

### After Refactor (3-Tier Model)

```
When user says "openspec"

TIER 1:  SKILLS.md (53 lines) — loaded once per session
TIER 2:  openspec/SKILL.md (59 lines) — loaded on trigger
TIER 3:  resources/phases.md (497 lines) — loaded ON DEMAND

Total BASELINE load: 112 lines
Total WITH detail: 609 lines

Example: 10 triggers × 112 lines = 1,120 lines baseline
         (Same trigger again = 0 extra lines, just reference)

Result: Efficient, focused, precise
Savings: 92% baseline, 13% with detail
```

---

## Maintenance Signals

### 🟢 GREEN — Healthy System

```
✅ Routers: 50-70 lines
✅ Skill files: 54-130 lines
✅ Resources: 150-500 lines
✅ No broken links
✅ All skills follow template
✅ Last maintenance: < 1 month ago
```

**Action:** Monthly check, no urgent fixes needed.

### 🟡 YELLOW — Needs Attention

```
⚠️  Router: 80-100 lines
⚠️  Skill file: 130-180 lines
⚠️  Resource: 450-550 lines
⚠️  1-2 broken links found
⚠️  Some skills don't follow template
⚠️  Last maintenance: 1-2 months ago
```

**Action:** Extract detail, verify links, schedule refactor.

### 🔴 RED — Critical Issues

```
🚨 Router: 100+ lines
🚨 Skill file: 180+ lines
🚨 Resource: 800+ lines
🚨 3+ broken links
🚨 Many skills off-template
🚨 Last maintenance: 3+ months ago
```

**Action:** Immediate extraction, refactoring sprint, reset practices.

---

## Quick Maintenance Commands

### Monthly Check (Copy & Paste)

```bash
#!/bin/bash
echo "=== TIER MAINTENANCE CHECK ==="
echo ""
echo "ROUTERS:"
wc -l internal/skills/packs/SKILLS.md internal/skills/packs/sdd/SKILLS.md
echo ""
echo "SKILL FILES (top 5 largest):"
find internal/skills/packs -name "SKILL.md" -exec wc -l {} + | sort -rn | head -5
echo ""
echo "RESOURCE FILES (count):"
find internal/skills/packs -path "*/resources/*.md" | wc -l
echo ""
echo "BROKEN LINKS CHECK:"
grep -r "→ resources/" internal/skills/packs --include="SKILL.md" | \
  while IFS: read file content; do
    link=$(echo "$content" | sed 's/.*resources\///' | awk '{print $1}')
    dir=$(dirname "$file")
    if [ ! -f "$dir/resources/$link" ]; then
      echo "❌ BROKEN: $file → resources/$link"
    fi
  done | head -10
```

### Find Files Over Limit

```bash
# Skill files over 130 lines
find internal/skills/packs -name "SKILL.md" -exec wc -l {} + | awk '$1 > 130 {print}'

# Resource files over 500 lines
find internal/skills/packs -path "*/resources/*.md" -exec wc -l {} + | awk '$1 > 500 {print}'

# Routers over 70 lines
wc -l internal/skills/packs/SKILLS.md internal/skills/packs/sdd/SKILLS.md | awk '$1 > 70 {print}'
```

---

## Training New Team Members

### Day 1: Understanding the Model

Show them this visualization:

```
"When someone says 'openspec'..."

Step 1: Agent reads SKILLS.md (master router)
        Agent: "I see 'openspec' → go to ./openspec/SKILL.md"

Step 2: Agent reads openspec/SKILL.md (skill instructions)
        Agent: "I see 7 phases. User wants detail on Phase 1"

Step 3: Agent reads resources/phases.md (detail file)
        Agent: "Here's the full Phase 1 questionnaire"

Agent never loads all 3 at once unless needed.
That's the entire philosophy: Load What You Need.
```

### Day 2: "I Want to Add a New Skill"

1. Create `internal/skills/packs/[name]/SKILL.md`
2. Write ~60 lines following the template
3. If detail needed, create `resources/[topic].md`
4. Add 1 line to the router (`SKILLS.md`)
5. Done — no bloat, easy to maintain

### Day 3: "The Skill is Growing"

If a skill file exceeds 130 lines:
1. Extract the excess to `resources/`
2. Link from the skill file
3. Verify skill file is under 100 lines
4. Done — system stays healthy

---

## One-Page Summary

**3-Tier Layered Context Model**

| Layer | File | Size | Content | Loaded | Grows |
|-------|------|------|---------|--------|-------|
| **1** | SKILLS.md | 50-70 lines | Trigger → file mappings | Every session | Never |
| **2** | SKILL.md | 54-130 lines | Steps, output, links | On trigger | Capped |
| **3** | resources/ | 150-500 lines | Detail, examples, templates | On demand | Unbounded |

**Golden Rule:** Load what you need, nothing else.

**Maintenance:** Monthly size check, extract when files grow.

**Result:** 92% baseline token reduction, full detail on-demand.

---

End of visual reference. Print this page for your team's reference.
