---
name: platform-spec
description: "Guides Platform Architects through creating a complete Platform Spec using SpecKit + MCP. Trigger when creating a new Platform Feature Spec, writing Constitution, or starting a new initiative."
---

# skill:platform-spec

## Does exactly this

Guides Platform Architects through producing a complete, MCP-grounded Platform Spec that defines the "what + why + UX" for a cross-domain initiative.

---

## When to use

- Starting a new platform initiative from scratch
- Need to create a Platform Feature Spec
- Writing or updating the Constitution
- Unsure what goes in a Platform Spec vs Component Spec

---

## Key Distinction

**Platform Spec:** "What + why + UX" for the initiative (cross-domain, owned by PM + Architect, lives in Platform Repo)
**Component Spec:** "How to implement" the Platform Spec in one service (owned by Component Team, lives in Component Repo)

---

## Before You Start — Prerequisites

Confirm each item is in place:

- [ ] You are in the Platform Repo (not a component repo)
- [ ] `specify init` has been run: `specify init platform-repo --ai claude`
- [ ] `specify check` passes
- [ ] You have an Initiative ID (e.g., ECO-124)
- [ ] MCP Router generated a Context Pack: `.specify/memory/context-<initiative-id>.md`
- [ ] `constitution.md` exists at `.specify/memory/constitution.md`

**If missing:**
- Constitution? → Run `Step 0: Constitution` first (see constitution skill)
- Context Pack? → Run MCP Router first (see `resources/mcp-router-guide.md`)

---

## Step-by-Step Workflow

### Step 1 — Establish Problem & Success Criteria

**Input:** Initiative ID, business context, user problem

**Define:**
- Problem statement (what user pain exists)
- Success metric (how we'll know it worked)
- Acceptance criteria in Given/When/Then format

See `resources/platform-spec-template.md` for template.

### Step 2 — Map Domain Model

**Input:** Context Pack (from MCP Router), Domain MCPs

**Define:**
- Which domains are involved
- Domain invariants that apply
- Cross-domain interactions
- Data ownership

### Step 3 — Design UX Flow

**Define:**
- User journey (step-by-step)
- Key screens or API endpoints
- Error handling flows
- Accessibility requirements (per Constitution)

### Step 4 — Specify Contracts

**Define:**
- Events this initiative emits (versioned)
- Events this initiative consumes (versioned)
- API contracts (versioned)
- Consumer/producer relationships

### Step 5 — Declare NFRs

**Define:**
- Logging requirements (per Constitution)
- Metrics and monitoring (per Constitution)
- Tracing requirements (per Constitution)
- PII handling (per Constitution)
- Performance targets (per Constitution)

### Step 6 — Gate Check

Run `/speckit.analyze` to verify all 5 gates pass. See `resources/gate-check-guide.md` if any fail.

---

## Fan-Out to Component Teams

When Platform Spec is approved:
1. Create Component Specs (1 per service implementing this feature)
2. Send fan-out task to each Component Team with:
   - Platform Spec ID + version
   - Context Pack version
   - Component ID (which service to implement)
   - Contracts this component must produce/consume

Component Teams will create Component Specs that implement your Platform Spec.

---

## Key Principles

- **Lead with problem** — "Why does this matter now?"
- **Declare all dependencies upfront** — Architecture, integrations, risks
- **Let component teams handle "how"** — Spec defines "what + why", not implementation details
- **Every spec must pass all 5 gates** — No exceptions

---

## If you need more detail

→ `resources/platform-spec-template.md` — Full template with all sections, worked examples, acceptance criteria patterns
→ `resources/mcp-router-guide.md` — How to run MCP Router and generate Context Packs
→ `resources/gate-check-guide.md` — How to debug failing gates
