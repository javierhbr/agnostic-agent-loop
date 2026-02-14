# Requirements Refinement Examples

## Welcome

This guide shows product managers, business analysts, and stakeholders how to transform ideas into detailed, AI-ready specifications using a systematic requirements pipeline.

**No technical background needed** - these examples use plain language and demonstrate realistic conversations with AI.

## Quick Links

- 📘 **New to the skills?** Start with [Skills Guide](SKILLS-GUIDE.md) - Learn how to use the four skills without technical jargon
- 📇 **Need a quick reminder?** Use the [Quick Reference Card](QUICK-REFERENCE.md) - Essential prompts and patterns
- 📋 **Ready to dive in?** Browse scenarios below for complete examples with realistic outputs

## The Requirements Pipeline

```
Vague Idea → Discovery → PRD → Development Plan → Tasks → Acceptance Tests
     ↓           ↓        ↓           ↓            ↓            ↓
  "Better    Questions  Structured  Phased    Specific   Executable
   notify"   & Answers  Document    Breakdown  Work Items  Validation
```

### Four Key Skills

1. **product-wizard** - Transforms ideas into Product Requirements Documents (PRDs)
2. **dev-plans** - Breaks PRDs into phased development tasks
3. **atdd** - Converts acceptance criteria into executable tests
4. **code-simplification** - Reviews code for clarity (supporting skill)

## Choose Your Scenario

Each scenario demonstrates the complete workflow from different starting points:

### [Scenario 1: From a Vague Idea](01-vague-idea/scenario.md)

**Start here if you have:** A fuzzy concept that needs definition

**Example:** "We need better notifications" → Complete PRD with metrics, user stories, and success criteria

**You'll learn:**
- How to articulate vague ideas
- What questions AI will ask
- How to define success metrics
- Creating testable requirements

[View Scenario](01-vague-idea/scenario.md) | [Copy-Paste Prompts](01-vague-idea/prompts.md)

---

### [Scenario 2: From a Requirements List](02-requirements-list/scenario.md)

**Start here if you have:** Bullet points, notes, or scattered requirements

**Example:** 15 bullet points about expense approval → Structured PRD with priorities and resolved conflicts

**You'll learn:**
- Organizing scattered requirements
- Resolving conflicts
- Prioritizing features
- Transforming vague items into testable criteria

[View Scenario](02-requirements-list/scenario.md) | [Copy-Paste Prompts](02-requirements-list/prompts.md)

---

### [Scenario 3: From Existing Logic](03-existing-logic/scenario.md)

**Start here if you have:** A manual process or workflow you want to automate

**Example:** Spreadsheet-based support ticket routing → Automated workflow with preserved business rules

**You'll learn:**
- Documenting existing processes
- Preserving business logic
- Identifying automation opportunities
- Managing migration risks

[View Scenario](03-existing-logic/scenario.md) | [Copy-Paste Prompts](03-existing-logic/prompts.md)

---

### [Scenario 4: From Codebase Documentation](04-codebase-documentation/scenario.md)

**Start here if you have:** Existing code that needs changes or refactoring

**Example:** Legacy payment system → Refactoring requirements with backward compatibility

**You'll learn:**
- Documenting what must be preserved
- Separating refactoring from features
- Risk management for legacy systems
- Testing strategies for existing behavior

[View Scenario](04-codebase-documentation/scenario.md) | [Copy-Paste Prompts](04-codebase-documentation/prompts.md)

---

## Quick Start Tools

### [Business Brief Template](templates/business-brief-template.md)

Fill out this template before starting - it gives AI the context needed for better results.

### [Requirements Checklist](templates/requirements-checklist.md)

Self-assessment tool to ensure you've covered the basics before creating a PRD.

### [Prompt Library](templates/prompt-library.md)

Master collection of 15-20 copy-paste prompts organized by workflow stage.

---

## How to Use These Examples

### 1. Choose Your Scenario

Pick the scenario that matches your starting point (vague idea, requirements list, existing logic, or codebase).

### 2. Read the Scenario

Understand the business context and stakeholders - this shows you what information to prepare.

### 3. Review the Prompts

Open the `prompts.md` file to see what you'll actually type to the AI. These are copy-paste ready.

### 4. Follow the Conversation

See how a real dialogue unfolds - questions, answers, and iterations.

### 5. Study the Outputs

Look at the PRD, development plan, and acceptance tests that result from the conversation.

### 6. Learn the Lessons

Read `learnings.md` for key takeaways you can apply to your own work.

---

## Integration with Agnostic Agent

After creating a PRD using the product-wizard skill, you can feed it into the agnostic agent workflow:

```bash
# Create an openspec proposal and tasks from your PRD
agentic-agent openspec init "<feature-name>" --from .agentic/spec/prd-<feature>.md

# This generates:
# - Proposal document
# - Development plan
# - Breakdown into individual tasks
# - Acceptance criteria for each task
```

The workflow becomes:
1. **Brainstorm** (optional) - Refine your idea
2. **product-wizard** → Create PRD
3. **agentic-agent openspec init** → Generate tasks
4. **dev-plans** → Break down into phases
5. **atdd** → Create acceptance tests for each task

---

## Tips for Success

### Be Specific About Context

**Bad:** "We need a dashboard"
**Good:** "We need a customer support dashboard showing ticket volume, average response time, and agent performance for our 12-person CS team"

### Define Success Upfront

**Bad:** "It should be fast"
**Good:** "Search results should load in under 200ms for our 10,000 customer database"

### Clarify Constraints

Always mention:
- Budget limitations
- Timeline requirements
- Technology constraints (e.g., "must integrate with Salesforce")
- Compliance requirements (e.g., GDPR, HIPAA)

### Ask Questions

If the AI's questions seem off-track, explain why. The dialogue is collaborative.

---

## What Makes a Good Requirement?

| Quality | Bad Example | Good Example |
|---------|-------------|--------------|
| **Specific** | "Easy to use" | "Users can complete expense submission in 3 clicks or fewer" |
| **Measurable** | "Fast performance" | "API responses under 200ms for 95th percentile" |
| **Testable** | "Works well" | "All user actions have confirmation feedback within 100ms" |
| **Scoped** | "Everything should sync" | "Customer data syncs to Salesforce every 15 minutes" |

---

## Common Questions

**Q: Do I need to know how to code?**
A: No. These examples show business-focused conversations. The AI handles technical translation.

**Q: How long does this process take?**
A: For most features: 30-60 minutes to create a PRD, 15-30 minutes for a development plan.

**Q: What if I don't know all the answers?**
A: That's expected. The AI will help surface what's missing through questions.

**Q: Can I skip the PRD and go straight to tasks?**
A: You can, but you'll likely waste time on rework. The PRD is your "source of truth" that prevents misunderstandings.

**Q: What if requirements change?**
A: Update the PRD, then regenerate the development plan. The PRD is a living document.

---

## Next Steps

1. **Choose your scenario** from the list above
2. **Open the prompts.md file** and copy the initial prompt
3. **Start a conversation** with your AI tool
4. **Reference the example conversation** if you get stuck
5. **Use the templates** for your next project

---

## Additional Resources

- [CLAUDE.md](../CLAUDE.md) - Project-specific workflow guidance
- [Product Wizard Skill](../internal/skills/packs/product-wizard/SKILL.md) - Full PRD skill documentation
- [Dev Plans Skill](../internal/skills/packs/dev-plans/SKILL.md) - Development planning skill
- [ATDD Skill](../internal/skills/packs/atdd/SKILL.md) - Acceptance test-driven development

---

**Ready to start?** Pick a scenario above and dive in. Remember: the best PRD is one that helps your team build the right thing, not a perfect document that never gets used.
