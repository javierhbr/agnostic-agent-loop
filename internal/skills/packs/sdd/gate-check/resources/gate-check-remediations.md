# Gate Check Remediations

Full failure tables, checklists, and fix instructions per gate. Referenced by the main gate-check skill.

---

## Gate 1 — Context Completeness

**What it checks:**
- Every spec section has a `Source:` line
- The `Context Pack:` field is declared in Metadata with a version
- `constitution.md` exists at `.specify/memory/constitution.md` and is non-empty

### Gate 1 Checklist

Walk through each section — every one must have a `Source:` line:

```
□ Metadata block has: Context Pack: cp-vN
□ Problem Statement has: Source: Platform MCP / Initiative [ID]
□ Goals has: Source: ...
□ User Experience has: Source: Platform MCP — UX guidelines vX
□ Domain Understanding has: Source: Domain MCP — [domain] vX
□ Cross-Domain Interactions has: Source: Domain MCP + Integration MCP
□ Contracts has: Source: Integration MCP — [contract names] vX
□ Technical Approach has: Source: Component MCP — [service] patterns vX
□ NFRs has: Source: Platform MCP — observability, security, performance vX
□ Observability has: Source: Platform MCP — logging/metrics/tracing standards vX
```

### Gate 1 Failure Fixes

| Symptom | Fix |
|---|---|
| Sections written without `Source:` | Add `Source: [MCP name + version]` to every section. Read `references/mcp-sources.md` in sdd-guide skill to know which MCP to cite per section |
| Missing Context Pack version in Metadata | Add `Context Pack: cp-vN` to the Metadata block. Get the version from the MCP Router output file |
| constitution.md missing or empty | Run `/speckit.constitution` in the Platform Repo first. Nothing else proceeds without it |
| Source cites "Platform MCP" without version | Add the version number: `Platform MCP v2.1`. Check the constitution's header for the current version |

---

## Gate 2 — Domain Validity

**What it checks:**
- No invariant violations (spec doesn't contradict rules in the Domain MCP)
- Domain ownership respected — no component accesses another domain's database directly
- All cross-domain communication is via versioned events or REST contracts only

### Gate 2 Failure Fixes

| Symptom | Fix |
|---|---|
| Spec shows Service A reading Service B's database | Change to event-driven: Service B emits an event, Service A consumes it. This is a domain boundary violation |
| Spec violates an invariant | Find the invariant in the Domain MCP file. Redesign the approach to respect it. If the invariant is wrong, update the Domain MCP first (Domain Owner must approve) |
| Spec assigns ownership of an entity to the wrong domain | Consult the Domain MCP to confirm ownership. Reassign responsibilities accordingly |
| "I don't have a domain MCP file" | Create `.specify/memory/domains/<domain>.md`. Read the Domain MCP section in sdd-guide's `references/mcp-sources.md` for the format |

### Key Question

**"Does your spec ever say one service will read or write data that belongs to another service's domain?"** If yes → that's the violation.

---

## Gate 3 — Integration Safety

**What it checks:**
- All contract consumers are identified (who depends on each event/API this component produces)
- A compatibility plan exists for any breaking changes
- If a contract version is bumped, a dual-publish strategy is defined

### Gate 3 Failure Fixes

| Symptom | Fix |
|---|---|
| Consumers not listed | Query the Integration MCP: "Who consumes [event/API name] v[N]?" List every consumer in the Contracts section |
| Breaking change with no compatibility plan | Decide: (a) dual-publish old+new versions for a defined period, or (b) add a new field without removing the old one. Document the plan in the Contracts section |
| No compatibility plan field in spec | Add a `### Compatibility Plan` subsection under Contracts |
| Contract change flagged but no Contract Spec in Platform Repo | Stop. Create the Contract Change Spec in Platform Repo (Integration Owner must approve) before continuing |

### Key Question

**"For every event or API your component produces, can you name every other service that consumes it?"** If not → check the Integration MCP.

---

## Gate 4 — NFR Compliance

**What it checks:**
- Logging is declared (what log events, at what level, in what format)
- Metrics are declared (what metric names, what dimensions/labels)
- Tracing is declared (what span names, what attributes)
- PII handling is specified (what fields are PII, how they are masked)
- Performance targets are set (p95 latency, throughput)

### Gate 4 Failure Fixes

| Symptom | Fix |
|---|---|
| NFR section is vague ("we'll add logging") | Be specific: "Order placement emits a structured JSON log event `order.placed` with fields: order_id, cart_id, amount (masked: no), user_email (masked: yes)" |
| No PII handling specified | Check each data field in your data model. Mark PII fields. State the masking rule: "user_email is masked at API boundary per Platform MCP security policy" |
| No performance target | State: "p95 latency < 300ms for synchronous flows" (or whatever the Platform MCP specifies) |
| Missing metrics | Define at minimum: a counter for success/failure, a latency histogram for the main operation |
| Missing traces | Define at minimum: a span for the main operation with attributes for the resource being operated on |

### Approach for Identifying PII Fields

1. List every data field in your Technical Approach section (user_id, user_email, order_total, etc.)
2. Mark each as PII or non-PII based on Platform MCP privacy policy
3. For every PII field, specify the masking rule — examples:
   - "user_email is masked at API boundary (never logged raw)"
   - "user_id is pseudonymized in logs (hashed with salt)"
   - "credit_card is never logged or stored; only the last 4 digits are visible"

---

## Gate 5 — Ready-to-Implement

**What it checks:**
- No open `BlockedBy` ADRs
- Spec is unambiguous — no vague or TBD sections
- All acceptance criteria are testable (can a QA engineer write a test from each criterion?)

### Gate 5 Failure Fixes

| Symptom | Fix |
|---|---|
| `BlockedBy: ADR-###` is still set | The ADR must reach `Approved` state. Use the `sdd-adr` skill to check status and unblock |
| Acceptance criteria say "works correctly" | Rewrite as testable: "Given a guest session token, GET /checkout/session returns 200 with order state PENDING within 300ms" |
| A section says "TBD" or "to be decided" | Make the decision now, or open an ADR if it needs more discussion. TBD is not acceptable |
| Tasks are not ordered by dependency | Run `/speckit.tasks` again. Tasks must show which ones can be done in parallel [P] and which ones must be sequential |

### Test for Gate 5

**Ask the human:** "Can a mid-level engineer pick up this spec tomorrow, with no other context, and know exactly what to build and how to verify it?"

If the answer is **no** → Gate 5 fails. Keep working until the spec is that clear.

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

Update `spec-graph.json` with `"status": "Approved"`.
