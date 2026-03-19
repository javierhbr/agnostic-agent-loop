# SDD Unified Methodology — Skills Directory (25+ Skills)

A complete reference of all SDD skills organized by the 5-phase model.

---

## Installation

```bash
agentic-agent skills ensure
# or install SDD skills explicitly
agentic-agent skills install sdd --tool claude-code
```

Skills are installed to: `.claude/skills/sdd/` (role-specific) and `.claude/skills/` (cross-cutting, orchestrators).

---

## Phase-at-a-Glance

```
Phase 0          Phase 1        Phase 2             Phase 3           Phase 4
INITIATIVE       DISCOVERY      ARCHITECTURE        IMPLEMENTATION    VERIFICATION
& RISK           (Full only)    & SPECIFICATION     & CODE            & PROOF

initiative ─┐                   architect ────────> developer ──────> verifier
risk ────────┤                  component-spec      tdd               gate-check
workflow ────┘    analyst       platform-const.     sdd-openspec
                                platform-context.   sdd-bmad
                                adr                 sdd-speckit

              ─── OFF-PATH: hotfix (lightweight, any time) ───────>

              ─── CROSS-CUTTING: explain-code, local-doc, ────────>
                  stakeholder-comm, superpowers-bridge

              ─── ORCHESTRATORS: process-guide, sdd-methodology ──>
                  uncle-sdd-agent, unified-sdd-skill

              ─── META: tier-enforcer ────────────────────────────>
```

---

## Quick Lookup Table

| # | Skill | Phase | Role | Key Artifact | Skill File |
|---|-------|-------|------|-------------|------------|
| 1 | initiative-definition | 0 | PM | initiative.md | `sdd/initiative-definition/SKILL.md` |
| 2 | risk-assessment | 0 | PM | Risk classification | `sdd/risk-assessment/SKILL.md` |
| 3 | workflow-router | 0 | PM / Team Lead | QUICK/STANDARD/FULL decision | `sdd/workflow-router/SKILL.md` |
| 4 | analyst | 1 | Business Analyst | discovery.md | `sdd/analyst/SKILL.md` |
| 5 | architect | 2 | Solution Architect | feature-spec.md, component-spec.md | `sdd/architect/SKILL.md` |
| 6 | component-spec | 2 | Component Team Lead | Component spec | `sdd/component-spec/SKILL.md` |
| 7 | platform-constitution | 2 | Platform Architect | constitution.md | `sdd/platform-constitution/SKILL.md` |
| 8 | platform-contextualizer | 2 | Platform Architect | Brownfield baseline | `platform-contextualizer-skill/SKILL.md` |
| 9 | adr | 2 | Architect / Tech Lead | ADR documents | `sdd/adr/SKILL.md` |
| 10 | developer | 3 | Developer | impl-spec.md, tasks.yaml | `sdd/developer/SKILL.md` |
| 11 | tdd | 3 | Developer | RED-GREEN-REFACTOR cycle | `tdd/SKILL.md` |
| 12 | sdd-openspec | 3 | Architect / Developer | proposal.md, design.md, tasks.md, archive | `sdd-openspec/SKILL.md` |
| 13 | sdd-bmad | 3 | PM / Architect / Dev | Progressive context, track selection | `sdd-bmad/SKILL.md` |
| 14 | sdd-speckit | 3 | PM / Architect / Dev | Constitution, specs, tasks (8-phase) | `sdd-speckit/SKILL.md` |
| 15 | verifier | 4 | QA | verify.md | `sdd/verifier/SKILL.md` |
| 16 | gate-check | 4 | Any | 5-gate enforcement report | `sdd/gate-check/SKILL.md` |
| 17 | hotfix | Off-path | Developer / Team Lead | Lightweight hotfix spec | `sdd/hotfix/SKILL.md` |
| 18 | explain-code | Cross-cut | Any | ASCII diagrams, analogies | `explain-code-skill/SKILL.md` |
| 19 | local-doc | Cross-cut | Any | SQLite FTS5 search results | (planned) |
| 20 | stakeholder-communication | Cross-cut | PM | Business-friendly status | `sdd/stakeholder-communication/SKILL.md` |
| 21 | superpowers-bridge | Cross-cut | Dev / Architect | Tool selection guidance | `superpowers-bridge/SKILL.md` |
| 22 | process-guide | Orchestrator | Any | Full lifecycle guide | `sdd/process-guide/SKILL.md` |
| 23 | sdd-methodology-tools | Orchestrator | Any | Central skill router | `sdd/SKILLS.md` |
| 24 | uncle-sdd-agent | Orchestrator | Any | 5-phase: Platform > Assess > Specify > Plan > Deliver | (planned) |
| 25 | unified-sdd-skill | Orchestrator | Any | Platform > Route > Specify > Plan > Deliver | `unified-sdd-skill/SKILL.md` |
| 26 | tier-enforcer | Meta | Skill Maintainer | 3-tier compliance report | `tier-enforcer/SKILL.md` |

---

## Skills by Phase

### Phase 0 — Initiative and Risk

Purpose: Define what to build and how risky it is, then route to the right workflow.

| Skill | Trigger | Role | Produces | When to Use |
|-------|---------|------|----------|-------------|
| initiative-definition | `sdd:initiative` | PM | initiative.md with problem, goals, success metrics | Starting any new feature or epic |
| risk-assessment | `sdd:risk` | PM | Risk classification (Low / Medium / High / Critical) | After defining the initiative, before architecture |
| workflow-router | `sdd:workflow` | PM / Team Lead | QUICK, STANDARD, or FULL decision | After risk classification, to determine how much ceremony |

Workflow routing outcome:

```
Risk Level     Workflow    Participants
---------------------------------------------
LOW            QUICK       Developer + Verifier
MEDIUM         STANDARD    Architect + Developers + Verifier
HIGH/CRITICAL  FULL        Analyst + Architect + Developers + Verifier + human approval
```

---

### Phase 1 — Discovery (Full Workflow Only)

Purpose: Deep investigation before committing to architecture. Skipped for QUICK and STANDARD.

| Skill | Trigger | Role | Produces | When to Use |
|-------|---------|------|----------|-------------|
| analyst | `sdd:analyst` | Business Analyst | discovery.md (interviews, constraints, unknowns) | HIGH/CRITICAL risk initiatives only |

---

### Phase 2 — Architecture and Specification

Purpose: Design the solution. Create specs that developers can implement independently.

| Skill | Trigger | Role | Produces | When to Use |
|-------|---------|------|----------|-------------|
| architect | `sdd:architect` | Solution Architect | feature-spec.md + component-spec.md | Designing the feature from the initiative |
| component-spec | `sdd:component-spec` | Component Team Lead | Detailed component specification | Breaking feature-spec into per-service specs |
| platform-constitution | `sdd:constitution` | Platform Architect | constitution.md (governance policies) | Establishing platform-wide rules and constraints |
| platform-contextualizer | `platform-contextualizer` | Platform Architect | Brownfield baseline capture | Documenting existing platform state before changes |
| adr | `sdd:adr` | Architect / Tech Lead | Architecture Decision Records | Recording and enforcing architectural decisions |

---

### Phase 3 — Implementation

Purpose: Build the code. Multiple implementation strategies available depending on project type.

| Skill | Trigger | Role | Produces | When to Use |
|-------|---------|------|----------|-------------|
| developer | `sdd:developer` | Developer | impl-spec.md + tasks.yaml + code | Implementing a component from its spec |
| tdd | `tdd` | Developer | RED-GREEN-REFACTOR test cycle | Writing tests first, then making them pass |
| sdd-openspec | `sdd-openspec` | Architect / Developer | proposal.md, design.md, tasks.md, archive | Spec-driven development from requirements |
| sdd-bmad | `sdd-bmad` | PM / Architect / Dev | Progressive context docs, track selection | Role-based execution with progressive planning |
| sdd-speckit | `sdd-speckit` | PM / Architect / Dev | Constitution, specs, tasks (8-phase workflow) | Turning ideas into executable specs and task lists |

---

### Phase 4 — Verification

Purpose: Prove the implementation meets all acceptance criteria before merge.

| Skill | Trigger | Role | Produces | When to Use |
|-------|---------|------|----------|-------------|
| verifier | `sdd:verifier` | QA | verify.md with evidence for every AC | After all developers finish, before merge |
| gate-check | `sdd:gate-check` | Any | 5-gate pass/fail report | At every phase boundary and before handoffs |

---

### Off-Path

| Skill | Trigger | Role | Produces | When to Use |
|-------|---------|------|----------|-------------|
| hotfix | `sdd:hotfix` | Developer / Team Lead | Lightweight hotfix spec | Production bugs requiring minimal ceremony |

---

### Cross-Cutting Skills

These skills are used at any phase and are not bound to a specific workflow step.

| Skill | Trigger | Role | Produces | When to Use |
|-------|---------|------|----------|-------------|
| explain-code | `explain-code` | Any | ASCII diagrams, analogies, visual explanations | Understanding unfamiliar code or architecture |
| local-doc | (planned) | Any | SQLite FTS5 search results across specs | Searching the spec graph locally |
| stakeholder-communication | `sdd:stakeholder` | PM | Business-friendly status updates | Reporting progress to non-technical stakeholders |
| superpowers-bridge | `superpowers-bridge` | Dev / Architect | CLI vs Superpowers tool selection | Deciding between CLI automation and Superpowers plugin |

---

### Orchestrators

Higher-level skills that coordinate multiple phase-specific skills into end-to-end workflows.

| Skill | Trigger | Role | Scope | When to Use |
|-------|---------|------|-------|-------------|
| process-guide | `sdd:process-guide` | Any | Full 5-phase lifecycle guide | Starting any SDD workflow, need step-by-step direction |
| sdd-methodology-tools | `sdd:*` | Any | Central skill router for all 15+ SDD sub-skills | Dispatching to the right SDD skill by trigger |
| uncle-sdd-agent | (planned) | Any | Platform-scale orchestrator (5 phases) | Large platform initiatives spanning multiple teams |
| unified-sdd-skill | `unified-sdd` | Any | Lighter orchestrator combining OpenSpec + BMAD + SpecKit | Single-team projects needing combined methodology |

---

### Meta

| Skill | Trigger | Role | Produces | When to Use |
|-------|---------|------|----------|-------------|
| tier-enforcer | `tier-check` / `tier-enforcer` | Skill Maintainer | 3-tier layered context compliance report | Auditing or creating skill files for tier compliance |

---

## Skill Interaction Map

```
                          ORCHESTRATORS
                    (process-guide, unified-sdd)
                              |
          +-------------------+-------------------+
          |                   |                   |
          v                   v                   v
     Phase 0             Phase 2             Phase 3
  initiative ──> risk ──> architect ──────> developer
                  |       component-spec    tdd
                  v       adr               sdd-openspec
            workflow-router                 sdd-bmad
                  |                         sdd-speckit
                  v                              |
            [QUICK?] ──yes──> skip to Phase 3    |
            [FULL?] ──yes──> Phase 1 (analyst)   |
                                                 v
                                            Phase 4
                                         verifier + gate-check
                                                 |
                                                 v
                                             [MERGE]

Cross-cutting skills (explain-code, stakeholder-communication,
superpowers-bridge) can be invoked at any point in the flow.

gate-check enforces quality at EVERY phase boundary:
  Phase 0 exit --> Phase 2 exit --> Phase 3 exit --> Phase 4 exit
```

---

## Workflow Decision Summary

| Question | Answer | Action |
|----------|--------|--------|
| New feature or epic? | Yes | Start at Phase 0: initiative-definition |
| Risk level unknown? | Yes | Use risk-assessment (5-question interview) |
| Low risk, small scope? | Yes | QUICK workflow: developer + verifier only |
| Medium risk? | Yes | STANDARD workflow: architect + developers + verifier |
| High or critical risk? | Yes | FULL workflow: analyst first, then architect + developers + verifier |
| Production bug? | Yes | Use hotfix (off-path, minimal spec) |
| Platform-wide change? | Yes | Use platform-constitution + platform-contextualizer |
| Need architectural decision recorded? | Yes | Use adr before implementation |
| Which implementation strategy? | Depends | developer (standard), tdd (test-first), sdd-openspec/bmad/speckit (spec-driven) |
| Ready to merge? | Maybe | Use verifier + gate-check to prove it |
| Not sure what to do? | -- | Start with process-guide |

---

## The Five Gates

Every phase boundary is enforced by gate-check. All five gates must PASS before handoff.

| Gate | Name | Checks For |
|------|------|-----------|
| 1 | Context Completeness | Metadata fields present (implements, context_pack, status) |
| 2 | Domain Validity | Domain invariants respected, no cross-domain violations |
| 3 | Integration Safety | Contract changes declared, consumers identified |
| 4 | NFR Compliance | Observability, security, performance requirements declared |
| 5 | Ready to Implement | No blocking ADRs, all ACs in GWT format, unambiguous |

Run gate checks with:

```bash
agentic-agent specify gate-check SPEC-[ID]
```

---

## File Locations

```
.claude/skills/
├── sdd/                                    # Phase-specific SDD skills
│   ├── SKILLS.md                           # SDD sub-router
│   ├── initiative-definition/SKILL.md      # Phase 0 — PM
│   ├── risk-assessment/SKILL.md            # Phase 0 — PM
│   ├── workflow-router/SKILL.md            # Phase 0 — PM / Team Lead
│   ├── analyst/SKILL.md                    # Phase 1 — Business Analyst
│   ├── architect/SKILL.md                  # Phase 2 — Solution Architect
│   ├── component-spec/SKILL.md             # Phase 2 — Component Team Lead
│   ├── platform-constitution/SKILL.md      # Phase 2 — Platform Architect
│   ├── adr/SKILL.md                        # Phase 2 — Architect / Tech Lead
│   ├── developer/SKILL.md                  # Phase 3 — Developer
│   ├── verifier/SKILL.md                   # Phase 4 — QA
│   ├── gate-check/SKILL.md                 # Phase 4 — Any
│   ├── hotfix/SKILL.md                     # Off-path
│   ├── stakeholder-communication/SKILL.md  # Cross-cutting — PM
│   └── process-guide/SKILL.md              # Orchestrator
├── tdd/SKILL.md                            # Phase 3 — Developer
├── sdd-openspec/SKILL.md                   # Phase 3 — Architect / Developer
├── sdd-bmad/SKILL.md                       # Phase 3 — PM / Architect / Dev
├── sdd-speckit/SKILL.md                    # Phase 3 — PM / Architect / Dev
├── explain-code-skill/SKILL.md             # Cross-cutting
├── platform-contextualizer-skill/SKILL.md  # Phase 2 — Platform Architect
├── superpowers-bridge/SKILL.md             # Cross-cutting
├── unified-sdd-skill/SKILL.md              # Orchestrator
├── tier-enforcer/SKILL.md                  # Meta
└── SKILLS.md                               # Master router (all skills)
```

---

## Quick Start

```bash
# 1. Install skills
agentic-agent skills ensure

# 2. Start with the process guide if unsure
#    Read: .claude/skills/sdd/process-guide/SKILL.md

# 3. Or jump directly to the skill for your role:
#    PM:        sdd:initiative -> sdd:risk -> sdd:workflow
#    Analyst:   sdd:analyst (Full workflow only)
#    Architect: sdd:architect -> sdd:component-spec -> sdd:gate-check
#    Developer: sdd:developer (or tdd, sdd-openspec, sdd-bmad, sdd-speckit)
#    QA:        sdd:verifier -> sdd:gate-check
#    Hotfix:    sdd:hotfix

# 4. Gate-check at every phase boundary
agentic-agent specify gate-check SPEC-[ID]
```

---

## Summary

This directory covers **26 skills** organized across a **5-phase model** with **5 quality gates** enforced at every boundary. Skills range from single-role phase tools (initiative-definition, verifier) to full-lifecycle orchestrators (process-guide, unified-sdd-skill).

**Start with:** `process-guide` -- it tells you which skill to use at each step.

**Key principle:** One trigger, one skill, one output. Gates must PASS before moving to the next phase.
