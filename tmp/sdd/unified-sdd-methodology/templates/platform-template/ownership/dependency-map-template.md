# Dependency Map

Platform version: `<yyyy.mm>`

## How to read this map

Each entry says:
- which component depends on which other component
- what the dependency is (contract, event, data)
- which impact tier applies

Impact tiers:

| Tier | Plain meaning | JIRA implication |
|------|---------------|-----------------|
| `must_change_together` | If one changes the shared contract, the other must change in the same release | Coordinate epics, same release window |
| `watch_for_breakage` | Consumer reads from the contract; a breaking change will break it | Open component epic when contract changes |
| `adapts_independently` | Loose coupling; consumer can absorb the change on its own schedule | No forced epic coordination needed |

---

## Dependencies

### `<component-a>`

| Depends on | Via | Impact tier |
|------------|-----|-------------|
| `<component-b>` | `<contract-ref or event>` | `must_change_together` |
| `<component-c>` | `<contract-ref or event>` | `watch_for_breakage` |

### `<component-b>`

| Depends on | Via | Impact tier |
|------------|-----|-------------|
| `<component-d>` | `<contract-ref or event>` | `adapts_independently` |

---

## Notes

Explain any non-obvious dependency decisions here.

---

## How to use this file

- Create one platform-level dependency map during the Platform phase.
- Update it when a new contract is added or an existing contract changes tier.
- The Team Lead reads this during Assess to determine:
  - which components are affected
  - what JIRA issue chain is needed
  - how many component epics are required
- Reference the impact tier in `platform-ref.yaml` under `impact`.
