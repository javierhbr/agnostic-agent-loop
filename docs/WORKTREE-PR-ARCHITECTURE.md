# Worktree & PR Architecture Diagram

## System Architecture

```
┌─────────────────────────────────────────────────────────────────────────┐
│                           AGENTIC AGENT CLI                              │
│  (agentic-agent binary with worktree + PR management)                   │
└─────────────────────────────────────────────────────────────────────────┘
                                    │
                  ┌─────────────────┼─────────────────┐
                  │                 │                 │
                  ↓                 ↓                 ↓
        ┌──────────────────┐  ┌──────────────────┐  ┌──────────────────┐
        │  Task Commands   │  │  PR Commands     │  │  Context         │
        ├──────────────────┤  ├──────────────────┤  ├──────────────────┤
        │ • claim          │  │ • create         │  │ • build          │
        │ • complete       │  │ • review         │  │ • generate       │
        │ • list           │  │ • status         │  │                  │
        └────────┬─────────┘  └────────┬─────────┘  └──────────────────┘
                 │                     │
                 ↓                     ↓
        ┌──────────────────────────────────────────────────────────────┐
        │              TASK WORKFLOW LAYER                              │
        │  (worktree.go, lock.go, manager.go)                          │
        ├──────────────────────────────────────────────────────────────┤
        │ • CreateWorktree() - setup isolated env + tests              │
        │ • VerifyGitignore() - auto-add .worktrees to .gitignore      │
        │ • RunProjectSetup() - auto-detect npm/go/python              │
        │ • RunBaselineTests() - REQUIRED to pass                      │
        │ • CleanupWorktree() - auto-delete on completion              │
        │ • CaptureCommits() - record all changes                      │
        └────────────────┬─────────────────────────────────────────────┘
                         │
                         ↓
        ┌──────────────────────────────────────────────────────────────┐
        │              GITHUB INTEGRATION LAYER                         │
        │  (github/cli.go, github/spec.go)                             │
        ├──────────────────────────────────────────────────────────────┤
        │ • CreatePR() - gh pr create with JSON parsing                │
        │ • GetPRInfo() - gh pr view for status                        │
        │ • RequestReview() - request reviewer                         │
        │ • MergePR() - merge strategy (squash/merge/rebase)           │
        │ • BuildPRTitle() - extract from spec                        │
        │ • BuildPRBody() - generate from spec + commits               │
        │ • ExtractPRNumber() - parse PR URL                           │
        └────────────────┬─────────────────────────────────────────────┘
                         │
                         ↓
        ┌──────────────────────────────────────────────────────────────┐
        │                  GITHUB API (via gh CLI)                      │
        │  • gh pr create, view, review, merge                          │
        │  • GitHub OAuth authentication                                │
        └──────────────────────────────────────────────────────────────┘
```

## Worktree Lifecycle

```
┌─────────────────────────────────────────────────────────────────────────┐
│                        TASK CLAIM WORKFLOW                               │
└─────────────────────────────────────────────────────────────────────────┘

1. agentic-agent task claim TASK-123
                │
                ↓
    ┌────────────────────────────────┐
    │ Find task in backlog            │
    └────────────────────────────────┘
                │
                ↓
    ┌────────────────────────────────────────┐
    │ Check: Are we in git repository?       │
    │ - YES → Continue                       │
    │ - NO → Use synthetic path (tests)      │
    └────────────────────────────────────────┘
                │
                ├─→ YES PATH: (Real git repo)
                │       │
                │       ↓
                │   ┌────────────────────────────────┐
                │   │ Verify .worktrees in .gitignore│
                │   │ (auto-add if missing)          │
                │   └────────────────────────────────┘
                │       │
                │       ↓
                │   ┌────────────────────────────────┐
                │   │ git worktree add               │
                │   │ .worktrees/feature/task-123 \  │
                │   │ -b feature/task-123            │
                │   └────────────────────────────────┘
                │       │
                │       ↓
                │   ┌────────────────────────────────┐
                │   │ Auto-detect setup              │
                │   │ • package.json → npm install   │
                │   │ • go.mod → go mod download     │
                │   │ • requirements.txt → pip       │
                │   └────────────────────────────────┘
                │       │
                │       ↓
                │   ┌────────────────────────────────┐
                │   │ Run baseline tests             │
                │   │ MUST PASS                      │
                │   └────────────┬───────────────────┘
                │       │
                │       ├─ PASS ───────┐
                │       │               │
                │       │               ↓
                │       │          ┌────────────────┐
                │       │          │ Task claimed   │
                │       │          │ Worktree ready │
                │       │          └────────────────┘
                │       │
                │       └─ FAIL ────┐
                │                   ↓
                │              ┌─────────────────────┐
                │              │ Worktree cleaned up │
                │              │ Claim failed        │
                │              │ (user can retry)    │
                │              └─────────────────────┘
                │
                └─→ NO PATH: (Test environment)
                        │
                        ↓
                    ┌────────────────────────────────┐
                    │ Return synthetic path          │
                    │ .worktrees/feature/task-123    │
                    │ (not created, no error)        │
                    └────────────────────────────────┘
```

## Task Complete Workflow

```
┌─────────────────────────────────────────────────────────────────────────┐
│         TASK COMPLETE WORKFLOW (Auto-cleanup + Commit capture)           │
└─────────────────────────────────────────────────────────────────────────┘

agentic-agent task complete TASK-123 --learnings "..."
                │
                ↓
    ┌────────────────────────────────┐
    │ Load task from in-progress      │
    └────────────┬───────────────────┘
                 │
                 ↓
    ┌────────────────────────────────────────────────┐
    │ Capture commits from feature/task-123          │
    │ Since: task.ClaimedAt timestamp                │
    │ Stores in: task.Commits[]                      │
    └────────────┬───────────────────────────────────┘
                 │
                 ↓
    ┌────────────────────────────────┐
    │ Record learnings                │
    │ task.Learnings = "..."          │
    └────────────┬───────────────────┘
                 │
                 ↓
    ┌────────────────────────────────────────────────┐
    │ Check if worktree exists                       │
    │ if task.WorktreePath != "" {                   │
    │   git worktree remove <path>                   │
    │ }                                              │
    └────────────┬───────────────────────────────────┘
                 │
                 ↓
    ┌────────────────────────────────┐
    │ Move task to "done" list        │
    │ Set: task.CompletedAt           │
    │ Set: task.Status = "done"       │
    └────────────┬───────────────────┘
                 │
                 ↓
    ✅ Task completed, worktree cleaned up
```

## PR Creation Workflow

```
┌─────────────────────────────────────────────────────────────────────────┐
│           PR CREATION WORKFLOW (Auto-generate from spec)                 │
└─────────────────────────────────────────────────────────────────────────┘

agentic-agent pr create --task TASK-123
                │
                ↓
    ┌────────────────────────────────┐
    │ Load task from in-progress      │
    │ (must have been claimed)        │
    └────────────┬───────────────────┘
                 │
                 ↓
    ┌────────────────────────────────┐
    │ Load spec file                  │
    │ .agentic/specs/<spec-id>        │
    └────────────┬───────────────────┘
                 │
                 ├─ Extract PR title ───────────────┐
                 │  (first # header in spec)        │
                 │                                  │
                 ├─ Extract PR body ────────────────┤
                 │  (summary + AC from spec)        │
                 │                                  │
                 └─ List commits ───────────────────┘
                   (from task.Commits)
                          │
                          ↓
    ┌────────────────────────────────────────────────┐
    │ Invoke: gh pr create                           │
    │  --title <extracted-title>                     │
    │  --body <generated-body>                       │
    │  --base main                                   │
    │  --head feature/task-123                       │
    │  --json url,number,state,title                 │
    └────────────┬───────────────────────────────────┘
                 │
                 ↓
    ┌────────────────────────────────────────────────┐
    │ Parse gh output: {"url": "...", "number": 42} │
    └────────────┬───────────────────────────────────┘
                 │
                 ↓
    ┌────────────────────────────────┐
    │ Record PR metadata in task      │
    │ task.GithubPR.URL = "..."       │
    │ task.GithubPR.Number = 42       │
    │ task.GithubPR.CreatedAt = now   │
    └────────────┬───────────────────┘
                 │
                 ↓
    ✅ PR created and linked to task
    Output: https://github.com/owner/repo/pull/42
```

## PR Review Workflow

```
┌─────────────────────────────────────────────────────────────────────────┐
│           PR REVIEW WORKFLOW (Spawn independent reviewer)                │
└─────────────────────────────────────────────────────────────────────────┘

agentic-agent pr review --task TASK-123 --pr-url https://github.com/.../42
                │
                ↓
    ┌────────────────────────────────┐
    │ Load original task              │
    │ Verify PR URL is valid          │
    │ Extract PR number from URL      │
    └────────────┬───────────────────┘
                 │
                 ↓
    ┌────────────────────────────────────────────────┐
    │ Create review task in backlog                  │
    │ ID: REVIEW-<PR-NUMBER>                         │
    │ Title: "Review PR #42: <original-title>"       │
    │ Type: "review"                                 │
    │ SpecRefs: [same as original]                   │
    │ Scope: [same as original]                      │
    └────────────┬───────────────────────────────────┘
                 │
                 ↓
    ┌────────────────────────────────────────────────┐
    │ Output step-by-step reviewer instructions:     │
    │  1. agentic-agent task claim REVIEW-42         │
    │  2. agentic-agent context build --task         │
    │  3. agentic-agent sdd gate-check <spec>        │
    │  4. agentic-agent validate                     │
    │  5. Review code, score quality                 │
    │  6. agentic-agent task complete REVIEW-42      │
    └────────────┬───────────────────────────────────┘
                 │
                 ↓
    ✅ Reviewer task created in backlog
    Output: REVIEW-42
```

## Multi-Agent Parallel Workflow

```
┌─────────────────────────────────────────────────────────────────────────┐
│          MULTI-AGENT COORDINATION (Parallel development)                 │
└─────────────────────────────────────────────────────────────────────────┘

TIME    DEVELOPER A                  DEVELOPER B            REVIEWER
────────────────────────────────────────────────────────────────────────
10:00   claim TASK-100
        → worktree-100/

10:01                                claim TASK-101
                                     → worktree-101/

10:05   Building in                  Building in
        worktree-100/                worktree-101/
        (main untouched)             (main untouched)

10:20   Commit & push to             (still building)
        feature/task-100

10:45   pr create TASK-100
        → PR #42 created

10:50   pr review --task TASK-100
        → REVIEW-42 created          claim TASK-101
                                     → Complete & cleanup

11:00                                                       claim REVIEW-42
                                                           → worktree-REVIEW-42/

11:15                                                       Complete REVIEW-42
                                                           → verdict announced

─────────────────────────────────────────────────────────────────────────

RESULT:
✅ 3 worktrees used in parallel
✅ Feature branches don't conflict
✅ Main branch stays clean
✅ Independent review context
✅ Full traceability of all changes
```

## Data Flow: Task → Worktree → PR → Review

```
┌──────────────────────────────────────────────────────────────────────────┐
│                                                                            │
│  TASK (in backlog)                                                        │
│  ├─ ID: TASK-123                                                         │
│  ├─ Title: "Implement auth"                                              │
│  ├─ SpecRefs: ["auth-feature"]                                           │
│  └─ Status: pending                                                      │
│         │                                                                 │
│         ↓                                                                 │
│  ┌──────────────────────────────────────────────────────────────────────┐│
│  │ agentic-agent task claim TASK-123                                    ││
│  │                                                                       ││
│  │ TASK (in-progress, with worktree info)                              ││
│  │ ├─ ID: TASK-123                                                      ││
│  │ ├─ Title: "Implement auth"                                           ││
│  │ ├─ SpecRefs: ["auth-feature"]                                        ││
│  │ ├─ Status: in-progress                                               ││
│  │ ├─ ClaimedAt: 2026-03-02T10:00:00Z                                   ││
│  │ ├─ Branch: feature/task-123                                          ││
│  │ └─ WorktreePath: .worktrees/feature/task-123                         ││
│  └──────────────────────────────────────────────────────────────────────┘│
│         │                                                                 │
│         │  Develop in worktree                                            │
│         │  .worktrees/feature/task-123/                                  │
│         │  • git add/commit/push                                          │
│         │  • Commits tracked                                              │
│         │                                                                 │
│         ↓                                                                 │
│  ┌──────────────────────────────────────────────────────────────────────┐│
│  │ agentic-agent pr create --task TASK-123                             ││
│  │                                                                       ││
│  │ TASK (in-progress, with PR info)                                     ││
│  │ ├─ ID: TASK-123                                                      ││
│  │ ├─ Status: in-progress                                               ││
│  │ ├─ Branch: feature/task-123                                          ││
│  │ ├─ Commits: [abc1234, def5678, ...]                                  ││
│  │ ├─ GithubPR:                                                          ││
│  │ │  ├─ URL: https://github.com/org/repo/pull/42                       ││
│  │ │  ├─ Number: 42                                                      ││
│  │ │  └─ CreatedAt: 2026-03-02T10:30:00Z                                ││
│  │ └─ WorktreePath: .worktrees/feature/task-123                         ││
│  └──────────────────────────────────────────────────────────────────────┘│
│         │                                                                 │
│         ↓                                                                 │
│  ┌──────────────────────────────────────────────────────────────────────┐│
│  │ agentic-agent pr review --task TASK-123 --pr-url https://...       ││
│  │                                                                       ││
│  │ REVIEW TASK (in backlog)                                             ││
│  │ ├─ ID: REVIEW-42                                                     ││
│  │ ├─ Title: "Review PR #42: Implement auth"                            ││
│  │ ├─ Type: review                                                      ││
│  │ ├─ SpecRefs: ["auth-feature"]                                        ││
│  │ ├─ Status: pending                                                   ││
│  │ └─ (linked to TASK-123 and PR #42)                                   ││
│  └──────────────────────────────────────────────────────────────────────┘│
│         │                                                                 │
│         ↓                                                                 │
│  ┌──────────────────────────────────────────────────────────────────────┐│
│  │ agentic-agent task claim REVIEW-42 (reviewer)                       ││
│  │                                                                       ││
│  │ REVIEW TASK (in-progress, with own worktree)                         ││
│  │ ├─ ID: REVIEW-42                                                     ││
│  │ ├─ Type: review                                                      ││
│  │ ├─ Status: in-progress                                               ││
│  │ ├─ Branch: feature/task-REVIEW-42                                    ││
│  │ ├─ WorktreePath: .worktrees/feature/task-REVIEW-42                   ││
│  │ └─ (completely independent from builder)                             ││
│  └──────────────────────────────────────────────────────────────────────┘│
│         │                                                                 │
│         │  Reviewer runs quality checks                                   │
│         │  • context build                                                │
│         │  • gate-check                                                   │
│         │  • validate                                                     │
│         │  • score & announce                                             │
│         │                                                                 │
│         ↓                                                                 │
│  ┌──────────────────────────────────────────────────────────────────────┐│
│  │ agentic-agent task complete REVIEW-42 --learnings "..."             ││
│  │                                                                       ││
│  │ REVIEW TASK (done)                                                   ││
│  │ ├─ Status: done                                                      ││
│  │ ├─ CompletedAt: 2026-03-02T11:15:00Z                                 ││
│  │ ├─ Learnings: "Verdict: APPROVE_WITH_CONDITIONS, 1 HIGH risk"       ││
│  │ └─ (worktree cleaned up)                                             ││
│  └──────────────────────────────────────────────────────────────────────┘│
│         │                                                                 │
│         ↓                                                                 │
│  ┌──────────────────────────────────────────────────────────────────────┐│
│  │ agentic-agent task complete TASK-123 --learnings "..."              ││
│  │                                                                       ││
│  │ TASK (done, with full history)                                       ││
│  │ ├─ ID: TASK-123                                                      ││
│  │ ├─ Status: done                                                      ││
│  │ ├─ ClaimedAt: 2026-03-02T10:00:00Z                                   ││
│  │ ├─ CompletedAt: 2026-03-02T11:30:00Z                                 ││
│  │ ├─ Commits: [abc1234, def5678, ...]                                  ││
│  │ ├─ Learnings: "Implemented with feature flags"                       ││
│  │ ├─ GithubPR.URL: https://github.com/org/repo/pull/42                ││
│  │ └─ (worktree cleaned up)                                             ││
│  └──────────────────────────────────────────────────────────────────────┘│
│                                                                            │
│  ✅ COMPLETE LIFECYCLE RECORDED WITH FULL TRACEABILITY                   │
│                                                                            │
└──────────────────────────────────────────────────────────────────────────┘
```

---

## Summary

The architecture enables:

1. **Isolated Development** — Each task in its own worktree
2. **Automatic Lifecycle** — Auto-create on claim, auto-delete on complete
3. **PR Generation** — Auto-generate from spec using `gh` CLI
4. **Independent Review** — Reviewer gets own worktree, own context
5. **Multi-Agent Coordination** — Parallel work, no conflicts, full traceability
6. **Test Safety** — Baseline tests required, synthetic path for non-git environments

All using project-local worktrees, `gh` CLI, and spec-driven automation.
