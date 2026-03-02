# OpenClaw Project Coordinator Setup

This directory contains two collaborating OpenClaw agents configured to work with the agnostic-agent CLI:

## Agents

### 🧩 TechLead (Dev Team Coordinator)
- **Role:** Routes technical work, manages tasks, runs gate-checks, spawns workers
- **Lives in:** `tech-lead/`
- **Commands:** `agentic-agent task`, `agentic-agent validate`, `agentic-agent sdd gate-check`
- **Coordination:** Reads spec-ready signals from ProductLead, announces task completions

### 🔭 ProductLead (Product Coordinator)
- **Role:** Defines requirements, creates specs, hands work to TechLead
- **Lives in:** `product-lead/`
- **Commands:** `agentic-agent openspec`, `product-wizard`, `sdd analyst`
- **Coordination:** Announces spec-ready to TechLead, listens for completion updates

## File Organization

Each agent has:
- `IDENTITY.md` — Name, emoji, vibe, avatar
- `SOUL.md` — Values, principles, behavioral boundaries
- `AGENTS.md` — Session startup checklist, coordination protocol
- `TOOLS.md` — CLI cheat sheet, environment paths, commands
- `BOOTSTRAP.md` — First-run onboarding script (delete after first activation)

Plus:
- `USER.md` — Shared human context (both agents read this)

## How to Activate

1. **Activate TechLead:**
   ```bash
   # Run BOOTSTRAP.md interactively
   # It will ask about your project and update IDENTITY.md, USER.md, SOUL.md
   # Then delete BOOTSTRAP.md
   ```

2. **Activate ProductLead:**
   ```bash
   # Same process — run BOOTSTRAP.md interactively
   # ProductLead will update the same USER.md with additional context
   # Then delete BOOTSTRAP.md
   ```

## How They Talk

Both agents communicate via `.agentic/coordination/announcements.yaml`:

- **ProductLead → TechLead:** Announces when a spec is approved and ready for development
  ```yaml
  from_agent: product-lead
  to_agent: tech-lead
  status: spec-ready
  ```

- **TechLead → ProductLead:** Announces when tasks are complete and ready for review
  ```yaml
  from_agent: tech-lead
  to_agent: product-lead
  status: complete
  ```

## Integration with agnostic-agent

Both agents are fully integrated with the agnostic-agent CLI:

### TechLead uses:
- `agentic-agent task list` → Read project state
- `agentic-agent task claim <ID>` → Reserve task
- `agentic-agent context build --task <ID>` → Get full context
- `agentic-agent sdd gate-check <spec-id>` → Validate quality gates
- `agentic-agent validate` → Run all validators
- `agentic-agent task complete <ID>` → Mark task done

### ProductLead uses:
- `agentic-agent openspec init` → Start new spec
- `agentic-agent openspec check` → Validate spec
- `product-wizard` skill → Generate PRDs
- `sdd analyst` skill → Requirements discovery
- `sdd architect` skill → Architecture design

## Next Steps

1. Fill in the `[Your name]` placeholders in `USER.md`
2. Run `tech-lead/BOOTSTRAP.md` interactively
3. Run `product-lead/BOOTSTRAP.md` interactively
4. Delete both BOOTSTRAP.md files after activation
5. Start coordinating: ProductLead creates specs, TechLead routes work

---

**Reference:** See the plan at `.claude/plans/agile-hugging-diffie.md` for complete setup details.
