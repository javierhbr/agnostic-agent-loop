# Platform and Component Alignment

Use this reference when the platform has:

- one master repository for shared platform truth
- many component repositories for local implementation
- JIRA for delivery coordination

Core model:

```text
[Platform master repo]
  canonical platform truth
        |
        v
[Component repo]
  versioned alignment + local OpenSpec artifacts
        |
        v
[JIRA]
  platform issue -> component epic -> stories -> PRs
```

Key rules:

- platform truth stays upstream in the master repo
- component repositories pin the platform version they align to
- OpenSpec artifacts stay local to the component repository
- JIRA links the delivery chain but does not replace the specs

Detailed methodology reference:

- `../../unified-sdd-methodology/canonical-platform-truth-and-component-alignment.md`

Template references:

- `../../unified-sdd-methodology/templates/README.md`
- `../../unified-sdd-methodology/templates/platform-template/README.md`
- `../../unified-sdd-methodology/templates/component-boilerplate/README.md`
- `../../unified-sdd-methodology/templates/platform-ref.yaml`
- `../../unified-sdd-methodology/templates/jira-traceability.yaml`
