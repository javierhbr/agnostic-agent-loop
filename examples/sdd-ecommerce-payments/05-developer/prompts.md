# Developer Phase: Prompts

## Overview

After the Architect resolves ADR-001 and all component-specs pass gates, the Developer agents work in **parallel per component**. Two or more Developer agents pick up different component-specs simultaneously and produce `impl-spec.md` + `tasks.yaml` for each.

**Key Point:** Developers run in parallel but wait if a component-spec has `blocked_by` non-empty.

---

## Prompt for Developer Agent (Billing Service)

```
You are a Developer agent in the SDD Full workflow. You will implement the Billing Service
component of the Subscription Billing feature.

INPUTS:
- Feature Spec: 04-architect/feature-spec.md
- Component Spec: 04-architect/component-spec-billing.md
- ADR-001: 04-architect/ADR-001-card-vault.md (assume resolved: Stripe chosen)

YOUR JOB:
1. Read component-spec-billing.md thoroughly
2. Call Component MCP: get_patterns() — how do we implement subscription state machines?
3. Call Component MCP: get_decisions() — any prior ADRs about billing systems?
4. Check blocked_by field — STOP if non-empty
5. Produce impl-spec-billing.md with:
   - Data model (tables, fields, constraints)
   - Code changes (functions, modules to create/modify)
   - Edge cases (4+ cases: customer cancels mid-cycle, payment fails, renewal retries, etc.)
   - Observability (logging, metrics, traces in code)
   - Rollout plan (feature flag strategy, canary, rollback)
6. Produce tasks.yaml with 8-10 actionable development tasks
7. Self-check all 5 gates before marking done

RULE: NEVER skip the component-spec. It is your contract.

Edge cases to cover:
- What if customer cancels mid-cycle? Do they get partial refund?
- What if payment fails 3 times? Auto-suspend subscription?
- What if merchant updates plan price? Affect active subscriptions?
- What if charge is processed but subscription.next_charge_date update fails (crash)?

OUTPUTS:
- impl-spec-billing.md (2000-3000 words)
- tasks.yaml (8-10 tasks with dependencies)

START NOW. Read the component-spec.
```

---

## Prompt for Developer Agent (Payments Service)

```
You are a Developer agent implementing the Payments Service component (recurring charges).

INPUTS:
- Feature Spec: 04-architect/feature-spec.md
- Component Spec: 04-architect/component-spec-payments.md
- ADR-001: Resolved (Stripe Vault chosen)

YOUR JOB:
1. Read component-spec-payments.md
2. Call Component MCP: get_patterns() — Stripe integration patterns?
3. Call Component MCP: get_decisions() — payment processing ADRs?
4. Produce impl-spec-payments.md with:
   - Stripe SDK integration (which library, version)
   - Code changes (Stripe tokenization endpoint, charge processing)
   - Edge cases (4+ cases: card expired, insufficient funds, Stripe timeout, double-charge prevention)
   - Observability (logging every Stripe API call)
   - Rollout plan (Stripe test mode, prod keys, feature flag)
5. Produce tasks.yaml (8-10 tasks)
6. Self-check gates

Edge cases:
- Stripe API timeout (30 seconds) — what's our fallback?
- Card expired but customer doesn't know — how do we notify?
- Stripe returns 500 error — retry or give up?
- Idempotency key collision — is Stripe smart enough?

OUTPUTS:
- impl-spec-payments.md
- tasks.yaml

START NOW.
```

---

## Expected Output Format

Both Developer agents will produce files like:
- `impl-spec-billing.md` (implementation-level spec for Billing service)
- `impl-spec-payments.md` (implementation-level spec for Payments service)
- `tasks.yaml` (YAML file with task decomposition)

See examples below (impl-spec-billing and tasks.yaml are included).
