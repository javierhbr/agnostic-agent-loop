---
name: workflow-router
description: "Workflow routing guide for SDD v3.0. Explains risk classification, three workflow types, and how to activate the right agents."
---

# Workflow Router

## Decision Tree: How to Choose Your Workflow

```
Does this change involve:
├─ Payment, authentication, or PII data?
│  └─ YES → High/Critical Risk → FULL Workflow
├─ Breaking an existing contract (API/event schema)?
│  └─ YES → Medium/High Risk → STANDARD or FULL Workflow
├─ Multiple services working in parallel?
│  └─ YES → Medium Risk → STANDARD Workflow
├─ Non-trivial rollback concerns?
│  └─ YES → Medium/High Risk → STANDARD or FULL Workflow
└─ Bug fix or single-service feature?
   └─ YES → Low Risk → QUICK Workflow
```

## Three Workflow Types

### QUICK Workflow (Low Risk, Bug Fixes)
**Agents:** Developer → Verifier

**When to use:**
- Bug fixes scoped to one service
- Performance improvements with no contract changes
- Internal refactors (no user-facing changes)
- Low blast radius

**Time to merge:** 1-2 days

### STANDARD Workflow (Medium Risk, Cross-Service)
**Agents:** Architect → Developer (parallel per component) → Verifier

**When to use:**
- New features touching 2-3 services
- Minor contract changes with compatibility plan
- Medium blast radius
- Well-understood domain

**Time to merge:** 3-5 days

### FULL Workflow (High/Critical Risk)
**Agents:** Analyst → Architect → Developer (parallel) → Verifier

**When to use:**
- Payment, authentication, or PII changes
- Breaking contract changes
- 4+ services involved
- High blast radius
- Novel or poorly understood domain

**Time to merge:** 5-10 days (plus human approval)

## CLI Commands to Activate a Workflow

```bash
# Determine risk level (interactive prompt or config)
agentic-agent specifyify start "Feature Name" --risk medium

# Show current workflow progress
agentic-agent specifyify workflow show <initiative-id>

# Run gate checks on a spec
agentic-agent specifyify gate-check <spec-id>

# Create ADRs to resolve blockers
agentic-agent specifyify adr create --title "Idempotency Strategy" --scope local
agentic-agent specifyify adr resolve ADR-001

# Sync spec graph to platform repo
agentic-agent specifyify sync-graph

# Install agent Markdown files
agentic-agent specifyify agents install
```

## Per-Workflow Agent Sequence

### Quick
1. **Developer** — reads change spec, produces impl-spec.md + tasks.yaml
2. **Verifier** — verifies all ACs with observable evidence, produces verify.md

### Standard
1. **Architect** — reads discovery/initiative, produces feature-spec.md + component-specs
2. **Developers** (parallel per component) — each produces impl-spec.md + tasks.yaml
3. **Verifier** — verifies all ACs across all components, updates spec-graph.json

### Full
1. **Analyst** — interviews team, produces discovery.md with evidence-based risk classification
2. **Architect** — reads discovery, produces feature-spec.md + component-specs
3. **Developers** (parallel) — each produces impl-spec.md + tasks.yaml
4. **Verifier** — verifies all ACs, updates spec-graph.json
5. **Human approval** (for critical changes) — final sign-off required

## Handoff Protocol

Each agent's output becomes the next agent's input. Gates must PASS before handoff.

```
QUICK:    change-spec → developer:impl-spec → verifier:verify.md → DONE
         (gates 1,4,5 check)           (gates 1,4,5 check)

STANDARD: discovery → architect:specs → developers:impl-specs → verifier:verify.md → DONE
         (no analyst) (all gates check)  (all gates check)        (all gates check)

FULL:     (initiative) → analyst:discovery → architect:specs → developers:impl → verifier → DONE
                       (gates check)        (all gates check)   (all gates)      (all gates)
```

## Key Rules

1. **Never skip an agent.** Each one adds critical value.
2. **Gates must PASS before handoff.** If any gate fails, the spec goes back to the prior agent.
3. **Risk level determines workflow.** Misclassified risk is the most common failure mode.
4. **ADRs block implementation.** If `blocked_by` is non-empty, implementation cannot proceed.
5. **Spec Graph is the audit trail.** It's the source of truth for what's been done.
