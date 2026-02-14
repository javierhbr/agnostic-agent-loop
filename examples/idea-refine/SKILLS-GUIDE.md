# Skills Guide: Refining Ideas Without Code

*A business-focused guide to transforming vague concepts into detailed specifications using AI skills*

---

## What This Guide Is For

This guide shows **product managers, business analysts, and non-technical stakeholders** how to use four powerful AI skills to refine ideas into actionable specifications - no coding or technical knowledge required.

**You'll learn:** How to have productive conversations with AI to turn fuzzy concepts into crystal-clear requirements.

---

## The Four Skills You'll Use

### 1. 🧙 product-wizard

**What it does:** Transforms vague ideas into structured Product Requirements Documents (PRDs)

**When to use:** You have an idea but aren't sure how to articulate it clearly

**Output:** A complete PRD with problem statement, user stories, success metrics, and acceptance criteria

---

### 2. 📋 dev-plans

**What it does:** Breaks PRDs into phased implementation tasks with review checkpoints

**When to use:** After creating a PRD, to understand what needs to be built and in what order

**Output:** A structured development plan organized into phases with clear deliverables

---

### 3. ✅ atdd

**What it does:** Converts acceptance criteria into testable specifications

**When to use:** To ensure your requirements are verifiable and not just wishful thinking

**Output:** Clear test scenarios showing how to verify the feature works

---

### 4. 🔍 code-simplification

**What it does:** Reviews specifications for clarity and simplicity

**When to use:** To ensure your requirements are easy to understand and implement

**Output:** Simplified, clearer version of your specifications

---

## The Refinement Pipeline

```
Vague Idea
    ↓
[product-wizard skill]
    ↓
Structured PRD
    ↓
[dev-plans skill]
    ↓
Phased Implementation Plan
    ↓
[atdd skill]
    ↓
Testable Requirements
    ↓
[code-simplification skill]
    ↓
Clear, Simple Specifications
```

---

## How to Use Each Skill

### Using product-wizard

#### The Pattern

1. **You describe your idea** (even if vague)
2. **AI asks clarifying questions** (problem, users, success criteria, constraints)
3. **You answer with specifics**
4. **AI generates a structured PRD**

#### Example Prompt

```
I want to use the product-wizard skill to create a PRD for [your idea].

Here's what I'm thinking:
- The problem: [what problem are you solving?]
- For whom: [who needs this?]
- Current state: [how do they handle this today?]
- Success looks like: [what outcome do you want?]

Can you help me turn this into a comprehensive PRD?
```

#### What Happens Next

The AI will ask 4-7 discovery questions about:
- Why this problem matters now
- Who your target users are
- How you'll measure success
- What constraints you have (time, budget, technology)
- What's explicitly out of scope

**Your job:** Answer with specifics, not generalities.

---

### Using dev-plans

#### The Pattern

1. **You provide a PRD** (created by product-wizard)
2. **AI analyzes requirements**
3. **AI creates phased breakdown** with tasks and checkpoints
4. **You review and adjust priorities**

#### Example Prompt

```
I have a PRD for [feature name]. Can you use the dev-plans skill to break
this into implementation phases with review checkpoints?

Here's the PRD: [paste PRD or reference file location]

Please organize into [2-4] phases with clear deliverables.
```

#### What You'll Get

- Phase-by-phase breakdown
- Dependencies between tasks
- Review checkpoints (places to pause and verify)
- QA checklist for quality gates
- Realistic timeline estimates

---

### Using atdd

#### The Pattern

1. **You provide acceptance criteria** (from your PRD or user story)
2. **AI creates test scenarios** showing how to verify each criterion
3. **You review to ensure criteria are testable**
4. **AI refines based on your feedback**

#### Example Prompt

```
I have this user story with acceptance criteria. Can you use the atdd skill
to show me how we'd test this?

User Story:
As a [user type]
I want to [action]
So that [benefit]

Acceptance Criteria:
- [Criterion 1]
- [Criterion 2]
- [Criterion 3]

Please show test scenarios for each criterion.
```

#### What You'll Get

- Clear test scenarios (Given-When-Then format)
- Examples showing how to verify each requirement
- Edge cases you might have missed
- Non-technical explanations of what's being tested

---

### Using code-simplification

#### The Pattern

1. **You provide specifications or requirements**
2. **AI identifies complexity hotspots**
3. **AI suggests simplifications**
4. **You choose which simplifications to apply**

#### Example Prompt

```
Can you use the code-simplification skill to review these requirements for
clarity and simplicity?

Requirements:
[paste your requirements]

Are there areas that are overly complex or could be simplified?
```

#### What You'll Get

- Identification of vague or complex language
- Suggestions for clearer phrasing
- Simplified versions that maintain intent
- Explanation of why simpler is better

---

## Complete Workflow Examples

### Example 1: From Vague Idea to Testable Requirements

**Step 1: Start with product-wizard**

```
**You:**
I want to use the product-wizard skill to create a PRD for a customer
notification system.

We currently email customers about order updates, but many miss these
notifications. I'm thinking we need something better, maybe push
notifications or SMS?

Can you help me turn this into a comprehensive PRD?
```

**AI asks:**
- What specific problems do missed notifications cause?
- Who are your customers? (B2B, B2C, demographics?)
- How many customers? How many notifications per day?
- What's your budget and timeline?
- Which notification types are most critical?

**You answer:**

```
**You:**
- Problem: 25% of customers miss delivery notifications, leading to failed
  deliveries and support tickets
- Customers: B2C e-commerce, 50,000 monthly orders
- Critical notifications: Order confirmation, shipped, out for delivery
- Timeline: Need MVP in 10 weeks
- Budget: $50K, team of 2 developers
- Must integrate with existing order system
```

**AI generates:**
- Complete PRD with problem statement
- User personas
- 3 user stories with acceptance criteria
- Success metrics (reduce missed deliveries by 40%)
- Technical approach (SMS for critical, email for non-critical)
- Phased rollout plan

---

**Step 2: Use dev-plans to break it down**

```
**You:**
Now use the dev-plans skill to create a development plan from this PRD.
Break it into phases with review checkpoints.
```

**AI generates:**
- Phase 1 (Weeks 1-4): SMS integration, order confirmation notifications
- Phase 2 (Weeks 5-7): Shipping notifications, customer preferences
- Phase 3 (Weeks 8-10): Testing, rollout, monitoring
- Each phase has checkpoint: "STOP: Review before proceeding"

---

**Step 3: Use atdd for verification**

```
**You:**
For the user story "As a customer, I want to receive SMS when my order
ships so I know when to expect delivery", use the atdd skill to show
how we'd test this.
```

**AI generates:**

```
Test Scenario 1: Happy Path
Given: Customer has opted in to SMS notifications
When: Order status changes to "Shipped"
Then: Customer receives SMS within 5 minutes
And: SMS includes tracking number and estimated delivery date

Test Scenario 2: Customer Opted Out
Given: Customer has not opted in to SMS
When: Order status changes to "Shipped"
Then: Customer receives email (not SMS)
And: No SMS is sent

Test Scenario 3: Invalid Phone Number
Given: Customer phone number is invalid
When: Order ships
Then: System logs error
And: Falls back to email notification
And: Alerts support team to fix phone number
```

---

**Step 4: Simplify with code-simplification**

```
**You:**
Use the code-simplification skill to review this requirement:

"The notification system shall transmit SMS messages via a telecommunications
service provider API integration, ensuring redundancy through fallback
mechanisms and maintaining compliance with regulatory frameworks."

Is this too complex?
```

**AI responds:**

```
Yes, this can be simplified:

**Before:**
"The notification system shall transmit SMS messages via a
telecommunications service provider API integration, ensuring redundancy
through fallback mechanisms and maintaining compliance with regulatory
frameworks."

**After:**
"Send SMS via Twilio API. If Twilio fails, retry 3 times then send email.
Follow TCPA compliance rules (no SMS without opt-in)."

**Why simpler is better:**
- Specific provider named (Twilio) instead of generic "telecommunications"
- Concrete retry logic (3 times) instead of vague "redundancy"
- Clear compliance reference (TCPA) instead of generic "regulatory frameworks"
- Actionable fallback (email) instead of abstract "mechanisms"
```

---

### Example 2: From Requirements List to Structured Spec

**You have:** 20 bullet points in an email

**Goal:** Organize into coherent specification

**Step 1: product-wizard to organize**

```
**You:**
I have a messy requirements list for an expense approval system. Can you
use the product-wizard skill to organize this into a structured PRD?

Requirements:
- Automatic approval for small expenses
- Manager approval for medium expenses
- CFO approval for large expenses
- Mobile app for submitting expenses
- Receipt photo upload
- QuickBooks integration
- Email notifications
- Fast approval turnaround
- Prevent duplicate submissions
- Audit trail of approvals
- Support multiple currencies
- Offline mode
- Expense categories
- Monthly spend reports
- Budget tracking
- Approval workflows
- Slack integration
- Easy to use
- Secure
- Scalable

Some of these conflict (email vs Slack). Some are vague ("fast", "easy").
Can you help organize and clarify?
```

**AI responds:**

```
**AI:**
I see several issues to resolve:

**Conflicts:**
1. "Email notifications" vs "Slack integration" - which is primary?
2. "Automatic approval" vs "Manager approval" - what's the threshold?

**Vague terms:**
1. "Small/medium/large expenses" - define amounts
2. "Fast turnaround" - how fast?
3. "Easy to use" - what does this mean specifically?

**Questions:**
1. What defines expense size thresholds?
2. Primary notification channel?
3. Success metric for approval speed?
4. What's MVP vs Phase 2?
5. What integration is most critical?
```

**You clarify:**

```
**You:**
Good questions!

1. Thresholds: <$100 auto, $100-$5000 manager, >$5000 CFO
2. Notifications: Slack primary (for approvers), email fallback
3. Success: Reduce approval time from 3 days to 4 hours
4. MVP: approval workflow, QuickBooks, receipt upload
5. Phase 2: Offline, multi-currency, Slack, budget tracking

Must-have integrations: QuickBooks (we use it for accounting)
```

**AI generates:**
- Organized PRD with clear sections
- Conflicts resolved with your decisions documented
- Vague terms replaced with specific metrics
- Priorities: MVP vs Phase 2 clearly separated
- 10 user stories with acceptance criteria

---

**Step 2: dev-plans for phases**

```
**You:**
Use dev-plans skill to break this into implementation phases.
```

**AI generates:**

```
Phase 1: Core Approval Workflow (Weeks 1-4)
- Build amount-based routing ($100/$5K thresholds)
- Manager/CFO approval interface
- Approval status tracking
CHECKPOINT: Manual approval flow working

Phase 2: Automation & Integration (Weeks 5-8)
- Automatic approval for <$100
- QuickBooks API integration
- Receipt upload (mobile + web)
CHECKPOINT: End-to-end expense to accounting

Phase 3: Notifications & Launch (Weeks 9-10)
- Slack notifications
- Email fallback
- Audit logging
- Load testing
CHECKPOINT: 80% expenses approved in <4 hours
```

---

### Example 3: From Existing Process to Automated Workflow

**You have:** Manual spreadsheet process you want to automate

**Goal:** Preserve business logic while automating

**Step 1: Document current process**

```
**You:**
I want to automate our customer support ticket routing process using the
product-wizard skill. Currently we do this manually with a spreadsheet.

Current Process:
1. Ticket comes in via email
2. Support manager reads it
3. Manager assigns based on:
   - Technical issue → Engineering team
   - Billing question → Finance team
   - Feature request → Product team
   - VIP customer → Senior support rep
   - Urgent → Escalation to team lead
4. Assignment takes 10-15 minutes per ticket
5. We handle 100 tickets/day = 25 hours of manual routing

We have edge cases:
- VIP customers skip the queue
- After-hours tickets wait until morning
- Escalations override normal routing

Can you help create a PRD that preserves our business logic but automates it?
```

**AI asks:**

```
**AI:**
Great context! Let me understand the edge cases:

1. How do you identify VIP customers? (tag, account type, manual list?)
2. What defines "urgent"? Who decides?
3. For after-hours, what timezone? What are business hours?
4. How do you handle tickets that fit multiple categories?
5. What's the current success rate? (% of tickets routed correctly)
6. What happens if assignment fails?
```

**You clarify:**

```
**You:**
1. VIP = account tagged in CRM, about 50 accounts
2. Urgent = customer marks ticket as "urgent" + manager validates
3. Business hours: 9am-6pm Pacific, Mon-Fri
4. Multiple categories: Take first match in priority order (VIP > Urgent >
   Topic)
5. Success rate: 95% correct routing (5% get reassigned)
6. Failed assignment: Goes to general queue, manager notified
```

**AI generates:**
- PRD with complete routing rules
- Decision tree showing logic flow
- Edge case handling documented
- Success metric: Match 95% manual routing accuracy
- Migration plan: Run parallel for 2 weeks

---

**Step 2: atdd to verify logic preserved**

```
**You:**
Use atdd skill to show how we'd test that VIP routing works correctly.
```

**AI generates:**

```
Test: VIP Customer Priority Routing

Given: A VIP customer submits a ticket
When: Ticket enters the system
Then:
- Ticket is tagged "VIP"
- Ticket bypasses normal queue
- Ticket is assigned to senior support rep within 1 minute
- Assigned rep is notified immediately
- Manager receives FYI notification

Test: VIP + After Hours

Given: VIP customer submits ticket at 8pm Pacific
When: Ticket enters the system
Then:
- Ticket is held until 9am Pacific next business day
- At 9am, VIP routing rules apply
- On-call engineer receives urgent alert (for critical VIPs)

Test: Non-VIP Customer

Given: Regular customer submits ticket
When: Ticket is technical issue
Then:
- Ticket goes to engineering team queue
- First available engineer picks it up
- No VIP tags applied
- Standard SLA applies (4 hours)
```

---

## Prompt Templates by Use Case

### Starting from Scratch

```
I have an idea for [feature/product]. I want to use the product-wizard
skill to create a PRD.

What I know:
- The problem: [describe problem]
- Who needs this: [user type]
- Current state: [how they handle it today]
- What success looks like: [outcome you want]

What I don't know yet:
- [questions you have]

Can you help me develop this into a complete PRD?
```

---

### Organizing Messy Requirements

```
I have a list of [number] requirements for [project]. Some conflict, some
are vague. Use the product-wizard skill to organize these into a structured
PRD.

Requirements:
[paste your list]

Please help me:
1. Identify conflicts and resolve them
2. Clarify vague items
3. Separate must-haves from nice-to-haves
4. Create a coherent specification
```

---

### Documenting What Exists

```
We have a manual [process/system] that we want to automate. Use the
product-wizard skill to create a PRD that preserves our business logic.

Current process:
[describe step-by-step]

Edge cases we handle:
[list special cases]

What must be preserved:
[critical rules that can't change]

What we want to improve:
[pain points to fix]
```

---

### Making Requirements Testable

```
I have these acceptance criteria. Use the atdd skill to show how we'd
verify each one.

Acceptance Criteria:
- [Criterion 1]
- [Criterion 2]
- [Criterion 3]

For each criterion, show:
1. How to test it
2. What "passing" looks like
3. Edge cases to consider
```

---

### Simplifying Complex Specs

```
Use the code-simplification skill to review this requirement:

"[paste complex or jargon-filled requirement]"

Questions:
1. Is this too complex?
2. Can it be stated more simply?
3. Are there vague terms that need definition?
4. Would a non-expert understand this?
```

---

## Tips for Effective Skill Usage

### ✅ Do This

**Be specific about what you don't know**
```
❌ "I'm not sure about the details"
✅ "I don't know our current user count, but I can find out. Should I get
   that before we continue?"
```

**Provide context upfront**
```
❌ "We need notifications"
✅ "We're a B2C e-commerce platform with 50K monthly orders. Customers miss
   delivery notifications via email. We want to reduce failed deliveries."
```

**Use real numbers**
```
❌ "It should be fast"
✅ "Current approval takes 3 days. We want to get to 4 hours for 80% of
   requests."
```

**Clarify constraints**
```
❌ "ASAP"
✅ "Launch by Q2 end (10 weeks). Team of 2 developers. $50K budget."
```

---

### ❌ Don't Do This

**Don't assume AI knows your context**
```
❌ "Fix the notification problem"
✅ "We currently email order confirmations but 25% of customers miss them"
```

**Don't use jargon without explaining**
```
❌ "Integrate with our ESP via API"
✅ "Integrate with our email service provider (SendGrid) using their REST
   API"
```

**Don't skip the "why"**
```
❌ "We need a mobile app"
✅ "Our customers are on-the-go and need to approve expenses from their
   phones. Desktop-only is a blocker."
```

**Don't rush through discovery**
```
❌ [Answers "yes" to every AI question without thinking]
✅ [Takes time to think through implications, asks clarifying questions
   back]
```

---

## Common Patterns

### Pattern 1: The Discovery Loop

```
You → Vague idea
AI → Clarifying questions
You → Specific answers
AI → More questions based on answers
You → More details
AI → Draft specification
You → Feedback
AI → Refined specification
```

**Key:** Each round gets more specific. Embrace the questions - they reveal gaps.

---

### Pattern 2: The Simplification Pass

```
You → Create initial PRD (with product-wizard)
AI → Generates comprehensive spec
You → Review and find complex sections
You → Use code-simplification on complex parts
AI → Suggests simpler versions
You → Choose clearest option
```

**Key:** First pass captures everything. Second pass makes it clear.

---

### Pattern 3: The Validation Cycle

```
You → PRD with acceptance criteria
You → Use atdd to generate tests
AI → Shows test scenarios
You → Realize criteria is vague
You → Refine criteria
You → Re-run atdd
AI → Tests now make sense
```

**Key:** If you can't test it, you can't verify it. atdd exposes vague requirements.

---

## Real Conversation Examples

### Example: product-wizard Discovery

```
**You:**
I want to use product-wizard to create a PRD for better customer
notifications.

**AI:**
Let's explore this. What problem are notifications solving?

**You:**
Customers miss important updates.

**AI:**
Which updates specifically? What happens when they miss them?

**You:**
Order shipped notifications. When missed, customers aren't home for delivery.

**AI:**
How often does this happen? What's the business impact?

**You:**
About 25% of deliveries fail first attempt. Costs us $15 per re-delivery.
That's $75K/month.

**AI:**
Excellent. So the problem is: 25% delivery failure rate due to missed
notifications, costing $75K/month. How are you notifying them now?

**You:**
Email only.

**AI:**
What have customers said they want instead?

**You:**
SMS. In surveys, 60% said they'd prefer text messages for delivery updates.

**AI:**
Perfect. A few more questions:
1. How many customers?
2. How many orders per month?
3. What's your budget for this?
4. Timeline?

**You:**
50,000 customers, 20,000 orders/month, $50K budget, need it in 10 weeks.

**AI:**
Great context. I'll create a PRD with:
- Problem: $75K/month in re-delivery costs from missed notifications
- Solution: SMS for critical delivery updates
- Users: 50K customers, prioritizing active orderers
- Success metric: Reduce delivery failures from 25% to 10%
- Timeline: 10 weeks, $50K budget
```

**Notice:** Each question got more specific. Final PRD writes itself from this conversation.

---

### Example: atdd Verification

```
**You:**
Use atdd skill to verify this acceptance criterion:

"User receives notification within 5 minutes of order shipment"

**AI:**
Here's how to test that:

Test: Notification Timing
Given: Customer has order #12345
And: Customer has opted in to SMS
When: Order #12345 ships at 2:00pm
Then: Customer receives SMS between 2:00pm and 2:05pm
And: SMS contains order number and tracking link

**But I notice this criterion raises questions:**

1. What if order ships at 2:04:30pm - do they have until 2:09:30pm or 2:05pm?
2. What if SMS fails - do we retry?
3. What if customer hasn't opted in?
4. What timezone for the 5-minute SLA?

Would you like me to help clarify these edge cases?

**You:**
Good catches! Let's clarify:

1. 5 minutes from ship time (so 2:09:30pm in your example)
2. Retry 3 times over 15 minutes, then email fallback
3. No SMS opt-in = email notification instead
4. Customer's local timezone (from shipping address)

**AI:**
Perfect. Updated tests:

Test 1: Happy Path (with retries)
Given: Customer opted in to SMS
When: Order ships at 2:04:30pm
Then: First SMS attempt within 5 minutes (by 2:09:30pm)
If: SMS fails, retry at 2:09:30pm, 2:12:30pm, 2:14:30pm
If: All SMS fail, send email by 2:20pm

Test 2: No SMS Opt-In
Given: Customer has NOT opted in to SMS
When: Order ships
Then: Email sent within 5 minutes
And: No SMS attempted

Test 3: Timezone Handling
Given: Customer in PST, order ships at 2pm EST
When: Order ships
Then: Notification shows 11am PST (customer's timezone)
And: "Shipped at" time is in customer's timezone
```

**Notice:** atdd exposed edge cases that weren't in the original criterion. Now the requirement is bulletproof.

---

## Measuring Success

### How to Know Skills Are Working

**After using product-wizard:**
- [ ] Problem statement is clear (not vague)
- [ ] Success metrics are measurable (not "better")
- [ ] User needs are specific (not "everyone")
- [ ] Constraints are documented (time, budget, tech)
- [ ] Scope boundaries are explicit (in/out lists)

**After using dev-plans:**
- [ ] You understand what gets built when
- [ ] Dependencies are clear
- [ ] Review checkpoints exist
- [ ] You could explain the plan to a stakeholder

**After using atdd:**
- [ ] Every acceptance criterion has a test
- [ ] Tests cover edge cases
- [ ] You know what "passing" looks like
- [ ] Vague criteria have been exposed and fixed

**After using code-simplification:**
- [ ] Requirements use plain language
- [ ] Jargon is explained or removed
- [ ] A non-expert could understand it
- [ ] Complex sections are broken down

---

## What to Do Next

### 1. Pick Your Starting Point

- **Vague idea?** → Start with product-wizard
- **Messy requirements list?** → product-wizard to organize
- **Existing process?** → product-wizard to document + automate
- **Need to verify requirements?** → atdd
- **Specs too complex?** → code-simplification

### 2. Use the Prompt Templates

Copy from the "Prompt Templates by Use Case" section above. Adapt `[placeholders]` to your situation.

### 3. Embrace the Questions

When AI asks clarifying questions, don't rush. Each question reveals a gap in your thinking.

### 4. Iterate

First draft won't be perfect. Use the skills in sequence:
1. product-wizard → draft PRD
2. code-simplification → make it clearer
3. atdd → verify it's testable
4. dev-plans → see if it's buildable

---

## Advanced Techniques

### Technique 1: Comparative Refinement

```
Use product-wizard to create TWO different approaches to [problem], then
help me compare trade-offs:

Approach A: [describe option A]
Approach B: [describe option B]

For each, show:
- Pros/cons
- Cost/complexity
- Timeline
- Risk factors
```

---

### Technique 2: Incremental Clarification

```
I have this vague requirement: "[paste vague req]"

Use product-wizard to ask me questions that will make this specific and
testable. Keep asking until it's clear.
```

---

### Technique 3: Example-Driven Specification

```
Use product-wizard to help me specify [feature] through examples:

Example 1: When [scenario], the system should [behavior]
Example 2: When [different scenario], the system should [different behavior]

From these examples, create acceptance criteria.
```

---

## Troubleshooting

### "AI doesn't understand my domain"

**Fix:** Provide a glossary upfront

```
I'm working in [industry]. Here are key terms:

- [Term 1]: [definition in plain English]
- [Term 2]: [definition]
- [Term 3]: [definition]

Now use product-wizard to create a PRD for [feature].
```

---

### "Specifications are too generic"

**Fix:** Provide specific examples

```
The PRD feels generic. Let me provide real examples from our business:

Example customer: Sarah, a Project Manager managing 12 projects
Example pain point: She missed a deadline because email notification was
buried
Example desired outcome: SMS alert 24 hours before deadline

Use these examples to make the PRD more specific.
```

---

### "AI suggests things out of scope"

**Fix:** Explicitly state boundaries

```
STOP: Before we continue, here's what's explicitly OUT OF SCOPE:

- [Thing 1]
- [Thing 2]
- [Thing 3]

Please only suggest solutions that avoid these.
```

---

### "Requirements conflict"

**Fix:** Ask AI to surface conflicts

```
I suspect these requirements conflict:

Requirement A: [description]
Requirement B: [description]

Use product-wizard to:
1. Confirm if they conflict
2. Explain the conflict
3. Suggest ways to resolve it
```

---

## Summary

### The Essential Pattern

1. **Start with product-wizard** - Transforms vague → structured
2. **Use dev-plans** - Breaks structured → actionable
3. **Apply atdd** - Makes actionable → testable
4. **Apply code-simplification** - Makes testable → clear

### The Core Skill

**Ask better questions.** The skills work by:
- Asking you what you don't know
- Exposing gaps in your thinking
- Forcing specificity where you're vague
- Testing whether requirements are realistic

### Your Superpower

You can now turn a 5-minute conversation into a 50-page specification that would have taken weeks of meetings, emails, and revisions.

**Use it wisely.**

---

*For complete examples with realistic business scenarios, see the main [examples/README.md](README.md) guide.*
