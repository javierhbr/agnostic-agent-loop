# Scenario 1: Example Conversation

This is a realistic dialogue between Sarah (Product Manager) and the AI using the product-wizard skill.

---

## Initial Request

**Sarah:**
> I have an idea for mobile push notifications for our project management platform. We currently only have email notifications and users are missing important updates. Can you help me create a comprehensive PRD using the product-wizard skill?

**AI:**
> I'll help you create a robust PRD for mobile push notifications. Let me ask some discovery questions to ensure we build the right specification.
>
> **Question 1: The Problem**
> - What specific problem are you solving? Why are email notifications insufficient?
> - How do you know this is a priority issue? (support tickets, user feedback, metrics?)
> - Why now? What's changed or what's the opportunity?
>
> **Question 2: Target Users**
> - Who will use push notifications? (all users, specific roles, segments?)
> - What actions trigger notifications they care about?
> - Are there different needs for different user types?
>
> **Question 3: Success Criteria**
> - How will you know this feature succeeded?
> - What metrics will you track?
> - What's the target/baseline?

---

## Discovery Round 1

**Sarah:**
> **Problem:** Users miss project deadlines and important updates because email notifications get lost in their inboxes. Our Customer Success team reports this in 30% of support tickets. Sales says "better notifications" is a top-3 feature request from prospects. We risk churn if we don't address this.
>
> **Users:** We have three main user types:
> 1. Project Managers - need updates on task completions, blockers, and deadlines
> 2. Team Members - need task assignments and @mentions  
> 3. Stakeholders - need high-level project status updates
>
> Project Managers are our priority for this release. They manage 5-20 projects each and get overwhelmed by email.
>
> **Success:** 
> - 50% of Project Managers enable notifications within 2 weeks of launch
> - 30% daily active notification interaction rate
> - Reduction in "missed deadline" complaints

**AI:**
> Excellent context. Now let me understand constraints and scope:
>
> **Question 4: Constraints**
> - Timeline? When do you need to launch?
> - Budget/resource constraints?
> - Technical limitations? (do you have mobile apps, or is this web-only?)
> - Compliance or privacy requirements?
>
> **Question 5: Scope Boundaries**
> - What's explicitly OUT of scope for this first release?
> - What notification types are must-haves vs. nice-to-haves?
> - Any existing notification infrastructure you'll build on?

---

## Discovery Round 2

**Sarah:**
> **Constraints:**
> - Need to launch MVP in 8 weeks
> - No budget for native iOS/Android apps right now
> - Must work with our existing web platform
> - Development team is 3 engineers
> - Must comply with browser notification permissions (no spam)
>
> **Out of Scope for MVP:**
> - SMS/text notifications
> - Slack/Teams integrations  
> - Notification scheduling/digest modes
> - Custom notification sounds
> - In-app notification center
>
> **Must-Have Notifications:**
> 1. Task assigned to me
> 2. Task due in 24 hours
> 3. Project deadline approaching
> 4. @mentioned in comments
>
> **Nice-to-Have (Phase 2):**
> - Project status changes
> - New comments on tasks I'm following
> - Weekly summary digest

**AI:**
> Perfect. One final clarification:
>
> **Question 6: Technical Approach**
> - Since you don't have native apps, should we use Progressive Web App (PWA) push notifications through the browser?
> - Do you want users to control notification preferences (which types they receive)?
> - Any specific browsers or devices to prioritize or exclude?

---

## Final Clarifications

**Sarah:**
> Yes, PWA push notifications through the browser is exactly right. Works on both desktop and mobile browsers.
>
> Users must be able to control preferences - which notification types they receive and quiet hours.
>
> Priority browsers: Chrome, Safari, Firefox. Edge is nice-to-have. No need to support Internet Explorer.

**AI:**
> Excellent! I have everything I need to create a comprehensive PRD. Here's what I'll include:
>
> **PRD Sections:**
> 1. Executive Summary
> 2. Problem Statement & Context
> 3. User Personas (focus on Project Managers)
> 4. User Stories with acceptance criteria
> 5. Functional Requirements
> 6. Technical Approach (PWA push)
> 7. Success Metrics
> 8. Non-Goals (scope boundaries)
> 9. Risks & Mitigation
> 10. Roadmap (8-week timeline with phases)
>
> I'll save this to `.agentic/spec/prd-mobile-notifications.md`. Would you like me to proceed?

**Sarah:**
> Yes, please create the PRD.

**AI:**
> ✓ PRD created and saved to `.agentic/spec/prd-mobile-notifications.md`
>
> The PRD is complete with:
> - 4 priority notification types for Project Managers
> - Success metrics: 50% adoption, 30% daily interaction
> - 8-week delivery timeline
> - Technical approach using service workers and Notification API
> - User preference controls
> - 12 user stories with acceptance criteria
>
> Would you like me to create a development plan next using the dev-plans skill?

**Sarah:**
> Yes, please create a development plan with phases and review checkpoints.

**AI:**
> Creating development plan...
>
> ✓ Development plan created with 3 phases:
> - **Phase 1 (Weeks 1-3):** Infrastructure - service worker, push service, permissions
> - **Phase 2 (Weeks 4-6):** Notification triggers and delivery
> - **Phase 3 (Weeks 7-8):** User preferences and testing
>
> Each phase ends with a review checkpoint. The plan includes QA checklist for notification testing.
>
> See `DEVELOPMENT_PLAN.md` for full details.

---

## Key Observations

### What Made This Effective

1. **Specific problem statement:** 30% of support tickets, top-3 feature request
2. **Clear user prioritization:** Project Managers first, not "everyone"
3. **Measurable success criteria:** 50% adoption, 30% interaction
4. **Explicit constraints:** 8 weeks, 3 engineers, no native apps
5. **Scope boundaries:** Listed what's NOT included
6. **Technical clarity:** PWA push, specific browsers

### What Sarah Provided

- Business context (company size, user count)
- Quantified pain points (30% of tickets)
- User segmentation (3 types, prioritized)
- Realistic constraints (time, budget, tech)
- Clear must-haves vs. nice-to-haves

### What the AI Extracted

- Problem that justifies investment
- Testable success criteria
- Technical approach that fits constraints
- Phased delivery plan
- User stories with acceptance criteria

### Time Investment

- Initial prompt: 2 minutes
- Discovery questions: 15 minutes (2 rounds)
- PRD generation: 5 minutes
- Review and refinement: 10 minutes
- Dev plan creation: 5 minutes
- **Total: ~37 minutes**

---

[← Back to Prompts](prompts.md) | [View PRD Output →](outputs/prd-mobile-notifications.md)
