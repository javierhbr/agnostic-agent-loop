# Component Spec Guide

Full section walkthroughs, gate checklists, spec-graph format, and local ADR guidance. Referenced by the main component-spec skill.

---

## Walking Through Each Section

For each section below, you must declare a `Source:` line before writing content. This ensures specs are anti-invention — no section is written from memory or assumption.

### Problem Statement

```
Source: Platform MCP / Initiative ECO-124
```

Summarize the user-facing problem this component is solving. **Copy the relevant part from the Platform Spec** — do not rewrite it from scratch.

### Goals / Non-Goals

```
Source: Platform Spec PLAT-124 v1 — Component Responsibilities section
```

What this component specifically does and does not do. Keep it tight to this service.

### Domain Understanding

```
Source: Domain MCP — <your domain> invariants v<version>
```

List the invariants from your domain MCP that are relevant to this spec. For example:
- "Cart owns session state — not Checkout"
- "Guest session TTL ≥ 30 minutes after last activity"

**If an invariant is not in your domain MCP file — STOP.** Ask the Domain Owner to add it first. Do not invent invariants.

### Cross-Domain Interactions

```
Source: Domain MCP + Integration MCP — <contract names + versions>
```

Describe how this component interacts with other domains. Reference the exact event or API versions from the Integration MCP. For example:
- "Emits: CartUpdated v2 (schema: Integration MCP — CartUpdated v2)"
- "Consumes: OrderPlaced v3 (schema: Integration MCP — OrderPlaced v3)"

### Contracts

```
Source: Integration MCP — <contract names> v<version>
```

List every event and API this component publishes or consumes. Include versions. If a new version is needed, flag it here and confirm the Contract Change Spec exists in Platform Repo.

### Technical Approach

```
Source: Component MCP — <your service> patterns v<version>
```

How this component implements the feature. This is the section unique to the Component Spec. Follow approved patterns from your Component MCP. Reference your service's tech stack, data model, and local constraints.

### NFRs

```
Source: Platform MCP — observability, security, performance v<version>
```

Copy the relevant NFR requirements from the Platform MCP. Do not reduce or soften them.
- Logging: confirm log format
- Metrics: confirm metric names and dimensions
- Tracing: confirm span names and attributes
- PII: confirm masking rules at API boundary
- Performance: confirm p95 target for this flow

### Observability

```
Source: Platform MCP — logging/metrics/tracing standards v<version>
```

Specify exact log events, metric names, and trace spans this implementation will produce.

---

## Gate Check — All 5 Gates Detailed

### Gate 1 — Context Completeness

**What it checks:**
- Every spec section has a `Source:` line
- The `Context Pack:` field is declared in Metadata with a version
- `constitution.md` exists at `.specify/memory/constitution.md` and is non-empty

**Gate 1 checklist:**
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

**Most common fixes:**
| Symptom | Fix |
|---|---|
| Sections written without `Source:` | Add `Source: [MCP name + version]` to every section |
| Missing Context Pack version in Metadata | Add `Context Pack: cp-vN` to the Metadata block |
| constitution.md missing or empty | Run `/speckit.constitution` in the Platform Repo first |
| Source cites "Platform MCP" without version | Add the version number: `Platform MCP v2.1` |

### Gate 2 — Domain Validity

**What it checks:**
- No invariant violations (spec doesn't contradict rules in the Domain MCP)
- Domain ownership respected — no component accesses another domain's database directly
- All cross-domain communication is via versioned events or REST contracts only

**Most common fixes:**
| Symptom | Fix |
|---|---|
| Spec shows Service A reading Service B's database | Change to event-driven: Service B emits an event, Service A consumes it |
| Spec violates an invariant | Find the invariant in the Domain MCP file. Redesign the approach to respect it. If the invariant is wrong, update the Domain MCP first (Domain Owner must approve) |
| Spec assigns ownership of an entity to the wrong domain | Consult the Domain MCP to confirm ownership. Reassign responsibilities accordingly |
| "I don't have a domain MCP file" | Create `.specify/memory/domains/<domain>.md`. Format details in sdd-guide's `references/mcp-sources.md` |

**Key question:** Does your spec ever say one service will read or write data that belongs to another service's domain? If yes → that's the violation.

### Gate 3 — Integration Safety

**What it checks:**
- All contract consumers are identified (who depends on each event/API this component produces)
- A compatibility plan exists for any breaking changes
- If a contract version is bumped, a dual-publish strategy is defined

**Most common fixes:**
| Symptom | Fix |
|---|---|
| Consumers not listed | Query the Integration MCP: "Who consumes [event/API name] v[N]?" List every consumer |
| Breaking change with no compatibility plan | Decide: (a) dual-publish old+new versions for a defined period, or (b) add a new field without removing the old one. Document the plan in the Contracts section |
| No compatibility plan field in spec | Add a `### Compatibility Plan` subsection under Contracts |
| Contract change flagged but no Contract Spec in Platform Repo | Stop. Create the Contract Change Spec in Platform Repo (Integration Owner must approve) before continuing |

**Key question:** For every event or API your component produces, can you name every other service that consumes it?

### Gate 4 — NFR Compliance

**What it checks:**
- Logging is declared (what log events, at what level, in what format)
- Metrics are declared (what metric names, what dimensions/labels)
- Tracing is declared (what span names, what attributes)
- PII handling is specified (what fields are PII, how they are masked)
- Performance targets are set (p95 latency, throughput)

**Most common fixes:**
| Symptom | Fix |
|---|---|
| NFR section is vague ("we'll add logging") | Be specific: "Order placement emits a structured JSON log event `order.placed` with fields: order_id, cart_id, amount (masked: no), user_email (masked: yes)" |
| No PII handling specified | Check each data field in your data model. Mark PII fields. State the masking rule: "user_email is masked at API boundary per Platform MCP security policy" |
| No performance target | State: "p95 latency < 300ms for synchronous flows" (or whatever the Platform MCP specifies) |
| Missing metrics | Define at minimum: a counter for success/failure, a latency histogram for the main operation |
| Missing traces | Define at minimum: a span for the main operation with attributes for the resource being operated on |

### Gate 5 — Ready-to-Implement

**What it checks:**
- No open `BlockedBy` ADRs
- Spec is unambiguous — no vague or TBD sections
- All acceptance criteria are testable (can a QA engineer write a test from each criterion?)

**Most common fixes:**
| Symptom | Fix |
|---|---|
| `BlockedBy: ADR-###` is still set | The ADR must reach `Approved` state. Use the `sdd-adr` skill to check status and unblock |
| Acceptance criteria say "works correctly" | Rewrite as testable: "Given a guest session token, GET /checkout/session returns 200 with order state PENDING within 300ms" |
| A section says "TBD" or "to be decided" | Make the decision now, or open an ADR if it needs more discussion. TBD is not acceptable |
| Tasks are not ordered by dependency | Run `/speckit.tasks` again. Tasks must show which ones can be done in parallel [P] and which ones must be sequential |

**Test for Gate 5:** Can a mid-level engineer pick up this spec tomorrow, with no other context, and know exactly what to build and how to verify it? If the answer is no → Gate 5 fails.

---

## Spec-Graph.json Format

After your spec is Approved, update `spec-graph.json`:

```json
{
  "SPEC-CART-01": {
    "implements": "PLAT-124 v1",
    "context_pack": "cp-v2",
    "contracts_referenced": ["CartUpdated-v2", "OrderPlaced-v3"],
    "blocked_by": [],
    "status": "Approved",
    "affects": ["cart-service", "CartUpdated event"]
  }
}
```

**Fields:**
- `implements` — The Platform Spec ID and version this component spec fulfills
- `context_pack` — The context pack version used (from Metadata)
- `contracts_referenced` — All events and APIs referenced in the Contracts section
- `blocked_by` — Array of ADR IDs (empty if all resolved)
- `status` — Approved, Draft, Blocked, Done
- `affects` — What this spec impacts (services, events, systems)

---

## Local ADRs

If during spec writing you encounter an ambiguity that is scoped only to your component (does not affect other components or platform policy), create a Local ADR:

```
component-repo/adr/ADR-LOCAL-001-<short-title>.md
```

Add to your spec Metadata:
```
Blocked By: ADR-LOCAL-001
```

Use the `sdd-adr` skill to draft and track it.

**If the ambiguity affects multiple components or platform policy → escalate to the Platform Repo.** Do not create a local ADR for a cross-cutting concern.
