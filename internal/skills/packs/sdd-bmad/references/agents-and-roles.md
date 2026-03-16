# Agents and roles reference

The BMAD docs list default agents and explicitly say that each default BMM agent installs as a **skill**. The agent reference also notes that skill IDs such as `bmad-dev` are used to invoke agents. citeturn1view3

## Practical role mapping for this skill

### PM agent
Use for:
- product brief
- PRD
- tech-spec
- requirement shaping
- scoping and acceptance criteria

### Architect agent
Use for:
- architecture workflow
- ADRs
- integration design
- technical constraints and tradeoffs

### Dev agent
Use for:
- dev-story
- implementation plans
- tests
- code-level execution notes

### Scrum / planning workflows
Use for:
- sprint-planning
- create-story
- sequencing work
- preparing execution packages

### QA / review
Use for:
- code-review
- validation checklists
- quality gates

## Skill-design implication

Because BMAD already models agents as skills, a Codex BMAD skill should act like a router and artifact generator rather than pretending to be every role at once. It should choose the role, then produce the correct artifact or next-step output.
