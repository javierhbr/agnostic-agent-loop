# Workflow and Commands

## Primary command flow
1. `/speckit.constitution`
2. `/speckit.specify`
3. `/speckit.clarify`
4. `/speckit.checklist` (optional but valuable)
5. `/speckit.plan`
6. `/speckit.tasks`
7. `/speckit.analyze`
8. `/speckit.implement`

## Command intent
### `/speckit.constitution`
Create or update governing principles for the project.
Use for code quality, testing standards, architecture preferences, UX consistency, security posture, and performance expectations.

### `/speckit.specify`
Describe the feature in terms of user value, business rules, constraints, and acceptance expectations.
Avoid technical stack details here.

### `/speckit.clarify`
Resolve ambiguity before planning. Use it to expose hidden decisions and turn vague requirements into explicit ones.

### `/speckit.checklist`
Generate a validation checklist for requirements completeness, clarity, and consistency.

### `/speckit.plan`
Create the technical implementation strategy: architecture, stack, data flow, integration points, and validation approach.

### `/speckit.tasks`
Generate concrete tasks from the plan. Keep them ordered, scoped, and traceable.

### `/speckit.analyze`
Check for coverage gaps, inconsistencies, and missing work across the generated artifacts.

### `/speckit.implement`
Execute the tasks. Prefer phased implementation for large work.

## Practical execution rules
- do not plan before the spec is sufficiently clear
- do not implement from a weak plan
- do not skip analyze for large or risky work
- for big features, split implementation into validated phases
