# Upgrade and Safety

## Upgrade the CLI
```bash
uv tool install specify-cli --force --from git+https://github.com/github/spec-kit.git
```

## Update project files
```bash
specify init --here --force --ai codex
```

## What updates
- slash command files
- scripts
- templates
- shared memory files

## What stays safe
- `specs/` contents
- plans and tasks already generated under `specs/`
- source code
- git history

## Important warning
A documented upgrade caveat is that the constitution memory file can be overwritten during project refresh.
Protect it by backing it up or restoring it from git after the upgrade.

## Safe upgrade routine
1. back up customized constitution/templates
2. run the project refresh
3. restore or merge local customizations
4. verify agent commands and templates
