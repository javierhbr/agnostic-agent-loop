# Quick Reference Card: Skills for Idea Refinement

*Keep this handy when talking to AI*

---

## The Four Skills

| Skill | Purpose | When to Use |
|-------|---------|-------------|
| 🧙 **product-wizard** | Vague → Structured PRD | You have an idea but it's fuzzy |
| 📋 **dev-plans** | PRD → Phased tasks | You have a PRD and need a plan |
| ✅ **atdd** | Criteria → Tests | You want to verify requirements |
| 🔍 **code-simplification** | Complex → Simple | Specs are too complicated |

---

## Essential Prompts (Copy & Adapt)

### Starting from Vague Idea

```
Use the product-wizard skill to create a PRD for [idea].

Problem: [what problem?]
Users: [who needs this?]
Current state: [how they do it now?]
Success: [what outcome do you want?]
```

---

### Organizing Messy Requirements

```
Use product-wizard to organize this requirements list:

[paste your bullets]

Please resolve conflicts, clarify vague items, and separate must-haves
from nice-to-haves.
```

---

### Creating Implementation Plan

```
Use dev-plans skill to break this PRD into phases with checkpoints.

[paste PRD or reference location]
```

---

### Making Requirements Testable

```
Use atdd skill to show how we'd test these acceptance criteria:

- [Criterion 1]
- [Criterion 2]

Show test scenarios for each.
```

---

### Simplifying Complex Specs

```
Use code-simplification skill to review this requirement:

"[paste complex requirement]"

Is this too complex? Can it be simpler?
```

---

## The Discovery Questions AI Will Ask

### From product-wizard

1. **Problem:** What specific problem? Why does it matter?
2. **Users:** Who needs this? Primary user segment?
3. **Success:** How will you measure success? What metrics?
4. **Constraints:** Timeline? Budget? Technical limits?
5. **Scope:** What's IN vs OUT for first release?

### Your Answers Should Include

- **Specific numbers:** "50K users" not "many users"
- **Real metrics:** "3 days → 4 hours" not "faster"
- **Evidence:** "30% of support tickets" not "users complain"
- **Clear boundaries:** "NOT including SMS" not "everything"

---

## Red Flags = Stop and Clarify

| You Said | Problem | Fix |
|----------|---------|-----|
| "Fast" | Too vague | "Under 200ms for 95% of requests" |
| "Easy to use" | Not measurable | "Complete task in ≤3 clicks" |
| "Everyone needs this" | No prioritization | "Project Managers are priority" |
| "ASAP" | No real timeline | "Launch by Q2 end (10 weeks)" |
| "Make it better" | No success criteria | "Reduce failures from 25% to 10%" |

---

## The Refinement Loop

```
1. Start → Vague idea
2. product-wizard → Ask questions
3. You → Answer with specifics
4. AI → Draft PRD
5. You → Review
6. code-simplification → Clarify complex parts
7. dev-plans → Break into phases
8. atdd → Verify it's testable
9. Done → Clear, actionable spec
```

---

## Success Checklist

After using the skills, you should have:

- [ ] **Clear problem statement** - Not "we need X" but "users can't Y because Z"
- [ ] **Specific users** - Named segments, not "everyone"
- [ ] **Measurable success** - Numbers with targets, not "better"
- [ ] **Documented constraints** - Time, budget, tech limits
- [ ] **Explicit scope** - IN and OUT lists
- [ ] **Testable criteria** - Can verify if it works
- [ ] **Simple language** - No jargon without explanation

---

## Time Estimates

| Task | Time |
|------|------|
| product-wizard discovery | 20-30 min |
| PRD generation | 5 min (automatic) |
| dev-plans creation | 5-10 min |
| atdd for one user story | 5 min |
| code-simplification review | 5 min |
| **Total for complete spec** | **40-60 min** |

Compare to traditional: Days or weeks of meetings and emails.

---

## Common Mistakes

### ❌ Don't

- Rush through AI questions
- Assume AI knows your context
- Use vague language
- Skip the "why"
- Treat all requirements as equal priority

### ✅ Do

- Take time to think through answers
- Provide specific context upfront
- Use concrete numbers and examples
- Explain the problem first
- Separate must-haves from nice-to-haves

---

## Emergency Troubleshooting

**AI doesn't understand my domain**
→ Provide a glossary: "[Term]: [definition in plain English]"

**Spec is too generic**
→ Give real examples: "Example user: Sarah manages 12 projects..."

**AI suggests out-of-scope items**
→ State boundaries: "Out of scope: SMS, Slack, offline mode"

**Requirements conflict**
→ Ask AI to surface: "Do requirements A and B conflict? How to resolve?"

---

## Next Steps

1. **Pick your scenario:**
   - Vague idea? Start with product-wizard
   - Messy list? product-wizard to organize
   - Existing process? product-wizard to document
   - Need verification? atdd

2. **Copy a prompt** from above

3. **Answer AI questions** with specifics

4. **Review output** and refine

5. **Move to next skill** in the pipeline

---

## Full Examples

See [SKILLS-GUIDE.md](SKILLS-GUIDE.md) for complete conversation examples and detailed patterns.

See [examples/README.md](README.md) for realistic business scenarios with full outputs.

---

**Remember:** Good specs come from good questions. Embrace the AI's discovery questions - they reveal what you haven't thought about yet.
