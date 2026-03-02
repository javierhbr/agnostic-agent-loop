# Target Operating Model
## SDD + MCP Enterprise Model — Version 3.0 · February 2026

> *"Nothing is implemented without a validated spec backed by governed context."*
> — Core Principle, v3.0

A practical guide to Spec-Driven Development with two MCP servers, four Markdown agents, and zero overhead services — the enterprise model for teams that need alignment, traceability, and zero tribal knowledge.

---

## Problem Statement

**Speed vs. correctness:** AI agents can write code in seconds — but they invent context. Without governed specs, every agent call introduces drift between what was intended and what was built.

**Tribal knowledge doesn't scale:** Contracts, invariants, and architectural decisions live in Slack threads and people's heads. When a new agent or developer touches the system, they guess — and guesses in trading systems are expensive.

**Process overhead kills velocity:** Heavy methodologies (router services, gate validators, contract registries as servers) add infrastructure that nobody maintains. Teams skip the process or drown in it.

**The answer:** Two MCP servers that serve governed knowledge. Four Markdown agents that carry the process. Zero overhead services — gates live as checklists inside templates, not in TypeScript validators. If a methodology component doesn't change what the agent writes, it's eliminated.

---

## 00 — The Big Idea

Your organisation is a knowledge system, not just a code system. Code is the output — specs, contracts, and decisions are the inputs that determine whether the code is right.

| ⚠ Without SDD | ✓ With SDD + MCP |
|---|---|
| Teams interpret requirements differently | Specs are the explicit, versioned source of truth |
| Service contracts break without warning | Contracts are governed with full consumer lists |
| Decisions live in Slack and people's heads | Every decision is an ADR with an audit trail |
| AI agents invent context from the codebase | Agents consume governed context — no invention |
| Rework after costly integration bugs | Gates block integration bugs before they happen |
| Onboarding requires tribal knowledge | Context Packs enable fast, safe onboarding |

---

## 00b — Why This Model — What We Took & What We Dropped

This model is not a new methodology. It is a deliberate composition of proven ideas — taking only the pieces that directly change what an agent or engineer writes, and cutting everything else.

Every component in this operating model was borrowed from an existing discipline. The selection criterion was ruthless: *"Does this piece change the output of the person or agent doing the work? If not, it's eliminated."* The result is a model where all process lives in Markdown files — versionable, readable, editable by anyone on the team via PR — and zero process lives in infrastructure.

Three frameworks contributed the core structural ideas: **BMAD** (Breakthrough Method for Agile AI Development) defined the four-role agent architecture and handoff protocol. **OpenSpec** defined the spec format, mandatory traceability metadata, and the Spec Graph linking all artifacts. **SpecKit** defined the template-driven assembly pattern — templates with embedded exit gate checklists and the Context Pack versioning mechanism. All three are documented below alongside SDD, MCP, ADR, Agile/XP, and Lean.

---

### 📐 SDD — Spec-Driven Development
*From: enterprise software delivery practices*

**What we took**

**Spec-before-code rule** — no implementation begins without a written, validated spec. Not a PR description. An actual spec with ACs, data model, and edge cases.
`→ prevents context invention`

**Acceptance Criteria in Given/When/Then** — every feature has verifiable, testable conditions that the Verifier agent checks with observable evidence.
`→ makes "done" unambiguous`

**Exit gates as self-verification** — instead of external validators, gates are checklists embedded in templates. The agent self-checks before handing off.
`→ removes overhead server`

~~**Centralised spec registry with automated snapshots**~~ — spec databases with change tracking, contract registries as standalone services, automated snapshot pipelines.
`✗ doesn't change what gets written`

---

### ⬟ MCP — Model Context Protocol
*From: Anthropic — structured agent context*

**What we took**

**Platform MCP** — one instance that serves context packs, templates, agent definitions, and workflow configs. The agent calls a tool and gets governed truth, not a blank context window.
`→ eliminates context invention`

**Component MCP (N instances, same image)** — one per component, serving that component's contracts, invariants, patterns, and decisions. Local truth scoped to where the agent is working.
`→ right context, right scope`

~~**MCP Router**~~ — a central dispatcher that routes agent queries to the correct MCP server, aggregating responses into a single context pack.
`✗ agents know their role — no router needed`

~~**Separate contracts-mcp and domain-mcp**~~ — dedicated servers for contract registries and domain invariants, each with their own deployment lifecycle.
`✗ contracts are component knowledge — get_contracts() on the component MCP`

---

### 📋 ADR — Architecture Decision Records
*From: Michael Nygard (2011) · widely adopted in microservices*

**What we took**

**Decisions as first-class artifacts** — every significant technical choice is a Markdown file with context, decision, and consequences. Component MCP surfaces them via `get_decisions()`.
`→ no repeated debates`

**ADR as a gate blocker** — an open, unresolved ADR blocks the Verifier from clearing the merge gate. Ambiguity is resolved before code is written, not after.
`→ no spec can escape an unresolved decision`

~~**ADR query service**~~ — a queryable database of decisions with cross-repo search, impact analysis APIs, and automated staleness detection.
`✗ Markdown files in git already provide this`

---

### 🔄 Agile / XP — Extreme Programming & Agile
*From: Beck, Manifesto for Agile (2001)*

**What we took**

**Risk-proportional process** — quick workflow for low-risk/bug-fix changes, full workflow only for critical. The weight of the process matches the weight of the risk.
`→ eliminates unnecessary overhead`

**Small deliverables per agent** — each agent produces one or two files and hands off. No monolithic analysis → implementation cycles; the Architect can produce specs and Developers can start in parallel.
`→ parallel execution per component`

~~**Story points, velocity tracking, sprint ceremonies**~~ — estimation rituals, retrospective formats, burndown charts as mandatory process.
`✗ tasks.yaml captures owner + points when needed — no ceremony`

---

### ✂️ Lean — Lean Software Thinking
*From: Poppendieck, Toyota Production System adapted*

**What we took**

**Golden rule — eliminate what doesn't change output** — every component in this methodology was evaluated against a single question: "Does this change what the agent or engineer writes?" If no, it was cut.
`→ the core filter for the whole model`

**Logic in files, not infrastructure** — gates, workflows, agent behaviour, and templates all live in Markdown files in a git repo. Zero TypeScript validators, zero gate-server deployments.
`→ no server = no maintenance cost`

**Knowledge feedback loop** — every implementation cycle updates component MCP sources (patterns, invariants, ADRs), making the next initiative start from a richer context. The system learns.
`→ compound knowledge accumulation`

---

### 🤖 BMAD — Breakthrough Method for Agile AI Development
*From: BMAD-METHOD · AI-native multi-agent workflow design*

**What we took — and why**

**Four bounded agent roles** — analyst, architect, developer, verifier. Each role has a single, non-overlapping responsibility and produces exactly one or two artifacts before handing off. An agent never does another agent's job. *Why: role-scoped context prevents agents from overreaching and generating untraceable outputs.*
`→ no scope creep per agent`

**Sequential handoff with hard stops** — each agent must complete its exit gate checklist before the next agent starts. The Verifier is a hard stop: no merge until every AC has observable evidence. *Why: handoffs make quality gates enforceable without an external service — the agent itself is the gate.*
`→ quality enforced at each transition`

**Parallel fan-out from Architect** — once `feature-spec.md` is approved, one Developer agent runs per affected component simultaneously. The Verifier waits for all. *Why: eliminates the sequential bottleneck that makes large cross-domain features slow — components that don't depend on each other can be built concurrently.*
`→ parallel execution without loss of traceability`

**Agent definitions as Markdown files** — each agent is a `.md` file in `.claude/agents/` loaded on demand. Behaviour lives in text, not code. *Why: the team can evolve agent behaviour via PR without a deployment — the same update path as specs and templates.*
`→ agent evolution without deploys`

~~**Orchestrator agent**~~ — a meta-agent that decides which specialist agent to invoke, routing requests based on request type classification.
`✗ Platform MCP get_workflow() replaces this — risk level determines the agent sequence, not a separate agent`

~~**Persona-based agent naming**~~ — agents named after team member personas ("Alex the Architect", "Dev Dana") to guide interaction style.
`✗ role names are sufficient — persona framing doesn't change what gets written`

---

### 📄 OpenSpec — Open Specification Format
*From: OpenSpec format · structured spec authoring for distributed systems*

**What we took — and why**

**Mandatory traceability metadata on every spec** — `implements` (parent spec ID + version), `context_pack` (the versioned snapshot used when the spec was authored), `blocked_by` (open ADR IDs), and `status`. *Why: without these fields, a spec is an island — you cannot tell what it is implementing, what context it was written against, or whether it is safe to act on.*
`→ spec is never an island`

**Spec Graph — `graph/index.yaml`** — a machine-readable index linking every artifact in a chain: Initiative → feature-spec → component-specs → impl-specs → verify.md → ADRs. Updated by the Verifier after every merge. *Why: the graph is the audit trail. It answers "what was built, why, against what context, and was it verified?" — questions that matter in regulated environments and post-incident analysis.*
`→ complete audit trail`

**MCP source citation per section** — every section of a spec declares where its content came from: which MCP tool, which version. *Why: a spec section without a source is an assertion without evidence. Citing the source makes the spec verifiable and makes it easy to see when a spec is stale relative to a newer MCP context.*
`→ every claim has a cited source`

**Specs are never deleted — only versioned** — status transitions from Approved → Implementing → Done, or → Paused, or → Blocked. A spec is never removed from the graph. *Why: deleted specs create invisible debt. A Paused spec with its context pack pinned can be rebased and resumed; a deleted spec is lost work that gets redone from scratch.*
`→ no invisible technical debt`

~~**Spec versioning server**~~ — a dedicated service managing spec versions with diff APIs, conflict detection, and branch-per-initiative spec trees.
`✗ git already handles versioning — the spec graph provides the links without a separate server`

---

### 📦 SpecKit — Template-Driven Spec Assembly
*From: SpecKit · template-driven spec assembly for AI-assisted development*

**What we took — and why**

**Templates with embedded exit gate checklists** — spec templates are not blank forms. Each section of the template includes the checklist the agent must self-verify before handing off. *Why: separating "how to write the spec" from "how to verify the spec" creates a gap where quality degrades. Embedding the gate in the template closes that gap — the agent writes and checks in one pass.*
`→ write and verify in one pass`

**Context Pack — versioned snapshot assembled before spec writing** — before any spec is written, the agent calls Platform MCP to assemble a pinned snapshot of all applicable policies, NFR baselines, and workflow config for this intent. *Why: writing a spec against a moving context is the root cause of spec drift. Pinning the context pack version makes the spec reproducible — another agent reading the same spec can retrieve the same context and understand exactly what constraints were in effect.*
`→ specs are reproducible`

**Template-driven spec types** — feature-spec, component-spec, impl-spec, and verify.md each have a dedicated template. Every instance of a spec type has the same structure. *Why: uniform structure means agents, humans, and CI can parse and validate specs without special-casing. The Verifier knows exactly where to find the ACs because every component-spec has them in the same place.*
`→ consistent structure → parseable by agents and CI`

**Templates served by Platform MCP (`get_template(name)`)** — templates are not stored locally with the agent. The agent calls the Platform MCP to fetch the current template version. *Why: this makes template evolution zero-friction — update the template in the specs repo and every agent immediately gets the new version on the next call, with no agent file changes required.*
`→ template updates propagate instantly`

~~**Spec scaffolding CLI**~~ — a command-line tool that generates spec file stubs, pre-fills metadata, and validates spec structure on creation.
`✗ Platform MCP get_template() replaces this — the agent fetches and fills the template inline, no CLI required`

---

## 01 — Two MCPs, Four Agents

Two MCP server types carry all the governed knowledge. Four Markdown agents carry all the process. No routers, no gate validators, no extra services — just logic in files.

### ⬟ Platform MCP — Platform Intelligence
*1 instance per org · reads openclaw-specs/ · 4 tools only*

The single source of platform-level truth. Provides context packs, templates with embedded exit gates, agent definitions, and workflow configs. Never needs redeployment — just update the specs repo.

| Tool | Returns |
|---|---|
| `get_context_pack(intent)` | risk_level + policies + workflow |
| `get_template(name)` | template with exit gate checklist embedded |
| `get_agent(name)` | analyst / architect / developer / verifier |
| `get_workflow(risk_level)` | which agents to activate |

### ◈ Component MCP — Component Knowledge
*N instances · same image · different config per component*

One instance per component. Serves that component's contracts, invariants, patterns, and technical decisions. The "local truth" each agent needs before writing a single line.

| Tool | Returns |
|---|---|
| `get_contracts()` | NATS topics, HTTP endpoints, event schemas |
| `get_invariants()` | immutable business rules for this component |
| `get_patterns()` | approved implementation patterns + canonical examples |
| `get_decisions()` | ADRs + prior technical decisions |

### ▤ Four Agents — Development Workflow
*Markdown files in `.claude/agents/` · loaded on demand by Claude Code*

Four Markdown agents define the entire development process. Gates live inside the agent as checklists — not as external validators. Agents self-verify before handing off.

| Agent | Workflow | Produces |
|---|---|---|
| `analyst.md` | Full only | `discovery.md` — evidence, affected components, success metrics |
| `architect.md` | Standard + Full | `feature-spec.md` + `component-spec.md` per domain |
| `developer.md` | All · parallel per component | `impl-spec.md` + `tasks.yaml` |
| `verifier.md` | All · hard stop | `verify.md` — AC verification, sign-offs, merge gate |

---

## 02 — How the Layers Flow

Knowledge at the top feeds specs in the middle, which drive code at the bottom. Each layer depends only on what is above it — never sideways or upward.

```
┌─────────────────────────────────────────────────────────────────┐
│  KNOWLEDGE LAYER  —  MCP Servers · governed context             │
│  Platform Policies · Domain Models · Contracts · Component Ctx  │
└────────────────────────────┬────────────────────────────────────┘
                             │
              ↓  Context Pack (versioned)
                 assembled by Platform MCP before any spec is written
                             │
┌────────────────────────────▼────────────────────────────────────┐
│  SPEC LAYER  —  Agents · 2 repos                                │
│  Platform Specs · Component Specs · Contract Specs · ADRs       │
└────────────────────────────┬────────────────────────────────────┘
                             │
              ↓  Execution — implementation against
                 approved, gate-checked specs
                             │
┌────────────────────────────▼────────────────────────────────────┐
│  CODE LAYER  —  Implementation                                  │
│  Services · APIs · Infrastructure                               │
└────────────────────────────┬────────────────────────────────────┘
                             ↓
┌─────────────────────────────────────────────────────────────────┐
│  OBSERVABILITY  —  Runtime signals                              │
│  Logs · Metrics · Traces                                        │
└─────────────────────────────────────────────────────────────────┘
```

---

## 03 — Repository Structure

Two repos. One for content (specs, agents, templates, policies) — editable by the whole team via PR. One for the system (code). Two types of MCP server. Zero overhead services.

### Specs Repo — Content · editable via PR by anyone

```
openclaw-specs/
├── constitution/        ← policies, NFRs
├── agents/              ← 4 Markdown agents
├── workflows/           ← quick / standard / full
├── templates/           ← with exit gates embedded
├── components/          ← specs + ADRs per component
└── graph/index.yaml     ← traceability index
```

### System Repo — Code · maintained by engineering

```
openclaw/
├── .claude/agents/      ← symlinks to specs/agents/
├── agents/mario_hugo/
│   └── docs/specs/      ← impl-spec + verify.md
├── internal/risk/       ← Guaripolo
└── config/              ← all overrides
```

### MCP Servers — 2 types · no router · no overhead

```
platform-intelligence/
├── platform-mcp/        ← 1 instance
│   └── 4 tools only     ← context, template, agent, workflow
└── component-mcp/       ← N instances
       same image · different config per component
```

---

## 04 — Development Lifecycle

Every initiative follows the same flow — from product intent to shipped code. The workflow activated depends on risk level.

| Workflow | Risk Level | Agent Sequence |
|---|---|---|
| **Quick** | low · bug_fix | Developer → Verifier |
| **Standard** | medium · new feature | Architect → Developer → Verifier |
| **Full** | high · critical | Analyst → Architect → Developer → Verifier |

---

### Step 1 — Analyst Agent
*`analyst.md` · Full workflow only · high / critical*

Defines the problem, evidence, affected components, and success metrics. One question at a time — real data from the team, not template examples. Produces `discovery.md`.

**MCP calls:**

| MCP | Tool | Returns |
|---|---|---|
| Platform MCP | `get_context_pack(intent)` | risk_level, applicable policies, workflow to follow |
| Platform MCP | `get_template("discovery")` | discovery.md template with embedded exit gate checklist |

**Exit gate — analyst self-checks before handoff:**
- Problem Statement with concrete metric
- Evidence with ≥1 real data point
- ≥2 User Stories
- ≥2 Out of Scope exclusions
- Affected Components table
- Success Metrics with timeframe

---

### Step 2 — Architect Agent
*`architect.md` · Standard + Full workflows · skipped for Quick*

Defines the "what": UX flow, domain responsibilities, NFRs, contract changes, risk flags. Produces one `feature-spec.md` (platform level) + one `component-spec.md` per affected component. Developers can start in parallel per component.

**MCP calls — per affected component:**

| MCP | Tool | Returns |
|---|---|---|
| Platform MCP | `get_context_pack(intent)` | policies, NFR baselines, constitution rules |
| Component MCP | `get_contracts()` | NATS topics, API endpoints, event schemas |
| Component MCP | `get_invariants()` | business rules that cannot be violated |
| Component MCP | `get_decisions()` | ADRs — Architect cannot contradict an approved ADR |

**Exit gate — architect self-checks before handoff:**
- ≥3 ACs in Given/When/Then
- NFRs with concrete numbers
- Contract Changes declared or "none" explicit
- Feature flag + rollback strategy defined
- Component spec per affected component
- If critical: "REQUIRES HUMAN APPROVAL" flag present

---

### Step 3 — Developer Agent ×N
*`developer.md` · All workflows · runs per component in parallel*

Defines the "how": exact data model, code changes, edge cases, observability. Produces `impl-spec.md` + `tasks.yaml` in the component's `docs/specs/` folder. Reads all context before writing a single line — MCP patterns are the starting point, not approximations.

**MCP calls:**

| MCP | Tool | Returns |
|---|---|---|
| Component MCP | `get_patterns()` | approved implementation patterns + canonical examples |
| Component MCP | `get_decisions()` | existing ADRs — flags contradiction to Architect before continuing |

**Exit gate — developer self-checks before handoff:**
- Data model with all fields + types
- All endpoints/functions with request/response/errors
- Edge cases table (≥4 rows)
- Observability section (metrics, alerts)
- Rollout plan with config flag
- `tasks.yaml` with owner + points + dependencies

---

### Step 4 — Verifier Agent
*`verifier.md` · All workflows · waits for all developers · **hard stop***

Verifies every AC with observable evidence. No merge without all checks marked. Produces `verify.md`. For critical changes (Guaripolo, schema migrations, live flag), merge is blocked until human sign-off is explicitly documented in `verify.md`.

**MCP calls:**

| MCP | Tool | Returns |
|---|---|---|
| Platform MCP | `get_context_pack(intent)` | correct DoD based on risk_level — critical requires 4 sign-offs, low requires 1 |

**Exit gate — merge blocked until ALL checks marked:**
- Every AC verified with observable evidence
- Tests passing (output attached)
- Observability confirmed in staging
- Rollback tested
- Sign-offs per risk level
- If critical: "HUMAN APPROVED" with name + date in verify.md

---

### Step 5 — Release + Knowledge Update
*Team / Domain Owner · post-merge*

After merge, update the Spec Graph index and component MCP sources. Any new invariant, pattern, or ADR discovered during implementation enriches the context available to the next initiative. The system gets smarter with each cycle.

**What gets updated post-release:**

| MCP | What | How |
|---|---|---|
| Component MCP | New patterns, runbooks, constraints | `components/{domain}/specs/` and `/adrs/` updated via PR |
| Platform MCP | New constitution rules, NFR baselines | `openclaw-specs/constitution/` updated via PR |
| Component MCP | Contracts shipped | `get_contracts()` returns updated versions; deprecations with sunset dates |
| Platform MCP | Spec Graph | `graph/index.yaml` updated — initiative → feature-spec → component-specs → verify.md fully linked |

---

### Fan-out: Architect → Developer Agents (parallel)

```
Architect · feature-spec.md ready
              Fan-out → per component
                 │
        ┌────────┼────────┐
        │        │        │
   mario-hugo  perform  ⛔ guaripolo
   ─────────   ────────   ───────────
   OC-047      OC-047     OC-048
   momentum    feedback   cb-rules
   no contract no contract contract: yes
   blocked: —  blocked: — HUMAN APPROVAL
```

---

## 05 — The 5 Gates

Non-negotiable checkpoints enforced at CLI, CI, PR, and agent levels. If any gate fails, progress is blocked — no exceptions.

| Gate | Name | Checks |
|---|---|---|
| **G1** | Context Completeness | All MCP sources cited with versions · Context Pack version pinned · `constitution.md` exists and non-empty |
| **G2** | Domain Validity | No invariant violations · Domain ownership respected · No direct cross-domain DB access |
| **G3** | Integration Safety | All contract consumers identified · Compatibility plan for breaking changes · Dual-publish strategy defined |
| **G4** | NFR Compliance | Logging, metrics, tracing declared · PII handling specified · Performance targets set (p95) |
| **G5** | Ready to Implement | No open BlockedBy ADRs · Spec is unambiguous · All acceptance criteria testable |

> 🚫 **If any gate FAILS →** the system must block progress. The failure must be resolved before moving forward. Self-enforced by the Verifier agent checklist at merge time, and by CI spec link validation in both repos.

---

## 06 — Example Prompts

Real prompts for a Spec-Oriented Agent — showing how each role interacts with the system in practice. The agent always runs Phase 1 (identify) before Phase 2 (produce spec).

> Send these to a Spec-Oriented Agent (Claude Code, Cursor, Copilot) configured with the SDD system prompt. The agent must identify the request type and missing context before generating any spec or code — jumping straight to implementation is a system violation.

---

### New Feature

**Product Manager — Kicking off a new cross-domain initiative**
> *Context: Starting Guest Checkout (ECO-124). No spec exists yet. Context Pack cp-v2 is ready.*

Prompt:
```
We need guest checkout on our e-commerce platform. Users complete a full purchase
without an account. Cart, Checkout, Payments, and Fulfillment are all involved.
Create the platform spec for initiative ECO-124 using context pack cp-v2.
```

Result: Agent runs Phase 1 — identifies cross-domain feature, maps domains (cart, checkout, payments, fulfillment), checks for missing context. Phase 2 produces a full Platform Spec: UX flow, domain responsibilities, contract baselines, NFRs, and complete Gate Check with all 5 gates evaluated.

---

**Product Manager — Surfacing ambiguities before planning begins**
> *Context: Architect agent has produced the feature spec. Running clarify pass before handing off to Developer agents.*

Prompt:
```
Run clarify on the guest checkout platform spec. Surface all ambiguities before we
write the plan. Flag anything that needs an ADR and assign owners.
```

Result: Agent surfaces open questions: idempotency strategy for payment retries, session TTL for abandoned carts, which service owns `guest_token` generation, Fulfillment dedup strategy. Creates ADR Drafts ADR-219, ADR-220, ADR-221. Marks spec `BlockedBy: ADR-219, ADR-220, ADR-221`. Prevents planning until ADRs are resolved.

---

### Contract Change

**Integration Owner — Proposing a breaking change to OrderPlaced**
> *Context: Need to add `payment_intent_id` and `guest_email` to OrderPlaced event for ECO-124.*

Prompt:
```
I need to add two new fields to OrderPlaced: payment_intent_id (string) and
guest_email (string, PII). Required for guest checkout. Check who consumes
OrderPlaced v2, assess whether this is a breaking change, and draft Contract
Change Spec SPEC-CONTRACT-77 with a compatibility plan.
```

Result: Agent queries Integration MCP — finds 4 consumers of OrderPlaced v2 (fulfillment, shipping, analytics, notification). Flags `guest_email` as PII, cites SEC-001 from Platform MCP (PII masking required at API boundary). Identifies this as a breaking change requiring dual-publish. Generates SPEC-CONTRACT-77: dual-publish v2 + v3 for 30 days, deprecation schedule, consumer notification plan.

---

**Integration Owner — Verifying consumer migration before sunsetting v1**
> *Context: CartUpdated v1 has been deprecated. Dual-publish window ends next week.*

Prompt:
```
Check whether all consumers of CartUpdated v1 have migrated to v2. If confirmed,
create the deprecation closure notice and mark v1 as sunset-ready in the contract
registry.
```

Result: Agent queries Integration MCP for CartUpdated v1 consumer list — finds `checkout-service` as the only consumer and confirms migration to v2 is complete. Generates deprecation closure: marks CartUpdated v1 status as "sunset", closes dual-publish window, creates Contract Spec update entry in the Spec Graph with migration audit trail.

---

### Component Spec

**Component Team — Cart — Picking up a fan-out task**
> *Context: Fan-out task received. Platform Spec OC-124 v1, context pack cp-v2, contract_change: yes.*

Prompt:
```
I have the Cart fan-out task for ECO-124. Platform spec OC-124 v1, context pack
cp-v2. Write the Cart Component Spec. Cart must support guest sessions via token
(no auth dependency), persist through checkout, and emit CartUpdated v2. Check
Cart domain invariants first.
```

Result: Agent reads context pack cp-v2, queries Domain MCP for Cart invariants — finds CART-INV-001 (auth-free sessions), CART-INV-002 (30-min TTL minimum), CART-INV-003 (no cross-domain DB access). Queries Component MCP for cart-service patterns (Redis for sessions, transactional outbox for events). Generates Component Spec declaring: `implements: OC-124 v1`, `context_pack: cp-v2`, `contracts_referenced: CartUpdated v2`. Gate Check: all 5 gates pass.

---

**Component Team — Payments — Unblocking a spec after ADR approval**
> *Context: Payments Component Spec was blocked on ADR-219. ADR-219 has just been approved.*

Prompt:
```
ADR-219 is approved — idempotency key is payment_intent_id, stored in DB before
each capture attempt. Unblock the Payments Component Spec and generate the
implementation plan.
```

Result: Agent updates Payments Component Spec: removes `BlockedBy: ADR-219`, status changes from Blocked → Approved. Incorporates the approved idempotency pattern (references ADR-219 as authoritative decision). Gate 5 now passes — no open ADRs. Generates `tasks.md` with dependency-ordered implementation tasks.

---

### Hotfix

**On-Call Engineer — Production payment timeout**
> *Context: Production incident. Payment authorisations timing out. Revenue impacted. Stripe EU degraded.*

Prompt:
```
URGENT: Production is failing. Payment authorisation calls to Stripe are timing out
after 30s. Orders dropping. Stripe status page shows EU region degraded. We route
all traffic there. Need immediate fix.
```

Result: Agent switches to Hotfix Path — creates minimal Hotfix Spec HOTFIX-12: issue = Stripe EU timeout, fix = add regional fallback routing to Stripe US with 15s timeout + 2 retries, rollback = revert routing config flag, observability = monitor `payment_authorization.timeout_rate`. Marks follow-up spec required for full circuit breaker implementation with ADR for resilience strategy.

---

**Component Team — Post-hotfix hardening spec**
> *Context: Hotfix HOTFIX-12 is stable and deployed. Time to harden the payments service properly.*

Prompt:
```
The HOTFIX-12 fix is stable. Now write the follow-up hardening spec for payments.
We need a proper circuit breaker for Stripe calls, full observability, and an ADR
documenting the resilience strategy decision.
```

Result: Agent creates Follow-up Spec referencing HOTFIX-12. Queries Component MCP for payments-service approved patterns — finds opossum circuit breaker is registered. Generates full Component Spec: circuit breaker config (5 failures/10s window, 30s half-open recovery), new metrics declared (`circuit_breaker.state`, `payment_provider.timeout_rate`), creates ADR-226 for resilience strategy with options analysed. Gate 4 check verifies observability section is complete.

---

### ADR

**Platform Architect — Resolving a blocking ADR before planning**
> *Context: ADR-220 (guest session expiry) is blocking Cart and Checkout component specs.*

Prompt:
```
Resolve ADR-220 on guest session expiry strategy. Three options: (A) 30-minute
sliding window TTL, (B) 24-hour fixed TTL, (C) TTL tied to checkout session
lifecycle. Consider Cart domain invariants, Redis cost, and abandonment rates.
Draft the ADR with a recommendation.
```

Result: Agent queries Domain MCP — finds CART-INV-002 (sessions must outlive checkout flow, minimum 30-min TTL). Analyses all 3 options: Option B violates no invariant but increases Redis cost 6×. Option C creates coupling between domains (violates boundary rules). Recommends Option A (sliding 30-min window). Generates ADR-220 in Proposed state with decision rationale, trade-offs, and consequence analysis. Marks Cart and Checkout specs as unblocked on approval.

---

**Platform Architect — Settling a cross-domain event versioning debate**
> *Context: Multiple teams are adopting inconsistent strategies for event schema versioning.*

Prompt:
```
Teams keep making different calls on event schema versioning — some use
field-addition only, others create new event types. I need a global ADR that
settles this for all domains once and for all. Consider Kafka compatibility
constraints and our existing consumer patterns.
```

Result: Agent queries Integration MCP and identifies the inconsistency pattern. Creates global ADR-225 scoped to Platform Repo. Documents two options with full Kafka compatibility analysis. Recommends semantic versioning: field additions are backward-compatible (MINOR bump), field removal/rename/type-change are breaking (MAJOR version, requires dual-publish). Marks ADR-225 as blocking all future Contract Specs until approved by Integration Owner.

---

### Gate Check

**AI Agent — Running gate analysis before fan-out**
> *Context: Architect has produced feature-spec.md OC-124 v1. Running Verifier gate check before fan-out.*

Prompt:
```
Run full gate analysis on OC-124 v1. I need a Gate 1–5 check before we fan out
to component teams. Show each gate result with specific evidence.
```

Result:
- G1 PASS — all MCP sources cited with versions, cp-v2 pinned
- G2 PASS — no invariant violations, domain ownership respected across Cart/Checkout/Payments/Fulfillment
- G3 FAIL — OrderPlaced v3 consumer list incomplete (shipping-service missing)
- G4 PASS — logging/metrics/tracing declared, PII masking plan for `guest_email` documented citing SEC-001
- G5 FAIL — `BlockedBy: ADR-219` still in Proposed state (no approved owner)

Fan-out BLOCKED. Agent outputs exact remediation steps for both failing gates.

---

**Component Team — Checking a local design decision against domain invariants**
> *Context: Cart engineer wants to store checkout session state in Cart's Redis instance.*

Prompt:
```
We're thinking of storing checkout session state in the Cart Redis — it's already
there for cart sessions and would simplify the architecture. Is this okay from a
domain perspective?
```

Result: Agent queries Domain MCP and immediately flags two invariant violations: CART-INV-001 (Cart owns cart session state — Checkout must not store or mutate it) and CHK-INV-001 (Checkout owns order lifecycle — its session state must live in Checkout's own storage). Returns Gate 2 FAIL with specific invariant IDs and rationale. Suggests correct pattern: Checkout creates its own Redis namespace (`checkout-sessions:*`) or Postgres table, referencing the cart via `cart_token` only.

---

## 07 — Roles & Responsibilities

Every role has a clear repo, clear ownership, and clear accountability. No overlap — no ambiguity about who decides what.

| Role | Repo | Owns | Accountable For |
|---|---|---|---|
| **Product Manager** | Platform Repo | Initiative definition, business goals, UX intent | Initiative (Epic), success criteria |
| **Platform Architect** | Specs Repo | `constitution/`, `templates/`, `workflows/`, `graph/index.yaml` — governed inputs that Platform MCP serves | "What the system must do" across all components · approves Architect agent output |
| **Domain Owner** | Specs + System Repo | `components/{name}/specs/` — invariants, contracts, patterns, ADRs served by Component MCP | Domain correctness of all specs touching their component · approves Developer agent output |
| **Integration Owner** | Specs Repo | `contracts/` inside each component spec — `get_contracts()` on Component MCP returns these | Approves contract changes declared in Architect exit gate · verifies no breaking changes |
| **Component Team** | System Repo | Implementation specs (`docs/specs/`), code, tests, local ADRs | "How the system works locally" — Developer + Verifier agents |
| **ADR Owner** | Platform or Component | Technical decisions — global or local scope | Resolving ambiguity before implementation proceeds |
| **AI Agents** | analyst · architect · developer · verifier | Discovery, specification, implementation, and verification — loaded from `.claude/agents/` as Markdown | Must call Platform + Component MCPs · self-check exit gates · produce traceable outputs per risk level |

---

## 08 — Scaling Guide

Start with one repo and add structure only when the team needs it. The two-MCP model and four agents never change — only the physical organisation of content does.

### 🟡 Early Stage — 1–2 people

One repo. Platform MCP reads a single `openclaw-specs/` directory. Component MCP reads the one monorepo component you have. Quick workflow only.

- Single repo: specs + code together
- Platform MCP reads local directory
- One component MCP for the whole system
- Quick workflow only (developer → verifier)

### 🔵 Growing Team — 3–6 people

Split repos. `openclaw-specs/` becomes its own repo so product and engineering can PR separately. One component MCP instance per bounded context.

- `openclaw-specs/` repo split from system repo
- Component MCP instances per major component
- Standard workflow introduced for new features
- `agents/` symlinked into `.claude/` from specs repo

### 🔴 Full Model — 6+ people

Full two-repo model as described in this document. Full workflow for high/critical. Human approval gates enforced at Verifier. Spec graph fully linked.

- Full workflow for high/critical risk
- Analyst agent activated for major initiatives
- Human approval gates enforced in `verifier.md`
- `graph/index.yaml` fully linked for audit trail

---

## Day 1 Checklist

Four phases. MCP servers first — then your first initiative — then parallel fan-out — then harden.

### Phase 1 · Foundation
- [ ] Create `openclaw-specs/` repo
- [ ] Write `constitution/policies.md`
- [ ] Deploy Platform MCP (reads `specs/`)
- [ ] Deploy Component MCP per component
- [ ] Symlink `.claude/agents/` to `specs/agents/`

### Phase 2 · First Initiative
- [ ] Select one low-risk pilot feature
- [ ] Run quick workflow: `developer.md` only
- [ ] Developer calls Component MCP
- [ ] Verifier produces `verify.md`
- [ ] Validate gate checklist end-to-end

### Phase 3 · Standard Flow
- [ ] Run standard workflow for a medium feature
- [ ] Architect produces `feature-spec.md`
- [ ] Fan out: `developer.md` per component (parallel)
- [ ] Verifier waits for all developers done
- [ ] Update `graph/index.yaml` post-merge

### Phase 4 · Harden
- [ ] Add full workflow for critical changes
- [ ] Analyst activated for Guaripolo/live flag
- [ ] Human approval gate in `verifier.md`
- [ ] CI enforces `verify.md` sign-off before merge

---

*Target Operating Model · SDD + MCP · v3.0*
*"Software is no longer just built — it is specified, validated, and executed as a system of knowledge."*
*February 2026*
