---
name: platform-contextualizer
summary: Use for starting the Platform phase of the unified SDD methodology on an existing platform with active teams by documenting current state, identifying gaps, and drafting durable shared context.
triggers:
  - platform contextualizer
  - platform phase
  - current state review
  - platform baseline
  - brownfield platform context
  - document platform context
  - iteration 1 platform
---

# Platform Contextualizer

Use this skill at the very beginning of Iteration 1 when the platform already
exists and teams are already working.

This skill does not assume a blank slate. It assumes:

- the platform has current conventions
- teams have existing habits
- documentation is incomplete or uneven
- some constraints are implicit rather than explicit

Its job is to turn that messy reality into a clear Platform-phase baseline the
teams can actually use.

## What this skill is for

Use this skill to:

- review the platform's current state
- document current conventions, constraints, artifacts, and team interactions
- identify gaps, drift, and missing standards
- draft a practical baseline for the Platform phase
- define where canonical platform truth should live
- define how component repositories should align to the platform truth
- define how JIRA should track platform and component execution
- define whether teams will use a local read-only platform MCP gateway
- prepare the platform for Route and Specify in Iteration 1

## Core idea

The skill combines:

- Sdd-Bmad for brownfield context, role framing, and project-context thinking
- Sdd-OpenSpec for durable context, reusable configuration, and artifact structure
- Sdd-Speckit for explicit principles and quality guardrails

## Workflow

```text
[Observe current state]
        |
        v
[Document what exists]
        |
        v
[Identify gaps and risks]
        |
        v
[Draft durable platform baseline]
        |
        v
[Validate with teams]
        |
        v
Ready for Platform phase
```

## Recommended sequence

1. Capture current state
2. Separate facts from assumptions
3. Identify gaps and pain points
4. Draft durable principles and context
5. Validate with platform stakeholders and team leads
6. Produce Platform-phase outputs

## Default skill mix

### 1. Start with Sdd-Bmad

Use `sdd-bmad` first to:

- treat the platform as brownfield
- inspect existing conventions before proposing new ones
- identify integration points and role boundaries
- frame project-context style outputs

### 2. Use Sdd-OpenSpec second

Use `sdd-openspec` to:

- encode stable context
- structure reusable configuration
- separate durable context from temporary change detail

### 3. Use Sdd-Speckit third

Use `sdd-speckit` to:

- turn vague values into explicit principles
- define quality and testing guardrails
- prepare the baseline for later spec-driven work

## Required outputs

The skill should normally produce:

- current-state snapshot
- platform gap register
- draft platform principles or constitution
- draft reusable project context / config notes
- draft platform versioning and ref model
- draft JIRA hierarchy and issue-link conventions
- draft the local platform MCP usage model when local query and validation is needed
- role and interaction map
- open questions and follow-up actions

## Output structure

When using this skill, structure the output as:

1. Current state
2. Known strengths
3. Gaps and risks
4. Durable context to keep
5. Durable rules to introduce
6. Platform/component alignment model
7. Artifacts to create or update
8. Open questions
9. Recommended next step for Platform phase

## Rules

- `rules/current-state-rules.md`
- `rules/gap-analysis-rules.md`
- `rules/durable-context-rules.md`

## References

- `references/workflow.md`
- `references/artifacts.md`
- `references/sources.md`

## Related skills

- `../sdd-bmad/SKILL.md`
- `../sdd-openspec/SKILL.md`
- `../sdd-speckit/SKILL.md`
- `../platform-truth-mcp-codex-skill/SKILL.md`
- `../unified-sdd-skill/SKILL.md`
