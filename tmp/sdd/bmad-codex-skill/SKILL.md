---
name: bmad-method
summary: Use for spec-driven software planning and implementation using the BMAD Method. Routes work by project scale, brownfield/new-build context, planning phase, agent role, and artifact maturity.
triggers:
  - bmad
  - BMAD Method
  - product brief
  - PRD
  - architecture workflow
  - tech-spec
  - brownfield
  - workflow-init
  - project-context
  - dev-story
  - sprint-planning
  - code-review
---

# BMAD Method

Use this skill when the user wants a structured, agent-friendly way to plan and implement software using BMAD Method concepts: guided workflows, role-specific agents, progressive context, and scale-aware planning.

This skill uses a **3-layer context model**:
- **Layer 1 — Routing and operating rules:** decide the right BMAD track and the next artifact.
- **Layer 2 — Execution playbooks:** create or refine the specific planning/implementation artifact.
- **Layer 3 — References:** consult workflow, agent, brownfield, and project-context guidance.

---

## Layer 1 — Routing and operating rules

### 1. Pick the correct track first
Route the request before writing anything.

Use **Quick Flow** when the work is:
- a bug fix
- a small enhancement
- a small feature with clear scope
- a brownfield change where speed matters

Use **Full BMAD** when the work is:
- a new product
- a major feature
- cross-team or cross-system work
- something that needs stakeholder alignment, PRD depth, or architecture decisions

### 2. Respect the BMAD phase model
BMAD builds context progressively. Work should usually move in this order:
1. discovery / setup
2. planning
3. architecture
4. implementation

Do not jump straight into implementation when the request is still ambiguous and missing core planning context.

### 3. Prefer the smallest sufficient artifact
Choose the lightest artifact that still gives enough context:
- **tech-spec** for quick flow, brownfield, bugs, and small features
- **PRD** for medium-to-large product work
- **architecture** when design impact is material
- **project-context** after architecture, or when onboarding to existing codebases

### 4. Treat brownfield as a first-class case
When the user is working on an existing project:
- anchor recommendations in existing conventions
- identify integration points explicitly
- avoid inventing new patterns without justification
- recommend project-context generation when conventions are unclear

### 5. Agent roles are task-specific
Map the request to the right BMAD persona:
- **PM** for briefs, PRDs, scope, requirements, and tech-specs
- **Architect** for system design, architecture docs, ADRs, and technical tradeoffs
- **Dev** for implementation, tests, and story execution
- **Scrum / planning workflows** for sequencing, story prep, and sprint organization
- **QA / review** for validation and code review outputs

### 6. Output must be implementation-friendly
Every artifact produced with this skill should:
- be concrete, not vague
- identify assumptions and open questions
- include success criteria
- preserve traceability from intent → scope → design → implementation

---

## Layer 2 — Execution playbooks

### A. Workflow router
Use this sequence every time:
1. Determine whether the project is **new-build** or **brownfield**.
2. Estimate scope: **small (1–5 stories)**, **medium (5–15)**, or **large (15+)**.
3. Pick the track:
   - small → usually **Quick Flow / tech-spec**
   - medium → usually **PRD**
   - large or architecture-heavy → **PRD + architecture**
4. Decide the next artifact to produce.
5. Produce only the next artifact unless the user explicitly wants the whole chain.

### B. Quick Flow playbook
Use for bugs, small features, prototypes, and clear-scope brownfield work.

Output a **tech-spec** that includes:
- problem / change summary
- current-state context
- affected components or modules
- constraints and assumptions
- integration points
- implementation approach
- testing approach
- risks / edge cases
- acceptance criteria
- story breakdown when helpful

### C. Full BMAD planning playbook
Use for larger or more ambiguous work.

Typical artifact progression:
1. product brief
2. PRD
3. architecture
4. project-context
5. sprint-planning / create-story
6. dev-story
7. code-review

When generating one of these, preserve continuity with previous artifacts.

### D. Architecture playbook
Create architecture output only when needed. Include:
- system context
- component boundaries
- integration points
- data flow
- operational concerns
- technical decisions and tradeoffs
- ADRs if introducing new patterns
- implementation implications for downstream stories

### E. Project-context playbook
Use when onboarding to an existing repo, after architecture, or when conventions need to be captured.

Project-context should summarize:
- established patterns and conventions
- stack and tooling
- repo structure
- testing approach
- architectural constraints
- naming and style conventions
- integration rules and “how we do things here”

### F. Story execution playbook
When the user asks for implementation-oriented output, produce a story package that includes:
- goal
- scope boundaries
- files or modules likely affected
- implementation notes
- tests to add/update
- definition of done
- review checklist

---

## Layer 3 — Reference map

Read the relevant references before generating output:
- `references/overview-and-philosophy.md` — what BMAD is and why progressive context matters
- `references/workflow-and-tracks.md` — quick flow vs full BMAD and phase sequencing
- `references/agents-and-roles.md` — default agent roles and when to use them
- `references/brownfield-and-project-context.md` — existing-project guidance and project-context generation
- `references/sources.md` — source notes and URLs

---

## Default output rules

- Recommend the **smallest sufficient workflow**.
- For existing projects, prefer **conforming to current conventions** over inventing new structure.
- For small work, bias toward **tech-spec**.
- For larger work, bias toward **PRD**, then architecture only if impact justifies it.
- If the user asks to “implement,” but planning context is missing, produce the missing artifact first or include a compact implementation-ready spec.
- Keep outputs tool-agnostic enough to be used by Codex or other coding agents.

## What good output looks like

A strong BMAD-based answer should:
- say which track was chosen and why
- identify the next artifact
- produce that artifact in a clean template
- include assumptions, constraints, and acceptance criteria
- keep a clear link to implementation
