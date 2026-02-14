# Scenario 1: Prompts for Vague Idea

## How to Use These Prompts

1. **Copy the initial prompt** and paste it into your AI conversation
2. **Answer the discovery questions** the AI asks
3. **Use follow-up prompts** to refine or clarify
4. **Adapt placeholders** like `[your context]` to your situation

---

## Initial Prompt

```
I have an idea for mobile push notifications for our project management platform. 
We currently only have email notifications and users are missing important updates.

Can you help me create a comprehensive PRD using the product-wizard skill?
```

### Why This Works

- **States the idea clearly:** Mobile push notifications
- **Provides context:** Current state (email-only) and problem (users missing updates)
- **Names the skill explicitly:** `product-wizard`
- **Asks for comprehensive output:** A full PRD, not just notes

---

## Discovery Phase Prompts

The AI will ask clarifying questions. Here are example responses:

### When asked about the problem:

```
Users miss project deadlines and important updates because email notifications 
get lost in their inboxes. Our Customer Success team reports this in 30% of 
support tickets. Sales says "better notifications" is a top-3 feature request 
from prospects.
```

### When asked about target users:

```
We have three main user types:
1. Project Managers - need updates on task completions, blockers, and deadlines
2. Team Members - need task assignments and @mentions
3. Stakeholders - need high-level project status updates

Project Managers are our priority for this release.
```

### When asked about success metrics:

```
Success means:
- 40% reduction in "missed deadline" support tickets within 3 months
- 60% of active users opt-in to push notifications within first month
- Average notification open rate above 25%

We'll measure this through our analytics platform and support ticket tagging.
```

### When asked about constraints:

```
Constraints:
- No budget for native iOS/Android apps right now
- Must work with our existing web platform
- Need to launch MVP in 8 weeks
- Must comply with browser notification permissions (no spam)
- Development team is 3 engineers
```

### When asked about scope boundaries:

```
Out of scope for MVP:
- SMS/text notifications
- Slack/Teams integrations
- Notification scheduling/digest modes
- Custom notification sounds
- In-app notification center

These may come later, but not in the first release.
```

---

## Follow-Up Prompts

### If you need to clarify technical approach:

```
For delivery, we want to use Progressive Web App (PWA) push notifications 
through the browser, not native apps. This works on both desktop and mobile 
browsers.
```

### If priorities are unclear:

```
Let's prioritize Project Managers first. They're our power users and feel the 
pain most acutely. Team Members can come in Phase 2.
```

### If success metrics need refinement:

```
Actually, let's change the success metric to focus on engagement: 
- 50% of Project Managers enable notifications within 2 weeks
- 30% daily active notification interaction rate

These are more actionable than support ticket reduction.
```

---

## Creating the Development Plan

After the PRD is created, use this prompt:

```
Now that we have a PRD, can you create a development plan using the dev-plans 
skill? Please break this into phases with clear review checkpoints.
```

---

## Generating Acceptance Tests

Once you have tasks, use this prompt:

```
For task [TASK-ID], please generate acceptance tests using the atdd skill. 
I want to ensure we test the acceptance criteria before implementing the feature.
```

Replace `[TASK-ID]` with the actual task identifier from your task list.

---

## Tips for Effective Prompting

### ✅ Do

- **Be specific about your context:** Company size, user count, current state
- **Mention constraints upfront:** Timeline, budget, technical limitations
- **Provide real numbers:** User counts, metrics, percentages
- **Name the skill explicitly:** "using the product-wizard skill"
- **Ask follow-up questions:** If something isn't clear, probe deeper

### ❌ Don't

- **Be vague:** "We need notifications" (which kind? for what?)
- **Assume technical knowledge:** Explain your tech stack if relevant
- **Skip the "why":** Always explain the problem you're solving
- **Ignore constraints:** Unlimited scope leads to unrealistic PRDs
- **Rush through discovery:** The questions exist for a reason

---

## Troubleshooting Prompts

### If the AI doesn't understand your domain:

```
Let me provide more context about our platform: 
[describe your product, users, and key workflows]
```

### If the PRD seems too generic:

```
This feels generic. Let me be more specific about our users: 
[detailed user personas and their specific pain points]
```

### If technical details are missing:

```
For technical context: we use React frontend, Node.js backend, and PostgreSQL. 
We need push notifications to integrate with our existing notification system.
```

### If scope is too broad:

```
That scope is too large for our 8-week timeline. Let's focus only on 
[specific subset] for MVP. The rest can be Phase 2.
```

---

## Example Full Conversation Flow

1. **You:** Initial prompt (above)
2. **AI:** Asks 4-6 discovery questions
3. **You:** Answer each question with specific details
4. **AI:** Presents draft PRD sections
5. **You:** "The success metrics section looks good, but can you make the user stories more specific to Project Managers?"
6. **AI:** Refines user stories
7. **You:** "Perfect. Please save this PRD to `.agentic/spec/prd-mobile-notifications.md`"
8. **AI:** Saves PRD
9. **You:** "Now create a development plan using the dev-plans skill"
10. **AI:** Creates phased development plan

**Total time:** 30-45 minutes for a complete, production-ready PRD

---

[← Back to Scenario](scenario.md) | [View Example Conversation →](conversation.md)
