---
name: product-wizard
description: Generate robust Product Requirements Documents. Use when asked to write a PRD, create product requirements, document a feature, or plan a feature for software products and AI-powered features.
---

# skill:product-wizard

## Does exactly this

Creates comprehensive, actionable PRDs that bridge business vision and technical execution. Synthesises best practices on problem framing, success metrics, and requirements.

---

## When to use

- Starting a new product or feature development cycle
- Translating a vague idea into a concrete specification
- Defining requirements for AI-powered features (evals, prompts, guardrails)
- Stakeholders need a unified "source of truth"
- User asks to "write a PRD", "document requirements", "plan a feature", or "spec out a product"

---

## Output Location

Save PRD files to `.agentic/spec/`:

```
mkdir -p .agentic/spec
.agentic/spec/prd-<feature-name>.md
```

Then use with openspec:
```bash
agentic-agent openspec init "<feature-name>" --from .agentic/spec/prd-<feature-name>.md
```

---

## Step-by-Step Workflow

### Phase 1: Discovery

**Before writing, ask 3–7 clarifying questions** with lettered options for quick answers:

**Required Discovery Areas:**
1. **The Problem** — What problem does this solve? Why now?
2. **Target Users** — Who is this for?
3. **Success Criteria** — How will you know this worked?
4. **Scope Boundaries** — What is explicitly OUT of scope?
5. **Constraints** — Budget, tech stack, timeline, regulations?
6. **Existing Context** — Prior art, prototypes, related products?

### Phase 2: Analysis & Scoping

- Identify dependencies and hidden complexities
- Map the core user flow
- Define non-goals to protect timeline
- Choose the right PRD format (Full, Lean, One-Pager, Technical, AI Feature)

### Phase 3: Drafting

Generate using the template. See `references/prd_template.md` for full schema.

Save to `.agentic/spec/prd-<feature-name>.md`.

### Phase 4: Validation

Before finalizing, verify all checkpoints (see below).

---

## PRD Formats

| Format | When to Use | Length |
|--------|-------------|--------|
| **Full PRD** | Major features, new products, strategic initiatives | 4–10 pages |
| **Lean PRD** | Agile features, well-understood problems | 2–3 pages |
| **One-Pager** | Small enhancements, executive briefs | 1 page |
| **Technical PRD** | Engineering-focused, infrastructure | 3–6 pages |
| **AI Feature PRD** | ML/LLM features requiring evals, guardrails | 4–8 pages |

---

## Core Principles

**Lead with Problem and Context** — The most important section is background. If timing can't be justified, the priority is questionable.

**Define Success Before Solutions** — Every PRD needs measurable success criteria upfront. "Fast" is not a requirement; "returns results within 200ms for 10k records" is.

**Keep It Lightweight** — Focus on key outcomes, not exhaustive detail. Lightweight PRDs focused on outcomes are more likely to be read.

**No Vague Language** — Never write "fast", "easy", or "intuitive" without concrete numbers. See `references/metrics_frameworks.md` for AARRR, HEART, North Star, and OKR frameworks.

**Requirements Must Be Testable** — QA can verify completion. Example: "Search returns results within 200ms for 10k records" (testable), not "search is fast" (vague).

---

## Validation Checklist

Before finalizing, verify:

- [ ] Problem is clear
- [ ] "Why now" is justified
- [ ] Users are identified
- [ ] Success is measurable (concrete KPIs)
- [ ] Scope is bounded (clear what's in AND out)
- [ ] Requirements are testable
- [ ] No vague language
- [ ] No placeholder text remains
- [ ] Risks are identified
- [ ] Open questions are captured

---

## Common Mistakes to Avoid

| Mistake | Fix |
|---------|-----|
| Starting with the solution | Lead with problem and context |
| No success criteria | Define KPIs before features |
| Vague requirements | Use concrete numbers |
| Missing "Why Now" | Justify timing explicitly |
| No out-of-scope section | Define non-goals aggressively |
| Writing "how" not "what" | Focus on outcomes, not implementation |
| Skipping discovery | Always ask clarifying questions first |

---

## If you need more detail

→ `references/prd_template.md` — Full PRD template with all sections and placeholders
→ `references/user_story_examples.md` — User story patterns, INVEST criteria, splitting techniques
→ `references/metrics_frameworks.md` — AARRR, HEART, North Star, OKRs guide with examples
