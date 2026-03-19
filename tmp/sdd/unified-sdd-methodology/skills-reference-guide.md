# SDD Skills Reference Guide

> Complete reference for all SDD skills: when to use them, who uses them, what they produce, and how they connect.

---

## SDD Phases at a Glance

```
Phase 0        Phase 1        Phase 2        Phase 3        Phase 4
─────────      ─────────      ─────────      ─────────      ─────────
Initiative     Discovery      Architecture   Implement      Verify
& Risk         (Full only)    & Spec         & Plan

initiative-    analyst        architect      developer      verifier
definition                    component-spec tdd            gate-check
risk-
assessment
workflow-
router
                              platform-      sdd-openspec
                              constitution   sdd-bmad
                              platform-      sdd-speckit
                              contextualizer
                              adr (blocks)
```

```
Off-Path (any time):  hotfix
Cross-cutting:        explain-code-skill · local-doc · stakeholder-communication · superpowers-bridge
Orchestrators:        process-guide · sdd-methodology-tools · uncle-sdd-agent · unified-sdd-skill
Meta:                 tier-enforcer
```

---

## Quick Lookup Table

| Skill | Phase | Role | Key Artifact Produced | Skill File |
|---|---|---|---|---|
| `initiative-definition` | Phase 0 | Product Manager | initiative.md | `sdd/initiative-definition/SKILL.md` |
| `risk-assessment` | Phase 0 | Product Manager | Risk section in initiative.md | `sdd/risk-assessment/SKILL.md` |
| `workflow-router` | Phase 0 (routing) | PM / Team Lead | Workflow routing guidance | `sdd/workflow-router/SKILL.md` |
| `analyst` | Phase 1 (Full only) | Business Analyst | discovery.md | `sdd/analyst/SKILL.md` |
| `architect` | Phase 2 | Solution Architect | feature-spec.md, component-spec.md | `sdd/architect/SKILL.md` |
| `component-spec` | Phase 2 | Component Team Lead | Component Spec file | `sdd/component-spec/SKILL.md` |
| `platform-constitution` | Phase 2 (Platform) | Platform Architect | constitution.md | `sdd/platform-constitution/SKILL.md` |
| `platform-contextualizer` | Phase 2 (Platform) | Platform Architect | Platform baseline snapshot | `sdd/platform-contextualizer-skill/SKILL.md` |
| `adr` | Phase 2 (gate) | Architect / Tech Lead | ADR-<n>-<title>.md | `sdd/adr/SKILL.md` |
| `developer` | Phase 3 | Developer | impl-spec.md, tasks.yaml | `sdd/developer/SKILL.md` |
| `tdd` | Phase 3 | Developer | Test suite, refactored code | `sdd/tdd/SKILL.md` |
| `sdd-openspec` | Phase 3 | Architect / Developer | proposal.md, design.md, tasks.md | `sdd/sdd-openspec/SKILL.md` |
| `sdd-bmad` | Foundation | PM / Architect / Dev | Tech-spec, PRD, Architecture | `sdd/sdd-bmad/SKILL.md` |
| `sdd-speckit` | Foundation | PM / Architect / Dev | Constitution, specs, tasks | `sdd/sdd-speckit/SKILL.md` |
| `verifier` | Phase 4 | QA / Verification Eng. | verify.md | `sdd/verifier/SKILL.md` |
| `gate-check` | All phase boundaries | Any role | Gate remediation, updated spec | `sdd/gate-check/SKILL.md` |
| `hotfix` | Off-path | Developer / Team Lead | Hotfix Spec, Follow-up Spec | `sdd/hotfix/SKILL.md` |
| `explain-code-skill` | Cross-cutting | Any role | Explanation docs, ASCII diagrams | `sdd/explain-code-skill/SKILL.md` |
| `local-doc` | Cross-cutting | Any role | Search results, spec references | `sdd/local-doc/SKILL.md` |
| `stakeholder-communication` | Cross-cutting | Product Manager | Status emails, progress reports | `sdd/stakeholder-communication/SKILL.md` |
| `superpowers-bridge` | Cross-cutting | Dev / Architect | Tool selection guidance | `sdd/superpowers-bridge/SKILL.md` |
| `process-guide` | Orchestrator | Any role | Phase-by-phase guidance | `sdd/process-guide/SKILL.md` |
| `sdd-methodology-tools` | Orchestrator (router) | Any role | Skill routing | `sdd/sdd-methodology-tools/SKILLS.md` |
| `uncle-sdd-agent` | Orchestrator | Any role | Change packages, ownership artifacts | `sdd/uncle-sdd-agent/SKILL.md` |
| `unified-sdd-skill` | Orchestrator | Any role | Routed proposals, specs, delivery | `sdd/unified-sdd-skill/SKILL.md` |
| `tier-enforcer` | Meta | Skill Maintainer | Compliant skill files | `sdd/tier-enforcer/SKILL.md` |

---

## Phase 0 — Initiative & Risk

### `initiative-definition`

**Description:** Define a new SDD Initiative (Epic) with a problem statement and measurable success criteria before any engineering work begins.

**When to use:**
- Starting a new product initiative or epic
- A PM needs to document business goals and non-goals
- Before handing off to engineers or risk assessment

**Role:** Product Manager

**Methodology step:** Phase 0 — Initiative definition

**Key actions:**
1. Write a measurable problem statement
2. List business goals and explicit non-goals
3. Write success criteria in Given/When/Then format
4. List affected components

**Artifacts consumed:** Business requirements, user feedback
**Artifacts produced:** `initiative.md` (problem statement, goals, non-goals, success criteria, affected components)

**Related skills:** `risk-assessment` (next step), `analyst` (Full workflow)

---

### `risk-assessment`

**Description:** Interview-based risk classification for PMs to determine the correct SDD workflow before engineering begins.

**When to use:**
- After writing the initiative definition
- Before handing off to engineers
- Unsure which workflow (Quick / Standard / Full) to use

**Role:** Product Manager

**Methodology step:** Phase 0.5 — Risk classification

**Key actions:**
1. Answer 5 questions: contract impact, data sensitivity, team scope, rollback strategy, ADR dependency
2. Score answers (0–5 YES responses)
3. Classify risk: Low / Medium / High / Critical
4. Document classification with rationale

**Artifacts consumed:** `initiative.md`
**Artifacts produced:** Risk Assessment section appended to `initiative.md` with classification and rationale

**Related skills:** `workflow-router` (uses classification to route), `initiative-definition` (prerequisite)

---

### `workflow-router`

**Description:** Route an initiative to the correct SDD workflow based on its risk classification.

**When to use:**
- Starting a new initiative after risk classification
- Unsure which workflow path to follow
- Onboarding a team to their workflow sequence

**Role:** Product Manager / Team Lead

**Methodology step:** Phase 0 (routing) — after risk-assessment

**Key actions:**
1. Read risk classification from initiative.md
2. Route to the appropriate workflow:
   - **QUICK** — Low risk, 1–2 days, minimal spec
   - **STANDARD** — Medium risk, 3–5 days, component spec required
   - **FULL** — High/Critical risk, 5–10 days, full analyst + architect + developer chain
3. Show the agent sequence for the chosen workflow
4. Explain handoff protocol between phases

**Artifacts consumed:** Risk classification in `initiative.md`
**Artifacts produced:** Workflow routing guidance, per-workflow agent sequence

**Related skills:** `risk-assessment` (prerequisite), `analyst` (Full), `architect` (Standard/Full), `developer`

---

## Phase 1 — Discovery (Full workflow only)

### `analyst`

**Description:** Deeply understand a business problem through evidence-based discovery before any architecture or implementation work.

**When to use:**
- Full workflow only (High / Critical risk initiatives)
- Before the architect phase
- When business context or affected components are unclear

**Role:** Business Analyst

**Methodology step:** Phase 1 — Discovery

**Key actions:**
1. Call Platform MCP: `get_context_pack(intent)`
2. Conduct evidence-based one-question-at-a-time interviews
3. Explicitly list all affected components
4. Classify risk (Low / Medium / High / Critical) with evidence
5. Produce `discovery.md`

**Artifacts consumed:** Initiative definition, Platform MCP context
**Artifacts produced:** `discovery.md` — problem statement, evidence, affected components, risk classification, key decisions needed

**Related skills:** `initiative-definition` (prerequisite), `architect` (consumes discovery.md), `risk-assessment`

---

## Phase 2 — Architecture & Specification

### `architect`

**Description:** Design feature specs and component specs grounded in full MCP context. Owns the WHAT and WHY — not the HOW.

**When to use:**
- Standard and Full workflows
- After discovery.md (Full) or initiative.md (Standard) is ready
- Before component teams write their implementation specs

**Role:** Solution Architect

**Methodology step:** Phase 2 — Architecture design

**Key actions:**
1. Read `discovery.md` (Full) or initiative definition (Standard)
2. Call Platform MCP: `get_context_pack(risk_level)`
3. Call Component MCP per service: `get_contracts()`, `get_invariants()`, `get_decisions()`
4. Write `feature-spec.md` (WHAT and WHY)
5. Write `component-spec.md` per affected component
6. Create fan-out tasks per component
7. Self-check all 5 gates before handoff

**Artifacts consumed:** `discovery.md` or initiative definition, Platform/Component MCP responses
**Artifacts produced:** `feature-spec.md`, N × `component-spec.md`, fan-out tasks with metadata

**Related skills:** `analyst` (provides discovery.md), `component-spec` (consumes fan-out tasks), `adr` (may block)

---

### `component-spec`

**Description:** Guide component teams through creating correct, MCP-grounded Component Implementation Specs — the "local how" for each service.

**When to use:**
- Receiving a fan-out task with a `platform_spec_id`
- Writing the implementation approach for a single component
- Creating local ADRs scoped to one service

**Role:** Component Team Lead / Architect

**Methodology step:** Phase 2 — Component-level specification

**Key actions:**
1. Verify required inputs: fan-out task, Platform Spec, Context Pack, domain MCP
2. Create Metadata block with traceability to platform spec
3. Write each section with a `Source:` line referencing the platform spec
4. Run Gate Check (all 5 gates)
5. Update `spec-graph.json` after approval

**Artifacts consumed:** Fan-out task from architect, Platform Spec, Context Pack, domain MCP responses
**Artifacts produced:** Component Spec file in `.specify/memory/specs/`, updated `spec-graph.json`

**Related skills:** `architect` (provides fan-out tasks), `gate-check` (validates), `developer` (consumes this spec)

---

### `platform-constitution`

**Description:** Author and maintain the Constitution — the foundational platform governance document that every spec must reference.

**When to use:**
- Starting a new platform that needs governance
- Updating platform policies (security, observability, performance)
- A gate check fails because a policy is missing or undefined
- Onboarding teams to platform non-negotiables

**Role:** Platform Architect

**Methodology step:** Platform phase — foundational governance

**Key actions:**
1. Define all 6 required sections: UX Rules, Security/PII, Observability, Performance, Domain Governance, Contract Versioning
2. Make every rule specific and testable
3. Get stakeholder sign-off
4. Version the document
5. Declare the version in every spec

**Artifacts consumed:** Platform requirements, stakeholder input
**Artifacts produced:** `constitution.md` (versioned, 6 sections)

**Related skills:** `gate-check` (enforces constitution compliance), `platform-contextualizer` (uses constitution as input), `adr` (may amend constitution)

---

### `platform-contextualizer`

**Description:** Document the current state of an existing platform to create the baseline before starting Iteration 1 of SDD adoption.

**When to use:**
- Beginning Iteration 1 with an already-active platform (brownfield)
- Capturing existing conventions, contracts, and pain points
- Creating the platform baseline before formal SDD phases begin
- Onboarding an AI agent to an existing platform

**Role:** Platform Architect

**Methodology step:** Platform phase — baseline capture (before Phase 1)

**Key actions:**
1. Observe and document current state (facts, not opinions)
2. Separate facts from assumptions explicitly
3. Identify gaps and pain points
4. Draft durable principles and context
5. Validate snapshot with stakeholders
6. Use BMAD → OpenSpec → Speckit in sequence to produce artifacts

**Artifacts consumed:** Existing platform docs, team interviews, conventions
**Artifacts produced:** Current-state snapshot, gap register, platform principles, config notes, versioning model, JIRA hierarchy, MCP usage model

**Related skills:** `sdd-bmad`, `sdd-openspec`, `sdd-speckit` (used in sequence), `platform-constitution` (output feeds into it)

---

### `adr`

**Description:** Create, track, and resolve Architecture Decision Records (ADRs) that gate implementation until a cross-cutting or local decision is resolved.

**When to use:**
- A spec surfaces an unresolved cross-cutting question
- Multiple components are affected by an architectural choice
- Unsure whether a decision is Global (platform-wide) or Local (component-scoped)
- Reviewing specs before fan-out to identify blocking decisions

**Role:** Architect / Tech Lead

**Methodology step:** Phase 2 — Gate-blocking decision management

**Key actions:**
1. Determine ADR scope: Global (platform-wide) vs Local (component-scoped)
2. Draft ADR using the standard template
3. Link ADR to all blocked specs (`BlockedBy` field)
4. Move ADR through states: Proposed → In Review → Approved / Rejected
5. Resolve ADR and unblock specs once approved

**Artifacts consumed:** Specs with `BlockedBy` references
**Artifacts produced:** `ADR-<number>-<short-title>.md`, updated `spec-graph.json`

**Related skills:** `gate-check` (enforces ADR resolution), `component-spec` and `developer` (blocked until ADR resolved), `architect`

---

## Phase 3 — Implementation

### `developer`

**Description:** Read the component spec and produce an implementation spec plus task decomposition. Owns the HOW.

**When to use:**
- All workflows (Quick, Standard, Full) — the implementation planning step
- After receiving an approved component spec
- Before writing any code

**Role:** Developer / Engineer

**Methodology step:** Phase 3 — Implementation planning

**Key actions:**
1. Read `component-spec.md` (the contract)
2. Call Component MCP: `get_patterns()`, `get_decisions()`
3. Verify no blocking ADRs remain
4. Write `impl-spec.md`: data model, code changes, edge cases, observability
5. Write `tasks.yaml`: decomposed implementation tasks
6. Optionally use Superpowers TDD + git worktrees
7. Self-check all 5 gates

**Artifacts consumed:** `component-spec.md`, Component MCP responses
**Artifacts produced:** `impl-spec.md`, `tasks.yaml`, test suite

**Related skills:** `component-spec` (prerequisite), `tdd` (implementation discipline), `verifier` (consumes impl-spec), `superpowers-bridge`

---

### `tdd`

**Description:** Enforce strict RED-GREEN-REFACTOR test-driven development discipline throughout implementation.

**When to use:**
- Implementing any feature or bug fix
- Needing a disciplined TDD cycle with explicit phase enforcement
- Coordinating multi-phase TDD workflows in the developer phase

**Role:** Developer

**Methodology step:** Phase 3 — Implementation discipline

**Key actions:**
1. Phase 1: Test spec and design
2. Phase 2: RED — write failing tests first
3. Phase 3: GREEN — write minimum code to pass
4. Phase 4: REFACTOR — improve without breaking tests
5. Phase 5: Integration tests
6. Phase 6: Continuous improvement

**Artifacts consumed:** Component specs, acceptance criteria
**Artifacts produced:** Test suite, refactored implementation code, coverage metrics

**Related skills:** `developer` (uses TDD within), `verifier` (consumes test results), `superpowers-bridge` (Superpowers TDD variant)

---

### `sdd-openspec`

**Description:** Spec-driven development that merges artifact-first design with CLI-first task execution. Manages the full lifecycle from proposal to archive.

**When to use:**
- Proposing changes starting from specs
- Planning work with proposal.md, design.md, tasks.md artifacts
- Combining artifact generation with traceable task tracking
- Component repos executing work from a platform change package

**Role:** Architect / Developer

**Methodology step:** Phase 3 — Spec-driven execution

**Key actions:**
1. Phase 1: Context sync → scaffold change → delta specs
2. Phase 2: Design → task definition
3. Phase 3: CLI ingestion (parse tasks into YAML)
4. Phase 4: Execution loop — claim → implement → validate → complete
5. Phase 5: Archive — merge delta specs into main specs

**Artifacts consumed:** Platform/component specs, context bundle
**Artifacts produced:** `proposal.md`, delta specs, `design.md`, `tasks.md`, archived spec updates

**Related skills:** `sdd-bmad`, `sdd-speckit` (complementary foundations), `unified-sdd-skill` (orchestrates OpenSpec), `developer`

---

### `sdd-bmad`

**Description:** Spec-driven software planning using the BMAD Method with progressive context building and scale-aware workflow selection.

**When to use:**
- Structured agent-friendly planning at any scale
- Guided workflows requiring clear role-specific agent personas
- Scale-aware planning (Quick Flow vs full BMAD)
- Building context progressively across discovery → planning → architecture → implementation

**Role:** PM / Architect / Developer (persona-driven)

**Methodology step:** Foundation methodology — used across all phases

**Key actions:**
1. Pick correct track: Quick Flow (small changes) vs Full BMAD (complex features)
2. Respect the phase model: discovery → planning → architecture → implementation
3. Prefer the smallest sufficient artifact
4. Treat brownfield as a first-class case
5. Map each request to the right persona (PM, Architect, Dev)
6. Ensure outputs are implementation-friendly

**Artifacts consumed:** Product briefs, requirements, existing code
**Artifacts produced:** Tech-spec (Quick), PRD, Architecture, Project-context, Stories

**Related skills:** `sdd-openspec`, `sdd-speckit` (used together), `uncle-sdd-agent`, `unified-sdd-skill` (orchestrate BMAD)

---

### `sdd-speckit`

**Description:** Turn product ideas into executable specs, plans, task lists, and workflows through an 8-phase structured approach.

**When to use:**
- Defining a feature before writing any code
- Creating specs, plans, tasks, and workflow artifacts
- Establishing project rules and engineering principles
- Both greenfield and brownfield work

**Role:** PM / Architect / Developer

**Methodology step:** Foundation methodology — specification and planning

**Key actions:**
1. Follow 8-phase workflow: constitution → specify → clarify → checklist → plan → tasks → analyze → implement
2. Define the constitution first for any meaningful work
3. Clarify all ambiguity before planning
4. Keep artifacts separated by responsibility
5. Protect the specs directory as durable source of truth

**Artifacts consumed:** Requirements, existing code, project context
**Artifacts produced:** Constitution, specs, plans, tasks, analysis, implementation guides

**Related skills:** `sdd-bmad`, `sdd-openspec` (complementary), `unified-sdd-skill`, `uncle-sdd-agent`

---

## Phase 4 — Verification

### `verifier`

**Description:** Hard stop before merge — verify that implementation matches every acceptance criterion with evidence. No unverified ACs allowed.

**When to use:**
- Implementation is complete and all tests pass
- Before marking any change ready for merge or production
- Needing documented evidence that every AC is satisfied

**Role:** QA / Verification Engineer

**Methodology step:** Phase 4 — Verification

**Key actions:**
1. Read `component-spec.md` and list all acceptance criteria
2. Read `impl-spec.md` and map code changes to each AC
3. Run tests and collect evidence per AC
4. Optionally use Superpowers verification
5. Write `verify.md` with per-AC evidence table
6. Update `spec-graph.json`: `status = Done`

**Artifacts consumed:** `component-spec.md`, `impl-spec.md`, test results
**Artifacts produced:** `verify.md` (AC-to-evidence mapping), updated `spec-graph.json`

**Related skills:** `developer` (produces impl-spec), `gate-check` (may block if ACs are untestable), `superpowers-bridge`

---

### `gate-check`

**Description:** Diagnose SDD gate failures at any phase boundary and provide specific remediation steps.

**When to use:**
- `/speckit.analyze` fails a gate
- A spec won't pass a specific gate
- Unsure what a gate requires
- Reviewing specs before fan-out to catch problems early
- Any phase handoff where quality must be confirmed

**Role:** Any role (cross-cutting)

**Methodology step:** All phase boundaries — quality gate enforcement

**Key actions:**
1. Ask 3 intake questions: which gates are failing, spec level, relevant section
2. Diagnose each of the 5 gates individually
3. Provide specific remediation per failing gate
4. Update the spec's Gates section after all pass
5. Update `spec-graph.json`: `status = Approved`

**Artifacts consumed:** Spec files, spec sections under review
**Artifacts produced:** Gate remediation guidance, updated spec `Gates` section, updated `spec-graph.json`

**Related skills:** `adr` (gate 4 checks ADR resolution), `component-spec`, `architect`, `verifier`

---

## Off-Path — Hotfix

### `hotfix`

**Description:** Guide production incidents and critical bugs through a lightweight SDD hotfix path that maintains traceability without full spec overhead.

**When to use:**
- Production incidents requiring immediate action
- Critical bugs needing urgent fixes
- Non-urgent bugs that still require structured tracking
- Any fix where you need traceability but can't wait for full SDD phases

**Role:** Developer / Team Lead

**Methodology step:** Off-path — parallel to main SDD workflow

**Key actions:**
1. Triage: does the fix change contracts or policies?
2. Create a lightweight Hotfix Spec or Bug Spec
3. Implement the fix
4. Verify with the stated Validation criterion
5. Create a required Follow-up Spec (within next sprint)

**Artifacts consumed:** Incident reports, bug descriptions
**Artifacts produced:** Hotfix Spec, Bug Spec, Follow-up Spec, updated `spec-graph.json`

**Related skills:** `verifier` (verification criterion), `gate-check` (lightweight gate on follow-up spec), `developer`

---

## Cross-cutting Support Skills

### `explain-code-skill`

**Description:** Explain how code works using visual ASCII diagrams, everyday analogies, step-by-step walkthroughs, and gotcha callouts.

**When to use:**
- Any time a team member needs to understand existing behavior
- Code review or onboarding to a module
- Any SDD phase where understanding existing code is required
- Teaching or pair-programming context

**Role:** Any role (support skill)

**Methodology step:** Cross-cutting — used throughout all phases

**Key actions:**
1. Provide an everyday analogy for the concept
2. Draw an ASCII diagram showing flow or structure
3. Give a step-by-step execution walkthrough
4. Call out one gotcha or common misconception

**Artifacts consumed:** Code, architecture documents, pull requests
**Artifacts produced:** Explanation documents with analogies, diagrams, walkthroughs

**Related skills:** `sdd-bmad`, `sdd-openspec`, `sdd-speckit` (can use explain-code at any phase)

---

### `local-doc`

**Description:** CLI tool for full-text search across spec files using SQLite FTS5. Ground all claims in spec text.

**When to use:**
- Searching specs or requirements for relevant content
- Answering questions about documented systems
- Finding the impact of proposed changes
- Discovering related specs before writing new ones

**Role:** Any role (cross-cutting)

**Methodology step:** Cross-cutting — reference tool for all phases

**Key actions:**
1. Extract search terms from the question
2. Run `local-doc search` queries
3. Read matched spec files
4. Reason over spec content
5. Ground every claim in quoted spec text

**Artifacts consumed:** Spec files (`.md`, `.mdx`, `.txt`)
**Artifacts produced:** Search results, spec references, grounded answers

**Related skills:** `component-spec`, `architect`, `developer`, `gate-check`

---

### `stakeholder-communication`

**Description:** Translate SDD workflow status into plain, business-friendly language for non-technical stakeholders.

**When to use:**
- Sending status updates during any SDD phase
- Explaining ADR blockers or gate failures to executives
- Escalating delays with clear business impact
- Producing regular progress reports

**Role:** Product Manager

**Methodology step:** Cross-cutting — parallel to all phases

**Key actions:**
1. Translate technical status to business language
2. Explain ADR blocking in accessible terms
3. Create weekly or bi-weekly progress reports
4. Identify escalation triggers and respond proactively
5. Provide plain-language explanations of gates, specs, and decisions

**Artifacts consumed:** Spec status, gate results, ADR decisions
**Artifacts produced:** Status emails, progress reports, escalation notices

**Related skills:** `adr`, `gate-check`, `workflow-router` (PM uses all three)

---

### `superpowers-bridge`

**Description:** Decision guide for choosing between CLI/SDD tools and the Superpowers plugin, and for combining both in a single workflow.

**When to use:**
- Deciding which tool to use for a given phase or task
- Combining Superpowers and CLI in a single workflow
- Setting up Superpowers alongside CLI tooling
- Unsure whether to use TDD via Superpowers or native SDD TDD skill

**Role:** Developer / Architect

**Methodology step:** Cross-cutting — tool selection throughout implementation phases

**Key actions:**
1. Follow the decision tree: requirements clear? PRD exists? tasks broken down?
2. Choose the right tool per phase
3. Map SDD skills to Superpowers equivalents
4. Use Superpowers for TDD, git worktrees, debugging, verification
5. Use CLI for task traceability, context bundles, SDD gates, multi-agent support

**Artifacts consumed:** Task requirements, implementation needs
**Artifacts produced:** Tool selection guidance, combined workflow examples

**Related skills:** `developer`, `tdd`, `verifier` (Superpowers variants available for each)

---

## Orchestrators

### `process-guide`

**Description:** Complete step-by-step guide for the full SDD v3.0 lifecycle from initiative through deployment and success measurement.

**When to use:**
- Starting a new initiative from scratch
- Onboarding a role to a specific phase
- Unsure what the next step is in an in-progress initiative
- Running a Full workflow for medium/high/critical risk

**Role:** Any role (master guide)

**Methodology step:** Master methodology guide — all phases

**Key actions:**
1. Walk through Phase 0–4 in sequence
2. Run gate checks at every phase handoff
3. Execute parallel implementation per component
4. Deploy with feature flags
5. Measure success metrics at Day 30

**Artifacts consumed:** Specs, task trackers, deployment configs
**Artifacts produced:** Phase-by-phase guidance for all 5 phases

**Related skills:** All phase-specific skills (references them all)

---

### `sdd-methodology-tools`

**Description:** Central skill router mapping 15 role-specific SDD skills to triggers and workflows. Use when you're not sure which skill applies.

**When to use:**
- Determining which SDD skill to use for your current task
- Routing based on risk level
- Selecting Quick / Standard / Full workflow agent sequence

**Role:** Any role (navigation/routing)

**Methodology step:** Routing and navigation — entry point

**Key actions:**
1. Match the trigger to the correct skill
2. Determine risk level from risk classification
3. Route to QUICK / STANDARD / FULL workflow
4. Point to the right agent sequence

**Artifacts consumed:** Initiative name, risk level
**Artifacts produced:** Routing guidance pointing to the correct skill(s)

**Related skills:** All 15 SDD skills (router for all)

---

### `uncle-sdd-agent`

**Description:** Platform-scale spec-driven development orchestrator that combines BMAD, OpenSpec, and Speckit. Uses an Assess phase to classify changes.

**When to use:**
- Unified SDD across multiple teams or components
- Platform-scale change management
- Iteration 1 / 2 planning
- Pull request reviews
- When a change needs full platform alignment before component work begins

**Role:** Any role (platform-scale orchestration)

**Methodology step:** Master orchestrator — Platform → Assess → Specify → Plan → Deliver

**Key actions:**
1. Use 5-phase workflow: Platform → Assess → Specify → Plan → Deliver
2. Roll out in 2 iterations:
   - Iteration 1: Platform + Assess + Specify
   - Iteration 2: Plan + Deliver
3. Use one change package per approved change
4. Define durable ownership artifacts
5. Map tasks to stories, keep platform alignment visible

**Artifacts consumed:** Requirements, platform context
**Artifacts produced:** Ownership artifacts, change packages, platform specs

**Related skills:** `sdd-bmad`, `sdd-openspec`, `sdd-speckit` (combined), `unified-sdd-skill` (lighter variant)

---

### `unified-sdd-skill`

**Description:** Practical platform-scale orchestrator using BMAD, OpenSpec, and Speckit. Uses a Route phase instead of Assess for simpler classification.

**When to use:**
- Unified SDD orchestration across a platform
- Platform SDD routing and change packages
- Iteration 1 / 2 planning (lighter-weight than `uncle-sdd-agent`)
- When the Assess phase feels heavy for the scope

**Role:** Any role (orchestration)

**Methodology step:** Master orchestrator — Platform → Route → Specify → Plan → Deliver

**Key actions:**
1. Use 5-phase workflow: Platform → Route → Specify → Plan → Deliver
2. Route classifies requests: component-only vs shared with platform
3. Roll out in 2 iterations
4. Use one change package per change
5. Keep platform/component alignment visible in planning and delivery

**Artifacts consumed:** Requests, platform context, component specs
**Artifacts produced:** Routed proposals, specs, tasks, delivery artifacts

**Related skills:** `sdd-bmad`, `sdd-openspec`, `sdd-speckit`, `uncle-sdd-agent` (heavier variant)

---

## Meta Skills

### `tier-enforcer`

**Description:** Audit and fix skills to comply with the 3-tier layered context model. Ensures all skills stay within line thresholds and reference resources correctly.

**When to use:**
- Creating a new skill (validate before shipping)
- Checking an existing skill for tier compliance violations
- A skill file exceeds line thresholds (Tier 2: 60–130 lines)
- Validation reports tier failures

**Role:** Skill Framework Maintainer

**Methodology step:** Meta-skill — skill governance (not part of delivery workflow)

**Key actions:**
1. Identify mode: creating / auditing / fixing
2. Create Tier 2 skills with required sections within 60–130 line budget
3. Audit skills for tier compliance
4. Fix violations by extracting content to resource files
5. Verify line counts and link resolution

**Artifacts consumed:** Existing skill files
**Artifacts produced:** Compliant skill files, extracted resource files, `packs.go` registration updates

**Related skills:** All skills (governs them all)

---

## Skill Interaction Map

```
initiative-definition ──► risk-assessment ──► workflow-router
                                                      │
                    ┌─────────────────────────────────┤
                    │              │                   │
                 QUICK          STANDARD             FULL
                    │              │                   │
                    │         architect ◄──── analyst ◄┘
                    │              │
                    └──────► component-spec ◄── adr (blocks)
                                   │
                              developer
                            ┌──────┴──────────┐
                           tdd          sdd-openspec
                                   │
                              verifier ◄── gate-check

Foundation:   sdd-bmad ◄──► sdd-openspec ◄──► sdd-speckit
                                   ▲
                          uncle-sdd-agent / unified-sdd-skill
                                   ▲
                            process-guide / sdd-methodology-tools

Cross-cutting (any phase):
  explain-code-skill · local-doc · stakeholder-communication · superpowers-bridge

Off-path: hotfix (bypasses Phase 0-3, rejoins at verifier)
Meta: tier-enforcer (governs skills, not delivery)
```

---

## Workflow Decision Summary

| Situation | Start here |
|---|---|
| Don't know which skill to use | `sdd-methodology-tools` |
| Starting a new initiative | `initiative-definition` → `risk-assessment` → `workflow-router` |
| Onboarding an existing platform | `platform-contextualizer` → `platform-constitution` |
| High/Critical risk feature | `analyst` → `architect` → `component-spec` → `developer` → `verifier` |
| Medium risk feature | `architect` → `component-spec` → `developer` → `verifier` |
| Low risk change | `developer` → `verifier` (lightweight) |
| Production incident | `hotfix` |
| Spec failing a gate | `gate-check` |
| Unresolved architecture decision | `adr` |
| Understanding existing code | `explain-code-skill` |
| Searching across all specs | `local-doc` |
| Explaining status to stakeholders | `stakeholder-communication` |
| Choosing Superpowers vs CLI | `superpowers-bridge` |
| Platform-scale orchestration | `unified-sdd-skill` or `uncle-sdd-agent` |
| Full SDD lifecycle walkthrough | `process-guide` |
| New or broken skill file | `tier-enforcer` |
