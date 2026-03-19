# OPERATING MODEL
## Unified Spec-Driven Development

> **Version:** Unified SDD -- March 2026
>
> **Core Principle:** *"Nothing is implemented without a validated spec backed by governed context."*

---

# 1) FORMAL OPERATING MODEL

## 1.0 Purpose

Define a repeatable, auditable, and scalable system to build software where:

- Knowledge is explicit (specs, ADRs, contracts, ownership artifacts)
- Context is governed (platform constitution, component knowledge)
- Execution is traceable (change packages, spec graph, jira-traceability)
- Humans and AI agents collaborate with clear role boundaries
- All process lives in Markdown and YAML -- versionable, editable, zero overhead infrastructure

---

## 1.1 Core Principle

> **"Nothing is implemented without a validated spec backed by governed context."**

Agents -- human or AI -- do not invent context. They retrieve governed knowledge, produce traceable artifacts, and hand off through explicit gates. If a methodology component does not change what the agent writes, it is eliminated.

---

## 1.2 Methodology Origins -- What We Took and Why

Unified SDD is a deliberate composition of proven frameworks, not a new methodology invented from scratch. Every practice was adopted against a single filter: *"Does this change what the agent or engineer writes? If not, it is eliminated."*

Three frameworks contributed the core structural ideas. Their specific contributions -- and the rationale behind each adoption -- are documented here for any team member who needs to understand, challenge, or evolve the model.

---

### BMAD -- Breakthrough Method for Agile AI Development

BMAD defines the AI-native multi-agent workflow architecture. Its strongest contributions fall in the **Assess** and **Plan** phases of the 5-phase model.

**1. Progressive context and role handoffs**
Each phase has a designated owner role. Context flows forward through explicit artifacts -- never through shared memory or implicit state. The output of one phase is the input contract for the next.

*Rationale:* Progressive context prevents context overreach. When a role writes outside its bounded responsibility, it invents context it does not have the right sources for. Strict role separation means each role calls only the tools it needs and produces outputs that are predictable, auditable, and attributable.

**2. Sequential handoff with hard stops**
Each role completes its exit gate checklist before the next role starts. The Verify step in Deliver is a hard stop -- no merge is possible until every AC has observable evidence.

*Rationale:* Hard stops at each handoff make quality gates enforceable without an external validator service. The role itself is the gate. Removing the hard stop turns the gate into a suggestion.

**3. Track selection and parallel fan-out**
Once design is approved, one developer instance runs per affected component simultaneously. Verification waits for all to complete before producing evidence.

*Rationale:* Sequential execution is the primary bottleneck in cross-domain features. Components that do not depend on each other can be built concurrently. Fan-out maintains full traceability -- every instance writes against the same approved design.

**4. Role definitions as Markdown files**
Each role is a `.md` file loaded on demand. Behaviour is expressed in text, not code.

*Rationale:* The team can evolve role behaviour via PR without a deployment. The update path for roles is identical to the update path for specs and templates.

*What we dropped from BMAD:*
- Orchestrator agent (meta-agent routing to specialists) -- replaced by phase ownership rules. The phase determines which role is active. A separate routing agent is unnecessary.
- Persona-based naming -- role names are sufficient. Persona framing does not change what gets written.

---

### OpenSpec -- Open Specification Format

OpenSpec defines the structured spec format for distributed systems. It provides the backbone from **Assess** through **Deliver** -- the change package, artifact chain, and archive pattern.

**1. Change package as canonical execution unit**
Every unit of work is a change package: a versioned bundle containing proposal, design, tasks, and evidence. The change package is opened in Assess and archived in Deliver.

*Rationale:* Without a bounded execution unit, work is scattered across tickets, PRs, and conversations. The change package makes every change self-contained, auditable, and recoverable.

**2. Mandatory traceability metadata on every spec**
Every spec must declare: `implements` (parent spec ID + version), `context_pack` (the versioned snapshot used when the spec was authored), `blocked_by` (open ADR IDs preventing implementation), and `status` (Draft / Approved / Implementing / Done / Paused / Blocked).

*Rationale:* Without these fields, a spec is an island -- you cannot tell what it is implementing, what context it was written against, or whether it is safe to act on.

**3. Spec Graph (graph/index.yaml)**
A machine-readable index linking every artifact in the chain: Initiative to proposal to design to tasks to evidence to ADRs. Updated after every merge.

*Rationale:* The Spec Graph is the audit trail. It answers "what was built, why, against what context, and was it verified?" -- questions critical in post-incident analysis, compliance reviews, and onboarding.

**4. Specs are never deleted -- only versioned**
Status transitions: Draft to Approved to Implementing to Done. Or to Paused (waiting on business priority), or to Blocked (waiting on ADR resolution). Paused specs preserve their context pack version so they can be safely rebased on resumption.

*Rationale:* Deleted specs create invisible debt. A Paused spec with its context pack pinned can be rebased and resumed with full awareness of what has changed.

*What we dropped from OpenSpec:*
- Spec versioning server -- git handles versioning, the Spec Graph handles links. No separate server.

---

### Speckit -- Template-Driven Spec Assembly

Speckit defines the template-driven pattern for AI-assisted spec authoring. Its strongest contributions fall in the **Platform**, **Specify**, and **Plan** phases -- constitution, clarify, and checklist patterns.

**1. Constitution as platform foundation**
The constitution document encodes all non-negotiable platform rules, NFR baselines, and governance policies. Every spec is written against the constitution as ground truth.

*Rationale:* Without a constitution, platform rules are tribal knowledge. Two teams writing specs two weeks apart against different assumptions produce inconsistent specs. The constitution makes constraints explicit and versionable.

**2. Clarify pattern with embedded exit gate checklists**
Spec templates are not blank forms. Each section includes the self-verification checklist the role must complete before handing off. The template and the gate are the same artifact.

*Rationale:* Separating "how to write the spec" from "how to verify the spec" creates a gap where quality degrades. Embedding the gate in the template closes that gap.

**3. Context Pack -- versioned snapshot assembled before spec writing**
Before any spec is written, the role assembles a pinned snapshot of all applicable policies, NFR baselines, and workflow config. This snapshot is recorded as `context_pack` in the spec metadata.

*Rationale:* Writing a spec against a moving context is the root cause of spec drift. Pinning the context pack version makes the spec reproducible.

**4. Templates served by platform knowledge**
Templates are not stored locally with each role. The role fetches the current template version from governed platform knowledge at runtime.

*Rationale:* Template evolution is zero-friction. Update the template and every role immediately gets the new version. No role file changes, no deployment, no coordination required.

*What we dropped from Speckit:*
- Spec scaffolding CLI -- platform knowledge serves templates inline. The role fetches and fills the template in the same context window.

---

# 2) THE 5-PHASE MODEL

The unified SDD model organizes all work into five sequential phases. Each phase has a designated owner, specific inputs and outputs, and an exit gate that must pass before the next phase begins.

```
  Phase 1        Phase 2        Phase 3        Phase 4        Phase 5
  PLATFORM  -->  ASSESS    -->  SPECIFY   -->  PLAN      -->  DELIVER
  Architect      Team Lead      Product        Architect      Team Lead
  (durable)      (per change)   (per change)   (per change)   (per change)
```

---

## 2.1 Phase 1: Platform (Owner: Architect)

**Goal:** Establish durable context that governs all downstream work.

Platform is not executed per change. It is established once and maintained continuously. Every other phase reads from Platform artifacts as ground truth.

**Artifacts produced:**
- `constitution/policies.md` -- non-negotiable platform rules, NFR baselines, governance
- `component-ownership.yaml` -- which team owns which bounded context
- `dependency-map.yaml` -- three-tier dependency classification (internal, partner, external)
- `shared-glossary.yaml` -- canonical term definitions across all components
- Role definitions (Markdown files loaded on demand)
- Spec templates with embedded exit gates

**Exit gate:** Constitution exists, ownership is declared for every component, dependency map covers all integration points, glossary has no ambiguous terms.

---

## 2.2 Phase 2: Assess (Owner: Team Lead)

**Goal:** Classify the change, open a change package, and establish traceability.

**Inputs:** Intent (feature request, bug report, incident, or initiative).

**Activities:**
1. Classify change type: Feature / Change / Bug / Hotfix
2. Determine risk level: low / medium / high / critical
3. Open change package (versioned directory with standard structure)
4. Create `platform-ref.yaml` -- pins the constitution version and applicable policies
5. Create `jira-traceability.yaml` -- links to external tracking (Jira, Linear, etc.)
6. Identify affected components from `component-ownership.yaml`
7. Check `shared-glossary.yaml` for any term ambiguity

**Artifacts produced:**
- Change package directory (opened)
- `platform-ref.yaml`
- `jira-traceability.yaml`
- Risk classification record

**Exit gate:** Change package opened. Risk classified. Affected components identified. Platform reference pinned. Traceability link established.

---

## 2.3 Phase 3: Specify (Owner: Product)

**Goal:** Define what must change in business terms, with traceable acceptance criteria.

**Inputs:** Change package from Assess, platform-ref.yaml, component ownership.

**Activities:**
1. Write `proposal.md` -- problem statement, goals, non-goals, user experience
2. Define acceptance criteria (Given/When/Then format, minimum 3)
3. Produce delta specs for each affected component -- what changes relative to current state
4. Run glossary check -- confirm all terms used in the proposal exist in `shared-glossary.yaml`
5. Declare contract changes (new/modified APIs, events, schemas)

**Artifacts produced:**
- `proposal.md`
- Delta specs per affected component
- Updated glossary entries (if new terms introduced)

**Exit gate:** Proposal has concrete acceptance criteria. Delta specs exist for every affected component. Glossary check passes. Contract changes declared or confirmed as "none".

---

## 2.4 Phase 4: Plan (Owner: Architect)

**Goal:** Translate the proposal into an implementable design with bounded tasks.

**Inputs:** Proposal.md, delta specs, platform-ref.yaml, component knowledge.

**Activities:**
1. Write `design.md` -- data model, code changes, integration points, edge cases, rollback strategy
2. Validate against component invariants -- no invariant may be violated
3. Validate against dependency map -- integration safety confirmed
4. Produce `tasks.md` -- ordered task list with owner, dependencies, and estimates
5. Perform story mapping -- confirm tasks cover all acceptance criteria
6. Run gate check (5 gates: context completeness, domain validity, integration safety, NFR compliance, ready to implement)

**Artifacts produced:**
- `design.md`
- `tasks.md`
- Gate check results (PASS/FAIL per gate)

**Exit gate:** All 5 gates pass. Design covers every acceptance criterion. Tasks are bounded and ordered. No open ADR blocking implementation. Rollback strategy defined.

---

## 2.5 Phase 5: Deliver (Owner: Team Lead)

**Goal:** Build, verify, and archive the change.

Deliver is a compound phase with six internal slices (see Section 6). The Team Lead coordinates execution across all slices.

**Inputs:** Approved design.md, tasks.md, platform-ref.yaml.

**Activities:** Build, Create PR, Review PR, Verify, Deploy, Archive.

**Artifacts produced:**
- Code changes (per task)
- Pull request(s)
- Review evidence
- `verify.md` -- AC verification with observable evidence
- Updated Spec Graph (`graph/index.yaml`)
- Archived change package

**Exit gate:** Every AC verified with evidence. All tests pass. Spec Graph updated. Change package archived with status Done.

---

# 3) CORE BOUNDARY RULE

The 5-phase model applies differently depending on repository scope.

| Scope | Phases Available | Methodology Sources |
|---|---|---|
| **Platform-side** (specs repo, constitution, governance) | All 5 phases | BMAD + OpenSpec + Speckit |
| **Component repos** (system code, local specs) | Assess, Specify (delta only), Plan, Deliver | OpenSpec only |

**Rule:** Component teams consume platform artifacts (constitution, templates, ownership, glossary) as read-only input. They contribute to them via PR to the platform repo, never by direct edit.

This boundary ensures that platform governance is centralized and versionable, while component teams retain autonomy over their implementation decisions.

---

# 4) KEY CONCEPTS

## 4.1 Change Package as Canonical Execution Unit

Every unit of work -- feature, bug, hotfix, or initiative -- is a change package. A change package is a versioned directory containing:

```
changes/CP-{ID}/
  platform-ref.yaml        # pinned constitution version
  jira-traceability.yaml   # external tracking links
  proposal.md              # what and why (Specify output)
  design.md                # how (Plan output)
  tasks.md                 # execution plan (Plan output)
  verify.md                # evidence (Deliver output)
```

The change package is opened in Assess and archived in Deliver. It is the single source of truth for "what happened, why, and was it verified."

## 4.2 Three DDD Ownership Artifacts

Three artifacts establish domain-driven ownership across the platform:

**1. Component Ownership (`component-ownership.yaml`)**
Maps every bounded context to its owning team. No component may be unowned. Ownership determines who reviews specs, who approves contract changes, and who is accountable for invariants.

**2. Dependency Map (`dependency-map.yaml`)**
Classifies every integration point into three tiers:

| Tier | Description | Change Protocol |
|---|---|---|
| Internal | Same team, same repo | Standard PR review |
| Partner | Different team, shared contract | Contract change declaration + consumer notification |
| External | Third-party API or service | Compatibility plan + dual-publish if breaking |

**3. Shared Glossary (`shared-glossary.yaml`)**
Canonical definitions for all domain terms used across components. Every proposal and design must use glossary terms. New terms must be added to the glossary before they appear in specs.

## 4.3 Traceability Artifacts

**platform-ref.yaml:** Pins the constitution version and applicable policies at the time the change package was opened. Ensures every downstream artifact is written against a known, reproducible policy baseline.

**jira-traceability.yaml:** Links the change package to external issue trackers (Jira, Linear, GitHub Issues). Enables bidirectional navigation between methodology artifacts and project management tools.

## 4.4 Two-Iteration Adoption Plan

Teams adopt the unified model in two iterations:

**Iteration 1 -- Foundation:**
- Establish Platform phase artifacts (constitution, ownership, dependency map, glossary)
- Run Assess and Specify for one pilot change
- Execute Plan and Deliver with gate checks
- Validate end-to-end traceability

**Iteration 2 -- Full Model:**
- All changes flow through the 5-phase model
- Deliver slices are fully operational (Build through Archive)
- Spec Graph updated after every merge
- CI enforces verify.md sign-off before merge

---

# 5) PHASE-BY-PHASE OPERATING TABLE

| Phase | Goal | Owner | Tools Used | Skills Invoked | Exit Output |
|---|---|---|---|---|---|
| **Platform** | Durable context, constitution, ownership | Architect | Platform knowledge, templates | platform-spec, platform-constitution | constitution, ownership.yaml, dependency-map.yaml, glossary.yaml |
| **Assess** | Classify, open change package, pin context | Team Lead | Component ownership, glossary | openspec (init), agentic-helper | platform-ref.yaml, jira-traceability.yaml, change package (opened) |
| **Specify** | Define what changes in business terms | Product | Proposal templates, glossary check | product-wizard, openspec (specify) | proposal.md, delta specs, glossary updates |
| **Plan** | Implementable design with bounded tasks | Architect | Component knowledge, invariants, patterns | dev-plans, gate-check, adr | design.md, tasks.md, gate results (5 gates PASS) |
| **Deliver** | Build, verify, archive | Team Lead | CI/CD, test runners, spec graph | tdd, verifier, hotfix | code, verify.md, archived change package, updated spec graph |

---

# 6) DELIVER INTERNAL SLICES

Deliver is a compound phase. Work progresses through six ordered slices:

```
  Build --> Create PR --> Review PR --> Verify --> Deploy --> Archive
```

### Slice 1: Build

Developer implements code changes against approved tasks.md. Each task is bounded, typed, and traceable to a specific acceptance criterion.

- Read design.md and tasks.md
- Implement per task, one at a time
- Write tests per task (TDD when applicable)
- Commit with task ID in message

### Slice 2: Create PR

Developer opens a pull request linking to the change package.

- PR title references change package ID
- PR body includes summary of changes and link to design.md
- All commits listed and mapped to tasks

### Slice 3: Review PR

Reviewer verifies code against design.md and component invariants.

- Check: code matches design intent
- Check: no invariant violations
- Check: contract changes match declaration in proposal
- Check: test coverage covers acceptance criteria

### Slice 4: Verify

Verifier produces `verify.md` with evidence for every acceptance criterion.

- Every AC must have observable evidence (test output, metric, screenshot)
- No unresolved "REQUIRES HUMAN APPROVAL" items
- All tests and linting pass
- Gate check: all 5 gates PASS

### Slice 5: Deploy

Team Lead coordinates deployment per the rollout plan in design.md.

- Feature flag activation (if applicable)
- Rollback plan confirmed and ready
- Observability validation (metrics, alerts, logs confirmed)

### Slice 6: Archive

Team Lead closes the change package.

- Update Spec Graph (`graph/index.yaml`) with final links
- Set change package status to Done
- Capture learnings (if any decisions were made under pressure, create follow-up ADR)
- Link archived package to jira-traceability.yaml

---

# 7) AGENT INTERACTION MODEL

## 7.1 Agents Are Role-Aligned Helpers

AI agents in the unified model are assistants bound to specific phase roles. They are not autonomous decision-makers.

| Principle | Description |
|---|---|
| **Agents produce drafts** | Every agent output is a draft until a human approves it |
| **Agents retrieve, never invent** | Agents call governed knowledge sources for context. They do not fabricate context |
| **Agents check gates** | Agents run exit gate checklists and report PASS/FAIL. They do not override FAIL |
| **Agents compare options** | When ambiguity exists, agents present alternatives with tradeoffs. They do not choose |

## 7.2 Human Responsibilities

| Responsibility | Description |
|---|---|
| **Intent** | Humans define what should be built and why |
| **Tradeoffs** | Humans make tradeoff decisions when agents present alternatives |
| **Approval** | Humans approve spec transitions (Draft to Approved) and merge gates |
| **Escalation** | Humans resolve blocked states (ADR decisions, priority conflicts, ambiguity) |

## 7.3 Interaction Pattern

```
  Human defines intent
       |
       v
  Agent retrieves context (platform knowledge, component knowledge)
       |
       v
  Agent produces draft artifact (proposal, design, tasks, verify)
       |
       v
  Human reviews, approves, or requests revision
       |
       v
  Agent updates artifact and re-checks gates
       |
       v
  Human approves gate passage --> next phase begins
```

Agents never skip a phase. Agents never approve their own output. Agents never merge without human sign-off on verify.md.

---

# 8) EVOLUTION PATH

## 8.1 Version Roadmap

**v1 (current): 5 phases with Deliver as one compound phase**
- Platform, Assess, Specify, Plan, Deliver
- Deliver contains 6 internal slices (Build through Archive)
- Suitable for teams up to 10 people with moderate deployment complexity

**v2 (future): Split Deliver into Build + Deploy**
- Platform, Assess, Specify, Plan, Build, Deploy
- Build phase: code, PR, review, verify
- Deploy phase: rollout, observability, archive
- Suitable for teams with dedicated platform engineering or SRE functions
- Trigger: when deploy coordination becomes a bottleneck separate from code review

## 8.2 Adoption Path

**Iteration 1: Foundation (weeks 1-4)**
- Establish Platform artifacts (constitution, ownership, dependency map, glossary)
- Select one low-risk pilot feature
- Run the full 5-phase cycle end-to-end
- Validate traceability: can you trace from code back to intent through every artifact?
- Retrospective: what was useful, what was overhead, what was missing?

**Iteration 2: Full Model (weeks 5-8)**
- All new features flow through the 5-phase model
- Bugs use abbreviated Assess-Plan-Deliver (Specify is optional for single-component bugs)
- Hotfixes use minimal Assess-Deliver (design captured retroactively)
- Spec Graph updated after every merge
- CI enforces verify.md sign-off
- Retrospective: measure rework rate, time-to-first-impl, contract break incidents

## 8.3 Scaling Stages

| Stage | Team Size | Model Adaptation |
|---|---|---|
| **Early** | 1-2 people | Single repo. Platform artifacts in local directory. Deliver slices informal. |
| **Growing** | 3-6 people | Split platform repo from component repos. Formal change packages. Standard gate checks. |
| **Full** | 6+ people | Full 2-repo model. CI enforcement. Human approval gates. Parallel fan-out per component. |

The 5-phase model and role definitions never change across stages. Only the physical organization of content and the degree of enforcement evolve.

---

# 9) GATES

Five gates are checked at the Plan exit and again at the Verify step in Deliver:

| Gate | Name | Checks |
|---|---|---|
| 1 | Context Completeness | Platform-ref pinned. Constitution referenced. Glossary check passed. |
| 2 | Domain Validity | No invariant violations. Ownership respected. No cross-boundary direct access. |
| 3 | Integration Safety | Consumers identified. Compatibility plan present. Dual-publish if breaking. |
| 4 | NFR Compliance | Logging, metrics, tracing declared. PII handling addressed. Performance targets set. |
| 5 | Ready to Implement | No open BlockedBy ADRs. Spec unambiguous. ACs testable in Given/When/Then format. |

> Progress MUST be blocked if any gate is not met. Gates are embedded as checklists in spec templates. No external gate validator services.

---

# 10) SUCCESS METRICS

| Metric | What It Measures |
|---|---|
| % of changes with a change package before Build starts | Spec discipline adoption |
| % of proposals with platform-ref.yaml pinned | Platform governance coverage |
| % of impactful changes with an ADR | Decision traceability |
| Contract break incidents in production | Integration Safety gate effectiveness |
| Hotfix frequency | Preventive gate quality |
| Rework rate | Spec clarity and gate completeness |
| Time from proposal.md to first task completion | End-to-end pipeline efficiency |
| Glossary coverage (terms used vs. terms defined) | Domain alignment |

---

*"Software is no longer just built -- it is specified, validated, and executed as a system of knowledge."*

---

*Unified SDD -- March 2026*
