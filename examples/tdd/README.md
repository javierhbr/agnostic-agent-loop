# TDD Skill Pack

The TDD skills have been consolidated and embedded into the CLI.

## Install

Project-level:

    agentic-agent skills install tdd --tool claude-code

Global/user-level:

    agentic-agent skills install tdd --tool claude-code --global

## Use with work command

    agentic-agent work --task <id> --follow-tdd

This decomposes the task into RED/GREEN/REFACTOR sub-tasks and verifies the TDD skill pack is installed.

## Supported tools

| Tool | Project directory |
|------|-------------------|
| claude-code | `.claude/skills/tdd/` |
| cursor | `.cursor/skills/tdd/` |
| gemini | `.gemini/skills/tdd/` |
| windsurf | `.windsurf/skills/tdd/` |
| antigravity | `.agent/skills/tdd/` |
| codex | `.codex/skills/tdd/` |

## List available packs

    agentic-agent skills list
