# SDD v3.0 Documentation Index

Welcome to the complete SDD (Spec-Driven Development) v3.0 documentation for your project.

---

## 📍 Start Here

**First time with SDD?**
1. Read: [**USING-SDD-PROCESS-GUIDE.md**](USING-SDD-PROCESS-GUIDE.md) — Quick orientation
2. Then: Install skills — `agentic-agent skills ensure`
3. Then: Open the process guide — `cat .claude/skills/sdd/process-guide/SKILL.md`
4. Finally: Start your first feature — `agentic-agent specifyify start "Name" --risk low`

---

## 📚 Documentation Files

### Overview & Orientation
- **[SDD-NEW-SKILL-SUMMARY.md](SDD-NEW-SKILL-SUMMARY.md)**
  What's new: the process-guide skill, how it works, what's included

- **[USING-SDD-PROCESS-GUIDE.md](USING-SDD-PROCESS-GUIDE.md)**
  How to reference the process guide while working, quick tips by role

- **[SDD-SKILLS-DIRECTORY.md](SDD-SKILLS-DIRECTORY.md)**
  Complete map of all 15 SDD skills, when to use each, decision matrix

### Practical Walkthroughs
- **[sdd-example-workflow.md](sdd-example-workflow.md)**
  Real end-to-end example: guest checkout feature through all 4 phases

### Technical Reference
- [OPERATING-MODEL.md](sdd-mcp/operation%20model/OPERATING-MODEL.md) — Full SDD v3.0 specification

---

## 🎯 By Role

### Product Manager
1. Read: [SDD-SKILLS-DIRECTORY.md](SDD-SKILLS-DIRECTORY.md) → "If You're a Product Manager"
2. Then: Process guide Phase 0 (Initiative Definition)
3. Use: initiative-definition, risk-assessment, stakeholder-communication skills

### Solution Architect
1. Read: [sdd-example-workflow.md](sdd-example-workflow.md) → Phase 1 section
2. Then: Process guide Phase 1 (Architecture Design)
3. Use: workflow-router, architect, gate-check skills

### Developer
1. Read: [sdd-example-workflow.md](sdd-example-workflow.md) → Phase 2 section
2. Then: Process guide Phase 2 (Implementation)
3. Use: developer, gate-check skills

### Verifier / QA
1. Read: [sdd-example-workflow.md](sdd-example-workflow.md) → Phase 3 section
2. Then: Process guide Phase 3 (Verification)
3. Use: verifier, gate-check skills

### DevOps / Platform Engineer
1. Read: Process guide Phase 4 (Deployment & Success)
2. Implement: Feature flag deployment, monitoring, rollout

---

## 🏗️ The Process at a Glance

```
Phase 0: Initiative Definition (PM)
  ↓ (gates check)
Phase 1: Architecture Design (Architect)
  ↓ (gates check)
Phase 2: Implementation (Developers, parallel)
  ↓ (gates check)
Phase 3: Verification (Verifier)
  ↓ (gates check)
Phase 4: Deployment & Metrics (DevOps/PM)
  ↓
SUCCESS ✓ Feature shipped with confidence
```

Each phase has 5-9 numbered steps. Each step shows:
- What to do
- Why you're doing it
- Exact CLI command
- Expected output

---

## 🛠️ Quick Command Reference

```bash
# Phase 0 (PM)
agentic-agent specifyify start "Feature Name" --risk [low|medium|high|critical]

# Phase 1 (Architect)
agentic-agent specifyify gate-check SPEC-[ID]

# Phase 2 (Developers)
agentic-agent task claim [TASK-ID]
agentic-agent task complete [TASK-ID]

# Phase 3 (Verifier)
agentic-agent validate
agentic-agent specifyify sync-graph

# Phase 4 (DevOps)
agentic-agent deploy --feature-flags-all-off
agentic-agent flags set FeatureName=100pct
```

---

## 📖 Reading Paths by Use Case

### "I need to ship a feature"
1. Read: [USING-SDD-PROCESS-GUIDE.md](USING-SDD-PROCESS-GUIDE.md)
2. Install: `agentic-agent skills ensure`
3. Follow: Process guide Phase 0 → 1 → 2 → 3 → 4

### "I'm new to SDD"
1. Read: [SDD-NEW-SKILL-SUMMARY.md](SDD-NEW-SKILL-SUMMARY.md) — What's new
2. Read: [sdd-example-workflow.md](sdd-example-workflow.md) — Real walkthrough
3. Then: Use [SDD-SKILLS-DIRECTORY.md](SDD-SKILLS-DIRECTORY.md) as reference

### "I'm confused about what to do next"
→ Open: `.claude/skills/sdd/process-guide/SKILL.md`
→ Find: Your current phase number (0, 1, 2, 3, or 4)
→ Follow: The numbered steps for that phase

### "I need to explain SDD to the team"
1. Share: [SDD-SKILLS-DIRECTORY.md](SDD-SKILLS-DIRECTORY.md) (complete overview)
2. Show: [sdd-example-workflow.md](sdd-example-workflow.md) (real example)
3. Reference: [USING-SDD-PROCESS-GUIDE.md](USING-SDD-PROCESS-GUIDE.md) (practical guide)

### "I want to understand the theory"
→ Read: [OPERATING-MODEL.md](sdd-mcp/operation%20model/OPERATING-MODEL.md) (full spec)

---

## 🎓 The 15 SDD Skills

```
process-guide ⭐ (START HERE — guides you through all 4 phases)
├── Phase 0 skills:
│   ├── initiative-definition
│   ├── risk-assessment
│   └── stakeholder-communication
├── Phase 1 skills:
│   ├── workflow-router
│   ├── architect
│   └── gate-check
├── Phase 2 skills:
│   ├── developer
│   └── gate-check
├── Phase 3 skills:
│   ├── verifier
│   └── gate-check
└── Platform skills:
    ├── platform-constitution
    ├── platform-spec
    ├── component-spec
    ├── adr
    └── hotfix
```

All installed via: `agentic-agent skills ensure`

---

## ✅ Success Checklist

You know SDD is working right when:

- ✅ Every initiative starts with Phase 0 (definition)
- ✅ Every phase ends with gate checks (`agentic-agent specifyify gate-check`)
- ✅ No phase skipped (gates enforce this)
- ✅ Acceptance criteria in GWT format (Given/When/Then)
- ✅ Observability (logging + metrics + tracing) working in staging
- ✅ Feature flags for safe production deploy (default OFF)
- ✅ Progressive rollout (10% → 25% → 50% → 100%)
- ✅ Final metrics measured at day 30
- ✅ Spec graph updated for audit trail
- ✅ Repeatable process for every feature

---

## 🚀 Getting Started Now

```bash
# 1. Ensure everything is installed
agentic-agent skills ensure

# 2. Verify process-guide is available
ls .claude/skills/sdd/process-guide/

# 3. Start your first feature
agentic-agent specifyify start "Your Feature Name" --risk low

# 4. Open the process guide
cat .claude/skills/sdd/process-guide/SKILL.md

# 5. Follow Phase 0 steps (0.1 through 0.5)
# 6. Then move to Phase 1, 2, 3, 4 as each completes

# 7. Share docs with your team
# You can distribute these docs from docs/ folder
```

---

## 📞 Help & Troubleshooting

### "Where do I find [skill]?"
See: [SDD-SKILLS-DIRECTORY.md](SDD-SKILLS-DIRECTORY.md)

### "A gate failed—what do I do?"
Process guide has a troubleshooting section. See: `.claude/skills/sdd/process-guide/SKILL.md`

### "I'm not sure which phase I'm in"
See: [USING-SDD-PROCESS-GUIDE.md](USING-SDD-PROCESS-GUIDE.md) → "How to Use It in Your Workflow"

### "I want to understand the full theory"
See: [OPERATING-MODEL.md](sdd-mcp/operation%20model/OPERATING-MODEL.md)

### "I want a real example"
See: [sdd-example-workflow.md](sdd-example-workflow.md) — Guest checkout walkthrough

---

## 📂 Project Artifacts

These are created as you use SDD:

```
.agentic/
├── sdd/
│   └── initiatives/              (created by Phase 0)
│       └── [feature-name].yaml
├── spec-graph.json              (updated by Phase 3)
└── tasks/
    ├── backlog.yaml
    ├── in-progress.yaml
    └── done.yaml

openclaw-specs/
├── constitution/
│   └── policies.md              (created once, updated by Platform Arch)
└── features/
    └── [feature-name]/
        ├── feature-spec.md      (created by Phase 1)
        ├── component-spec-*.md
        ├── impl-spec-*.md       (created by Phase 2)
        └── verify.md            (created by Phase 3)
```

---

## ✨ Next Steps

1. **Read this page completely** ← You are here
2. **Pick your role** (PM, Architect, Developer, Verifier, DevOps)
3. **Follow the "By Role" section** above
4. **Install skills**: `agentic-agent skills ensure`
5. **Start your first feature**: `agentic-agent specifyify start "Name" --risk low`
6. **Open the process guide**: `cat .claude/skills/sdd/process-guide/SKILL.md`
7. **Follow Phase 0 steps**: numbered 0.1 through 0.5
8. **Share docs with your team**: reference this README

---

## 🎯 Core Principle

**Every feature goes through the same 4-phase process with 5 gates enforced at each phase.**

This ensures consistency, quality, and confidence in everything you ship.

---

**You're ready! Start with Phase 0.** 🚀
