# Scenario 1: Key Learnings

## Main Takeaway

**Vague ideas become actionable through systematic questioning.** The product-wizard skill transforms "we need X" into a complete specification by extracting the real problem, users, and success criteria.

---

## Lesson 1: Start with the "Why"

### The Pattern

❌ **Weak:** "We need push notifications"
✅ **Strong:** "Users miss deadlines because email notifications get buried. 30% of support tickets mention this."

### Why It Matters

- Justifies investment to stakeholders
- Prevents building the wrong solution
- Provides context for trade-off decisions
- Makes success measurable

### How to Apply

Always answer: **"What problem does this solve, and how do we know it's real?"**

---

## Lesson 2: Prioritize Users, Not Features

### The Pattern

❌ **Weak:** "All users need notifications"
✅ **Strong:** "Project Managers are priority—they manage 5-20 projects and feel the pain most"

### Why It Matters

- Prevents scope creep
- Focuses development effort
- Makes design decisions easier
- Enables phased rollout

### How to Apply

Identify your **primary user segment** and design for them first. Others can come in Phase 2.

---

## Lesson 3: Define Success Upfront

### The Pattern

❌ **Weak:** "Users will love it"  
✅ **Strong:** "50% of Project Managers enable notifications within 2 weeks; 30% daily interaction rate"

### Why It Matters

- Makes the feature testable
- Prevents endless refinement
- Aligns team on goals
- Enables data-driven decisions

### How to Apply

Ask: **"How will we know this worked?"** Pick 2-3 measurable metrics with targets.

---

## Lesson 4: Constraints Are Your Friend

### The Pattern

❌ **Weak:** "Build the best notification system possible"
✅ **Strong:** "8 weeks, 3 engineers, no native apps, must use browser push"

### Why It Matters

- Forces creative solutions (PWA instead of native)
- Prevents over-engineering
- Makes plans realistic
- Protects team capacity

### How to Apply

Always specify: **Timeline, budget, technical limitations, team size**

---

## Lesson 5: Scope Boundaries Prevent Bloat

### The Pattern

❌ **Weak:** List only what you're building
✅ **Strong:** Explicitly list what you're NOT building (SMS, Slack, scheduling, etc.)

### Why It Matters

- Prevents feature creep during development
- Manages stakeholder expectations
- Documents future possibilities
- Speeds up decisions

### How to Apply

Create a **"Non-Goals"** section. If someone asks "what about X?", you have an answer.

---

## Lesson 6: Questions Reveal Gaps

### What Happened

Sarah's initial request was vague. The AI's discovery questions surfaced:
- Missing success metrics
- Unclear user prioritization
- Unspecified technical constraints
- Ambiguous scope

### Why It Matters

**You don't know what you don't know.** Structured questions reveal blind spots before they become problems.

### How to Apply

Embrace the discovery phase. If AI asks a question you can't answer, that's valuable information—go find out.

---

## Lesson 7: Iteration Is Fast

### Timeline Breakdown

- Initial prompt: 2 minutes
- Discovery (2 rounds): 15 minutes
- PRD generation: 5 minutes (automated)
- Review: 10 minutes
- Dev plan: 5 minutes
- **Total: 37 minutes**

### Why It Matters

**37 minutes** to go from "we need notifications" to a complete PRD with user stories, acceptance criteria, and a development plan. This is 10x faster than traditional requirements gathering.

### How to Apply

Don't try to perfect your initial request. Start with what you know, answer questions, and iterate.

---

## Lesson 8: Specificity Compounds

### The Progression

1. **Vague:** "We need notifications"
2. **Better:** "We need mobile push notifications"
3. **Good:** "Project Managers need push notifications for task assignments"
4. **Great:** "Project Managers need browser push notifications for 4 event types, measured by 50% adoption in 2 weeks"

Each layer of specificity makes the next question clearer.

### How to Apply

**Add one specific detail at a time.** You don't need all the answers upfront—the conversation builds them.

---

## Common Pitfalls (And How to Avoid Them)

### Pitfall 1: "I'll figure out success metrics later"

❌ **Leads to:** Features that ship but don't move the needle

✅ **Instead:** Define success criteria during discovery, not after launch

### Pitfall 2: "Let's build it for everyone"

❌ **Leads to:** Diluted features that satisfy no one well

✅ **Instead:** Pick your primary user, nail it for them, expand later

### Pitfall 3: "We don't have technical constraints"

❌ **Leads to:** Unrealistic plans and disappointed teams

✅ **Instead:** Be honest about time, budget, and tech limitations

### Pitfall 4: "The AI should just know our context"

❌ **Leads to:** Generic PRDs that don't fit your needs

✅ **Instead:** Provide specific context: company size, user count, current state

---

## What to Do Next

### Immediate Actions

1. **Try it yourself:** Pick a vague idea and use the prompts from this scenario
2. **Review the PRD output:** See how discovery answers become structured requirements
3. **Compare to your current process:** How long does your team take to create similar specs?

### Level Up

1. **Use the Business Brief Template:** Fill it out before starting to organize your thoughts
2. **Customize the prompts:** Adapt them to your domain and company
3. **Build a prompt library:** Save prompts that work well for your team

---

## Success Pattern Summary

```
Vague Idea
    ↓
Ask: Why? (Problem statement)
    ↓
Ask: Who? (User prioritization)
    ↓
Ask: How will we know? (Success metrics)
    ↓
Ask: What constraints? (Time, budget, tech)
    ↓
Ask: What's out of scope? (Boundaries)
    ↓
Actionable PRD with testable requirements
```

**Remember:** The best requirements emerge from questions, not from trying to write a perfect brief upfront.

---

[← Back to Conversation](conversation.md) | [View Next Scenario →](../02-requirements-list/scenario.md)
