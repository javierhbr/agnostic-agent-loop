# SDD Methodology — Unified 5-Phase Model

Spec-Driven Development (SDD) is a methodology for shipping changes with
traceability, clear ownership, and repeatable quality gates. This page explains
the current unified model, its phases, the artifacts it produces, and how teams
adopt it.

---

## Why SDD Exists

Most engineering teams lose context between "someone had an idea" and "the code
is in production." Requirements live in chat threads, architecture decisions
live in someone's head, and verification is ad-hoc. SDD fixes this by making
the **change package** the canonical unit of work: a self-contained folder of
artifacts that travels through five phases, each owned by a named role, until
the change is deployed and archived.

---

## The Five Phases

```text
Phase 1: Platform    (Architect)   -- Shared context, constitution, ownership
Phase 2: Assess      (Team Lead)   -- Classify request, open change package
Phase 3: Specify     (Product)     -- Write proposal, delta specs, glossary
Phase 4: Plan        (Architect)   -- Design, tasks, story mapping
Phase 5: Deliver     (Team Lead)   -- Build, PR, review, verify, deploy, archive
```

Each phase produces specific artifacts and ends with a gate check before the
next phase can begin.

### Phase 1 — Platform (Architect owns)

Create durable shared context that outlives any single change. This phase runs
once when the platform is set up and is revisited whenever the architecture
evolves.

**Artifacts produced:**

- `constitution.md` — Guiding principles, tech-radar decisions, non-negotiables
- `component-ownership-<name>.md` — One file per bounded context, lists owning
  team, public API surface, and allowed consumers
- `dependency-map.md` — Three-tier dependency map (Tier 1 critical, Tier 2
  standard, Tier 3 convenience) with upgrade policy per tier
- `glossary.md` — Shared vocabulary across all teams; every term used in specs
  must appear here

**Tooling boundary:** Platform-side repos use BMAD + OpenSpec + Speckit.
Component repos use OpenSpec only. This is the core boundary rule — it keeps
component teams focused on specs and delivery without platform tooling overhead.

### Phase 2 — Assess (Team Lead owns)

Classify the incoming request and open a change package.

**Steps:**

1. Classify the request (feature, bug, hotfix, tech-debt, spike)
2. Create the change package directory under `.agentic/changes/<change-id>/`
3. Populate `platform-ref.yaml` — links to the platform constitution, relevant
   component ownership files, and glossary version
4. Populate `jira-traceability.yaml` — maps the change to external tracking
   (epic, story, ticket IDs)

**Artifacts produced:**

- `platform-ref.yaml` — Platform context reference for this change
- `jira-traceability.yaml` — External tracker cross-references

### Phase 3 — Specify (Product owns)

Write the proposal and delta specs that describe what changes.

**Steps:**

1. Write `proposal.md` — Problem statement, success criteria, scope boundaries
2. Write delta specs — Each spec is tagged ADDED, MODIFIED, or REMOVED to
   describe the change to existing components
3. Cross-check every domain term against `glossary.md`; add missing terms

**Artifacts produced:**

- `proposal.md` — The "what and why" of the change
- Delta spec files — One per affected component, tagged with change type

### Phase 4 — Plan (Architect owns)

Turn the specification into an executable plan.

**Steps:**

1. Write `design.md` — Technical approach, component interactions, risk
   mitigations, ADR references
2. Write `tasks.md` — Ordered task list with effort estimates, dependencies,
   and acceptance criteria
3. Map tasks to stories in the external tracker (updating
   `jira-traceability.yaml`)

**Artifacts produced:**

- `design.md` — Technical design document
- `tasks.md` — Deliverable task breakdown

### Phase 5 — Deliver (Team Lead owns)

Execute the plan through build, review, verification, and deployment.

**Steps:**

1. **Build** — Claim tasks, develop in isolated worktrees
2. **Create PR** — Generate pull request from spec context
3. **Review PR** — Independent reviewer in separate worktree
4. **Verify** — Run validators, gate checks, test suites
5. **Deploy** — Feature-flagged rollout (default off, progressive enablement)
6. **Archive** — Move change package to completed, capture learnings

**CLI commands used:**

```bash
agentic-agent task claim TASK-ID
agentic-agent task complete TASK-ID
agentic-agent pr create --task TASK-ID
agentic-agent pr review --task TASK-ID --pr-url URL
agentic-agent validate
agentic-agent task list
```

---

## The Change Package

Every change flows through the five phases inside a single directory:

```
.agentic/changes/<change-id>/
  platform-ref.yaml          # Phase 2 — platform context links
  jira-traceability.yaml     # Phase 2 — external tracker mapping
  proposal.md                # Phase 3 — what and why
  delta-specs/               # Phase 3 — ADDED/MODIFIED/REMOVED specs
    component-auth.md
    component-payments.md
  design.md                  # Phase 4 — technical approach
  tasks.md                   # Phase 4 — task breakdown
```

The change package is the single source of truth for a unit of work. It
replaces the older pattern of scattered `feature-spec.md`, `impl-spec.md`, and
`verify.md` files.

---

## DDD Ownership Artifacts

SDD uses three artifacts to encode domain-driven ownership at the platform
level:

1. **Component Ownership** (`component-ownership-<name>.md`) — Defines who owns
   each bounded context, what its public API is, and which other components may
   depend on it.

2. **Dependency Map** (`dependency-map.md`) — Classifies every dependency into
   three tiers:
   - Tier 1 (Critical) — Core business logic dependencies; strict upgrade policy
   - Tier 2 (Standard) — Infrastructure and framework dependencies; regular cadence
   - Tier 3 (Convenience) — Utilities and helpers; upgrade at will

3. **Shared Glossary** (`glossary.md`) — Canonical definitions of domain terms.
   Every term used in a proposal or delta spec must have an entry here. This
   prevents the "same word, different meaning" problem across teams.

---

## Traceability

Two YAML files maintain traceability for every change:

- **`platform-ref.yaml`** — Points to the constitution version, relevant
  ownership files, and glossary snapshot used when the change was specified.
  This ensures the change was authored against a known platform baseline.

- **`jira-traceability.yaml`** — Maps the change package to external systems:
  epic IDs, story IDs, PR numbers, deployment tickets. This closes the loop
  between spec artifacts and project management tools.

---

## Adopting SDD — Two-Iteration Approach

Teams adopt SDD incrementally over two iterations:

### Iteration 1 — Platform + Assess + Specify (Phases 1-3)

Focus on getting the foundation right. Set up the constitution and ownership
artifacts. Start classifying requests and writing proposals with delta specs.
This iteration builds the habit of "specify before you build" without changing
how teams plan or deliver.

### Iteration 2 — Plan + Deliver (Phases 4-5)

Once specification is routine, add structured planning with `design.md` and
`tasks.md`, and adopt the full delivery workflow with worktrees, PR generation,
independent review, and validation gates.

This two-step adoption avoids the "big bang" rollout failure mode.

---

## SDD Skills

SDD includes 25+ skills organized by role and phase. Install all skills with:

```bash
agentic-agent skills ensure
```

Skills are grouped into:

| Category      | Skills                                                      | Used In |
| ------------- | ----------------------------------------------------------- | ------- |
| Platform      | platform-constitution, platform-spec, component-spec        | Phase 1 |
| Assessment    | workflow-router, risk-assessment                            | Phase 2 |
| Specification | openspec, product-wizard, brainstorming                     | Phase 3 |
| Planning      | architect, dev-plans, adr                                   | Phase 4 |
| Delivery      | developer, tdd, verifier, gate-check, hotfix                | Phase 5 |
| Cross-cutting | process-guide, superpowers-bridge, diataxis, agentic-helper | All     |

The full skill directory is maintained in the skills packs under
`internal/skills/packs/`. Each skill follows the three-tier context model:
router (Tier 1) points to slim skill files (Tier 2) which reference detailed
resources (Tier 3) on demand.

---

## By Role

### Product Manager

- Owns Phase 3 (Specify)
- Writes `proposal.md` and delta specs
- Uses: openspec, product-wizard, brainstorming skills
- Reference: Process guide Phase 3 steps

### Architect

- Owns Phase 1 (Platform) and Phase 4 (Plan)
- Maintains constitution, ownership artifacts, dependency map
- Writes `design.md` and `tasks.md`
- Uses: platform-constitution, platform-spec, architect, adr, dev-plans skills

### Team Lead

- Owns Phase 2 (Assess) and Phase 5 (Deliver)
- Opens change packages, coordinates delivery
- Uses: workflow-router, risk-assessment, gate-check skills

### Developer

- Executes tasks in Phase 5 (Deliver)
- Works in isolated worktrees, creates PRs
- Uses: developer, tdd, component-spec skills

### Reviewer / Verifier

- Reviews PRs and runs verification in Phase 5
- Independent worktree, different perspective than builder
- Uses: verifier, gate-check skills

---

## Getting Started

```bash
# 1. Install SDD skills
agentic-agent skills ensure

# 2. View available tasks
agentic-agent task list

# 3. Claim and start work
agentic-agent task claim TASK-ID

# 4. Validate before completing
agentic-agent validate

# 5. Complete the task
agentic-agent task complete TASK-ID

# 6. Create a PR from spec context
agentic-agent pr create --task TASK-ID
```

---

## Further Reading

- [OPERATING-MODEL.md](sdd-mcp/operation%20model/OPERATING-MODEL.md) — Full SDD specification
- [SDD-SKILLS-DIRECTORY.md](SDD-SKILLS-DIRECTORY.md) — Complete skill reference
- [sdd-example-workflow.md](sdd-example-workflow.md) — End-to-end walkthrough
- [USING-SDD-PROCESS-GUIDE.md](USING-SDD-PROCESS-GUIDE.md) — Practical usage tips
- `docs/README-LAYERED-CONTEXT.md` — Layered context model documentation

---

## Core Principle

Every change flows through five phases with clear ownership at each step:
Architect sets context, Team Lead opens the package, Product specifies the
change, Architect plans the execution, and Team Lead drives delivery. The
change package is the thread that connects intent to production.
