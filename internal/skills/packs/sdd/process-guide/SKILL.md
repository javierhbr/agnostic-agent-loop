---
name: sdd-process-guide
description: "Complete step-by-step guide for executing the SDD v3.0 methodology from initiative definition through deployment and verification."
---

# skill:sdd-process-guide

## Does exactly this

Walks a team through the full **SDD v3.0 lifecycle** — from problem statement and risk assessment through architecture, parallel implementation, verification, and progressive production rollout — with gate checks enforced at every phase handoff.

---

## Use this skill when

- Starting a new initiative from scratch and needing the full SDD workflow
- Onboarding a role (PM, Architect, Developer, Verifier, DevOps) to their specific phase
- Unsure what the next step is in an in-progress SDD initiative
- Running the full workflow for a medium, high, or critical risk change

## Do not use this skill when

- The change is low-risk and a single task (use `agentic-agent task claim` directly)
- You only need one phase — consult the role-specific skill (analyst, architect, developer, verifier)

---

## Phases — in order, no skipping

**Phase 0 — Initiative Definition (PM)** — Define problem + measurable metric, run 5-question risk interview, create initiative, write GWT success criteria.

**Phase 1 — Architecture Design (Architect)** — Read initiative, check Platform Constitution, produce feature-spec + component-specs, run gate-checks, create dev tasks.

**Phase 2 — Implementation (Developers, parallel)** — Claim task, read component-spec, implement with tests + observability, produce impl-spec, gate-check, push + complete task.

**Phase 3 — Verification (Verifier)** — Verify every AC with test evidence, check observability in staging, run full test suite, produce verify.md, sync spec-graph, merge.

**Phase 4 — Deployment (DevOps/PM)** — Deploy with flag OFF, canary 2h, enable at 10% → 25% → 50% → 100%, measure success metrics at Day 30.

---

## CLI cheat sheet

```bash
# Phase 0
agentic-agent specifyify start "Feature Name" --risk [low|medium|high|critical]
agentic-agent specifyify workflow show [initiative-id]

# Phase 1
agentic-agent specifyify gate-check SPEC-[ID]
agentic-agent task create --title "Implement [Service]"

# Phase 2
agentic-agent task claim [TASK-ID]
agentic-agent specifyify gate-check SPEC-[SERVICE]-IMPL
agentic-agent task complete [TASK-ID]

# Phase 3
agentic-agent validate
agentic-agent specifyify gate-check SPEC-[ID]
agentic-agent specifyify sync-graph

# Phase 4
agentic-agent deploy --environment staging --feature-flags-all-off
agentic-agent flags set FeatureFlag=10pct
agentic-agent metrics watch [metric-name]
```

---

## Non-negotiable rules

1. Never skip a gate check — run `sdd gate-check` before every phase handoff
2. All ACs must be in Given/When/Then format and proven by runnable tests
3. Observability (logging + metrics + tracing) must be working in staging before merge
4. Feature flags default to OFF — production deploys are always safe
5. Rollback strategy must be declared before implementation begins
6. Blocked ADRs must be resolved before Gate 5 can pass

---

## Done when

- Phase 0: Initiative file created, risk classified, GWT success criteria written
- Phase 1: All 5 gates pass, dev tasks exist in backlog
- Phase 2: All component tasks complete, all gates pass per component
- Phase 3: verify.md produced, spec-graph synced, merged to main
- Phase 4: 100% rollout complete, Day-30 metrics meet or exceed targets

---

## If you need more detail

→ `resources/process-detail.md` — Full phase walkthroughs (Phase 0–4), risk scoring, spec templates, impl-spec template, verify.md template, deployment commands, troubleshooting
