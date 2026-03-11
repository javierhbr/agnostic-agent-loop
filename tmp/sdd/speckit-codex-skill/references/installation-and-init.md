# Installation and Init

## Prerequisites
- Python 3.11+
- Git
- uv
- a supported AI coding agent

## Common initialization patterns
### One-shot use
```bash
uvx --from git+https://github.com/github/spec-kit.git specify init <PROJECT_NAME>
```

### Initialize current directory
```bash
uvx --from git+https://github.com/github/spec-kit.git specify init .
# or
uvx --from git+https://github.com/github/spec-kit.git specify init --here
```

### Specify the agent
```bash
uvx --from git+https://github.com/github/spec-kit.git specify init <PROJECT_NAME> --ai codex
```

### Force script type
```bash
uvx --from git+https://github.com/github/spec-kit.git specify init <PROJECT_NAME> --script sh
uvx --from git+https://github.com/github/spec-kit.git specify init <PROJECT_NAME> --script ps
```

## Operational notes
- script variants exist for both Bash and PowerShell
- active feature detection follows the current Git branch naming pattern like `001-feature-name`
- initialization installs the command and template scaffolding used by the workflow
