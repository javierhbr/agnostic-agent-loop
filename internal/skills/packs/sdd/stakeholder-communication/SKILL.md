---
name: stakeholder-communication
description: "Guide for PMs to communicate SDD workflow status to non-technical stakeholders in plain language."
---

# Stakeholder Communication (Product Managers)

## Purpose

SDD uses technical terms (specs, gates, ADRs). Your job is to translate them into business language
when talking to executives, customers, and other non-engineering stakeholders.

## Status Translation Table

| Technical Status | Business Translation | What It Means |
|---|---|---|
| Draft | "Design in progress" | Engineers are figuring out how to build it |
| Approved | "Ready for development" | Design is locked, developers can start |
| Implementing | "In development" | Developers are writing code |
| Blocked (by ADR) | "Waiting on a decision" | Technical decision needed before we can proceed |
| Done | "Complete and verified" | Built, tested, and ready for customers |

## ADR Blocking Explanation

When a change is blocked by an ADR, here's how to explain it:

**Technical phrase:** "Spec SPEC-001 is blocked by ADR-219: Idempotency Strategy."

**For stakeholders:** "We're waiting on a technical decision about how to handle retries.
Once [Owner Name] approves the decision (expected by [Date]), we can proceed with implementation."

**Status update template:**
```
Status: Waiting on Decision

Initiative: [Name]
Blocker: ADR-219 (Idempotency Strategy)
Owner: [Name] (expected decision by [Date])

What this means: Our team has identified that we need to make a technical decision
before we can safely build this feature. Once [Owner] approves, we'll unlock implementation.
```

## Progress Report Template

Send to stakeholders every 1-2 weeks:

```markdown
# Status Update: [Initiative Name]

## Progress

- **Phase:** [Discovery / Design / Development / Verification]
- **Completion:** [X% - based on agents and gates passed]
- **Timeline:** [On track / At risk / Delayed]
- **Expected completion:** [Date]

## Current Status

- **What's done:** [e.g., "Design approved by Architect, ready for development"]
- **What's in progress:** [e.g., "4 developers implementing across payment, cart, and order services"]
- **What's blocked:** [If anything, explain simply]

## Key Metrics (if applicable)

- [e.g., Cart abandonment improvement target: 35% → 15%]
- [e.g., Payment processing: target p95 < 2s (currently 5s)]

## Risks & Mitigations

- **Risk:** [If any risk, e.g., "Payment provider integration complexity"]
  - **Mitigation:** [How we handle it, e.g., "Using feature flag for safe rollback"]

## Next Steps

- [e.g., "Finish development by Friday"]
- [e.g., "Start verification testing next week"]
- [e.g., "Waiting for decision on ADR-219"]

---
Questions? Contact [PM Name] or your engineering lead.
```

## Escalation Triggers

Escalate to leadership if ANY of these happen:

1. **Gate failure:** "Design did not pass inspection - needs rework"
2. **Unresolved ADR past 5 days:** "Technical decision still pending, blocking development"
3. **Critical risk with no rollback:** "We cannot safely undo this if it fails - rethinking approach"
4. **Timeline slip > 1 week:** "Revised estimate: delivery date is now [New Date]"
5. **Scope creep:** "New requirement emerged - impacts timeline/risk"

## Plain Language Explanations

| Technical | Business Speak |
|---|---|
| "Gate 4 (NFR Compliance) FAILED" | "Code doesn't have enough monitoring - can't see if it's working well in production" |
| "blocked_by: [ADR-001]" | "Waiting on a technical decision" |
| "Spec graph inconsistency" | "Our records show conflicting information - auditing the system" |
| "Contract change: breaking" | "This changes how services talk to each other - requires coordination" |
| "GWT acceptance criteria" | "Clear test scenarios: given this situation, when we do this, we expect this result" |

## Example Email to Exec

```
Subject: [Project Name] - Status Update

Hi [Exec],

[Project] is progressing on schedule. Here's where we stand:

**Current Phase:** Implementation (Developers writing code across 3 teams)

**On Track:** Yes - 70% complete, target delivery [Date]

**Key Metrics:**
- [Metric 1] improving from [Current] to [Target]
- [Metric 2] improving from [Current] to [Target]

**No Blockers:** All technical decisions are approved. Implementation proceeding as planned.

**Safety:** We have feature flags in place to safely rollback if needed.

Next milestone: Complete all code reviews by [Date], then verification testing.

Questions? Happy to discuss.

—[PM Name]
```
