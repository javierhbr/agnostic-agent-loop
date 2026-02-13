```markdown
---
name: robust-prd-wizard
description: >
  Creates robust, execution-ready Product Requirement Documents (PRDs) for human teams and AI agents.
  Use this skill when you need to define or refine a feature, product, or system with clear goals,
  behavioral design, and (optionally) an agent-executable PRD format.
---

# Robust PRD Wizard

## What this skill does

- Interactively gathers the minimum critical context to write a high-quality PRD (problem, users, goals, scope, constraints). [web:2][web:6]  
- Supports multiple PRD formats: classic human-readable PRD, lean one-pager, technical PRD, and agent-executable (ralph-tui style) PRD. [page:2][page:4][page:5]  
- Optionally layers in behavioral product design so the PRD explicitly targets behavior change and applies behavioral principles. [page:1]  

Use this skill whenever you want a repeatable, opinionated workflow for turning fuzzy feature ideas into clear, shippable requirements.

---

## How the agent should think

Always follow this three-phase workflow:

1. **Discovery first, output later**  
   Never jump straight into writing a PRD.  
   Ask targeted questions to clarify:
   - Problem and “why now”  
   - Target users and target behaviors  
   - Business goals and success metrics  
   - Scope, non-goals, constraints  
   - For agent mode: quality gates and commands to run  
   The goal is to capture just enough detail to produce a high-signal PRD without overwhelming the user. [web:2][page:2][page:4][page:5]  

2. **Select PRD mode**  
   Based on the user’s answers, choose one PRD mode:
   - **Standard PRD** – for substantial features/products.  
   - **Lean PRD** – for minor changes or incremental improvements.  
   - **Technical / Design-heavy PRD** – for infra, refactors, or UX-driven work.  
   - **Agent PRD (ralph-tui)** – when the user wants an agent to implement work from the PRD. [page:2][page:3][page:4][page:5]  

   If unclear, ask:  
   > “What format fits best: standard, lean, technical, or agent-executable (ralph) PRD?”

3. **Decide whether to add behavioral design**  
   If the work touches activation, engagement, retention, monetization, or recurring behaviors, enable the behavioral overlay. [page:1]  
   Otherwise, skip it to keep the PRD concise.

---

## Step 1 – Discovery questions

Ask questions in small batches, then summarize back what you heard before moving on.

### 1. Problem and context

- “In 2–3 sentences, what problem are we solving?”  
- “Why is this important to solve now instead of later?”  
- “What existing solutions or workarounds do people use today?” [page:3][page:4]  

### 2. Users and behaviors

- “Who is this for? (Primary users, secondary users)”  
- “What specific behavior do you want to encourage, change, or prevent?” [page:1]  
- “Where in the current flow do users get stuck or drop off?”  

### 3. Goals and success metrics

- “What business or product goals should this PRD support?”  
- “How will we know this worked? What metrics should move (and by how much, if known)?” [page:3][page:4]  

### 4. Scope, non-goals, constraints

- “What’s explicitly **in scope** for this iteration?”  
- “What’s explicitly **out of scope** or a ‘future iteration’?” [page:3][page:4][page:5]  
- “Any deadlines, legal/privacy/security constraints, or tech constraints?” [page:2][page:4]  

### 5. Agent execution (if requested)

Only ask these if the user mentions agents, ralph-tui, or automated implementation.

- “Which commands must pass for every story? (e.g., `pnpm typecheck && pnpm lint`, `npm test`, etc.)” [page:5]  
- “Should UI changes be verified in a browser (e.g., via dev-browser)? (Yes/No)” [page:5]  
- “Are there particular directories, frameworks, or patterns the agent must follow?” [page:5]  

Summarize back in a short bullet list and confirm before writing the PRD.

---

## Step 2 – Mode selection

After discovery, decide the PRD mode.

### If the user does not specify

Infer mode from context:

- Large project / new product → **Standard PRD**  
- Small tweak / clear change → **Lean PRD**  
- Deep infra / platform / architecture / heavy UX → **Technical/Design PRD**  
- User explicitly wants code agents to implement work → **Agent PRD (ralph-tui)** [page:2][page:3][page:4][page:5]  

Ask for confirmation if ambiguous.

---

## Step 3 – PRD schemas

### A. Standard PRD (human-readable)

Use this structure: [page:2][page:3][page:4]  

1. **Executive summary**  
   - Problem (1–2 sentences).  
   - Proposed solution (1–2 sentences).  
   - Top 3–5 success metrics.  

2. **Problem, context, and why now**  
   - Background and current state.  
   - Pain points and opportunities.  
   - Why this is important at this moment.  

3. **Goals and objectives**  
   - Business goals.  
   - Product objectives tied to metrics.  

4. **Users and use cases**  
   - Primary and secondary personas.  
   - Key use cases / jobs-to-be-done.  

5. **User stories and acceptance criteria**  
   - `As a [user], I want [action] so that [benefit].`  
   - Acceptance criteria that are concrete and testable. [page:2][page:4]  

6. **Behavioral design (optional overlay)**  
   Include if behavior change is important: [page:1]  
   - Target behaviors and desired outcomes.  
   - Barriers (friction, uncertainty, status quo bias, present bias).  
   - Interventions (defaults, loss aversion, habit loops, celebration moments).  

7. **Functional requirements**  
   - Numbered FR-1, FR-2… with clear, testable descriptions.  

8. **Non-functional requirements**  
   - Performance, security, privacy, reliability, accessibility, compatibility. [page:2][page:4]  

9. **Scope and non-goals**  
   - Explicit in-scope/out-of-scope lists.  

10. **Dependencies and assumptions**  

11. **Risks and mitigations**  

12. **Timeline and milestones**  

13. **Open questions**  

---

### B. Lean PRD

Use this when the user wants a one-pager or a small feature spec. [page:3][page:4]  

Sections (all concise):

1. Problem and why now  
2. Target users and main use case  
3. Goals and success metrics  
4. Proposed solution (high level, 3–7 bullets)  
5. Scope and non-goals  
6. Risks / unknowns  
7. Open questions  

Behavioral notes can be collapsed into 1–2 bullets under “Proposed solution” if relevant.

---

### C. Technical / Design-heavy PRD

Add extra depth for engineering/UX work. [page:2][page:3][page:4]  

1. Executive summary  
2. Problem, goals, constraints  
3. Detailed technical context (architecture, current limitations)  
4. Proposed technical approach (patterns, systems, interfaces)  
5. Data model and API changes (if any)  
6. UX flows / IA / design requirements (if applicable)  
7. Risks, trade-offs, alternatives considered  
8. Migration, rollout, and observability plan  

Keep “what and why” primary, but it’s acceptable to outline candidate approaches when the user explicitly asks for technical guidance.

---

### D. Agent-executable PRD (ralph-tui style)

Use only when the user wants a PRD that an agent will execute (e.g., ralph-tui). [page:5]  

Wrap the entire PRD with:

```text
[PRD]
...
[/PRD]
```

Inside, use:

1. **Overview**  
   - Short description and problem.  

2. **Goals**  
   - Measurable objectives.  

3. **Quality gates (required)**  
   - Global commands that must pass for every story (e.g., `pnpm typecheck && pnpm lint && pnpm test`).  
   - UI verification rules if requested. [page:5]  

4. **User stories**  
   For each story:  
   - Title.  
   - Description.  
   - Checklist-style acceptance criteria (independent, machine-verifiable). [page:5]  

5. **Functional requirements**  
   - Numbered FR-1, FR-2…  

6. **Non-goals**  

7. **Technical considerations** (only if needed):  
   - Frameworks, directories, patterns to follow. [page:5]  

8. **Success metrics**  

9. **Open questions**  

Guidelines:

- Do not repeat quality gates within individual stories; define once globally. [page:5]  
- Keep each story small enough for a single agent run. [page:5]  
- Reference specific files/paths when known to minimize ambiguity. [page:5]  

---

## Behavioral overlay (when enabled)

Insert or expand a “Behavioral product design” section into any PRD mode. [page:1]  

Include:

- Target behavior(s).  
- Current barriers (friction, confusion, loss of motivation, etc.).  
- Selected principles (loss aversion, defaults, habit loops, protective friction, social proof, etc.).  
- How each principle is applied in the solution.  
- Ethical guardrails (no dark patterns, respect for user autonomy).  

---

## Style and quality rules

- Prefer concise bullet points and numbered lists over long prose. [web:2]  
- Avoid vague adjectives like “simple”, “fast”, “intuitive”; replace with measurable thresholds or examples. [page:2][page:4]  
- Keep `SKILL.md` under ~500 lines and avoid time-sensitive details. [web:2][web:11]  
- Focus on **what** and **why**; leave detailed implementation decisions to engineers or agents unless the user explicitly wants implementation guidance. [page:2][page:4]  

---

## How to respond to the user

1. Ask 5–10 discovery questions in small batches.  
2. Summarize what you learned and propose a PRD mode (standard, lean, technical, or agent).  
3. Confirm the mode and whether to apply the behavioral overlay.  
4. Generate the PRD following the chosen schema.  
5. Offer to revise specific sections (goals, metrics, scope, user stories, behavioral design) based on feedback.

```