# Development Plan Template

This file contains the full template and detailed guidelines for creating a `DEVELOPMENT_PLAN.md` file. Reference this when executing the dev-plans skill.

---

## Development Plan Structure

Create a new file called `DEVELOPMENT_PLAN.md` with this exact structure:

```markdown
# Development Plan for [PROJECT_NAME]

## Project Purpose and Goals

[Clear statement of what this project aims to achieve and why]

## Context and Background

[Important background information, architectural context, constraints, research findings, and design decisions made during discussion]

## Development Tasks

### Phase 1: [Phase Name]

- [ ] Task 1
  - [ ] Sub-task 1.1 (if needed)
  - [ ] Sub-task 1.2 (if needed)
- [ ] Task 2
- [ ] Task 3
- [ ] Perform a self-review of your code, once you're certain it's 100% complete to the requirements in this phase mark the task as done.
- [ ] STOP and wait for human review # (Unless the user has asked you to complete the entire implementation)

### Phase 2: [Phase Name]

- [ ] Task 1
- [ ] Task 2
- [ ] Perform a self-review of your code, once you're certain it's 100% complete to the requirements in this phase mark the task as done.
- [ ] STOP and wait for human review # (Unless the user has asked you to complete the entire implementation)

[Additional phases as needed]

## Important Considerations & Requirements

- [ ] Do not over-engineer the solution
- [ ] Do not add placeholder or TODO code
- [ ] [Additional requirements from conversation]
- [ ] [Architectural constraints]
- [ ] [Integration requirements]

## Technical Decisions

[Document any key technical decisions, trade-offs considered, and rationale for chosen approaches]

## Testing Strategy

[Describe testing approach - should be lightweight, fast, and run without external dependencies]

## Debugging Protocol

If issues arise during implementation:

- **Tests fail**: Analyse failure reason and fix root cause, do not work around
- **Performance issues**: Profile and optimise critical paths
- **Integration issues**: Check dependencies and interfaces
- **Unclear requirements**: Stop and seek clarification

## QA Checklist

- [ ] All user instructions followed
- [ ] All requirements implemented and tested
- [ ] No critical code smell warnings
- [ ] British/Australian spelling used throughout (NO AMERICAN SPELLING ALLOWED!)
- [ ] Code follows project conventions and standards
- [ ] Documentation is updated and accurate if needed
- [ ] Security considerations addressed
- [ ] Integration points verified (if applicable)
- [ ] [Project-specific QA criteria based on technology stack]
- [ ] [Additional QA criteria from user requirements]
```

---

## Task Granularity Rules

Each task must be **atomic** — one focused unit of work that can be completed and tested independently.

| Guideline | Example |
|-----------|---------|
| One concern per task | "Add user model" not "Add user model and auth endpoints" |
| Single layer | Prefer "Add API routes" + "Add UI components" over "Add full feature" |
| Testable in isolation | Each task has its own verifiable acceptance criteria (3-5 per task) |
| Split aggressively | Frontend + backend = two tasks. Infrastructure + logic = two tasks. |
| Target 10-20 tasks | Medium feature (3 screens + API + storage) = ~15 tasks |

Do NOT bundle multiple features, layers, or components into one task to reduce the count.

---

## Phase Structure

Each phase should:

1. **Be independently reviewable** — human can review and approve before next phase starts
2. **Have a self-review checkpoint** — agent reviews its own work before marking complete
3. **Have a STOP checkpoint** — explicit pause for human review (unless user requested auto-complete)
4. **Have 3-8 tasks** — neither too granular nor too chunky

---

## Writing Guidelines

- Use dashes with single spaces for markdown lists: `- [ ] Task`
- Do not include dates or time estimates in the plan itself
- Be clear, concise, and actionable
- Write in British English
- Use technical terminology consistently
- Avoid vague language — be specific about what needs to be done
- Example: ✅ "Create User model with email, password_hash, and created_at fields" not ❌ "Add user stuff"

---

## Quality Gates by Risk Level

Adjust QA checklist based on project risk tolerance:

| Level | Focus | QA Depth |
|-------|-------|----------|
| **High-risk production** | Strict QA, extensive testing, security audits | Comprehensive |
| **Internal tools/local dev** | Lighter QA, focus on functionality | Medium |
| **Open source contributions** | Follow project's contribution guidelines precisely | Per guidelines |
| **Prototypes/experiments** | Minimal QA, emphasis on learning and iteration | Light |

---

## Testing Philosophy

- Lightweight and fast — tests run in seconds, not minutes
- No external dependencies — no network, databases, or external services
- Tests run in isolation — no test order dependencies
- Cover critical paths and edge cases — not every line, but every behavior
- Integration tests for key workflows (if applicable) — test system pieces working together

---

## Technical Decision Documentation

When documenting technical decisions, include:

1. **The decision**: What was chosen and why
2. **Alternatives considered**: What else was evaluated
3. **Trade-offs**: What was gained and sacrificed
4. **Constraints**: Why this was best choice given constraints
5. **Implications**: What this affects downstream

Example:

```markdown
### Authentication Approach

**Decision:** Implemented JWT-based stateless authentication with refresh tokens.

**Alternatives considered:**
- Session-based (would require server-side state, not suitable for horizontal scaling)
- API keys (would not support user logout)

**Trade-offs:**
- Gained: Stateless design, horizontal scalability
- Sacrificed: Immediate revocation (refresh tokens can have short expiry)

**Constraints:** Project must support multiple API instances without shared session storage.

**Implications:** All clients must handle token refresh; logout requires blacklist or token expiry.
```
