# Documentation Index & Map

**Navigate the SDD/OpenSpec documentation by your needs.**

---

## By Document Type

### 📋 **README.md** (Your Starting Point)
**390 lines | Navigation hub**

What's inside:
- Decision tree by project type (Tiny → Small → Monorepo → Complex)
- CLI command cheat sheet
- When to use what workflow
- Real examples
- FAQ

**Start here if:** You're new or need overall orientation

**Go here next:** Pick your scenario below

---

### ⚡ **CHEATSHEET.md** (Quick Lookup)
**500 lines | One-page reference (print-friendly)**

What's inside:
- Workflow comparison table
- All CLI commands on one page
- Config templates (minimal + monorepo)
- 7 implementation prompts (quick version)
- Decision trees
- Real examples (bug fix, small feature, monorepo)
- Common gotchas

**Start here if:** You need quick answers or want to print

**Bookmark this:** It's your quick reference while coding

---

### 🎨 **DIAGRAMS.md** (Visual Explanation)
**800 lines | 12 ASCII diagrams**

What's inside:
1. Workflow selection decision tree
2. OpenSpec lifecycle
3. Monorepo multi-path spec discovery
4. Task vs OpenSpec vs SDD comparison
5. Single package workflow
6. Monorepo workflow (parallel teams)
7. Spec to code flow
8. Gate check flow (SDD)
9. Implementation prompt flow
10. Multi-service dependency graph
11. Time estimates by workflow
12. Config comparison (single vs monorepo)

**Start here if:** You're a visual learner or need to understand flow

**Refer to:** When explaining to team members

---

### 📖 **SMALL-PROJECTS.md** (Single Package Guide)
**586 lines | For solo devs and small teams**

What's inside:
- 3 workflow tiers (Tiny/Small/Medium)
- When to use which approach
- Spec templates (copy-paste)
- Task templates (copy-paste)
- Real example (product search, payment flow)
- Superpowers integration (when to use TDD, gates, etc.)
- FAQ for small projects
- Decision matrix

**Start here if:** You're building a single feature or package

**Use this for:** Templates, examples, Superpowers guidance

---

### 🏢 **MONOREPO-OPENSPEC.md** (Multi-Package Guide)
**669 lines | For monorepo teams with parallel development**

What's inside:
- Root config setup (multi-path spec resolution)
- Per-package spec organization
- How specs are discovered
- **7 Implementation Prompts** (detailed):
  1. Spec Review (decompose to tasks)
  2. Package-Specific (guide single package)
  3. Integration (cross-package)
  4. Testing (test structure)
  5. Verification (completeness check)
  6. Dependency Analysis (multi-service order)
  7. Documentation (generate from specs)
- Real examples:
  - OAuth integration (3 packages)
  - Payment flow (3 services)
- Multi-service coordination patterns
- CLI commands reference
- FAQ for monorepos

**Start here if:** You're in a monorepo with multiple teams

**Use this for:** 7 prompts, multi-service examples, dependency planning

---

## By Your Scenario

### 🐛 "I found a bug" → **Fix it in 1 hour**
```
Quick fix workflow:
agentic-agent task create "Fix: XYZ"
agentic-agent task claim <id>
[code]
agentic-agent task complete <id>
```

**Read:** [CHEATSHEET.md - Example 1: Quick Bug Fix](./CHEATSHEET.md#example-1-quick-bug-fix)

**Or:** [SMALL-PROJECTS.md - For Tiny Projects](./SMALL-PROJECTS.md#for-tiny-projects-solo-developer-low-risk)

---

### ✨ "I'm building a new feature" → **1-2 week feature, single package**
```
Feature workflow:
agentic-agent openspec init "Feature" --from requirements.md
agentic-agent task claim <id>
[code]
agentic-agent task complete <id>
agentic-agent openspec complete <id>
```

**Read:** [SMALL-PROJECTS.md - For Small Projects](./SMALL-PROJECTS.md#for-small-projects-1-2-people-1-service)

**Templates:** [SMALL-PROJECTS.md - Template: Small Project Spec](./SMALL-PROJECTS.md#template-small-project-spec-copy--paste)

**Real example:** [CHEATSHEET.md - Example 2: Small Feature](./CHEATSHEET.md#example-2-small-feature)

---

### 🏢 "I'm in a monorepo" → **Multiple packages, parallel teams**
```
Monorepo workflow:
# Each package has specs/ directory
# Root config with multi-path specDirs
agentic-agent openspec init "Feature" --from packages/service-a/specs/...
agentic-agent task claim <id>
[code in packages/service-a]
agentic-agent task complete <id>
```

**Read:** [MONOREPO-OPENSPEC.md](./MONOREPO-OPENSPEC.md)

**Setup:** [MONOREPO-OPENSPEC.md - Setup](./MONOREPO-OPENSPEC.md#setup-one-time)

**Real examples:**
- [OAuth (3 packages)](./MONOREPO-OPENSPEC.md#multi-service-example-payment-flow)
- [Payment Flow (3 services)](./MONOREPO-OPENSPEC.md#workflow-from-spec-to-production)

**7 Prompts:** [MONOREPO-OPENSPEC.md - Prompts for Implementation](./MONOREPO-OPENSPEC.md#prompts-for-implementation)

**Visual:** [DIAGRAMS.md #6: Monorepo Workflow](./DIAGRAMS.md) + [#10: Multi-Service Dependencies](./DIAGRAMS.md)

---

### 💳 "This is critical/high-risk" → **Payments, auth, regulatory**
```
SDD workflow:
agentic-agent sdd start "Feature" --risk critical
agentic-agent sdd workflow show <id>
agentic-agent sdd gate-check <id>
[implementation with gates]
```

**Read:** README.md or SDD-specific docs (linked from README)

**Gates explained:** [DIAGRAMS.md #8: Gate Check Flow](./DIAGRAMS.md)

---

### 🎯 "I need quick answers" → **Don't read everything**
```
Fastest path:
1. Scan CHEATSHEET.md (5 min)
2. Pick your scenario
3. Jump to specific section
```

**Start:** [CHEATSHEET.md](./CHEATSHEET.md)

---

### 👥 "I need to explain this to my team"
```
Best resources:
1. Show DIAGRAMS.md #1 (decision tree)
2. Point to CHEATSHEET.md (all commands)
3. Share relevant example
```

**For visuals:** [DIAGRAMS.md](./DIAGRAMS.md)

**For reference:** [CHEATSHEET.md](./CHEATSHEET.md)

**For examples:** [SMALL-PROJECTS.md](./SMALL-PROJECTS.md) or [MONOREPO-OPENSPEC.md](./MONOREPO-OPENSPEC.md)

---

## By Question Type

| Question | Answer Location |
|----------|------------------|
| "What workflow should I use?" | [README.md - Workflow Comparison](./README.md#workflow-comparison) or [DIAGRAMS.md #4](./DIAGRAMS.md) |
| "What are all the CLI commands?" | [CHEATSHEET.md - Quick Commands](./CHEATSHEET.md#quick-commands) or [README.md - CLI Command Cheat Sheet](./README.md#cli-command-cheat-sheet) |
| "How do I set up my config?" | [CHEATSHEET.md - Config Templates](./CHEATSHEET.md#config-templates) or [MONOREPO-OPENSPEC.md - Setup](./MONOREPO-OPENSPEC.md#setup-one-time) |
| "How do I structure specs?" | [SMALL-PROJECTS.md - Template](./SMALL-PROJECTS.md#template-small-project-spec-copy--paste) or [MONOREPO-OPENSPEC.md - Creating Specs](./MONOREPO-OPENSPEC.md#creating-specs-per-package) |
| "What are the 7 prompts?" | [MONOREPO-OPENSPEC.md - Prompts for Implementation](./MONOREPO-OPENSPEC.md#prompts-for-implementation) or [CHEATSHEET.md - 7 Implementation Prompts](./CHEATSHEET.md#7-implementation-prompts) |
| "How do I do multi-service?" | [MONOREPO-OPENSPEC.md](./MONOREPO-OPENSPEC.md) or [DIAGRAMS.md #10](./DIAGRAMS.md) |
| "What takes how long?" | [DIAGRAMS.md #11: Time Estimates](./DIAGRAMS.md) |
| "How do gates work?" | [DIAGRAMS.md #8: Gate Check Flow](./DIAGRAMS.md) |
| "I'm stuck on my task" | Use [7 Prompts](./MONOREPO-OPENSPEC.md#prompts-for-implementation) or [DIAGRAMS.md #9: Prompt Flow](./DIAGRAMS.md) |
| "Can I print something?" | [CHEATSHEET.md](./CHEATSHEET.md) (designed for print) |
| "How is this different from X?" | [CHEATSHEET.md - Workflow Comparison](./CHEATSHEET.md#workflow-comparison) |

---

## Document Statistics

| Document | Lines | Key Sections | Best For |
|----------|-------|--------------|----------|
| README.md | 390 | Navigation, overview, FAQ | Starting point |
| CHEATSHEET.md | 500 | Commands, templates, examples | Quick lookup, printing |
| DIAGRAMS.md | 800 | 12 ASCII diagrams, flows | Visual explanation |
| SMALL-PROJECTS.md | 586 | Tiny/Small/Medium workflows, templates | Single package features |
| MONOREPO-OPENSPEC.md | 669 | Setup, 7 prompts, multi-service examples | Monorepo teams |
| **Total** | **2,945** | **50+ sections** | **Complete reference** |

---

## Quick Links

**Getting started?**
- [README.md - Quick Start](./README.md#quick-start-by-project-type)
- [CHEATSHEET.md](./CHEATSHEET.md)

**Building a feature?**
- [SMALL-PROJECTS.md](./SMALL-PROJECTS.md)
- [MONOREPO-OPENSPEC.md](./MONOREPO-OPENSPEC.md)

**Visual learner?**
- [DIAGRAMS.md](./DIAGRAMS.md)

**Need examples?**
- [CHEATSHEET.md - Real Examples](./CHEATSHEET.md#real-examples)
- [SMALL-PROJECTS.md - Real Example](./SMALL-PROJECTS.md#real-example-small-team-adding-search)
- [MONOREPO-OPENSPEC.md - Real Examples](./MONOREPO-OPENSPEC.md#workflow-from-spec-to-production)

**Need templates?**
- [SMALL-PROJECTS.md - Templates](./SMALL-PROJECTS.md#template-small-project-spec-copy--paste)
- [CHEATSHEET.md - Spec Template](./CHEATSHEET.md#spec-template-copy-paste)

**Need the 7 prompts?**
- [MONOREPO-OPENSPEC.md - Prompts for Implementation](./MONOREPO-OPENSPEC.md#prompts-for-implementation)
- [CHEATSHEET.md - 7 Implementation Prompts](./CHEATSHEET.md#7-implementation-prompts)

---

## Reading Paths by Role

### Product Manager
**Understand the methodology:**
1. [README.md - Workflow Comparison](./README.md#workflow-comparison)
2. [DIAGRAMS.md #1 - Decision Tree](./DIAGRAMS.md)
3. [DIAGRAMS.md #11 - Time Estimates](./DIAGRAMS.md)

### Developer (Single Package)
**Build features efficiently:**
1. [CHEATSHEET.md](./CHEATSHEET.md) (bookmark this)
2. [SMALL-PROJECTS.md - For Small Projects](./SMALL-PROJECTS.md#for-small-projects-1-2-people-1-service)
3. Use prompts as needed from [7 Implementation Prompts](./CHEATSHEET.md#7-implementation-prompts)

### Developer (Monorepo)
**Coordinate across teams:**
1. [MONOREPO-OPENSPEC.md - Setup](./MONOREPO-OPENSPEC.md#setup-one-time)
2. [MONOREPO-OPENSPEC.md - 7 Implementation Prompts](./MONOREPO-OPENSPEC.md#prompts-for-implementation)
3. [DIAGRAMS.md #10 - Dependencies](./DIAGRAMS.md)
4. [CHEATSHEET.md](./CHEATSHEET.md) (for quick commands)

### Tech Lead
**Plan and oversee work:**
1. [README.md](./README.md)
2. [DIAGRAMS.md](./DIAGRAMS.md) (all 12 diagrams for planning)
3. [MONOREPO-OPENSPEC.md](./MONOREPO-OPENSPEC.md) (if multi-team)
4. Share relevant excerpts with team

### Architect
**Design and review specs:**
1. [MONOREPO-OPENSPEC.md - Creating Specs](./MONOREPO-OPENSPEC.md#creating-specs-per-package)
2. [MONOREPO-OPENSPEC.md - Prompts #6 (Dependencies)](./MONOREPO-OPENSPEC.md#6-dependency-analysis-prompt-multi-service)
3. [DIAGRAMS.md #10 - Multi-Service Dependencies](./DIAGRAMS.md)

---

## How to Use This Index

**1. Find your role/scenario above**
2. **Read the documents in order** (they're linked)
3. **Bookmark CHEATSHEET.md** for quick lookups
4. **Share diagrams** with your team
5. **Use prompts** when stuck on implementation

**Questions?** Check the [FAQ sections](./README.md#faq-small-projects) in each document.

---

## File Structure

```
docs/sdd-mcp/
├── INDEX.md                      ← You are here
├── README.md                     ← Start here
├── CHEATSHEET.md                 ← Print this
├── DIAGRAMS.md                   ← Show your team
├── SMALL-PROJECTS.md             ← Single package guide
├── MONOREPO-OPENSPEC.md          ← Monorepo guide
├── SMALL-PROJECTS.md.bak         ← (previous version, can delete)
└── operation model/              ← SDD operating model docs
    ├── OPERATING-MODEL.md
    ├── one-pager-01.md
    ├── six-pager-01.md
    └── target-operating-model.md
```

---

## Next Steps

1. **Bookmark:** [CHEATSHEET.md](./CHEATSHEET.md)
2. **Read:** [README.md](./README.md)
3. **Choose your scenario** and jump to relevant guide
4. **Use the 7 prompts** when implementing
5. **Share diagrams** with your team

**That's it. You're ready to go.**
