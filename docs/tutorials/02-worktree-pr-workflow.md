# Worktree & PR Workflow - Quick Start Guide

**TL;DR:** Claim task → develop in worktree → create PR from spec → independent review → complete task

---

## 5-Minute Setup

### Install & Build

```bash
go build ./cmd/agentic-agent
export PATH="$PATH:$(pwd)"
```

### Verify gh CLI

```bash
gh --version          # Must be installed
gh auth status        # Must be authenticated to GitHub
```

---

## Step-by-Step Workflow

### 1️⃣ Claim Task (Auto-Creates Worktree)

```bash
agentic-agent task claim TASK-123
```

**Output:**
```
✅ Claimed task TASK-123: Implement feature X

🔧 Isolated worktree created:
   cd .worktrees/feature/task-123

📝 When done, run:
   agentic-agent task complete TASK-123
```

**What happens:**
- ✅ Worktree created at `.worktrees/feature/task-123/`
- ✅ Dependencies auto-installed (npm/go/python detected)
- ✅ Baseline tests run and verified
- ✅ Branch `feature/task-123` created

### 2️⃣ Develop in Worktree

```bash
cd .worktrees/feature/task-123
# Build code, write tests, commit
git add .
git commit -m "feat: implement feature"
git push -u origin feature/task-123
```

**Why this is safe:**
- ✅ Main branch untouched
- ✅ Your feature branch is isolated
- ✅ Multiple agents can work in parallel
- ✅ No merge conflicts with other branches

### 3️⃣ Create PR (Auto-Generate from Spec)

```bash
agentic-agent pr create --task TASK-123
```

**Output:**
```
Creating PR: feature/task-123 → main
✅ PR created: https://github.com/owner/repo/pull/42
```

**What happens:**
- ✅ Title extracted from spec's first `# ` header
- ✅ Body generated from spec sections + commits
- ✅ Acceptance criteria added to PR
- ✅ PR URL saved in task metadata

### 4️⃣ Spawn Reviewer (Independent Review)

```bash
agentic-agent pr review --task TASK-123 \
  --pr-url https://github.com/owner/repo/pull/42
```

**Output:**
```
✅ Review task created: REVIEW-42

🔍 Next steps for reviewer:
  1. agentic-agent task claim REVIEW-42
  2. agentic-agent context build --task REVIEW-42
  3. agentic-agent specifyify gate-check <spec-id>
  4. agentic-agent validate
  5. Review code, score quality
  6. agentic-agent task complete REVIEW-42
```

**What happens:**
- ✅ New task created in backlog: `REVIEW-42`
- ✅ Reviewer can claim and review independently
- ✅ Reviewer gets own worktree
- ✅ No review bias (reviewer didn't build the code)

### 5️⃣ Reviewer Reviews Code

```bash
# Reviewer claims review task
agentic-agent task claim REVIEW-42

# Reviewer loads context
agentic-agent context build --task REVIEW-42

# Reviewer runs quality checks
agentic-agent specifyify gate-check auth-feature
agentic-agent validate

# Reviewer reads code, scores 0-10
# ...review process...

# Reviewer completes with verdict
agentic-agent task complete REVIEW-42 \
  --learnings "Verdict: APPROVE (9/10). Ready to merge."
```

### 6️⃣ Complete Task (Auto-Cleanup)

```bash
agentic-agent task complete TASK-123 \
  --learnings "Implemented with full test coverage. Reviewed: 9/10"
```

**Output:**
```
✅ Captured 3 commits
✅ Completed task TASK-123
📚 Learnings recorded
✅ Worktree removed: .worktrees/feature/task-123
```

**What happens:**
- ✅ All commits captured and linked to task
- ✅ Learnings recorded for knowledge base
- ✅ Worktree auto-deleted (cleanup)
- ✅ Task moved to "done" list

---

## Common Commands Reference

| Command | Purpose |
|---------|---------|
| `agentic-agent task claim TASK-ID` | Start work (creates worktree) |
| `agentic-agent task complete TASK-ID` | Finish work (delete worktree) |
| `agentic-agent pr create --task TASK-ID` | Generate PR from spec |
| `agentic-agent pr review --task ID --pr-url URL` | Spawn reviewer task |
| `agentic-agent pr status --pr-url URL` | Check PR status |
| `agentic-agent context build --task TASK-ID` | Load context for review |
| `agentic-agent specifyify gate-check SPEC-ID` | Run quality gates |
| `agentic-agent validate` | Run validation rules |
| `agentic-agent task list` | See all tasks |

---

## Troubleshooting

### ❌ Worktree creation failed: "Not a git repository"

**Cause:** You're not in a git repository
**Solution:** Initialize git or use the code in a git repo
```bash
git init
git add .
git commit -m "initial"
```

### ❌ Baseline tests failed

**Cause:** Tests don't pass in the current environment
**Solution:** Fix tests, then retry claiming
```bash
agentic-agent task claim TASK-ID
```
The worktree will be auto-cleaned if tests fail.

### ❌ PR creation failed: "gh command not found"

**Cause:** GitHub CLI not installed
**Solution:** Install `gh` CLI
```bash
brew install gh  # macOS
apt-get install gh  # Linux
# or visit: https://cli.github.com/
```

### ❌ PR creation failed: "not authenticated"

**Cause:** GitHub CLI not authenticated
**Solution:** Authenticate with GitHub
```bash
gh auth login
# Follow prompts to authenticate
```

### ❌ Worktree path shows but no tests ran

**Cause:** Probably in a non-git environment (tests)
**Solution:** This is OK - you get a synthetic path that works anyway

---

## Tips & Tricks

### Tip 1: View Task Details

```bash
agentic-agent task show TASK-ID
# Shows: status, worktree path, commits, PR URL, etc.
```

### Tip 2: Continue Work After Interruption

```bash
agentic-agent task claim TASK-ID --continue
# Reuses same worktree if still exists
```

### Tip 3: Record Learnings for Knowledge Base

```bash
agentic-agent task complete TASK-ID \
  --learnings "Pattern: Always use mutex for concurrent access. Gotcha: Need defer unlock()."
```

### Tip 4: Check PR Before Completing

```bash
agentic-agent pr status --pr-url https://github.com/owner/repo/pull/42
# Verify PR is merged before completing task
```

### Tip 5: Parallel Development

Multiple agents can claim different tasks:
```bash
# Agent A
agentic-agent task claim TASK-100

# Agent B (at same time)
agentic-agent task claim TASK-101

# Both work in parallel: .worktrees/feature/task-100, task-101
# Zero conflicts!
```

---

## Multi-Agent Example

**Scenario:** 2 developers, 1 reviewer, working in parallel

```bash
# Time 10:00
Developer A: agentic-agent task claim TASK-100
            → Worktree A created

Developer B: agentic-agent task claim TASK-101
            → Worktree B created

# Time 10:00-11:00
# Both develop in parallel
A: cd .worktrees/feature/task-100 && npm test && git push
B: cd .worktrees/feature/task-101 && npm test && git push

# Time 11:00
A: agentic-agent pr create --task TASK-100
   → PR #50 created

B: agentic-agent pr create --task TASK-101
   → PR #51 created

# Time 11:05
Orchestrator: agentic-agent pr review --task TASK-100 --pr-url ...
             agentic-agent pr review --task TASK-101 --pr-url ...
             → REVIEW-50, REVIEW-51 created

# Time 11:10
Reviewer: agentic-agent task claim REVIEW-50
         agentic-agent task claim REVIEW-51
         → Both reviewed in parallel

# Time 11:30
# All done, both PRs reviewed and approved!
```

**Result:** 2 features shipped in parallel with independent review, zero conflicts.

---

## Full Documentation

For detailed information, see:

- **[WORKTREE-AND-PR-IMPLEMENTATION.md](WORKTREE-AND-PR-IMPLEMENTATION.md)** — Complete reference
- **[WORKTREE-PR-ARCHITECTURE.md](WORKTREE-PR-ARCHITECTURE.md)** — System design & diagrams
- **[WORKTREE-PR-EXAMPLES.md](WORKTREE-PR-EXAMPLES.md)** — Real-world scenarios

---

## Next Steps

1. **Clone or cd into the agentic-agent repo**
   ```bash
   cd /path/to/agnostic-agent-loop
   ```

2. **Build the binary**
   ```bash
   go build ./cmd/agentic-agent
   ```

3. **Create a test task**
   ```bash
   agentic-agent task create --title "Test feature"
   ```

4. **Claim and develop**
   ```bash
   agentic-agent task claim TASK-ID
   cd .worktrees/feature/task-ID
   ```

5. **Try the full workflow**
   ```bash
   agentic-agent pr create --task TASK-ID
   agentic-agent pr review --task TASK-ID --pr-url <URL>
   agentic-agent task complete TASK-ID
   ```

---

**You're ready to go! Happy developing! 🚀**
