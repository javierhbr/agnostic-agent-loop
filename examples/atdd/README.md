# ATDD Workflow: Acceptance Criteria → Executable Tests → Implementation

Write acceptance tests from openspec task criteria before implementation. Every acceptance criterion becomes a failing test; implementation makes them pass.

---

## What You'll Learn

- Extract acceptance criteria from openspec tasks
- Translate criteria into executable test cases
- Follow the ATDD cycle: criteria → failing tests → implementation → green
- Combine ATDD with TDD for full test coverage
- Use `skill_refs: [atdd]` on tasks for targeted skill inclusion

---

## 0. Setup

```bash
# From the project root
go build -o examples/atdd/agentic-agent ./cmd/agentic-agent
cd examples/atdd

# Initialize a project
rm -rf .agentic agnostic-agent.yaml
./agentic-agent init --name "atdd-demo"

# Install the ATDD skill pack
./agentic-agent skills install atdd --tool claude-code
```

---

## 1. Create an OpenSpec Change with Tasks

```bash
# Initialize a change from requirements
./agentic-agent openspec init "User Auth" --from requirements.md
```

Fill in `proposal.md` and `tasks.md`, then the tasks auto-import:

```bash
./agentic-agent task list --no-interactive
```

```text
--- BACKLOG ---
[TASK-001] [user-auth] Implement login endpoint
[TASK-002] [user-auth] Add session persistence
[TASK-003] [user-auth] Write integration tests
```

---

## 2. View Task Acceptance Criteria

```bash
./agentic-agent task show TASK-001 --no-interactive
```

```text
ID: TASK-001
Title: [user-auth] Implement login endpoint
Acceptance Criteria:
  - User can log in with email and password
  - Invalid credentials return 401 with error message
  - Successful login returns JWT token
  - Token expires after 24 hours
```

These four criteria become your test cases.

---

## 3. Claim and Start ATDD Cycle

```bash
./agentic-agent task claim TASK-001
```

### Phase 1: RED — Write Failing Acceptance Tests

Tell your agent:

```text
Read the acceptance criteria from .agentic/openspec/changes/user-auth/tasks/01-login-endpoint.md.
Write one test per criterion. All tests must fail. Do NOT implement any production code.
```

The agent writes:

```typescript
describe('Login endpoint', () => {
  test('user can log in with valid credentials', async () => {
    const res = await request(app).post('/auth/login')
      .send({ email: 'user@test.com', password: 'valid123' });
    expect(res.status).toBe(200);
    expect(res.body.token).toBeDefined();
  });

  test('invalid credentials return 401', async () => {
    const res = await request(app).post('/auth/login')
      .send({ email: 'user@test.com', password: 'wrong' });
    expect(res.status).toBe(401);
    expect(res.body.error).toContain('Invalid credentials');
  });

  test('successful login returns JWT', async () => {
    const res = await request(app).post('/auth/login')
      .send({ email: 'user@test.com', password: 'valid123' });
    const decoded = jwt.decode(res.body.token);
    expect(decoded).toBeTruthy();
    expect(decoded.email).toBe('user@test.com');
  });

  test('token expires after 24 hours', async () => {
    const res = await request(app).post('/auth/login')
      .send({ email: 'user@test.com', password: 'valid123' });
    const decoded = jwt.decode(res.body.token);
    const expiry = decoded.exp - decoded.iat;
    expect(expiry).toBe(86400); // 24 hours in seconds
  });
});
```

Run tests — all 4 fail (no login endpoint exists yet):

```bash
npm test
# FAIL: 4 tests failed
```

**Gate:** All tests must fail for the right reason (missing implementation, not syntax errors).

### Phase 2: GREEN — Implement to Pass Tests

Tell your agent:

```text
Write the minimum code to make all 4 acceptance tests pass. No extra features.
```

Run tests after implementation:

```bash
npm test
# PASS: 4 tests passed
```

### Phase 3: REFACTOR

Improve code quality. Run tests after every change:

```bash
npm test
# PASS: 4 tests passed (still green)
```

---

## 4. Complete and Move On

```bash
./agentic-agent task complete TASK-001
./agentic-agent openspec status user-auth
```

Repeat the ATDD cycle for TASK-002, TASK-003, etc.

---

## 5. Combining ATDD + TDD

ATDD and TDD work at different levels:

| Level | Skill | What it tests |
|-------|-------|---------------|
| **Acceptance** | atdd | Feature works as specified (user-facing criteria) |
| **Unit** | tdd | Functions work correctly (internal logic) |

Use both on a task:

```yaml
tasks:
  - id: "TASK-001"
    title: "Implement login endpoint"
    skill_refs:
      - atdd
      - tdd
```

Workflow: Write acceptance tests (ATDD) → then for each function needed, write unit tests (TDD) → implement → all tests green.

---

## 6. Combining ATDD + Ralph

Use ATDD tests as Ralph's convergence signal:

```
/ralph-loop "Implement TASK-001: Login endpoint.

Read: .agentic/openspec/changes/user-auth/tasks/01-login-endpoint.md

Acceptance tests already written in tests/auth.test.ts.
Run npm test each iteration.
When all 4 acceptance tests pass: <promise>TASK COMPLETE</promise>
" --max-iterations 10 --completion-promise "TASK COMPLETE"
```

Ralph iterates until the acceptance tests pass — clear, measurable convergence.

---

## Quick Reference

| Step | Command / Action |
|------|-----------------|
| Show criteria | `task show TASK-ID --no-interactive` |
| Claim task | `task claim TASK-ID` |
| Write failing tests | One test per acceptance criterion |
| Verify RED | `npm test` — all fail |
| Implement | Minimum code to pass |
| Verify GREEN | `npm test` — all pass |
| Refactor | Improve, keep green |
| Complete | `task complete TASK-ID` |

## The ATDD Cycle

```text
task show → extract criteria → write failing tests → verify RED
    ↓
implement → verify GREEN → refactor → task complete
    ↓
next task...
```
