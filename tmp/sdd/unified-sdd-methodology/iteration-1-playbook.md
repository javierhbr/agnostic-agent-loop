# Iteration 1 Playbook

## Purpose

This playbook turns Iteration 1 of the unified SDD methodology into a simple,
practical guide for teams.

Iteration 1 covers the first three phases:

1. Platform
2. Route
3. Specify

The purpose of Iteration 1 is to improve the quality of shared context and
specification before teams change planning and delivery behavior.

## How to use this playbook

Use this document as the day-to-day guide for the first rollout of the
methodology.

- start with `Platform` once for the shared environment
- run `Route` for each new request or change
- run `Specify` before any technical planning starts

Helpful starter skill for existing platforms:

- [../platform-contextualizer-codex-skill/SKILL.md](../platform-contextualizer-codex-skill/SKILL.md)
- Detailed alignment guide: [canonical-platform-truth-and-component-alignment.md](canonical-platform-truth-and-component-alignment.md)
- Local MCP guide: [local-platform-mcp-model.md](local-platform-mcp-model.md)
- Worked examples: [example/README.md](example/README.md)
- Templates:
  - [templates/platform-ref.yaml](templates/platform-ref.yaml)
  - [templates/jira-traceability.yaml](templates/jira-traceability.yaml)
  - [templates/platform-mcp-boilerplate/README.md](templates/platform-mcp-boilerplate/README.md)

The rule is simple:

- do not plan before `Specify` is ready
- do not implement from a weak or ambiguous spec
- treat the platform repo as the canonical shared truth and component repos as version-aligned local truth

## Iteration 1 at a glance

```text
[PLATFORM]
  create shared context and durable rules
        |
        v
[ROUTE]
  turn a request into a scoped change package
        |
        v
[SPECIFY]
  define the required behavior and remove ambiguity
        |
        v
Ready for Iteration 2
        |
        v
[PLAN] -> [DELIVER]
```

## Phase 1: Platform

If the platform already exists and teams are already working, start by using
the `Platform Contextualizer` skill to review the current state before trying
to define the durable baseline.

### Phase flow

```text
[Platform stakeholders]
  Architect + Product + Team Lead
            |
            v
[Define durable rules and shared context]
            |
            +--> Speckit: draft constitution and quality bar
            +--> OpenSpec: encode reusable config and rules
            +--> BMAD: frame roles, context, and later routing needs
            |
            v
[Platform baseline]
  constitution + config + role map + quality language
            |
            v
Ready for Route
```

### 1. Main objectives and outcomes

Objectives:

- define the shared rules that all later changes must follow
- create one common baseline for quality, testing, and documentation
- align the platform and the teams on how agents should support the workflow

Outcomes:

- a durable constitution or principles document
- an OpenSpec project config with stable context and reusable rules
- a clear role model for later phases

### 2. Concepts and activities to learn and apply

Core concepts:

- durable context vs change-specific detail
- explicit principles vs vague values
- progressive context before execution
- platform guardrails before team-level change work

Main activities:

- identify the rules that should apply to every future change
- separate permanent standards from temporary project details
- define the minimum quality bar for specs and delivery
- agree on the phase owners and support roles
- define the platform versioning and ref model for component repositories
- define the JIRA hierarchy or issue-link conventions for platform and component work
- define whether teams will use a local read-only platform MCP gateway

### 3. Agent roles and responsibilities

Human roles:

- Architect owns the phase and decides what becomes durable platform context
- Product contributes business constraints, UX expectations, and quality concerns
- Team Lead contributes team conventions, delivery realities, and adoption needs

Agent roles:

- BMAD Architect agent frames shared context, role boundaries, and future routing needs
- OpenSpec config agent turns durable context into reusable project configuration
- Speckit constitution agent drafts explicit, testable platform principles

### 4. Skills used and how they are applied

- `speckit-codex-skill`
  - use first to draft the constitution and quality bar
  - use it to make principles explicit, testable, and durable
- `openspec-codex-skill`
  - use second to encode stable project context and reusable rules
  - use it to create or refine `openspec/config.yaml`
- `bmad-codex-skill`
  - use third to align role boundaries, planning expectations, and progressive context
  - use it to prepare later routing and architecture work

### 5. Rules that govern interactions and outputs

Apply these rules:

- `speckit-codex-skill/rules/constitution-rules.md`
  - principles must be explicit and testable
  - quality, testing, security, and UX expectations must be visible
- `openspec-codex-skill/rules/project-config-template.yaml`
  - keep stable context in config
  - keep artifact rules concise and reusable
- `bmad-codex-skill/rules/artifact-rules.md`
  - every durable artifact should state purpose, constraints, risks, and success criteria

### 6. Expected artifacts and deliverables

Expected outputs:

- constitution or principles document
- `openspec/config.yaml` or equivalent reusable config
- shared role map for the workflow
- common language for quality and artifact expectations
- platform versioning and ref model
- JIRA hierarchy conventions for platform issue, component epic, and stories
- local platform MCP usage model when teams need local query and validation
- adoption-ready templates for `platform-ref.yaml` and `jira-traceability.yaml`

### 7. Criteria for moving to the next phase

Move to `Route` when:

- principles are explicit enough for teams to follow
- durable context is separated from temporary change detail
- role ownership is clear for Platform, Route, and Specify
- there is a reusable config or rules baseline for future changes

### 8. Potential challenges and mitigation strategies

- Challenge: principles are too vague
  - Mitigation: rewrite them as short, testable statements
- Challenge: teams try to put everything into platform rules
  - Mitigation: keep only durable standards here and move change detail to later phases
- Challenge: teams disagree on standards
  - Mitigation: start with a minimum viable baseline and refine after real usage

### 9. Feedback and iteration process

- review the first few routed changes against the constitution
- track which rules teams actually use and which ones they ignore
- refine the constitution and config after each pilot or retro
- keep the platform baseline small, practical, and easy to maintain

## Phase 2: Route

### Phase flow

```text
[Incoming request]
  initiative / requirement / team proposal
            |
            v
[Team Lead owns routing]
  Support: Product + Architect + Engineering Manager
            |
            +--> BMAD: classify greenfield/brownfield, size, and track
            +--> OpenSpec: open or frame the change package
            +--> Speckit: clarify only if the request is too vague
            |
            v
[Routed change package]
  scope + size + impact + path + next artifact + known unknowns
            |
            v
Ready for Specify
```

### 1. Main objectives and outcomes

Objectives:

- turn a request into a scoped change package
- decide how much workflow depth is required
- identify the next artifact and the right owner

Outcomes:

- one routed change package
- a size and impact classification
- a selected track for the next step
- a short list of known unknowns

### 2. Concepts and activities to learn and apply

Core concepts:

- intake normalization
- change package as the canonical execution unit
- size vs impact
- quick path vs deeper path
- greenfield vs brownfield
- platform-only vs component-only vs shared change
- platform issue -> component epic -> story traceability

Main activities:

- review the incoming request and its entry point
- identify affected scope, teams, and dependencies
- identify affected platform refs and affected component repositories
- classify size and impact separately
- classify the change as local-only, shared, or platform rule adoption
- choose the BMAD path depth
- open the change package and define the next artifact
- create the initial JIRA issue chain and alignment metadata

### 3. Agent roles and responsibilities

Human roles:

- Team Lead owns the phase and makes the routing decision
- Product clarifies business priority, intent, and non-goals
- Architect identifies architecture, integration, and cross-team impact
- Engineering Manager supports prioritization, staffing, and escalation

Agent roles:

- BMAD routing agent classifies project type, size, and track
- OpenSpec change agent frames the change package and the next artifact
- Speckit clarify agent is optional here when the request is too vague to route safely

### 4. Skills used and how they are applied

- `bmad-codex-skill`
  - use first to classify greenfield or brownfield, size, and path depth
  - use it to decide quick flow, PRD-first, or architecture-heavy work
- `openspec-codex-skill`
  - use second to frame the change package and prepare the next artifact
  - use `/opsx:explore` when the request is still fuzzy
  - use `/opsx:propose` when the team is ready to open the change
- `speckit-codex-skill`
  - use only when ambiguity blocks safe routing
  - keep it lightweight at this phase

### 5. Rules that govern interactions and outputs

Apply these rules:

- `bmad-codex-skill/rules/track-selection-rules.md`
  - classify project type
  - classify size
  - pick the correct planning track
  - escalate to architecture when required
- `openspec-codex-skill/rules/artifact-rules.md`
  - keep the change small, clear, and reviewable
  - call out unknowns instead of hiding them
- `bmad-codex-skill/rules/artifact-rules.md`
  - the routed output must include scope, assumptions, risks, and success criteria

### 6. Expected artifacts and deliverables

Expected outputs:

- routed change package
- intake summary
- size and impact classification
- selected path and next artifact
- known unknowns and open questions
- initial `platform-ref.yaml`
- initial `jira-traceability.yaml`

### 7. Criteria for moving to the next phase

Move to `Specify` when:

- the change package exists and has one clear owner
- size and impact are explicit
- the affected scope is bounded enough to specify
- the next artifact is clear
- major unknowns are visible, even if not all are resolved yet
- the initial platform refs and issue chain are visible

### 8. Potential challenges and mitigation strategies

- Challenge: the team jumps straight to implementation
  - Mitigation: require a routed change package before any planning starts
- Challenge: size and impact are mixed together
  - Mitigation: classify them separately every time
- Challenge: brownfield constraints are missed
  - Mitigation: identify existing conventions and integration points early
- Challenge: everything is treated as high priority
  - Mitigation: force explicit scope, affected teams, and path choice

### 9. Feedback and iteration process

- compare route decisions with actual downstream effort
- track where teams under-classified or over-classified changes
- refine size and impact examples after each pilot
- update routing guidance when repeated mistakes appear

## Phase 3: Specify

### Phase flow

```text
[Routed change package]
            |
            v
[Product owns behavior definition]
  Support: Team Lead + Architect + Developers
            |
            +--> OpenSpec: draft proposal.md, delta specs, and local traceability
            +--> platform-side review: clarify and scope support only when needed
            |
            v
[Spec package]
  proposal + delta specs + clarification notes + checklist result
            |
            v
Ready for Plan
```

### 1. Main objectives and outcomes

Objectives:

- define the required behavior before planning
- remove ambiguity that would create rework later
- produce a spec package another team or agent can use safely

Outcomes:

- approved proposal
- delta specs with concrete requirements
- clarified assumptions and edge cases
- a readiness signal for planning

### 2. Concepts and activities to learn and apply

Core concepts:

- define the "what" before the "how"
- behavior over implementation
- goals and non-goals
- testable requirements and scenarios
- clarify before plan
- component-local behavior vs shared platform truth
- platform version and platform refs as part of spec context

Main activities:

- write the problem statement in plain language
- define goals, non-goals, and affected behavior
- confirm the platform version and platform refs that constrain the change
- draft delta specs using explicit `ADDED`, `MODIFIED`, and `REMOVED` sections
- decide whether the change needs a linked platform delta or only a component delta
- run clarify to expose hidden assumptions
- run a checklist pass before planning

### 3. Agent roles and responsibilities

Human roles:

- Product owns the phase and approves the behavior definition
- Team Lead protects scope and confirms readiness for planning
- Architect identifies hard constraints without turning the spec into a design doc
- Developers raise edge cases, failure behavior, and hidden dependencies

Agent roles:

- OpenSpec artifact agent drafts `proposal.md` and delta specs
- platform-side review agents may expose ambiguity and readiness gaps before the
  component package is approved

### 4. Skills used and how they are applied

- `openspec-codex-skill`
  - use first and use it as the only component-repo skill
  - use it to create `proposal.md`, delta specs, `platform-ref.yaml`, and
    `jira-traceability.yaml`
  - keep behavior changes explicit and comparable to current reality
- platform-side support
  - BMAD or Speckit may still help upstream when the shared change needs scope
    control or clarification before the component package is approved
  - do not mix them into the local component change package

### 5. Rules that govern interactions and outputs

Apply these rules:

- `openspec-codex-skill/rules/artifact-rules.md`
  - proposal must separate goals from non-goals
  - specs must use concrete behavior and scenarios
- component rule
  - once the work is inside the component repo, use OpenSpec only
  - do not mix platform-side methodology support into the local artifact chain

### 6. Expected artifacts and deliverables

Expected outputs:

- `proposal.md`
- delta specs with `ADDED`, `MODIFIED`, and `REMOVED` sections
- clarification notes or resolved ambiguity log
- checklist results
- ready-for-plan decision
- confirmed `platform-ref.yaml`
- updated `jira-traceability.yaml`
- linked platform delta when shared truth changes

### 7. Criteria for moving to the next phase

Move to `Plan` when:

- goals and non-goals are explicit
- major behavior changes are testable
- the spec is understandable by Product and engineering
- important ambiguity is resolved or explicitly tracked
- the team agrees that planning can start without guessing
- the component spec is aligned to explicit platform refs

### 8. Potential challenges and mitigation strategies

- Challenge: the spec becomes a design doc too early
  - Mitigation: move technical detail to Plan unless it is a hard requirement
- Challenge: requirements are vague or promotional
  - Mitigation: rewrite them as specific behaviors and scenarios
- Challenge: scope keeps expanding during specification
  - Mitigation: keep non-goals visible and route out-of-scope work into a separate change
- Challenge: edge cases are missed
  - Mitigation: include Developers and Architect in clarify review

### 9. Feedback and iteration process

- review implementation rework caused by missing or weak requirements
- collect examples of good and bad specs from the first pilots
- tighten the checklist when recurring gaps appear
- refine proposal and spec templates based on real change packages

## What success looks like after Iteration 1

Iteration 1 is successful when teams can do the following consistently:

- define shared context once instead of repeating it in every change
- route work into one clear change package
- produce specs that are clear enough to plan without guesswork

At that point, the organization is ready to move into Iteration 2:

- Plan
- Deliver
