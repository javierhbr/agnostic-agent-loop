# Scenario 4: Prompts for Codebase Documentation

## Initial Prompt

```
Our payment system code is complex and needs refactoring. I've described the 
current system and what we want to change. Can you help create a PRD for this 
refactoring project using the product-wizard skill?

Current system: [describe payment flow]
Must preserve: [existing payment methods]
Want to add: [new capabilities]
```

## Preservation vs. Change

```
These payment methods must continue working exactly as they do:
- Credit/debit cards (Stripe integration)
- ACH bank transfers
- PayPal

After refactoring, we want to add:
- Apple Pay
- Google Pay
- Klarna (buy now pay later)
```
