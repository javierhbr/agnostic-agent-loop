# Scenario 5: Working in an Existing Codebase

## Business Context

TaskFlow is a project management platform (3 years old, 150K lines of TypeScript). The team needs to:

1. **Add a new feature:** Export projects to CSV for reporting
2. **Modify an existing feature:** Change how task filtering works (currently: AND logic, want: OR logic with saved filters)

Both are in the same codebase. How to decide what approach to take?

---

## The Decision Point

### Feature 1: Export to CSV (Adding New)

**Scope:**
- 4 files to create/modify: `export-service.ts`, `api/export-endpoint.ts`, `ui/export-button.tsx`, `tests/export.test.ts`
- Affects: Service layer, API layer, UI, tests
- Dependencies: CSV library, existing permission system
- Risk: Moderate (new code, but integrates with existing permission logic)

**Decision:** Use **openspec** pipeline (4 files, cross-layer, new capability)

### Feature 2: Change Filtering Logic (Modifying Existing)

**Scope:**
- 6 files would need changing: `filters.ts`, `query-builder.ts`, `ui/filter-panel.tsx`, `api/filters-endpoint.ts`, `types/filter-types.ts`, `tests/filters.test.ts`
- Affects: Core query logic, API contracts, database queries
- Risk: HIGH (changes existing behavior, affects queries across the platform)
- Breaking change: Yes — existing saved filters with AND logic break

**Decision:** Use **openspec** (6 files, behavior change, breaking change, needs migration plan)

---

## Decision Rule

| Situation | Use |
|-----------|-----|
| Bug fix, typo, single file | Direct edit (just describe it) |
| Small feature (1–3 files, no breaking changes) | Direct edit (describe it, I implement) |
| New feature (4+ files, new capability) | **openspec** |
| Modify existing feature (5+ files, behavior change) | **openspec** |
| Change with side effects across the codebase | **openspec** (for traceability) |
| Have a PRD/spec already | **openspec init --from file** |

---

## How It Works: The Export Feature Example

### Phase 1: Brainstorm (15 min)

Tell Claude: *"I want to add CSV export to TaskFlow. Let me brainstorm this before spec-ing it out."*

Claude asks:
- Who needs to export? (Project managers, team leads, analysts?)
- What data? (All tasks? Selected tasks? Custom columns?)
- Export from where? (Task list, dashboard, project view?)
- Format requirements? (Column order, date format, custom fields?)
- Bulk export? (One project or many?)

### Phase 2: PRD (20 min)

Tell Claude: *"Use the product-wizard skill to create a PRD for CSV export based on my answers above."*

Claude generates:
- Problem: "Project managers can't report on projects; they export to Excel manually"
- Users: Project managers, team leads, analysts (50% of users)
- Success: "70% of power users use export within first month"
- Acceptance criteria:
  - Export single project as CSV
  - Include task name, status, assignee, due date
  - Maintain order from task list
  - Works on Chrome, Firefox, Safari

### Phase 3: Proposal & Tasks (15 min)

Tell Claude: *"Create a proposal, dev plan, and tasks using openspec. Call it 'csv-export'."*

Claude runs:
```bash
agentic-agent openspec init "CSV Export Feature" --from <prd-file>
```

You get:
- `proposal.md` (problem, approach, scope, acceptance)
- `tasks.md` (4 tasks: Service, API endpoint, UI, tests)
- Auto-imported into backlog

### Phase 4: Implement (2–3 hours)

For each task:
```bash
agentic-agent task claim TASK-001
# I implement...
agentic-agent task complete TASK-001
```

After all tasks:
```bash
agentic-agent openspec complete csv-export
agentic-agent openspec archive csv-export
```

---

## How It Works: The Filtering Change Example

### The Complication

Changing filtering from AND → OR logic is a **breaking change**:
- Old saved filters stop working
- Queries change behavior
- Users see different results

This requires:
- Migration script (convert AND filters to OR)
- Feature flag (roll out safely)
- Parallel testing (old vs new)
- User communication

### The Approach

Phase 0: **Brainstorm the migration strategy**
```
"I want to change task filtering from AND to OR logic. This is a breaking change.
Help me brainstorm the migration and rollout strategy."
```

Phase 1: **PRD the complete picture**
```
"Use product-wizard to create a PRD that covers:
- Why we're changing
- What breaks
- Migration plan
- Feature flag strategy
- Rollback plan
- Timeline"
```

Phase 2: **OpenSpec the decomposed tasks**
```
agentic-agent openspec init "Filtering Logic Change" --from <prd-file>
```

Tasks become:
1. Write data migration (convert AND to OR)
2. Update query builder (new logic)
3. Add feature flag (toggle old/new)
4. Update API contract (document breaking change)
5. Update UI (show new behavior)
6. Write comprehensive tests (verify migration)
7. Deploy with flag OFF (safety first)
8. Monitor, then flip flag in production

---

## Key Learnings

### ✅ Do This

1. **Always read existing code first** before proposing changes
   - Understand current patterns
   - Check how similar features are built
   - Identify side effects

2. **Use openspec for anything that touches multiple layers or has side effects**
   - Even if it seems small at first
   - Better to have structure than regret later

3. **Run `agentic-agent context generate <dir>` before editing a directory**
   - Gives me understanding of the codebase
   - Helps me write idiomatic code

4. **Use feature flags for breaking changes**
   - Safe rollout
   - Easy rollback
   - Measure impact

5. **Validate before completing**
   - Run tests
   - Check migrations work
   - Verify no regressions

### ❌ Don't Do This

- Skip brainstorming/PRD for "obvious" features
- Underestimate scope (4 files becomes 6 when you start)
- Mix breaking changes with new features
- Commit changes without running `agentic-agent validate`
- Bypass openspec to "save time" (it actually saves time)

---

## Related Examples

See the corresponding `prompts.md` for copy-paste prompts for both scenarios.

See the parent `QUICK-REFERENCE.md` and `SKILLS-GUIDE.md` for skill usage.
