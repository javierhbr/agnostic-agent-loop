# SDD v3.0 Skills Directory

A complete map of all 15 SDD skills available in your project and when to use each one.

---

## Installation

First, ensure SDD skills are installed:

```bash
agentic-agent skills ensure
# or
agentic-agent skills install sdd --tool claude-code
```

Skills are then available at: `.claude/skills/sdd/`

---

## The Complete Skill Map

### 🎯 **START HERE: Process Guide**

| Skill | Audience | Purpose | Use When |
|-------|----------|---------|----------|
| **process-guide** | Everyone | Complete 4-phase walkthrough | Starting any new feature |

**This skill is your map.** It orchestrates all others by telling you which skill to use at each phase.

**Reference:** `docs/USING-SDD-PROCESS-GUIDE.md`

---

### 🏢 **PHASE 0: Initiative Definition (Product Manager)**

| Skill | Audience | Purpose | Use When |
|-------|----------|---------|----------|
| **initiative-definition** | Product Manager | Define problem, goals, and success metrics | Starting a new initiative |
| **risk-assessment** | Product Manager | Interview guide for classifying risk level | After defining problem, before handing to engineers |
| **stakeholder-communication** | Product Manager | Translate SDD terms to business language | Communicating status to non-technical stakeholders |

**Example workflow:**
```bash
# 1. Use initiative-definition to write the problem statement
#    Output: Problem statement with metrics

# 2. Use risk-assessment to classify the work
#    Answer: 5 questions about risk
#    Output: Risk classification (Low/Medium/High/Critical)

# 3. Use stakeholder-communication to report status
#    Output: Plain-language updates for executives
```

---

### 🏛️ **PHASE 0.5: Platform Governance (Platform Architect)**

| Skill | Audience | Purpose | Use When |
|-------|----------|---------|----------|
| **platform-constitution** | Platform Architect | Author and maintain platform governance policies | Establishing platform-wide rules and constraints |

**Outputs:**
- `openclaw-specs/constitution/policies.md` — The governance document all specs must follow

---

### 🏗️ **PHASE 1: Architecture Design (Solution Architect)**

| Skill | Audience | Purpose | Use When |
|-------|----------|---------|----------|
| **workflow-router** | Architect | Decide if you need QUICK/STANDARD/FULL workflow | Starting architecture based on initiative |
| **architect** | Architect | Create feature-spec and component-spec files | Designing the solution from the initiative |
| **gate-check** | Architect | Validate architecture against 5 gates | After creating specs, before handing to developers |

**Example workflow:**
```bash
# 1. Use workflow-router to understand which workflow to follow
#    Decision: QUICK/STANDARD/FULL based on risk

# 2. Use architect skill to create specs
#    Outputs: feature-spec.md + component-spec.md (x4)
#    Reference: Platform Constitution for constraints

# 3. Use gate-check to validate before handoff
#    Validation: All 5 gates PASS
```

---

### 💻 **PHASE 2: Implementation (Developers, Parallel per Component)**

| Skill | Audience | Purpose | Use When |
|-------|----------|---------|----------|
| **developer** | Developer | Implement component from component-spec | Assigned a component task |

**Example workflow (per developer):**
```bash
# 1. Read component-spec for your service
# 2. Use developer skill to understand what to produce
#    Output: impl-spec.md + tasks.yaml + code
# 3. Gate check your implementation
#    Validation: Gates 4 & 5 PASS
```

---

### 🔍 **PHASE 3: Verification & Deployment (Verifier)**

| Skill | Audience | Purpose | Use When |
|-------|----------|---------|----------|
| **verifier** | Verifier | Verify all ACs with evidence, create verify.md | All developers finished, ready to merge |

**Example workflow:**
```bash
# 1. Use verifier skill to verify every AC
#    Check: GWT acceptance criteria all PASS
#    Evidence: Test results, metrics, logs
#
# 2. Create verify.md proof document
#    Output: Observable evidence for each AC
#
# 3. Gate check final state
#    Validation: All 5 gates PASS
#
# 4. Sync spec graph and merge
```

---

### 🏛️ **PHASE 1 SUPPORT: Original SDD Skills (Maintained)**

| Skill | Purpose | Use When |
|-------|---------|----------|
| **platform-spec** | Define platform-level specifications | Defining cross-team platform changes |
| **component-spec** | Define component-level specifications | Architecture specifies a new service/boundary |
| **adr** | Create Architecture Decision Records | Need to document and block decisions |
| **hotfix** | Quick fixes with minimal spec | Bug fix scoped to one service |

---

## Decision Matrix: Which Skill to Use?

```
I'm a...                    My current task is...                        Use skill...
────────────────────────────────────────────────────────────────────────────────────────
Product Manager             Define a new initiative                      initiative-definition
Product Manager             Classify risk level                          risk-assessment
Product Manager             Report status to executives                  stakeholder-communication

Platform Architect          Establish governance policies                platform-constitution
Platform Architect          Check if specs meet platform rules           gate-check

Solution Architect          Decide workflow type (QUICK/STD/FULL)        workflow-router
Solution Architect          Design the feature                           architect
Solution Architect          Check if design is ready for dev             gate-check

Developer (any role)        Implement a component                        developer
Developer (any role)        Check if my code is gate-ready               gate-check

Verifier / QA Engineer      Prove feature works                          verifier
Verifier / QA Engineer      Check overall readiness for merge            gate-check
Verifier / QA Engineer      Deploy feature safely                        process-guide (Phase 4)

Anyone                      Confused about what to do next               process-guide
```

---

## The Five Gates (Checked by Each Skill)

Every skill's exit criteria includes gate checks. The five gates are:

| Gate | Name | Checks For | Run Via |
|------|------|-----------|---------|
| **1** | Context Completeness | Metadata fields present (implements, context_pack, status) | `agentic-agent sdd gate-check` |
| **2** | Domain Validity | Domain invariants respected, no cross-domain DB access | `agentic-agent sdd gate-check` |
| **3** | Integration Safety | Contract changes declared, consumers identified | `agentic-agent sdd gate-check` |
| **4** | NFR Compliance | Observability, security, performance declared | `agentic-agent sdd gate-check` |
| **5** | Ready to Implement | No blocking ADRs, all ACs in GWT format, unambiguous | `agentic-agent sdd gate-check` |

---

## File Locations

All skills available at:

```
.claude/skills/sdd/
├── analyst/SKILL.md                    (Full workflow only)
├── architect/SKILL.md
├── developer/SKILL.md
├── verifier/SKILL.md
├── workflow-router/SKILL.md
├── initiative-definition/SKILL.md      (PM)
├── risk-assessment/SKILL.md            (PM)
├── stakeholder-communication/SKILL.md  (PM)
├── platform-constitution/SKILL.md      (Platform Arch)
├── platform-spec/SKILL.md
├── component-spec/SKILL.md
├── gate-check/SKILL.md
├── adr/SKILL.md
├── hotfix/SKILL.md
└── process-guide/SKILL.md              (START HERE)
```

---

## Quick Start by Role

### If You're a Product Manager

```bash
1. Read: .claude/skills/sdd/initiative-definition/SKILL.md
2. Read: .claude/skills/sdd/risk-assessment/SKILL.md
3. Run: agentic-agent sdd start "Feature Name" --risk [level]
4. Reference: .claude/skills/sdd/stakeholder-communication/SKILL.md for updates
```

### If You're a Solution Architect

```bash
1. Read: .claude/skills/sdd/process-guide/SKILL.md (Phase 1)
2. Read: .claude/skills/sdd/architect/SKILL.md
3. Create: feature-spec.md + component-spec.md
4. Run: agentic-agent sdd gate-check SPEC-[ID]
5. When all gates PASS: Create parallel dev tasks
```

### If You're a Developer

```bash
1. Read: .claude/skills/sdd/process-guide/SKILL.md (Phase 2)
2. Run: agentic-agent task claim [TASK-ID]
3. Read: Your assigned component-spec.md (from architect)
4. Reference: .claude/skills/sdd/developer/SKILL.md
5. Produce: impl-spec.md + code + tests
6. Run: agentic-agent sdd gate-check SPEC-[SERVICE]-IMPL
```

### If You're a Verifier/QA

```bash
1. Read: .claude/skills/sdd/process-guide/SKILL.md (Phase 3)
2. Reference: .claude/skills/sdd/verifier/SKILL.md
3. Run: agentic-agent validate (all tests)
4. Create: verify.md with evidence
5. Run: agentic-agent sdd gate-check SPEC-[ID] (final check)
6. Run: agentic-agent sdd sync-graph
```

### If You're Confused

```bash
1. Read: .claude/skills/sdd/process-guide/SKILL.md
   (This guides you through all 4 phases with exact steps)
```

---

## Example: Complete Journey

```bash
# Phase 0: PM uses process-guide + initiative-definition + risk-assessment
agentic-agent sdd start "Guest Checkout" --risk medium
# Initiative created ✓

# Phase 1: Architect uses process-guide + architect
# (Reads initiative, Platform Constitution, creates specs)
agentic-agent sdd gate-check SPEC-CHECKOUT-GUEST
# All gates PASS ✓

# Phase 2: 4 Developers use process-guide + developer (parallel)
agentic-agent task claim TASK-001
# (Implements their component)
agentic-agent task complete TASK-001
# Code merged ✓

# Phase 3: Verifier uses process-guide + verifier
# (Verifies every AC, creates verify.md)
agentic-agent validate
agentic-agent sdd gate-check SPEC-CHECKOUT-GUEST
agentic-agent sdd sync-graph
# Merged to main ✓

# Phase 4: Deploy using process-guide (Phase 4)
agentic-agent deploy --feature-flags-all-off
agentic-agent flags set GuestCheckoutEnabled=100pct
# Feature shipped ✓
```

---

## Summary

You now have **15 integrated skills** that guide you through **4 complete phases** with **5 quality gates** enforced at each step.

**Start with:** `process-guide/SKILL.md` — it tells you which skill to use when.

**Key principle:** Each skill has clear exit criteria and gate checks. You don't move to the next phase until gates PASS.

Happy shipping! 🚀
