
# SDD Methodology: Three Usage Levels

The SDD (Spec-Driven Development) methodology supports **three independent levels of usage** depending on your needs, risk level, and team structure.

---

## Overview

| Level | Scope | Use Case | Risk | Phases | Timeline |
|-------|-------|----------|------|--------|----------|
| **OpenSpec** | Single service/team | Component changes, bug fixes | Low | Specify + Deliver only | 1-3 weeks |
| **Platform** | Multi-service coordination | Architecture decisions, API contracts | Medium | Platform + Assess + Specify + Plan + Deliver | 1-2 weeks |
| **SDD Full** | Critical multi-service | Payment systems, auth, PII, high-risk | Critical | All 5 phases (Platform, Assess, Specify, Plan, Deliver) | 3-4 weeks |

**Core boundary rule:** Component repos use OpenSpec ONLY. Platform-side tooling (BMAD, Speckit) stays in the platform repo.

---

## Level 1: OpenSpec (Component Level)

**Use when:** Your change is isolated to one service/team and doesn't require multi-team coordination.

### What It Is

OpenSpec is a lightweight specification + task decomposition system for individual teams. You write a proposal and decompose it into tasks, then track completion.

**Core boundary rule:** Component repos use OpenSpec ONLY. They do not use BMAD, Speckit, or platform-level tooling.

### Commands

```bash
# Create a change proposal + decompose into tasks
agentic-agent openspec init "Auth Feature" --from requirements.md

# This creates:
# - proposal.md (what you're building, requirements)
# - tasks.md (how you're breaking down the work)
# - metadata.yaml (contracts, blocked_by, status)

# Check gates (basic validation)
agentic-agent openspec check <change-id>

# List all changes
agentic-agent openspec list

# Show details
agentic-agent openspec show <change-id>

# Track progress
agentic-agent openspec status <change-id>

# Mark complete
agentic-agent openspec complete <change-id>

# Archive after completion
agentic-agent openspec archive <change-id>
```

### Gate Checks (Basic)

OpenSpec enforces basic gates:

- Metadata exists (change ID, status, blocked_by)
- No unresolved ADRs blocking the change
- Contracts listed if referencing other services

### Directory Structure

```
.agentic/specs/<change-id>/
├── proposal.md          # What you're building
├── tasks.md             # How you're breaking it down
└── metadata.yaml        # ID, status, blocked_by, contracts
```

### Example: Bug Fix (Low Risk)

```bash
# Team working on payments-service fixes a race condition
cd payments-service/

agentic-agent openspec init "Fix card expiry race condition" \
  --from bug-report.md \
  --type bug

# Team decomposes into tasks
# > Adds 3 tasks to tasks.md with acceptance criteria

# Team implements
agentic-agent task claim <task-id>
# ... code changes ...
agentic-agent task complete <task-id>

# Team marks change done
agentic-agent openspec complete card-expiry-fix
agentic-agent openspec archive card-expiry-fix
```

### When to Use OpenSpec Only

- Single service/component
- No other teams depend on your changes
- Low risk (bug fix, internal refactor, small feature)
- Timeline: 1-3 weeks
- Team-level autonomy (no platform governance needed)

---

## Level 2: Platform (Architecture/Governance Level)

**Use when:** Platform team defines specs for multiple component teams, or you need to enforce architecture decisions across services.

### What It Is

Platform level provides governance, initiative tracking, and fan-out task distribution using BMAD + OpenSpec + Speckit. Platform team designs specs, then distributes them to component teams who pick up OpenSpec changes.

**Note:** Platform-side tooling (BMAD, Speckit) stays in the platform repo. Component repos receive fan-out tasks and use OpenSpec ONLY (core boundary rule).

### Commands

```bash
# Initialize platform repo structure
agentic-agent platform init

# Creates:
# - constitution/ (governance policies)
# - initiatives/ (initiative tracking)
# - platform-specs/ (platform-level specs)
# - contracts/ (contract definitions)
# - adr/ (architecture decision records)

# Create a platform feature spec with fan-out
agentic-agent platform add-feature \
  --name "Subscription Billing" \
  --initiative "ECO-123" \
  --context-pack "cp-v2.1" \
  --fanout "billing-service,payments-service,email-service"

# This creates:
# - platform-specs/subscription-billing/
#   ├── spec.md (platform-level WHAT)
#   ├── metadata.yaml
#   └── fanout.yaml (fan-out to 3 component teams)

# Update platform spec status
agentic-agent platform change-feature \
  --id "subscription-billing" \
  --field "status" \
  --value "Implementing"

# Update platform spec context pack
agentic-agent platform change-feature \
  --id "subscription-billing" \
  --field "context_pack" \
  --value "cp-v2.2"

# Track initiative priority/status
agentic-agent platform change-priority \
  --initiative "ECO-123" \
  --priority "Approved" \
  --status "Ready for team pickup"
```

### Gate Checks (All 5 Gates)

Platform specs enforce all 5 SDD gates:
1. **Gate 1 — Context Completeness:** Every section has Source: line (MCP citations)
2. **Gate 2 — Domain Validity:** No cross-domain DB access, boundaries respected
3. **Gate 3 — Integration Safety:** All contracts identified, consumers listed
4. **Gate 4 — NFR Compliance:** Logging, metrics, tracing, PII, performance specified
5. **Gate 5 — Ready to Implement:** No blocking ADRs, ACs testable, unambiguous

### Directory Structure

```
Platform Repo:
├── constitution/ (governance policies)
├── initiatives/ (initiative tracking)
├── platform-specs/
│   └── <feature-id>/
│       ├── spec.md (platform WHAT)
│       ├── metadata.yaml
│       └── fanout.yaml (fan-out tasks)
├── contracts/ (contract definitions)
├── adr/ (architecture decision records)
└── spec-graph.json (traceability)

Component Repo (receives fan-out):
└── .agentic/specs/<component-spec-id>/
    ├── proposal.md
    ├── tasks.md
    └── metadata.yaml
```

### Example: Multi-Service API (Medium Risk)

```bash
# Platform team designs the feature
cd platform-repo/

agentic-agent platform add-feature \
  --name "Webhooks for Subscriptions" \
  --initiative "ECO-124" \
  --context-pack "cp-v2.1" \
  --fanout "billing-service,notifications-service,auth-service"

# Platform team writes spec
vim platform-specs/webhooks/spec.md
# - Problem: "Customers need to be notified of subscription events"
# - Contracts: "subscription.created", "subscription.renewed", "subscription.cancelled"
# - NFRs: "Webhook delivery within 60 seconds", "retry up to 3 times"

# Each component team picks up their fan-out task
cd ../billing-service/

agentic-agent openspec init "Implement webhook events" \
  --from ../platform-specs/webhooks/spec.md \
  --implements "webhooks"

# Team implements their part (emitting events)
# Deliverable: webhook event emission in billing service

# All teams complete
agentic-agent openspec complete webhooks-events

# Platform tracks completion
agentic-agent platform change-priority \
  --initiative "ECO-124" \
  --priority "Done" \
  --status "All 3 teams completed webhook implementation"
```

### When to Use Platform Level

- Multiple services affected
- Cross-team coordination needed
- Architecture decisions matter (contracts, APIs, patterns)
- Platform team owns governance
- Medium risk (shared patterns, breaking changes)
- Timeline: 1-2 weeks
- Need to track initiative progress across teams

---

## Level 3: SDD Full Workflow (Critical Features)

**Use when:** Risk is high/critical and you need full discovery → verification rigor.

### What It Is

The complete SDD methodology with five phases and a two-iteration adoption model:

1. **Platform** — Define governance, constitution, and platform-level specs (BMAD + Speckit)
2. **Assess** — Discover the problem with team evidence (Analyst role)
3. **Specify** — Design delta specs with gates (Architect role, produces proposal.md + delta specs)
4. **Plan** — Create design.md and tasks.md (Architect role)
5. **Deliver** — Implement specs in parallel, then verify all ACs with evidence (Developer + Verifier roles)

**Two-iteration adoption:** Teams start with Specify + Deliver (iteration 1), then add Platform + Assess + Plan (iteration 2) as they mature.

### Commands

```bash
# Start an SDD initiative
agentic-agent specify start "Subscription Billing" --risk critical

# This creates:
# - .agentic/initiatives/<id>/ (workflow state)
# - Sets workflow = Full (for critical risk)
# - Next agent: Analyst

# Show workflow progress
agentic-agent specify workflow show <initiative-id>

# Install agent skills
agentic-agent specify agents install

# Create ADR for blocking decisions
agentic-agent specify adr create --title "Where to store card tokens?"

# Resolve ADR when decision made
agentic-agent specify adr resolve ADR-001

# Sync spec graph
agentic-agent specify sync-graph

# Complete initiative
agentic-agent openspec complete <change-id>
```

### SDD Phases

**Phase 1: Platform** (skip for component-only changes)

- Define governance and constitution
- Create platform-level specs using BMAD + Speckit
- Set up platform-ref.yaml for component alignment
- Produce: constitution, platform-specs, contracts

**Phase 2: Assess** (lightweight for medium risk, full for critical)
- Interview team one question at a time
- Extract evidence (real data, not assumptions)
- Classify risk with rationale
- Identify blocking ADRs
- Produce: `discovery.md`

**Phase 2.5: ADR Resolution**
- CTO/Security team reviews ADR options
- Make decision (e.g., Stripe vault chosen)
- Mark ADR approved
- Unblock implementation

**Phase 3: Specify** (Architect role)

- Read discovery / requirements
- Design proposal.md (what you are building)
- Create delta specs (per service changes)
- Create ADRs for blocking decisions
- Self-check all 5 gates
- Produce: proposal.md + delta specs + ADRs

**Phase 4: Plan** (Architect role)

- Read proposal and delta specs
- Create design.md (technical design)
- Decompose into tasks.md (8-10 tasks)
- Self-check all 5 gates
- Produce: design.md + tasks.md

**Phase 5: Deliver** (Developer + Verifier roles)

- Developer reads delta specs and tasks.md
- Implement in parallel per component
- Verifier reads all ACs from specs
- For each AC, gather evidence:
  - Test name + output
  - Log showing feature works
  - Metric/trace evidence
- Check code quality (tests, lint, coverage)
- Write verify.md with evidence
- Blocking condition: Any unverified AC --> BLOCK
- Produce: implementation + verify.md with sign-off

### Gate Checks (All 5 Gates + Evidence)

SDD Full enforces:
1. Gate 1: Context Completeness (all sections sourced)
2. Gate 2: Domain Validity (no DB crossing)
3. Gate 3: Integration Safety (contracts + consumers)
4. Gate 4: NFR Compliance (logging, metrics, tracing, PII, perf)
5. Gate 5: Ready to Implement (no blocking ADRs, ACs testable)

**Plus:** Verification with observable evidence

### Workflow Paths by Risk

| Risk | Workflow | Phases | Timeline |
|------|----------|--------|----------|
| Low | Quick | Specify + Deliver | 1-2 days |
| Medium | Standard | Assess (lightweight) + Specify + Plan + Deliver | 3-5 days |
| High/Critical | **Full** | **Platform + Assess + Specify + Plan + Deliver** | **5-10 days** |

### Example: Critical Payment Feature

See [`examples/sdd-ecommerce-payments/`](../../examples/sdd-ecommerce-payments/) for a complete, worked example with all artifacts.

**Quick summary:**

```bash
# Start SDD
agentic-agent specify start "Subscription Billing" --risk critical

# Phase 1: Platform sets governance (if needed)
# Produces: constitution, platform-ref.yaml

# Phase 2: Assess - Analyst interviews team
# Produces: discovery.md (evidence of 43 support tickets, 8 lost deals)

# Phase 3: Specify - Architect designs
# Produces: proposal.md, delta specs per service
# Creates ADR-001: "Where to store card tokens?" (self-hosted vs Stripe)

# Phase 3.5: CTO/Security decide
# ADR-001 approved: "Use Stripe Vault"

# Phase 4: Plan - Architect creates design + tasks
# Produces: design.md + tasks.md per component
# Billing: 10 tasks, Payments: 8 tasks

# Phase 5: Deliver - Developers implement, Verifier checks
# Developers implement in parallel
# Verifier checks all ACs with evidence
# Produces: verify.md
# All ACs verified: tests pass, logs show behavior, metrics collected

# Complete
agentic-agent openspec complete SBL-2025-Q1
```

### When to Use SDD Full

- Risk = High or Critical (payment, auth, PII, regulatory)
- Multiple services affected (5+)
- Blocking decisions needed (ADRs)
- Team needs evidence-based verification
- Data consistency is critical
- Rollback is complex
- Timeline: 3-4 weeks (justified by risk)
- Full audit trail needed

---

## Comparison: Which Level to Use?

### Question 1: Is your change isolated to one service?

**YES** → Use **OpenSpec** (Level 1)
- Team owns the spec and tasks
- Low ceremony, fast execution
- Example: Fix a bug, add internal feature

**NO** → Continue to Question 2

### Question 2: Does the change require architecture decisions?

**YES** → Use **Platform** (Level 2) or **SDD Full** (Level 3)
- Platform team designs contracts
- Component teams implement in parallel
- Example: New API endpoint, shared pattern

**NO** → Continue to Question 3

### Question 3: What's the risk level?

**Low/Medium Risk** → Use **Platform** (Level 2)
- Analyst phase not needed (requirements clear)
- Architect designs specs
- Developers implement
- Timeline: 1-2 weeks

**High/Critical Risk** → Use **SDD Full** (Level 3)
- Analyst phase discovers evidence
- Architect designs from evidence
- Developers implement
- Verifier checks all ACs with evidence
- Timeline: 3-4 weeks

---

## Decision Tree

```
Start: What's your change?

├─ Single service, low risk?
│  └─ YES → OpenSpec (Level 1)
│
└─ Multi-service or architecture decision?
   ├─ YES, Low/Medium risk?
   │  └─ YES → Platform (Level 2)
   │
   └─ High/Critical risk?
      └─ YES → SDD Full (Level 3)
```

---

## Integration: How They Work Together

```
SDD Full Workflow (When Risk = Critical) — 5 Phases
│
├─ Phase 1: Platform — governance, constitution, platform-ref.yaml
├─ Phase 2: Assess — discovers evidence (discovery.md)
├─ Phase 3: Specify — designs from evidence
│  │
│  ├─ Creates proposal.md (what you are building)
│  └─ Creates delta specs (per service changes)
│     │
│     ├─ Maps to Platform level:
│     │  └─ platform-specs/<feature-id>/spec.md + fanout.yaml
│     │
│     └─ Maps to OpenSpec level (core boundary rule):
│        └─ Each team's .agentic/specs/<component>/
│           ├─ proposal.md
│           ├─ tasks.md
│           └─ implementation
│
├─ Phase 4: Plan — design.md + tasks.md per component
├─ Phase 5: Deliver — Developer(s) implement, Verifier checks with evidence
└─ Complete
```

---

## Commands by Level

### OpenSpec Commands

```bash
agentic-agent openspec init <name> --from <file>
agentic-agent openspec list
agentic-agent openspec show <id>
agentic-agent openspec status <id>
agentic-agent openspec check <id>
agentic-agent openspec complete <id>
agentic-agent openspec archive <id>
```

### Platform Commands

```bash
agentic-agent platform init
agentic-agent platform add-feature --name "..." --initiative "..." --context-pack "..." --fanout "..."
agentic-agent platform change-feature --id "..." --field "status" --value "..."
agentic-agent platform change-priority --initiative "..." --priority "..."
```

### SDD Commands

```bash
agentic-agent specify start <name> --risk low|medium|high|critical
agentic-agent specify workflow show <id>
agentic-agent specify agents install
agentic-agent specify gate-check <spec-id>
agentic-agent specify adr create --title "..."
agentic-agent specify adr resolve <adr-id>
agentic-agent specify sync-graph
```

---

## Example: Refund Feature (3-Service Change)

### Scenario
Add credit card refunds to payments system. Touches: payments, accounting, customer-notifications.

### Risk Assessment
- Payment-related: YES
- PII involved: NO (card already tokenized)
- Multiple services: YES (3)
- Rollback hard: YES (if charge succeeded but refund record missing)
- **Risk = Medium** → Use **Platform** level

### Workflow

```bash
# Step 1: Platform team designs
agentic-agent platform add-feature \
  --name "Card Refunds" \
  --initiative "PAY-2025" \
  --context-pack "cp-v2.1" \
  --fanout "payments-service,accounting-service,notifications-service"

# Step 2: Platform team writes spec.md
# - Problem: "Merchants need ability to refund customers"
# - Contracts: "refund.initiated", "refund.completed" events
# - NFRs: "Refund within 5 seconds", "Webhook to accounting within 60s"

# Step 3: Each service team picks up their part
cd payments-service/
agentic-agent openspec init "Process refunds" \
  --from ../platform-specs/card-refunds/spec.md \
  --implements "card-refunds" \
  --type "feature"

# Team writes proposal + tasks
vim .agentic/specs/refund-processor/proposal.md
# - Implement POST /refunds/:id endpoint
# - Call Stripe refund API
# - Emit refund.initiated + refund.completed events

# Step 4: Teams implement independently
agentic-agent task claim <task-1>
# ... code ...
agentic-agent task complete <task-1>

# Step 5: Teams mark done
agentic-agent openspec complete refund-processor

# Step 6: Platform tracks completion
agentic-agent platform change-priority \
  --initiative "PAY-2025" \
  --priority "Done"
```

---

## Best Practices

### OpenSpec Level (Specify + Deliver)

- Write proposal that your team understands
- Decompose into 3-10 concrete tasks
- Use blocked_by if waiting on other teams
- Mark complete only when tested

### Platform Level (Platform + Assess + Specify + Plan + Deliver)

- Write platform spec with clear contracts
- Use fanout to distribute to teams
- Run all 5 gate checks before fan-out
- Track initiative priority/status
- Update component teams with changes

### SDD Full Level (All 5 Phases)

- Use for critical features only
- Get Assess phase buy-in (evidence matters)
- Specify phase must produce proposal.md + delta specs
- Resolve ADRs quickly (they block development)
- Plan phase produces design.md + tasks.md
- Deliver phase: Developer runs in parallel, Verifier checks EVERY AC with evidence

---

## See Also

- [TARGET-OPERATING-MODEL.md](./operation%20model/TARGET-OPERATING-MODEL.md) — Complete SDD methodology
- [OPERATING-MODEL.md](./operation%20model/OPERATING-MODEL.md) — Operational details
- [examples/sdd-ecommerce-payments/](../../examples/sdd-ecommerce-payments/) — Full worked example
  - [README.md](../../examples/sdd-ecommerce-payments/README.md) — Start here for complete walkthrough
  - [AGENT-GUIDE.md](../../examples/sdd-ecommerce-payments/AGENT-GUIDE.md) — Role-based instructions

---

## Questions?

**Q: Can I use OpenSpec and Platform together?**
A: Yes. OpenSpec for component work, Platform for architecture governance. They integrate naturally through spec references.

**Q: Do I have to use all three levels?**
A: No. Use the minimum needed for your risk level:
- Low risk: OpenSpec only
- Medium risk: Platform only
- Critical risk: SDD Full (includes Platform + OpenSpec)

**Q: What if I start with OpenSpec and later realize I need Platform/SDD?**
A: You can upgrade anytime. Reference your existing OpenSpec files in Platform specs. Continue to next phase (Architect, then Developer/Verifier).

**Q: Can a small team use SDD Full?**
A: Yes. The five phases (Platform, Assess, Specify, Plan, Deliver) can be roles for 2-3 people. Team Lead handles Assess and Deliver, Product handles Specify, Architect handles Platform and Plan. It takes longer (3-4 weeks) but the evidence trail is valuable for critical features.
