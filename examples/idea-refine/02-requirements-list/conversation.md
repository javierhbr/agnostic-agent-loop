# Scenario 2: Example Conversation (Condensed)

## Initial Request

**CFO:**
> I have a requirements list for expense approval automation. [Pastes 18 bullet points from initial-requirements.md]
> Can you help turn this into a structured PRD using the product-wizard skill?

## AI Response - Organizing & Clarifying

**AI:**
> I see several conflicts and vague items. Let me organize these:
>
> **Conflicts:**
> - Email (#7) vs. Slack (#12) notifications
> - Automatic approval (#1) vs. Manager approval (#5)
>
> **Vague Terms:**
> - "Small amounts" - what threshold?
> - "Fast" - how fast?
> - "Right person" - what routing rules?
>
> **Questions:**
> 1. For notifications: Email, Slack, or both?
> 2. What defines automatic vs. manual approval amounts?
> 3. What are your success metrics?
> 4. Timeline and resource constraints?

## CFO Clarifies

**CFO:**
> **Notifications:** Slack primary (approvers live in Slack), email as fallback
> **Approval thresholds:** <$100 auto, $100-$5K manager, >$5K CFO
> **Success:** Reduce approval time from 3 days to 4 hours
> **Timeline:** Need by end of Q2 (12 weeks)
> **Team:** 2 developers

**AI:**
> Perfect. I'll create a PRD with:
> - Clear approval workflow rules
> - Slack + email notification strategy
> - QuickBooks integration (requirement #2)
> - Priorities: MVP vs. Phase 2
> - Success metric: 3 days → 4 hours

**Result:** Scattered bullets → Structured 15-page PRD in 35 minutes

[View Full Scenario →](scenario.md)
