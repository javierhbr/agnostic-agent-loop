# Development Tools

This directory contains tool-specific configurations and integrations for various development environments and AI coding assistants.

## Directory Structure

```
tools/
├── .claude/          # Claude Code configurations
├── .cursor/          # Cursor IDE configurations
└── README.md         # This file
```

## Claude Code (.claude/)

[Claude Code](https://claude.ai/claude-code) is Anthropic's official CLI for Claude AI. The `.claude/` directory contains:

- **Skills** - Custom Claude Code skills for this project
- **Keybindings** - Custom keyboard shortcuts
- **Configuration** - Claude-specific settings

### Using Claude Code Skills

The skills in this project help Claude understand the Agentic Agent framework:

```bash
# Example: Use a custom skill
/skill-name
```

See [.claude/README.md](.claude/README.md) for available skills.

## Cursor IDE (.cursor/)

[Cursor](https://cursor.sh/) is an AI-powered code editor. The `.cursor/` directory contains:

- **Rules** - Cursor-specific coding rules and conventions
- **Configuration** - Editor settings optimized for this project

### Using Cursor

Open the project in Cursor and it will automatically load configurations from `.cursor/`.

## Why This Directory?

Development tool configurations are grouped here to:

1. **Reduce root clutter** - Keep the project root clean and focused
2. **Clear organization** - All tool configs in one place
3. **Easy discovery** - New contributors can find tool setups quickly
4. **Maintainability** - Update tool configs without searching the entire project

## Adding New Tool Configurations

When adding support for a new development tool:

1. Create a subdirectory: `tools/.toolname/`
2. Add tool-specific configurations
3. Update this README with usage instructions
4. Add a tool-specific README if needed: `tools/.toolname/README.md`

## Related Documentation

- [Main README](../README.md) - Project overview
- [CLI Tutorial](../docs/guide/CLI_TUTORIAL.md) - Using the agentic-agent CLI
- [Project Layout](../docs/development/project-layout.md) - Directory structure explanation

## Supported Tools

### AI Coding Assistants
- ✅ Claude Code - Official Anthropic CLI
- ✅ Cursor - AI-powered code editor

### Coming Soon
- GitHub Copilot configuration
- JetBrains IDE integration
- VS Code extension settings

## Notes

- Tool configurations are **optional** - the framework works without them
- Configurations here are **project-specific** - they won't affect your global settings
- Feel free to customize these configs for your workflow
