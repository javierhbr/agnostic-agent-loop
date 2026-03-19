# OpenSpec + SDD Visual Diagrams

**ASCII diagrams and flowcharts for understanding workflows.**

---

## 1. Workflow Selection Decision Tree

```
┌─────────────────────────────────────────────────────────┐
│         START: What are you building?                   │
└────────────────────┬────────────────────────────────────┘
                     │
        ┌────────────┴────────────┐
        │                         │
   Bug/Hotfix?                Feature/Change?
   (< 1 hour)                 (1-4 weeks)
        │                         │
        ▼                         ▼
   ┌─────────────┐         ┌──────────────┐
   │    TASKS    │         │   OPENSPEC   │
   │ (Simplest)  │         │   (Balanced) │
   │             │         │              │
   │ • Create    │         │ • Init       │
   │ • Claim     │         │ • Proposal   │
   │ • Complete  │         │ • Tasks      │
   │             │         │ • Archive    │
   └─────────────┘         └──────┬───────┘
                                  │
                        ┌─────────┴──────────┐
                        │                    │
                   High Risk?           Monorepo?
                   (Payment/Auth)       (Multi-Svc)
                        │                    │
                        ▼                    ▼
                   ┌──────────┐         ┌──────────────┐
                   │   SDD    │         │ Multi-Path   │
                   │  (Full)  │         │   Config     │
                   │          │         │              │
                   │ • Gates  │         │ specDirs:    │
                   │ • ADRs   │         │  - packages/*│
                   │ • Phases │         │  - specs     │
                   │          │         │              │
                   └──────────┘         └──────────────┘
```

---

## 2. OpenSpec Lifecycle

```
┌──────────────────────────────────────────────────────────────────┐
│                    OPENSPEC LIFECYCLE                             │
└──────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│ 1. PLAN PHASE (Day 1)                                           │
│                                                                 │
│   Write spec file (requirements.md)                             │
│   ↓                                                             │
│   agentic-agent openspec init "Name" --from requirements.md     │
│   ↓                                                             │
│   Creates:                                                      │
│   • .agentic/openspec/changes/<feature>/proposal.md            │
│   • .agentic/openspec/changes/<feature>/tasks.md               │
└─────────────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────────────┐
│ 2. BUILD PHASE (Days 2-N)                                       │
│                                                                 │
│   ┌────────────────────────────────────────────────────────┐   │
│   │ For each task:                                         │   │
│   │                                                        │   │
│   │   agentic-agent task claim <task-id>                  │   │
│   │   ↓                                                    │   │
│   │   [Implement in your package]                         │   │
│   │   ├─ Edit code files                                  │   │
│   │   ├─ Add tests                                        │   │
│   │   └─ Commit changes                                   │   │
│   │   ↓                                                    │   │
│   │   agentic-agent task complete <task-id>              │   │
│   │   ↓                                                    │   │
│   │   Task marked done, move to next                      │   │
│   └────────────────────────────────────────────────────────┘   │
│   ↓                                                             │
│   Repeat until all tasks completed                              │
└─────────────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────────────┐
│ 3. VERIFY PHASE (Final Day)                                     │
│                                                                 │
│   agentic-agent openspec status <feature-id>                   │
│   ↓                                                             │
│   Check all acceptance criteria ✓                              │
│   Run tests one final time ✓                                   │
│   Code review complete ✓                                       │
│   ↓                                                             │
│   agentic-agent openspec complete <feature-id>                 │
│   ↓                                                             │
│   agentic-agent openspec archive <feature-id>                  │
└─────────────────────────────────────────────────────────────────┘
                            ↓
                    ✅ FEATURE DONE
```

---

## 3. Monorepo Multi-Path Spec Discovery

```
┌──────────────────────────────────────────────────────────────────┐
│         MONOREPO: Spec Discovery & Resolution                   │
└──────────────────────────────────────────────────────────────────┘

agnostic-agent.yaml (root config):
┌────────────────────────────────────┐
│ paths:                             │
│   specDirs:                        │
│     - packages/*/specs     ← Rule 1│
│     - apps/*/specs         ← Rule 2│
│     - specs                ← Rule 3│
└────────────────────────────────────┘

When you run:
  agentic-agent openspec init "Name" --from auth-oauth.md

CLI searches in order (FIRST MATCH WINS):
┌──────────────────────────────────────────────────────────────────┐
│ 1. packages/*/specs/auth-oauth.md                               │
│    ├─ packages/auth-service/specs/auth-oauth.md   ← FOUND! ✓   │
│    ├─ packages/web-ui/specs/auth-oauth.md         (skip)       │
│    └─ packages/api/specs/auth-oauth.md            (skip)       │
│                                                                  │
│ 2. apps/*/specs/auth-oauth.md        (skipped, rule 1 matched) │
│                                                                  │
│ 3. specs/auth-oauth.md               (skipped, rule 1 matched) │
└──────────────────────────────────────────────────────────────────┘

Result:
  ✓ Found: packages/auth-service/specs/auth-oauth.md
  ✓ Creates: .agentic/openspec/changes/auth-oauth/
  ✓ Tracks centrally, specs live with packages
```

---

## 4. Task vs OpenSpec vs SDD Comparison

```
┌─────────────────────────────────────────────────────────────────┐
│                                                                 │
│                      WORKFLOW TIERS                             │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘

TIER 1: TASKS (Simplest)
┌──────────────────────────────────────────────────────────────┐
│ $ agentic-agent task create "Fix X"                          │
│ $ agentic-agent task claim <id>                              │
│ $ agentic-agent task complete <id>                           │
│                                                              │
│ ├─ Use for:      Bug fixes, hotfixes, small refactors      │
│ ├─ Team size:    Solo developer                            │
│ ├─ Duration:     < 1 day                                   │
│ ├─ Spec file:    NO                                        │
│ ├─ Gates:        NO                                        │
│ └─ Setup time:   5 minutes                                 │
└──────────────────────────────────────────────────────────────┘
                            ↓ Upgrade if...
                   (> 1 week, clear requirements)

TIER 2: OPENSPEC (Balanced)
┌──────────────────────────────────────────────────────────────┐
│ $ agentic-agent openspec init "Name" --from spec.md          │
│ $ agentic-agent task claim <id>                              │
│ $ agentic-agent task complete <id>                           │
│ $ agentic-agent openspec complete <id>                       │
│                                                              │
│ ├─ Use for:      Features with acceptance criteria          │
│ ├─ Team size:    1-3 people                                │
│ ├─ Duration:     1-2 weeks                                 │
│ ├─ Spec file:    proposal.md + tasks.md                    │
│ ├─ Gates:        Optional (manual checklist)               │
│ └─ Setup time:   10 minutes                                │
└──────────────────────────────────────────────────────────────┘
                            ↓ Upgrade if...
                   (High risk: payments, auth, PII)

TIER 3: SDD (Full 5-Phase Methodology)
┌──────────────────────────────────────────────────────────────┐
│ Platform -> Assess -> Specify -> Plan -> Deliver             │
│ $ agentic-agent task list                                    │
│ $ agentic-agent task claim <id>                              │
│ $ agentic-agent validate                                     │
│                                                              │
│ ├─ Use for:      Critical features, regulatory compliance   │
│ ├─ Team size:    3+ people, multiple teams                 │
│ ├─ Duration:     3-4 weeks                                 │
│ ├─ Spec file:    Full SDD spec (10+ pages)                 │
│ ├─ Gates:        5 mandatory gates (Context→Domain→...)    │
│ └─ Setup time:   1+ hours                                  │
└──────────────────────────────────────────────────────────────┘
```

---

## 5. Single Package Workflow

```
┌──────────────────────────────────────────────────────────────┐
│         SINGLE PACKAGE: OpenSpec Workflow                    │
└──────────────────────────────────────────────────────────────┘

myproject/
├── agnostic-agent.yaml
├── requirements.md          ← Write spec
├── src/
└── .agentic/
    └── openspec/changes/    ← Tracks lifecycle

WORKFLOW:

  Day 1: PLAN
  ┌──────────────────────────────────────┐
  │ cat > requirements.md                │
  │ # Feature Name                       │
  │ ## Problem: ...                      │
  │ ## Solution: ...                     │
  │ ## Acceptance Criteria: ...          │
  └──────────────────────────────────────┘
            ↓
  ┌──────────────────────────────────────┐
  │ agentic-agent openspec init          │
  │   "Feature" --from requirements.md   │
  └──────────────────────────────────────┘
            ↓
        Creates:
  .agentic/openspec/changes/feature/
    ├── proposal.md       ← Auto-generated summary
    └── tasks.md          ← Auto-decomposed tasks

  Days 2-N: BUILD
  ┌──────────────────────────────────────┐
  │ For each task:                       │
  │                                      │
  │ agentic-agent task claim <id>        │
  │ [Code in src/]                       │
  │ agentic-agent task complete <id>     │
  └──────────────────────────────────────┘
            ↓
  Final Day: VERIFY & RELEASE
  ┌──────────────────────────────────────┐
  │ agentic-agent openspec status        │
  │   Check: all criteria met? ✓         │
  │   Check: tests pass? ✓               │
  │   Check: code reviewed? ✓            │
  │                                      │
  │ agentic-agent openspec complete      │
  │ agentic-agent openspec archive       │
  └──────────────────────────────────────┘
            ↓
        ✅ DONE
```

---

## 6. Monorepo Workflow (Multiple Packages)

```
┌──────────────────────────────────────────────────────────────┐
│       MONOREPO: Multi-Package Parallel Workflow              │
└──────────────────────────────────────────────────────────────┘

monorepo/
├── agnostic-agent.yaml      ← Root config (multi-path specDirs)
├── specs/                   ← Shared specs
├── packages/
│   ├── service-a/
│   │   ├── specs/
│   │   │   ├── auth-oauth.md        ← Package-specific spec
│   │   │   └── logging-upgrade.md
│   │   └── src/
│   ├── service-b/
│   │   ├── specs/
│   │   │   ├── webhook-retry.md
│   │   │   └── db-migration.md
│   │   └── src/
│   └── service-c/
│       ├── specs/
│       │   └── api-v2.md
│       └── src/
└── .agentic/
    └── openspec/changes/    ← Centralized tracking
        ├── auth-oauth/
        ├── webhook-retry/
        └── api-v2/

WORKFLOW (PARALLEL):

  Day 1: PLAN ALL FEATURES
  ┌──────────────────────────────────────┐
  │ agentic-agent openspec init "Auth"   │
  │   --from packages/service-a/specs/.. │
  │                                      │
  │ agentic-agent openspec init "Webhook"│
  │   --from packages/service-b/specs/.. │
  │                                      │
  │ agentic-agent openspec init "API"    │
  │   --from packages/service-c/specs/.. │
  └──────────────────────────────────────┘
            ↓
  Days 2-5: BUILD IN PARALLEL
  ┌────────────────────────────────────────────────────┐
  │ Team A (service-a):                               │
  │   cd packages/service-a                           │
  │   agentic-agent task claim <task-1>              │
  │   [code] → task complete                          │
  │                                                   │
  │ Team B (service-b):                    PARALLEL   │
  │   cd packages/service-b                           │
  │   agentic-agent task claim <task-2>              │
  │   [code] → task complete                          │
  │                                                   │
  │ Team C (service-c):                               │
  │   cd packages/service-c                           │
  │   agentic-agent task claim <task-3>              │
  │   [code] → task complete                          │
  └────────────────────────────────────────────────────┘
            ↓
  Day 6: INTEGRATE & TEST
  ┌──────────────────────────────────────┐
  │ Service A provides OAuth tokens      │
  │   ↓                                  │
  │ Service B & C consume tokens         │
  │   ↓                                  │
  │ Run integration tests                │
  └──────────────────────────────────────┘
            ↓
  Day 7: VERIFY & RELEASE
  ┌──────────────────────────────────────┐
  │ agentic-agent openspec complete auth │
  │ agentic-agent openspec complete ..   │
  │ agentic-agent openspec archive auth  │
  └──────────────────────────────────────┘
            ↓
        ✅ ALL DONE
```

---

## 7. Spec to Code Flow

```
┌──────────────────────────────────────────────────────────────┐
│          FROM SPEC TO IMPLEMENTED CODE                       │
└──────────────────────────────────────────────────────────────┘

SPEC FILE (requirements.md):
┌────────────────────────────────────┐
│ # OAuth Integration                │
│                                    │
│ ## Problem                         │
│ Users can't login with OAuth       │
│                                    │
│ ## Solution                        │
│ Add OAuth 2.0 flow (Google/GitHub) │
│                                    │
│ ## Acceptance Criteria             │
│ - [ ] Users can login with Google  │
│ - [ ] Users can login with GitHub  │
│ - [ ] Tokens refresh automatically │
│ - [ ] Tests pass                   │
└────────────────────────────────────┘
            ↓
  agentic-agent openspec init "OAuth" --from requirements.md
            ↓

AUTO-GENERATED PROPOSAL.MD:
┌────────────────────────────────────┐
│ # OAuth Integration                │
│                                    │
│ OAuth 2.0 support for Google and   │
│ GitHub to enable passwordless auth │
│ with automatic token refresh       │
└────────────────────────────────────┘
            ↓

AUTO-GENERATED TASKS.MD:
┌────────────────────────────────────┐
│ ## Tasks                           │
│                                    │
│ 1. Design OAuth provider interface │
│    - Files: src/oauth/provider.ts  │
│    - Acceptance: tests pass        │
│                                    │
│ 2. Implement Google OAuth          │
│    - Files: src/oauth/google.ts    │
│    - Acceptance: login test passes │
│                                    │
│ 3. Implement GitHub OAuth          │
│    - Files: src/oauth/github.ts    │
│    - Acceptance: login test passes │
│                                    │
│ 4. Add token refresh logic         │
│    - Files: src/oauth/refresh.ts   │
│    - Acceptance: refresh works     │
│                                    │
│ 5. Update UI with OAuth buttons    │
│    - Files: src/pages/login.tsx    │
│    - Acceptance: buttons appear    │
│                                    │
│ 6. Write integration tests         │
│    - Files: tests/oauth.test.ts    │
│    - Acceptance: all tests pass    │
└────────────────────────────────────┘
            ↓
        agentic-agent task claim task-1
            ↓
  FOR EACH TASK (Prompt Guide: "Package-Specific Implementation")
  ┌────────────────────────────────────┐
  │ Using Prompt #2 (Package-Specific) │
  │                                    │
  │ Claude suggests:                   │
  │ ├─ Interface design                │
  │ ├─ Code skeleton                   │
  │ ├─ Files to create                 │
  │ └─ Test structure                  │
  │                                    │
  │ Developer:                         │
  │ ├─ Writes code based on prompt     │
  │ ├─ Runs tests: npm test            │
  │ └─ Commits changes                 │
  └────────────────────────────────────┘
            ↓
        agentic-agent task complete task-1
            ↓
        ... repeat for all tasks ...
            ↓
        agentic-agent openspec complete oauth
            ↓
    ✅ ACCEPTANCE CRITERIA MET
```

---

## 8. Gate Check Flow (SDD Only)

```
┌──────────────────────────────────────────────────────────────┐
│       SDD GATE CHECKS (High-Risk Features)                   │
└──────────────────────────────────────────────────────────────┘

Platform -> Assess -> Specify -> Plan -> Deliver
            ↓
GATE 1: CONTEXT GATE
┌───────────────────────────────────┐
│ ✓ Spec is in code repo            │
│ ✓ Team has read spec              │
│ ✓ Dependencies identified         │
│ ✓ Tech stack chosen               │
└───────────────────────────────────┘
      ✅ PASS → Continue
      ❌ FAIL → Fix and re-check
            ↓
GATE 2: DOMAIN GATE
┌───────────────────────────────────┐
│ ✓ Acceptance criteria clear       │
│ ✓ Edge cases documented           │
│ ✓ Error handling designed         │
│ ✓ Contracts between services OK   │
└───────────────────────────────────┘
      ✅ PASS → Continue
      ❌ FAIL → Update spec
            ↓
GATE 3: INTEGRATION GATE
┌───────────────────────────────────┐
│ ✓ External APIs documented        │
│ ✓ Mock strategy defined           │
│ ✓ Test data prepared              │
│ ✓ Observability points chosen     │
└───────────────────────────────────┘
      ✅ PASS → Continue
      ❌ FAIL → Add design details
            ↓
GATE 4: NFR GATE (Non-Functional)
┌───────────────────────────────────┐
│ ✓ Performance targets set         │
│ ✓ Security requirements clear     │
│ ✓ Scalability approach defined    │
│ ✓ Compliance checklist done       │
└───────────────────────────────────┘
      ✅ PASS → Continue
      ❌ FAIL → Add NFR details
            ↓
GATE 5: READY GATE
┌───────────────────────────────────┐
│ ✓ All tasks broken down           │
│ ✓ Dependencies resolved           │
│ ✓ Team assigned                   │
│ ✓ Timeline realistic              │
└───────────────────────────────────┘
      ✅ PASS → Implementation starts
      ❌ FAIL → Refine plan
            ↓
  agentic-agent validate
        (checks all 5 gates)
            ↓
    Either: ✅ ALL PASS (proceed)
    Or:     ❌ FAILURES (fix then retry)
```

---

## 9. Implementation Prompt Flow

```
┌──────────────────────────────────────────────────────────────┐
│       7 IMPLEMENTATION PROMPTS: When to Use Each             │
└──────────────────────────────────────────────────────────────┘

PROJECT TIMELINE:
┌─────────────────────────────────────────────────────────────┐
│                                                             │
│  DAY 1: PLANNING                                            │
│  ┌─────────────────────────────────────────────────────┐   │
│  │ Prompt #1: Spec Review                              │   │
│  │ "Break spec into concrete tasks"                    │   │
│  │ → Generates task checklist                          │   │
│  └─────────────────────────────────────────────────────┘   │
│                                                             │
│  DAYS 2-N: IMPLEMENTATION                                  │
│  ┌─────────────────────────────────────────────────────┐   │
│  │ For EACH task:                                      │   │
│  │                                                     │   │
│  │ Prompt #2: Package-Specific                        │   │
│  │ "Design + code example for this task in my package"│   │
│  │ → Generates design + code skeleton                 │   │
│  │                                                     │   │
│  │ Prompt #3: Integration (if multi-service)          │   │
│  │ "How does my service integrate with others?"       │   │
│  │ → Generates contract + integration points          │   │
│  │                                                     │   │
│  │ Prompt #4: Testing                                 │   │
│  │ "Structure tests for these acceptance criteria"    │   │
│  │ → Generates test outline + mocks                   │   │
│  │                                                     │   │
│  │ [Developer implements]                             │   │
│  │ [Runs tests]                                       │   │
│  │ [Commits changes]                                  │   │
│  └─────────────────────────────────────────────────────┘   │
│                                                             │
│  DEPENDENCY PHASE (if multi-service):                      │
│  ┌─────────────────────────────────────────────────────┐   │
│  │ Prompt #6: Dependency Analysis                      │   │
│  │ "What's the implementation order? Rollout strategy?"│   │
│  │ → Generates dependency graph + sequence            │   │
│  └─────────────────────────────────────────────────────┘   │
│                                                             │
│  VERIFICATION PHASE:                                       │
│  ┌─────────────────────────────────────────────────────┐   │
│  │ Prompt #5: Verification                             │   │
│  │ "Is everything complete? All criteria met?"         │   │
│  │ → Generates checklist + blocking issues            │   │
│  │                                                     │   │
│  │ Prompt #7: Documentation                           │   │
│  │ "Generate docs for this feature"                   │   │
│  │ → Generates API docs + guides                      │   │
│  └─────────────────────────────────────────────────────┘   │
│                                                             │
└─────────────────────────────────────────────────────────────┘

WHEN STUCK:
┌─────────────────────────────────────────────────────────────┐
│ Design unclear?        → Use Prompt #2 (Package-Specific)   │
│ Cross-service problem? → Use Prompt #3 (Integration)        │
│ Testing strategy?      → Use Prompt #4 (Testing)            │
│ Implementation order?  → Use Prompt #6 (Dependency)         │
│ "Is it done?"          → Use Prompt #5 (Verification)       │
│ "How do I document?"   → Use Prompt #7 (Documentation)      │
└─────────────────────────────────────────────────────────────┘
```

---

## 10. Multi-Service Dependency Graph

```
┌──────────────────────────────────────────────────────────────┐
│       EXAMPLE: OAuth Feature Across 3 Services              │
└──────────────────────────────────────────────────────────────┘

INITIAL SPECS (in each package):
┌──────────────┐      ┌──────────────┐      ┌──────────────┐
│   Auth Svc   │      │   API Svc    │      │   Web UI     │
│              │      │              │      │              │
│ OAuth        │      │ Token        │      │ Login Form   │
│ Provider     │      │ Verification │      │ OAuth Buttons│
│              │      │              │      │              │
│ Provides:    │      │ Needs:       │      │ Needs:       │
│ - POST /auth │      │ - OAuth      │      │ - OAuth      │
│ - tokens     │      │   tokens     │      │   endpoints  │
└──────────────┘      └──────────────┘      └──────────────┘

DEPENDENCIES (Prompt #6 analysis):
┌──────────────────────────────────────────────────────────────┐
│                                                              │
│  Auth Svc (START HERE)                                      │
│  ├─ Implement OAuth endpoints                               │
│  ├─ Add token generation                                    │
│  ├─ Tests pass ✓                                            │
│  └─ Deploy to staging                                       │
│           ↓                                                  │
│  API Svc (WAIT FOR AUTH)                                    │
│  ├─ Add token verification middleware                       │
│  ├─ Call Auth Svc /verify endpoint                         │
│  ├─ Tests pass ✓                                            │
│  └─ Deploy to staging                                       │
│           ↓                                                  │
│  Web UI (LAST)                                              │
│  ├─ Add login form with OAuth buttons                       │
│  ├─ Call Auth Svc endpoints                                │
│  ├─ Store tokens in localStorage                           │
│  ├─ Tests pass ✓                                            │
│  └─ Deploy                                                  │
│                                                              │
└──────────────────────────────────────────────────────────────┘

BLOCKERS & RISKS:
┌──────────────────────────────────────────────────────────────┐
│ ❌ If Auth Svc not ready                                    │
│    → API Svc can mock /verify endpoint                      │
│                                                              │
│ ❌ If API Svc not ready                                     │
│    → Web UI can call Auth Svc directly (temporary)          │
│                                                              │
│ ✅ Contract first                                           │
│    → Agree API contracts before implementation              │
│    → Mock endpoints until real service ready                │
└──────────────────────────────────────────────────────────────┘
```

---

## 11. Time Estimates by Workflow

```
┌──────────────────────────────────────────────────────────────┐
│           REALISTIC TIME BREAKDOWN                           │
└──────────────────────────────────────────────────────────────┘

SINGLE PACKAGE (OPENSPEC):

Feature: "Product Search" (1-2 week feature)

  ├─ Setup & Spec             1 hour
  │  └─ Write requirements.md + openspec init
  │
  ├─ Implementation            4-5 days
  │  ├─ Task 1 (Backend API)       1 day
  │  ├─ Task 2 (Frontend UI)       1 day
  │  ├─ Task 3 (Filters)           1 day
  │  └─ Task 4 (Tests + Polish)    1-2 days
  │
  └─ Verification + Release    1 day
     └─ Final tests + code review + archive

  TOTAL: ~1.5 weeks


MONOREPO (MULTI-PACKAGE):

Feature: "OAuth Integration" (2-3 week feature, 3 teams)

  ├─ Planning                  2 hours
  │  ├─ Spec for Auth Svc         30 min
  │  ├─ Spec for API Svc          30 min
  │  └─ Spec for Web UI           30 min
  │
  ├─ Implementation            6-8 days (PARALLEL)
  │  ├─ Team A (Auth Svc)          2-3 days
  │  ├─ Team B (API Svc)           2-3 days  PARALLEL
  │  └─ Team C (Web UI)            2-3 days
  │
  ├─ Integration & Testing    2-3 days
  │  └─ Cross-service tests + deployment
  │
  └─ Verification + Release    1 day

  TOTAL: ~2-3 weeks (would be 4-5 weeks if sequential)


SDD (CRITICAL FEATURE):

Feature: "Payment Processing" (high-risk, large team)

  ├─ Planning & Specification  5-7 days
  │  ├─ Analyst phase             2 days
  │  ├─ Architect phase           2 days
  │  └─ Gate checks               1-2 days
  │
  ├─ Implementation            8-10 days
  │  ├─ Core payment logic         3-4 days
  │  ├─ Error handling + edge cases 2-3 days
  │  ├─ Tests                      2-3 days
  │  └─ Security audit            1-2 days
  │
  ├─ Verification             3-5 days
  │  ├─ Verifier phase            2-3 days
  │  └─ Final gate checks          1-2 days
  │
  └─ Release                   1 day

  TOTAL: ~3-4 weeks


GUIDELINE:
┌──────────────────────────────────────────────────────────────┐
│ • Add 20-30% buffer for unknowns                            │
│ • Parallel work reduces total time (monorepo advantage)      │
│ • High-risk features need more planning time upfront        │
│ • Each gate check adds 1-2 days                             │
│ • Testing often takes 30-40% of implementation time         │
└──────────────────────────────────────────────────────────────┘
```

---

## 12. Config Comparison

```
┌──────────────────────────────────────────────────────────────┐
│        AGNOSTIC-AGENT.YAML: Single vs Monorepo              │
└──────────────────────────────────────────────────────────────┘

SINGLE PACKAGE:
┌────────────────────────────────────────┐
│ project:                               │
│   name: MyProject                      │
│   version: 0.1.0                       │
│                                        │
│ agents:                                │
│   defaults:                            │
│     model: claude-3-5-sonnet-20241022 │
│                                        │
│ paths:                                 │
│   openSpecDir: .agentic/openspec/..   │
└────────────────────────────────────────┘
                ↓
        (specs in .agentic/spec/ or root)


MONOREPO:
┌────────────────────────────────────────┐
│ project:                               │
│   name: MyMonorepo                     │
│   version: 0.1.0                       │
│                                        │
│ agents:                                │
│   defaults:                            │
│     model: claude-3-5-sonnet-20241022 │
│                                        │
│ paths:                                 │
│   openSpecDir: .agentic/openspec/..   │
│   specDirs:                            │
│     - packages/*/specs     ← KEY!      │
│     - apps/*/specs                     │
│     - specs                            │
└────────────────────────────────────────┘
                ↓
  (specs auto-discovered in multiple places)
  (first match wins in order above)
```

---

## Summary

Use these diagrams to:
- **Choose workflow** → Diagram #1
- **Track progress** → Diagrams #2, #5, #6
- **Monorepo setup** → Diagrams #3, #4, #6, #10
- **Understand gates** → Diagram #8
- **When to use prompts** → Diagram #9
- **Estimate time** → Diagram #11
- **Compare configs** → Diagram #12
