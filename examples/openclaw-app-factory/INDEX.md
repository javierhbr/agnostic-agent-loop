# OpenClaw Examples Index

Complete guides for implementing autonomous agent factories with `agnostic-agent-loop` and OpenClaw.

---

## 📖 Documentation

### For First-Time Users
**Start here:** [QUICKSTART.md](QUICKSTART.md)
- 5-minute setup
- Run your first worker
- Essential commands only

### For Complete Understanding
**Read next:** [README.md](README.md)
- Full scenario: Trailmate app factory
- Architecture and workflows
- Detailed phase-by-phase implementation
- File coordination patterns
- Troubleshooting

### For Implementation Details
**Reference:** `.agentic/skills/openclaw/resources/`
- `orchestrator.md` — Sheldon pattern (full orchestrator playbook)
- `worker.md` — Depth-2 agent (full worker playbook)
- `researcher.md` — Shan pattern (market research)
- `reviewer.md` — QC agent (code review + quality gates)
- `coordination.md` — YAML schema documentation

---

## 🎯 Key Concepts

### The Pattern

```
Main Agent (User)
  │ requests
  ├─→ Orchestrator (reads state, spawns workers)
  │     ├─→ Worker A (implement feature A)
  │     ├─→ Worker B (implement feature B)
  │     └─→ Worker C (implement feature C)
  │
  └─→ (Workers poll kill-signals, announce results)
```

### Coordination (No Central Server)

Three YAML files in `.agentic/coordination/`:
- **reservations.yaml** — File locks (TTL 10 min)
- **announcements.yaml** — Result propagation (child → parent)
- **kill-signals.yaml** — Graceful cascade shutdown

### Key Guarantees

| Guarantee | Mechanism |
|-----------|-----------|
| No concurrent file writes | File reservations with TTL locks |
| Results propagate up | Announce chain (YAML appends) |
| Graceful shutdown | Cascade kill signals + polling |
| Resumable work | Checkpoints with iteration state |
| Clear role boundaries | Orchestrator reads, workers do work |

---

## 🚀 Quick Navigation

### I want to...

**...understand the pattern**
→ Read [QUICKSTART.md](QUICKSTART.md) (5 min) then [README.md](README.md) (Architecture section)

**...run a minimal test**
→ Follow [QUICKSTART.md](QUICKSTART.md) steps 1-5

**...build a real factory**
→ Follow [README.md](README.md) → "Running the Factory" section

**...debug coordination issues**
→ See [README.md](README.md) → "File Coordination (Under the Hood)" and "Troubleshooting"

**...understand file locking**
→ See [README.md](README.md) → "Reservations (Soft Locks)"

**...understand result flow**
→ See [README.md](README.md) → "Announcements (Result Propagation)"

**...add a new phase**
→ See [README.md](README.md) → "Common Patterns" → Pattern 1 (Parallel Workers)

**...learn orchestrator loop**
→ See `.agentic/skills/openclaw/resources/orchestrator.md` → "The Loop"

**...learn worker responsibilities**
→ See `.agentic/skills/openclaw/resources/worker.md` → "Step 1-10"

---

## 📁 File Structure

```
examples/openclaw-app-factory/
├── INDEX.md                     ← You are here
├── QUICKSTART.md                ← 5-minute setup
└── README.md                    ← Full scenario example (800+ lines)

.agentic/skills/openclaw/         ← Installed by: agentic-agent skills install openclaw
├── SKILL.md                      ← Slim trigger
├── AGENT.md                      ← Claude agent identity
├── agents/
│   ├── orchestrator.md           ← .claude/agents/openclaw-orchestrator.md
│   ├── worker.md                 ← .claude/agents/openclaw-worker.md
│   ├── researcher.md             ← .claude/agents/openclaw-researcher.md
│   └── reviewer.md               ← .claude/agents/openclaw-reviewer.md
└── resources/
    ├── orchestrator.md           ← 470-line playbook
    ├── worker.md                 ← 440-line playbook
    ├── researcher.md             ← 300-line playbook
    ├── reviewer.md               ← 380-line playbook
    └── coordination.md           ← 330-line schema guide
```

---

## 🔄 Example Workflow (Phase by Phase)

### Phase 1: Research
```
Task: TASK-1-1 (Research opportunity)
Worker: Researcher (openclaw-researcher agent)
Output: docs/research/pitch.md
Time: 30 min
Announcement: Market score (1-10) + opportunity summary
```

### Phase 2: Validation
```
Task: TASK-2-1 (Validate feasibility)
Worker: Validator (custom agent)
Output: docs/validation/report.md
Time: 20 min
Announcement: Feasibility (go/no-go) + risks
```

### Phase 3: Development
```
Task: TASK-3-1 (Implement core features)
Workers: Builder A, Builder B, Builder C (parallel, 3 file scopes)
Output: src/ (Swift code)
Time: 4 hours (parallel = 1.3 hours actual)
Announcement: Feature complete + AC pass/fail
```

### Phase 4: Review & QC
```
Task: TASK-4-1 (Code review)
Worker: Reviewer (different model than builders)
Output: docs/review-report.md
Time: 1 hour
Announcement: Score (0-10), verdict (APPROVE / REJECT)
```

### Phase 5: Packaging
```
Task: TASK-5-1 (App Store packaging)
Workers: Packager (1 agent for listing + icon)
Output: AppStore/ (metadata, images)
Time: 1 hour
Announcement: App ready for submission
```

### Phase 6: Marketing
```
Task: TASK-6-1 (Demo video + social media)
Workers: Promo Video Agent, Social Media Agent (parallel)
Output: videos/, assets/
Time: 1 hour
Announcement: Marketing complete
```

### Phase 7: Submission
```
Final human approval + manual App Store submit
(Agent prepared everything, human just presses submit)
```

---

## 🛠️ Setup Command Reference

```bash
# One-time setup
agentic-agent init --name "My Factory"
agentic-agent skills install openclaw --tool claude-code

# Create tasks
agentic-agent task create --title "Research" --scope docs --acceptance "One-pager written"

# Claim and work
agentic-agent task claim TASK-123-1
agentic-agent context build --task TASK-123-1

# Spawn workers (via Claude Code)
/subagents spawn openclaw-researcher "..."
/subagents spawn openclaw-worker "..."
/subagents spawn openclaw-orchestrator "..."

# Monitor
watch -n 5 'cat .agentic/coordination/announcements.yaml | tail -20'
agentic-agent task list

# Complete
agentic-agent task complete TASK-123-1
```

---

## 📊 Example Metrics (Trailmate Factory)

| Phase | Duration | Workers | Output Files | AC Pass Rate |
|-------|----------|---------|--------------|--------------|
| Research | 30 min | 1 | pitch.md | 3/3 (100%) |
| Validation | 20 min | 1 | report.md | 2/2 (100%) |
| Development | 1.3 hrs | 3 parallel | 8 .swift files | 12/12 (100%) |
| Review | 1 hr | 1 | review.md | 9/10 score |
| Packaging | 1 hr | 1 | store-assets/ | 4/4 (100%) |
| Marketing | 1 hr | 2 parallel | video.mp4 + assets | complete |
| **Total** | **5 hours** | **9 agents** | **20+ files** | **>95%** |

---

## 🎓 Learning Path

1. **Day 1:** Read QUICKSTART → Run one task with one worker
2. **Day 2:** Read README Phase 3 section → Implement 3 parallel workers
3. **Day 3:** Add reviewer + quality gate (Phase 4)
4. **Day 4:** Run full orchestrator (Phases 1-7)
5. **Day 5:** Deploy on cron for 24/7 autonomous operation

---

## ❓ FAQ

**Q: Do I need an MCP server?**
A: No. Workers and orchestrator manage `.agentic/coordination/` YAML files directly.

**Q: How many workers can run in parallel?**
A: Max 5 per orchestrator (OpenClaw's `maxChildrenPerAgent: 5`). Spawn in waves for more.

**Q: What happens if a worker crashes?**
A: File reservations TTL-expire after 10 min. Orchestrator cascade kill detects timeout after 30 min.

**Q: How does context compaction affect workers?**
A: Worker saves checkpoint before context fills. Orchestrator resumes it from iteration N+1. No work lost.

**Q: Can I run this on my laptop?**
A: Yes. Coordination is purely file-based (no network needed).

**Q: What if I want to pause the factory?**
A: Write kill signal to `.agentic/coordination/kill-signals.yaml`. Workers detect and exit gracefully.

---

## 🔗 Related Documentation

- **OpenClaw pattern:** See `docs/openclaw-pattern.md` (if it exists)
- **Skill system:** See `docs/TIER-MAINTENANCE-GUIDE.md`
- **Task management:** See `agentic-agent task --help`
- **Context bundling:** See `agentic-agent context build --help`

---

**Version:** 1.0
**Last Updated:** 2026-03-01
**Status:** Production-ready
