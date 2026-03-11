---
name: verifier
description: "Hard stop before merge. Verifies implementation matches acceptance criteria with evidence, produces verify.md, updates Spec Graph."
---

# skill:verifier

## Does exactly this

Hard stop before merge. Verify implementation matches spec acceptance criteria with observable evidence, then update Spec Graph to Done.

---

## When to use

- Implementation is complete and all tests pass
- Before marking a change as ready for merge or production
- Need to produce evidence that every acceptance criterion is satisfied
- Updating Spec Graph to Done status

---

## Core Rules (Always Enforce These)

1. **NEVER merge without verify.md** — this is the hard stop
2. **EVERY acceptance criterion must have evidence** — test, log, or metric result
3. **Update Spec Graph with status = Done ONLY after verification**
4. **If ANY AC is untestable or unverified, BLOCK** — send back to Architect
5. **Mark REQUIRES HUMAN APPROVAL if touching payment, auth, or PII**

---

## Steps — in order, no skipping

1. **Read component-spec.md** — List every acceptance criterion verbatim (AC1, AC2, etc.).

2. **Read impl-spec.md** — Verify the mapping: each "Code Changes" entry → which ACs does it satisfy?

3. **Run tests and gather evidence** — For each AC: unit tests pass? integration tests pass? metrics recorded? lint/build pass?

4. **(Optional) Run Superpowers verification** — Use `superpowers:verification-before-completion` for final validation (if available).

5. **Write verify.md** — File: `.agentic/specs/[component-spec-id]/verify.md`. See resources for full template and examples.

6. **Update Spec Graph** — `agentic-agent specifyify sync-graph`, then `agentic-agent openspec complete [change-id]`.

---

## Blocking Conditions (Don't Merge If Any Are True)

- [ ] Any AC is UNTESTABLE
- [ ] Any test FAILS
- [ ] Lint or build FAILS
- [ ] Coverage < 80%
- [ ] Change touches payment/auth/PII without human approval
- [ ] Spec Graph not updated to Done

**If ANY are true: BLOCK.** Return to Architect with remediation.

---

## If you need more detail

→ `resources/verifier-templates.md` — verify.md full template, evidence examples, AC-to-evidence mapping, REQUIRES HUMAN APPROVAL escalation path
