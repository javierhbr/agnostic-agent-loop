# OpenSpec in Monorepos

**Working with OpenSpec as the primary workflow tool in a monorepo structure.**

This guide shows how to organize specs per-package, track changes centrally, and use prompts to guide implementation—all with the `agentic-agent` CLI.

---

## TL;DR

1. **One root config** (agnostic-agent.yaml) with multi-path spec resolution
2. **Specs live per-package** (packages/*/specs/)
3. **OpenSpec tracks changes** (.agentic/openspec/changes/)
4. **CLI auto-discovers** all specs across packages
5. **Prompts guide implementation** for each package/task

**Core boundary rule:** Component repos (individual packages) use OpenSpec ONLY. Platform-level tooling (BMAD, Speckit) stays in the platform repo. Use `platform-ref.yaml` in component repos to reference platform alignment without importing platform tools.

---

## Setup (One-Time)

### 1. Create Root Config

In your monorepo root, create or update `agnostic-agent.yaml`:

```yaml
project:
  name: MyMonorepo
  version: 0.1.0
  roots:
    - .

agents:
  defaults:
    max_tokens: 4000
    model: claude-3-5-sonnet-20241022

# Multi-directory spec resolution:
# Specs are searched in order — first match wins
paths:
  openSpecDir: .agentic/openspec/changes
  specDirs:
    - packages/*/specs           # Package-level specs
    - apps/*/specs               # App-level specs (if applicable)
    - specs                      # Root-level specs (shared)
  contextDirs:
    - .agentic/context

workflow:
  validators:
    - context-check
    - task-scope
```

### 2. Create Package Spec Directories

```bash
# Service A
mkdir -p packages/service-a/specs

# Service B
mkdir -p packages/service-b/specs

# Shared specs (optional)
mkdir -p specs
```

---

## Creating Specs Per Package

### Structure

Each package has its own `specs/` directory with one spec per feature/change:

```
packages/
├── service-a/
│   ├── specs/
│   │   ├── auth-oauth.md
│   │   ├── logging-upgrade.md
│   │   └── cache-layer.md
│   ├── src/
│   └── package.json
├── service-b/
│   ├── specs/
│   │   ├── webhook-retry.md
│   │   └── db-migration.md
│   └── src/
└── service-c/
    ├── specs/
    │   └── api-v2.md
    └── src/
```

### Example Spec File

**packages/service-a/specs/auth-oauth.md**

```markdown
# OAuth Integration

## Problem
Users can only authenticate with username/password. We need OAuth support (Google, GitHub).

## Solution
Integrate OAuth 2.0 flow using industry-standard libraries.

## What Changes
- Add OAuth provider abstraction
- Create /auth/callback endpoint
- Store provider tokens in database
- Update login UI with provider buttons

## Acceptance Criteria
- [ ] Users can login with Google
- [ ] Users can login with GitHub
- [ ] Tokens refresh automatically
- [ ] Existing password auth still works
- [ ] Tests cover both flows

## Notes
- Avoid vendor lock-in to single OAuth library
- Keep token storage encrypted
- Set refresh token expiry to 30 days
```

---

## Initialize OpenSpec from Specs

When you're ready to track a spec as a "change", initialize it:

```bash
# CLI searches all specDirs for the file and reads it
agentic-agent openspec init "OAuth Integration" \
  --from packages/service-a/specs/auth-oauth.md

# Creates: .agentic/openspec/changes/oauth-integration/
#  ├── proposal.md          (extracted from spec)
#  └── tasks.md             (decomposed into actionable tasks)
```

The CLI:
1. Finds `auth-oauth.md` in `packages/service-a/specs/`
2. Reads its content
3. Generates `proposal.md` (problem + solution summary)
4. Generates `tasks.md` (breakdown of acceptance criteria into tasks)

---

## Working with Tasks

### List All Changes Across Monorepo

```bash
# Shows all OpenSpec changes discovered
agentic-agent openspec list

# Output:
# oauth-integration       In Progress
# webhook-retry           Pending
# cache-layer             Completed
```

### Show Spec Details

```bash
agentic-agent openspec show oauth-integration

# Output includes:
# - Original spec location (packages/service-a/specs/auth-oauth.md)
# - Proposal summary
# - Tasks breakdown
# - Change status
```

### Create Tasks from OpenSpec

After initializing OpenSpec, create tasks:

```bash
# Option A: Let CLI decompose tasks automatically
agentic-agent task create "Implement OAuth provider abstraction"
agentic-agent task create "Build /auth/callback endpoint"
agentic-agent task create "Add token refresh logic"
agentic-agent task create "Update login UI"
agentic-agent task create "Write integration tests"

# Option B: Track via OpenSpec
agentic-agent task list  # See all tasks
agentic-agent task claim <task-id>
# ... implement ...
agentic-agent task complete <task-id>
```

### Work on a Package

```bash
# Move into the package directory
cd packages/service-a

# Claim a task
agentic-agent task claim <task-id>

# Implement changes in this package
# - Edit src/auth/*.ts
# - Add tests in tests/auth/*.test.ts

# Complete the task
agentic-agent task complete <task-id>

# Return to repo root
cd ../..
```

---

## Monorepo Directory Structure

Complete example:

```
mymonorepo/
├── agnostic-agent.yaml              # Root config (multi-path spec resolution)
├── README.md
├── turbo.json                        # (if using Turbo)
├── package.json
│
├── specs/                            # Shared/cross-package specs
│   ├── platform-strategy.md
│   └── security-requirements.md
│
├── packages/
│   ├── service-a/                    # Microservice
│   │   ├── specs/
│   │   │   ├── auth-oauth.md
│   │   │   ├── logging-upgrade.md
│   │   │   └── cache-layer.md
│   │   ├── src/
│   │   │   ├── auth/
│   │   │   ├── logging/
│   │   │   └── cache/
│   │   ├── tests/
│   │   ├── package.json
│   │   └── tsconfig.json
│   │
│   ├── service-b/                    # Microservice
│   │   ├── specs/
│   │   │   ├── webhook-retry.md
│   │   │   └── db-migration.md
│   │   ├── src/
│   │   ├── tests/
│   │   └── package.json
│   │
│   └── service-c/                    # Microservice
│       ├── specs/
│       │   └── api-v2.md
│       ├── src/
│       └── package.json
│
├── apps/                             # (Optional) Applications
│   └── web-ui/
│       ├── specs/
│       │   └── dark-mode.md
│       └── src/
│
└── .agentic/                         # CLI-managed
    ├── openspec/
    │   └── changes/                  # OpenSpec change tracking (lifecycle)
    │       ├── oauth-integration/
    │       │   ├── proposal.md
    │       │   └── tasks.md
    │       ├── webhook-retry/
    │       │   ├── proposal.md
    │       │   └── tasks.md
    │       └── cache-layer/
    │           ├── proposal.md
    │           └── tasks.md
    ├── context/                      # Context bundles
    └── tasks/
        ├── backlog.yaml
        ├── in-progress.yaml
        └── done.yaml
```

---

## Prompts for Implementation

Use these prompts with Claude Code / Cursor / your agent to guide implementation:

### 1. **Spec Review Prompt** (Before Starting)

Use this when initializing a spec:

```
Review this spec and break it into concrete implementation tasks.

Spec: packages/service-a/specs/auth-oauth.md

For each acceptance criterion, create one task:
- What needs to be implemented (code change)
- Which files will be touched
- How to verify it works (test, manual check)

Format as a checklist.
```

### 2. **Package-Specific Implementation Prompt**

Use this when working on a task in a package:

```
I'm implementing the OAuth integration in packages/service-a/.

Current task: "Implement OAuth provider abstraction"

Requirements from spec (packages/service-a/specs/auth-oauth.md):
- Add OAuth provider abstraction
- Support multiple providers (Google, GitHub)
- Keep token storage encrypted
- Avoid vendor lock-in

Please:
1. Design the provider interface/class
2. Show example implementation for Google OAuth
3. List the files I need to create/modify
4. Provide test structure

Keep it simple and production-ready.
```

### 3. **Integration Prompt** (Cross-Package)

Use this when coordinating between packages:

```
I'm integrating OAuth across the monorepo.

Current status:
- packages/service-a: OAuth provider abstraction complete ✓
- packages/service-b: Needs to consume OAuth tokens
- packages/web-ui: Needs login UI update

What needs to happen in packages/service-b to consume the OAuth tokens from service-a?
- API contract (what endpoints to call)
- Token refresh strategy
- Error handling (token expired, invalid, etc.)

Format: Brief design + code example.
```

### 4. **Testing Prompt**

Use this when writing tests for a spec:

```
I'm writing tests for the OAuth integration (packages/service-a/specs/auth-oauth.md).

Acceptance criteria to test:
- [ ] Users can login with Google
- [ ] Users can login with GitHub
- [ ] Tokens refresh automatically
- [ ] Existing password auth still works

Please provide:
1. Unit test structure (for provider abstraction)
2. Integration test structure (full OAuth flow)
3. Mock strategy for external OAuth providers
4. CI/CD considerations

Use Jest/Vitest. Keep tests maintainable.
```

### 5. **Verification Prompt** (Before Completion)

Use this before marking a task/spec complete:

```
Verify that the OAuth integration (packages/service-a/specs/auth-oauth.md) is complete.

Checklist:
- [ ] All acceptance criteria met?
- [ ] Tests passing?
- [ ] Code review ready?
- [ ] Performance acceptable?
- [ ] Security audit passed? (token storage, refresh logic)
- [ ] Documentation updated?

Show me:
1. Test output (passing tests)
2. Files changed (git diff)
3. Any blocking issues

If all clear, I'll mark the spec as complete.
```

### 6. **Dependency Analysis Prompt** (Multi-Service)

Use this when specs in different packages depend on each other:

```
I have two specs:
- packages/service-a/specs/auth-oauth.md (provides OAuth)
- packages/service-b/specs/webhook-retry.md (consumes OAuth for webhooks)

What's the implementation order?
- Which package starts first?
- What's the contract between them?
- What happens if service-b is deployed before service-a?

Show dependency graph and rollout strategy.
```

### 7. **Documentation Prompt**

Use this when creating implementation docs:

```
Generate documentation for the OAuth integration based on this spec:
- File: packages/service-a/specs/auth-oauth.md
- Target audience: developers integrating OAuth in other services

Sections:
1. Overview (what was implemented)
2. API reference (endpoints, payload examples)
3. Token lifecycle (refresh, expiry, rotation)
4. Error handling (common errors + solutions)
5. Testing (how to test OAuth flow locally)
6. Troubleshooting

Format as Markdown, include code examples.
```

---

## Workflow: From Spec to Production

### Day 1: Define Spec

```bash
# Write spec in package
cat > packages/service-a/specs/auth-oauth.md << 'EOF'
# OAuth Integration
...
EOF

# Initialize OpenSpec to start tracking
agentic-agent openspec init "OAuth Integration" \
  --from packages/service-a/specs/auth-oauth.md
```

**Use Prompt #1 (Spec Review)** to decompose into tasks.

### Days 2-5: Implement

```bash
# For each task:
agentic-agent task claim <task-id>

# In package directory, implement changes
cd packages/service-a
# ... edit src/auth/*.ts, add tests ...

# Verify with Prompt #4 (Testing)
npm test

# Complete task
cd ../..
agentic-agent task complete <task-id>
```

**Use Prompt #2 (Package-Specific)** to guide each implementation.

### Day 6: Integrate (if multi-package)

```bash
# Check if other packages need updates
agentic-agent openspec list

# For dependent specs (e.g., service-b consuming OAuth):
cd packages/service-b
# ... implement integration ...

# Verify integration with Prompt #3 (Integration)
```

### Day 7: Verify & Complete

```bash
# Run full verification
agentic-agent openspec show oauth-integration

# Use Prompt #5 (Verification) to ensure all criteria met

# Mark as complete
agentic-agent openspec complete oauth-integration

# Archive for reference
agentic-agent openspec archive oauth-integration
```

---

## CLI Commands Reference

| Command | Purpose | Example |
|---------|---------|---------|
| `openspec init <name> --from <file>` | Create change from spec | `agentic-agent openspec init "OAuth" --from packages/service-a/specs/auth-oauth.md` |
| `openspec list` | List all changes across monorepo | `agentic-agent openspec list` |
| `openspec show <id>` | Show spec details & original location | `agentic-agent openspec show oauth-integration` |
| `openspec status <id>` | Show change progress | `agentic-agent openspec status oauth-integration` |
| `openspec complete <id>` | Mark change as done | `agentic-agent openspec complete oauth-integration` |
| `openspec archive <id>` | Archive completed change | `agentic-agent openspec archive oauth-integration` |
| `task create <name>` | Create a task | `agentic-agent task create "Implement OAuth provider"` |
| `task claim <id>` | Start working on task | `agentic-agent task claim task-123` |
| `task complete <id>` | Mark task done | `agentic-agent task complete task-123` |
| `task list` | List all tasks | `agentic-agent task list` |

---

## Multi-Service Example: Payment Flow

Real monorepo example:

```
packages/
├── billing-service/
│   └── specs/
│       └── subscription-billing.md
├── payments-service/
│   └── specs/
│       └── payment-processing.md
└── customer-api/
    └── specs/
        └── billing-endpoints.md
```

### Workflow

```bash
# Day 1: Define specs for all three services
agentic-agent openspec init "Subscription Billing" \
  --from packages/billing-service/specs/subscription-billing.md

agentic-agent openspec init "Payment Processing" \
  --from packages/payments-service/specs/payment-processing.md

agentic-agent openspec init "Billing API Endpoints" \
  --from packages/customer-api/specs/billing-endpoints.md

# Days 2-5: Implement in parallel
# Team A works on billing-service
cd packages/billing-service && npm run dev

# Team B works on payments-service
cd packages/payments-service && npm run dev

# Team C works on customer-api
cd packages/customer-api && npm run dev

# Day 6: Integration testing
# Use Prompt #3 (Integration) to verify contracts between services

# Day 7: Verify all three
agentic-agent openspec status subscription-billing
agentic-agent openspec status payment-processing
agentic-agent openspec status billing-endpoints

# Mark complete
agentic-agent openspec complete subscription-billing
agentic-agent openspec complete payment-processing
agentic-agent openspec complete billing-endpoints
```

---

## FAQ: Monorepo + OpenSpec

**Q: Can I have specs at root level too?**

A: Yes! The config includes `- specs` in `specDirs`, so root-level specs are found automatically. Use for cross-package specs.

```bash
# Root spec (applies to multiple packages)
cat > specs/shared-security.md << 'EOF'
# Security Standards
All packages must:
- Validate input
- Use HTTPS only
- Encrypt sensitive data
EOF

agentic-agent openspec init "Security Standards" \
  --from specs/shared-security.md
```

**Q: What if two packages have specs with the same name?**

A: The first match wins (based on `specDirs` order). To avoid collisions, use unique names or namespace:

```
packages/service-a/specs/auth-oauth.md      (found first)
packages/service-b/specs/cache-oauth.md     (different name)
```

**Q: Can I organize specs differently?**

A: Yes! Update `specDirs` in agnostic-agent.yaml:

```yaml
paths:
  specDirs:
    - docs/specs                     # Centralized specs
    - packages/*/specs               # Package-level overrides
```

**Q: Do I need one change per spec?**

A: No. Create OpenSpec changes only when you want to track lifecycle. Small fixes can use `agentic-agent task` directly.

**Q: How do I handle cross-package dependencies?**

A: Use Prompt #6 (Dependency Analysis) to determine order. Then:
- Service A starts first (provider)
- Service B waits for A's API to stabilize
- Service C integrates both

**Q: What if requirements change during implementation?**

A: Update the spec file directly:
```bash
vim packages/service-a/specs/auth-oauth.md
```

Then regenerate tasks:
```bash
agentic-agent openspec init "OAuth Integration" \
  --from packages/service-a/specs/auth-oauth.md
```

---

## Summary

| Aspect | What to Do |
|--------|-----------|
| **Config** | One root `agnostic-agent.yaml` with multi-path specDirs |
| **Specs** | Write in `packages/*/specs/` (live with code) |
| **Tracking** | Use `agentic-agent openspec` for lifecycle |
| **Tasks** | Create from specs, claim, complete, verify |
| **Prompts** | Use 7 prompts to guide implementation |
| **Verification** | Check acceptance criteria before marking complete |
| **Parallel Work** | Multiple teams/packages work independently |

**Result:** Specs stay with packages, OpenSpec tracks changes centrally, CLI discovers everything automatically.

---

## Next Steps

1. Create root `agnostic-agent.yaml` with your monorepo structure
2. Write specs in `packages/*/specs/` directories
3. Initialize OpenSpec from specs: `agentic-agent openspec init`
4. Use the 7 prompts to guide implementation
5. Track progress with `agentic-agent task` and `openspec`
6. Archive completed changes: `agentic-agent openspec archive`

**Don't over-engineer:** Start simple, add specs only when needed (> 1 week of work).
