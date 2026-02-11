# Claude Code Custom Rules

These rules are injected into CLAUDE.md via `{{ .AgentRules }}` during generation.

- When starting a task, run `agentic-agent task claim <TASK_ID>`.
- Before editing files in a directory, run `agentic-agent context generate <DIR>`.
- After completing work, run `agentic-agent task complete <TASK_ID>`.
- Use the Bash tool for `go test` and `go build` commands.
- Prefer editing existing files over creating new ones.
- When writing tests, use table-driven test patterns.
