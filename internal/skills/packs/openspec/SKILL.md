---
name: openspec
description: Spec-driven development from requirements files. Scaffolds changes, writes specs, and tracks implementation. Use when starting from requirements, implementing specs, or managing change proposals.
---

# skill:openspec

## Does exactly this

Guides an agent through spec-driven development: from requirements document through task execution and verification. Phases 0-7 take a user from raw ideas to shipped features with full traceability.

---

## Use this skill when

- Starting a new feature from a requirements file
- User says "openspec", "implement from spec", "start from requirements"
- User provides a requirements/spec file and wants structured implementation

## Do not use this skill when

- The task is a quick fix or single-file change
- There are no requirements to work from

---

## Example prompts

- "Start a project from requirements.md following openspec"
- "Implement the features described in docs/auth-spec.md using openspec"
- "openspec init from .agentic/spec/payment-requirements.md"
- "Continue implementing change auth-feature" (resumes Phase 6)

---

## MANDATORY: Use CLI commands only

**You MUST use `agentic-agent` CLI commands for ALL openspec and task operations.**

Do NOT manually create, edit, or modify any files under `.agentic/` directly.
Do NOT skip CLI commands to "save time" or "be more efficient".
The CLI maintains state consistency — bypassing it causes data corruption.

---

## Steps — in order, no skipping

**Phase 0** — Run `agentic-agent skills ensure` (verify CLI and skills installed)

**Phase 0.5 — REQUIRED** — Get/create formal requirements doc (brainstorm, PRD, or existing file)
→ [See resources/phases.md#phase-05 for full questionnaire and template prompts]

**Phase 1 — REQUIRED** — Read requirements, ask 3-5 clarifying questions, confirm understanding
→ [See resources/phases.md#phase-1 for question template and confirmation format]

**Phase 2** — Run `agentic-agent openspec init "<name>" --from <file>`, fill proposal.md

**Phase 3 — REQUIRED** — Define and document tech stack in tech-stack.md; update proposal.md
→ [See resources/phases.md#phase-3 for questionnaire and tech-stack.md template]

**Phase 4 — REQUIRED** — Create tasks.md + detail files (if 4+ tasks); run `agentic-agent task list` to import
→ [See resources/phases.md#phase-4 for task templates, granularity rules, and file structure]

**Phase 5** — Write specs/ files (data-model.md, api-design.md, architecture.md, ui-design.md) for 4+ task changes

**Phase 6** — Claim → implement → test → complete each task sequentially
→ Optional: use superpowers:executing-plans, superpowers:test-driven-development, superpowers:using-git-worktrees

**Phase 7** — Verify (superpowers or sdd/verifier), complete + archive via CLI

---

## Output

- Proposal file at `.agentic/spec/[change-id]/proposal.md`
- Tech stack at `.agentic/context/tech-stack.md`
- Tasks auto-imported into `.agentic/tasks/backlog.yaml`
- Spec files (optional, for complex changes) in `.agentic/spec/[change-id]/specs/`

---

## Done when

- Phase 0: Skills are installed
- Phase 0.5: Requirements document exists and user approves
- Phase 1: User answers all clarifying questions with confidence
- Phase 2: `proposal.md` file is created and populated
- Phase 3: Tech stack is documented and user confirms it
- Phase 4: All tasks imported, visible in backlog, `task list` succeeds
- Phase 5: Spec files written (if applicable)
- Phase 6: All tasks claimed, implemented, tested, completed
- Phase 7: Final verification passed, change archived

---

## Rules (non-negotiable)

1. Always use `agentic-agent` CLI commands for task and change operations
2. Never write directly to `.agentic/` YAML files
3. Always claim a task before working on it
4. Always complete a task after finishing it
5. Write specs in `specs/` for changes with 4+ tasks or multiple components
6. Do NOT run `openspec import` — tasks auto-import on `task list` or `task claim`

---

## If you need more detail

→ `resources/phases.md` — Full Phase 0-7 workflows with questionnaires, templates, and examples
