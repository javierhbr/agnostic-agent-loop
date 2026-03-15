---
name: unified-sdd-orchestrator
summary: Use for platform-scale spec-driven development that combines BMAD's progressive planning and roles, OpenSpec's change package and artifact flow, and Speckit's executable specification discipline.
triggers:
  - unified sdd
  - platform sdd
  - change package
  - platform route
  - platform specify
  - platform plan
  - platform deliver
  - pull request review
  - iteration 1
  - iteration 2
---

# Unified SDD Orchestrator

Use this skill when the team wants one practical workflow across multiple
teams, artifacts, and delivery phases without losing the strengths of BMAD,
OpenSpec, and Speckit. Use `../explain-code-codex-skill/SKILL.md` as a support
skill when roles need to explain existing code, planned behavior, or PR impact.

This skill is phase-first, not artifact-first.

It uses a 3-layer context model:

- Layer 1 - operating model and routing
- Layer 2 - phase playbooks and rules
- Layer 3 - references and source skills

## Layer 1 - Operating model

### 1. Use one 5-phase workflow

The default workflow is:

1. Platform
2. Route
3. Specify
4. Plan
5. Deliver

Deliver includes:

- Build
- Create PR
- Review PR
- Verify
- Deploy
- Archive

### 2. Roll out in two iterations

Use the methodology in this order:

- Iteration 1: Platform, Route, Specify
- Iteration 2: Plan, Deliver

### 3. Use one change package per approved change

Every approved request should normalize into one change package.

The change package is the canonical execution unit for:

- proposal
- delta specs
- design
- tasks
- delivery state
- verification evidence
- archive history

### 4. Use canonical platform truth + versioned component alignment

When the platform uses one master repository plus many component repositories:

- keep shared platform truth upstream in the platform repo
- keep local OpenSpec artifacts in the component repo
- pin platform version and platform refs in each affected component repo
- use JIRA for issue hierarchy and delivery status, not as the full spec store

Use:

- `platform-ref.yaml` for platform version and platform refs
- `jira-traceability.yaml` for platform issue, component epic, and stories
- a local read-only platform MCP gateway when teams need fast local access to
  platform truth without hosted infrastructure

### 5. Keep humans accountable and agents supportive

Humans own:

- intent
- tradeoffs
- approvals
- release decisions

Agents support:

- drafting
- routing
- artifact generation
- ambiguity checks
- review guidance

### 6. Route by size and impact

Use size to choose planning depth.
Use impact to choose validation and control depth.

Do not mix them.

### 7. Prefer the smallest sufficient workflow

- small work -> compact planning artifacts
- medium work -> standard path
- large or architecture-heavy work -> deeper planning and phased delivery

### 8. Keep delivery reviewable

Deliver in slices.
Each slice should normally produce one reviewable pull request.

### 9. Update artifacts as reality changes

Do not let specs, design, tasks, PR state, or archive drift from what was
actually implemented.

## Layer 2 - Phase router

### Platform

- Owner: Architect
- Goal: define durable context, constitution, config, and role expectations
- Also define the platform truth location, versioning, and alignment conventions
- Use first in Iteration 1
- Primary rules: `rules/platform-rules.md`

### Route

- Owner: Team Lead
- Goal: classify the request, open the change package, and select the next artifact
- Also classify whether the change is component-only or shared with platform truth
- Use for every new change
- Primary rules: `rules/route-rules.md`

### Specify

- Owner: Product
- Goal: define behavior, remove ambiguity, and produce a ready-for-plan spec package
- Also confirm platform refs and whether a linked platform delta is needed
- Finish Iteration 1 here
- Primary rules: `rules/specify-rules.md`

### Plan

- Owner: Architect
- Goal: convert the approved spec into design, tasks, and delivery slices
- Also map tasks to stories and keep platform alignment visible in planning
- Start Iteration 2 here
- Primary rules: `rules/plan-rules.md`

### Deliver

- Owner: Team Lead
- Goal: execute slices through PR, review, verification, deploy, and archive
- Also keep PR, story, epic, and platform alignment traceability current
- Final phase in v1
- Primary rules: `rules/deliver-rules.md`
- PR-specific rules: `rules/pr-review-rules.md`

## Default skill mix by phase

### Platform

- Speckit first
- OpenSpec second
- BMAD third

### Route

- BMAD first
- OpenSpec second
- Speckit only when ambiguity blocks routing

### Specify

- OpenSpec first
- Speckit second
- BMAD third

### Plan

- BMAD first
- OpenSpec second
- Speckit third

### Deliver

- BMAD for implementation and review support
- OpenSpec for apply and archive
- Speckit for task discipline and phased execution

## Output structure when applying this skill

When using this skill in a response, structure the output as:

1. Current phase
2. Goal
3. Owner and support roles
4. Skills to use
5. Rules to apply
6. Artifact or deliverable to produce
7. Exit gate
8. Risks
9. Next step

## Layer 3 - Reference map

### Methodology docs

- `../unified-sdd-methodology/team-proposal.md`
- `../unified-sdd-methodology/iteration-1-playbook.md`
- `../unified-sdd-methodology/iteration-2-playbook.md`
- `../unified-sdd-methodology/canonical-platform-truth-and-component-alignment.md`
- `../unified-sdd-methodology/local-platform-mcp-model.md`
- `../unified-sdd-methodology/example/README.md`

### Rules

- `rules/platform-rules.md`
- `rules/alignment-and-traceability-rules.md`
- `rules/route-rules.md`
- `rules/specify-rules.md`
- `rules/plan-rules.md`
- `rules/deliver-rules.md`
- `rules/pr-review-rules.md`

### References

- `references/overview-and-philosophy.md`
- `references/phase-model.md`
- `references/platform-component-alignment.md`
- `references/agent-interaction-model.md`
- `references/sources.md`
- `../platform-contextualizer-codex-skill/SKILL.md`
- `../platform-truth-mcp-codex-skill/SKILL.md`
- `../explain-code-codex-skill/SKILL.md`

### Role agents

- `agents/README.md`
- `agents/architect-agent.md`
- `agents/team-lead-agent.md`
- `agents/product-agent.md`
- `agents/developer-agent.md`

### Source skills

- `../bmad-codex-skill/SKILL.md`
- `../openspec-codex-skill/SKILL.md`
- `../speckit-codex-skill/SKILL.md`
