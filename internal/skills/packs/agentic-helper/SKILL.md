---
name: agentic-helper
description: "Master guide and orchestrator for the agentic-agent CLI and all SDD methodologies. Use when asking about workflow selection, tools, how to use the SDD skills, or how to chain commands."
---

# skill:agentic-helper

## The Unified Workflow Consultant

Guides you through agentic-agent CLI commands and helps you select the correct methodology framework for your task. Whether you are doing a quick fix, starting a new project, or managing a large brownfield platform, this skill directs you to the right tools.

---

## 1. The 3-Pillar Mental Model

Before starting, identify your current role, the artifact you need, and the phase you are in:

| **Pillar** | Key Elements |
|---|---|
| **Roles (The Who)** | Product (Why), Architect (How), Team Lead (Route), Developer (Implement), Verifier (Evidence) |
| **Artifacts (The What)** | PRD, Feature Spec, Proposal, Design, Tasks, Delta Specs, Acceptance Tests |
| **Flow (The When)** | Platform/Route → Specify → Plan → Deliver → Archive |

---

## 2. Methodology Routing Matrix

Before starting, assess your project scale (Platform vs. Component) and goal:

| Your Goal | Scope | Recommended Skill | What it Does |
|---|---|---|---|
| **Multi-team, existing platform?** | **Platform** | `unified-sdd` & `platform-contextualizer` | Platform-scale approach combining Sdd-Bmad, Sdd-OpenSpec, and Sdd-Speckit. |
| **New feature or service?** | **Platform/Component** | `sdd-speckit` | Turns ideas into executable specs, plans, and task lists. |
| **Defined change package?** | **Component** | `sdd-openspec` (or `op`) | Spec-driven changes using proposals and delta specs. |
| **Progressive planning & roles?** | **Platform/Component** | `sdd-bmad` | Routes work by scale, roles (Analyst/Architect/Dev), and artifacts. |
| **Understand architecture?** | **Platform/Component** | `explain-code` | Visual diagrams and analogies for codebase architecture. |
| **Isolated bug fix?** | **Component** | Tiny Workflow (No skill) | Direct execution with basic CLI commands. |

---

## 2. CLI Integration Guide

Methodologies provide the plan; the CLI provides the physical execution across both platform and component levels.

1. **Setup & Context (Platform/Component):**
   - `agentic-agent init` - Initializes the `.agentic` directory for the platform.
   - `agentic-agent context generate <DIR>` - Generates `AGENTS.md` context for a specific component or directory. **Mandatory before editing.**

2. **Execution (Component Level):**
   - `task create` - Converts your component-level plan into an executable task.
   - `task claim <ID>` - **Mandatory.** Locks the task to your branch and starts traceability.
   - Implement using `run-with-ralph` or manual execution.

3. **Finalization:**
   - `agentic-agent validate` - Run gate checks to ensure scope was not violated at the component or platform level.
   - `task complete <ID>` - Captures commits and prepares the component for PR.

---

## 3. Hard Stops (Never Break These)

- Never implement when `blocked_by` is non-empty.
- Never skip `task claim <ID>` (loses traceability).
- Never edit `.agentic/` YAML files directly.
- Never merge without running `agentic-agent validate`.

---

## If you need more detail

→ `resources/workflow-cheatsheet.md` — Visual map of the relationships between Sdd-Bmad, Sdd-OpenSpec, and Sdd-Speckit under the Unified SDD umbrella.
→ `resources/workflow-commands.md` — Full command sequences for traditional tiers (TINY, SMALL, OpenSpec+, Full SDD).