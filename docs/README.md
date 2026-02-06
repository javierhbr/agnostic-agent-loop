# Agentic Agent Documentation

Welcome to the Agentic Agent framework documentation! This directory contains comprehensive guides, architecture documentation, and integration references.

## 📚 Documentation Index

### Getting Started

- **[CLI Tutorial](guide/CLI_TUTORIAL.md)** - Step-by-step guide to using the agentic-agent CLI
- **[Main README](../README.md)** - Project overview and quick start

### User Guides

Located in [`guide/`](guide/):
- [CLI Tutorial](guide/CLI_TUTORIAL.md) - Command-line interface usage
- Workflow guides (coming soon):
  - Beginner workflow
  - Intermediate workflow
  - Advanced workflow

### Architecture & Design

Located in [`architecture/`](architecture/):
- [Implementation Plan](architecture/plan.md) - Overall project architecture and roadmap
- [Architecture Decision Records](architecture/decisions/) - ADRs documenting key design decisions

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

### Development

Located in [`development/`](development/):
- Contributing guidelines (coming soon)
- Testing guide (coming soon)
- [Project Layout](development/project-layout.md) - Directory structure explanation (coming soon)

## 🗂️ Documentation Structure

```
docs/
├── README.md                      # This file - documentation index
├── guide/                         # User guides and tutorials
│   ├── CLI_TUTORIAL.md
│   └── workflows/                 # Step-by-step workflows
├── architecture/                  # Architecture and design
│   ├── plan.md
│   └── decisions/                 # Architecture Decision Records
├── bdd/                          # BDD testing documentation
│   ├── BDD_IMPLEMENTATION_SUMMARY.md
│   ├── BDD_QUICK_REFERENCE.md
│   └── BDD_GUIDE.md
├── integrations/                  # Third-party integrations
│   ├── ralph/                     # Ralph PDR integration
│   └── agentskills/               # AgentSkills.io compatibility
└── development/                   # Developer documentation
    ├── contributing.md            # (coming soon)
    ├── testing.md                 # (coming soon)
    └── project-layout.md          # (coming soon)
```

## 🔍 Finding What You Need

### I want to...

**Learn how to use the CLI**
→ Start with [CLI Tutorial](guide/CLI_TUTORIAL.md)

**Understand the architecture**
→ Read [Implementation Plan](architecture/plan.md)

**Write or run tests**
→ Check out the [BDD Guide](bdd/BDD_GUIDE.md)

**Integrate with Ralph PDR**
→ See [Ralph PDR Workflow](integrations/ralph/RALPH_PDR_WORKFLOW.md)

**Contribute to the project**
→ Read [Contributing Guide](development/contributing.md) (coming soon)

**Understand the codebase structure**
→ Check [Project Layout](development/project-layout.md) (coming soon)

## 📦 Package-Level Documentation

In addition to this centralized documentation, each major internal package has its own README explaining its specific functionality:

- [cmd/agentic-agent/README.md](../cmd/agentic-agent/README.md) - CLI commands overview (coming soon)
- [internal/tasks/README.md](../internal/tasks/README.md) - Task management system (coming soon)
- [internal/context/README.md](../internal/context/README.md) - Context generation and management (coming soon)
- [internal/validator/README.md](../internal/validator/README.md) - Validation rules (coming soon)
- [internal/ui/README.md](../internal/ui/README.md) - UI components (coming soon)

## 🤝 Contributing to Documentation

Documentation improvements are always welcome! When contributing:

1. Keep documentation close to the code it describes
2. Update cross-references when moving or renaming docs
3. Follow the existing structure and style
4. Add links to the appropriate index files

## 📝 License

This documentation is part of the Agentic Agent project. See the main [README](../README.md) for license information.
