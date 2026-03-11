# Platform MCP Tool Catalog

## Implemented v1 tools

- `get_platform_version`
- `list_platform_refs`
- `get_platform_ref`
- `get_jira_mapping`
- `validate_component_alignment`
- `validate_component_jira_chain`
- `detect_platform_drift_from_pinned_version`

## Rules

- query tools read platform truth
- validation tools use the component's pinned platform version by default
- no tool writes to the platform or component repo in v1

## Later additions

- `get_contract`
- `get_capability`
- `get_platform_adr`
- `list_platform_changes_since_version`
- `validate_platform_ref_exists`
- `explain_constraints_for_component`
