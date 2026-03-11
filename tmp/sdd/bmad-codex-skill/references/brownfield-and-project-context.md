# Brownfield and project-context reference

The BMAD docs include a dedicated guide for established projects. That guide says to use BMAD effectively on existing codebases by grounding work in current documentation and conventions, not by replacing them. It also recommends installation plus an AI IDE such as Claude Code or Cursor as prerequisites. citeturn1view4

## Brownfield feature flow
The feature-addition guide recommends:
1. run `workflow-init`
2. choose approach by scope
3. create the planning document
4. add architecture work only if the feature affects system architecture
5. follow implementation workflows such as `sprint-planning`, `create-story`, `dev-story`, and `code-review` citeturn1view7

## Project context
The project-context docs say to run `bmad-generate-project-context` after architecture, or on existing projects, so the workflow scans the architecture and project files to capture established decisions and conventions. The docs explain that this matters because without project context, agents default to generic patterns and may miss project-specific constraints. citeturn1view5

## Codex-skill implication
When working with an existing repo:
- inspect established patterns first
- call out integration points
- avoid introducing new conventions casually
- recommend or generate a project-context summary when conventions are unclear
