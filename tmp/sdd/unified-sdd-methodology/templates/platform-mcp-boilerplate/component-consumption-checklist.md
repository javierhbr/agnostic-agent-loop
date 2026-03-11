# Component Consumption Checklist

Use this checklist before a component team relies on the local platform MCP.

- the developer has a local clone of the platform repo
- the local MCP server points to the correct platform repo path
- `platform-ref.yaml` exists in the component repo
- `jira-traceability.yaml` exists in the component repo
- the component is validating against the pinned platform version by default
- the team understands that the MCP server is read-only
- the team still keeps local behavior in OpenSpec artifacts inside the component repo
