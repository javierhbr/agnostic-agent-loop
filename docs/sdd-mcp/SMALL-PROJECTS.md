# SDD for Small Projects

**Small projects don't need the full SDD ceremony.** This guide shows how to use SDD *principles* and the CLI with minimal overhead for teams of 1-3 people working on simple features.

---

## TL;DR

| Project Size | Team | Risk | Approach | Timeline |
|--------------|------|------|----------|----------|
| **Tiny** | 1 person | Low (bug, small feature) | Skip specs, use `agentic-agent task` only | Hours-1 day |
| **Small** | 1-2 people | Low-Medium (new feature, 1 service) | Use `agentic-agent openspec init` + tasks | 1-2 weeks |
| **Medium** | 2-3 people | Medium (2-3 services) | Use Platform level (skip Analyst) | 2-3 weeks |

---

## For Tiny Projects (Solo Developer, Low Risk)

**Examples:** Bug fix, simple refactor, internal feature, small optimization

### Skip Everything, Use Task Tracking

```bash
# Just track your work with tasks
agentic-agent task create "Fix login timeout bug"
agentic-agent task claim <task-id>
# ... code ...
agentic-agent task complete <task-id>
```

**That's it.** No specs needed. Move fast.

### When to Add Minimal Ceremony

If the bug/feature is complex (> 3 tasks), write a 1-paragraph summary:

```bash
mkdir -p .agentic/specs/my-feature
cat > .agentic/specs/my-feature/README.md << 'EOF'
# Login Timeout Fix

## Problem
Users are logged out after 15 minutes of inactivity, even on idle tabs.

## Solution
Keep refresh token alive with silent renewal on focus.

## Tasks
- [ ] Add focus listener to renew token
- [ ] Test with 30-min idle
- [ ] Update docs
EOF
```

**Done.** No proposal, no gates, no verification.

---

## For Small Projects (1-2 People, 1 Service)

**Examples:** New feature (checkout flow), API endpoint, dashboard, internal tool

### Minimal OpenSpec Workflow

```bash
# 1. Create a lightweight spec
agentic-agent openspec init "User Dashboard" \
  --from README.md \
  --type feature

# 2. This creates:
# - proposal.md (what you're building)
# - tasks.md (how you're breaking it down)

# 3. You write a simple proposal (1-2 paragraphs)
vim .agentic/specs/user-dashboard/proposal.md

---
# User Dashboard

## What
Display user profile, recent orders, preferences in one view.

## Why
Users currently click through 3 pages to see their data.
Goal: Reduce time-to-profile-view from 4 clicks to 1.

## Acceptance
- [ ] Dashboard loads in < 500ms
- [ ] Shows profile, orders, preferences
- [ ] Mobile responsive
---

# 4. Decompose into tasks
vim .agentic/specs/user-dashboard/tasks.md

---
# Tasks

## Backend (4 tasks)
1. Create `/api/user/dashboard` endpoint
   - GET /api/user/dashboard
   - Returns: profile + recent_orders + preferences
   - Status code: 200

2. Add caching (Redis, 5-min TTL)
3. Write integration tests
4. Performance test (should be < 500ms)

## Frontend (5 tasks)
5. Create Dashboard page component
6. Add profile section
7. Add orders section
8. Add preferences section
9. Mobile responsive styles

## QA (2 tasks)
10. Test on mobile devices
11. Test with slow network (DevTools throttle)
---

# 5. Work through tasks
agentic-agent task claim <task-1>
# ... implement ...
agentic-agent task complete <task-1>

# Continue for all tasks

# 6. Mark feature done
agentic-agent openspec complete user-dashboard
agentic-agent openspec archive user-dashboard
```

### That's Your Entire Small Project Process

**No gates. No verification. No specs.**

Just:
1. ✅ Problem statement (1-2 sentences)
2. ✅ Tasks decomposed (what you're building)
3. ✅ Tasks completed
4. ✅ Archive

### When to Add Minimal Gates

If your feature touches **external systems** (payments, auth, integrations), add one gate check:

```bash
# Before marking complete, ask:
# ✓ Did you test all tasks?
# ✓ Any blocking issues from other teams?
# ✓ Safe to deploy?

agentic-agent openspec check user-dashboard
```

---

## For Small Teams (2-3 People, 2-3 Services)

**Examples:** Multi-service feature (payment flow, auth integration), cross-team work

### Lightweight Platform + OpenSpec

Don't do full SDD. Skip Analyst. Do light Architect + Developer.

```bash
# 1. Team lead writes a quick platform spec (1 page)
agentic-agent platform init

# 2. Create the feature at platform level
agentic-agent platform add-feature \
  --name "Guest Checkout" \
  --initiative "CHECKOUT-2025" \
  --context-pack "cp-v2" \
  --fanout "checkout-service,payments-service"

# 3. Team lead fills in spec.md (brief, < 500 words)
vim platform-specs/guest-checkout/spec.md

---
# Guest Checkout

## Problem
30% of cart abandonment is from users unwilling to create account.

## Solution
Enable checkout without account creation.

## What Changes
- checkout-service: Add guest flow, collect email only
- payments-service: Process guest charges

## Acceptance Criteria
- [ ] Guest can checkout without account
- [ ] Email collected and sent confirmation
- [ ] Charge processed successfully
- [ ] Cart abandonment drops < 25% (measure in 2 weeks)

## No ADRs, No Complex Design
Just: Problem + Solution + Who Changes + How to Verify

---

# 4. Each component team uses OpenSpec
cd checkout-service/

agentic-agent openspec init "Guest checkout flow" \
  --from ../platform-specs/guest-checkout/spec.md

# Team writes a simple proposal + tasks
vim .agentic/specs/guest-checkout/proposal.md

---
# Guest Checkout Flow

## UI Changes
1. Add "Guest Checkout" button on login prompt
2. Skip account creation, collect email
3. Show confirmation message

## Backend Changes
1. Create POST /checkout/guest endpoint
2. Accept: email, cart_id
3. Skip user_id requirement
4. Call payments-service

## Acceptance
- [ ] Guest flow complete
- [ ] Email confirmation sent
- [ ] Tests pass
---

# 5. Teams implement independently
agentic-agent task claim <task-1>
# ... code ...
agentic-agent task complete <task-1>

# 6. Mark features complete
agentic-agent openspec complete guest-checkout
agentic-agent openspec complete guest-payments

# 7. Platform team tracks completion
agentic-agent platform change-priority \
  --initiative "CHECKOUT-2025" \
  --priority "Done"
```

### Key Differences from Full SDD

| Full SDD | Small Project SDD |
|----------|-------------------|
| Analyst phase (discovery) | Skip (assume requirements known) |
| Architect writes detailed specs | Team lead writes 1-page spec |
| 5 gates enforced | 1 gate: "Did everyone test?" |
| Verifier checks all ACs | Skip (team tests as they go) |
| 3-4 weeks | 2-3 weeks |

---

## Pattern: The "Light SDD" Workflow

Perfect for small teams, quick features:

```
┌─────────────────────────────────────┐
│  1. Write 1-page platform spec      │
│     (if multi-service)              │
└─────────────────────────────────────┘
              ↓
┌─────────────────────────────────────┐
│  2. Each team: openspec init        │
│     + simple proposal + tasks       │
└─────────────────────────────────────┘
              ↓
┌─────────────────────────────────────┐
│  3. Implement tasks in parallel     │
│     (no gates, just get it done)    │
└─────────────────────────────────────┘
              ↓
┌─────────────────────────────────────┐
│  4. Quick verification              │
│     "Did we build what was asked?"  │
└─────────────────────────────────────┘
              ↓
┌─────────────────────────────────────┐
│  5. Archive & move to next feature  │
└─────────────────────────────────────┘
```

---

## Real Example: Small Team Adding Search

**Team:** 2 developers
**Feature:** Add search to product catalog
**Risk:** Low (single service)
**Timeline:** 1 week

### Day 1: Write Spec & Tasks

```bash
agentic-agent openspec init "Product Search" --from requirements.md

# Write proposal
vim .agentic/specs/product-search/proposal.md

---
# Product Search

## What
Add search box to product listing page.
Users can search by name, category, price range.

## Why
Users spend 5+ minutes scrolling to find products.
Goal: Reduce to < 1 minute with search.

## Acceptance
- [ ] Search box appears on product page
- [ ] Results filter in < 500ms
- [ ] Works with name, category, price filters
- [ ] Mobile responsive
---

# Write tasks
vim .agentic/specs/product-search/tasks.md

---
# Tasks

## Backend (3 tasks)
1. Add Elasticsearch index for products (or use existing DB full-text search)
2. Create GET /api/products/search?q=...&category=...&price_min=...&price_max=...
3. Write tests (unit + integration)

## Frontend (4 tasks)
4. Add search box component
5. Add filter inputs (category, price range)
6. Connect to API, show results in real-time
7. Mobile responsive styles

## QA (1 task)
8. Manual testing on mobile + desktop
---
```

### Days 2-6: Build

```bash
# Dev 1 works on backend
agentic-agent task claim backend-elasticsearch
# ... implement ...
agentic-agent task complete backend-elasticsearch

# Dev 2 works on frontend in parallel
agentic-agent task claim frontend-search-box
# ... implement ...
agentic-agent task complete frontend-search-box

# Continue...
```

### Day 7: Finish & Deploy

```bash
# Quick verification
# - Search works? ✓
# - Filters work? ✓
# - Mobile ok? ✓
# - Tests pass? ✓

agentic-agent openspec complete product-search
agentic-agent openspec archive product-search

# Done!
```

---

## CLI Shortcuts for Small Projects

### Don't Use These (Overkill)

```bash
# Skip these for small projects:
agentic-agent sdd start ...          # (full SDD only)
agentic-agent platform init ...      # (only if multi-service)
agentic-agent sdd adr create ...     # (only if blocking decision)
agentic-agent sdd gate-check ...     # (only if high-risk)
```

### Do Use These (Useful)

```bash
# Use these:
agentic-agent openspec init <name>   # Create spec + tasks
agentic-agent task list              # See all tasks
agentic-agent task claim <id>        # Start working
agentic-agent task complete <id>     # Mark done
agentic-agent openspec status <id>   # Check progress
agentic-agent openspec complete <id> # Finish feature
agentic-agent openspec archive <id>  # Clean up
```

---

## Boosting Small Projects with Superpowers Plugin

The optional **Superpowers plugin** adds methodology gates, TDD enforcement, and systematic debugging to any small project workflow. Use it when you want stricter quality control without full SDD complexity.

### When to Use Superpowers for Small Projects

| Scenario | Use Superpowers? | Why |
|----------|------------------|-----|
| Solo developer, quick bug fix | ❌ No | Overkill. Use plain `task` workflow |
| 1-2 person team, unclear requirements | ✅ Yes | Brainstorming + planning gates help |
| New feature with critical logic (payments, auth) | ✅ Yes | TDD enforcement prevents bugs |
| Parallel tasks (2-3 devs) | ✅ Yes | Git worktrees for isolation |
| Subtle bug that's hard to debug | ✅ Yes | Systematic debugging (4-phase protocol) |
| Final verification before deploy | ✅ Yes | Hard gate requiring evidence |
| Fast iteration on clear requirements | ❌ No | CLI Ralph loop is faster |

### Small Project Workflow with Superpowers

**For features with **unclear requirements** or **critical logic**:**

```bash
# 1. Brainstorm to clarify requirements
/brainstorming

# 2. Create plan with human checkpoints
/writing-plans

# 3. Initialize spec-driven change
agentic-agent openspec init "Feature" --from prd.md

# 4. Claim task
agentic-agent task claim TASK-ID

# 5. Implement with TDD for critical sections
/test-driven-development

# 6. Iterative implementation
/executing-plans

# 7. Final verification gate
/verification-before-completion

# 8. Complete
agentic-agent openspec complete
```

**For features with **clear specs** (faster path):**

```bash
# Skip Superpowers brainstorming, use CLI directly
agentic-agent openspec init "Feature" --from requirements.md
agentic-agent task claim TASK-ID

# Choose:
# - For simple logic: /ralph-loop
# - For critical logic: /test-driven-development
# - For bugs: /systematic-debugging

agentic-agent task complete
agentic-agent openspec complete
```

### Superpowers Skills for Small Projects

| Skill | Use Case | Time Cost |
|-------|----------|-----------|
| `superpowers:brainstorming` | Vague idea → clear spec | 20-30 min |
| `superpowers:writing-plans` | Create detailed task plan | 15-30 min |
| `superpowers:test-driven-development` | Write tests before code (critical logic) | +5-10 min |
| `superpowers:using-git-worktrees` | Parallel isolated work | Setup only |
| `superpowers:executing-plans` | Iterative task execution with checkpoints | Built-in verification |
| `superpowers:systematic-debugging` | Debug subtle issues | Variable |
| `superpowers:verification-before-completion` | Hard gate before marking done | 5-10 min |

### Example: Small Team Using Superpowers

**Team:** 2 developers
**Feature:** Add payment processing (critical logic)
**Approach:** Use Superpowers for brainstorming + TDD

```bash
# Day 1: Planning
/brainstorming                      # Clarify payment flow (20 min)
/product-wizard                     # Create PRD (15 min)
/writing-plans                      # Create implementation plan (25 min)

# Day 2-3: Implementation
agentic-agent openspec init "Payment Processing" --from prd.md
agentic-agent task claim TASK-1

# For critical payment logic:
/test-driven-development            # TDD enforcement
# Dev writes tests first
# Then implements only what tests require
# Refactors for quality

agentic-agent task complete TASK-1

# Continue for all payment-related tasks

# Day 4: Verification
/verification-before-completion     # Hard gate
# Checks: all tests pass, coverage good, no edge cases missed

agentic-agent openspec complete
```

### Decision Matrix: CLI vs Superpowers

```
┌─────────────────────────────────────────────────────┐
│  Clear spec + simple logic?                         │
│  YES → Use CLI (agentic-agent + /ralph-loop)        │
│  NO → Continue...                                   │
└─────────────────────────────────────────────────────┘
        ↓
┌─────────────────────────────────────────────────────┐
│  Critical business logic (payment/auth/PII)?         │
│  YES → Use TDD (superpowers or in-house)            │
│  NO → Continue...                                   │
└─────────────────────────────────────────────────────┘
        ↓
┌─────────────────────────────────────────────────────┐
│  Vague requirements or unclear scope?                │
│  YES → Use Superpowers brainstorming                │
│  NO → CLI workflow is fine                          │
└─────────────────────────────────────────────────────┘
        ↓
┌─────────────────────────────────────────────────────┐
│  Parallel work (2+ devs)?                           │
│  YES → Use Superpowers worktrees for isolation      │
│  NO → CLI task tracking is enough                   │
└─────────────────────────────────────────────────────┘
        ↓
┌─────────────────────────────────────────────────────┐
│  Subtle bug / hard to debug?                        │
│  YES → Use Superpowers systematic-debugging         │
│  NO → CLI workflow + tests usually sufficient       │
└─────────────────────────────────────────────────────┘
```

### Key Insight: Superpowers for Quality, CLI for Speed

- **CLI tools:** Fast automation, multi-agent support, git traceability
- **Superpowers skills:** Strict gates, TDD enforcement, human checkpoints

**Use both:** CLI for speed + Superpowers for quality on risky parts.

---

## When Small Projects Need More Ceremony

Upgrade from "Light SDD" if:

✅ **Risk increases:** Feature touches payments, auth, PII
✅ **Team grows:** Add a third developer → coordination matters
✅ **Complexity grows:** > 10 tasks or circular dependencies
✅ **Cross-team:** Multiple services owned by different teams
✅ **Blocking decisions:** Need ADR (e.g., "which database?")

**Then move to:** Platform level (skip Analyst) or SDD Full (if critical risk)

---

## FAQ: Small Projects

**Q: Do I really need a spec for a 3-task bug fix?**
A: No. Use `agentic-agent task` only. Specs are for features that take > 1 week.

**Q: Can I skip tasks.md and just use agentic-agent task?**
A: Yes! Use either:
- Option A: `openspec init` (creates proposal + tasks.md)
- Option B: `agentic-agent task create` directly (skip spec entirely)

Choose what feels right for your team.

**Q: What if I'm a solo developer?**
A: Use `agentic-agent task` to track your work. Skip specs entirely unless you need to remember context later.

**Q: Can I add gates to a small project?**
A: Yes, but it's overkill. Just test before marking complete. Gates are for high-risk features.

**Q: Should I use feature branches with small projects?**
A: Yes, same as always. OpenSpec doesn't change your git workflow.

**Q: How do I integrate small project specs with larger platform specs?**
A: Reference them:
```bash
agentic-agent openspec init "My Feature" \
  --from ../platform-specs/feature-name/spec.md \
  --implements "feature-name"
```

**Q: What if requirements change mid-project?**
A: Update proposal.md and tasks.md as you go. No approval needed. This is agile + SDD.

**Q: Can I skip testing on small projects?**
A: Testing isn't part of SDD (it's part of software engineering). Write tests regardless of project size. Verification just automates checking them.

---

## Template: Small Project Spec (Copy & Paste)

```markdown
# [Feature Name]

## Problem
[What user problem does this solve? 1-2 sentences]
[Metric: how will you measure success?]

## Solution
[What you're building]

## What Changes
- Service 1: [what]
- Service 2: [what]
(optional, skip if single service)

## Acceptance Criteria
- [ ] [Testable outcome]
- [ ] [Testable outcome]
- [ ] [Testable outcome]

## Tasks
1. [Task description]
   - Acceptance: [how to verify]
2. [Task description]
   - Acceptance: [how to verify]
(Continue for all tasks)

## Notes
(anything else the team should know)
```

**That's all you need.**

---

## Template: Small Project Tasks (Copy & Paste)

```markdown
# Tasks

## Backend (N tasks)
1. [Task name]
   - Files: [what you'll modify]
   - Acceptance: [pass test: X]

2. [Task name]
   - Files: [what you'll modify]
   - Acceptance: [pass test: Y]

## Frontend (N tasks)
3. [Task name]
   - Files: [what you'll modify]
   - Acceptance: [visual check: X]

4. [Task name]
   - Files: [what you'll modify]
   - Acceptance: [visual check: Y]

## Testing (N tasks)
5. [Task name]
   - Acceptance: [tests pass]
```

---

## Small Project Decision Tree

```
What are you building?

├─ Bug fix or small refactor?
│  └─ Just use: agentic-agent task create + task complete
│     (No spec needed)
│
├─ New feature (1 service, 1-2 weeks)?
│  └─ Use: agentic-agent openspec init
│     └─ Write: 1-paragraph proposal + 5-10 tasks
│     └─ Implement: task by task
│     └─ Complete: openspec complete + archive
│
└─ Multi-service or complex feature?
   └─ Use: Light Platform level
      └─ Write: 1-page platform spec
      └─ Fanout: to component teams
      └─ Each team: openspec init
      └─ Complete: each component + platform priority
```

---

## Summary: SDD for Small Projects

| Aspect | Big Project | Small Project |
|--------|-------------|---------------|
| Specs | Detailed, gated | 1 paragraph |
| Phases | 4+ (analyst → verifier) | 1 (build) |
| Gates | 5 gates enforced | Skip gates |
| Verification | Evidence-based | Team testing |
| Timeline | 3-4 weeks | 1-2 weeks |
| Ceremony | High (justified by risk) | Low (move fast) |

**Use SDD principles (break into tasks, track progress, verify) without the overhead.**

---

## Next Steps

1. **Tiny projects:** Skip this guide, use `agentic-agent task` only
2. **Small projects:** Use template above, create openspec spec + tasks
3. **Growing complexity:** Move to Platform level when needed
4. **Monorepo projects:** See [MONOREPO-OPENSPEC.md](./MONOREPO-OPENSPEC.md) for multi-package workflows

**Don't over-engineer.** Start simple, add ceremony only when risk/complexity justifies it.

---

## For Monorepo Teams

If you're working in a monorepo with multiple packages/services, see the dedicated guide: **[MONOREPO-OPENSPEC.md](./MONOREPO-OPENSPEC.md)**

It covers:
- Root config with multi-path spec resolution
- Per-package spec organization
- 7 implementation prompts
- Multi-service coordination
- Parallel team workflows
