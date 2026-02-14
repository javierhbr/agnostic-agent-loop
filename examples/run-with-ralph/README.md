# Run with Ralph: Iterative Task Implementation

Execute tasks using the Ralph Wiggum iterative loop methodology. Works in **any AI chat tool** (Claude Code, Cursor, Copilot, Windsurf) — just type `/ralph-loop`.

---

## What You'll Learn

- Use the `/ralph-loop` skill in AI chat for iterative task completion
- Configure checkpoint/resume for long-running tasks
- Alternative: Use autopilot mode for batch processing in terminal
- Process tasks with automatic git tracking and convergence detection
- Handle complex tasks with checkpoint safety nets

---

## 0. Prerequisites

- Agentic Agent CLI installed and configured
- Tasks in `.agentic/tasks/backlog.yaml` (can be from openspec or manual)
- AI chat tool with skill support (Claude Code, Cursor, etc.)
- **Optional:** Ralph Wiggum plugin (provides advanced `/ralph-loop` features)

---

## 1. List Available Tasks

```bash
agentic-agent task list --no-interactive
```

```text
--- BACKLOG ---
[TASK-936281-1] [todo-pwa] Set up project structure with Vite, React, TypeScript
[TASK-936281-2] [todo-pwa] Configure TypeScript and ESLint
[TASK-936281-3] [todo-pwa] Set up Tailwind with dark mode configuration
...
```

---

## 2. Claim a Task (BEFORE Ralph)

```bash
agentic-agent task claim TASK-936281-2
```

This records your git branch and timestamp. **Never skip this** — `task complete` needs the claim timestamp to capture commits.

---

## 3. Get Task Details

```bash
agentic-agent task show TASK-936281-2 --no-interactive
```

```text
ID: TASK-936281-2
Title: [todo-pwa] Configure TypeScript and ESLint
Description: Set up TypeScript with strict mode and ESLint with React/TypeScript rules...
Spec Refs:
  - todo-pwa/proposal.md
  - todo-pwa/tasks/02-typescript-eslint.md
Acceptance Criteria:
  - tsconfig.json configured with strict mode enabled
  - ESLint installed and configured for React + TypeScript
  - No ESLint errors on existing code
  - Path aliases configured (@/ for src/)
  - VS Code settings recommended for team consistency
```

Extract three things:
1. **Spec Refs** → files Ralph reads each iteration
2. **Acceptance Criteria** → Ralph's completion conditions
3. **Description** → implementation hints

---

## 4. Build the Ralph Prompt

Template — fill from `task show` output:

```text
You are implementing TASK-936281-2: Configure TypeScript and ESLint.

## Context — read these files first:
- .agentic/openspec/changes/todo-pwa/proposal.md
- .agentic/openspec/changes/todo-pwa/tasks/02-typescript-eslint.md
- .agentic/context/tech-stack.md

## Acceptance Criteria — ALL must pass:
- tsconfig.json configured with strict mode enabled
- ESLint installed and configured for React + TypeScript
- No ESLint errors on existing code
- Path aliases configured (@/ for src/)
- VS Code settings recommended for team consistency

## On each iteration:
1. Read the task detail file for technical notes
2. Implement what's missing
3. Run npx tsc --noEmit && npx eslint . to verify
4. If ALL criteria pass, output <promise>TASK COMPLETE</promise>
```

---

## 5. Launch Ralph Loop in AI Chat

### Option A: Simple Ralph Loop (Recommended)

Just type in your AI chat:

```
/ralph-loop
```

The skill will automatically:
1. Get your claimed task with `agentic-agent task continue`
2. Read all spec files referenced in the task
3. Understand acceptance criteria
4. Iteratively implement until ALL criteria pass
5. Complete task with `agentic-agent task complete`

**Example in Claude Code:**

```
You: /ralph-loop

AI: I'm using the ralph-loop skill to iteratively complete this task.

─── Step 1: Get Task ───
$ agentic-agent task continue

Working on: TASK-936281-2 - Configure TypeScript and ESLint

Acceptance criteria:
  - tsconfig.json configured with strict mode
  - ESLint installed and configured
  - No ESLint errors
  - Path aliases configured (@/ for src/)
  - VS Code settings recommended

Specs: todo-pwa/proposal.md, todo-pwa/tasks/02-typescript-eslint.md

─── Step 2: Read Specs ───
Reading .agentic/openspec/changes/todo-pwa/proposal.md...
Reading .agentic/openspec/changes/todo-pwa/tasks/02-typescript-eslint.md...

─── Iteration 1 ───
Creating tsconfig.json with strict mode...
Installing ESLint dependencies...

$ npm install -D eslint @typescript-eslint/parser @typescript-eslint/eslint-plugin
$ npm install -D eslint-plugin-react eslint-plugin-react-hooks

Creating .eslintrc.json...

$ npx eslint .
ERROR: 'React' is defined but never used

Status: 3/5 criteria met
Continuing...

─── Iteration 2 ───
Fixing ESLint errors...
Adding React import to components...

$ npx eslint .
✓ No errors

$ npx tsc --noEmit
✓ Type check passed

Status: 5/5 criteria met ✅

$ agentic-agent task complete TASK-936281-2 --learnings "Configured TypeScript strict mode and ESLint"

<promise>TASK COMPLETE</promise>
```

### Option B: Ralph Loop with Custom Prompt

If you need more control:

```
You: I need to implement TASK-936281-2. Use /ralph-loop but make sure to:
- Enable strict null checks in tsconfig
- Use Airbnb ESLint config
- Add prettier integration

AI: I'm using the ralph-loop skill with your custom requirements...
```

### Option C: Autopilot CLI (Non-Interactive)

For batch processing in terminal:

```bash
# Configure checkpoints (optional - see step 5a)
# Edit agnostic-agent.yaml first

# Run autopilot with agent execution
agentic-agent --agent claude-code autopilot start --execute-agent --max-iterations 1
```

This will:
1. Automatically claim the next task from backlog
2. Execute the AI agent to complete it
3. Auto-complete on success
4. Save checkpoints during execution

---

## 5a. Configure Checkpoints (For Long Tasks)

For complex tasks that might take many iterations, configure checkpoint behavior in `agnostic-agent.yaml`:

```yaml
# agnostic-agent.yaml
checkpoint:
  # Checkpoint every N iterations (default: 5)
  iteration_interval: 5

  # Checkpoint at these token percentages (default: [0.5, 0.75, 0.9])
  token_thresholds: [0.5, 0.75, 0.9]
```

**Why Use Checkpoints?**

- ✅ Never lose progress if interrupted (Ctrl+C, network issues, etc.)
- ✅ Automatically resume from last checkpoint on restart
- ✅ Warning when approaching token limits (80%, 90%)
- ✅ Auto-cleanup when task completes successfully

**Example Scenarios:**

```yaml
# Aggressive checkpointing (critical tasks)
checkpoint:
  iteration_interval: 2
  token_thresholds: [0.1, 0.25, 0.5, 0.75, 0.9]

# Minimal checkpointing (simple tasks)
checkpoint:
  iteration_interval: 10
  token_thresholds: [0.9]

# High token limit agent (Gemini)
checkpoint:
  iteration_interval: 5
  token_thresholds: [0.8, 0.95]
```

**How It Works:**

```bash
# Start autopilot with checkpoints
$ agentic-agent --agent claude autopilot start --execute-agent

--- Iteration 5/10 ---
🤖 Executing claude-code agent...
  ✅ Agent completed (tokens: 18000, total: 89000)
  💾 Checkpoint saved (iteration 5, 66.7% complete)

# If interrupted (Ctrl+C), resume later:
$ agentic-agent --agent claude autopilot start --execute-agent

--- Iteration 1/10 ---
Next task: [TASK-936281-2] Configure TypeScript and ESLint

🤖 Executing claude-code agent...
📌 Resuming from checkpoint (iteration 5, 89000 tokens used)
  ✅ Agent completed (tokens: 15000, total: 104000)
```

**Checkpoint Storage:**

```
.agentic/checkpoints/
├── TASK-936281-2-001.json    # Iteration 1
├── TASK-936281-2-005.json    # Iteration 5 checkpoint
├── TASK-936281-2-010.json    # Iteration 10 checkpoint
└── TASK-936281-2-latest.json # Latest checkpoint
```

**For more details:** See [checkpoint-resume.md](../../docs/checkpoint-resume.md)

---

## 6. Complete the Task

```bash
agentic-agent task complete TASK-936281-2
agentic-agent openspec status todo-pwa
```

```text
Change: todo-pwa
  Total: 22  Done: 1  In Progress: 0  Pending: 21
```

---

## 7. Repeat for Next Task

```bash
agentic-agent task claim TASK-936281-3
agentic-agent task show TASK-936281-3 --no-interactive
# Build ralph prompt from output...
/ralph-loop "..." --max-iterations 10 --completion-promise "TASK COMPLETE"
agentic-agent task complete TASK-936281-3
```

---

## Full Change Lifecycle

```text
agentic-agent task list
    ↓
For each task:
    ├── task claim TASK-ID
    ├── task show TASK-ID → extract criteria + spec refs
    ├── /ralph-loop "prompt" --max-iterations 10
    ├── task complete TASK-ID
    └── openspec status CHANGE-ID
    ↓
agentic-agent openspec complete CHANGE-ID
agentic-agent openspec archive CHANGE-ID
```

---

## 8. Autopilot vs Ralph Loop: When to Use What

### Use `/ralph-loop` in AI Chat When:
- ✅ Working interactively in your IDE
- ✅ Want to see iteration progress in real-time
- ✅ Task requires judgment or clarification
- ✅ Learning or experimenting
- ✅ Single task focus

**Example:**
```
You: /ralph-loop

AI: [Shows each iteration with visible progress]
  ─── Iteration 1 ───
  Implementing...
  Testing...
  Status: 2/4 criteria met

  ─── Iteration 2 ───
  Fixing remaining issues...
  Status: 4/4 criteria met ✅
  <promise>TASK COMPLETE</promise>
```

### Use `autopilot start --execute-agent` in Terminal When:
- ✅ Processing multiple tasks in batch
- ✅ Running in CI/CD or automation
- ✅ Want structured log output
- ✅ Non-interactive environment
- ✅ Background processing

**Example:**
```bash
$ agentic-agent --agent claude autopilot start --execute-agent --max-iterations 5

╔═══════════════════════════════════════════════════════════╗
║           Agentic Agent - Autopilot Mode                 ║
╚═══════════════════════════════════════════════════════════╝

🤖 Agent execution: ENABLED (claude-code)
🔄 Max iterations: 5

--- Iteration 1/5 ---
Next task: [TASK-001] Configure TypeScript
  ✅ Agent completed (tokens: 3456)
  ✅ All acceptance criteria met!
  ✅ Task TASK-001 completed successfully

--- Iteration 2/5 ---
All tasks complete. Autopilot finished.
```

### Detecting CLI Commands in Chat

If you paste an autopilot CLI command in your AI chat:

```
You: agentic-agent --agent claude autopilot start --execute-agent --max-iterations 1

AI: I see you've pasted an autopilot CLI command. In AI chat, we use the /ralph-loop
skill instead, which provides the same iterative behavior with better visibility.

Let me run the ralph-loop for you:

─── Step 1: Get Task ───
...
```

The skill automatically detects this and uses ralph-loop instead!

**For more details:** See [autopilot-vs-ralph-loop.md](../../docs/autopilot-vs-ralph-loop.md)

---

## Combining with ATDD (Recommended for Critical Features)

### Why ATDD + Ralph Loop = ✅

**ATDD (Acceptance Test-Driven Development)** ensures you're building the right thing.
**Ralph Loop** ensures you build it correctly through iteration.

Together they create an unbreakable workflow:
1. **RED** - Write failing acceptance tests from criteria
2. **GREEN** - Ralph-loop iterates until all tests pass
3. **REFACTOR** - Optional cleanup after ralph completes

### ATDD Workflow with Ralph Loop

#### Step 1: Claim and Analyze Task

```bash
$ agentic-agent task claim TASK-001
$ agentic-agent task show TASK-001

ID: TASK-001
Title: Implement user authentication
Acceptance Criteria:
  - User can login with email/password
  - JWT token is returned on success
  - Invalid credentials return 401
  - Password is hashed with bcrypt
  - Tests have 80%+ coverage
```

#### Step 2: Write Acceptance Tests (ATDD RED Phase)

In your AI chat:

```
You: Write acceptance tests for TASK-001 based on the criteria.
Follow the ATDD pattern - tests should fail initially.

Don't implement anything yet, just create failing tests that define success.

AI: I'll create acceptance tests for the authentication feature.

Creating tests/auth/login.test.ts:

describe('User Authentication', () => {
  describe('POST /api/login', () => {
    it('should return JWT token for valid credentials', async () => {
      const response = await request(app)
        .post('/api/login')
        .send({ email: 'test@example.com', password: 'password123' })
        .expect(200);

      expect(response.body).toHaveProperty('token');
      expect(response.body.token).toMatch(/^eyJ/); // JWT format
    });

    it('should return 401 for invalid credentials', async () => {
      await request(app)
        .post('/api/login')
        .send({ email: 'test@example.com', password: 'wrong' })
        .expect(401);
    });

    it('should hash passwords with bcrypt', async () => {
      const user = await User.findOne({ email: 'test@example.com' });
      const isHashed = await bcrypt.compare('password123', user.password);
      expect(isHashed).toBe(true);
    });
  });
});

$ npm test

FAIL tests/auth/login.test.ts
  ✗ should return JWT token for valid credentials
    → POST /api/login route not found
  ✗ should return 401 for invalid credentials
    → POST /api/login route not found
  ✗ should hash passwords with bcrypt
    → User model not defined

Test Suites: 1 failed
Tests: 0 passed, 3 failed

✓ Tests are failing as expected (RED phase)
```

#### Step 3: Use Ralph-Loop to Make Tests Pass (GREEN Phase)

```
You: /ralph-loop

AI: I'm using the ralph-loop skill to iteratively complete this task.

─── Step 1: Get Task ───
$ agentic-agent task continue

Working on: TASK-001 - Implement user authentication

Acceptance criteria:
  - User can login with email/password
  - JWT token is returned on success
  - Invalid credentials return 401
  - Password is hashed with bcrypt
  - Tests have 80%+ coverage

─── Step 2: Understand Acceptance Tests ───
Reading tests/auth/login.test.ts...

I can see 3 failing tests that define success:
1. JWT token returned for valid credentials
2. 401 for invalid credentials
3. Password hashing with bcrypt

─── Iteration 1 ───
Creating User model with bcrypt password hashing...
Implementing POST /api/login route...
Adding JWT token generation...

$ npm test

PASS tests/auth/login.test.ts (partial)
  ✓ should hash passwords with bcrypt
  ✗ should return JWT token for valid credentials
    → JWT_SECRET not configured
  ✗ should return 401 for invalid credentials
    → Error handling not implemented

Status: 1/5 criteria met (tests: 1/3 passing)
Continuing...

─── Iteration 2 ───
Adding JWT_SECRET to environment config...
Implementing error handling for invalid credentials...

$ npm test

PASS tests/auth/login.test.ts
  ✓ should return JWT token for valid credentials
  ✓ should return 401 for invalid credentials
  ✓ should hash passwords with bcrypt

$ npm test -- --coverage

Coverage: 85% (above 80% threshold)

Status: 5/5 criteria met ✅

$ agentic-agent task complete TASK-001 \
  --learnings "Implemented JWT auth with bcrypt" \
  --learnings "100% acceptance test coverage"

<promise>TASK COMPLETE</promise>
```

### ATDD + Unit Tests: Complete Coverage Strategy

#### Layer 1: Acceptance Tests (End-to-End)

Test the **full feature** from user perspective:

```typescript
// tests/acceptance/auth.test.ts
describe('Authentication Feature', () => {
  it('complete login flow works', async () => {
    // User registers
    await request(app).post('/api/register')
      .send({ email: 'new@example.com', password: 'pass123' });

    // User logs in
    const response = await request(app).post('/api/login')
      .send({ email: 'new@example.com', password: 'pass123' });

    expect(response.status).toBe(200);
    expect(response.body.token).toBeDefined();

    // User accesses protected route with token
    const profile = await request(app).get('/api/profile')
      .set('Authorization', `Bearer ${response.body.token}`);

    expect(profile.status).toBe(200);
    expect(profile.body.email).toBe('new@example.com');
  });
});
```

#### Layer 2: Integration Tests (Component Interaction)

Test **components working together**:

```typescript
// tests/integration/auth-service.test.ts
describe('AuthService', () => {
  it('authenticates user and returns token', async () => {
    const authService = new AuthService();
    const userRepo = new UserRepository();

    await userRepo.create({ email: 'test@example.com', password: 'pass' });
    const result = await authService.login('test@example.com', 'pass');

    expect(result.token).toBeDefined();
    expect(result.user.email).toBe('test@example.com');
  });
});
```

#### Layer 3: Unit Tests (Individual Functions)

Test **individual functions** in isolation:

```typescript
// tests/unit/password-hasher.test.ts
describe('PasswordHasher', () => {
  it('hashes password with bcrypt', async () => {
    const hasher = new PasswordHasher();
    const hashed = await hasher.hash('mypassword');

    expect(hashed).not.toBe('mypassword');
    expect(await bcrypt.compare('mypassword', hashed)).toBe(true);
  });

  it('validates correct password', async () => {
    const hasher = new PasswordHasher();
    const hashed = await hasher.hash('mypassword');

    expect(await hasher.verify('mypassword', hashed)).toBe(true);
    expect(await hasher.verify('wrongpass', hashed)).toBe(false);
  });
});
```

### Ralph Loop Test Strategy

When ralph-loop iterates, it verifies at **all three levels**:

```
─── Iteration N ───
Making changes...

Verification (from broadest to most specific):
1. $ npm run test:acceptance  → End-to-end flows work?
2. $ npm run test:integration → Components integrate?
3. $ npm run test:unit        → Functions correct?
4. $ npm run test:coverage    → Coverage > threshold?

Status: [X/Y criteria met]
```

### Example: Full ATDD + Ralph Workflow

**Task: Add password reset feature**

```bash
# 1. Claim task
$ agentic-agent task claim TASK-042

# 2. In AI chat - write tests first
You: Write acceptance, integration, and unit tests for TASK-042 password reset.

Tests should cover:
- Acceptance: Full reset flow (request → email → reset → login)
- Integration: ResetService + EmailService + UserRepo
- Unit: Token generation, expiration validation

Don't implement yet - RED phase only.

AI: [Creates 3 layers of failing tests]

$ npm test
Tests: 0 passed, 12 failed ✓ (RED phase complete)

# 3. Ralph loop to make tests pass
You: /ralph-loop

AI: [Iterates until all 12 tests pass]

─── Iteration 5 ───
All tests passing:
  ✓ Acceptance: 4/4 tests
  ✓ Integration: 5/5 tests
  ✓ Unit: 3/3 tests
  ✓ Coverage: 92%

<promise>TASK COMPLETE</promise>

# 4. Complete task
$ agentic-agent task complete TASK-042
```

### Benefits of ATDD + Ralph Loop

| Benefit | How It Helps |
|---------|--------------|
| **Clear Definition of Done** | Acceptance tests = criteria = done |
| **No Scope Creep** | Ralph stops when tests pass, nothing more |
| **Regression Safety** | Tests prevent breaking changes later |
| **Confidence** | Green tests = feature works as specified |
| **Documentation** | Tests show how feature should behave |
| **Faster Iteration** | Ralph verifies quickly with `npm test` |

### Pro Tips for ATDD + Ralph

#### 1. Write Tests from Criteria

```yaml
# Task acceptance criteria
acceptance:
  - User receives reset email with token
  - Token expires after 1 hour
  - Reset updates password
```

Maps directly to tests:

```typescript
it('user receives reset email with token', ...);
it('token expires after 1 hour', ...);
it('reset updates password', ...);
```

#### 2. Use Test-First for Complex Logic

For algorithms, edge cases, or business rules:

```
You: Before implementing password strength validator, write unit tests for:
- Minimum 8 characters
- At least 1 uppercase, 1 lowercase, 1 number
- No common passwords (from list)
- Edge cases: empty string, null, special chars

Then /ralph-loop to implement.

AI: [Writes 10+ failing unit tests, then implements until all pass]
```

#### 3. Layer Tests Appropriately

- **Acceptance**: 1-2 tests per feature (happy path + critical error)
- **Integration**: 3-5 tests per component interaction
- **Unit**: Many tests per function (cover all branches)

#### 4. Ralph Verifies All Layers

```
─── Iteration 3 ───
Running verification suite...

$ npm run test:acceptance
✓ 2/2 acceptance tests pass

$ npm run test:integration
✗ 3/5 integration tests fail
  → AuthService not handling expired tokens

Status: Criteria not met, continuing...

─── Iteration 4 ───
Fixed: Added token expiration check in AuthService

$ npm run test:integration
✓ 5/5 integration tests pass

$ npm test -- --coverage
✓ Coverage: 87% (threshold: 80%)

Status: 5/5 criteria met ✅
```

### Testing Coverage Metrics

Ralph can enforce coverage thresholds:

```json
// package.json
{
  "jest": {
    "coverageThreshold": {
      "global": {
        "branches": 80,
        "functions": 80,
        "lines": 80,
        "statements": 80
      }
    }
  }
}
```

Ralph verifies each iteration:

```
$ npm test -- --coverage

Coverage not met:
  ✗ Branches: 75% (threshold: 80%)
  ✓ Functions: 85%
  ✓ Lines: 82%
  ✓ Statements: 83%

Status: Coverage criteria not met, continuing...
```

---

## Quick Reference

### AI Chat Commands (Interactive)

| Task | Command |
|------|---------|
| List tasks | `agentic-agent task list` |
| Claim task | `agentic-agent task claim TASK-ID` |
| Show details | `agentic-agent task show TASK-ID` |
| **Start ralph-loop** | `/ralph-loop` (in AI chat) |
| Cancel ralph | `/cancel-ralph` (if using Ralph Wiggum plugin) |
| Complete task | `agentic-agent task complete TASK-ID` |
| Check progress | `agentic-agent openspec status CHANGE-ID` |

### Terminal Commands (Batch/CI)

| Task | Command |
|------|---------|
| Run autopilot | `agentic-agent --agent <name> autopilot start --execute-agent` |
| Dry run | `agentic-agent autopilot start --dry-run` |
| Set max iterations | `agentic-agent autopilot start --execute-agent --max-iterations 5` |
| View checkpoints | `ls -la .agentic/checkpoints/` |
| Resume from checkpoint | `agentic-agent autopilot start --execute-agent` (auto-resumes) |

### Configuration Files

| File | Purpose |
|------|---------|
| `agnostic-agent.yaml` | Agent config, checkpoint settings |
| `.agentic/tasks/backlog.yaml` | Available tasks |
| `.agentic/tasks/in-progress.yaml` | Currently claimed tasks |
| `.agentic/checkpoints/` | Checkpoint storage (auto-managed) |
| `.agentic/openspec/changes/` | OpenSpec proposals and tasks |

## Critical Rules

| Rule | Why |
|------|-----|
| Always claim before starting | Git tracking needs the claim timestamp |
| Use `/ralph-loop` in AI chat | Automatic iteration with convergence detection |
| Configure checkpoints for long tasks | Prevents progress loss on interruption |
| One task per ralph-loop | Multiple tasks = no convergence |
| Run `task complete` after finishing | Captures commits, updates progress tracking |
| Check token usage warnings | Pause/resume before hitting limits |

---

## Complete Example: From Task to Done

### Scenario: Implement TypeScript Configuration

**Step 1: Claim the task**

```bash
$ agentic-agent task claim TASK-936281-2

✓ Task TASK-936281-2 claimed
  Branch: feature/typescript-config
  Timestamp: 2026-02-13T15:30:00Z
```

**Step 2: In your AI chat (Claude Code, Cursor, etc.)**

```
You: /ralph-loop

AI: I'm using the ralph-loop skill to iteratively complete this task.

─── Step 1: Get Task ───
$ agentic-agent task continue

Working on: TASK-936281-2 - Configure TypeScript and ESLint

Acceptance criteria:
  - tsconfig.json configured with strict mode
  - ESLint installed and configured
  - No ESLint errors
  - Path aliases configured
  - VS Code settings recommended

Specs:
  - todo-pwa/proposal.md
  - todo-pwa/tasks/02-typescript-eslint.md

─── Step 2: Read Specs ───
[AI reads both spec files...]

─── Iteration 1 ───
Creating tsconfig.json...
Installing ESLint...

$ npm install -D eslint @typescript-eslint/parser @typescript-eslint/eslint-plugin

$ npx eslint .
ERROR: 'React' is defined but never used

Status: 3/5 criteria met
Continuing...

─── Iteration 2 ───
Adding React import to components...
Configuring path aliases in tsconfig...

$ npx eslint .
✓ No errors

$ npx tsc --noEmit
✓ Type check passed

Status: 5/5 criteria met ✅

$ agentic-agent task complete TASK-936281-2 \
  --learnings "Configured TypeScript strict mode" \
  --learnings "Set up ESLint with React/TS rules"

✓ Task TASK-936281-2 completed successfully
  Commits captured: 3 commits on feature/typescript-config

<promise>TASK COMPLETE</promise>
```

**Step 3: Verify completion**

```bash
$ agentic-agent task list

--- DONE ---
[TASK-936281-2] Configure TypeScript and ESLint
  Completed: 2026-02-13T15:45:00Z
  Commits: 3 commits on feature/typescript-config

--- BACKLOG ---
[TASK-936281-3] Set up Tailwind with dark mode
[TASK-936281-4] Create layout components
...

$ agentic-agent openspec status todo-pwa

Change: todo-pwa
  Total: 22  Done: 2  In Progress: 0  Pending: 20
  Progress: 9%
```

**Step 4: Continue to next task**

```
You: /ralph-loop

AI: [Automatically picks up TASK-936281-3 and starts iterating...]
```

---

## Pro Tips

### 1. Configure Checkpoints Before Long Tasks

```yaml
# agnostic-agent.yaml
checkpoint:
  iteration_interval: 3    # More frequent for complex tasks
  token_thresholds: [0.25, 0.5, 0.75, 0.9]
```

### 2. Use Autopilot for Batch Processing

```bash
# Process multiple simple tasks overnight
agentic-agent --agent claude autopilot start \
  --execute-agent \
  --max-iterations 10
```

### 3. Combine ATDD + Ralph Loop

```
You: Before starting, write acceptance tests for TASK-001 based on the criteria.
Then /ralph-loop to implement until tests pass.

AI: [Writes failing tests first, then uses ralph-loop to make them pass]
```

### 4. Monitor Token Usage in Chat

Ralph-loop automatically warns you:
- 80% tokens: ⚠️ Keep going carefully
- 90% tokens: 🛑 Consider pausing and resuming

### 5. Resume After Interruption

```bash
# If autopilot was interrupted (Ctrl+C)
$ agentic-agent --agent claude autopilot start --execute-agent

📌 Resuming from checkpoint (iteration 5, 89000 tokens used)
```

---

## Troubleshooting

### Problem: Ralph-loop skill not found

**Solution:**
```bash
# Ensure skill pack is installed
agentic-agent skills ensure

# Or manually install
cp -r .agentic/skill-packs/ralph-loop ~/.claude/skills/
```

### Problem: Task not claimed before ralph-loop

**Symptom:**
```
$ agentic-agent task continue
Error: No in-progress tasks found
```

**Solution:**
```bash
# Always claim first!
agentic-agent task claim TASK-ID

# Then run ralph-loop
```

### Problem: Checkpoint not resuming

**Check:**
```bash
# Verify checkpoint exists
ls -la .agentic/checkpoints/TASK-*-latest.json

# View checkpoint data
cat .agentic/checkpoints/TASK-001-latest.json | jq .

# If corrupted, remove and restart
rm .agentic/checkpoints/TASK-001-*.json
```

### Problem: Too many iterations without completion

**Symptom:**
```
⚠ Reached max iterations (10)
Task remains in-progress
```

**Solutions:**
1. Check if acceptance criteria are too broad
2. Break task into smaller sub-tasks
3. Increase `--max-iterations` if needed
4. Review checkpoint to see what's blocking

---

## What's Next?

- **Learn about Tracks:** [track-based-workflow](../track-based-workflow/README.md)
- **Explore OpenSpec:** [spec-driven-workflow](../spec-driven-workflow/README.md)
- **Deep dive on Checkpoints:** [checkpoint-resume.md](../../docs/checkpoint-resume.md)
- **Compare approaches:** [autopilot-vs-ralph-loop.md](../../docs/autopilot-vs-ralph-loop.md)
