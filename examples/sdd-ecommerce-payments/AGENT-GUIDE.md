# SDD Agent Guide: Role-Based Instructions

For AI agents operating autonomously, this guide provides concise step-by-step instructions per agent role.

---

## Analyst Agent

**Trigger:** Initiative defined, Risk Level = CRITICAL or HIGH

**Files to Read:**
- [01-initiative/initiative.md](./01-initiative/initiative.md)
- [03-analyst/prompts.md](./03-analyst/prompts.md) (example prompt structure)

**Your Job:**
1. Read initiative.md thoroughly
2. Call Platform MCP: `get_context_pack(initiative_name)`
3. Interview the team ONE QUESTION AT A TIME:
   - Problem evidence? (real data, not opinion)
   - Affected services? (names)
   - Blocking ADRs? (what decisions needed?)
   - Rollback concern? (how to undo?)
   - Observability? (what metrics matter?)
4. Produce `discovery.md` with:
   - Problem statement (metric + evidence)
   - Evidence (data points, not assumptions)
   - Affected components (explicit list)
   - Risk classification (with rationale)
   - Blocking ADRs identified
   - Recommended workflow (Quick/Standard/Full)

**Exit Gate (Self-Check):**
- [ ] Problem statement has concrete metric
- [ ] Evidence is real data
- [ ] Affected components named
- [ ] Risk classified with rationale
- [ ] Blocking ADRs identified
- [ ] Workflow matches risk level

**Output:**
- `03-analyst/discovery.md` (2000-3000 words)

**Time:** 3-5 days

**Next Agent:** Architect

---

## Architect Agent

**Trigger:** Discovery.md complete, Risk Level = STANDARD or FULL

**Files to Read:**
- [03-analyst/discovery.md](./03-analyst/discovery.md)
- [04-architect/prompts.md](./04-architect/prompts.md) (example prompt structure)

**Your Job:**
1. Read discovery.md thoroughly
2. Call Platform MCP: `get_context_pack(risk_level)`
3. For each affected service, call Component MCP:
   - `get_contracts()` — APIs, events, schemas
   - `get_invariants()` — business rules
   - `get_decisions()` — prior ADRs
4. Produce `feature-spec.md` (platform-level WHAT):
   - Metadata (implements, context_pack, blocked_by, status)
   - Problem statement (with Source:)
   - Goals / Non-Goals
   - User experience
   - Domain responsibilities
   - Cross-domain interactions
   - NFRs (performance, logging, metrics, tracing, PII)
   - Feature flag & rollback strategy
   - Acceptance criteria (5+ in Given/When/Then)
   - Blocking decisions (ADRs)
5. Produce `component-spec-<service>.md` (per-service HOW):
   - One per affected service
   - API contracts with examples
   - Data model
   - Acceptance criteria (testable)
   - Edge cases
6. Produce `ADR-###.md` for blocking decisions
   - Options with pros/cons
   - Recommendation
   - Decision timeline
7. Self-check all 5 gates (see gate-check skill)

**Gate Check (Self-Check):**
- [ ] Gate 1: Context Completeness (every section has Source:)
- [ ] Gate 2: Domain Validity (no cross-domain DB access)
- [ ] Gate 3: Integration Safety (contracts listed, consumers identified)
- [ ] Gate 4: NFR Compliance (logging, metrics, tracing, PII, performance)
- [ ] Gate 5: Ready to Implement (no blocking ADRs, ACs testable)

**Output:**
- `04-architect/feature-spec.md`
- `04-architect/component-spec-<service>.md` (multiple)
- `04-architect/ADR-###.md` (blocking decisions)

**Time:** 5-7 days

**Next Steps:**
1. ADR resolution (CTO, Security)
2. Developer agents (parallel per component)

---

## ADR Reviewer

**Trigger:** ADR created by Architect

**Files to Read:**
- `04-architect/ADR-###.md`

**Your Job:**
1. Review ADR options and recommendation
2. Gather team (CTO, Security, Finance, etc. as needed)
3. Make decision: Option A / B / C
4. Update ADR:
   - Status: APPROVED
   - Decision: [selected option]
   - Decision Date: [today]
5. Unblock all specs: Set `blocked_by: []`
6. Notify developers ADR is resolved

**Output:**
- ADR marked APPROVED
- All component-specs with this ADR unblocked

**Time:** 3-5 days

**Next Agent:** Developer

---

## Developer Agent

**Trigger:** Component-spec ready, `blocked_by: []` (no blocking ADRs)

**Files to Read:**
- [05-developer/prompts.md](./05-developer/prompts.md) (example structure)
- `04-architect/component-spec-<service>.md` (your assigned component)
- `04-architect/feature-spec.md` (overall context)

**Your Job:**
1. Read component-spec thoroughly
2. Check `blocked_by` field — STOP if non-empty, wait for ADR
3. Call Component MCP:
   - `get_patterns()` — canonical implementation patterns
   - `get_decisions()` — prior ADRs, lessons learned
4. Produce `impl-spec-<service>.md`:
   - Data model (database schema with constraints)
   - Code changes (exact functions/modules to create/modify)
   - Edge cases (table with 4+ scenarios and handling)
   - Observability (logging, metrics, tracing, alerts in code)
   - Rollout plan (feature flag, canary, rollback procedure)
5. Produce `tasks.yaml`:
   - 8-10 actionable development tasks
   - Dependencies between tasks
   - Effort estimates (hours/days)
   - Acceptance criteria per task
6. Self-check all 5 gates

**Gate Check (Self-Check):**
- [ ] Gate 1: Context Completeness
- [ ] Gate 2: Domain Validity
- [ ] Gate 3: Integration Safety
- [ ] Gate 4: NFR Compliance
- [ ] Gate 5: Ready to Implement

**Output:**
- `05-developer/impl-spec-<service>.md` (2000-3000 words)
- `05-developer/tasks.yaml` (8-10 tasks)

**Time:** 3-4 days per component

**Parallelization:** Multiple developers can work simultaneously (one per component)

**Next Agent:** Verifier (after code is written + tested)

---

## Verifier Agent

**Trigger:** Code merged, all tests passing

**Files to Read:**
- [06-verifier/prompts.md](./06-verifier/prompts.md)
- `04-architect/feature-spec.md` (all ACs listed)
- `04-architect/component-spec-<service>.md` (all ACs)
- `05-developer/impl-spec-<service>.md` (code structure)
- Test results, logs, metrics (from CI/CD)

**Your Job:**
1. For every acceptance criterion, gather evidence:
   - Test that proves it (test name + line + output)
   - Log showing feature works
   - Metric/trace showing behavior
2. Check code quality:
   - Tests: # passing, # failing
   - Coverage: % (must be > 80%)
   - Lint: # errors, # warnings (must be 0 errors)
   - Build: success/failure
3. Check blocking conditions:
   - Any AC untestable? → BLOCK
   - Any test FAILED? → BLOCK
   - Coverage < 80%? → BLOCK
   - Payment/PII touched without human approval? → BLOCK
4. Produce `verify.md`:
   - Metadata (spec ID, verified by, date, status)
   - AC verification (AC1-AC7 with evidence)
   - Code quality checks (tests, lint, coverage, build)
   - Observability verification (logging, metrics, tracing, alerts)
   - Blocking conditions (human approval required)
   - Sign-off (ready for production?)

**Blocking Conditions (Don't Merge If):**
- [ ] Any AC is UNTESTABLE
- [ ] Any test FAILS
- [ ] Lint or build FAILS
- [ ] Coverage < 80%
- [ ] Change touches payment/auth/PII without human approval
- [ ] Spec Graph not updated to Done

**Output:**
- `06-verifier/verify.md`
- `agentic-agent specify sync-graph` (update spec-graph)
- `agentic-agent openspec complete <change-id>` (mark Done)

**Time:** 2-3 days

**Next Steps:**
1. If BLOCKED: Return to Architect/Developer with remediation
2. If READY: Request human approvals (Security, Payments, Legal if payment feature)
3. After approvals: Deploy to production

---

## Quick Reference: File Structure

```
examples/sdd-ecommerce-payments/
├── 00-setup/
│   └── agnostic-agent.yaml
│
├── 01-initiative/              ← Analyst reads, PM writes
│   └── initiative.md
│
├── 02-risk-assessment/         ← Tech Lead determines workflow
│   └── risk-classification.md
│
├── 03-analyst/                 ← ANALYST ROLE
│   ├── prompts.md              (example prompt to invoke analyst)
│   └── discovery.md            (analyst output)
│
├── 04-architect/               ← ARCHITECT ROLE
│   ├── prompts.md
│   ├── feature-spec.md         (architect output: platform WHAT)
│   ├── component-spec-billing.md
│   ├── component-spec-payments.md
│   └── ADR-001-card-vault.md   (blocking decision)
│
├── 05-developer/               ← DEVELOPER ROLE (parallel)
│   ├── prompts.md
│   ├── impl-spec-billing.md    (developer output: implementation HOW)
│   └── tasks.yaml              (developer output: task breakdown)
│
├── 06-verifier/                ← VERIFIER ROLE
│   ├── prompts.md
│   └── verify.md               (verifier output: evidence)
│
├── 07-gate-checks/
│   └── gate-check-output.md    (gate 1-5 validation results)
│
├── README.md                   (human walkthrough)
└── AGENT-GUIDE.md             (this file)
```

---

## Workflow State Machine

```
START
  ↓
Analyst reads discovery
  ↓
Analyst interviews team → discovery.md
  ↓
Architect reads discovery
  ↓
Architect designs → feature-spec + component-specs
  ↓
Architect creates ADR-001 (blocking)
  ↓
ADR Reviewer decides → ADR APPROVED
  ↓
Developer reads component-spec (blocked_by empty)
  ↓
Developer implements → impl-spec + tasks
  ↓
Engineers code + test
  ↓
Verifier checks evidence → verify.md
  ↓
Human approvals? (Security, Payments, Legal)
  ├─ YES → Deploy
  └─ NO  → Return to Developer
```

---

## Gate Checks: What Agents Do

Each agent self-checks gates before handing off:

**Analyst** → Gate 1-5 on discovery.md
**Architect** → Gate 1-5 on feature-spec + component-specs
**Developer** → Gate 1-5 on impl-specs
**Verifier** → Verify all ACs with evidence

Run: `agentic-agent specify gate-check <spec-id> --format text`

---

## Key Files to Know

| File | What It Is | Read When |
|------|-----------|-----------|
| [01-initiative/initiative.md](./01-initiative/initiative.md) | Problem statement | Before Analyst phase |
| [02-risk-assessment/risk-classification.md](./02-risk-assessment/risk-classification.md) | Determines workflow type | Before Analyst phase |
| [03-analyst/discovery.md](./03-analyst/discovery.md) | Team insights + evidence | Before Architect phase |
| [04-architect/feature-spec.md](./04-architect/feature-spec.md) | Overall solution design | Before Developer phase |
| [04-architect/component-spec-*.md](./04-architect/) | Per-service design | Before Developer phase (component-specific) |
| [04-architect/ADR-001-card-vault.md](./04-architect/ADR-001-card-vault.md) | Blocking decision | Before Developer phase (blocks implementation) |
| [05-developer/impl-spec-*.md](./05-developer/) | Code structure | During development |
| [05-developer/tasks.yaml](./05-developer/tasks.yaml) | Task decomposition | During development (assign tasks) |
| [06-verifier/verify.md](./06-verifier/verify.md) | Test evidence | Before deployment |
| [07-gate-checks/gate-check-output.md](./07-gate-checks/gate-check-output.md) | Gate validation | Before handoff (architect), before deployment (verifier) |

---

## Commands You'll Use

```bash
# Setup
agentic-agent specify agents install

# Start initiative
agentic-agent specify start "Subscription Billing" --risk critical

# View workflow
agentic-agent specify workflow show <initiative-id>

# Gate check
agentic-agent specify gate-check SBL-2025-Q1 --format text

# List ADRs
agentic-agent specify adr list

# Resolve ADR
agentic-agent specify adr resolve ADR-001

# Sync graph
agentic-agent specify sync-graph

# Complete initiative
agentic-agent openspec complete SBL-2025-Q1
```

---

## Success Criteria Per Agent

**Analyst:** Discovery.md complete, all exit gates pass
**Architect:** Feature-spec + component-specs complete, all gates pass (or blocked by ADR)
**Developer:** impl-spec + tasks.yaml complete, code written, all tests passing
**Verifier:** verify.md complete, all ACs verified, human approvals obtained

---

**Use this guide to execute each phase autonomously. Start with your role above.**
