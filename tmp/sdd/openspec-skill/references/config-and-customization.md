# Config and Customization

## Project config
OpenSpec project configuration lives in `openspec/config.yaml`.

Typical fields:
- `schema:` default schema name
- `context:` durable project context injected into artifact generation
- `rules:` per-artifact constraints injected only for matching artifact IDs

## Recommended use of context
Put stable information in `context:` such as:
- tech stack
- API conventions
- testing expectations
- backward-compatibility rules
- language preference

Do not overload context with temporary ticket chatter.

## Recommended use of rules
Use `rules:` for artifact-specific instructions, for example:
- proposal must include rollback plan
- specs must use Given/When/Then scenarios
- design must include sequence diagrams for complex flows

## Schema precedence
When multiple schema hints exist, resolution order is:
1. CLI flag
2. change-local metadata
3. project config
4. default built-in schema

## Custom schemas
Create a custom schema only when the workflow itself must change.
A practical approach:
1. fork the built-in `spec-driven` schema
2. edit `schema.yaml`
3. update artifact templates
4. version the schema in the repo
