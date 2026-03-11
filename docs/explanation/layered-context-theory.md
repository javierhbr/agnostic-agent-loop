# 3-Tier Layered Context Model — Complete Documentation

This directory contains comprehensive documentation for understanding, implementing, and maintaining the 3-tier layered context model used throughout this project.

---

## What Is This Model?

The 3-tier layered context model is an architecture for organizing agent guidance that **minimizes context bloat while maximizing precision**. Instead of loading entire skill files at once, the system loads content in three layers:

1. **Tier 1: Routers** (53-59 lines) — Master routing files loaded at session start
2. **Tier 2: Slim Skill Files** (54-130 lines) — Instructions loaded when triggered
3. **Tier 3: Resource Files** (150-500 lines) — Detail loaded on demand

**Result:** 92% baseline token reduction while keeping all detail immediately accessible.

---

## Quick Start (5 Minutes)

### Understanding the Model

**Start here:** [LAYERED-CONTEXT-VISUAL.md](LAYERED-CONTEXT-VISUAL.md)

This document contains:
- ASCII diagrams showing the 3-tier architecture
- Visual examples of how content flows through the system
- Size targets and health signals (Green/Yellow/Red)
- Quick reference commands for monitoring
- One-page summary

**Time to read:** 5 minutes

### Deep Understanding (20 Minutes)

**Then read:** [LAYERED-CONTEXT-MODEL.md](LAYERED-CONTEXT-MODEL.md)

This document contains:
- Full explanation of the core problem and how the model solves it
- Token math showing efficiency gains (before/after)
- How this project implements the model
- Key principles for long-term maintenance
- Real-world example of how the model prevents bloat
- Troubleshooting guide for common issues
- Maintenance dashboard template

**Time to read:** 20 minutes

### Practical Maintenance (Referenced as Needed)

**For operational tasks:** [TIER-MAINTENANCE-GUIDE.md](TIER-MAINTENANCE-GUIDE.md)

This document contains:
- Monthly 15-minute maintenance routine with exact bash commands
- Step-by-step procedures for extracting detail from bloated files
- Procedures for splitting large resource files
- Procedures for adding new skills (template included)
- Quarterly 60-minute deep dive audit
- Git workflow for tracking maintenance
- Troubleshooting problems that arise
- Maintenance checklist template

**Time to reference:** 5-10 minutes per operation

---

## Which Document Do I Read?

### "I want to understand the architecture"
→ **[LAYERED-CONTEXT-MODEL.md](LAYERED-CONTEXT-MODEL.md)**
- Explains the philosophy and principles
- Shows how token efficiency works
- Provides troubleshooting for conceptual questions

### "I want a quick visual reference"
→ **[LAYERED-CONTEXT-VISUAL.md](LAYERED-CONTEXT-VISUAL.md)**
- ASCII diagrams and visual flows
- Size targets and health signals
- Decision trees for where content belongs
- Print this for your desk

### "I need to do maintenance"
→ **[TIER-MAINTENANCE-GUIDE.md](TIER-MAINTENANCE-GUIDE.md)**
- Step-by-step procedures with bash commands
- Extract detail from a skill file (20-30 min)
- Add a new skill (45-60 min)
- Quarterly audit (60 min)
- Troubleshooting specific problems

### "I'm implementing this in a new project"
→ **[LAYERED-CONTEXT-MODEL.md](LAYERED-CONTEXT-MODEL.md#long-term-maintenance-checklist) → [TIER-MAINTENANCE-GUIDE.md](TIER-MAINTENANCE-GUIDE.md#adding-a-new-skill)**
- Start with the model document for philosophy
- Use maintenance guide for the setup procedures
- Reference the visual guide for quick checks

---

## The Model in This Project

### Current State

**Tier 1: Routers**
- `internal/skills/packs/SKILLS.md` (53 lines)
- `internal/skills/packs/sdd/SKILLS.md` (59 lines)

**Tier 2: Skill Files**
- 19 total skills refactored
- Size range: 54-118 lines
- Average: 85 lines
- All follow the standard template

**Tier 3: Resource Files**
- 16 resource files created
- Total: 5,986 lines
- Organized by skill with clear sections and anchors

**Status:** ✅ Production-ready with established maintenance procedures

### File Location Map

```
docs/
├─ README-LAYERED-CONTEXT.md (this file)
├─ LAYERED-CONTEXT-MODEL.md (theory + principles)
├─ TIER-MAINTENANCE-GUIDE.md (operational procedures)
├─ LAYERED-CONTEXT-VISUAL.md (quick reference diagrams)
│
internal/skills/packs/
├─ SKILLS.md (Tier 1 router)
├─ sdd/
│  ├─ SKILLS.md (Tier 1 sub-router)
│  ├─ [15 skills]/
│  │  ├─ SKILL.md (Tier 2)
│  │  └─ resources/ (Tier 3)
│  │
│  ├─ hotfix/
│  │  ├─ SKILL.md (118 lines)
│  │  └─ resources/
│  │     └─ hotfix-templates.md (340 lines)
│  │
│  ├─ platform-spec/
│  │  ├─ SKILL.md (128 lines)
│  │  └─ resources/
│  │     └─ platform-spec-template.md (331 lines)
│  │
│  ├─ platform-constitution/
│  │  ├─ SKILL.md (122 lines)
│  │  └─ resources/
│  │     └─ constitution-template.md (402 lines)
│  │
│  └─ [12 other skills]/
│     ├─ SKILL.md (50-85 lines)
│     └─ resources/ (if needed)
│
├─ openspec/
│  ├─ SKILL.md (59 lines)
│  └─ resources/
│     └─ phases.md (497 lines)
│
├─ api-docs/
│  ├─ SKILL.md (58 lines)
│  └─ resources/
│     └─ api-docs-examples.md (728 lines)
│
└─ [9 other flat-pack skills]/
   ├─ SKILL.md (54-85 lines)
   └─ resources/ (if needed)
```

---

## Key Metrics to Monitor

### Monthly (15 minutes)

Run the commands in [TIER-MAINTENANCE-GUIDE.md#monthly-maintenance-routine](TIER-MAINTENANCE-GUIDE.md#monthly-maintenance-routine):

```bash
# Router sizes
wc -l internal/skills/packs/SKILLS.md internal/skills/packs/sdd/SKILLS.md
# Should be: 53 and 59 lines

# Skill file sizes
find internal/skills/packs -name "SKILL.md" -exec wc -l {} + | sort -n
# All should be < 130 lines

# Count resources
find internal/skills/packs -path "*/resources/*.md" | wc -l
# Should be around 15-20

# Check for broken links
grep -r "→ resources/" internal/skills/packs --include="SKILL.md"
# All links should point to existing files
```

### Quarterly (60 minutes)

Run the commands in [TIER-MAINTENANCE-GUIDE.md#quarterly-deep-dive](TIER-MAINTENANCE-GUIDE.md#quarterly-deep-dive):

- Size audit (find files exceeding limits)
- Content review (identify extraction opportunities)
- Link verification (check all references work)
- Duplication check (consolidate redundant content)
- Staleness check (update outdated resources)

---

## Common Workflows

### "I'm adding a new skill"

1. Read: [TIER-MAINTENANCE-GUIDE.md#adding-a-new-skill](TIER-MAINTENANCE-GUIDE.md#adding-a-new-skill)
2. Follow the 6-step procedure (45-60 minutes)
3. Results:
   - New skill folder created
   - SKILL.md written (60-70 lines)
   - Router updated (1 line added)
   - Resources created if needed
   - Tested and verified

### "A skill file is getting too big"

1. Read: [TIER-MAINTENANCE-GUIDE.md#extracting-detail-from-a-skill-file](TIER-MAINTENANCE-GUIDE.md#extracting-detail-from-a-skill-file)
2. Follow the 6-step procedure (20-30 minutes)
3. Results:
   - Skill file trimmed to < 100 lines
   - Detail extracted to resource file
   - Links added
   - Tested and verified

### "A resource file is over 500 lines"

1. Read: [TIER-MAINTENANCE-GUIDE.md#splitting-a-large-resource-file](TIER-MAINTENANCE-GUIDE.md#splitting-a-large-resource-file)
2. Follow the 6-step procedure (30-45 minutes)
3. Results:
   - Resource file split into logical pieces
   - Multiple smaller files created
   - Skill files updated with multiple links
   - Old file deleted

### "Something isn't working with a skill"

1. Read: [TIER-MAINTENANCE-GUIDE.md#responding-to-this-skill-isnt-working](TIER-MAINTENANCE-GUIDE.md#responding-to-this-skill-isnt-working)
2. Diagnose the issue type (bloat, clarity, missing content, etc.)
3. Apply the appropriate fix
4. Verify with an agent
5. Document lessons learned

---

## Health Indicators

### 🟢 GREEN — System is healthy

✅ All routers < 70 lines
✅ All skill files 54-130 lines
✅ Resources organized with clear sections
✅ Zero broken links
✅ Monthly maintenance completed
✅ No files over size limits

**Action:** Continue monthly checks, nothing urgent.

### 🟡 YELLOW — Attention needed

⚠️ Router 80-100 lines
⚠️ Skill file 130-180 lines
⚠️ Resource file 450-550 lines
⚠️ 1-2 broken links
⚠️ Monthly maintenance overdue

**Action:** Schedule extraction/refactoring in next 2 weeks.

### 🔴 RED — Critical issues

🚨 Router > 100 lines
🚨 Skill file > 180 lines
🚨 Resource file > 800 lines
🚨 Multiple broken links
🚨 Monthly maintenance overdue > 1 month

**Action:** Immediate intervention required. Follow [TIER-MAINTENANCE-GUIDE.md](TIER-MAINTENANCE-GUIDE.md).

---

## The Golden Rule

> **An agent that loads everything knows nothing useful.**

Keep this principle in mind when:
- Adding new content
- Maintaining existing files
- Training team members
- Designing new skills

The moment a file starts to grow beyond its tier's size limit, extract detail and link from it. It's always the right move.

---

## For Different Roles

### Product Manager
- Understand the model: [LAYERED-CONTEXT-MODEL.md#the-core-problem-it-solves](LAYERED-CONTEXT-MODEL.md#the-core-problem-it-solves)
- Why it matters: Reduced token waste = faster, cheaper AI operations
- Reference for skill requests: [TIER-MAINTENANCE-GUIDE.md#adding-a-new-skill](TIER-MAINTENANCE-GUIDE.md#adding-a-new-skill)

### Developer Adding a Skill
- Quick start: [LAYERED-CONTEXT-VISUAL.md](LAYERED-CONTEXT-VISUAL.md)
- Detailed procedure: [TIER-MAINTENANCE-GUIDE.md#adding-a-new-skill](TIER-MAINTENANCE-GUIDE.md#adding-a-new-skill)
- Template: [TIER-MAINTENANCE-GUIDE.md#step-2-write-the-skill-file-use-template](TIER-MAINTENANCE-GUIDE.md#step-2-write-the-skill-file-use-template)

### Architect Designing New Features
- Model principles: [LAYERED-CONTEXT-MODEL.md](LAYERED-CONTEXT-MODEL.md)
- Size targets: [LAYERED-CONTEXT-VISUAL.md#size-targets](LAYERED-CONTEXT-VISUAL.md#size-targets)
- Decision tree: [LAYERED-CONTEXT-VISUAL.md#decision-tree-where-does-this-content-belong](LAYERED-CONTEXT-VISUAL.md#decision-tree-where-does-this-content-belong)

### Team Lead / Maintainer
- Everything in this file
- Monthly routine: [TIER-MAINTENANCE-GUIDE.md#monthly-maintenance-routine-15-minutes](TIER-MAINTENANCE-GUIDE.md#monthly-maintenance-routine-15-minutes)
- Quarterly audit: [TIER-MAINTENANCE-GUIDE.md#quarterly-deep-dive-60-minutes](TIER-MAINTENANCE-GUIDE.md#quarterly-deep-dive-60-minutes)
- Troubleshooting: [TIER-MAINTENANCE-GUIDE.md#troubleshooting-problems](TIER-MAINTENANCE-GUIDE.md#troubleshooting-problems)

### New Team Member (Onboarding)
1. Day 1: Read [LAYERED-CONTEXT-VISUAL.md](LAYERED-CONTEXT-VISUAL.md) (5 min)
2. Day 2: Read [LAYERED-CONTEXT-MODEL.md#the-3-tier-architecture](LAYERED-CONTEXT-MODEL.md#the-3-tier-architecture) (10 min)
3. Day 3: Read [TIER-MAINTENANCE-GUIDE.md#adding-a-new-skill](TIER-MAINTENANCE-GUIDE.md#adding-a-new-skill) (20 min)
4. Week 1: Participate in monthly maintenance routine

---

## Frequently Asked Questions

### Q: Why do we need all three tiers?

**A:** Each tier serves a purpose:
- Tier 1 (routers) never grow and provide consistent entry points
- Tier 2 (skill files) are lean and focused on action
- Tier 3 (resources) can grow unbounded with examples and templates

Separating them prevents Tier 1 and 2 from bloating while allowing Tier 3 to be comprehensive.

### Q: What if I just want to keep everything in one file?

**A:** You get:
- Tier 1 bloat: Routers grow to 200-500 lines
- Tier 2 bloat: Skill files grow to 1,000+ lines
- Context explosion: Every trigger loads everything
- 80% token waste: Most loaded content unused
- Confusion: Agents can't find what they need
- Maintenance nightmare: Everything needs updating

### Q: When should I split a resource file?

**A:** When it exceeds 500 lines OR when it covers 3+ logically distinct topics.

See [TIER-MAINTENANCE-GUIDE.md#splitting-a-large-resource-file](TIER-MAINTENANCE-GUIDE.md#splitting-a-large-resource-file) for the procedure.

### Q: Can I have nested resource files?

**A:** No. Keep the structure simple:
```
skill/
├─ SKILL.md
└─ resources/
   ├─ topic1.md
   ├─ topic2.md
   └─ topic3.md
```

Not:
```
skill/
└─ resources/
   └─ guides/
      └─ advanced/
         └─ deep/
            └─ file.md  ← Too nested
```

### Q: How often should I do maintenance?

**A:** Minimum: Monthly (15 minutes) for size checks. Recommended: Monthly + Quarterly (60 minutes) for deep audit.

See [TIER-MAINTENANCE-GUIDE.md](TIER-MAINTENANCE-GUIDE.md) for schedules.

### Q: What if my team doesn't follow the model?

**A:** That's the most common problem. Solution:
1. Make it part of code review (size checks in PR review)
2. Automate the checks (put bash script in CI)
3. Educate regularly (link docs in every PR)
4. Lead by example (apply model to all new work)

### Q: Can I apply this model to other projects?

**A:** Yes! The model is framework-agnostic. See [TIER-MAINTENANCE-GUIDE.md#to-tailor-this-to-your-repo](TIER-MAINTENANCE-GUIDE.md#to-tailor-this-to-your-repo) for adaptation guidance.

---

## Maintenance Calendar

### Monthly (Start of Month)

- Run size audit (15 min)
- Check for broken links (5 min)
- Update maintenance report
- **Commit:** `docs: Monthly tier maintenance audit — [month]`

### Quarterly (Start of Quarter: Jan, Apr, Jul, Oct)

- Deep dive audit (60 min)
- Content review
- Duplication check
- Staleness review
- **Commit:** `docs: Q[N] [year] tier maintenance audit — comprehensive review`

### When Adding Skills

- Follow [TIER-MAINTENANCE-GUIDE.md#adding-a-new-skill](TIER-MAINTENANCE-GUIDE.md#adding-a-new-skill)
- Update MEMORY.md with new skill info
- **Commit:** `feat: Add skill:[name]`

### When Extracting/Splitting

- Follow the appropriate procedure
- Update maintenance report
- **Commit:** `refactor: [Extract/Split] [skill-name] resources`

---

## Summary

The 3-tier layered context model is a discipline, not a one-time setup. Success comes from:

1. **Understanding the philosophy:** Load only what you need
2. **Following the structure:** Tier 1 routers → Tier 2 skills → Tier 3 resources
3. **Monitoring metrics:** Monthly size checks, quarterly deep dives
4. **Maintaining discipline:** Extract when files grow, link when referencing detail
5. **Educating the team:** Make it part of your culture

This directory contains everything needed to implement and maintain the model long-term.

---

## Document Checklist

- [x] **LAYERED-CONTEXT-MODEL.md** — Complete explanation with principles and troubleshooting
- [x] **TIER-MAINTENANCE-GUIDE.md** — Step-by-step operational procedures
- [x] **LAYERED-CONTEXT-VISUAL.md** — Visual diagrams and quick reference
- [x] **README-LAYERED-CONTEXT.md** — This navigation document

**Estimated reading time:**
- Quick overview: 5-10 minutes
- Full understanding: 40-50 minutes
- Implementation/maintenance: Referenced as needed

**Bookmark this file and the visual guide for quick reference.**

---

Start here: [LAYERED-CONTEXT-VISUAL.md](LAYERED-CONTEXT-VISUAL.md)
