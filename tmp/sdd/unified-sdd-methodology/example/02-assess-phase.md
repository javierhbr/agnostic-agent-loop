# Example: Assess Phase

Goal:

- turn an incoming request into one assessed change package with the right scope

## Example scenario

An incoming request says: customers must be able to update their email address with validation and clear failure feedback.

## Flow

```text
[Incoming request]
  product requirement
        |
        v
[Team Lead assesses]
        |
        +--> BMAD: classify size, impact, and path
        +--> OpenSpec: open change package and issue chain
        +--> Speckit: clarify only if assessment is blocked
        |
        v
[Assessed change]
  PLAT-123 + PROF-456 + AUTH-234
```

## `Team Lead`

When to act:

- first for every new request

Skills in order:

1. `bmad-skill`
2. `openspec-skill`
3. `speckit-skill` only if assessment is blocked
4. `explain-code-skill` if code-path impact is unclear

Example prompts:

- "Using the BMAD skill, classify this validated-email-update request by size, impact, and architecture depth, and tell me whether it is platform-only, component-only, or shared."
- "Using the OpenSpec skill, open the change package for validated customer email updates, identify the affected component repos, and draft the initial platform-ref.yaml and jira-traceability.yaml."
- "Using the OpenSpec skill, create the initial delivery chain for PLAT-123, PROF-456, and AUTH-234, and record the next artifact and owner."
- "Using the explain-code skill, explain the current blast radius of an email update across profile-service and auth-service, and call out one assessment risk."

Expected outputs:

- assessed change package
- platform issue and component epics
- initial `platform-ref.yaml`
- initial `jira-traceability.yaml`

## `Product`

When to act:

- during Assess when business value and non-goals affect scope

Skills in order:

1. `openspec-skill`
2. `explain-code-skill`

Example prompts:

- "Using the OpenSpec skill, clarify the business goal, urgency, and non-goals for validated customer email updates so the request can be assessed safely."
- "Using the OpenSpec skill, identify which customer-facing behaviors must be in scope for the first release and which ones must be deferred."
- "Using the explain-code skill, explain the current behavior to product and highlight the one behavior that matters most for scoping this request."

Expected outputs:

- bounded business scope
- explicit non-goals
- product input to assessment decision

## `Architect`

When to act:

- during Assess when contracts or cross-team dependencies may change

Skills in order:

1. `bmad-skill`
2. `explain-code-skill`

Example prompts:

- "Using the BMAD skill, assess whether validated customer email updates require a shared platform change package, a contract update, or only local component changes."
- "Using the BMAD skill, identify the main architectural dependencies and cross-team risks for PLAT-123."
- "Using the explain-code skill, explain the current system path for customer email updates and highlight one dependency that increases impact."

Expected outputs:

- architecture impact decision
- dependency list
- shared vs local change recommendation

## `Developer`

When to act:

- during Assess only when feasibility or code coupling is unclear

Skills in order:

1. `speckit-skill`
2. `explain-code-skill`

Example prompts:

- "Using the Speckit skill, identify the technical unknowns, edge cases, or hidden dependencies that could change the assessment decision for validated customer email updates."
- "Using the explain-code skill, explain the current implementation path for email updates and call out one hidden coupling risk."

Expected outputs:

- technical unknowns
- early feasibility notes

## Output targets

Use these concrete artifacts as the Assess target:

- `component-repo/platform-ref.yaml`
- `component-repo/jira-traceability.yaml`
