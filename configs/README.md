# Configuration Directory

This directory contains configuration files and templates for the Agentic Agent framework.

## Files

### `agnostic-agent.yaml`

The default configuration file template. When users run `agentic-agent init`, this file is copied to their project directory and serves as a starting point for project-specific configuration.

**Configuration Options:**
- `prd_output_path`: Path for PRD (Product Requirements Document) output
- `progress_file`: Path to track task progress
- `archive_dir`: Directory for completed task archives
- Custom paths for context, tasks, and other framework components

## Templates Directory

The `templates/` directory contains initialization templates that are used when setting up a new project with `agentic-agent init`.

### Template Structure

```
templates/
└── init/
    ├── agnostic-agent.yaml    # Project configuration template
    ├── agent-rules/           # Base agent rules
    ├── context/               # Context file templates
    └── tasks/                 # Task YAML templates
```

### Important Note on Templates

**For Developers**: The templates in this directory (`configs/templates/`) serve as the **source of truth** for the project structure. However, due to Go's `//go:embed` limitation (which cannot embed files outside the package directory), these templates are **copied** to [internal/project/templates/](../internal/project/templates/) for embedding into the compiled binary.

**When modifying templates:**
1. Edit templates in `configs/templates/` (this directory)
2. Copy changes to `internal/project/templates/` to update the embedded version
3. Or use the update script: `scripts/sync-templates.sh` (if available)

This dual-location approach ensures:
- ✅ Clean project structure (configs are in `configs/`)
- ✅ Self-contained binary (templates embedded in the executable)
- ✅ Single source of truth (`configs/templates/` is the canonical location)

## Usage

### For End Users

When you run `agentic-agent init`, the CLI will:
1. Create the `.agentic/` directory structure in your project
2. Copy the default `agnostic-agent.yaml` configuration file
3. Set up initial context and task files from templates

### For Contributors

When adding new templates or modifying existing ones:
1. Add/modify files in `configs/templates/`
2. Sync to `internal/project/templates/` for embedding
3. Update this README if adding new template types
4. Test with `agentic-agent init` to ensure templates are applied correctly

## See Also

- [Main README](../README.md) - Project overview
- [CLI Tutorial](../docs/guide/CLI_TUTORIAL.md) - Command-line usage guide
- [Project Layout](../docs/development/project-layout.md) - Directory structure explanation
