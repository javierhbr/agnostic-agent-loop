# Example: Entry Points for Changes

The methodology allows multiple entry points, but all of them normalize into one change package model.

## Entry point 1: Platform initiative

Use when:

- a shared capability or shared contract changes
- multiple teams or repositories are affected

Example:

- `PLAT-123`: validated customer email updates across the platform

Who starts:

- `Architect` + `Product`

Skills and prompts:

- `Architect` with `BMAD`
  - "Using the BMAD skill, assess the platform-level impact of this initiative and identify the shared capabilities, contracts, and teams involved."
- `Product` with `OpenSpec`
  - "Using the OpenSpec skill, define the platform-level business outcome and the high-level acceptance expectations for this initiative."
- `Team Lead` with `OpenSpec`
  - "Using the OpenSpec skill, create the initial platform issue and component epic chain for the impacted repositories."

Expected outputs:

- platform issue
- affected component epic list
- initial shared refs

## Entry point 2: Product requirement

Use when:

- a business request starts from product or customer needs
- the affected components are not fully known yet

Example:

- customers must be able to update their email with validation and clear error feedback

Who starts:

- `Product`

Skills and prompts:

- `Product` with `OpenSpec`
  - "Using the OpenSpec skill, define user stories, acceptance criteria, and non-goals for validated customer email updates."
- `Team Lead` with `BMAD` then `OpenSpec`
  - "Using the BMAD skill, classify the requirement by size, impact, and affected repositories, then use the OpenSpec skill to open the change package."
- `Architect` with `BMAD`
  - "Using the BMAD skill, assess whether this product requirement changes shared contracts or only local component behavior."

Expected outputs:

- user stories
- initial scope
- affected component list

## Entry point 3: Component or team proposal

Use when:

- a team sees a technical opportunity, refactor, or component-level enhancement

Example:

- profile-service should add stricter email validation before sending profile update events

Who starts:

- `Team Lead` or `Developer`

Skills and prompts:

- `Team Lead` with `BMAD`
  - "Using the BMAD skill, assess whether this component proposal is local-only or whether it affects shared platform contracts."
- `Team Lead` with `OpenSpec`
  - "Using the OpenSpec skill, create a component change package for profile-service and record the current platform refs."
- `Developer` with `Speckit`
  - "Using the Speckit skill, identify the executable behavior changes and edge cases this component proposal introduces."
- `Architect` with `BMAD`
  - "Using the BMAD skill, review whether the proposal introduces design drift from platform principles."

Expected outputs:

- local change package
- platform alignment check
- component spec

## Entry point 4: Bug fix

Use when:

- a defect is known and speed matters

Example:

- email update fails when the customer address contains uppercase characters

Who starts:

- `Team Lead` or `Developer`

Skills and prompts:

- `Team Lead` with `BMAD`
  - "Using the BMAD skill, classify this bug fix by size and impact and decide whether it can take the compact path."
- `Developer` with `Explain Code`
  - "Using the explain-code skill, explain the current failing code path and the likely root cause with one gotcha."
- `Product` with `OpenSpec`
  - "Using the OpenSpec skill, define the corrected expected behavior and the minimum acceptance criteria for the bug fix."
- `Architect` with `BMAD`
  - "Using the BMAD skill, assess whether this bug fix touches shared contracts or can stay entirely local."

Expected outputs:

- bug-fix assessment decision
- compact spec or full spec
- impacted component and platform refs
