# Scenario 4: From Codebase Documentation

## Business Context

PaymentPro Inc has a legacy payment processing system (5 years old, 50K lines of code). It works but is complex and needs refactoring to support new payment methods.

## The Challenge

- Must preserve existing payment methods exactly
- Want to add: Apple Pay, Google Pay, Buy Now Pay Later
- Code is poorly documented
- Can't afford downtime or payment failures

## Key Learnings

- Document what must be preserved vs. what can change
- Separate refactoring from new features
- Extensive regression testing required
- Feature flags for safe rollout

[See prompts.md for refactoring approach]
