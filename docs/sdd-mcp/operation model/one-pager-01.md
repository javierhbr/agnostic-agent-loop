# One Pager
## Spec-Driven Development (SDD) v4.0
**Five Phases · Governed Knowledge · Zero Overhead Services**

> *Version 4.0 — March 2026*
> *"Nothing is implemented without a validated spec backed by governed context."*

---

## Problem

Modern distributed platforms have a **knowledge problem**, not just a code problem:

- Teams interpret the same requirement differently → misalignment and rework
- Service contracts break silently → costly integration bugs
- Decisions live in Slack and people's heads → no audit trail, tribal knowledge
- AI agents invent context from the codebase → unpredictable, unsafe outputs
- Onboarding requires human explanation → doesn't scale

---

## Solution

Adopt **SDD v4.0** with a 5-phase model, role-based ownership, and zero overhead infrastructure:

| Component | What it does |
|---|---|
| **5-Phase Model** | Platform --> Assess --> Specify --> Plan --> Deliver |
| **Role Assignments** | Architect (Platform, Plan) · Team Lead (Assess, Deliver) · Product (Specify) |
| **Change Package** | Canonical execution unit: versioned directory containing proposal, design, tasks, and evidence |
| **Spec Graph** | Traceable chain: Initiative --> proposal --> design --> tasks --> verify.md --> ADRs |
| **Ownership Artifacts** | `component-ownership.yaml` · `dependency-map.yaml` · `shared-glossary.yaml` |

---

## Core Principle

> *"Nothing is implemented without a validated spec backed by governed context."*

Agents — human or AI — do not invent context. They **call MCP tools** to get governed truth.

---

## Core Boundary Rule

The 5-phase model applies differently depending on repository scope.

| Scope | Phases Available | Methodology Sources |
|---|---|---|
| **Platform-side** (specs repo, constitution, governance) | All 5 phases | BMAD + OpenSpec + SpecKit |
| **Component repos** (system code, local specs) | Assess, Specify (delta only), Plan, Deliver | OpenSpec only |

**Rule:** Component teams consume platform artifacts (constitution, templates, ownership, glossary) as read-only input. They contribute to them via PR to the platform repo, never by direct edit.

---

## Methodology Origins

SDD v4.0 is a deliberate composition of proven ideas. Three frameworks provided the structural core -- every practice listed here was adopted because it directly changes what an agent or engineer produces. Everything else was cut.

**BMAD — Breakthrough Method for Agile AI Development**
The four-agent architecture (analyst → architect → developer → verifier) comes directly from BMAD. Each agent has a single bounded responsibility, produces one or two files, and enforces a self-checked exit gate before handing off. The parallel fan-out pattern — one Developer agent per component running simultaneously after the Architect completes the feature-spec — is also a BMAD practice. *Why we adopted it: role-scoped agents prevent context overreach. Parallel fan-out eliminates the bottleneck that makes large cross-domain features slow.*

**OpenSpec — Open Specification Format**
The traceability metadata on every spec (`implements`, `context_pack`, `blocked_by`, `status`) comes from OpenSpec. So does the Spec Graph (`graph/index.yaml`) — the machine-readable index that links every artifact from Initiative through verify.md, updated by the Verifier after every merge. The rule that specs are never deleted, only versioned or paused, is also OpenSpec. *Why we adopted it: without mandatory metadata, a spec is an island. The Spec Graph is the audit trail — it answers "what was built, why, against what context, and was it verified?"*

**SpecKit — Template-Driven Spec Assembly**
The template-with-embedded-exit-gate pattern comes from SpecKit. Rather than writing a spec and separately checking it, the checklist lives inside the template — the agent writes and self-verifies in one pass. The Context Pack mechanism — a pinned, versioned snapshot of applicable policies, NFR baselines, and workflow config assembled before any spec is written — is also a SpecKit practice. Templates are served by Platform MCP (`get_template(name)`) so they evolve without agent file changes. *Why we adopted it: embedding the gate in the template closes the gap between writing and verifying. Pinning the context pack makes specs reproducible — another agent can retrieve the same context and understand what constraints were in effect.*

---

## System Layers

```mermaid
flowchart TD
    KL["KNOWLEDGE LAYER\nPlatform MCP · Component MCPs\nPolicies · Contracts · Invariants · Patterns · ADRs"]
    CP["Context Pack — versioned snapshot\nassembled before any spec is written"]
    SL["SPEC LAYER\nfeature-spec · component-spec · impl-spec · ADRs · verify.md"]
    CL["CODE LAYER\nServices · APIs · Infrastructure"]
    OL["OBSERVABILITY\nLogs · Metrics · Traces"]

    KL -->|"get_context_pack()"| CP
    CP --> SL
    SL -->|"implementation against gate-checked specs"| CL
    CL --> OL

    style KL fill:#fff3dc,stroke:#c47d0e,color:#1a1c2e
    style CP fill:#fdf8ee,stroke:#c47d0e,color:#1a1c2e
    style SL fill:#eef0fa,stroke:#1e2d6b,color:#1a1c2e
    style CL fill:#eef7f5,stroke:#1a7a6d,color:#1a1c2e
    style OL fill:#f0faf4,stroke:#2a7d4f,color:#1a1c2e
```

```
┌─────────────────────────────────────────────────────┐
│              KNOWLEDGE LAYER  (MCP Servers)         │
│   Platform MCP · Component MCPs (×N)               │
│   Policies · Contracts · Invariants · Patterns      │
└─────────────────────┬───────────────────────────────┘
                      │  get_context_pack()
                      │  Context Pack (versioned snapshot)
                      ▼
┌─────────────────────────────────────────────────────┐
│                   SPEC LAYER  (2 repos)             │
│   feature-spec · component-spec · impl-spec         │
│   ADRs · verify.md                                  │
└─────────────────────┬───────────────────────────────┘
                      │  implementation against
                      │  gate-checked specs
                      ▼
┌─────────────────────────────────────────────────────┐
│                   CODE LAYER                        │
│          Services · APIs · Infrastructure           │
└─────────────────────┬───────────────────────────────┘
                      ▼
┌─────────────────────────────────────────────────────┐
│                  OBSERVABILITY                      │
│              Logs · Metrics · Traces                │
└─────────────────────────────────────────────────────┘
```

---

## 5-Phase Model

```
  Phase 1        Phase 2        Phase 3        Phase 4        Phase 5
  PLATFORM  -->  ASSESS    -->  SPECIFY   -->  PLAN      -->  DELIVER
  Architect      Team Lead      Product        Architect      Team Lead
  (durable)      (per change)   (per change)   (per change)   (per change)
```

Each phase has a designated owner, specific inputs and outputs, and an exit gate that must pass before the next phase begins.

### Change Package

Every unit of work is a **change package** -- a versioned directory opened in Assess and archived in Deliver:

```
changes/CP-{ID}/
  platform-ref.yaml        # pinned constitution version
  jira-traceability.yaml   # external tracking links
  proposal.md              # what and why (Specify output)
  design.md                # how (Plan output)
  tasks.md                 # execution plan (Plan output)
  verify.md                # evidence (Deliver output)
```

---

## Two MCPs, Four Agents

```mermaid
graph LR
    subgraph PMCP["Platform MCP — 1 instance"]
        P1["get_context_pack(intent)"]
        P2["get_template(name)"]
        P3["get_agent(name)"]
        P4["get_workflow(risk_level)"]
    end

    subgraph CMCP["Component MCP — N instances"]
        C1["get_contracts()"]
        C2["get_invariants()"]
        C3["get_patterns()"]
        C4["get_decisions()"]
    end

    subgraph AGENTS["4 Agents — .claude/agents/"]
        A1["analyst.md\nfull only"]
        A2["architect.md\nstandard + full"]
        A3["developer.md\nall · parallel"]
        A4["verifier.md\nall · hard stop"]
    end

    P1 --> A2
    P2 --> A2
    P3 --> A1
    P4 --> A2
    C1 --> A2
    C2 --> A2
    C3 --> A3
    C4 --> A3
    A1 --> A2
    A2 --> A3
    A3 --> A4
```

```
┌─────────────────────────────┐   ┌──────────────────────────────┐
│   Platform MCP  (1 instance)│   │  Component MCP  (N instances)│
│                             │   │                              │
│  get_context_pack(intent)───┼─┐ │  get_contracts() ───────────┼─┐
│  get_template(name) ────────┼─┤ │  get_invariants()───────────┼─┤
│  get_agent(name)────────────┼─┤ │  get_patterns() ────────────┼─┼──┐
│  get_workflow(risk_level)───┼─┘ │  get_decisions()────────────┼─┼──┤
└─────────────────────────────┘   └──────────────────────────────┘  │
          │                                  │                       │
          ▼                                  ▼                       │
┌─────────────────────────────────────────────────────────────────┐ │
│                  4 Agents  (.claude/agents/)                    │ │
│                                                                 │ │
│  analyst.md ──► architect.md ──► developer.md ──► verifier.md  │ │
│  (full only)    (std + full)     (all · parallel) (all · stop) │◄┘
└─────────────────────────────────────────────────────────────────┘
```

---

## Three Workflows

```mermaid
flowchart LR
    RL{Risk Level?}

    RL -->|"low · bug_fix"| Q["QUICK\nDeveloper --> Verifier"]
    RL -->|medium| S["STANDARD\nArchitect --> Developer --> Verifier"]
    RL -->|"high · critical"| F["FULL\nAnalyst --> Architect\n--> Developer --> Verifier\n+ Human Approval"]

    style Q fill:#fff3dc,stroke:#c47d0e,color:#1a1c2e
    style S fill:#eef0fa,stroke:#1e2d6b,color:#1a1c2e
    style F fill:#fdf0f0,stroke:#b93030,color:#1a1c2e
    style RL fill:#f4f1ec,stroke:#888,color:#1a1c2e
```

```
            ┌─────────────────────┐
            │   Risk Level?       │
            │  (Platform MCP      │
            │   decides)          │
            └──┬──────┬───────────┘
               │      │           │
         low / │      │ medium    │ high / critical
        bug_fix│      │           │
               ▼      ▼           ▼
  ┌────────────────┐ ┌──────────────────────┐ ┌────────────────────────────┐
  │  QUICK         │ │  STANDARD            │ │  FULL                      │
  │                │ │                      │ │                            │
  │  Developer     │ │  Architect           │ │  Analyst                   │
  │     │          │ │     │                │ │     │                      │
  │  Verifier      │ │  Developer (║ parallel) │  Architect               │
  │                │ │     │                │ │     │                      │
  │                │ │  Verifier            │ │  Developer (║ parallel)   │
  └────────────────┘ └──────────────────────┘ │     │                      │
                                              │  Verifier + Human Gate     │
                                              └────────────────────────────┘
```

Developers fan out **in parallel per component** once `feature-spec.md` is ready.

---

## Roles

| Role | Phase Ownership | Responsibility |
|---|---|---|
| **Architect** | Platform, Plan | Constitution, templates, workflows, Spec Graph, design.md, tasks.md, gate checks |
| **Team Lead** | Assess, Deliver | Opens change package, classifies risk, coordinates build/verify/archive |
| **Product** | Specify | Proposal.md, acceptance criteria, delta specs, glossary check |
| **Domain Owner** | -- | Component invariants, contracts, patterns, ADRs |
| **Component Team** | -- | Implementation, code, tests, local ADRs |

---

## Comparison

| Without SDD | With SDD v4.0 |
|---|---|
| Implicit, undocumented specs | Explicit versioned specs with ACs |
| Integration bugs in QA or prod | Gates block bugs before implementation |
| Decisions in Slack | Every decision is an ADR with audit trail |
| AI agents invent context | Agents consume governed MCP context |
| Onboarding requires tribal knowledge | Context Packs enable fast, safe onboarding |

---

## Benefits

- Cross-domain consistency — all agents work from the same governed truth
- Fewer integration bugs — contracts versioned, gate-checked before any code is written
- Full traceability — Initiative → feature-spec → component-specs → verify.md → ADRs
- Scales from 1 person to 50+ — same two MCP types and four agents throughout
- AI-agent ready — structured context eliminates hallucinated decisions
- Compound knowledge — each initiative enriches the MCP for the next one
- Zero infrastructure overhead — all process lives in Markdown, not servers

---

## Final Insight

> *"Software is no longer just built — it is specified, validated, and executed as a system of knowledge."*

---

*Target Operating Model · SDD + MCP · v4.0 — March 2026*
