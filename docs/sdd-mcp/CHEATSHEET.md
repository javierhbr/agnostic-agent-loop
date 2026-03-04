# OpenSpec + SDD Quick Reference

**One-page cheat sheet for the agentic-agent CLI.**

---

## Choose Your Workflow

```
START HERE: What are you building?

┌────────────────────────────────────────────────────────┐
│  Bug fix or small refactor? (< 1 day)                  │
│  👉 Use: agentic-agent task only                       │
│  Time: 5 minutes setup                                 │
└────────────────────────────────────────────────────────┘

┌────────────────────────────────────────────────────────┐
│  Single feature (1-2 weeks)?                           │
│  👉 Use: agentic-agent openspec init                   │
│  Time: 10 minutes setup                                │
└────────────────────────────────────────────────────────┘

┌────────────────────────────────────────────────────────┐
│  Monorepo (multi-package)?                             │
│  👉 Use: OpenSpec + multi-path config                  │
│  Time: 20 minutes setup                                │
└────────────────────────────────────────────────────────┘

┌────────────────────────────────────────────────────────┐
│  Critical / high-risk (payments, auth)?                │
│  👉 Use: Full SDD with gates                           │
│  Time: 1 hour planning + gates                         │
└────────────────────────────────────────────────────────┘
```

---

## Quick Commands

### Tasks (Simplest)

```bash
agentic-agent task create "Your task"        # Create
agentic-agent task list                      # View all
agentic-agent task claim <id>                # Start working
agentic-agent task complete <id>             # Mark done
```

### OpenSpec (Lightweight)

```bash
agentic-agent openspec init "Name" \
  --from requirements.md                     # Create from spec

agentic-agent openspec list                  # View all changes
agentic-agent openspec show <id>             # Show details
agentic-agent openspec status <id>           # Check progress
agentic-agent openspec complete <id>         # Mark done
agentic-agent openspec archive <id>          # Archive
```

### SDD Full (High-Risk Only)

```bash
agentic-agent sdd start "Name" \
  --risk critical                            # Create initiative

agentic-agent sdd workflow show <id>         # View phase
agentic-agent sdd gate-check <id>            # Verify gates
agentic-agent sdd adr create \
  --title "Decision" --scope global          # Create decision
agentic-agent sdd adr list                   # View decisions
agentic-agent sdd adr resolve <id>           # Approve decision
```

---

## Workflow Comparison

```
╔═══════════════════════════════════════════════════════╗
║              TASK              OPENSPEC       SDD     ║
╠═══════════════════════════════════════════════════════╣
║ Scope      | Bug fix        | Feature        | Critical║
║ Timeline   | Hours-1 day    | 1-2 weeks      | 3+ weeks║
║ Team       | Solo           | 1-3 people     | Large   ║
║ Ceremony   | None           | Light          | Heavy   ║
║ Gates      | None           | Optional       | Required║
║ Setup      | 5 min          | 10 min         | 1 hour  ║
║ Spec File  | No             | proposal.md    | Full    ║
╚═══════════════════════════════════════════════════════╝
```

---

## Project Structure

### Single Package

```
myproject/
├── agnostic-agent.yaml          ← Config
├── requirements.md              ← Spec
├── src/
└── .agentic/
    ├── openspec/changes/        ← OpenSpec lifecycle
    │   └── my-feature/
    │       ├── proposal.md
    │       └── tasks.md
    └── tasks/
        ├── backlog.yaml
        ├── in-progress.yaml
        └── done.yaml
```

### Monorepo

```
mymonorepo/
├── agnostic-agent.yaml          ← Config (with multi-path)
├── specs/                       ← Shared specs
│   └── platform-strategy.md
├── packages/
│   ├── service-a/
│   │   ├── specs/               ← Per-package specs
│   │   │   ├── auth-oauth.md
│   │   │   └── logging.md
│   │   └── src/
│   └── service-b/
│       ├── specs/
│       │   └── cache.md
│       └── src/
└── .agentic/
    ├── openspec/changes/        ← Centralized tracking
    │   ├── oauth-integration/
    │   └── cache-layer/
    └── tasks/
```

---

## Config Templates

### Minimal (Single Package)

```yaml
# agnostic-agent.yaml
project:
  name: MyProject
  version: 0.1.0

agents:
  defaults:
    model: claude-3-5-sonnet-20241022

paths:
  openSpecDir: .agentic/openspec/changes
```

### Monorepo (Multi-Path)

```yaml
# agnostic-agent.yaml
project:
  name: MyMonorepo
  version: 0.1.0

agents:
  defaults:
    model: claude-3-5-sonnet-20241022

paths:
  openSpecDir: .agentic/openspec/changes
  specDirs:
    - packages/*/specs           # Auto-discover
    - apps/*/specs
    - specs                      # Root-level
```

---

## Implementation Workflow

### Day 1: Plan

```
Write spec file
      ↓
agentic-agent openspec init
      ↓
Generates: proposal.md + tasks.md
```

### Days 2-N: Build

```
agentic-agent task claim <task-1>
      ↓
Code in your package
      ↓
agentic-agent task complete <task-1>
      ↓
Repeat for each task
```

### Final Day: Verify & Release

```
agentic-agent openspec status <id>
      ↓
Check all acceptance criteria ✓
      ↓
agentic-agent openspec complete <id>
      ↓
agentic-agent openspec archive <id>
```

---

## 7 Implementation Prompts

**Use with Claude Code / Cursor / your agent:**

### 1️⃣ Spec Review
```
Review this spec and break into concrete tasks.

Spec: [file path]

For each acceptance criterion:
- What needs to be implemented
- Which files will be touched
- How to verify it works
```

### 2️⃣ Package-Specific Implementation
```
I'm implementing [task] in packages/[service]/.

Requirements from spec:
[paste spec content]

Please:
1. Design the solution
2. Show example code
3. List files to create/modify
4. Provide test structure
```

### 3️⃣ Integration (Cross-Package)
```
I'm integrating [feature] across:
- packages/service-a (provider)
- packages/service-b (consumer)

Current status:
- Service A: [status]
- Service B: [status]

What needs to happen in service-b?
- API contract
- Error handling
- Token refresh strategy
```

### 4️⃣ Testing
```
Write tests for [feature].

Acceptance criteria:
- [ ] [criterion 1]
- [ ] [criterion 2]
- [ ] [criterion 3]

Provide:
1. Unit test structure
2. Integration test structure
3. Mock strategy
4. Test file locations
```

### 5️⃣ Verification (Before Completion)
```
Verify [feature] is complete.

Checklist:
- [ ] All acceptance criteria met?
- [ ] Tests passing?
- [ ] Code review ready?
- [ ] Performance OK?
- [ ] Security audit passed?
- [ ] Docs updated?

Show me:
1. Test output (passing)
2. Files changed (git diff)
3. Any blocking issues
```

### 6️⃣ Dependency Analysis
```
Specs across packages:
- packages/a/specs/[spec-1].md (provider)
- packages/b/specs/[spec-2].md (consumer)

What's the implementation order?
- Which package starts first?
- What's the contract?
- What if B deployed before A?

Show dependency graph + rollout strategy.
```

### 7️⃣ Documentation
```
Generate documentation for [feature].

Based on: packages/[service]/specs/[spec].md

Sections:
1. Overview
2. API reference
3. Token lifecycle
4. Error handling
5. Testing locally
6. Troubleshooting

Format: Markdown with code examples.
```

---

## CLI Command Map

| What | Command |
|------|---------|
| **Create** | `agentic-agent task create "name"` |
| **List** | `agentic-agent task list` |
| **Claim** | `agentic-agent task claim <id>` |
| **Complete** | `agentic-agent task complete <id>` |
| **New Spec** | `agentic-agent openspec init "name" --from file.md` |
| **List Specs** | `agentic-agent openspec list` |
| **Show Spec** | `agentic-agent openspec show <id>` |
| **Mark Done** | `agentic-agent openspec complete <id>` |
| **Archive** | `agentic-agent openspec archive <id>` |

---

## Decision Trees

### Should I use OpenSpec?

```
Is this a bug fix?
├─ YES → Just use: agentic-agent task create
└─ NO → Continue...

Will it take > 1 week?
├─ NO → Use: agentic-agent task create
└─ YES → Continue...

Do you have clear acceptance criteria?
├─ NO → Write spec first, then: agentic-agent openspec init
└─ YES → agentic-agent openspec init --from requirements.md
```

### Single Package or Monorepo?

```
Single package
├─ Use: SMALL-PROJECTS.md workflow
└─ Config: minimal agnostic-agent.yaml

Monorepo (multiple packages)
├─ Use: MONOREPO-OPENSPEC.md workflow
└─ Config: multi-path specDirs
```

### When to Upgrade Workflows

```
Current: Using tasks only
├─ Scale up if: > 2 week feature
└─ Next: agentic-agent openspec

Current: Using OpenSpec
├─ Scale up if: High risk (payment/auth/PII)
└─ Next: Full SDD (agentic-agent sdd start)

Current: Using SDD
├─ Add gates if: Regulatory/critical
└─ Command: agentic-agent sdd gate-check
```

---

## Real Examples

### Example 1: Quick Bug Fix

```bash
# Create task
agentic-agent task create "Fix login timeout"

# Work
git checkout -b fix/login-timeout
# ... fix code ...
git commit -m "fix: increase session timeout to 30 min"

# Mark done
agentic-agent task claim <id>
agentic-agent task complete <id>

# Done! No spec needed.
```

### Example 2: Small Feature

```bash
# Write spec
cat > requirements.md << 'EOF'
# Product Search

## What
Add search + filters to product page

## Acceptance
- [ ] Search box appears
- [ ] Results filter in < 500ms
- [ ] Mobile responsive
EOF

# Initialize
agentic-agent openspec init "Product Search" --from requirements.md

# Work
agentic-agent task claim search-box-ui
# ... code ...
agentic-agent task complete search-box-ui

# Continue for all tasks, then mark feature done
agentic-agent openspec complete product-search
agentic-agent openspec archive product-search
```

### Example 3: Monorepo (OAuth)

```bash
# Root config already has multi-path specDirs

# Create specs in packages
mkdir -p packages/auth-service/specs
cat > packages/auth-service/specs/oauth.md << 'EOF'
# OAuth Provider
...
EOF

mkdir -p packages/web-ui/specs
cat > packages/web-ui/specs/oauth.md << 'EOF'
# OAuth Login UI
...
EOF

# Initialize all three
agentic-agent openspec init "OAuth Provider" \
  --from packages/auth-service/specs/oauth.md

agentic-agent openspec init "OAuth UI" \
  --from packages/web-ui/specs/oauth.md

# Teams work in parallel
cd packages/auth-service && agentic-agent task claim <id>
# ... code ...
agentic-agent task complete <id>

# Mark all done
agentic-agent openspec complete oauth-provider
agentic-agent openspec complete oauth-ui
```

---

## Spec Template (Copy-Paste)

```markdown
# [Feature Name]

## Problem
[What user problem does this solve? 1-2 sentences]

## Solution
[What you're building]

## Acceptance Criteria
- [ ] [Testable outcome]
- [ ] [Testable outcome]
- [ ] [Testable outcome]

## Tasks
1. [Task description]
   - Files: [what you'll modify]
   - Acceptance: [how to verify]

2. [Task description]
   - Files: [what you'll modify]
   - Acceptance: [how to verify]

## Notes
[Anything else the team should know]
```

---

## Help & Support

```
agentic-agent --help                    # Top-level help
agentic-agent task --help               # Task commands
agentic-agent openspec --help           # OpenSpec commands
agentic-agent sdd --help                # SDD commands
agentic-agent [command] --help          # Help for specific command
```

---

## Common Gotchas

| Problem | Solution |
|---------|----------|
| "Where do specs go?" | `packages/*/specs/` (monorepo) or `.agentic/spec/` (single) |
| "Can I have specs and tasks?" | Yes! Use both: specs for planning, tasks for tracking |
| "What if I change the spec?" | Update the file, re-run `openspec init` |
| "How do I integrate specs?" | Use `--from` flag: `openspec init --from other-spec.md` |
| "Can teams work in parallel?" | Yes! Use monorepo workflow + multi-path config |
| "Do I need gates?" | Only for high-risk (payments, auth, regulatory) |
| "Can I skip OpenSpec?" | Yes! Use `task` commands directly |

---

## TL;DR Matrix

| Scenario | Use | Command |
|----------|-----|---------|
| Bug fix | Task | `task create "fix: X"` |
| Small feature | OpenSpec | `openspec init --from spec.md` |
| Monorepo | Multi-path + OpenSpec | See MONOREPO-OPENSPEC.md |
| Critical feature | SDD | `sdd start --risk critical` |
| Parallel teams | Monorepo + worktrees | `agentic-agent task claim` × N |
| Final verification | Hard gate | `openspec complete` (after tests pass) |

---

## Next Step

1. Choose your scenario (Bug / Small / Monorepo / Critical)
2. Run the command above
3. Use the 7 prompts to guide implementation
4. Mark done: `openspec complete` or `task complete`

**Don't over-engineer. Start simple, add ceremony only when needed.**
