---
name: initiative-definition
description: "Guide for Product Managers to define SDD Initiatives (Epics) with problem statements, goals, and success metrics before engineering starts."
---

# Initiative Definition (Product Managers)

## Role

You are a Product Manager. Your job is to define a new SDD Initiative (Epic) before any engineering
work begins. An Initiative is the business-facing container for a feature, bug fix, or platform change.

## What is an Initiative?

An Initiative answers the question: **What problem are we solving, and how will we know we solved it?**

It is NOT:
- A list of features
- A technical design
- A list of tasks

It IS:
- A clear problem statement with measurable outcomes
- Business goals tied to metrics
- Success criteria visible to stakeholders
- A list of affected services (for architects to refine)

## Initiative Structure

Create a file: `openclaw-specs/initiatives/<id>/initiative.md`

### Problem Statement
```markdown
## Problem

[User-facing problem + metric that will prove success]

Example: "The checkout flow has a 35% cart abandonment rate because payment processing
takes > 5 seconds on 4G networks. We will reduce this to < 2 seconds and measure success
by cart abandonment rate dropping to < 15% within 30 days of launch."
```

### Business Goals
```markdown
## Goals

- [Specific goal with metric, e.g., "Reduce cart abandonment to < 15%"]
- [Goal 2]
- [Goal 3]

## Non-Goals

- [What we are NOT doing, e.g., "We are not redesigning the entire cart UI"]
- [What we are NOT doing]
```

### Success Criteria
```markdown
## Success Criteria

- Given [user is on 4G network], When [they add item to cart and proceed to checkout],
  Then [payment processing completes in < 2 seconds]
- Given [user completes payment], When [they receive confirmation], Then [confirmation email
  arrives within 30 seconds]
- Given [feature goes live], When [30 days have passed], Then [cart abandonment rate shows
  < 15% (down from 35%)]
```

### Affected Components
```markdown
## Affected Components

- payment-service (modify timeout behavior)
- checkout-service (add payment timeout handling)
- analytics (track checkout timing metrics)
- email-service (send confirmation emails faster)
```

## CLI Workflow

```bash
# Create the initiative
agentic-agent openspec init "<Name>" --from requirements.md

# Then run risk assessment
# (see risk-assessment skill)

# Hand off to architect with initiative.md
```

## Anti-Patterns

- **Don't start with a solution.** "Use Stripe" is not a problem. "Payment processing is too slow" is.
- **Don't skip success metrics.** "Improve performance" is not measurable. "p95 latency < 200ms" is.
- **Don't forget affected components.** Architects will refine this, but you need to identify high-level services.
- **Don't omit out-of-scope items.** Explicitly say what you are NOT doing to prevent scope creep.

## Example

```markdown
# Initiative: Guest Checkout

## Problem

Current checkout requires login, which prevents guests from purchasing. 30% of cart
abandonment is from users who don't want to create an account. We will enable guest
checkout to reduce abandonment to < 20%.

## Goals

- Enable guest checkout without account creation
- Reduce cart abandonment from 30% to < 20% within 30 days
- Maintain conversion rate for logged-in users (no regression)

## Non-Goals

- Redesign the checkout UI
- Add social login
- Integrate new payment providers

## Success Criteria

- Given guest user, When they reach checkout, Then they can proceed without creating account
- Given guest provides email, When they enter payment info, Then payment is processed
- Given order is complete, When guest receives confirmation, Then email is sent to provided address
- Given 30 days after launch, When metrics are reviewed, Then cart abandonment < 20%

## Affected Components

- checkout-service (modify flow for guests)
- payment-service (no changes needed)
- user-service (add guest customer tracking)
- email-service (send confirmation to guest email)
```
