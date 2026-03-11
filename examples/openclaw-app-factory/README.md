# OpenClaw Autonomous App Factory Example

A complete, production-grade example of using `agnostic-agent-loop` with the OpenClaw orchestrator pattern to build apps autonomously **across multiple projects**.

## Scenario: Trailmate (Offline Fitness Tracker)

An autonomous agent factory that researches app opportunities, validates ideas, builds apps, reviews code, and ships to the App Store—all without human intervention except the final "submit" button. The factory also manages **multiple app projects** via TechLead and ProductLead coordinators.

---

## Architecture: Single-Project Classic (with Multi-Project Extension)

### Option 1: Single Project (Trailmate only)

```
Main Agent (Human Interface)
  │
  └─→ Orchestrator
        ├─→ Worker: Builder
        ├─→ Worker: Reviewer
        ├─→ Worker: QC
        └─→ Worker: Packager
```

### Option 2: Multi-Project (NEW - using TechLead & ProductLead)

```
ProductLead (Product Coordinator)              TechLead (Dev Team Coordinator)
  (Define specs across projects)                (Build specs across projects)
  │                                             │
  ├─ Read: .openclaw/PROJECTS.md          ├─ Read: .openclaw/PROJECTS.md
  ├─ Activate: proj-001, proj-002...      ├─ Activate: proj-001, proj-002...
  ├─ Create specs + announce              ├─ Claim tasks + spawn workers
  └─→ Include project_id in each          └─→ Filter announcements by project_id
       announcement                             │
       │                                        └─→ .agentic/coordination/announcements.yaml
       └─────────────────────────────────────────→ (shared YAML, project_id-tagged)
```

**Flow (both options):**
1. **Researcher** scans for app opportunities (5-min cron)
2. **Validator** checks feasibility
3. **Builder** codes the app (Swift + SwiftUI)
4. **Reviewer** verifies code quality
5. **QADev** runs automated checks with 10-point rubric (must score **≥8/10 to pass**)
   - AC Coverage, Unit Tests, Integration, E2E/Visual, Edge Cases, Performance, Security, Accessibility, Regression, Documentation
   - Score <8 → requests fixes back to Builder
   - Score ≥8 → approves and passes to next stage
6. **Packager** generates App Store listing
7. **Onboarding Agent** creates user education screens
8. **Promo Video Agent** generates demo video
9. **Submit** to App Store (human approval only)

---

## Setup

### Option A: Single-Project Setup (Trailmate only)

#### 1. Initialize Project

```bash
mkdir trailmate-factory && cd trailmate-factory
agentic-agent init --name "Trailmate Factory" --no-interactive

# Install OpenClaw skill pack
agentic-agent skills install openclaw --tool claude-code --no-interactive
```

#### 2. Create Project State File

The orchestrator reads this to decide what phase to run next.

```bash
cat > .agentic/project-state.md << 'EOF'
# Trailmate Project State

## Current Phase: Research

- [ ] Phase 1: Research & Opportunity
  - [ ] Scan platforms for pain points
  - [ ] Identify market opportunity
  - [ ] Write one-pager pitch

- [ ] Phase 2: Validation
  - [ ] Check technical feasibility
  - [ ] Validate business model

- [ ] Phase 3: Development
  - [ ] Design architecture
  - [ ] Implement core features
  - [ ] Add observability

- [ ] Phase 4: Review & QC
  - [ ] Code review (8+ score)
  - [ ] Quality checks (all pass)

- [ ] Phase 5: Packaging
  - [ ] App Store listing
  - [ ] Screenshots & icon
  - [ ] Onboarding screens

- [ ] Phase 6: Marketing
  - [ ] Demo video
  - [ ] Social media content

- [ ] Phase 7: Submission
  - [ ] Submit to App Store
EOF
```

### 3. Create Initial Task (Research)

```bash
agentic-agent task create \
  --title "Research: Offline Fitness Tracker Opportunity" \
  --description "Scan platforms for pain points in fitness tracking. Identify market gap for offline-first GPS fitness app targeting trail runners." \
  --scope docs/research \
  --outputs "docs/research/trailmate-pitch.md" \
  --acceptance "One-pager written", "Market demand > 5K searches/month", "Competition assessed"
```

---

### Option B: Multi-Project Setup (Trailmate + FitHub + more)

#### 1. Create Shared OpenClaw Directory (Outside Projects)

```bash
# Create a shared coordinator space
mkdir -p ~/app-factory/.openclaw
cd ~/app-factory/.openclaw
```

#### 2. Create Project Registry

```bash
cat > PROJECTS.md << 'EOF'
# Project Registry

Both ProductLead and TechLead use this to manage multiple app projects.

| ID | Name | Root Path | Stack | Status |
|----|------|-----------|-------|--------|
| proj-trailmate | Trailmate | ~/app-factory/trailmate | Swift/SwiftUI | active |
| proj-fithub | FitHub | ~/app-factory/fithub | Swift/SwiftUI | active |
| proj-fitjournal | FitJournal | ~/app-factory/fitjournal | React/TypeScript | active |

## Quick Switch

```bash
cd ~/app-factory/trailmate && agentic-agent status     # Switch to Trailmate
cd ~/app-factory/fithub && agentic-agent status        # Switch to FitHub
```

See `.openclaw/active-project.yaml` to track current workspace.
EOF
```

#### 3. Create Shared User Context

```bash
cat > USER.md << 'EOF'
# User

**Name:** Your Name
**Role:** Product + Engineering Lead
**Timezone:** UTC

## Context

Managing an autonomous app factory that builds fitness apps across multiple platforms.
Both ProductLead (product coordinator) and TechLead (dev team coordinator) use this workspace.

## Projects

- **Trailmate:** Offline GPS fitness tracker (iOS, Go backend)
- **FitHub:** Social fitness app (iOS + Android, Node.js)
- **FitJournal:** Fitness logging + analytics (Web, Python)
EOF
```

#### 4. Create Per-Project Directories

```bash
mkdir -p ~/app-factory/{trailmate,fithub,fitjournal}
cd ~/app-factory/trailmate

# Initialize each project
agentic-agent init --name "Trailmate" --no-interactive

# Install OpenClaw skill pack (optional, once at root)
# agentic-agent skills install openclaw --tool claude-code --no-interactive
```

Repeat for `fithub/` and `fitjournal/`.

#### 5. Initialize ProductLead

```bash
cd ~/app-factory/.openclaw

# ProductLead will read PROJECTS.md and manage specs across all projects
# She will announce spec-ready with project_id included
cat > product-lead-init.md << 'EOF'
# ProductLead Session

## Load Order

1. Read PROJECTS.md — discover Trailmate, FitHub, FitJournal
2. Check active-project.yaml — which project is active?
3. If none, ask human: "Which project should I focus on?"
4. cd ~/app-factory/<project-name>
5. agentic-agent status

## Example: Create Trailmate Spec

```bash
cd ~/app-factory/trailmate
agentic-agent openspec init "GPS Recording Feature" --risk high
```

## Example: Create FitHub Spec

```bash
cd ~/app-factory/fithub
agentic-agent openspec init "Social Sharing" --risk medium
```

## Announcement with project_id

```yaml
announcements:
  - from_agent: product-lead
    to_agent: tech-lead
    project_id: proj-trailmate        # ← Specify which project
    status: spec-ready
    summary: "GPS Recording approved. 5 tasks in backlog."
    data:
      spec_path: .agentic/spec/gps-recording/proposal.md
      task_count: 5
      priority: high
```

TechLead will see this, cd to trailmate project, and start work.
EOF
```

#### 6. Initialize TechLead

```bash
cd ~/app-factory/.openclaw

# TechLead will read PROJECTS.md and manage development across all projects
# He will filter announcements by project_id
cat > tech-lead-init.md << 'EOF'
# TechLead Session

## Load Order

1. Read PROJECTS.md — discover all projects
2. Check active-project.yaml — which project am I in?
3. cd ~/app-factory/<project-name>
4. agentic-agent status

## Filter Announcements by project_id

When I read announcements.yaml, I only respond to entries where:
- `project_id: proj-trailmate` (if currently in trailmate/)
- `from_agent: product-lead`
- `status: spec-ready`

Example:
```yaml
# ✓ I respond to this (I'm in proj-trailmate)
- from_agent: product-lead
  to_agent: tech-lead
  project_id: proj-trailmate
  status: spec-ready

# ✗ I ignore this (I'm in proj-trailmate, not proj-fithub)
- from_agent: product-lead
  to_agent: tech-lead
  project_id: proj-fithub
  status: spec-ready
```

## Workflow

When I see a spec-ready announcement for my active project:

```bash
# 1. Verify current project
agentic-agent status

# 2. List tasks for that spec
agentic-agent task list

# 3. Claim a task
agentic-agent task claim TASK-ID

# 4. Build context
agentic-agent context build --task TASK-ID

# 5. Spawn workers (same as single-project)
/subagents spawn openclaw-builder "Task: TASK-ID..."

# 6. Announce completion with project_id
announcements.yaml:
  - from_agent: tech-lead
    to_agent: product-lead
    project_id: proj-trailmate
    status: complete
    summary: "GPS spec shipped. 4 tasks done."
```
EOF
```

---

## Running the Factory

### Single-Project Mode

#### Manual Mode (For Testing)

```bash
# Step 1: Research
agentic-agent task list
# Pick TASK-123456-1 (research task)

agentic-agent task claim TASK-123456-1
agentic-agent context build --task TASK-123456-1

# Spawn Researcher agent
/subagents spawn openclaw-researcher "
Task: TASK-123456-1
Scan Reddit, X, App Store for fitness app complaints.
Find niche market opportunity.
Write one-pager: problem, audience, solution, market insight, revenue model.
When done: <promise>COMPLETE</promise>
"

# Wait for researcher to finish
# Then complete the task
agentic-agent task complete TASK-123456-1
```

#### Orchestrator Mode (Fully Autonomous)

```bash
# Spawn the orchestrator which manages everything
/subagents spawn openclaw-orchestrator "
You are Sheldon, the orchestrator for Trailmate Factory.

Read: .agentic/project-state.md

1. Check current phase
2. List backlog tasks for that phase
3. For each task:
   - Validate gates: agentic-agent specify gate-check <spec>
   - Claim: agentic-agent task claim <ID>
   - Build context: agentic-agent context build --task <ID>
   - Spawn 2-4 workers (Builder, Reviewer, QC, Packager)
   - Poll .agentic/coordination/announcements.yaml every 10s
   - Complete task when all workers announce
4. Update .agentic/project-state.md to next phase
5. Repeat for next phase

Max 30 iterations. Stop when Phase 7 (Submission) complete.
Output: <promise>COMPLETE</promise>
"
```

---

## Detailed Example Flow: Phase 3 (Development)

### Task Setup

```bash
cat > .agentic/tasks/backlog-phase3.yaml << 'EOF'
tasks:
  - id: TASK-500-1
    title: "Implement: GPS Recording & Offline Storage"
    description: "Record GPS in real-time, store offline, sync when online"
    scope:
      - internal/tracking/
      - internal/storage/
    spec_refs:
      - .agentic/spec/trailmate-feature/component-spec.md
    acceptance:
      - "GPS recording works offline"
      - "Data persists after app restart"
      - "Sync works when online"
      - "Battery efficiency < 15% per hour"
    skill_refs:
      - sdd/developer
      - tdd
EOF
```

### Orchestrator Spawns Workers

```bash
# Worker A: Implements GPS module
/subagents spawn openclaw-worker "
Task: TASK-500-1
Scope: internal/tracking/gps.go

Read spec: .agentic/spec/trailmate-feature/component-spec.md

Implement:
- GPSRecorder class
- 1-second GPS polling
- Offline queue management
- Sync on connectivity

AC must pass:
- GPS records without crashing
- Survives app restart
- Syncs when online
- <15% battery per hour (measure with Instruments)

Reserve files:
- internal/tracking/gps.go
- internal/tracking/gps_test.go

When complete: <promise>COMPLETE</promise>
"

# Worker B: Implements storage + sync
/subagents spawn openclaw-worker "
Task: TASK-500-1
Scope: internal/storage/

Read spec: .agentic/spec/trailmate-feature/component-spec.md

Implement:
- GPXDataStore (CoreData + realm)
- Offline queue (unsent records)
- CloudKit sync logic
- Conflict resolution

AC must pass:
- Data persists after restart
- CloudKit sync works
- Handles network errors gracefully

Reserve files:
- internal/storage/datastore.go
- internal/storage/sync.go

When complete: <promise>COMPLETE</promise>
"

# Orchestrator polls announcements
while true; do
  announcements=$(cat .agentic/coordination/announcements.yaml 2>/dev/null || echo "")
  worker_a=$(echo "$announcements" | grep -c "from_agent: worker-a.*status: complete")
  worker_b=$(echo "$announcements" | grep -c "from_agent: worker-b.*status: complete")

  if [ "$worker_a" -eq 1 ] && [ "$worker_b" -eq 1 ]; then
    # Both done
    agentic-agent task complete TASK-500-1
    break
  fi

  sleep 10
done
```

### Workers Coordinate via YAML

**Worker A announces:**
```yaml
# .agentic/coordination/announcements.yaml
announcements:
  - from_agent: worker-a
    to_agent: orchestrator
    task_id: TASK-500-1
    status: complete
    summary: "GPS recorder: 4/4 AC pass. Battery: 12% per hour. Ready for integration."
    data:
      files_changed: ["internal/tracking/gps.go", "internal/tracking/gps_test.go"]
      iterations: 3
      learnings: ["CLLocationManager needs background mode", "Use Timer at 1Hz for smoothing"]
    timestamp: "2026-03-01T14:30:00Z"
```

**Worker B announces:**
```yaml
announcements:
  - from_agent: worker-b
    to_agent: orchestrator
    task_id: TASK-500-1
    status: complete
    summary: "Storage + sync: 4/4 AC pass. CloudKit integration tested. Error handling robust."
    data:
      files_changed: ["internal/storage/datastore.go", "internal/storage/sync.go"]
      iterations: 2
      learnings: ["CloudKit CKModifyRecordsOperation for batch sync", "Conflict resolution: last-write-wins"]
    timestamp: "2026-03-01T14:35:00Z"
```

**Orchestrator synthesizes:**
```yaml
announcements:
  - from_agent: orchestrator
    to_agent: main
    task_id: TASK-500-1
    status: complete
    summary: "GPS + Sync complete. 2 workers, 4 files changed, all AC pass. Ready for QC phase."
    data:
      workers_completed: 2
      total_files_changed: 4
      time_elapsed_minutes: 35
      phase_complete: "Development (Phase 3)"
      next_phase: "Review & QC (Phase 4)"
    timestamp: "2026-03-01T14:36:00Z"
```

---

## File Coordination (Under the Hood)

### Reservations (Soft Locks)

Workers reserve files before editing. TTL 10 minutes auto-expires if a worker crashes.

```yaml
# .agentic/coordination/reservations.yaml
reservations:
  - reservation_id: res-worker-a-001
    file_path: internal/tracking/gps.go
    owner: worker-a
    task_id: TASK-500-1
    created_at: "2026-03-01T14:00:00Z"
    expires_at: "2026-03-01T14:10:00Z"

  - reservation_id: res-worker-b-001
    file_path: internal/storage/datastore.go
    owner: worker-b
    task_id: TASK-500-1
    created_at: "2026-03-01T14:00:00Z"
    expires_at: "2026-03-01T14:10:00Z"
```

If a third worker tries to edit `gps.go`:
```
ERROR: Cannot reserve internal/tracking/gps.go
Held by worker-a until 2026-03-01T14:10:00Z
```

### Kill Signals (Cascade Stop)

If orchestrator times out (no progress in 30 minutes):

```yaml
# .agentic/coordination/kill-signals.yaml
kill_signals:
  - signal_id: kill-001
    orchestrator_id: main
    target: all
    reason: "Task timeout: TASK-500-1 workers stuck for 30 min"
    active: true
    issued_at: "2026-03-01T15:00:00Z"
```

Each worker checks this at iteration start:
```
[Worker A] Checking kill signal...
[Worker A] Kill signal detected! Orchestrator timeout.
[Worker A] Releasing files...
[Worker A] Announcing failed status...
[Worker A] Exiting.
<promise>COMPLETE</promise>
```

---

## Phase 4: Review & QC

### Reviewer Workflow

```bash
/subagents spawn openclaw-reviewer "
Task: TASK-500-1

Read spec: .agentic/spec/trailmate-feature/component-spec.md
Read code: internal/tracking/gps.go, internal/storage/datastore.go

Verify:
1. All acceptance criteria met (test results)
2. No crash risks (nil checks, error handling)
3. Edge cases handled (app killed, network down, user interrupts)
4. Observability (logging, metrics, tracing)
5. Gates pass: agentic-agent specify gate-check trailmate-feature

Score quality (0-10):
- Completeness: all sections implemented?
- Correctness: algorithms match spec?
- Safety: error handling robust?

If score >= 8: APPROVE
If score 7: APPROVE_WITH_CONDITIONS (request fixes)
If score <= 6: REJECT (major rework)

Announce: verdict + score + findings
When done: <promise>COMPLETE</promise>
"
```

### Reviewer Output

```yaml
# .agentic/coordination/announcements.yaml
announcements:
  - from_agent: reviewer-codex
    to_agent: orchestrator
    task_id: TASK-500-1
    status: complete
    summary: "Code review: 9/10. Excellent implementation. Minor logging improvement suggested."
    data:
      score: 9
      verdict: "APPROVE"
      gates_passed: "5/5"
      critical_risks: 0
      high_risks: 0
      medium_risks: 1
      findings: "Add request ID to sync logs for tracing"
    timestamp: "2026-03-01T15:30:00Z"
```

If score < 8, code goes back to workers for fixes:
```
verdict: "APPROVE_WITH_CONDITIONS"
message: "Fix: request ID logging (1 HIGH risk). Resubmit for re-review."
```

---

## Running Continuously (24/7)

### Cron Job

```bash
# Start orchestrator every hour
0 * * * * cd /path/to/trailmate-factory && \
  /subagents spawn openclaw-orchestrator \
  "Read .agentic/project-state.md. Run next phase. Update state. Output <promise>COMPLETE</promise>"
```

### Monitoring Dashboard

```bash
watch -n 30 'echo "=== Task Status ===" && \
  agentic-agent task list && \
  echo "=== Coordination Status ===" && \
  cat .agentic/coordination/announcements.yaml | tail -20 && \
  echo "=== Kill Signals ===" && \
  cat .agentic/coordination/kill-signals.yaml | tail -5'
```

---

## Common Patterns

### Pattern 1: Parallel Workers (No Conflict)

```
Phase: UI Implementation
├─ Worker A: HomeScreen.swift (scope: screens/home/)
├─ Worker B: DetailScreen.swift (scope: screens/detail/)
└─ Worker C: SettingsScreen.swift (scope: screens/settings/)

No file conflicts → all 3 can run in parallel
Orchestrator waits for all 3 announcements → task complete
```

### Pattern 2: Sequential Workers (With Handoff)

```
Phase: Backend API
├─ Worker A: Auth endpoints (reserves internal/api/auth.go)
├─ Worker B: Data endpoints (wait for A complete → use A's auth)
└─ Worker C: Admin endpoints (wait for A+B → integrate both)

Worker B: polls kill-signals + checks if A announced complete
If A announced: proceed with implementation
If kill-signal: exit gracefully
```

### Pattern 3: Reviewer + Builder Loop

```
Phase: Code Quality Gate
┌─ Builder: implements feature
├─ Reviewer: runs gate-check, score < 8
├─ Builder: fixes issues (respawn)
└─ Reviewer: re-reviews (repeat until score >= 8)
```

---

## Multi-Project Mode (NEW)

### Running ProductLead Across Projects

ProductLead starts at `~/app-factory/.openclaw/` and manages specs for all projects:

```bash
cd ~/app-factory/.openclaw

/subagents spawn product-lead "
You are ProductLead.

Load order:
1. Read PROJECTS.md — find Trailmate, FitHub, FitJournal
2. Check active-project.yaml — which project are you in?
3. For each active project:
   - cd ~/app-factory/<project-name>
   - agentic-agent status
   - Review .agentic/spec/ for pending specs
   - Create/approve specs as needed
   - Announce to TechLead with project_id

Example announcement:
  from_agent: product-lead
  to_agent: tech-lead
  project_id: proj-trailmate
  status: spec-ready
  summary: 'GPS Recording: approved. 5 tasks.'

When done: Output <promise>COMPLETE</promise>
"
```

### Running TechLead Across Projects

TechLead starts at `~/app-factory/.openclaw/` and manages development for all projects:

```bash
cd ~/app-factory/.openclaw

/subagents spawn tech-lead "
You are TechLead.

Load order:
1. Read PROJECTS.md — find all projects
2. Check active-project.yaml — which project are you in?
3. Check .agentic/coordination/announcements.yaml
4. For each spec-ready announcement with project_id:
   - cd ~/app-factory/<project-id-name>
   - agentic-agent task list
   - Claim tasks for that spec
   - Spawn workers (same as single-project)
   - Announce completion with project_id

Example: Seeing proj-trailmate spec-ready:
  cd ~/app-factory/trailmate
  agentic-agent task claim TASK-500-1
  /subagents spawn worker 'Task: TASK-500-1...'

When done: Output <promise>COMPLETE</promise>
"
```

### Full Multi-Project Orchestration

Both agents work in parallel across projects:

```
Timeline:
─────────────────────────────────────────────────────────

T=0m: ProductLead activates
  → cd trailmate, create GPS spec, announce proj-trailmate

T=1m: TechLead receives announcement
  → cd trailmate, spawn 2 builders, wait for workers

T=5m: ProductLead activates for FitHub
  → cd fithub, create Social Sharing spec, announce proj-fithub

T=6m: TechLead finishes Trailmate workers
  → Read announcements, filter by proj-trailmate ✓, complete task
  → Read announcements for proj-fithub, spawn workers

T=20m: Both projects shipping in parallel
  → Trailmate in QC phase, FitHub in Development phase
```

### Multi-Project Announcement Format

Always include `project_id` so agents know which workspace to cd to:

```yaml
# .agentic/coordination/announcements.yaml (shared file, all projects)

announcements:
  # Trailmate: GPS spec approved
  - from_agent: product-lead
    to_agent: tech-lead
    project_id: proj-trailmate
    task_id: spec-gps-001
    status: spec-ready
    summary: "GPS Recording approved. 5 tasks in backlog."
    data:
      spec_path: .agentic/spec/gps/proposal.md
      task_count: 5
      priority: high
    timestamp: "2026-03-01T10:00:00Z"

  # FitHub: Social Sharing spec approved (different project)
  - from_agent: product-lead
    to_agent: tech-lead
    project_id: proj-fithub
    task_id: spec-social-001
    status: spec-ready
    summary: "Social Sharing approved. 3 tasks in backlog."
    data:
      spec_path: .agentic/spec/social/proposal.md
      task_count: 3
      priority: medium
    timestamp: "2026-03-01T10:05:00Z"

  # Trailmate: GPS builders finished
  - from_agent: tech-lead
    to_agent: product-lead
    project_id: proj-trailmate
    task_id: spec-gps-001
    status: complete
    summary: "GPS spec shipped. 5 tasks done. Ready for QC."
    data:
      files_changed: 12
      tests_pass: true
      ready_for_review: true
    timestamp: "2026-03-01T10:30:00Z"
```

### Monitoring Multi-Project

```bash
cd ~/app-factory

# See all projects
cat .openclaw/PROJECTS.md

# Monitor active project
cat .openclaw/active-project.yaml

# Watch announcements (all projects mixed)
watch -n 5 'tail -20 .agentic/coordination/announcements.yaml'

# Check which project each task belongs to
for dir in trailmate fithub fitjournal; do
  echo "=== $dir ==="
  cd $dir && agentic-agent task list
  cd ..
done
```

---

## Troubleshooting

| Problem | Cause | Solution |
|---------|-------|----------|
| Worker hangs indefinitely | Network issue, infinite loop | Orchestrator cascade kill after 30 min |
| Two workers edit same file | Scope overlap | Design task scope to avoid conflicts |
| Reviewer score < 8 three times | Code quality issues | Escalate to human expert, don't respawn |
| Announcements not found | Workers crashed before announcing | Check `.agentic/coordination/kill-signals.yaml` for reasons |
| Context compaction mid-task | Token limit reached | Worker saves checkpoint, orchestrator resumes from iteration N+1 |

---

## Next Steps

### Single-Project Path
1. **Study the playbooks:** Read `internal/skills/packs/openclaw/resources/*.md`
2. **Try manual mode:** Claim a task, spawn one worker, verify announcement flow
3. **Move to orchestrator:** Spawn `openclaw-orchestrator`, let it manage the full pipeline
4. **Monitor coordination:** Watch `.agentic/coordination/` files as workers coordinate
5. **Scale up:** Add more workers, more phases, more agents

### Multi-Project Path
1. **Create project registry:** Set up `PROJECTS.md` with all your apps
2. **Try single project first:** Get one app working end-to-end
3. **Activate ProductLead:** Have her create specs for one project
4. **Activate TechLead:** Have him build that spec, verify announcements include project_id
5. **Add second project:** Repeat for another app, watch both agents coordinate
6. **Scale to N projects:** Let both agents manage 3+ apps in parallel

---

## Developer Agents (from openclaw-coordinator skill pack)

For full implementation workflows, the factory can use the complete developer agent ecosystem:

### The 6-Agent System

**Coordinators:**
- **TechLead** — Orchestrates work, spawns developers by layer, enforces QA gates
- **ProductLead** — Defines specs, owns API contracts

**Developer Workers:**
- **BackendDev** — API/database implementation
- **FrontendDev** — Web UI with visual testing (agent-browser)
- **MobileDev** — Flutter cross-platform (Dart/Flutter MCP server)
- **QADev** — Quality gatekeeper with 10-point scoring rubric

### Installation

```bash
# Install the complete openclaw-coordinator pack
agentic-agent skills install openclaw-coordinator --tool claude-code

# All 6 agents install to .claude/agents/
# - openclaw-tech-lead.md, openclaw-product-lead.md (coordinators)
# - openclaw-backend-dev.md, openclaw-frontend-dev.md, openclaw-mobile-dev.md, openclaw-qa-dev.md (workers)
```

### QA Scoring Rubric (used by QADev)

All code must score ≥8/10 before shipping:

| # | Criterion | What QADev checks |
|----|-----------|------------------|
| 1 | AC Coverage | All acceptance criteria have tests |
| 2 | Unit Tests | ≥80% line coverage |
| 3 | Integration | API contracts match spec, no deviations |
| 4 | E2E/Visual | Happy-path E2E test with screenshots |
| 5 | Edge Cases | ≥3 error conditions tested |
| 6 | Performance | No regression in response time/bundle |
| 7 | Security | Auth, validation, injection prevention |
| 8 | Accessibility | ARIA labels, keyboard nav, contrast |
| 9 | Regression | All existing tests still pass |
| 10 | Documentation | Test intent clearly explained |

**Score <8:** QADev requests fixes back to developer (builder, frontend, mobile)
**Score ≥8:** QADev approves and TechLead announces complete to ProductLead

---

## Files Reference

- Task definitions: `.agentic/tasks/{backlog,in-progress,done}.yaml`
- Project state: `.agentic/project-state.md` (orchestrator reads this)
- Coordination: `.agentic/coordination/{reservations,announcements,kill-signals}.yaml`
- Specs: `.agentic/spec/trailmate-feature/`
- Checkpoints: `.agentic/checkpoints/{TASK-ID}-{iteration}.json`
- Context: `.agentic/context/tech-stack.md`, `global-context.md`

---

## Real-World Adaptations

- **Change app type:** Modify specs to target Android (Kotlin), web (React), or desktop
- **Change team size:** Start with 2 workers, scale to 5 (OpenClaw max), then spawn in waves
- **Add approval gates:** Add human checkpoints before shipping (gate-check, review approval)
- **Add monitoring:** Integrate logs into Datadog/CloudWatch for alerts on worker crashes
- **Add rollback:** Create a `rollback-to-previous-build.sh` script for failed deployments
