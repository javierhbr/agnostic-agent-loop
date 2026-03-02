---
name: hotfix
description: "Guides humans through the SDD hotfix and bug path. Trigger for production incidents, critical bugs, urgent fixes, or when asking 'how do I handle a bug in SDD'. Speed supported, governance never skipped."
---

# skill:hotfix

## Does exactly this

Routes production incidents and bugs through the SDD hotfix path with minimal specs required, enabling fast fixes while maintaining full traceability.

---

## When to use

- Production is down or critical bug discovered
- Payment, auth, or data corruption issue
- Urgent fix needed, but traceability is non-negotiable
- Non-urgent bug and you need structured bug tracking

---

## Step 1 — Triage: Routing Decision

Ask: **"Does this fix change any event schema, API contract, or require breaking a platform policy?"**

| Answer | Route | Action |
|--------|-------|--------|
| No — component-only fix | Component Repo | Create Hotfix Spec locally |
| Yes — contract-changing | Platform Repo FIRST | Create Contract Change Spec, update Integration MCP, THEN implement |
| Yes — policy-violating | Platform Repo FIRST | Create policy exception ADR, THEN implement |

---

## Step 2a — Normal Bug (Non-Urgent)

Create a lightweight Component Spec: `BUG-<number> — <Description>`

Sections:
- Metadata (ID, Component, Status: Draft)
- Reproduction (steps, expected vs actual, impact)
- Root Cause Hypothesis
- Fix Plan (minimal, scoped change only)
- Tests (unit, integration, regression)
- Gate Check (quick checklist: component-only? no contract changes? logging preserved?)

See `resources/hotfix-templates.md` for full template and examples.

---

## Step 2b — Hotfix (Production, Urgent)

Create a minimal Hotfix Spec: `HOTFIX-<number> — <Description>`

Sections:
- Metadata (ID, Component, Status: In Progress, Follow-up Spec: TBD)
- Issue (what's failing, impact, when started)
- Root Cause Hypothesis
- Fix (exact change: file, function, config)
- Rollback (exact revert steps)
- Validation (specific metric or log confirming recovery)

**Design principle:** Get this created FAST. Perfection does not block speed. See `resources/hotfix-templates.md` for examples.

---

## Step 3 — Implement the Fix

For component-only hotfix:
- Implement directly (one-liner config/flag change) OR
- Run `/speckit.implement` with Hotfix Spec as context

For contract-changing hotfix:
- Wait for Contract Change Spec approval in Platform Repo
- Then implement component side

Verify using the `Validation` criterion from the Hotfix Spec.

---

## Step 4 — Verify and Close

After fix is deployed:
1. Confirm `Validation` metric/log shows recovery
2. Update Hotfix Spec status to `Done`
3. Update `spec-graph.json` with Hotfix ID
4. Notify affected teams if contract changed

---

## Step 5 — Follow-up Spec (Required)

**Every hotfix MUST produce a Follow-up Spec within the next sprint.** This is not optional.

Follow-up Spec covers:
- [ ] Full test coverage for fixed behaviour
- [ ] Refactor if tech debt was introduced
- [ ] ADR if decision was made under pressure
- [ ] Contract Spec if a contract was touched
- [ ] Post-mortem: what invariant was violated? What gate would have caught this?

See `resources/hotfix-templates.md` for Follow-up Spec template.

---

## Governance Reminders

**On skipping the Hotfix Spec:**
> The Operating Model allows fast implementation. The minimal spec takes 10 minutes and preserves traceability permanently. We write it first, then fix.

**On skipping the Follow-up Spec:**
> The Follow-up Spec is where the system learns. Without it, the same bug recurs, decisions stay undocumented, and gaps stay open. It goes in the next sprint.

---

## If you need more detail

→ `resources/hotfix-templates.md` — Full Bug Spec and Hotfix Spec templates, Follow-up Spec template, triage decision tree, incident response examples
