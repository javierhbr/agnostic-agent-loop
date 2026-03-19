# Example: Platform Phase

Goal:

- create shared context and durable rules before change-level work starts

## Example scenario

The platform wants to support validated customer email updates across multiple components.
Before any component team writes specs or code, the platform needs a common baseline.

## Flow

```text
[Existing platform]
  current repos + contracts + team conventions
        |
        v
[Architect leads Platform]
        |
        +--> BMAD: inspect brownfield architecture and constraints
        +--> OpenSpec: encode durable context and refs
        +--> Speckit: turn principles into explicit rules
        |
        v
[Platform baseline]
  principles + versioning + refs + JIRA conventions
```

## `Architect`

When to act:

- first in Platform
- whenever the platform baseline is missing or stale

Skills in order:

1. `bmad-skill`
2. `openspec-skill`
3. `speckit-skill`
4. `explain-code-skill` when teams need the current architecture explained

Example prompts:

- "Using the BMAD skill, review the current platform as a brownfield system, identify the architectural constraints for customer identity flows, and list the integration points across profile-service, auth-service, and notification-service."
- "Using the OpenSpec skill, turn those platform constraints into durable context, versioning rules, and reusable refs for the customer-identity capability."
- "Using the Speckit skill, convert the customer identity platform principles into explicit, testable rules for validation, contracts, observability, and security."
- "Using the explain-code skill, explain the current customer identity architecture with an analogy, an ASCII diagram, a step-by-step walkthrough, and one architecture gotcha."

Expected outputs:

- platform baseline
- capability refs
- contract refs
- platform versioning model
- JIRA conventions

## `Product`

When to act:

- during Platform when business context must become durable
- before detailed feature stories are written in Specify

Skills in order:

1. `openspec-skill`
2. `speckit-skill`
3. `explain-code-skill` when current behavior needs clarification

Example prompts:

- "Using the OpenSpec skill, document the durable business context for validated customer email updates, including business goals, customer expectations, and the likely impacted components."
- "Using the OpenSpec skill, identify the customer experiences and components affected by validated email updates, and record the high-level acceptance expectations that should influence later specs."
- "Using the Speckit skill, turn the business expectations for email updates into explicit rules for acceptance, failure behavior, and customer communication."
- "Using the explain-code skill, explain how the platform behaves today when a customer updates an email address, and highlight one gap the future feature must address."

Expected outputs:

- durable business constraints
- impacted component list
- platform-level acceptance language

## `Team Lead`

When to act:

- during Platform when delivery realities must become part of the baseline

Skills in order:

1. `bmad-skill`
2. `openspec-skill`
3. `explain-code-skill`

Example prompts:

- "Using the BMAD skill, document the current team conventions, repo boundaries, and delivery constraints that affect profile-service, auth-service, and notification-service."
- "Using the OpenSpec skill, capture the durable workflow expectations for platform-ref.yaml, jira-traceability.yaml, and component-level OpenSpec change packages."
- "Using the explain-code skill, explain the current handoff flow between platform rules, component planning, and PR review using an analogy, an ASCII diagram, a walkthrough, and one coordination gotcha."

Expected outputs:

- team constraints
- adoption notes for component repos
- handoff rules for platform and component work

## Output target

Use `platform-repo/platform-baseline.md` as the concrete target for this phase.
