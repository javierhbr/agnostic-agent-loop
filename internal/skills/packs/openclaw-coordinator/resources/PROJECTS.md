# Project Registry

All known projects that TechLead and ProductLead manage. Both agents read this at session startup to orient themselves.

| ID | Name | Root Path | Stack | Status |
|----|------|-----------|-------|--------|
| proj-001 | agnostic-agent-loop | ~/others/ai-agents/agnostic-agent-loop | Go/Cobra/SDD | active |
| proj-002 | foyer-app | ~/others/metuur/foyer-app | React/TypeScript | active |

## How to Switch Workspace

The `agentic-agent` CLI has no native workspace switching — all paths are cwd-relative. To switch projects:

```bash
# Step 1: Read the registry (you're already in .openclaw/)
cat PROJECTS.md

# Step 2: cd to the target project root
cd ~/path/to/project

# Step 3: Verify orientation by checking project status
agentic-agent status

# Step 4: Update active-project.yaml (optional but recommended)
# Write or update .openclaw/active-project.yaml with:
# - project_id: proj-XXX
# - root: /absolute/path/to/project
# - name: project-name
# - switched_at: [ISO timestamp]
```

## Coordination Across Projects

Both TechLead and ProductLead **share a single announcements file** regardless of which project is active. To keep announcements organized, **always include `project_id` in every announcement entry**.

### Example Announcement (with project_id)

```yaml
announcements:
  - from_agent: product-lead
    to_agent: tech-lead
    project_id: proj-001              # Which project this announcement belongs to
    task_id: spec-auth-login-001
    status: spec-ready
    summary: "Auth login spec approved. Tasks created in backlog."
    data:
      spec_path: .agentic/spec/auth-login/proposal.md
      task_count: 5
      priority: high
```

## Adding a New Project

1. Add a row to this table with a unique ID (proj-NNN format)
2. Set Status to `active` or `archived`
3. Commit PROJECTS.md
4. Both agents will auto-discover the new project on next session startup
