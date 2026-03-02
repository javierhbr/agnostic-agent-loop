# SDD & OpenSpec Documentation

Complete guides for using OpenSpec and SDD methodologies with the `agentic-agent` CLI.

---

## 📖 Quick Navigation

**Just want the essentials?** → [CHEATSHEET.md](./CHEATSHEET.md) (one page, print-friendly)

**Visual learner?** → [DIAGRAMS.md](./DIAGRAMS.md) (12 ASCII diagrams)

**Have a specific scenario?** → Jump to section below ⬇️

---

## Quick Start by Project Type

### 🚀 **Tiny Projects** (Bug fixes, small features, < 1 day)

**Skip everything.** Use `agentic-agent task` only:

```bash
agentic-agent task create "Your feature"
agentic-agent task claim <id>
# ... code ...
agentic-agent task complete <id>
```

→ See: [SMALL-PROJECTS.md](./SMALL-PROJECTS.md#for-tiny-projects-solo-developer-low-risk)

---

### 📦 **Small Projects** (Single package, 1-2 weeks, 1-3 people)

Use OpenSpec for lightweight spec + task tracking:

```bash
agentic-agent openspec init "Feature Name" --from requirements.md
agentic-agent task claim <task-id>
# ... code ...
agentic-agent task complete <task-id>
agentic-agent openspec complete <feature-id>
```

→ See: [SMALL-PROJECTS.md](./SMALL-PROJECTS.md#for-small-projects-1-2-people-1-service)

---

### 🏢 **Monorepo Projects** (Multiple packages, parallel teams)

Use OpenSpec with multi-path spec resolution:

```bash
# Root config with spec auto-discovery
# specs live in: packages/*/specs/

agentic-agent openspec init "Feature" --from packages/service-a/specs/auth.md
agentic-agent task claim <task-id>
# ... code in packages/service-a ...
agentic-agent task complete <task-id>
```

→ See: [MONOREPO-OPENSPEC.md](./MONOREPO-OPENSPEC.md)

Includes:
- Root config setup
- Per-package spec organization
- 7 implementation prompts
- Multi-service coordination
- Real examples (payment flow, auth integration)

---

### 🏗️ **Complex Projects** (High risk, multi-team, regulatory)

Use full SDD methodology with gates, ADRs, verification:

```bash
agentic-agent sdd start "Feature" --risk critical
agentic-agent sdd workflow show <id>
# ... runs through Analyst → Architect → Developer → Verifier phases ...
agentic-agent sdd gate-check <id>  # Verify before moving forward
```

→ See: [SDD-FULL.md](./SDD-FULL.md) (if available) or main SDD docs

---

## Documentation Files

| File | For | Use Case |
| --- | --- | --- |
| [README.md](./README.md) (this file) | Everyone | Navigation & quick reference |
| [CHEATSHEET.md](./CHEATSHEET.md) | Everyone | One-page quick reference (print-friendly) |
| [DIAGRAMS.md](./DIAGRAMS.md) | Visual learners | ASCII diagrams & flowcharts |
| [SMALL-PROJECTS.md](./SMALL-PROJECTS.md) | Solo devs, small teams | Bug fixes, simple features, lightweight specs |
| [MONOREPO-OPENSPEC.md](./MONOREPO-OPENSPEC.md) | Monorepo teams | Multi-package workflows, parallel development, 7 prompts |

---

## When to Use What

```
What are you building?

├─ Bug fix or small refactor?
│  └─ Just use: agentic-agent task create + complete
│     (No spec needed)
│
├─ Single feature (1-2 weeks)?
│  └─ Use: agentic-agent openspec init
│     + write simple proposal + tasks
│
├─ Multi-package feature (monorepo)?
│  └─ Use: MONOREPO-OPENSPEC.md workflow
│     + per-package specs + 7 prompts
│
└─ Critical/high-risk (payments, auth, regulatory)?
   └─ Use: Full SDD
      + all gates + verification + ADRs
```

---

## CLI Command Cheat Sheet

### Task Management

```bash
# Create & track simple tasks
agentic-agent task create "Task name"
agentic-agent task list
agentic-agent task claim <id>
agentic-agent task complete <id>
```

### OpenSpec (Lightweight)

```bash
# Create change from spec
agentic-agent openspec init "Name" --from requirements.md

# Track across project
agentic-agent openspec list
agentic-agent openspec show <id>
agentic-agent openspec status <id>

# Complete & archive
agentic-agent openspec complete <id>
agentic-agent openspec archive <id>
```

### SDD (Full Methodology)

```bash
# Start risk-classified initiative
agentic-agent sdd start "Name" --risk low|medium|high|critical

# View workflow state
agentic-agent sdd workflow show <id>

# Create architectural decisions
agentic-agent sdd adr create --title "Decision" --scope global|local
agentic-agent sdd adr list
agentic-agent sdd adr resolve <id>

# Gate checks before moving forward
agentic-agent sdd gate-check <spec-id>
```

---

## 7 Implementation Prompts (For All Projects)

Use these with Claude Code / Cursor / your agent:

1. **Spec Review** — Decompose spec into tasks
2. **Package-Specific** — Guide implementation in one package
3. **Integration** — Coordinate across packages
4. **Testing** — Structure tests for acceptance criteria
5. **Verification** — Ensure all criteria met before completion
6. **Dependency Analysis** — Multi-service implementation order
7. **Documentation** — Generate from specs

→ See: [MONOREPO-OPENSPEC.md - Prompts for Implementation](./MONOREPO-OPENSPEC.md#prompts-for-implementation)

(These work for any project size, not just monorepos!)

---

## Configuration

### Minimal Config (Single Package)

```yaml
project:
  name: MyProject
  version: 0.1.0

agents:
  defaults:
    model: claude-3-5-sonnet-20241022

paths:
  openSpecDir: .agentic/openspec/changes
```

### Monorepo Config (Multiple Packages)

```yaml
project:
  name: MyMonorepo
  version: 0.1.0

agents:
  defaults:
    model: claude-3-5-sonnet-20241022

paths:
  openSpecDir: .agentic/openspec/changes
  specDirs:
    - packages/*/specs      # Auto-discover package specs
    - apps/*/specs          # App specs
    - specs                 # Shared specs
```

---

## Real Examples

### Example 1: Single Service (Small Project)

**Feature:** Add search to product catalog

```bash
# Write spec
cat > requirements.md << 'EOF'
# Product Search

## Problem
Users spend 5+ min scrolling to find products.

## Solution
Add search + filters (name, category, price).

## Acceptance
- [ ] Search box on product page
- [ ] Results in < 500ms
- [ ] Mobile responsive
EOF

# Initialize
agentic-agent openspec init "Product Search" --from requirements.md

# Implement
agentic-agent task claim <task-1>
# ... code ...
agentic-agent task complete <task-1>

# Verify & finish
agentic-agent openspec complete product-search
agentic-agent openspec archive product-search
```

### Example 2: Monorepo (Multi-Service)

**Feature:** OAuth integration across 3 services

```bash
# Create specs in each package
cat > packages/auth-service/specs/oauth.md << 'EOF'
# OAuth Provider
Implement OAuth 2.0 provider for Google + GitHub
EOF

cat > packages/api-service/specs/oauth.md << 'EOF'
# OAuth Consumer
Consume OAuth tokens from auth-service
EOF

cat > packages/web-ui/specs/oauth.md << 'EOF'
# OAuth Login UI
Login buttons for Google + GitHub
EOF

# Initialize all three
agentic-agent openspec init "OAuth Provider" --from packages/auth-service/specs/oauth.md
agentic-agent openspec init "OAuth Consumer" --from packages/api-service/specs/oauth.md
agentic-agent openspec init "OAuth Login UI" --from packages/web-ui/specs/oauth.md

# Teams work in parallel
# Team A: packages/auth-service
# Team B: packages/api-service
# Team C: packages/web-ui

# Mark complete when all done
agentic-agent openspec complete oauth-provider
agentic-agent openspec complete oauth-consumer
agentic-agent openspec complete oauth-login-ui
```

---

## Decision Tree

**Choose your workflow based on:**

1. **How complex is it?**
   - Simple (< 1 week) → Use OpenSpec or just tasks
   - Complex (> 2 weeks) → Use full SDD

2. **How many people?**
   - Solo → Tasks only, or lightweight OpenSpec
   - Team (2-3) → OpenSpec
   - Large team → Full SDD with gates

3. **How risky is it?**
   - Low (internal features) → Tasks or OpenSpec
   - Medium (user-facing) → OpenSpec
   - High (payments, auth, regulatory) → Full SDD

4. **Is it a monorepo?**
   - Single package → SMALL-PROJECTS.md
   - Multiple packages → MONOREPO-OPENSPEC.md

---

## FAQ

**Q: Do I need specs for everything?**

A: No. Use specs only for:
- Features taking > 1 week
- Anything with acceptance criteria
- Coordination across teams/packages

Skip specs for: bug fixes, small refactors, urgent hotfixes

**Q: Can I switch workflows mid-project?**

A: Yes. Start simple (tasks), upgrade to OpenSpec if needed, escalate to SDD if risk increases.

**Q: What if specs conflict with my existing workflow?**

A: Use only what you need. The CLI is designed to be optional:
- Just use `agentic-agent task` for tracking
- Add OpenSpec when you need spec lifecycle
- Add SDD gates when risk/complexity increases

**Q: How do I migrate existing specs to OpenSpec?**

A: `agentic-agent openspec init --from <existing-spec-file>` will read it and create the OpenSpec structure.

**Q: Can I use OpenSpec without the full SDD?**

A: Yes! OpenSpec is lightweight and standalone. SDD is optional.

---

## Next Steps

1. **Read your scenario:**
   - Small project? → [SMALL-PROJECTS.md](./SMALL-PROJECTS.md)
   - Monorepo? → [MONOREPO-OPENSPEC.md](./MONOREPO-OPENSPEC.md)

2. **Set up config:**
   - Create `agnostic-agent.yaml` in your project root

3. **Create your first spec:**
   - Write simple requirements file
   - Run `agentic-agent openspec init`

4. **Start tracking work:**
   - Use `agentic-agent task` to claim & complete tasks
   - Use 7 prompts to guide implementation

5. **Finish feature:**
   - Verify acceptance criteria
   - Run `agentic-agent openspec complete`
   - Archive: `agentic-agent openspec archive`

---

## Support

- Questions about small projects? → [SMALL-PROJECTS.md](./SMALL-PROJECTS.md#faq-small-projects)
- Questions about monorepos? → [MONOREPO-OPENSPEC.md](./MONOREPO-OPENSPEC.md#faq-monorepo--openspec)
- CLI help: `agentic-agent --help` or `agentic-agent <command> --help`

---

## Summary

| Project Type | Use This | Tools | Timeline |
| --- | --- | --- | --- |
| Tiny | Tasks only | `task create/claim/complete` | Hours-1 day |
| Small | OpenSpec | `openspec init` + tasks | 1-2 weeks |
| Monorepo | Multi-path OpenSpec | See MONOREPO-OPENSPEC.md | 1-3 weeks |
| Complex/High-Risk | Full SDD | `sdd start` + gates + verification | 3-4 weeks |

**Use the right tool for your project. Don't over-engineer.**
