# Architect Phase: Prompts

## How to Invoke the Architect Agent

**Context:** You have Analyst's `discovery.md`. Now invoke Architect to design the platform feature-spec and per-component architecture.

### Prompt for Architect Agent

```
You are running the Architect phase of the SDD Full workflow. Your job is to read
discovery.md, call Platform MCP and Component MCPs, then produce:
1. feature-spec.md (platform-level WHAT and WHY)
2. component-spec.md × N (per-service HOW)
3. ADR documents (blocking decisions)

INITIATIVE: SBL-2025-Q1 - Subscription Billing Platform
DISCOVERY: See 03-analyst/discovery.md

YOUR JOB:
1. Read discovery.md thoroughly
2. Call Platform MCP: get_context_pack("Subscription Billing", risk="critical")
3. Call Component MCP for each service:
   - Payments Service: get_contracts(), get_invariants(), get_decisions()
   - User Service: get_contracts(), get_invariants()
   - Email Service: get_contracts()
4. Produce feature-spec.md with 5+ acceptance criteria in Given/When/Then format
5. Produce component-spec.md for each affected service (Billing, Payments, Email, User)
6. Produce ADR-001: Decide on card storage strategy (self-hosted vault vs. Stripe)
7. Self-check all 5 gates before handoff to Developer

REQUIRED OUTPUTS:
- 04-architect/feature-spec.md (2000-3000 words)
- 04-architect/component-spec-billing.md
- 04-architect/component-spec-payments.md
- 04-architect/component-spec-email.md
- 04-architect/component-spec-user.md (minimal)
- 04-architect/ADR-001-card-vault.md

REQUIRED SECTIONS (EACH WITH SOURCE: LINE):
- Metadata (implements, context_pack, blocked_by, status)
- Problem Statement
- Goals / Non-Goals
- User Experience / API Examples
- Domain Responsibilities
- Cross-Domain Interactions
- NFRs (from Platform MCP)
- Feature Flag & Rollback Strategy
- Acceptance Criteria (min 3 in Given/When/Then)
- Blocking Decisions
- 5 Gates validation checklist

RULES:
- EVERY section has a Source: line citing the MCP call or discovery.md
- Feature-spec focuses on WHAT (user experience, business goals)
- Component-specs focus on HOW (technical contracts, data models, internal logic)
- If a section requires a decision, create an ADR and block with blocked_by: [ADR-001, ...]
- Run gate-check mentally before handing off

START NOW. Read discovery.md and begin the feature-spec.
```

---

## Example Architect Artifacts (Shown Below)

See the actual spec files in this directory:
- `feature-spec.md` — Platform-level WHAT
- `component-spec-billing.md` — Billing Service HOW
- `component-spec-payments.md` — Payments Service HOW
- `component-spec-email.md` — Email Service HOW
- `component-spec-user.md` — User Service changes
- `ADR-001-card-vault.md` — Blocking decision on card storage

These were produced by the Architect agent reading discovery.md and calling MCPs.
