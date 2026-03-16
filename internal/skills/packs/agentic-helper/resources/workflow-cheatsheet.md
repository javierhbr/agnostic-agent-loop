# Unified SDD Workflow Cheatsheet

The Unified Spec-Driven Development (SDD) umbrella combines multiple methodologies to handle work of any scale. This guide shows how they relate and connect to the CLI.

## The Methodologies

1. **Sdd-Bmad Method:** The overarching framework for routing work based on scale, phase, and role (Analyst, Architect, Developer, Verifier). Use it for progressive planning.
2. **Sdd-Speckit:** The discipline of turning vague ideas into executable, highly structured specifications, plans, and task lists.
3. **Sdd-OpenSpec (op):** The artifact engine. It manages change packages (proposals, delta specs) and enforces strict artifact flow across directories.
4. **Platform Contextualizer:** Used at the start of Iteration 1 to map existing (brownfield) platforms.
5. **Unified SDD Orchestrator:** The mega-skill that combines Sdd-Bmad, Sdd-Speckit, and Sdd-OpenSpec for complex, multi-team platform development.

---

## How to Choose

### Scenario A: Large Brownfield Project
**Scope:** **Platform**
**Trigger:** `unified-sdd`
**Flow:**
1. Run `platform-contextualizer` to document current state of the entire platform.
2. Use `unified-sdd` to generate the overall platform plan and shared context.
3. Route sub-components using `sdd-bmad` to individual teams or tasks.

### Scenario B: New Feature or Service
**Scope:** **Platform/Component**
**Trigger:** `sdd-speckit`
**Flow:**
1. Provide the product brief to `sdd-speckit` to define the feature boundary.
2. Sdd-Speckit generates requirements, a dev plan, and a task list for the component.
3. Move to the CLI execution phase (`task claim`).

### Scenario C: Small Well-Defined Change
**Scope:** **Component**
**Trigger:** `sdd-openspec` (or `op`)
**Flow:**
1. Propose the specific component change with `sdd-openspec`.
2. Generate delta specs for the affected files.
3. Apply the tasks and verify at the component level.

---

## The CLI Bridge (From Plan to Code)

Regardless of the methodology used to generate your plan, the execution is managed by the `agentic-agent` CLI across platform and component levels:

| Step | Command | Level | Why it matters |
|---|---|---|---|
| **1. Set Context** | `agentic-agent context generate <DIR>` | Component | Applies Read-Before-Write rule to ensure safe changes. |
| **2. Track Work** | `task create` (if not generated) | Component | Creates the YAML record for traceability. |
| **3. Begin Code** | `task claim <ID>` | Component | **Mandatory.** Locks the task and traces branch start. |
| **4. Implement** | `run-with-ralph` (or manual) | Component | Iterative code changes and testing. |
| **5. Verify** | `agentic-agent validate` | Platform/Component | Runs gate checks and prevents scope creep. |
| **6. Finish** | `task complete <ID>` | Component | Archives task and captures commits. |