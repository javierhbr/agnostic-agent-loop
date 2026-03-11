# Agentic Helper — Workflow Commands

Full command sequences for each workflow tier (TINY, SMALL, OpenSpec+, Full SDD) with examples. Referenced by the main agentic-helper skill.

---

## TINY Workflow (Bug Fix, Single Task)

**When to use:** Bug fix, < 1 day, 1 file, low risk.

**Setup:** 5 minutes.

### Full Command Sequence

```bash
# 1. Create the task
agentic-agent task create --title "Fix: [description]"

# Output shows TASK-ID, copy it

# 2. Claim the task (records git branch + timestamp)
agentic-agent task claim TASK-ID

# 3. Generate context for your directory
agentic-agent context generate <DIR>

# 4. [... implement the fix ...]

# 5. Validate before completing
agentic-agent validate

# 6. Complete the task (captures commits)
agentic-agent task complete TASK-ID
```

### Example: Fix a Typo

```bash
agentic-agent task create --title "Fix: typo in README"
agentic-agent task claim TASK-123
agentic-agent context generate docs/
# Edit README.md
agentic-agent validate
agentic-agent task complete TASK-123
```

---

## SMALL Workflow (Feature, OpenSpec)

**When to use:** Feature, 1-2 weeks, 1 service, medium risk.

**Setup:** 10 minutes.

### Full Command Sequence

```bash
# 1. Ensure skills are set up
agentic-agent skills ensure

# 2. Initialize the spec from requirements file (PRD or doc)
agentic-agent openspec init "<feature-name>" --from <requirements-file>

# Output shows change-id, copy it

# 3. List all tasks in the spec
agentic-agent task list

# 4. For each task (loop):
agentic-agent task claim <TASK-ID>
agentic-agent context generate <DIR>
# ... implement ...
agentic-agent validate
agentic-agent task complete <TASK-ID>

# 5. After all tasks complete:
agentic-agent openspec complete <change-id>

# 6. Archive the change
agentic-agent openspec archive <change-id>
```

### Example: Add Dark Mode Feature

```bash
agentic-agent skills ensure
agentic-agent openspec init "dark-mode" --from docs/prd-dark-mode.md
agentic-agent task list
# Shows TASK-1: Add theme context, TASK-2: Create theme CSS, TASK-3: Add toggle UI

agentic-agent task claim TASK-1
agentic-agent context generate src/
# ... implement context ...
agentic-agent validate
agentic-agent task complete TASK-1

# ... repeat for TASK-2, TASK-3 ...

agentic-agent openspec complete dark-mode-change-001
agentic-agent openspec archive dark-mode-change-001
```

---

## OpenSpec+ Workflow (Monorepo, Multiple Packages)

**When to use:** Multi-package monorepo, 2-4 weeks, cross-service feature, medium-high risk.

**Setup:** 20 minutes.

### Full Command Sequence

```bash
# 1. Ensure skills
agentic-agent skills ensure

# 2. Initialize per-package specs
agentic-agent openspec init "<package-1-feature>" --from <prd-file>
agentic-agent openspec init "<package-2-feature>" --from <prd-file>
agentic-agent openspec init "<package-3-feature>" --from <prd-file>

# 3. List all tasks across all specs
agentic-agent task list

# 4. For each task (can work in parallel):
agentic-agent task claim <TASK-ID>
agentic-agent context build --task <TASK-ID>
# ... implement per package ...
agentic-agent validate
agentic-agent task complete <TASK-ID>

# 5. Run integration tests across all packages
# (custom test command, e.g., npm run test:integration)

# 6. Complete all specs
agentic-agent openspec complete <change-id-1>
agentic-agent openspec complete <change-id-2>
agentic-agent openspec complete <change-id-3>

# 7. Archive
agentic-agent openspec archive <change-id-1>
agentic-agent openspec archive <change-id-2>
agentic-agent openspec archive <change-id-3>
```

### Example: Payment Integration Across 3 Services

```bash
agentic-agent skills ensure

# Create per-service specs
agentic-agent openspec init "payment-checkout" --from docs/prd-payments.md
agentic-agent openspec init "payment-fulfillment" --from docs/prd-payments.md
agentic-agent openspec init "payment-accounting" --from docs/prd-payments.md

agentic-agent task list
# Shows 12 tasks across 3 specs

# Work on tasks in parallel
agentic-agent task claim TASK-1  # checkout service
agentic-agent context build --task TASK-1
# ... implement checkout payment ...
agentic-agent validate
agentic-agent task complete TASK-1

# ... repeat for other tasks ...

agentic-agent openspec complete payment-checkout-change-001
agentic-agent openspec complete payment-fulfillment-change-001
agentic-agent openspec complete payment-accounting-change-001
```

---

## Full SDD Workflow (Critical Risk)

**When to use:** Payment, auth, PII, contract breaking, 4+ services, high/critical risk.

**Setup:** 1+ hour.

### Full Command Sequence

```bash
# Phase 0: Initiative + Risk Assessment
agentic-agent specify start "<name>" --risk critical
agentic-agent specify workflow show <id>

# Phase 1: Architecture + Specs
agentic-agent specify gate-check <spec-id>

# If blocked by ADR:
agentic-agent specify adr list --blocked
agentic-agent specify adr create --title "<decision>"
agentic-agent specify adr resolve <ADR-ID>

# Phase 2: Development (parallel per component)
agentic-agent task claim <ID>
agentic-agent context build --task <ID>
# ... code, observability, edge cases, tests ...
agentic-agent specify gate-check <component-spec-id>
agentic-agent task complete <ID>

# Phase 3: Verification
agentic-agent validate
agentic-agent specify gate-check <spec-id>
agentic-agent specify sync-graph

# Phase 4: Deploy (feature flags, progressive rollout)
# ... deployment strategy ...
```

### The 5 Gates (Validation Checkpoints)

1. **Context Completeness** — Every spec section has a Source: line
2. **Domain Validity** — No invariant violations, no cross-domain DB access
3. **Integration Safety** — All contract consumers identified
4. **NFR Compliance** — Logging, metrics, tracing, PII, performance declared
5. **Ready-to-Implement** — No ambiguity, no blocking ADRs, all ACs testable

**Run gate checks frequently:**
```bash
agentic-agent specify gate-check <spec-id>  # Before fan-out, before implementation
agentic-agent specify gate-check <component-spec-id>  # After component implementation
```

### Example: Critical Payment Flow

```bash
agentic-agent specify start "subscription-billing" --risk critical
agentic-agent specify workflow show id-001

# Architecture phase
agentic-agent specify gate-check SPEC-PAY-001

# Hit Gate 5 failure: ADR needed on idempotency
agentic-agent specify adr list --blocked
# Shows: ADR-219 blocks SPEC-PAY-001

agentic-agent specify adr create --title "Idempotency strategy for payment retries"
agentic-agent specify adr resolve ADR-219

# Re-check gate
agentic-agent specify gate-check SPEC-PAY-001  # Now PASS

# Development phase (parallel)
agentic-agent task claim TASK-1  # Billing service
agentic-agent context build --task TASK-1
# ... implement with observability, edge cases ...
agentic-agent specify gate-check SPEC-PAY-001-BILLING  # Component-level gate check
agentic-agent task complete TASK-1

# ... repeat for other components ...

# Verification phase
agentic-agent validate
agentic-agent specify gate-check SPEC-PAY-001
agentic-agent specify sync-graph

# Deploy with feature flags
```

---

## Error Recovery

### Gate Fails on Context Completeness

**Symptom:** `Gate 1 FAIL: Problem Statement missing Source:`

**Fix:**
```bash
# 1. Add Source: line to the missing section
# Edit spec file, add: Source: Platform MCP v2.1

# 2. Re-run gate
agentic-agent specify gate-check <spec-id>
```

### Blocked By Non-Empty ADR

**Symptom:** Spec shows `blocked_by: [ADR-219]`

**Fix:**
```bash
# 1. Check ADR status
agentic-agent specify adr list --blocked

# 2. Resolve the ADR
agentic-agent specify adr resolve ADR-219

# 3. Re-check gate
agentic-agent specify gate-check <spec-id>
```

### Validation Fails

**Symptom:** `agentic-agent validate` shows scope violation or context stale

**Fix:**
```bash
# 1. Regenerate context
agentic-agent context generate <DIR>

# 2. Re-validate
agentic-agent validate
```

### "Which Workflow Should I Use?"

**Fix:** Walk the decision tree in the main agentic-helper skill. Risk classification is the first step — assess whether the task touches payment, auth, PII, breaks contracts, or involves 4+ services. Each escalator moves you up one tier.

---

## Key Command Reference

```bash
# Task management
agentic-agent task create --title "..."
agentic-agent task list
agentic-agent task claim <ID>
agentic-agent task complete <ID>

# OpenSpec (Small to OpenSpec+)
agentic-agent openspec init "<name>" --from <file>
agentic-agent openspec complete <change-id>
agentic-agent openspec archive <change-id>

# SDD (Critical risk)
agentic-agent specify start "<name>" --risk critical
agentic-agent specify workflow show <id>
agentic-agent specify gate-check <spec-id>
agentic-agent specify adr create --title "..."
agentic-agent specify adr resolve <ADR-ID>
agentic-agent specify adr list --blocked
agentic-agent specify sync-graph

# Context and validation
agentic-agent context generate <DIR>
agentic-agent context build --task <ID>
agentic-agent validate
agentic-agent status
agentic-agent skills ensure
```
