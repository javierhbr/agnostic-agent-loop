# JIRA Conventions

## Recommended hierarchy

```text
[Platform issue]
        |
        +--> [Component epic A]
        |         |
        |         +--> [Story A1]
        |
        +--> [Component epic B]
                  |
                  +--> [Story B1]
```

## Rules

- platform issue tracks the shared outcome or shared truth change
- component epic tracks one repository or component delivery stream
- story tracks one reviewable slice or task group
- PRs should reference story keys
- stories should reference task IDs when possible
