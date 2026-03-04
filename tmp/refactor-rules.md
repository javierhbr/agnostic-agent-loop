## The Philosophy

Both systems are built on the same single idea:

**An agent that loads everything knows nothing useful.**

When you give an agent a huge file, it doesn't get smarter — it gets noisier. It starts averaging across everything it loaded instead of focusing on the one thing the task actually needs. The output gets generic, inconsistent, and slow.

So the entire philosophy is built around one constraint: **an agent should only ever hold what the current task requires, and nothing else.**

---

## Rules — *how we work here*

Rules define behavior, process, and non-negotiables. The always-on file (AGENTS.md) is tiny on purpose — it only contains a router and hard constraints. When a task arrives, it matches one route and loads one playbook. That playbook has the full detail for that task type. Nothing else enters the context.

**Example:**

Someone asks the agent to build a login form.

```
AGENTS.md sees: UI work
AGENTS.md says: load /docs/playbooks/ui.md
Agent loads: ui.md only
Agent works: with Atomic Design rules, accessibility checklist, responsive behavior
Agent ignores: API rules, DB rules, AI rules — never loaded, never noise
```

The agent didn't become less capable. It became more precise because it carried less.

---

## Skills — *what we produce here*

Skills define a repeatable callable capability — one input, one output, one clear definition of done. A skill file is tiny because it only describes that one thing. Examples, tech rules, and contracts live elsewhere and are only referenced when needed.

**Example:**

Someone triggers `skill:write-spec` for a chat feature.

```
SKILLS.md sees: skill:write-spec
SKILLS.md says: load /skills/write-spec.md
Agent loads: write-spec.md only
Agent produces: /docs/specs/chat.spec.md with Goal, Actor, Requirements, API contract, UI behavior, Constraints
Agent stops: does not implement, does not load ui.md, does not load api.md
```

The skill did exactly one thing and produced exactly one artifact. No drift, no invention, no bloat.

---

## The core principle they share

Both files that are always loaded — AGENTS.md and SKILLS.md — are **routers, not manuals.** Their only job is to say: given what you're doing, go here and only here. Detail lives one level deeper, loaded on demand, never embedded.

This means the agent's context at any moment contains exactly three things:

```
The router  →  the one matched file  →  the task
```

Nothing else ever enters unless the task explicitly requires it.

---

## Why this works at a deeper level

Most agent systems fail because they're written for humans to read, not for agents to load. A human reads a long document once and remembers it. An agent loads it fresh every single time — and every token it loads is a token it has to reason across. The more it loads, the more it reconciles, and the more likely it produces something averaged and imprecise.

This system is written for agents to load, not humans to read. That's the entire design intent behind every decision — short routers, self-contained files, stateless skills, zero inline examples, zero duplicated rules.

**Small context, full precision. That's the whole philosophy.**





Got it. The problem is **token bloat** — every prompt loads too much, agents get noisy, and quality drops. You want a strict, tiny always-on file that forces agents to load *only* what the current task needs.

Here's the system:

---

## The Two Files You Need

---

### FILE 1: `AGENTS.md` — Always-on, always loaded, never grows

```markdown
# AGENTS.md

## Prime Directive
Do not load what you don't need.
Do not invent what is not specified.
Be correct over fast.

## Every single prompt, no exceptions:
1. Read the task
2. Identify the type (one only): → see TASK ROUTER below
3. Load ONLY the file for that type
4. Do the work
5. Report using OUTPUT FORMAT

## TASK ROUTER
UI work        → load /docs/playbooks/ui.md
API work       → load /docs/playbooks/api.md
AI/streaming   → load /docs/playbooks/ai.md
Data/DB        → load /docs/playbooks/data.md
Tests/QA       → load /docs/playbooks/quality.md
Architecture   → load /docs/playbooks/architecture.md

## OUTPUT FORMAT (every response, no exceptions)
1. Task type identified: [one of the above]
2. File loaded: [exact path]
3. What I will do: [1 sentence]
4. Changes: [files touched + why]
5. Verify with: [command or manual step]

## HARD RULES
- Never load more than ONE playbook per task
- Never read entire folders
- Never implement without reading the playbook first
- Never assume behavior not written in a spec
- If the task touches two types → split into two tasks
```

---

### FILE 2: One playbook per task type (loaded only when needed)

Each playbook is **complete and self-contained** — everything the agent needs for that task type, nothing else. Example:

```markdown
# /docs/playbooks/ui.md

## When to load this
Only when the task involves components, layouts, or styling.

## Step by step — do this in order, no skipping
1. Identify component level: atom | molecule | organism | page
2. Check /components/ for existing similar components — reuse first
3. Use Flex or Grid — no absolute positioning
4. No global CSS — encapsulate all styles
5. Props must be typed and minimal — no large object props
6. Cover all states: loading / empty / error
7. If behavior changes → update /docs/specs/<feature>.spec.md

## Done when
- Renders correctly on mobile, tablet, desktop
- No logic duplicated from existing components
- All visual states handled
- No new global styles introduced
```

Each other playbook follows the exact same structure — loaded only when the router sends the agent there.

---

## Why This Works

The `AGENTS.md` is **always tiny** because it contains zero detail. It only contains three things: the directive, the router, and the output format. All detail lives in playbooks that are only loaded when matched. The agent never reads more than two files per task.

---

## Next Step

Share your folder structure (just top-level folders) and I'll generate all the playbooks with your exact paths so the team can drop it in and use it immediately.



# Layered Context System for AI Agents (React + Vercel AI)

## The Problem

When you dump everything into README.md, AGENTS.md, or a system prompt, you get context overflow, slow/inconsistent responses, and agents ignoring critical parts. The fix is a layered context architecture where **the agent knows where to look, not everything at once**.

---

## Core Concept: 3-Level Context Hierarchy

| Level | Purpose | Examples |
|---|---|---|
| Level 1 | Minimal entrypoint | README.md, AGENTS.md |
| Level 2 | Reusable rules | conventions.md, ui-guidelines.md |
| Level 3 | On-demand specs | chat.spec.md, auth.spec.md |

**Golden rule: The agent shouldn't know everything — it should know where to find things.**

---

## Recommended File Structure

```
/docs
  /playbooks
    ui.md
    api.md
    data.md
    vercel-ai.md
    quality.md
    architecture.md
  /context
    architecture.md
    tech-stack.md
    conventions.md
    ui-guidelines.md
  /specs
    README.md          ← specs index
    auth.spec.md
    chat.spec.md
    payments.spec.md
  /adr
    README.md          ← ADR index
    ADR-0001.md

README.md              ← thin entrypoint only
AGENTS.md              ← compact brain
```

---

## README.md — Thin Entrypoint Only

No logic here. Just routing.

```markdown
# Project Overview

This is a React + Vercel AI application.

## For AI Agents

Follow instructions in:
- /docs/agents/AGENTS.md → main rules
- /docs/context/conventions.md → coding rules
- /docs/context/architecture.md → system design
- /docs/specs → feature specs

DO NOT assume behavior not defined in specs.
ALWAYS load the relevant spec before coding.
```

---

## AGENTS.md — Compact Brain

```markdown
# Always-On Rules (MUST FOLLOW)

## Prime Directive
Do not invent requirements. Be correct over fast.

## Workflow
1. State ONE objective (1 sentence)
2. Plan (3–7 steps)
3. Load ONLY the minimum docs using the Task Router below
4. Implement the smallest change set
5. Pass the Quality Gate
6. Reply using the Output Format

## Non-Negotiables
- If behavior/UX/API/data changes → read/create the relevant spec BEFORE coding
- No repo-wide refactors unless explicitly requested
- Minimal context loading: never read entire folders

## Quality Gate (before finishing)
- Types/lint/format pass
- Tests updated/added if behavior changed
- No secrets or sensitive logging

## Output Format
1. Objective
2. Plan
3. Changes (files + brief description)
4. Verification (commands + what to check)

---

## Task Router — Load Only ONE Playbook

A) UI / React / styling / Atomic Design → /docs/playbooks/ui.md
B) API routes / server actions / backend → /docs/playbooks/api.md
C) Data model / DB / migrations        → /docs/playbooks/data.md
D) Vercel AI SDK / streaming / tools   → /docs/playbooks/vercel-ai.md
E) Testing / CI / lint / formatting    → /docs/playbooks/quality.md
F) Architecture / cross-cutting        → /docs/playbooks/architecture.md + /docs/adr/

## Context Tag DSL (optional, for teams)
Prefix any prompt with tags to auto-route:

  task:ui       → loads ui.md
  task:api      → loads api.md
  task:data     → loads data.md
  task:ai       → loads vercel-ai.md
  task:quality  → loads quality.md
  task:arch     → loads architecture.md
```

---

## Playbooks

### /docs/playbooks/ui.md
```markdown
# UI Playbook

## Use when
Creating/updating React components, layouts, styles, responsive behavior.

## Conventions
- App Router: /app/**
- Components: /components/**
- Atomic Design: atoms / molecules / organisms / templates

## Steps
1. Identify level: atom | molecule | organism | page
2. Locate existing similar components and reuse patterns
3. Use Flex/Grid; avoid brittle absolute positioning
4. Encapsulate styling — no global CSS unless requested
5. Add accessibility basics (labels, aria, focus states)
6. Keep props small and typed; avoid large object props
7. If UI behavior changes, update the related spec in /docs/specs/

## Acceptance Checklist
- Responsive on mobile/tablet/desktop
- No duplicated component logic (prefer composition)
- No new global styles
- Visual states covered: loading / empty / error
```

### /docs/playbooks/api.md
```markdown
# API Playbook

## Use when
Implementing /app/api/**, server actions, backend logic, integrations.

## Steps
1. Find the feature spec in /docs/specs/ (create if missing)
2. Define request/response contracts (types + error shape)
3. Implement smallest route/action; keep reusable logic in /lib/**
4. Validate input (zod or existing validator)
5. Handle errors with stable error codes/messages
6. Add tests or a minimal verification path
7. Ensure no secrets logged; use env vars properly

## Acceptance Checklist
- Input validation exists
- Errors are deterministic and typed
- No stack traces leaking to the client
```

### /docs/playbooks/data.md
```markdown
# Data Playbook

## Use when
Changing schema, migrations, DB access patterns, data contracts.

## Steps
1. Read the relevant spec (or create it)
2. Confirm data ownership and source of truth
3. Apply schema change with migration if applicable
4. Update data access layer (/lib/db/** or equivalent)
5. Update types and any dependent API/UI contracts
6. Add/adjust tests for new behavior
7. Provide rollback/forward notes for risky migrations

## Acceptance Checklist
- Migration included if needed
- Backward compatibility considered
- Queries indexed/efficient where relevant
```

### /docs/playbooks/vercel-ai.md
```markdown
# Vercel AI Playbook

## Use when
Using Vercel AI SDK: streaming responses, tools, model calls, chat flows.

## Steps
1. Read the feature spec for the AI behavior (create if missing)
2. Identify where streaming happens (route/handler)
3. Keep prompts minimal — move long rules to specs/docs and reference them
4. Define tool contracts clearly (name, input schema, output)
5. Add guardrails: timeouts, error fallbacks, retries if allowed
6. Ensure no sensitive data is sent to the model unless explicitly permitted
7. Provide a local verification method (manual test steps)

## Acceptance Checklist
- Streaming works end-to-end
- Tool schemas validated
- Errors handled (network / model / tool)
- No prompt bloat in code
```

### /docs/playbooks/quality.md
```markdown
# Quality Playbook

## Use when
Adding/adjusting tests, CI, linting, formatting, type checks.

## Steps
1. Identify project tooling (eslint / prettier / tsc / test runner)
2. Make minimal changes consistent with existing config
3. Prefer adding tests near the changed behavior
4. Ensure deterministic tests (no flaky timers or random values)
5. Update docs if the dev workflow changed

## Acceptance Checklist
- lint and typecheck pass
- Tests pass and cover new behavior
```

### /docs/playbooks/architecture.md
```markdown
# Architecture Playbook

## Use when
Cross-cutting changes, new subsystems, changing conventions, major refactors.

## Steps
1. Write/update an ADR in /docs/adr/ADR-xxxx.md
2. Describe: Context → Decision → Consequences → Alternatives considered
3. Update any affected specs and conventions
4. Implement incrementally (avoid giant PRs)
5. Add verification steps and migration notes

## Acceptance Checklist
- ADR exists and is linked in relevant specs
- Change is incremental and reviewable
```

---

## Index Files

**`/docs/specs/README.md`**
```markdown
# Specs Index
One spec per feature: `<feature>.spec.md`
Specs define behavior and contracts. Code must follow specs, not assumptions.
```

**`/docs/adr/README.md`**
```markdown
# ADR Index
Records architectural decisions. Use ADR-xxxx.md format.
```

---

## Anti-Patterns to Avoid

- A giant AGENTS.md with everything in it
- README containing all business logic
- Rules duplicated across multiple files
- Specs mixed in with code
- Agent loading the entire repo into context

---

## To Tailor This to Your Repo

Share your folder structure (just the top level) and whether/where you use Atomic Design today, e.g.:

```
app/
components/
  atoms/
  molecules/
lib/
docs/
```

That's all that's needed to give you exact paths and a finalized version.'




Perfect. Same philosophy as AGENTS.md — but for writing skills. Here it is:

---

## The Two Files You Need

---

### FILE 1: `SKILLS.md` — Always-on, always loaded, never grows

```markdown
# SKILLS.md

## Prime Directive
A skill must do ONE thing.
A skill must load ONLY what that one thing needs.
A skill must never carry context it doesn't use.

## Every skill, no exceptions:
1. One trigger → one skill
2. One skill → one file
3. One file → one clear output
4. If it needs detail → link to it, never embed it

## SKILL ROUTER
skill:write-spec      → /skills/write-spec.md
skill:create-component → /skills/create-component.md
skill:review-code     → /skills/review-code.md
skill:generate-tests  → /skills/generate-tests.md
skill:write-adr       → /skills/write-adr.md
skill:scaffold        → /skills/scaffold.md

## TRIGGER FORMAT
Human:  skill:write-spec → [what you need]
Agent:  invoke: skill:write-spec, input: [what you need]

## HARD RULES
- One skill per invocation — never chain inside a skill file
- A skill file must never exceed 1 page
- No examples inside skill files — examples live in /skills/examples/
- No tech stack details inside skill files — those live in playbooks
- If a skill needs another skill → the caller chains them, not the skill itself
```

---

### FILE 2: The skill file template — every skill follows this exact structure

```markdown
# skill:[name]

## Trigger
skill:[name]

## Does exactly this
[One sentence. What it produces.]

## Input required
- [field]: [what it is]
- [field]: [what it is]

## Steps — in order, no skipping
1. [step]
2. [step]
3. [step]

## Output
[Exact format of what gets produced — file, text block, structure]

## Done when
- [condition]
- [condition]

## If you need more detail
→ [link to playbook or spec, never inline the content]
```

---

### Example: `/skills/write-spec.md`

```markdown
# skill:write-spec

## Trigger
skill:write-spec

## Does exactly this
Produces a single feature spec file in /docs/specs/

## Input required
- feature: name of the feature
- goal: what it does in one sentence
- actor: who uses it

## Steps — in order, no skipping
1. Check /docs/specs/ — does a spec for this feature already exist?
2. If yes → update it, do not create a duplicate
3. If no → create /docs/specs/[feature].spec.md
4. Fill: Goal, Actor, Requirements, API contract, UI behavior, Constraints
5. Save — do not implement anything

## Output
/docs/specs/[feature].spec.md

## Done when
- File exists in /docs/specs/
- All 6 sections are filled
- No implementation code included

## If you need more detail
→ /docs/playbooks/api.md for API contracts
→ /docs/playbooks/ui.md for UI behavior rules
```

---

## Why This Stays Small

| What bloats a skill | What this system does instead |
|---|---|
| Inline examples | Lives in `/skills/examples/` — loaded only if asked |
| Tech stack rules | Lives in playbooks — linked, never embedded |
| Chained logic | Caller chains skills, skills never call each other |
| Context from other skills | Each skill is stateless — input in, output out |
| Long explanations | One sentence per step, period |

---

## The Mental Model

```
AGENTS.md       → routes to playbooks (how to work)
SKILLS.md       → routes to skills   (what to produce)
Playbook        → detail on HOW
Skill           → detail on WHAT
```

They never overlap. They never duplicate. They only reference each other.

---

Share your actual skill list and I'll write each skill file for you following this exact template.