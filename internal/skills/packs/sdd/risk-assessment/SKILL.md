---
name: risk-assessment
description: "Risk assessment interview guide for PMs to classify initiatives as Quick/Standard/Full before handing to engineers."
---

# Risk Assessment (Product Managers)

## Purpose

Before any engineering work begins, classify your initiative's risk level using this interview guide.
Risk level determines the workflow type: Quick (low) → Standard (medium) → Full (high/critical).

## Interview Flow

Ask ONE question at a time. Wait for answer before next question. Document responses.

### Question 1: Contract Impact
**"Does this change any cross-service contract (API, event schema, or database interface)?"**

- YES → Risk increases
- NO → Continue

### Question 2: Data Sensitivity
**"Does this change touch payment, authentication, PII (email/phone/address), or health data?"**

- YES → Risk increases significantly
- NO → Continue

### Question 3: Team Scope
**"Will more than one component team work on this simultaneously?"**

- YES → Risk increases (coordination complexity)
- NO → Continue

### Question 4: Rollback Strategy
**"If this breaks in production, how quickly can you rollback? Do you have a rollback plan?"**

- NO ROLLBACK PLAN → Risk increases (can't undo mistakes)
- HAS ROLLBACK PLAN → Risk decreases
- FEATURE FLAG WITH DEFAULT OFF → Risk decreases significantly

### Question 5: ADR Dependency
**"Are there any existing ADRs (Architecture Decision Records) that govern this area, or do we need to make new design decisions?"**

- OPEN/UNCERTAIN ADRs → Risk increases (blocked on decisions)
- CLEAR ADRs → Risk decreases
- NEW DECISIONS NEEDED BUT SIMPLE → Risk increases (adds time)

## Risk Scoring

Count "YES" answers to high-risk questions (1, 2) and "UNCERTAIN/OPEN" answers to ADR question:

| Score | Risk Level | Workflow | Timeline |
|-------|-----------|----------|----------|
| 0-1 | Low | Quick | 1-2 days |
| 2-3 | Medium | Standard | 3-5 days |
| 4+ | High/Critical | Full | 5-10 days |

## Documenting the Classification

Update your initiative.md with a "Risk Assessment" section:

```markdown
## Risk Assessment

- Contract change: [YES / NO]
  - If yes: [which contracts / APIs]
- Data sensitivity: [YES / NO]
  - If yes: [payment / auth / PII / health]
- Team scope: [single service / 2-3 services / 4+ services]
- Rollback strategy: [feature flag / database rollback / revert deploy / NO ROLLBACK PLAN]
- ADR dependencies: [none / ADR-123 (Idempotency) / ADR-456 (Session TTL)]

**Classification: [LOW / MEDIUM / HIGH / CRITICAL]**

**Recommended Workflow: [QUICK / STANDARD / FULL]**

**Rationale:**
[Explain why this risk level, referencing the answers above]
```

## After Assessment: Hand Off to Engineering

```bash
# Determine risk level from assessment
agentic-agent sdd start "<initiative-name>" --risk [low|medium|high|critical]

# This activates the appropriate workflow (Quick/Standard/Full)
# The first agent in the sequence is notified to begin work
```

## Example

```markdown
## Risk Assessment

- Contract change: YES (CartUpdated event schema v2 → v3)
- Data sensitivity: NO (cart contents, not payment/PII)
- Team scope: 2-3 services (cart-service, order-service, analytics)
- Rollback strategy: Feature flag (CartV3Enabled, defaults to OFF)
- ADR dependencies: ADR-219 (Idempotency Strategy) — RESOLVED

**Classification: MEDIUM**

**Recommended Workflow: STANDARD**

**Rationale:**
Contract change (CartUpdated event) increases risk, but we have feature flag rollback
and only 2-3 services involved. No payment/PII involved. ADR-219 is already approved.
This is a good fit for Standard workflow (Architect + parallel developers + Verifier).
```
