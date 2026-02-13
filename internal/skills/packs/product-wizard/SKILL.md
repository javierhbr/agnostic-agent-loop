---
name: product-wizard
description: 'Generate robust, production-grade Product Requirements Documents (PRDs) for software products and AI-powered features. Use when asked to "write a PRD", "create product requirements", "document a feature", "document an intent", "plan a feature", "spec out a product", "intent definition", "prd definition", "prd", or when structuring product specifications. Supports full PRDs, lean PRDs, one-pagers, technical PRDs, and AI feature PRDs with executive summaries, user stories, success metrics, and risk analysis.'
---

# Product Wizard

Create comprehensive, actionable Product Requirements Documents that bridge business vision and technical execution. Synthesizes best practices from 11+ product leaders on problem framing, success metrics, behavioral design, and AI-era requirements.

## When to Use This Skill

- Starting a new product or feature development cycle
- Translating a vague idea into a concrete specification
- Defining requirements for AI-powered features (evals, prompts, guardrails)
- Stakeholders need a unified "source of truth" for project scope
- User asks to "write a PRD", "document requirements", "plan a feature", "spec out a product", "define an intent", or "define a PRD"
- Adapting PRD depth: full PRD, lean PRD, one-pager, technical PRD, or AI feature PRD

## Prerequisites

- Understanding of the problem or feature to be documented
- Access to user context: who the users are, what pain exists
- Knowledge of constraints (tech stack, timeline, budget) — or willingness to mark TBD

## Step-by-Step Workflow

### Phase 1: Discovery (The Interview)

Before writing anything, **interrogate the user** to fill knowledge gaps. Never skip this.

Ask 3–7 clarifying questions with lettered options for quick answers:

**Required Discovery Areas:**

1. **The Problem**: What problem does this solve, and why does it matter *now*?
2. **Target Users**: Who is this for? What does their life look like today vs. after?
3. **Success Criteria**: How will you know this worked? What metric moves?
4. **Scope Boundaries**: What is explicitly *out of scope*?
5. **Constraints**: Budget, tech stack, timeline, or regulatory requirements?
6. **Existing Context**: Prior art, competitive products, or internal prototypes?

**Format questions for fast iteration:**

```
1. What is the primary goal?
   A. Improve onboarding   B. Increase retention
   C. Reduce support load  D. Other: [specify]

2. Who is the target user?
   A. New users   B. All users
   C. Enterprise  D. Other: [specify]
```

Users respond with "1A, 2C" for speed. Use 2–4 adaptive rounds.

**Note**: If the user provides a detailed brief upfront, skip redundant questions. Always clarify missing critical info.

### Phase 2: Analysis & Scoping

Synthesize input before drafting:

- Identify dependencies and hidden complexities
- Map the core user flow
- Define non-goals to protect the timeline
- Choose the right PRD format (see below)

### Phase 3: Drafting

Generate the document using the PRD template from `references/prd_template.md`. Follow the schema for the selected format. Present a draft and ask for feedback on specific sections.

### Phase 4: Validation

Run `scripts/validate_prd.sh` or manually verify the self-review checklist before finalizing.

## Core Principles

### Lead with Problem and Context
The most important section is background — what is the problem, why does it matter, and why *now*? If timing can't be justified, the priority is questionable.

### Define Success Before Solutions
Every PRD needs measurable success criteria upfront. "Fast" is not a requirement; "returns results within 200ms for 10k records" is. Consult `references/metrics_frameworks.md` for AARRR, HEART, North Star, and OKR frameworks.

### Keep It Lightweight for Action
Focus on key outcomes, not exhaustive detail. Lightweight PRDs focused on outcomes are more likely to be read and used.

### Prototypes Over Prose (When Appropriate)
For AI and UI features, live prototypes communicate more than documentation. Prompt sets are the new PRDs for AI features.

### Evals as Living PRDs (For AI Products)
Translate requirements into executable evaluations — an eval judge running constantly is the purest form of a PRD.

### Apply Behavioral Science
Consider loss aversion, friction reduction, and defaults when designing requirements. Remove friction from desired paths; add appropriate friction to prevent mistakes.

### Focus on "Why" and "What", Not "How"
Let engineers decide implementation. PRDs articulate the problem and desired outcome.

## PRD Formats

| Format | When to Use | Length |
|--------|-------------|--------|
| **Full PRD** | Major features, new products, strategic initiatives | 4–10 pages |
| **Lean PRD** | Agile features, well-understood problems | 2–3 pages |
| **One-Pager** | Small enhancements, executive briefs | 1 page |
| **Technical PRD** | Engineering-focused, infrastructure | 3–6 pages |
| **AI Feature PRD** | ML/LLM features requiring evals and guardrails | 4–8 pages |

Specify format: "Create a lean PRD for..." or "Generate a technical PRD for..."

## Full PRD Schema

Use the complete template from `references/prd_template.md`. Key sections:

1. **Executive Summary** — Problem, solution, why now, 3–5 measurable KPIs
2. **Background & Context** — Current state, prior art, competitive landscape
3. **User Personas** — Role, goals, frustrations, quotes
4. **User Stories** — `As a [user], I want [action], so that [benefit]` with acceptance criteria. See `references/user_story_examples.md` for patterns.
5. **Functional Requirements** — Numbered (FR-1, FR-2), unambiguous, testable
6. **AI System Requirements** *(if applicable)* — Models, evals, guardrails, human-in-the-loop
7. **Design & UX** *(if applicable)* — Wireframes, accessibility, responsive
8. **Technical Specifications** — Architecture, integrations, performance, security
9. **Success Metrics** — North Star + supporting metrics with targets. See `references/metrics_frameworks.md`.
10. **Risks & Mitigation** — Risk table with likelihood, impact, mitigation
11. **Roadmap & Milestones** — MVP → v1.1 → v2.0 phases
12. **Non-Goals** — Explicit exclusions
13. **Open Questions** — Unresolved items needing decisions

### Requirements Quality Standard

```
# BAD (vague):
- The search should be fast and return relevant results.

# GOOD (concrete):
+ Search returns results within 200ms for 10k records.
+ Search achieves >= 85% Precision@10 in benchmarks.
+ UI achieves 100% Lighthouse Accessibility score.
```

## Lean PRD Schema

1. **Problem & Why Now** (2–3 sentences)
2. **Success Criteria** (3–5 measurable KPIs)
3. **User Stories** (3–7 stories with acceptance criteria)
4. **Scope** (In / Out bullets)
5. **Technical Considerations**
6. **Open Questions**

## One-Pager Schema

1. **Problem** (1 sentence)
2. **Solution** (1 sentence)
3. **Success Metric** (1 KPI with target)
4. **Acceptance Criteria** (3–5 bullets)
5. **Timeline**

## Common Scenarios

**Feature from customer request**: Document verbatim → analyze underlying problem → generalize → validate against strategy → scope.

**Strategic initiative**: Link to OKRs → market analysis → multi-phase rollout → strategy-aligned metrics.

**AI-powered feature**: Lead with user problem → define eval criteria → include prompt sets → specify guardrails → plan human-in-the-loop.

**Technical debt**: Explain user impact → document current limitations with data → define measurable improvements → include engineering input.

**Compliance/regulatory**: Reference specific regulations → include legal review → minimum viable compliance → audit trail requirements.

## Validation Checklist

Before finalizing, verify:

- [ ] Problem is clear — anyone can understand what we're solving
- [ ] "Why now" is justified — timing is explained
- [ ] Users are identified — we know who this is for
- [ ] Success is measurable — concrete KPIs with numeric targets
- [ ] Scope is bounded — clear what's in AND out
- [ ] Requirements are testable — QA can verify completion
- [ ] No vague language — no "fast", "easy", "intuitive" without numbers
- [ ] No placeholder text remains — everything is filled or marked TBD
- [ ] Risks are identified — failure modes considered
- [ ] Open questions are captured — not hidden

Run `scripts/validate_prd.sh <prd_file.md>` for automated checks.

## Common Mistakes to Flag

| Mistake | Fix |
|---------|-----|
| Starting with the solution | Lead with problem and context |
| No success criteria | Define KPIs before features |
| Vague requirements ("fast", "easy") | Use concrete numbers |
| Missing "Why Now" | Justify timing explicitly |
| No out-of-scope section | Define non-goals aggressively |
| Writing "how" not "what" | Focus on outcomes, not implementation |
| Skipping discovery | Always ask 2+ clarifying questions first |
| Assuming constraints | Ask or mark TBD |

## Troubleshooting

| Issue | Solution |
|-------|----------|
| PRD is too long/detailed | Use Lean PRD or One-Pager format |
| Requirements too vague | Add specific examples, concrete numbers, visual references |
| Stakeholders not aligned | Share as draft early, present in person, get explicit sign-off |
| Scope keeps expanding | Use "Non-Goals" section aggressively, separate PRDs for future phases |
| Engineers say not feasible | Involve engineering earlier, be flexible on approach, focus on problem |
| PRD sits unread | Keep lightweight, focus on outcomes, present don't just share |

## References

This skill includes bundled resources:

**references/**
- `prd_template.md` — Full PRD template with all sections and placeholders
- `user_story_examples.md` — User story patterns, INVEST criteria, splitting techniques
- `metrics_frameworks.md` — AARRR, HEART, North Star, OKRs guide with examples

**scripts/**
- `validate_prd.sh` — Validates PRD completeness, user story format, metrics, and scope

Agent Skills specification: https://agentskills.io/specification
