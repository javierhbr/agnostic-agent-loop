# Agentic Agent Documentation

Welcome to the Agentic Agent framework documentation! This directory contains comprehensive guides, architecture documentation, and integration references.

## Documentation Index

### Getting Started

- **[CLI Tutorial](guide/CLI_TUTORIAL.md)** - Step-by-step guide to using the agentic-agent CLI
- **[Main README](../README.md)** - Project overview and quick start

### User Guides

Located in [`guide/`](guide/):
- [CLI Tutorial](guide/CLI_TUTORIAL.md) - Command-line interface usage, including task-level skill refs, simplify command, and ATDD/BDD workflow
- [Spec-Driven Development](SPEC_DRIVEN_DEVELOPMENT.md) - Multi-framework spec resolution, context bundles, autopilot mode

### Architecture & Design

Located in [`architecture/`](architecture/):
- [Implementation Plan](architecture/plan.md) - Overall project architecture and roadmap
- [TODO-01](architecture/TODO-01.md) - Foundational framework: agent adapters, context control, session protocol
- [TODO-02](architecture/TODO-02.md) - Advanced features: skill layer, orchestrator state machine, CUOC protocol
- [TODO-03](architecture/TODO-03.md) - TOON integration, skill versioning, wire format strategy
- [Task Templates Feature](architecture/TASK_TEMPLATES_FEATURE.md) - `task sample-task` and `task from-template` design

### BDD Testing

Located in [`bdd/`](bdd/):
- [BDD Implementation Summary](bdd/BDD_IMPLEMENTATION_SUMMARY.md) - Complete BDD framework overview
- [BDD Quick Reference](bdd/BDD_QUICK_REFERENCE.md) - Quick reference for writing BDD tests
- [BDD Guide](bdd/BDD_GUIDE.md) - Detailed guide to BDD testing practices

### Integrations

Located in [`integrations/`](integrations/):

#### Ralph PDR
- [Ralph PDR Workflow](integrations/ralph/RALPH_PDR_WORKFLOW.md) - Integration with Ralph PDR methodology
- [Ralph Integration Summary](integrations/ralph/RALPH_INTEGRATION_SUMMARY.md) - Summary of Ralph integration features

#### AgentSkills.io
- [AgentSkills Compliance](integrations/agentskills/AGENTSKILLS_COMPLIANCE.md) - Compatibility with agentskills.io standard

## Documentation Structure

```
docs/
├── README.md                          # This file - documentation index
├── SPEC_DRIVEN_DEVELOPMENT.md         # Spec-driven development guide
├── guide/                             # User guides and tutorials
│   └── CLI_TUTORIAL.md
├── architecture/                      # Architecture and design
│   ├── plan.md                        # Implementation plan and roadmap
│   ├── TODO-01.md                     # Foundational framework spec
│   ├── TODO-02.md                     # Advanced features spec
│   ├── TODO-03.md                     # TOON integration spec
│   └── TASK_TEMPLATES_FEATURE.md      # Task template system design
├── bdd/                               # BDD testing documentation
│   ├── BDD_IMPLEMENTATION_SUMMARY.md
│   ├── BDD_QUICK_REFERENCE.md
│   └── BDD_GUIDE.md
└── integrations/                      # Third-party integrations
    ├── ralph/                         # Ralph PDR integration
    │   ├── RALPH_PDR_WORKFLOW.md
    │   └── RALPH_INTEGRATION_SUMMARY.md
    └── agentskills/                   # AgentSkills.io compatibility
        └── AGENTSKILLS_COMPLIANCE.md
```

## Finding What You Need

### I want to...

**Learn how to use the CLI**
> Start with [CLI Tutorial](guide/CLI_TUTORIAL.md)

**Understand the architecture**
> Read [Implementation Plan](architecture/plan.md)

**Use spec-driven development with Spec Kit, OpenSpec, or native specs**
> See [Spec-Driven Development](SPEC_DRIVEN_DEVELOPMENT.md)

**Use task-level skill refs or the simplify command**
> See [CLI Tutorial — Skill Refs and Simplify](guide/CLI_TUTORIAL.md#scenario-4-skill-refs-and-code-simplification)

**Write or run tests**
> Check out the [BDD Guide](bdd/BDD_GUIDE.md)

**Integrate with Ralph PDR**
> See [Ralph PDR Workflow](integrations/ralph/RALPH_PDR_WORKFLOW.md)

## Package-Level Documentation

Key internal packages:

| Package | Purpose |
|---------|---------|
| [internal/skills/](../internal/skills/) | Agent detection, skill packs, installer, registry, drift detection |
| [internal/tasks/](../internal/tasks/) | Task CRUD, lifecycle, readiness checks, progress tracking |
| [internal/tracks/](../internal/tracks/) | Feature/bug track management with brainstorm, spec, plan files |
| [internal/plans/](../internal/plans/) | Markdown plan parser and updater (checkbox status tracking) |
| [internal/context/](../internal/context/) | Directory context generation, global context, rolling summary |
| [internal/encoding/](../internal/encoding/) | Context bundle assembly (YAML, JSON, TOON) |
| [internal/orchestrator/](../internal/orchestrator/) | Autopilot loop, state machine, task archival |
| [internal/status/](../internal/status/) | Project status dashboard (counts, blockers, readiness) |
| [internal/gitops/](../internal/gitops/) | Read-only git integration (branch, commits, changed files) |
| [internal/specs/](../internal/specs/) | Multi-directory spec resolution |
| [internal/config/](../internal/config/) | YAML config loading and agent-specific overrides |
| [internal/simplify/](../internal/simplify/) | Code simplification bundle builder |
| [internal/validator/](../internal/validator/) | Validation framework and rules |
| [internal/ui/](../internal/ui/) | Bubble Tea TUI components |
| [internal/token/](../internal/token/) | Token counting and budget estimation |

## Contributing to Documentation

Documentation improvements are always welcome! When contributing:

1. Keep documentation close to the code it describes
2. Update cross-references when moving or renaming docs
3. Follow the existing structure and style
4. Add links to the appropriate index files

## License

This documentation is part of the Agentic Agent project. See the main [README](../README.md) for license information.
