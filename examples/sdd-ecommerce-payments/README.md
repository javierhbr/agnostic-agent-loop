# SDD Methodology: Full End-to-End Example

**Scenario:** ShopFlow e-commerce platform adds **Subscription Billing** feature (recurring payments + card storage)

**Risk Level:** `CRITICAL` (payment + PII) → **Full Workflow**

**Timeline:** 3-4 weeks (Analyst → Architect → Developer × 5 → Verifier)

---

## Table of Contents

1. [What This Example Shows](#what-this-example-shows)
2. [Quick Start: How to Read This](#quick-start-how-to-read-this)
3. [Step-by-Step Walkthrough](#step-by-step-walkthrough)
4. [Key Artifacts](#key-artifacts)
5. [For Humans vs. AI Agents](#for-humans-vs-ai-agents)
6. [Try It Yourself](#try-it-yourself)

---

## What This Example Shows

This directory contains a **complete, realistic worked example** of the SDD (Spec-Driven Development) methodology in action. It shows:

- **Every step** from problem definition to verification
- **Every artifact** produced at each phase
- **Exact prompts** to invoke each agent role (Analyst, Architect, Developer, Verifier)
- **Real gate checks** showing what passing and failing look like
- **How decisions (ADRs) block and unblock implementation**
- **How parallel development works** (multiple developer agents simultaneously)
- **Verification with observable evidence** (tests, logs, metrics)

**This is not a toy example.** It includes realistic:
- Multi-service coordination (5 services affected)
- Payment system integration (Stripe vault decision)
- PII/compliance requirements (PCI-DSS)
- Observability specs (logging, metrics, tracing, alerts)
- Edge cases (charge failures, retries, rollback)
- Risk assessment (critical payment data)

---

## Quick Start: How to Read This

### Option 1: Timeline (Humans)

Read in this order to understand the workflow:

1. [01-initiative/initiative.md](./01-initiative/initiative.md) — PM defines the problem & goals
2. [02-risk-assessment/risk-classification.md](./02-risk-assessment/risk-classification.md) — Determines workflow type (why Full?)
3. [03-analyst/prompts.md](./03-analyst/prompts.md) & [discovery.md](./03-analyst/discovery.md) — Analyst interviews team & produces discovery
4. [04-architect/prompts.md](./04-architect/prompts.md), [feature-spec.md](./04-architect/feature-spec.md) & component specs — Architect designs the solution
5. [04-architect/ADR-001-card-vault.md](./04-architect/ADR-001-card-vault.md) — Blocking decision resolved
6. [05-developer/prompts.md](./05-developer/prompts.md), [impl-spec-billing.md](./05-developer/impl-spec-billing.md), [tasks.yaml](./05-developer/tasks.yaml) — Developers implement
7. [07-gate-checks/gate-check-output.md](./07-gate-checks/gate-check-output.md) — Verify all specs pass gates
8. [06-verifier/prompts.md](./06-verifier/prompts.md) & [verify.md](./06-verifier/verify.md) — Verifier confirms all ACs with evidence

### Option 2: Role-Based (AI Agents)

See [AGENT-GUIDE.md](./AGENT-GUIDE.md) for concise instructions per agent role.

### Option 3: Deep Dive (Architects/Engineers)

- Read [04-architect/feature-spec.md](./04-architect/feature-spec.md) to understand the complete feature design
- Read [05-developer/impl-spec-billing.md](./05-developer/impl-spec-billing.md) to see implementation detail
- Read [06-verifier/verify.md](./06-verifier/verify.md) to see how verification works

---

## Step-by-Step Walkthrough

### Step 0: Setup

Configure the project:

```bash
# Install SDD agent skills
agentic-agent specify agents install --dir .claude/agents

# View config (see 00-setup/agnostic-agent.yaml)
cat agnostic-agent.yaml
```

---

### Step 1: Initiative Definition

**Who:** Product Manager

**Input:** Business problem

**Output:** [01-initiative/initiative.md](./01-initiative/initiative.md)

**What Happens:**
- PM uses the `initiative-definition` skill to define the problem: "40% of support tickets are about subscriptions"
- PM sets goals: "Enable subscriptions, reduce support to < 5%"
- PM lists affected components: Billing (new), Payments, Email, User Service

**Read:** [01-initiative/initiative.md](./01-initiative/initiative.md)

---

### Step 2: Risk Assessment & Workflow Selection

**Who:** Tech Lead or Architect

**Input:** Initiative

**Output:** [02-risk-assessment/risk-classification.md](./02-risk-assessment/risk-classification.md)

**What Happens:**
- Team answers 5 questions: Payment? PII? Multiple services? Rollback hard? Observability clear?
- Answers → Payment: YES, Multiple services: YES, Rollback: HARD → Risk = CRITICAL
- Risk = CRITICAL → Workflow = FULL (Analyst → Architect → Developer(s) → Verifier)
- CLI command: `agentic-agent specify start "Subscription Billing" --risk critical`

**Read:** [02-risk-assessment/risk-classification.md](./02-risk-assessment/risk-classification.md)

---

### Step 3: Analyst Phase

**Who:** Analyst Agent (or human analyst in collaborative mode)

**Input:** Initiative + Risk Classification

**Output:** [03-analyst/discovery.md](./03-analyst/discovery.md)

**Prompt:** [03-analyst/prompts.md](./03-analyst/prompts.md)

**What Happens:**
- Analyst calls Platform MCP to get payment policies
- Analyst interviews team one question at a time:
  - "What evidence shows this is really a problem?" → 43 support tickets in Jan
  - "Which services touch subscriptions?" → Billing, Payments, Email, User
  - "What ADRs are blocking?" → Card vault storage location (ADR-001)
  - "How will we observe success?" → Metrics, logs, traces
- Analyst produces discovery.md with problem statement, evidence, affected components, risk classification, blocking ADRs

**Read:** [03-analyst/discovery.md](./03-analyst/discovery.md)

**Exit Gate Check:**
- ✓ Problem has metric (43 tickets, 8 lost deals)
- ✓ Evidence is real data (not assumptions)
- ✓ Affected components named
- ✓ Risk classified with rationale
- ✓ Blocking ADRs identified (ADR-001)

---

### Step 4: Architect Phase

**Who:** Architect Agent

**Input:** Discovery.md

**Output:**
- [04-architect/feature-spec.md](./04-architect/feature-spec.md) — Platform-level WHAT
- [04-architect/component-spec-billing.md](./04-architect/component-spec-billing.md) — Billing service HOW
- [04-architect/component-spec-payments.md](./04-architect/component-spec-payments.md) — Payments service HOW
- [04-architect/ADR-001-card-vault.md](./04-architect/ADR-001-card-vault.md) — Blocking decision (card storage)

**Prompt:** [04-architect/prompts.md](./04-architect/prompts.md)

**What Happens:**
- Architect reads discovery.md
- Architect calls Platform MCP for payment standards, PCI requirements
- Architect calls Component MCPs for Billing, Payments, Email patterns & prior ADRs
- Architect designs feature-spec with 7 acceptance criteria in Given/When/Then format
- Architect designs 4 component-specs (one per service: Billing, Payments, Email, User)
- Architect creates ADR-001 to block implementation until card vault decision is made (self-hosted vs. Stripe)
- Architect runs self-check on all 5 gates

**Read These in Order:**
1. [04-architect/feature-spec.md](./04-architect/feature-spec.md) — Overall solution design
2. [04-architect/component-spec-billing.md](./04-architect/component-spec-billing.md) — Billing Service design
3. [04-architect/component-spec-payments.md](./04-architect/component-spec-payments.md) — Payments Service design
4. [04-architect/ADR-001-card-vault.md](./04-architect/ADR-001-card-vault.md) — Card storage decision (blocks everything)

**Gate Check Result:** Feature-spec and component-specs PASS all 5 gates, but implementation is BLOCKED until ADR-001 is resolved

---

### Step 5a: ADR Resolution (Parallel to Development)

**Who:** CTO + Security Team

**Input:** ADR-001 (card vault options)

**Output:** ADR-001 APPROVED with decision = Stripe Vault

**What Happens:**
- CTO and Security review 3 options:
  - A: Self-hosted vault (full control, high PCI burden)
  - B: Stripe vault (PCI handled by Stripe, lower time-to-market)
  - C: Hybrid (Stripe tokenize, ShopFlow store)
- Decision: B (Stripe Vault) chosen for Phase 1 (faster, lower risk)
- ADR marked APPROVED
- Blocking ADR removed: `blocked_by: []`

**Read:** [04-architect/ADR-001-card-vault.md](./04-architect/ADR-001-card-vault.md)

---

### Step 5b: Developer Phase (Parallel per Component)

**Who:** Developer Agents (2-5 running in parallel, one per component)

**Input:** Component-spec (e.g., component-spec-billing.md)

**Output:** impl-spec + tasks.yaml per component

**Prompt:** [05-developer/prompts.md](./05-developer/prompts.md)

**What Happens (Billing Developer):**
- Developer reads component-spec-billing.md
- Developer calls Component MCP for subscription state machine patterns, prior ADRs
- Developer checks `blocked_by` field — ADR-001 was blocking, but now resolved ✓
- Developer produces impl-spec-billing.md with:
  - Exact database schema (plans, subscriptions, subscription_events tables)
  - Code structure (services, repositories, models)
  - Edge cases (customer cancels mid-cycle, charge fails 3 times, merchant updates price mid-cycle)
  - Observability (structured JSON logging, metrics, tracing)
  - Rollout plan (feature flag strategy, canary, rollback procedure)
- Developer produces tasks.yaml with 10 development tasks (decomposed, ordered, dependencies)

**Read:**
- [05-developer/prompts.md](./05-developer/prompts.md) — How to invoke Developer
- [05-developer/impl-spec-billing.md](./05-developer/impl-spec-billing.md) — Billing implementation detail
- [05-developer/tasks.yaml](./05-developer/tasks.yaml) — Task decomposition

**Key Point:** Multiple developers work in parallel:
- Developer 1: Billing Service (payment/subscription logic)
- Developer 2: Payments Service (Stripe integration)
- Developer 3: Email Service (lifecycle notifications)
- Developer 4: User Service (subscription tracking)

All blocked until ADR-001 resolved. All unblocked simultaneously. All deliver in parallel.

---

### Step 6: Gate Checks (Before Merge)

**Who:** Architect or Engineer

**Input:** Feature-spec, component-specs, impl-specs

**Output:** [07-gate-checks/gate-check-output.md](./07-gate-checks/gate-check-output.md)

**What Happens:**
```bash
# Run gate check on feature spec
agentic-agent specify gate-check SBL-2025-Q1 --format text

# Run gate check on component specs
agentic-agent specify gate-check SBL-BILLING-001
agentic-agent specify gate-check SBL-PAYMENTS-001
...

# Result: All gates pass (or specific failures shown with remediation)
```

**Gate Checks:**
1. **Gate 1 — Context Completeness:** Every section has `Source:` line (MCP, Initiative, etc.)
2. **Gate 2 — Domain Validity:** No cross-domain DB access (all event-driven)
3. **Gate 3 — Integration Safety:** All contracts identified, consumers listed, compatibility plan if breaking
4. **Gate 4 — NFR Compliance:** Logging, metrics, tracing, PII, performance all specified
5. **Gate 5 — Ready to Implement:** No blocking ADRs, acceptance criteria testable, unambiguous

**Read:** [07-gate-checks/gate-check-output.md](./07-gate-checks/gate-check-output.md)

---

### Step 7: Verifier Phase

**Who:** Verifier Agent

**Input:** Merged code, all tests passing

**Output:** [06-verifier/verify.md](./06-verifier/verify.md)

**Prompt:** [06-verifier/prompts.md](./06-verifier/prompts.md)

**What Happens:**
- Verifier reads all acceptance criteria from feature-spec + component-specs
- Verifier gathers evidence for each AC:
  - **AC1 (Merchants can create plans)** → Test: `TestCreatePlan` PASSED, log shows plan created
  - **AC2 (Customers can start subscriptions)** → Test: `TestCreateSubscription` PASSED, subscription in DB with status ACTIVE
  - **AC3 (Auto-charges on schedule)** → Test: `TestScheduler` PASSED, 142 subscriptions processed, charges executed
  - **AC4 (Events emitted)** → Test: `TestEventPublishing` PASSED, 4 event types published
  - **AC5 (PII masked)** → Test: `TestPIIMasking` PASSED, no plaintext email/card in logs
  - **AC6 (Failures retry)** → Test: `TestFailureRetry` PASSED, failed charges suspended and retried
  - **AC7 (Rollback tested)** → Test: `TestRollback` PASSED, kill switch disables new subscriptions, existing continue
- Verifier checks code quality:
  - Tests: 127 tests, 100% passing ✓
  - Coverage: 87% (> 80%) ✓
  - Lint: 0 errors ✓
  - Build: Success ✓
- **BLOCKING CONDITION:** Payment/PII requires human sign-off (Security + Payments + Legal)
- Verifier marks: "READY FOR PRODUCTION (pending human approval)"

**Read:** [06-verifier/verify.md](./06-verifier/verify.md)

---

### Step 8: Completion & Knowledge Loop

**Who:** DevOps/SRE

**Input:** verify.md (all ACs verified, human approvals signed)

**Output:** Spec marked Done, knowledge updated in Component MCP

**What Happens:**
```bash
# Update spec graph
agentic-agent specify sync-graph

# Mark complete
agentic-agent validate SBL-2025-Q1

# Result: Spec status = Done, spec-graph updated
```

**Knowledge Loop:**
- Billing Service patterns learned → Component MCP updated for Phase 2
- Stripe integration patterns → Component MCP updated
- Edge cases handled (charge retry logic, idempotency) → Component MCP
- Observability setup (Datadog + OpenTelemetry) → Platform MCP updated

Next initiative benefits from this documented knowledge.

---

## Key Artifacts

| File | Purpose | Owner | Status |
|------|---------|-------|--------|
| [01-initiative/initiative.md](./01-initiative/initiative.md) | Problem statement, goals, success metrics | PM | ✓ Final |
| [02-risk-assessment/risk-classification.md](./02-risk-assessment/risk-classification.md) | Risk analysis, workflow selection | Tech Lead | ✓ Final |
| [03-analyst/discovery.md](./03-analyst/discovery.md) | Evidence-based problem analysis | Analyst Agent | ✓ Final |
| [04-architect/feature-spec.md](./04-architect/feature-spec.md) | Platform-level WHAT (7 ACs, NFRs) | Architect Agent | ✓ Approved |
| [04-architect/component-spec-billing.md](./04-architect/component-spec-billing.md) | Billing Service design | Architect Agent | ✓ Approved |
| [04-architect/component-spec-payments.md](./04-architect/component-spec-payments.md) | Payments Service design | Architect Agent | ✓ Approved |
| [04-architect/ADR-001-card-vault.md](./04-architect/ADR-001-card-vault.md) | Card vault storage decision | CTO + Security | ✓ Approved (Stripe chosen) |
| [05-developer/impl-spec-billing.md](./05-developer/impl-spec-billing.md) | Billing implementation detail | Developer Agent | ✓ Final |
| [05-developer/tasks.yaml](./05-developer/tasks.yaml) | 10 development tasks, dependencies, effort | Developer Agent | ✓ Ready |
| [06-verifier/verify.md](./06-verifier/verify.md) | Test evidence for all ACs | Verifier Agent | ✓ Final (human approval pending) |
| [07-gate-checks/gate-check-output.md](./07-gate-checks/gate-check-output.md) | Gate 1-5 validation results | Architect/Engineer | ✓ All passing |

---

## For Humans vs. AI Agents

### Humans Should Read:

1. **Introduction** (this README) — understand the flow
2. [01-initiative/initiative.md](./01-initiative/initiative.md) — problem statement
3. [02-risk-assessment/risk-classification.md](./02-risk-assessment/risk-classification.md) — understand why Full workflow
4. [04-architect/feature-spec.md](./04-architect/feature-spec.md) — understand the overall solution
5. [06-verifier/verify.md](./06-verifier/verify.md) — understand how we verified it works

### AI Agents Should Read:

Start with [AGENT-GUIDE.md](./AGENT-GUIDE.md) for concise role-based instructions.

Or follow the timeline above, reading prompts.md + artifact files for each phase.

---

## Try It Yourself

### Option 1: Replay This Example in Your Project

```bash
# 1. Copy this example to your project
cp -r examples/sdd-ecommerce-payments/ my-project/

# 2. Update initiative.md with your problem
# 3. Run risk assessment
agentic-agent specify start "Your Feature" --risk critical

# 4. Run analyst agent with your team
# (Use 03-analyst/prompts.md as template)

# 5. Continue through each phase...
```

### Option 2: Customize for Your Feature

Use the artifacts in this example as templates:

- **Feature-spec template:** [04-architect/feature-spec.md](./04-architect/feature-spec.md)
- **Component-spec template:** [04-architect/component-spec-billing.md](./04-architect/component-spec-billing.md)
- **Implementation-spec template:** [05-developer/impl-spec-billing.md](./05-developer/impl-spec-billing.md)
- **Verification template:** [06-verifier/verify.md](./06-verifier/verify.md)

### Option 3: Study for Understanding

Read the spec files to understand:
- How to write Given/When/Then acceptance criteria (see feature-spec AC1-AC7)
- How to design edge cases (see impl-spec-billing Edge Cases table)
- How to document observability (see NFRs section in feature-spec)
- How to write blocking ADRs (see ADR-001-card-vault.md)
- How to verify with evidence (see verify.md)

---

## Summary: SDD in 8 Phases

| Phase | Agent/Role | Input | Output | Duration |
|-------|-----------|-------|--------|----------|
| 1 | Initiative | Problem | [initiative.md](./01-initiative/initiative.md) | 1-2d |
| 2 | Tech Lead | Initiative | [risk-classification.md](./02-risk-assessment/risk-classification.md) | 1d |
| 3 | Analyst | Discovery interview | [discovery.md](./03-analyst/discovery.md) | 3-4d |
| 4 | Architect | Discovery | [feature-spec.md](./04-architect/feature-spec.md) + [component-specs](./04-architect/) | 5d |
| 5 | CTO + Security | ADRs | [ADR-001 APPROVED](./04-architect/ADR-001-card-vault.md) | 3-5d |
| 6 | Developers (parallel) | Component-specs | [impl-specs](./05-developer/) + [tasks.yaml](./05-developer/tasks.yaml) | 8-10d |
| 7 | Engineers | Code + tests | [gate checks PASS](./07-gate-checks/gate-check-output.md) | 1d |
| 8 | Verifier | Merged code | [verify.md](./06-verifier/verify.md) ✓ APPROVED | 2d |

**Total:** 3-4 weeks for critical payment feature, with evidence at every step.

---

## Questions?

- **What is SDD?** → Read [docs/sdd-mcp/TARGET-OPERATING-MODEL.md](../../docs/sdd-mcp/operation%20model/TARGET-OPERATING-MODEL.md)
- **How do I run the Analyst phase?** → See [03-analyst/prompts.md](./03-analyst/prompts.md)
- **What's an ADR?** → See [04-architect/ADR-001-card-vault.md](./04-architect/ADR-001-card-vault.md)
- **How do I write acceptance criteria?** → See feature-spec ACs
- **How do I verify?** → See [06-verifier/verify.md](./06-verifier/verify.md)

---

**Made with SDD. Start your next feature the right way.**
