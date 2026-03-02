---
name: platform-constitution
description: "Guide for Platform Architects to author and maintain the Constitution — the governance policy document for the entire platform."
---

# skill:platform-constitution

## Does exactly this

Helps Platform Architects author and maintain the Constitution — the foundational governance document that defines platform-wide policies, enforcement rules, and non-negotiables.

---

## When to use

- Starting a new platform and need to define governance
- Platform policies need to be updated or clarified
- A gate check failed because a policy was missing
- Onboarding new teams and need to communicate non-negotiables

---

## What Is the Constitution?

The Constitution defines:
- How services interact (contracts, versioning)
- What developers must do (observability, security, testing)
- What we care about most (performance, reliability, UX)
- What's non-negotiable (PII handling, compliance, data retention)

**It is NOT:** Style guide, best practices list, team-specific rules
**It IS:** Platform-wide policies, enforcement rules in gates, source of truth for all specs

---

## Constitution Structure (6 Sections)

**1. UX Rules** — Accessibility, design system, mobile-first. Gate 4 enforces.

**2. Security & PII Policy** — Encryption, masking, auth, session TTL, cross-service rules. Gate 4 enforces.

**3. Observability Standards** — Log format, metrics, tracing, alerts. Gate 4 enforces.

**4. Performance Baselines** — p95 latency, throughput, caching. Gate 4 enforces.

**5. Domain Governance** — Service ownership, invariants, cross-domain communication. Gate 2 enforces.

**6. Contract Versioning** — Semver, dual-publish window, deprecation SLA. Gate 3 enforces.

---

## How Gates Use the Constitution

| Gate | Enforces |
|------|----------|
| Gate 1 (Context) | Constitution version declared in every spec |
| Gate 2 (Domain) | Data ownership per Constitution; no direct cross-service DB access |
| Gate 3 (Integration) | Contract versioning and dual-publish per Constitution |
| Gate 4 (NFR) | Observability, security, PII, performance per Constitution |
| Gate 5 (Ready) | No ambiguity; Constitution referenced for all policy decisions |

---

## Creating Your Constitution

**File location:** `.specify/memory/constitution.md`

See `resources/constitution-template.md` for full template with examples per section.

**Key steps:**
1. Define 6 sections (UX, Security, Observability, Performance, Domain, Contracts)
2. Make rules specific and testable (not vague)
3. Get stakeholder sign-off (Architects, Security, Domain Owners)
4. Version the Constitution (v1.0, v1.1, etc.)
5. Declare version in every spec Metadata

---

## Using the Constitution in Specs

Every spec must pin the Constitution version:

```markdown
## Metadata
- Constitution Version: v3.0
```

And reference it per section:
```markdown
Source: Constitution/Observability Standards v3.0

Logging: JSON structured logs (request_id, timestamp, level, service)
Metrics: cart.duration_ms (histogram), cart.count (counter)
Tracing: Span "AddItem" with tags: item_id, qty
Alerts: IF p95_latency > 200ms THEN alert
```

---

## Key Rules

1. **Constitution is law** — Specs failing to comply fail Gate 2
2. **You can grant exceptions** — Exception = ADR (document why)
3. **Never delete rules** — Create ADR explaining why it's no longer enforced
4. **Version the Constitution** — Each change increments version (v1.0 → v1.1)

---

## Update Workflow

Constitution changes require:
1. PR to `.specify/memory/constitution.md`
2. Approval by Platform Architect + at least one Domain Owner
3. After merge: Platform MCP serves new content
4. All Approved/Implementing specs must recheck Gates 1 & 4
5. ADR if policy contradicts existing specs

---

## If you need more detail

→ `resources/constitution-template.md` — Full template with all sections, worked examples per section, gate enforcement rules, policy writing guidelines
