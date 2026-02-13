# Product Wizard: PRD Generation

Generate production-grade Product Requirements Documents from vague ideas. Supports full PRDs, lean PRDs, one-pagers, technical PRDs, and AI feature PRDs.

---

## What You'll Learn

- Use the product-wizard skill to generate PRDs
- Choose the right PRD format for your scope
- Walk through the discovery → analysis → drafting → validation workflow
- Validate PRDs with the bundled validation script
- Feed PRDs into the openspec pipeline

---

## 0. Setup

```bash
# From the project root
go build -o examples/product-wizard/agentic-agent ./cmd/agentic-agent
cd examples/product-wizard

# Initialize a project
rm -rf .agentic agnostic-agent.yaml
./agentic-agent init --name "prd-demo"

# Install the product-wizard skill pack
./agentic-agent skills install product-wizard --tool claude-code
```

Installed files:

```text
.claude/skills/product-wizard/
├── SKILL.md                    # Main skill instructions
├── prd_template.md             # Full PRD template
├── user_story_examples.md      # User story patterns and INVEST criteria
├── metrics_frameworks.md       # AARRR, HEART, North Star, OKRs
└── validate_prd.sh             # Automated validation script
```

---

## 1. Tell Your Agent to Write a PRD

```text
Write a lean PRD for a todo app with offline support and dark mode.
```

The agent activates the product-wizard skill and starts with **Phase 1: Discovery**:

```text
Before I write the PRD, let me clarify a few things:

1. What is the primary platform?
   A. Web (PWA)  B. Native mobile  C. Desktop  D. All

2. Who is the target user?
   A. Individual productivity  B. Small teams  C. Enterprise  D. Other

3. What's the timeline?
   A. MVP in 2 weeks  B. V1 in 1 month  C. Full product in 3 months

4. Any existing tech stack constraints?
   A. React  B. Vue  C. Flutter  D. No preference
```

You respond: `1A, 2A, 3A, 4A`

---

## 2. Agent Generates the PRD

After 2-3 rounds of discovery, the agent writes a **lean PRD**:

```markdown
# Todo App with Offline Support — Lean PRD

## Problem & Why Now
Users need a lightweight task manager that works without internet.
Mobile-first users lose productivity when connectivity drops.
PWA technology now makes offline-first viable without native app overhead.

## Success Criteria
- 90% of core actions work offline (add, edit, complete, delete)
- Page loads in < 2 seconds on 3G
- 50% of users enable PWA install prompt
- Dark mode adoption > 40% of active users

## User Stories
- As a mobile user, I want to add tasks offline so that I don't lose ideas
- As a commuter, I want my changes to sync when I'm back online
- As a night-time user, I want dark mode so the screen doesn't strain my eyes

## Scope
**In:** Task CRUD, offline storage (localStorage), auto-sync, dark mode, PWA install
**Out:** Team collaboration, recurring tasks, integrations, native apps

## Technical Considerations
- React + Vite + TypeScript
- Tailwind CSS with dark mode via class strategy
- Service worker via vite-plugin-pwa
- localStorage for persistence (IndexedDB for future)

## Open Questions
- Should we support sub-tasks in MVP?
- Push notifications when online?
```

---

## 3. Validate the PRD

```bash
# Run the bundled validation script
bash .claude/skills/product-wizard/validate_prd.sh docs/todo-pwa-prd.md
```

```text
╔════════════════════════════════════════╗
║        PRD Validation Report           ║
╚════════════════════════════════════════╝

━━━ Required Sections ━━━
✓ Problem Statement found
✓ Goals & Objectives found
✓ User Stories found
✓ Success Metrics found
✓ Scope found

━━━ Content Quality ━━━
✓ Placeholder text passed
✓ TBD markers passed

━━━ Vague Language Check ━━━
✓ No vague language detected

✅ Validation passed — ready for stakeholder review
```

---

## 4. Feed PRD into OpenSpec

Once the PRD is approved, use it as the source for an openspec change:

```bash
agentic-agent openspec init "Todo PWA" --from docs/todo-pwa-prd.md
```

This creates the change directory with proposal and tasks templates seeded from the PRD content. Fill in the tasks, and they auto-import into the backlog.

---

## 5. PRD Formats

| Format | Prompt | When |
|--------|--------|------|
| **Full PRD** | "Write a full PRD for..." | New products, strategic initiatives |
| **Lean PRD** | "Write a lean PRD for..." | Agile features, well-understood problems |
| **One-Pager** | "Write a one-pager for..." | Small enhancements, executive briefs |
| **Technical PRD** | "Write a technical PRD for..." | Infrastructure, engineering-focused |
| **AI Feature PRD** | "Write an AI feature PRD for..." | ML/LLM features needing evals |

---

## End-to-End Pipeline

```text
product-wizard (PRD) → openspec (tasks) → atdd (tests) → run-with-ralph (implementation)
```

1. **product-wizard** generates the PRD with measurable criteria
2. **openspec** breaks it into scoped tasks with acceptance criteria
3. **atdd** writes executable tests from acceptance criteria
4. **run-with-ralph** iterates until tests pass

---

## Quick Reference

| Action | Command / Prompt |
|--------|-----------------|
| Generate PRD | Tell agent: "Write a [format] PRD for [topic]" |
| Validate PRD | `bash .claude/skills/product-wizard/validate_prd.sh <file>` |
| Feed to OpenSpec | `openspec init "Name" --from <prd-file>` |
| Install skill | `skills install product-wizard --tool claude-code` |
| Check drift | `skills check` |
