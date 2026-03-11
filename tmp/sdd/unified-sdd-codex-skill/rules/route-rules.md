# Route Rules

Use these rules in the Route phase.

## Goal

Turn an incoming request into a scoped change package with the right next step.

## Apply these skills in this order

1. `bmad-codex-skill`
2. `openspec-codex-skill`
3. `speckit-codex-skill` only if routing is blocked by ambiguity

## Must do

- classify greenfield or brownfield
- classify size and impact separately
- classify the change as component-only, shared, or platform-rule adoption
- choose the lightest workflow that is still safe
- create or frame one change package
- identify affected platform refs and component repositories
- create the initial JIRA issue chain when the work is material enough to track
- identify the next artifact and owner
- call out unknowns explicitly

## Avoid

- jumping directly into plan or implementation
- mixing size and impact into one score
- hiding integration or architecture risks

## Required outputs

- routed change package
- size and impact classification
- selected path depth
- next artifact
- known unknowns
- initial platform/component alignment metadata
- initial JIRA traceability metadata

## Exit gate

Move to Specify only when:

- the change package exists
- the owner is clear
- the scope is bounded enough to specify
- the next artifact is obvious
