# SDD with Prompts: AI Agent Integration Guide

This guide shows how to use SDD methodology with AI agents (Claude, ChatGPT, etc.) via prompts. Use these prompts with `USAGE-LEVELS.md` and `SMALL-PROJECTS.md` to automate specification, planning, and verification.

---

## TL;DR: Prompt Patterns by Level

| Level | Agent Role | Prompt Pattern | Output |
|-------|-----------|---|--------|
| **OpenSpec** | PM/Developer | "Decompose feature into tasks" | proposal.md + tasks.md |
| **Platform** | Architect | "Design multi-service spec" | platform-spec.md + fanout.yaml |
| **SDD Full** | Analyst | "Interview team, extract evidence" | discovery.md |
| **SDD Full** | Architect | "Design from discovery" | feature-spec + component-specs |
| **SDD Full** | Developer | "Implement from spec" | impl-spec + tasks.yaml |
| **SDD Full** | Verifier | "Verify all ACs with evidence" | verify.md |

---

## Level 1: OpenSpec with Prompts

### Use Case: Small Feature (1-2 People, 1 Service)

#### Prompt 1: Generate Proposal & Tasks

**Invoke with:** Claude, ChatGPT, or your preferred AI

```
You are a Product Manager helping decompose a feature into actionable tasks.

FEATURE: [Paste your feature description here]

Generate:
1. A clear proposal.md (1-2 paragraphs):
   - Problem: What user problem does this solve?
   - Why: Why now? What's the metric for success?
   - What: What are we building?
   - Acceptance: 3-5 acceptance criteria (Given/When/Then if possible)

2. A tasks.md breakdown:
   - Group tasks by component (Backend, Frontend, Testing)
   - List 5-10 actionable tasks
   - For each task, write acceptance criteria
   - Include dependencies (Task 3 depends on Task 1, etc.)

Format as valid markdown so I can copy into proposal.md and tasks.md directly.

Example feature (for reference):
---
FEATURE: Add dark mode toggle to settings

OUTPUT:
# Proposal

## Problem
Users report eye strain when using the app at night.
Metric: Measure daily active users staying > 2 hours (goal: 10% increase).

## Solution
Add dark mode theme toggle in settings, persist preference.

## Acceptance
- [ ] Dark mode toggle appears in settings
- [ ] Theme persists across sessions
- [ ] All UI readable in dark mode
- [ ] Mobile responsive

# Tasks

## Backend (2 tasks)
1. Add theme preference to user model
   - Save to database
   - API endpoint: PATCH /users/preferences {theme: "light"|"dark"}
   - Acceptance: API returns 200, persists in DB

2. Add theme endpoint test
   - Test: GET /users/preferences returns current theme
   - Acceptance: Test passes

## Frontend (4 tasks)
3. Create Settings page component
   - Acceptance: Page renders without errors

4. Add dark mode toggle button
   - Acceptance: Toggle appears, is clickable

5. Connect to API (fetch user preference, save on toggle)
   - Acceptance: Theme preference saves to DB

6. Add dark mode CSS styles
   - Acceptance: All components readable in dark mode

## Testing (1 task)
7. Manual test on mobile devices
   - Acceptance: Dark mode works on iOS and Android
---
```

#### Prompt 2: Create OpenSpec CLI Command

After AI generates proposal + tasks, use the CLI:

```bash
# Create the spec directory and files
agentic-agent openspec init "My Feature" --from requirements.md

# Update the generated proposal
vim .agentic/specs/my-feature/proposal.md
# (Paste AI's generated proposal)

# Update the generated tasks
vim .agentic/specs/my-feature/tasks.md
# (Paste AI's generated tasks)

# Track progress
agentic-agent task list
agentic-agent task claim <task-id>
# ... code ...
agentic-agent task complete <task-id>

# When all tasks done
agentic-agent openspec complete my-feature
agentic-agent openspec archive my-feature
```

---

## Level 2: Platform with Prompts

### Use Case: Multi-Service Feature (2-3 Services, Architecture Decisions)

#### Prompt 1: Generate Platform Spec

```
You are a Platform Architect designing a multi-service feature.

INITIATIVE: [Paste initiative description]
SERVICES AFFECTED: [List: auth-service, payment-service, etc.]
TIMELINE: [e.g., 2 weeks]

Generate a concise platform-spec.md (max 1 page) with:

1. Problem: Business problem in 2-3 sentences
2. Solution: What we're building in 2-3 sentences
3. Services: List each affected service and what changes
4. Contracts: List new/changed APIs or events
5. Acceptance Criteria: 5+ testable outcomes
6. Blockers: Any decisions needed before implementation?

Format as markdown for platform-specs/<feature-id>/spec.md

Example output (reference):
---
# Platform Spec: Guest Checkout

## Problem
30% of cart abandonment is from users unwilling to create account.

## Solution
Enable checkout without account creation, collect email only.

## Services
- checkout-service: Add guest flow, skip account creation
- payment-service: Process guest charges
- email-service: Send confirmation to guest email

## Contracts
NEW EVENT: guest.checkout_completed
{
  "guest_email": "...",
  "order_id": "...",
  "total_cents": 9999
}

## Acceptance Criteria
- [ ] Guest can checkout without account
- [ ] Email confirmation sent within 60 seconds
- [ ] Charge processed successfully
- [ ] Cart abandonment rate < 25% (measure in 2 weeks)

## Blockers
None - can proceed immediately
---
```

#### Prompt 2: Generate Fan-Out Component Specs

```
For each service from the platform spec, generate a component-spec:

SERVICE: [auth-service, payment-service, etc.]
PLATFORM SPEC: [Paste platform spec above]

Generate component-spec-[service].md with:

1. API Changes: What endpoints are new/modified?
2. Data Model: What fields added to existing models?
3. Dependencies: What other services does this call?
4. Acceptance Criteria: 3-5 testable per-component outcomes
5. Edge Cases: 2-3 edge cases and how you handle them

Format as markdown for platform-specs/fanout/component-spec-[service].md

Keep it brief (< 500 words).
```

#### Prompt 3: Create CLI Commands

```bash
# Initialize platform repo
agentic-agent platform init

# Create platform spec with fan-out
agentic-agent platform add-feature \
  --name "Guest Checkout" \
  --initiative "CHECKOUT-2025" \
  --context-pack "cp-v2" \
  --fanout "checkout-service,payment-service,email-service"

# Update with AI-generated spec
vim platform-specs/guest-checkout/spec.md
# (Paste AI's generated platform spec)

# Each service team creates OpenSpec
cd checkout-service/
agentic-agent openspec init "Guest checkout flow" \
  --from ../platform-specs/guest-checkout/spec.md

# Team implements their component
vim .agentic/specs/guest-checkout/proposal.md
# (Paste AI's generated component-spec)

agentic-agent task claim <task-id>
# ... code ...
agentic-agent task complete <task-id>

# Mark components complete
agentic-agent openspec complete guest-checkout

# Platform team tracks
agentic-agent platform change-priority \
  --initiative "CHECKOUT-2025" \
  --priority "Done"
```

---

## Level 3: SDD Full with Prompts

### Use Case: Critical Feature (High Risk, Multiple Services)

#### Prompt 1: Analyst Phase - Discovery Interview Generator

```
You are an SDD Analyst helping a team discover the problem space.

INITIATIVE: [Feature name]
TEAM CONTEXT: [Who's involved, what's the team size]

Generate an interview guide with 8-10 questions to ask the team ONE AT A TIME:

For each question:
1. Write the exact question to ask
2. Explain what evidence you're looking for
3. Provide example good/bad answers

Focus on:
- Problem evidence (real data, not assumptions)
- Affected components
- Risk factors
- Rollback concerns
- Observability needs

Example output:
---
## Interview Question 1: Problem Evidence

**Question to Ask:**
"What evidence shows this is really a problem? What data do we have?"

**What You're Looking For:**
- Real metrics: support tickets, error rates, user feedback
- NOT: "I think users want this" or "It would be nice"

**Good Answer Examples:**
- "We have 43 support tickets in January asking about subscriptions"
- "8 customers explicitly said they'd switch if we had this"
- "Our NPS dropped 5 points because of this limitation"

**Bad Answer Examples:**
- "Everyone needs this" (vague, no evidence)
- "It would improve user experience" (assumption, not data)
---

Continue for 8-10 questions...
```

#### Prompt 2: Architect Phase - Spec Generator

```
You are an SDD Architect designing from team discovery.

DISCOVERY SUMMARY:
[Paste team evidence and findings from Analyst phase]

Generate TWO specs:

1. feature-spec.md (Platform Level - WHAT)
   - Problem (from discovery)
   - Goals / Non-Goals
   - User Experience (walkthrough)
   - Domain Responsibilities (who owns what)
   - Cross-Domain Interactions (event flows)
   - NFRs (performance, logging, metrics, tracing, PII)
   - Acceptance Criteria (7+ in Given/When/Then format)
   - Feature Flag & Rollback Strategy

2. component-spec-[service].md (Per Service - HOW)
   - For EACH affected service, generate one:
   - API Contracts (endpoints, events, schemas)
   - Data Model (tables, fields, constraints)
   - Acceptance Criteria (3-5 per service)
   - Edge Cases (3-4 scenarios)

Format as markdown for feature-spec.md and component-spec-*.md

Make specs detailed enough to pass these gates:
- Gate 1: Every section has Source: line (cite discovery or requirements)
- Gate 2: No cross-domain DB access (event-driven only)
- Gate 3: All contracts identified, consumers listed
- Gate 4: Logging, metrics, tracing, PII handling specified
- Gate 5: No TBD sections, acceptance criteria testable
```

#### Prompt 3: Developer Phase - Implementation Spec Generator

```
You are an SDD Developer reading a component-spec.

COMPONENT SPEC:
[Paste component-spec-[service].md here]

Generate:

1. impl-spec.md
   - Data Model: SQL schema for any new tables
   - Code Changes: Exact functions/modules to create/modify
   - Edge Cases: Table with 4+ edge case scenarios and handling
   - Observability: Logging, metrics, tracing specs
   - Rollout Plan: Feature flag strategy, canary rollout

2. tasks.yaml
   - Decompose into 8-10 actionable tasks
   - Group by component (Backend, Frontend, Testing)
   - For each task: what you'll modify, acceptance criteria
   - List dependencies (Task 3 depends on Task 1)

Format as markdown and YAML for impl-spec.md and tasks.yaml.

Each task should be completable in 1-3 days by one developer.
```

#### Prompt 4: Verifier Phase - Test Plan Generator

```
You are an SDD Verifier creating a verification plan.

FEATURE SPEC:
[Paste feature-spec.md with all ACs]

IMPLEMENTATION:
[Paste impl-spec.md]

Generate verify.md with:

1. For EACH acceptance criterion:
   - What test will you write?
   - What log/metric shows success?
   - Example test code or assertion

2. Code Quality Checklist:
   - Test coverage target (80%+)
   - Lint checks
   - Build checks
   - Performance targets

3. Observability Verification:
   - What logs should exist?
   - What metrics should we see?
   - What traces should appear?

4. Blocking Conditions:
   - If X test fails, we can't merge
   - If Y is untestable, we can't ship
   - If Z touches payment/PII, requires human approval

Format as markdown for verify.md.

Include example test output and log excerpts.
```

#### Prompt 5: Create CLI Commands for SDD Full

```bash
# Start SDD workflow
agentic-agent specifyify start "Subscription Billing" --risk critical

# Analyst phase
# Use Prompt 1 to generate interview guide
# Interview team and document answers

# Architect phase
# Use Prompt 2 to generate feature + component specs
mkdir -p .agentic/specs/subscription-billing
# (Paste AI-generated feature-spec.md)
# (Paste AI-generated component-spec-*.md)

# Check gates
agentic-agent specifyify gate-check subscription-billing --format text

# If ADR needed, create it
agentic-agent specifyify adr create --title "Where to store card tokens?"

# Resolve when decision made
agentic-agent specifyify adr resolve ADR-001

# Developer phase (parallel)
# Use Prompt 3 to generate impl-specs + tasks
cd billing-service/
# (Paste AI-generated impl-spec.md)
# (Paste AI-generated tasks.yaml)

agentic-agent task claim <task-id>
# ... code ...
agentic-agent task complete <task-id>

# Verifier phase
# Use Prompt 4 to generate test plan
# Run tests, collect evidence

# Complete
agentic-agent openspec complete subscription-billing
```

---

## Combining Prompts with USAGE-LEVELS.md

### Decision: Which Level Should I Use?

Use **USAGE-LEVELS.md** to decide, then use **Prompts** to implement:

```
1. Read USAGE-LEVELS.md → Decision Tree
   (Answer: OpenSpec / Platform / SDD Full)

2. Use appropriate Prompt from above:
   - OpenSpec → Use Prompt 1-2
   - Platform → Use Prompt 1-3
   - SDD Full → Use Prompt 1-5

3. Use CLI commands to track progress
```

**Example:** Feature affects 3 services, medium risk

```
DECISION (from USAGE-LEVELS.md):
"Multi-service + architecture decisions + medium risk = Use Platform"

IMPLEMENTATION:
1. Use Prompt 1 to generate platform-spec.md
2. Use Prompt 2 to generate component-specs
3. Use CLI: agentic-agent platform add-feature
4. Each team uses Prompt 1-2 for OpenSpec
5. Each team uses Prompt 3 for impl-specs
6. Verify with Prompt 4
7. CLI: agentic-agent openspec complete
```

---

## Combining Prompts with SMALL-PROJECTS.md

### For Tiny Projects (Skip Prompts)

```bash
# Don't use prompts - just use CLI
agentic-agent task create "Fix bug"
agentic-agent task claim <task-id>
# ... code ...
agentic-agent task complete <task-id>
```

### For Small Projects (Use Minimal Prompts)

```
1. Read SMALL-PROJECTS.md → Template section
2. Use Prompt 1 (OpenSpec) to generate proposal + tasks
3. Use CLI: agentic-agent openspec init
4. Implement manually or with prompts

Example:
"Generate proposal + tasks for user dashboard feature"
→ Use Prompt 1 from this guide
→ Paste into proposal.md + tasks.md
→ Implement
```

### For Growing Complexity (Upgrade Prompts)

```
Initial: Small project with OpenSpec prompts
Growing: Multiple services needed

Action:
1. Use SMALL-PROJECTS.md → "When to upgrade" section
2. Switch to Platform level
3. Use Prompt 1-3 (Platform prompts)
4. Re-scope existing work into component-specs
```

---

## Best Practices: Using Prompts with CLI

### ✅ Do This

```bash
# 1. Prompt generates spec
# (Paste AI output into file)

# 2. CLI tracks it
agentic-agent openspec init "Feature"

# 3. Team implements
agentic-agent task claim <task-id>

# 4. CLI marks complete
agentic-agent openspec complete <feature-id>

# AI + CLI = best of both worlds
```

### ❌ Don't Do This

```bash
# ❌ Use prompts but never save specs
# ❌ Use CLI but never generate proposals
# ❌ Use prompts for everything (overkill for simple bugs)
# ❌ Use CLI for everything (no planning)

# Balance: Prompts for design, CLI for tracking
```

---

## Prompt Workflow by Role

### Product Manager / Architect

```
1. Use USAGE-LEVELS.md → Decide level
2. Use Prompt 1-3 → Generate specs
3. Use CLI → Create/update
4. Team implements
```

### Developer

```
1. Read component-spec from PM
2. Use Prompt 3 → Generate impl-spec + tasks
3. Use CLI → Track progress
4. Code implementation
```

### Verifier / QA

```
1. Read feature-spec + impl-spec
2. Use Prompt 4 → Generate test plan
3. Run tests, gather evidence
4. Use CLI → Mark complete
```

---

## Real Example: Feature with Prompts

### Scenario: Add Card Refunds (2 Services)

**Step 1: Decide Level (USAGE-LEVELS.md)**
- 2 services (payments, accounting)
- Medium risk (financial but no payment processing)
- Result: **Use Platform level**

**Step 2: Generate Platform Spec (Prompt 1-2)**

```
Prompt:
"Generate a 1-page platform spec for adding card refunds to our payment system.
Services: payments-service, accounting-service.
Timeline: 2 weeks."

AI Output:
# Platform Spec: Card Refunds

## Problem
Merchants need ability to refund customers. Currently, refunds are manual (Finance team involvement).

## Solution
Implement automated refund API. Merchants can refund via dashboard. Accounting automatically notified.

[Rest of spec...]
```

**Step 3: CLI Setup**

```bash
agentic-agent platform init
agentic-agent platform add-feature \
  --name "Card Refunds" \
  --initiative "PAY-2025" \
  --context-pack "cp-v2" \
  --fanout "payments-service,accounting-service"
```

**Step 4: Each Team Generates Component Spec (Prompt 2)**

```
Prompt for payments-service:
"Generate component spec for payments-service from this platform spec: [paste]"

AI Output:
# Component Spec: Payments Service

## API Changes
POST /refunds
{
  "charge_id": "...",
  "amount_cents": 9999
}

[Rest of spec...]
```

**Step 5: CLI Creates OpenSpec**

```bash
cd payments-service/
agentic-agent openspec init "Process refunds" \
  --from ../platform-specs/card-refunds/spec.md
```

**Step 6: Each Team Generates Impl Spec (Prompt 3)**

```
Prompt:
"Generate impl-spec + tasks for implementing card refunds in payments-service
from this component-spec: [paste]"

AI Output:
# Implementation Spec

## Data Model
ALTER TABLE charges ADD COLUMN refund_id VARCHAR(32);
INSERT INTO refunds (charge_id, amount_cents, status)

## Code Changes
- Create POST /refunds handler
- Call Stripe refund API
- Emit refund.completed event

## Tasks
1. Add refund endpoint
   - Acceptance: curl -X POST /refunds returns 200
2. Add Stripe integration
   - Acceptance: Refund appears in Stripe dashboard
[Rest...]
```

**Step 7: CLI Tracks Implementation**

```bash
agentic-agent task claim <task-1>
# ... code ...
agentic-agent task complete <task-1>
# Continue for all tasks
```

**Step 8: Verifier Generates Test Plan (Prompt 4)**

```
Prompt:
"Generate a verification plan (verify.md) for card refunds feature
from this feature-spec and impl-spec: [paste both]"

AI Output:
# Verification Report: Card Refunds

## AC1: Refund API accepts valid request
Status: PASS
Evidence: Test test_refund_charge_success PASSED
Output: curl -X POST /refunds -d '{"charge_id":"...", "amount_cents":9999}' returns 200

[Rest...]
```

**Step 9: CLI Completion**

```bash
agentic-agent openspec complete card-refunds
agentic-agent platform change-priority \
  --initiative "PAY-2025" \
  --priority "Done"
```

---

## FAQ: Using Prompts with SDD

**Q: Should I always use prompts?**
A: No. Use prompts for specs (saving time), use CLI for tracking. Don't use prompts for simple bugs.

**Q: Can I use ChatGPT instead of Claude?**
A: Yes. Adjust prompts slightly for ChatGPT style, but logic is the same.

**Q: What if AI-generated spec is wrong?**
A: Edit before using CLI. Specs are templates, adjust to your needs.

**Q: Can I use prompts without CLI?**
A: Yes, but you lose progress tracking. Using both is best.

**Q: How detailed should prompts be?**
A: More context → better output. Always paste examples and context.

**Q: What if I don't know what risk level is?**
A: Start with USAGE-LEVELS.md decision tree, then use appropriate prompts.

---

## Summary: Prompts + CLI

| Task | Tool |
|------|------|
| Decide usage level | USAGE-LEVELS.md |
| Choose level for small project | SMALL-PROJECTS.md |
| Generate proposal + tasks | Prompt 1 (OpenSpec) |
| Generate platform spec | Prompt 1-2 (Platform) |
| Generate impl spec + tasks | Prompt 3 (Developer) |
| Generate test plan | Prompt 4 (Verifier) |
| Create/track specs | CLI commands |
| Mark complete | CLI commands |

**Use both for best results: AI for design, CLI for tracking.**

---

## Next Steps

1. **For your first feature:**
   - Read USAGE-LEVELS.md → Decide level
   - Use appropriate prompt(s) above
   - Use CLI to create + track
   - Reference SMALL-PROJECTS.md if worried about overhead

2. **For ongoing features:**
   - Keep PROMPTS-GUIDE.md handy
   - Copy relevant prompts
   - Customize with your context
   - Track with CLI

3. **For your team:**
   - Share these three docs:
     - USAGE-LEVELS.md (which level to use)
     - SMALL-PROJECTS.md (don't over-engineer)
     - PROMPTS-GUIDE.md (how to generate specs)
   - Let team choose tool based on project size
