# Scenario 2: Prompts for Requirements List

## Initial Prompt

```
I have a requirements list for expense approval automation. Some items conflict 
and some are vague. Can you help me turn this into a structured PRD using the 
product-wizard skill?

[Attach or paste initial-requirements.md]
```

## When AI Asks About Conflicts

```
For the email vs. Slack conflict:
- Primary: Slack notifications for approvers (they're always in Slack)
- Secondary: Email as fallback for people not in Slack
- Not needed: SMS

For automatic vs. manager approval:
- Automatic approval for expenses <$100
- Manager approval for $100-$5,000
- CFO approval for >$5,000
```

## When AI Asks for Priorities

```
Must-have (MVP):
1. Approval workflows based on amount
2. QuickBooks integration
3. Receipt attachment
4. Basic reporting

Nice-to-have (Phase 2):
1. Slack notifications
2. Offline mode
3. Foreign currency
4. HR system integration
```

[View Full Prompt Library →](../templates/prompt-library.md)
