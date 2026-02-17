# Requirements Checklist

*Use this before creating a PRD to ensure you've covered the basics.*

---

## Problem Definition

- [ ] **Problem is clearly stated** - Not "we need X" but "users can't Y because Z"
- [ ] **Evidence provided** - Support tickets, user feedback, metrics, or competitive pressure
- [ ] **Impact quantified** - Costs, time waste, churn risk, or revenue opportunity
- [ ] **"Why now" is explained** - What changed or what's the urgency

**Example:**
✅ "30% of support tickets report missed deadlines because email notifications get buried. Competitors offer push notifications. We risk churn."

---

## User Understanding

- [ ] **Target users identified** - Specific roles, not "everyone"
- [ ] **Primary user chosen** - If multiple user types, which is priority?
- [ ] **User needs documented** - What do they need to accomplish?
- [ ] **User context provided** - How many users? What's their workflow?

**Example:**
✅ "Project Managers (400 users) manage 5-20 projects each. They need real-time alerts for blockers to respond within 1 hour."

---

## Success Measurement

- [ ] **Metrics defined** - 2-3 measurable indicators of success
- [ ] **Targets set** - Specific numbers, not "more" or "better"
- [ ] **Baseline known** - Current state to compare against
- [ ] **Measurement method clear** - How/where you'll track these metrics

**Example:**
✅ "50% of PMs enable notifications within 2 weeks (tracked via analytics). 30% daily interaction rate (clicks/notifications sent)."

---

## Constraints Documented

- [ ] **Timeline specified** - When do you need to launch?
- [ ] **Budget/resources known** - Team size, budget limits
- [ ] **Technical constraints listed** - Platform, integrations, limitations
- [ ] **Compliance requirements** - GDPR, HIPAA, security, etc.

**Example:**
✅ "8-week timeline, 3 engineers, must use browser push (no native apps), must comply with browser permission standards."

---

## Scope Boundaries

- [ ] **Must-haves listed** - Core features for MVP
- [ ] **Nice-to-haves separated** - Phase 2 features
- [ ] **Explicit non-goals** - What you're NOT building
- [ ] **Priorities clear** - If you had to cut scope, what goes first?

**Example:**
✅ "Must-have: 4 notification types. Nice-to-have: notification center. Non-goals: SMS, Slack integration, scheduling."

---

## Context Provided

- [ ] **Current state described** - How do users do this today?
- [ ] **Existing systems listed** - What tools/platforms are involved?
- [ ] **Company context** - Size, industry, user count
- [ ] **Prior attempts** - Have you tried to solve this before?

**Example:**
✅ "Currently email-only notifications. 2,000 users, B2B SaaS. Previous attempt: built Slack integration but only 10% of users have Slack."

---

## Requirements Quality

Use these tests for each requirement:

- [ ] **Specific** - No vague words like "fast", "easy", "better"
- [ ] **Measurable** - Can be tested objectively
- [ ] **Achievable** - Realistic given constraints
- [ ] **Relevant** - Supports the core problem/goal
- [ ] **Testable** - Clear pass/fail criteria

**Examples:**

| Vague ❌ | Specific ✅ |
|---------|-----------|
| "Fast performance" | "API responses <200ms for 95th percentile" |
| "Easy to use" | "Users complete expense submission in ≤3 clicks" |
| "Works on mobile" | "Supports iOS 15+ and Android 11+ browsers" |
| "Good notifications" | "Notifications deliver within 30 seconds of trigger" |

---

## Red Flags

Watch out for these warning signs:

### 🚩 Vague Success Criteria

❌ "Users will love it"
❌ "Increase engagement"
❌ "Make it better"

✅ **Fix:** Add specific numbers and measurement method

### 🚩 No User Prioritization

❌ "For all users" or "Everyone needs this"

✅ **Fix:** Choose your primary user segment

### 🚩 Undefined Constraints

❌ "As soon as possible" or "Whatever it takes"

✅ **Fix:** Set realistic timeline, budget, technical limits

### 🚩 Scope Creep Language

❌ "And also...", "While we're at it...", "It would be nice if..."

✅ **Fix:** Move to "Phase 2" or "Non-goals" section

### 🚩 Solution Before Problem

❌ Starting with "We need a dashboard" instead of the problem

✅ **Fix:** Start with user pain point, then discuss solutions

---

## Self-Assessment Scoring

Count your checkmarks:

- **20-25 checks:** Excellent! You're ready to create a PRD
- **15-19 checks:** Good start. Fill in the gaps before proceeding
- **10-14 checks:** More discovery needed. Use the business brief template
- **<10 checks:** Too early. Do more user research and problem validation

---

## What to Do Based on Score

### If you scored 20-25

✅ You're ready! Use the [Prompt Library](prompt-library.md) to start your PRD.

### If you scored 15-19

📝 Fill out the [Business Brief Template](business-brief-template.md) to clarify gaps, then proceed.

### If you scored 10-14

🔍 You need more discovery:
- Talk to users about their pain points
- Gather metrics on current state
- Define success criteria with stakeholders
- Document technical constraints

### If you scored <10

⏸️ Pause and do problem validation:
- Is this a real problem worth solving?
- Do you have evidence (not just opinions)?
- Have you talked to target users?
- Do you know what success looks like?

---

## Next Steps

1. **Complete this checklist** for your project
2. **Fill any gaps** using the business brief template
3. **Choose your scenario** from the [main README](../README.md)
4. **Use the prompts** from the [Prompt Library](prompt-library.md)

---

**Remember:** A good PRD starts with good requirements. Spend time here to save time later.
