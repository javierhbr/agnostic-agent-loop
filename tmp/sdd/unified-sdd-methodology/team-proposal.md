# Unified Spec-Driven Development Methodology

## Purpose

This document proposes a simplified, phased workflow that combines BMAD,
OpenSpec, and Speckit into one operating model. The goal is to explain the
methodology from the beginning of platform interaction down to development,
deploy, and archive.

The workflow is phase-first, not artifact-first. The artifacts still matter,
but teams should understand the delivery model through a small number of
phases.

## Why simplify the workflow

The earlier draft was accurate, but it exposed too many internal steps at the
top level. That makes the method harder to teach across teams and harder to
use for phased implementation.

The simplified model keeps the strengths of each source skill:

- BMAD gives progressive context, track selection, and role-based handoffs
- OpenSpec gives the change package, proposal/spec/design/tasks flow, and archive
- Speckit gives constitution, clarify, checklist, analyze, and executable delivery discipline

## Proposed v1 model: 5 phases

The v1 methodology uses 5 phases:

1. Platform
2. Route
3. Specify
4. Plan
5. Deliver

In v1, Deliver is owned by the Team Lead. Build, verify, deploy, and archive
remain inside Deliver as delivery slices, with explicit pull request creation
and review between build and verification. If release complexity grows later,
Deliver can split into separate Build and Deploy phases.

## Adoption plan: 2 iterations

We should adopt the 5-phase model in two iterations.

Iteration 1 focuses on the front half of the workflow:

- Platform
- Route
- Specify

Iteration 2 focuses on the back half of the workflow:

- Plan
- Deliver

This rollout sequence makes adoption safer. Teams first learn how to define
shared context, route work, and produce better specifications. Then they apply
those stronger inputs to planning and delivery.

```text
TARGET MODEL

  [Platform] -> [Route] -> [Specify] -> [Plan] -> [Deliver]


ADOPTION ITERATIONS

  Iteration 1
  [Platform] -> [Route] -> [Specify]

  Iteration 2
                                   [Plan] -> [Deliver]

  Future option
                                   [Plan] -> [Build] -> [Deploy]
```

Detailed operational guide:

- Iteration 1 playbook: [iteration-1-playbook.md](iteration-1-playbook.md)
- Iteration 2 playbook: [iteration-2-playbook.md](iteration-2-playbook.md)
- Platform/component alignment guide: [canonical-platform-truth-and-component-alignment.md](canonical-platform-truth-and-component-alignment.md)
- Local platform MCP model: [local-platform-mcp-model.md](local-platform-mcp-model.md)
- Worked examples by phase, role, and skill: [example/README.md](example/README.md)
- Unified agent/skill package: [../unified-sdd-codex-skill/SKILL.md](../unified-sdd-codex-skill/SKILL.md)
- Role agents: [../unified-sdd-codex-skill/agents/README.md](../unified-sdd-codex-skill/agents/README.md)
- Existing-platform starter skill: [../platform-contextualizer-codex-skill/SKILL.md](../platform-contextualizer-codex-skill/SKILL.md)
- Local platform MCP skill: [../platform-truth-mcp-codex-skill/SKILL.md](../platform-truth-mcp-codex-skill/SKILL.md)
- Go MCP server scaffold: [../platform-truth-mcp-server/README.md](../platform-truth-mcp-server/README.md)
- Templates:
  - [templates/README.md](templates/README.md)
  - [templates/platform-template/README.md](templates/platform-template/README.md)
  - [templates/component-boilerplate/README.md](templates/component-boilerplate/README.md)
  - [templates/platform-mcp-boilerplate/README.md](templates/platform-mcp-boilerplate/README.md)
  - [templates/platform-ref.yaml](templates/platform-ref.yaml)
  - [templates/jira-traceability.yaml](templates/jira-traceability.yaml)

## End-to-end flow

```text
[PLATFORM]
  durable context, platform guardrails, shared roles, project rules
  Owner: Architect
  Support: Product, Team Lead
        |
        v
[ROUTE]
  open the change package, assess size and impact, select track
  Owner: Team Lead
  Support: Product, Architect, Engineering Manager
        |
        v
[SPECIFY]
  define the required behavior and remove ambiguity
  Owner: Product
  Support: Team Lead, Architect, Developers
        |
        v
[PLAN]
  convert the approved spec into a technical plan and delivery slices
  Owner: Architect
  Support: Team Lead, Developers, Product
        |
        v
[DELIVER]
  build, create PR, review PR, verify, deploy, archive
  Owner: Team Lead
  Support: Developers, QA / Validation, Architect, Product
```

## How agents interact with the platform and teams

The methodology uses agents as role-aligned helpers, not as isolated actors.
Agents do not replace platform governance or team ownership. They help draft,
route, refine, and verify the work inside each phase.

```text
[Platform stakeholders and team leads]
        |
        v
[Phase owner]
        |
        +--> [BMAD agent role]
        +--> [OpenSpec artifact agent]
        +--> [Speckit quality/check agent]
        |
        v
[Drafted artifacts and recommendations]
        |
        v
[Team review, correction, approval]
        |
        v
[Next phase]
```

The basic interaction pattern is the same in every phase:

- humans own intent, tradeoffs, and approval
- agents produce structured drafts, comparisons, and checks
- the phase owner decides when the work is ready to move on

Inside component repositories, the local artifact chain stays OpenSpec-only even
when the platform side used the combined methodology earlier in the flow.

## Platform spec and component spec interaction

The methodology now assumes:

- one master platform repository for shared truth
- one or more component repositories for local implementation truth
- JIRA as the delivery tracking chain
- an optional local read-only platform MCP gateway for developer access to the
  platform truth

The model is:

```text
[Platform master repo]
  platform spec
  shared contracts
  platform ADRs
  platform issue
        |
        | publishes version + refs
        v
[Component repo]
  platform-ref.yaml
  jira-traceability.yaml
  OpenSpec proposal/spec/design/tasks
  component epic
        |
        v
[Stories] -> [PRs] -> [Verification] -> [Archive]
```

Use this rule:

- platform repo = canonical shared truth
- component repo = local implementation truth
- JIRA = workflow and ownership truth
- local platform MCP = read-only developer access layer to platform truth

Only shared changes should update both platform and component truth.

```text
Local component change

  [Platform version 2026.03]
              |
              v
  [Component repo aligns to 2026.03]
              |
              v
  [Component epic] -> [Stories] -> [PRs]


Shared platform + component change

  [Platform issue]
        |
        +--> [Platform change package]
        |
        +--> [Component epic A] -> [Stories] -> [PRs]
        |
        +--> [Component epic B] -> [Stories] -> [PRs]
```

Detailed reference:

- [canonical-platform-truth-and-component-alignment.md](canonical-platform-truth-and-component-alignment.md)
- [local-platform-mcp-model.md](local-platform-mcp-model.md)
- [templates/platform-ref.yaml](templates/platform-ref.yaml)
- [templates/jira-traceability.yaml](templates/jira-traceability.yaml)
- [templates/platform-mcp-boilerplate/README.md](templates/platform-mcp-boilerplate/README.md)
- [example/08-platform-plan-to-component-openspec.md](example/08-platform-plan-to-component-openspec.md)

## Platform Plan to component OpenSpec handoff

Use this boundary rule:

- platform phases may use BMAD, OpenSpec, and Speckit together
- once the work enters a component repository, use OpenSpec only

The handoff should look like this:

```text
[Platform Plan]
  platform version
  platform refs
  shared decisions
  rollout constraints
  platform issue
        |
        v
[Component repo]
  platform-ref.yaml
  jira-traceability.yaml
  OpenSpec proposal/spec/design/tasks
        |
        v
[Component delivery]
  stories -> PRs -> verification -> archive
```

The platform handoff must give the component team:

- the platform version to pin
- the platform refs that constrain the change
- the shared design and contract decisions
- the rollout or migration expectations
- the JIRA chain that the component epic belongs to

The component team then uses OpenSpec to:

- pin alignment in `platform-ref.yaml`
- pin the issue chain in `jira-traceability.yaml`
- write the local `proposal.md`
- write the local component spec delta
- write `design.md`
- write `tasks.md`
- keep delivery traceability current through PR and archive

Do not treat the platform plan as the component implementation plan. The
platform plan sets shared boundaries. The component repo turns those
boundaries into local executable OpenSpec artifacts.

## Local platform MCP gateway

When hosted infrastructure is not available, teams may use a local read-only
platform MCP gateway.

```text
[Local platform clone]
        |
        v
[Local platform MCP]
  read-only query and validation
        |
        v
[Component repo]
  OpenSpec local artifacts
```

Use this rule:

- the MCP gateway reads platform truth and JIRA-linked metadata locally
- the MCP gateway stays read-only in v1
- component work validates against the pinned platform version by default
- component repos still keep local truth in OpenSpec artifacts

Current v1 implementation:

- Go
- stdio
- newline-delimited JSON-RPC
- read-only only
- local binary scaffold under `../platform-truth-mcp-server/`

## Phase-by-phase operating model

| Phase | Goal | Primary owner | Human interaction | Agent interaction | Skills used | Rules applied | Exit output |
| --- | --- | --- | --- | --- | --- | --- | --- |
| Platform | Define durable context for all future changes | Architect | Platform stakeholders, Product, and Team Leads agree on guardrails, quality standards, common context, versioning, and JIRA conventions | Speckit creates or refreshes the constitution, OpenSpec encodes durable project config, BMAD frames context for later routing and roles | `speckit-codex-skill`, `openspec-codex-skill`, `bmad-codex-skill` | `speckit-codex-skill/rules/constitution-rules.md`, `openspec-codex-skill/rules/project-config-template.yaml` | constitution, project config, common context, platform ref model |
| Route | Turn a request into a routed change package | Team Lead | Product, Architect, and Manager classify the request, affected scope, platform refs, and issue chain | BMAD selects the right track, OpenSpec frames the change package, agents call out size, impact, unknowns, and the next artifact | `bmad-codex-skill`, `openspec-codex-skill` | `bmad-codex-skill/rules/track-selection-rules.md`, `openspec-codex-skill/rules/artifact-rules.md` | scoped change package, selected path, initial platform/component/JIRA traceability |
| Specify | Define the required behavior before implementation | Product | Product works with impacted teams to define goals, non-goals, acceptance behavior, and whether the change is local or shared | Platform-side support may use BMAD and Speckit, but component repos use OpenSpec only for `proposal.md`, delta specs, and local traceability files | platform: `openspec-codex-skill`, `speckit-codex-skill`, `bmad-codex-skill`; component: `openspec-codex-skill` | `openspec-codex-skill/rules/artifact-rules.md`, platform-side review rules as needed | approved proposal, delta specs, resolved ambiguity, confirmed platform refs |
| Plan | Convert the approved spec into a technical execution plan | Architect | Architect, Team Lead, and Developers agree on design, dependencies, rollout approach, and story mapping | Platform planning may use BMAD and Speckit, but component repos turn the handoff into local OpenSpec `design.md` and `tasks.md` | platform: `bmad-codex-skill`, `openspec-codex-skill`, `speckit-codex-skill`; component: `openspec-codex-skill` | `openspec-codex-skill/rules/artifact-rules.md`, platform-side planning rules as needed | design, ADRs when needed, tasks, delivery slices, task-to-story traceability |
| Deliver | Execute the plan in controlled slices through review, deploy, and archive | Team Lead | Team Lead coordinates Developers, QA, Architect, and Product across execution and release | Component repos execute through OpenSpec task updates, PR traceability, verification notes, and archive-ready closure | component: `openspec-codex-skill`; platform-side review support optional outside the component repo | `openspec-codex-skill/rules/artifact-rules.md` | implemented slices, reviewed PRs, verification evidence, deployed change, archived package with traceability |

## What each phase contains

### 1. Platform

Platform is the shared starting point. It captures the rules that should not
be reinvented for every change.

Typical outputs:

- constitution or equivalent platform principles
- durable project config
- common quality and documentation rules
- role expectations for later phases

This phase is where Speckit's constitution mindset is strongest.

### 2. Route

Route is the intake gateway. It does not fully specify the work. It decides
how much workflow is needed and opens the change package.

Route applies BMAD routing rules:

- greenfield or brownfield
- small, medium, or large
- quick flow, PRD-first, or architecture-heavy path

Route also applies the size and impact split already adopted in ADR-004.

### 3. Specify

Specify defines the "what" before the "how."

This phase combines:

- OpenSpec proposal and delta specs
- platform-side clarify and scope support when needed

The result is a change package that is clear enough for planning.

Inside component repositories, use OpenSpec only for the local spec package.

### 4. Plan

Plan defines the technical execution model.

This phase combines:

- platform-side planning support where needed
- OpenSpec design and task artifacts in the component repo

The plan should end with implementation slices, not one large execution block.

### 5. Deliver

Deliver is the execution phase in v1, and it is owned by the Team Lead.

Deliver contains six internal slices:

1. Build
2. Create PR
3. Review PR
4. Verify
5. Deploy
6. Archive

This keeps release activity inside the main delivery phase while the teams are
still adopting the methodology.

Inside component repositories, Deliver is an OpenSpec-driven execution flow.

## Deliver phase: internal slices

```text
[DELIVER]
  Owner: Team Lead
        |
        +--> [Build]
        |      Developers implement the scoped tasks
        |      OpenSpec apply keeps tasks and artifacts current
        |
        +--> [Create PR]
        |      Developers open a reviewable pull request for the current slice
        |      The PR links back to the change package and completed tasks
        |
        +--> [Review PR]
        |      Team Lead assigns reviewers
        |      Review stays linked to the local OpenSpec package
        |
        +--> [Verify]
        |      QA / Validation and Developers collect evidence
        |      OpenSpec keeps verification notes explicit and traceable
        |
        +--> [Deploy]
        |      Team Lead coordinates release timing, dependencies, and rollback readiness
        |      Product and Architect support business and technical decisions
        |
        +--> [Archive]
               OpenSpec archive closes the change package and promotes the new truth
```

## Why this phased model fits the three skills

### BMAD contribution

BMAD is strongest in Route and Plan. It gives the methodology:

- progressive context
- track selection
- role-based handoffs
- planning depth that scales with complexity

### OpenSpec contribution

OpenSpec is the backbone from Route through Deliver. It gives the methodology:

- one change package as the canonical execution unit
- a stable artifact chain
- artifact revision during implementation
- archive as the end of the workflow

### Speckit contribution

Speckit strengthens Platform, Specify, and Plan. It gives the methodology:

- constitution and guardrails
- clarify before plan
- checklist and analysis before execution
- phased implementation discipline

## Initial decisions

- ADR-001 proposes federated intake with a single canonical workflow
- ADR-002 proposes the change package as the canonical execution unit
- ADR-003 proposes stage-based role accountability for each change package
- ADR-004 proposes separating change size from change impact in routing
- ADR-005 proposes the 5-phase v1 model with Team Lead owning Deliver
- ADR-006 proposes rolling out the 5-phase model in two iterations
- ADR-007 proposes making pull request creation and review explicit in Deliver
- ADR-010 proposes canonical platform truth with versioned component alignment and JIRA-linked execution
- ADR-011 proposes OpenSpec-only local workflows inside component repositories
- ADR-012 proposes a local read-only platform MCP gateway for developer access to platform truth
- ADR-013 proposes implementing the first local MCP server as a self-contained Go stdio server

## Evolution path

The next evolution is clear:

- adoption iteration 1 = Platform -> Route -> Specify
- adoption iteration 2 = Plan -> Deliver

- v1 = Platform -> Route -> Specify -> Plan -> Deliver
- v2 = Platform -> Route -> Specify -> Plan -> Build -> Deploy

We should only split Deliver when release coordination becomes complex enough
to justify a separate top-level phase.
