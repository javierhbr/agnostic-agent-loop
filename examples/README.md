# Examples

This directory contains example projects and demonstrations of the Agentic Agent framework.

## Available Examples

### `foo/` - Simple Example Package
A minimal Go package demonstrating basic project structure. Moved from `src/foo/` to clarify it's example code, not main framework source.

**Purpose:** Shows how a simple Go package might be organized when using the Agentic Agent framework.

### `test-sandbox/` - Test Environment
An isolated environment for testing framework features without affecting the main codebase.

**Purpose:** Safe space to experiment with agentic-agent commands and workflows.

## Creating Your Own Example

To create a new example project:

1. Create a directory: `examples/my-example/`
2. Initialize with agentic-agent:
   ```bash
   cd examples/my-example
   agentic-agent init
   ```
3. Add your example code and documentation
4. Update this README with a description

## Example Project Structure

A typical example project using the Agentic Agent framework:

```
examples/my-example/
├── README.md                    # Example documentation
├── agnostic-agent.yaml          # Framework configuration
├── .agentic/                    # Framework runtime
│   ├── tasks/                   # Task definitions
│   ├── context/                 # Context files
│   └── agent-rules/             # Agent rules
└── src/                         # Example source code
    └── main.go
```

## Running Examples

### Initialize an Example
```bash
cd examples/foo
agentic-agent init
```

### Work with Tasks
```bash
# List tasks
agentic-agent task list backlog

# Claim a task
agentic-agent task claim task-001

# Complete a task
agentic-agent task complete task-001
```

### Generate Context
```bash
# Generate context for current directory
agentic-agent context generate .
```

## Example Use Cases

### Basic Project
Demonstrates:
- Simple Go package structure
- Task management workflow
- Context generation

### Multi-Agent Project (Coming Soon)
Demonstrates:
- Multiple agents working in parallel
- Task decomposition
- Conflict resolution
- Shared context management

### CI/CD Integration (Coming Soon)
Demonstrates:
- Running agentic-agent in CI/CD pipelines
- Automated task validation
- Coverage tracking

## Notes for Framework Development

**Important:** Examples in this directory serve two purposes:

1. **User Education** - Show how to use the framework
2. **Framework Testing** - Provide realistic test scenarios

When modifying examples:
- Keep them simple and focused on one concept
- Document clearly what each example demonstrates
- Ensure examples work with the latest framework version
- Add tests if the example includes complex logic

## Related Documentation

- [Main README](../README.md) - Project overview
- [CLI Tutorial](../docs/guide/CLI_TUTORIAL.md) - Command-line usage
- [BDD Guide](../docs/bdd/BDD_GUIDE.md) - Testing with examples

## Contributing Examples

We welcome new examples! When contributing:

1. Ensure the example is clear and well-documented
2. Test that it works with the current framework version
3. Update this README with your example
4. Consider adding a demo script or walkthrough
