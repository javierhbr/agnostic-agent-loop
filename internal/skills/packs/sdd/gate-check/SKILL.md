---
name: gate-check
description: "Diagnoses SDD gate failures and provides specific remediation. Trigger when /speckit.analyze fails a gate or when unsure what a gate requires."
---

# skill:gate-check

## Does exactly this

Diagnoses which of the 5 SDD gates is failing and tells you exactly what is missing and how to fix it.

---

## When to use

- `/speckit.analyze` returns a FAIL on any gate
- A spec won't pass a specific gate (Gate 2, Gate 3, etc.)
- You're unsure what a gate requires
- Reviewing a spec before fan-out to catch issues proactively

---

## Intake — Ask These 3 Questions

1. Which gate(s) are failing?
2. Are you at Platform level (before fan-out) or Component level (before implementation)?
3. Paste the relevant section of your spec that relates to the failing gate.

---

## The 5 Gates — One-Line Summaries

**Gate 1 — Context Completeness:** Every section has a `Source:` line, Context Pack version is declared in Metadata, constitution.md exists.

**Gate 2 — Domain Validity:** No invariant violations, no cross-domain direct database access, all communication via events/APIs only.

**Gate 3 — Integration Safety:** All contract consumers identified, compatibility plan exists for breaking changes.

**Gate 4 — NFR Compliance:** Logging, metrics, tracing, PII handling, and performance targets all explicitly declared.

**Gate 5 — Ready-to-Implement:** No open `BlockedBy` ADRs, no TBD sections, all acceptance criteria are testable by a QA engineer.

---

## Quick Questions to Diagnose Each Gate

**Gate 1:** Does every section start with `Source: [MCP name + version]`?

**Gate 2:** Does your spec say any service reads another service's database? If yes → that's the violation.

**Gate 3:** For every event/API your component produces, can you name every other service that consumes it?

**Gate 4:** Does the spec declare specific logging events, metric names, trace spans, and PII masking rules? If vague → fails.

**Gate 5:** Can a mid-level engineer pick this up tomorrow with no other context and know exactly what to build?

---

## After All Gates Pass

Update the spec's Gates section:
```
## Gates
- Context completeness: PASS
- Domain validity: PASS
- Integration safety: PASS
- NFR compliance: PASS
- Ready to implement: PASS
```

Update `spec-graph.json`: `"status": "Approved"`

---

## If you need more detail

→ `resources/gate-check-remediations.md` — Full failure tables per gate, section-by-section checklist, worked examples, fix instructions
