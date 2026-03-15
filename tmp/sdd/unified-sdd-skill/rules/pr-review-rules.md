# Pull Request Review Rules

Use these rules inside Deliver whenever a slice is ready for review.

## Goal

Make review a standard part of delivery instead of a local optional habit.

## Must do

- create reviewable PRs
- link each PR back to the change package and relevant tasks
- keep the PR scope aligned to one slice or a small related slice set
- assign the right reviewers
- resolve or explicitly defer review feedback before deploy

## Recommended reviewers

- Team Lead for slice coordination and scope fit
- Architect when the PR touches design integrity, interfaces, or ADR-relevant choices
- Developers for code quality and implementation correctness
- QA / Validation when the PR changes tests, evidence, or release risk
- Product when business behavior needs explicit confirmation

## PR description minimum

- change package reference
- tasks included
- summary of behavior changed
- validation performed
- known risks or deferred items

## Avoid

- giant PRs with mixed concerns
- PRs with no task or change package reference
- review that happens after deploy
- unresolved feedback hidden from the delivery record

## Exit gate

Move from PR review to Verify only when:

- required reviewers have responded
- feedback is resolved or explicitly deferred
- the team agrees the slice is ready for verification
