# Agent Rules - Agnostic Agent Loop

**These rules apply to all AI agents working on this codebase, regardless of tool (Claude, OpenCode, Copilot, etc.).**

---

## Context Management Rules

**Violating the letter of these rules is violating the spirit of the rules.**

### Mandatory Workflow: Read-Before-Write

Before modifying ANY file in a directory, you **must**:

1. Check if `context.md` exists in that directory
2. If it exists, **read it** - identify `Allowed` and `Forbidden` dependencies, architectural role, and constraints
3. If it does NOT exist and the directory contains source files, generate context documentation
4. Only after reading context may you edit files

After completing changes that alter architecture, exports, or dependencies:
5. Update the `context.md` to reflect changes

### New Directory Workflow

When creating a new directory:

1. Create the directory and files
2. **Immediately** generate context documentation for the directory
3. Review the generated context and enrich it with architectural role and dependency rules

### Hexagonal Architecture Boundaries

Enforce these dependency rules when working in specific layers:

| Layer | Can Depend On | Cannot Depend On |
|-------|--------------|-------------------|
| **Core/Domain** | Nothing | Application, Infrastructure, Config |
| **Core/Application** | Domain only | Infrastructure, Config |
| **Infrastructure/Adapters** | Domain, Application | Other adapters directly |
| **Infrastructure/Config** | All layers (DI glue) | — |

If a `context.md` specifies dependency rules, those override the defaults above.

### context.md Template

```markdown
# Context: [Directory Name]

## Responsibility
[One sentence description]

## Architectural Role
- **Type:** [Core / Infrastructure / Adapter / Config]
- **Direction:** [Inbound / Outbound / Internal]

## Dependency Rules
- **Allowed:** [List of packages/layers this can import]
- **Forbidden:** [List of packages/layers this must NOT import]
```

### Red Flags — STOP and Fix

- ❌ Editing a file without reading that directory's `context.md` first
- ❌ Creating a directory without generating `context.md`
- ❌ Importing a package listed as Forbidden in `context.md`
- ❌ "The task is too small to need context" — no task is too small
- ❌ "I already understand the code" — code shows structure, context.md shows rules
- ❌ "The user told me where to put it" — user instructions do not override architectural constraints
- ❌ "I'll update context.md later" — later means never; update immediately

---

## Code Quality Rules

### Always

- ✅ Read context before writing
- ✅ Keep tasks small and focused
- ✅ Update documentation when changing behavior
- ✅ Run tests before completing work
- ✅ Follow existing patterns in the codebase
- ✅ Write clear commit messages

### Never

- ❌ Skip reading context.md
- ❌ Violate architectural boundaries
- ❌ Mix concerns in a single file
- ❌ Leave TODO comments without tracking
- ❌ Commit commented-out code
- ❌ Ignore linter warnings

---

## Task Management Rules

### Before Starting Work

1. Understand the task description and acceptance criteria
2. Read relevant context.md files
3. Identify files that need changes
4. Plan the approach

### During Work

1. Make incremental, testable changes
2. Update context.md if architecture changes
3. Run tests frequently
4. Keep track of what's done vs. what remains

### Before Completing Work

1. Run validation checks
2. Verify all acceptance criteria are met
3. Update documentation
4. Review your changes for quality
5. Ensure context.md is current

---

## Documentation Rules

### What to Document

- **Architectural decisions** - Why this pattern/structure?
- **Dependency rules** - What can/cannot depend on what?
- **Constraints** - Performance, security, compatibility requirements
- **Non-obvious behavior** - Edge cases, special handling
- **Integration points** - How components connect

### What Not to Document

- **Obvious code behavior** - Don't repeat what code says
- **Implementation details** - Code is the documentation for "how"
- **Temporary notes** - Use task tracking instead

### Documentation Locations

| Type | Location |
|------|----------|
| Architectural rules | `context.md` in each directory |
| Project overview | `.agentic/context/global-context.md` |
| Architecture decisions | `docs/architecture/decisions/` |
| Feature specs | `.agentic/openspec/changes/` |
| API documentation | inline code comments |

---

## Communication Rules

### With Users

- Be clear and concise
- Explain what you're doing and why
- Ask clarifying questions when requirements are ambiguous
- Confirm understanding before large changes
- Report blockers immediately

### In Code

- Write self-explanatory code
- Use descriptive names
- Comment only non-obvious logic
- Keep functions small and focused
- Follow language idioms

### In Commits

- Write clear commit messages
- One logical change per commit
- Reference task IDs when applicable
- Explain "why" not "what" in commit body

---

## Agent-Specific Tool Integration

Different agent tools have different commands. Refer to your tool-specific instructions:

- **Claude**: See `CLAUDE.md`
- **OpenCode**: See `OPENCODE.md` 
- **Copilot**: See `COPILOT.md`

Common patterns across tools:
- Generate context for a directory
- List available tasks
- Claim/start a task
- Mark a task complete
- Run validation checks
- Build context bundles

---

## Emergency Protocols

### If You Realize You Violated a Rule

1. **STOP immediately**
2. Assess the damage
3. Revert if necessary
4. Fix properly
5. Document what went wrong
6. Continue

### If Requirements Conflict with Architecture

1. **Raise the conflict** to the user
2. Explain the architectural constraint
3. Propose compliant alternatives
4. Let the user decide: refactor architecture or change requirements
5. Document the decision

### If You're Unsure

1. **Ask** - Don't guess
2. Read more context
3. Look for similar examples in the codebase
4. Propose an approach for validation
5. Proceed only after confirmation

---

**Remember: These rules exist to maintain code quality and architectural integrity. Following them makes the codebase better for everyone.**
