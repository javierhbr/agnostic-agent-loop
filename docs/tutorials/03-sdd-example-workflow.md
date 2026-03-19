# SDD Example Workflow: Validated Customer Email Updates

This tutorial walks through a complete SDD workflow using the unified 5-phase methodology. The running example is a real-world change: allowing customers to update their email address with proper validation, contract versioning, and cross-service coordination.

---

## Context

The customer identity platform currently does not allow email updates after account creation. Product has identified this as a top-5 support ticket driver. The change touches profile-service (owner of customer data) and auth-service (owner of login credentials), requires a shared contract update, and must coordinate deployment across both services.

**Key identifiers used throughout:**

| ID | Type | Description |
| ---- | ------ | ----------- |
| PLAT-123 | Platform issue | Customer identity platform review |
| PROF-456 | Epic (profile-service) | Email update capability |
| AUTH-234 | Epic (auth-service) | Credential sync for email changes |
| chg-profile-email-validation | Change package | Tracks the full change across services |
| contracts.customer-profile.v2 | Shared contract | Updated profile contract with email history |

---

## Phase 1: Platform (Architect)

The architect establishes platform-level context before any feature work begins. This phase produces ownership artifacts, a constitution, and encoding that all downstream phases reference.

### Step 1.1: BMAD Brownfield Review

The architect conducts a brownfield review of the customer identity platform. This is not a greenfield design -- both profile-service and auth-service exist and carry legacy constraints.

Key findings:

- profile-service owns `customer` table with `email` column (no history)
- auth-service owns `credentials` table keyed by email
- No existing contract between the two services for email changes
- Email is used as a login identifier in auth-service (high coupling)

### Step 1.2: Speckit Constitution

The architect writes the platform constitution encoding validation rules, contract policies, observability standards, and API versioning conventions.

Key constitution decisions:

- **Validation**: Email changes require verification of the new address before commit
- **Contracts**: All cross-service data flows use versioned shared contracts
- **Observability**: Every state transition emits a structured event with correlation ID
- **API versioning**: Breaking changes require a new contract version; consumers get a migration window

### Step 1.3: OpenSpec config.yaml Encoding

```yaml
# openspec/config.yaml
project:
  name: customer-identity-platform
  id: PLAT-123

services:
  profile-service:
    owner: team-profile
    repo: github.com/org/profile-service
  auth-service:
    owner: team-auth
    repo: github.com/org/auth-service

contracts:
  customer-profile:
    current_version: v1
    consumers: [auth-service, notification-service]

conventions:
  branching: feature/PROF-{id}
  jira_prefix: PROF (profile-service), AUTH (auth-service)
```

### Step 1.4: Versioning and JIRA Conventions

- Platform issue: PLAT-123
- Profile-service epic: PROF-456, stories: PROF-789, PROF-790, PROF-791
- Auth-service epic: AUTH-234
- Contract versioning: customer-profile.v1 (current) to customer-profile.v2 (target)

### Step 1.5: Ownership Artifacts

The architect writes three ownership artifacts.

**component-ownership-profile-service.md** (excerpt):

```markdown
# Component Ownership: profile-service

Owner: team-profile
Domain: Customer profile data (name, email, address, preferences)
Key invariant: One active email per customer at any time
Dependencies: auth-service (downstream consumer of email changes)
Contracts published: contracts.customer-profile.v2
```

**component-ownership-auth-service.md** (excerpt):

```markdown
# Component Ownership: auth-service

Owner: team-auth
Domain: Authentication credentials and login flows
Key invariant: Credential lookup by email must resolve in < 50ms
Dependencies: profile-service (upstream publisher of email changes)
Contracts consumed: contracts.customer-profile.v2
```

**dependency-map.md** (excerpt):

```markdown
# Dependency Map

## Impact Tiers

Tier 1 (hard constraint -- must not break):
  - auth-service credential lookup by email
  - login flow availability (99.9% SLA)

Tier 2 (should not degrade):
  - notification-service email delivery
  - profile-service API response time (< 200ms p95)

Tier 3 (acceptable temporary degradation):
  - analytics dashboards (eventual consistency OK)
```

**glossary.md** (excerpt):

```markdown
# Glossary

- **email change**: A customer-initiated request to update their login email
- **verification**: Proof of ownership of the new email address via token
- **credential sync**: The process of updating auth-service after a verified email change
- **contract**: A versioned schema defining the data shape exchanged between services
```

---

## Phase 2: Assess (Team Lead)

The team lead classifies the change and opens a formal change package.

### Step 2.1: BMAD Classification

| Dimension | Value | Rationale |
| --------- | ----- | --------- |
| Greenfield/Brownfield | Brownfield | Both services exist with production traffic |
| Size | Medium | 2 services, 1 contract, 3 stories |
| Impact | High | Email is a login identifier; breaking it breaks login |
| Scope | Shared change | Crosses service boundaries with a contract update |

Classification result: **Standard workflow** with architect impact review required.

### Step 2.2: Open Change Package

```yaml
# openspec/changes/chg-profile-email-validation/platform-ref.yaml
change_package: chg-profile-email-validation
platform: PLAT-123
status: open
created: 2026-03-19
services:
  - profile-service
  - auth-service
contracts_affected:
  - contracts.customer-profile.v2
epics:
  - PROF-456
  - AUTH-234
```

```yaml
# openspec/changes/chg-profile-email-validation/jira-traceability.yaml
epics:
  PROF-456:
    stories: [PROF-789, PROF-790, PROF-791]
    service: profile-service
  AUTH-234:
    stories: []  # to be created after design review
    service: auth-service
links:
  platform_issue: PLAT-123
  change_package: chg-profile-email-validation
```

### Step 2.3: Architect Impact Review

The architect reviews the change package against the dependency map:

- **Tier 1 impact identified**: auth-service credential lookup by email is a hard constraint. The email change flow must guarantee that auth-service credentials are updated atomically or the old email continues to work until sync completes.
- **Mitigation**: Use an event-driven approach. Profile-service publishes `EmailChanged` event. Auth-service consumes it and updates credentials. During the window, both old and new email work for login.

---

## Phase 3: Specify (Product)

Product writes the proposal and delta specs that define exactly what changes.

### Step 3.1: Check Glossary, Write proposal.md

Product checks the glossary to use consistent terminology, then writes the proposal.

**proposal.md** (excerpt):

```markdown
# Proposal: Validated Customer Email Updates

## Problem
Customers cannot update their email address after account creation.
This generates 2,400 support tickets per month (top-5 driver).

## Solution
Allow customers to request an email change through the profile UI.
The new email must be verified before it replaces the old one.
Auth-service credentials sync automatically via the customer-profile.v2 contract.

## Success Criteria
- Given a customer requests an email change,
  When they verify the new email via token,
  Then their profile and login credentials update within 30 seconds.
- Given a customer has a pending email change,
  When they log in with the old email,
  Then login succeeds until the new email is verified.
- Support tickets for email changes drop by 80% within 60 days.

## Scope
- IN: Email change request, verification flow, credential sync, audit trail
- OUT: Username changes, email merge across accounts, bulk migration
```

### Step 3.2: Write Delta Specs

Delta specs describe what is ADDED, MODIFIED, or REMOVED.

**delta-spec.md** (excerpt):

```markdown
# Delta Spec: chg-profile-email-validation

## profile-service

ADDED:
  - POST /api/v1/customers/{id}/email-change-request
    - Accepts: { new_email: string }
    - Returns: { request_id: string, verification_token_sent: boolean }
  - POST /api/v1/customers/{id}/email-change-verify
    - Accepts: { request_id: string, token: string }
    - Returns: { verified: boolean, email_updated: boolean }
  - Table: email_change_requests (customer_id, old_email, new_email, token, status, created_at, verified_at)

MODIFIED:
  - contracts.customer-profile: v1 -> v2
    - ADDED field: email_change_history (array of { old, new, changed_at })
    - ADDED event: EmailChanged { customer_id, old_email, new_email, timestamp }

## auth-service

ADDED:
  - Consumer: EmailChanged event handler
    - Updates credential lookup index for new email
    - Maintains old email as alias until TTL expires (72h)

REMOVED:
  - Nothing
```

### Step 3.3: Scope and Alignment Check

Product confirms:

- Delta spec aligns with proposal success criteria
- No items outside declared scope
- Contract version bump is justified (new field + new event)
- Auth-service team (team-auth) has reviewed and accepted the consumer contract

---

## Phase 4: Plan (Architect)

The architect translates the spec into a concrete design and implementable tasks.

### Step 4.1: Write design.md

The architect reads the impact tiers from the dependency map. Tier 1 (auth-service credential lookup) is a hard constraint that shapes the design.

**design.md** (excerpt):

```markdown
# Design: Validated Customer Email Updates

## Architecture Decision
Event-driven email sync with dual-email grace period.

Profile-service is the source of truth for email changes.
Auth-service is a consumer that updates credentials asynchronously.

## Sequence
1. Customer submits email change request via profile-service API
2. Profile-service sends verification email to new address
3. Customer clicks verification link -> profile-service verifies token
4. Profile-service updates customer record and publishes EmailChanged event
5. Auth-service consumes event, adds new email to credential index
6. Auth-service maintains old email alias for 72h (Tier 1 safety)
7. After 72h TTL, old email alias is removed

## Tier 1 Constraint: auth-service
- Old email MUST continue to work for login during sync window
- If EmailChanged event processing fails, retry with exponential backoff
- Dead letter queue after 5 retries; alert team-auth on-call
- Rollback: old email always works; new email only works after sync

## Contract: contracts.customer-profile.v2
- Backward compatible: v1 consumers ignore new fields
- v2 consumers get email_change_history and EmailChanged event
- Migration window: 90 days before v1 deprecation

## Observability
- profile-service: email_change_requested, email_change_verified, email_change_failed
- auth-service: credential_sync_completed, credential_sync_failed
- Correlation ID: request_id flows through all events
```

### Step 4.2: Write tasks.md

Three tasks mapped to JIRA stories under PROF-456.

**tasks.md** (excerpt):

```markdown
# Tasks: chg-profile-email-validation

## PROF-789: Email change request and verification endpoints
Service: profile-service
Type: feature
Estimated: 3 points
Depends on: nothing
Delivers:
  - POST /email-change-request endpoint
  - POST /email-change-verify endpoint
  - email_change_requests table migration
  - Verification email sending
  - Unit and integration tests
Acceptance criteria:
  - Given valid customer ID and new email, When POST /email-change-request, Then verification email sent
  - Given valid token, When POST /email-change-verify, Then email updated and EmailChanged event published

## PROF-790: Contract v2 and EmailChanged event publishing
Service: profile-service
Type: contract
Estimated: 2 points
Depends on: PROF-789
Delivers:
  - contracts.customer-profile.v2 schema definition
  - EmailChanged event publisher
  - email_change_history field population
  - Contract compatibility tests
Acceptance criteria:
  - Given email verified, When event published, Then EmailChanged contains old_email, new_email, timestamp
  - Given v1 consumer, When reading contract, Then no breaking changes

## PROF-791: Audit trail and observability
Service: profile-service
Type: observability
Estimated: 1 point
Depends on: PROF-789
Delivers:
  - Structured logging for all email change state transitions
  - Metrics: email_change_requested, email_change_verified, email_change_failed
  - Correlation ID propagation
Acceptance criteria:
  - Given any email change operation, When it completes, Then structured log emitted with correlation_id
  - Given email change failure, When alert threshold reached, Then on-call notified
```

### Step 4.3: Developer Feasibility Review

Developer reviews tasks and confirms:

- Migration is additive (no column drops)
- Event schema is straightforward
- 72h alias TTL in auth-service is achievable with existing TTL infrastructure
- No blockers identified

---

## Phase 5: Deliver (Team Lead)

The team lead coordinates build, review, verification, and deployment.

### Step 5.1: Build

Developer claims and implements the first task.

```bash
# View available tasks
agentic-agent task list

# Claim the first task
agentic-agent task claim PROF-789

# The CLI creates a worktree at .worktrees/feature/task-PROF-789/
# Developer implements the email change endpoints, writes tests

# Validate before completing
agentic-agent validate

# Complete the task
agentic-agent task complete PROF-789
```

Repeat for PROF-790 (depends on PROF-789) and PROF-791 (parallel with PROF-790).

### Step 5.2: Create PR

```bash
agentic-agent pr create --task PROF-789
```

**PR description** (excerpt):

```markdown
## Validated Customer Email Updates - Request and Verification Endpoints

### Change Package
chg-profile-email-validation

### JIRA
- Epic: PROF-456
- Story: PROF-789

### What
- Added POST /api/v1/customers/{id}/email-change-request
- Added POST /api/v1/customers/{id}/email-change-verify
- Added email_change_requests table migration
- Added verification email integration

### Testing
- 12 unit tests covering happy path and edge cases
- 3 integration tests for the full request-verify flow
- Verified old email continues to work during pending state

### Contract Impact
- Prepares for contracts.customer-profile.v2 (shipped in PROF-790)
- No breaking changes to existing consumers

### Traceability
Platform: PLAT-123 | Epic: PROF-456 | Story: PROF-789
Change Package: chg-profile-email-validation
```

### Step 5.3: Review PR

A reviewer (different agent or team member) reviews against:

- Delta spec compliance: do the endpoints match the spec?
- Design constraints: is the Tier 1 safety (old email works) maintained?
- Observability: are structured logs present?
- Test coverage: are acceptance criteria proven by tests?

### Step 5.4: Verify

The verifier checks every acceptance criterion with evidence:

```text
AC1: Given valid customer and new email, When POST /email-change-request
     Then verification email sent
     Evidence: Integration test EmailChangeRequest_SendsVerification PASS
     Log: {"event":"email_change_requested","customer_id":"C-100","correlation_id":"req-abc"}

AC2: Given valid token, When POST /email-change-verify
     Then email updated and EmailChanged event published
     Evidence: Integration test EmailChangeVerify_UpdatesAndPublishes PASS
     Event: {"type":"EmailChanged","old":"old@example.com","new":"new@example.com"}
```

### Step 5.5: Deploy

Deployment requires coordination with AUTH-234 (auth-service must be ready to consume EmailChanged events before profile-service starts publishing them).

Deployment order:

1. Deploy auth-service with EmailChanged consumer (events will be no-ops until profile-service publishes)
2. Deploy profile-service with email change endpoints and event publisher
3. Monitor: credential_sync_completed metric should appear within minutes of first real email change
4. Verify: old email login still works during sync window

### Step 5.6: Archive

Once all stories under PROF-456 are complete and deployed:

```bash
# Archive the change package
# Move openspec/changes/chg-profile-email-validation/ to openspec-archive/
```

The archived change package preserves the full history: platform-ref, jira-traceability, proposal, delta specs, design, tasks, PR links, and verification evidence.

---

## Complete Traceability Chain

Every artifact links back to the platform issue and forward to the deployed code.

```text
PLAT-123 (platform issue)
  |
  +-- openspec/config.yaml (platform encoding)
  +-- component-ownership-profile-service.md
  +-- component-ownership-auth-service.md
  +-- dependency-map.md (impact tiers)
  +-- glossary.md
  |
  +-- chg-profile-email-validation (change package)
       |
       +-- platform-ref.yaml --> PLAT-123
       +-- jira-traceability.yaml --> PROF-456, AUTH-234
       |
       +-- proposal.md (what and why)
       +-- delta-spec.md (what changes, precisely)
       |
       +-- design.md (how, with Tier 1 constraints)
       +-- tasks.md --> PROF-789, PROF-790, PROF-791
       |
       +-- PROF-789 --> PR #42 --> commit abc1234
       +-- PROF-790 --> PR #43 --> commit def5678
       +-- PROF-791 --> PR #44 --> commit ghi9012
       |
       +-- verification evidence (per AC)
       +-- deployment record (auth-service first, then profile-service)
       |
       +-- openspec-archive/ (final resting place)
```

Reading the chain from any point:

- From a commit: which story? which task? which change package? which platform issue?
- From a platform issue: which change packages? which services? which contracts? which PRs?
- From a contract version: which change package introduced it? which consumers were notified?

---

## CLI Commands Summary

```bash
# Phase 5 -- the commands you use most often

agentic-agent task list                    # see backlog and in-progress
agentic-agent task claim PROF-789          # claim task, creates worktree
agentic-agent validate                     # run validators before completing
agentic-agent task complete PROF-789       # mark done, captures commits
agentic-agent pr create --task PROF-789    # generate PR from spec
agentic-agent pr review --task PROF-789    # spawn independent reviewer
agentic-agent pr status                    # check PR state
```

---

## Key Takeaways

1. **Platform first.** Phase 1 ownership artifacts (dependency map, glossary, constitution) prevent Phase 4/5 surprises. The Tier 1 constraint on auth-service shaped the entire design.

2. **Change packages are the unit of traceability.** Everything -- proposal, delta specs, design, tasks, PRs, verification -- lives under one change package ID that links to the platform issue and JIRA epics.

3. **Delta specs are precise.** ADDED/MODIFIED/REMOVED eliminates ambiguity about what the change does. Reviewers and verifiers reference the delta spec, not verbal descriptions.

4. **Deploy order matters for shared changes.** The consumer (auth-service) deploys before the publisher (profile-service). This is derived from the Tier 1 constraint identified in Phase 1.

5. **Archive preserves institutional knowledge.** Six months from now, when someone asks "why does auth-service support dual-email login?", the archived change package has the answer.
