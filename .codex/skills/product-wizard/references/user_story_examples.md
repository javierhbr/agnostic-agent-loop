# User Story Examples & Best Practices

## Format

### Standard Template
```
As a [user type/persona],
I want to [perform action/use feature],
So that [achieve benefit/value].

Acceptance Criteria:
- [ ] Given [context], when [action], then [expected outcome]
- [ ] Given [context], when [action], then [expected outcome]
```

### Alternative: Job Story (JTBD)
```
When [situation],
I want to [motivation],
So I can [expected outcome].
```

## INVEST Criteria

Good user stories are:
- **I**ndependent: Developed and delivered separately
- **N**egotiable: Details can be discussed
- **V**aluable: Clear value to users or business
- **E**stimable: Can be sized by the team
- **S**mall: Completed in one sprint
- **T**estable: Has clear acceptance criteria

## Examples by Domain

### E-Commerce: Product Search

> As an online shopper,
> I want to filter products by price range,
> So that I can find items within my budget.

**Acceptance Criteria:**
- [ ] Given I'm on the listing page, when I set min/max price, then only products in range display
- [ ] Given I've applied a filter, when I clear it, then all products show
- [ ] Given min > max, when I apply, then I see an error message
- [ ] Filter persists across pagination
- [ ] Displays count of matching products

**Priority**: P0 | **Points**: 5

### SaaS: Team Collaboration

> As a project manager,
> I want to assign tasks to team members,
> So that everyone knows their responsibilities.

**Acceptance Criteria:**
- [ ] Given I'm viewing a task, when I click "Assign", then I see team members
- [ ] Given I select a member, when I confirm, then they receive a notification
- [ ] Given a task is assigned, when viewing the list, then assignee is visible
- [ ] I can assign multiple people to one task
- [ ] I can change or remove assignments

**Priority**: P0 | **Points**: 5

### Mobile: Offline Mode

> As a mobile user with unreliable connectivity,
> I want to access recently viewed content offline,
> So that I can continue using the app without internet.

**Acceptance Criteria:**
- [ ] Given I viewed content online, when offline, then last 50 items are accessible
- [ ] Given I'm offline making changes, when I reconnect, then changes sync automatically
- [ ] Offline indicator appears when connectivity is lost
- [ ] Cached content auto-clears after 7 days

**Priority**: P1 | **Points**: 13

### Security: Two-Factor Authentication

> As a security-conscious user,
> I want to enable two-factor authentication,
> So that my account is protected from unauthorized access.

**Acceptance Criteria:**
- [ ] Given I enable 2FA, then I choose between SMS and authenticator app
- [ ] Given authenticator, when I scan QR code, then I must enter verification code
- [ ] Given 2FA enabled, when I log in, then I'm prompted for second factor
- [ ] I receive backup codes when activating 2FA
- [ ] I can disable 2FA with current password + 2FA code

**Priority**: P0 | **Points**: 13

## Acceptance Criteria Patterns

### Given-When-Then (for complex logic)
```
Given [initial context/state],
When [action/event],
Then [expected outcome].
```

### Checklist (for straightforward requirements)
```
- [ ] Requirement 1
- [ ] Requirement 2
- [ ] Edge case handling
```

### Table (for multiple scenarios)
| Condition | Action | Expected Result |
|-----------|--------|----------------|
| Valid email | Click "Send" | Confirmation message |
| Invalid email | Click "Send" | Error message |
| Empty field | Click "Send" | "Required field" error |

## Common Mistakes

| Mistake | Bad Example | Good Example |
|---------|------------|-------------|
| Too technical | "Use Redis caching with 10-min TTL" | "Pages load in under 2 seconds" |
| Too vague | "App should be fast" | "Search results appear in under 1 second" |
| Missing "why" | "I want to upload profile pictures" | "...so other users can recognize me" |
| Multiple actions | "Create account, set profile, invite team" | Split into 3 stories |
| No acceptance criteria | Just the story | Add 3–5 testable criteria |

## Story Splitting Techniques

1. **By workflow steps**: "Book a flight" → Search, Select, Enter details, Pay, Confirm
2. **By persona**: "Manage subscriptions" → Free user views plans, Paid user upgrades, Admin manages team
3. **By business rules**: "Apply discounts" → Percentage off, Fixed amount, Free shipping, Expired codes
4. **By data variations**: "Import contacts" → CSV, Google Contacts, Outlook, LinkedIn
5. **By CRUD**: "Manage projects" → Create, Read, Update, Delete
6. **By happy path / edge cases**: MVP happy path first, edge cases as follow-up stories

## Story Sizing Reference

| Points | Effort | Example |
|--------|--------|---------|
| 1–2 | Few hours | Simple UI change, copy update |
| 3–5 | 1–2 days | New form with validation, simple API |
| 8 | 3–5 days | Complex form with logic, third-party integration |
| 13+ | 1+ weeks | **Too large — split the story** |
