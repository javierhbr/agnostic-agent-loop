# Agentic Skill Packs Overview

This directory contains the core intelligence of the Agnostic Agent Loop. These "Skill Packs" provide specialized methodology and workflow guidance, ranging from product discovery to high-scale platform engineering.

---

## 🛠 The Skill Ecosystem

The ecosystem is designed as a modular toolkit where skills can be chained together or used in isolation depending on the project's scale and risk.

### 🌟 The Main Entry Point: `agentic-helper`
If you are unsure where to start, trigger the **`agentic-helper`**. 
It acts as the **Universal Workflow Consultant** and orchestrator. It will interview you about your goal (Platform vs. Component scale) and recommend the specific skills and CLI commands required for your journey.

### 🗺 Methodology Routing

| Skill | Primary Use Case | Scope |
|---|---|---|
| **`unified-sdd`** | Multi-team, platform-scale development combining Sdd-Bmad, Sdd-OpenSpec, and Sdd-Speckit. | **Platform** |
| **`sdd-speckit`** | Turning product briefs and vague ideas into executable plans and task lists. | **Platform/Component** |
| **`sdd-openspec` (op)** | Proposing and implementing well-defined changes using delta specs. | **Component** |
| **`sdd-bmad`** | Progressive planning that routes work by role (Architect/Dev/Verifier). | **Platform/Component** |
| **`atdd` / `tdd`** | Driving implementation through automated acceptance and unit tests. | **Component** |
| **`explain-code`** | Generating visual diagrams and analogies to understand existing architecture. | **Knowledge** |

## 🔄 The 3-Pillar Mental Model

To use this ecosystem effectively, you must distinguish between **Who** is doing the work, **What** is being produced, and **When** it happens.

### 1. Roles (The Who)
The methodology assigns responsibility to specific personas (even if one person/agent plays multiple roles):
- **Product:** Defines "What" and "Why" (Goals, Success Metrics).
- **Architect:** Defines "How" (Tech Stack, Data Model, API Design).
- **Team Lead:** Routes work, manages change packages, and ensures delivery.
- **Developer:** Implements the code against the specifications.
- **Verifier:** Provides evidence that implementation matches the spec before merge.

### 2. Artifacts & Tools (The What & How)
These are the physical deliverables and the skills that generate them:
- **`product-wizard`** → PRD / Initiative.
- **`sdd-speckit`** → Feature Spec, Requirements, and Task Lists.
- **`sdd-openspec`** → Proposal, Design, and Delta Specs.
- **`atdd`** → Executable Acceptance Tests.
- **`agentic-agent` CLI** → The physical execution engine (tasks, context, validation).

### 3. Step/Flow (The When)
The lifecycle of a change follows a 5-phase "Golden Path":
1.  **Platform/Route:** Establish context and classify the request scale.
2.  **Specify:** Define requirements and remove ambiguity (Produce Specs).
3.  **Plan:** Break the spec into design and actionable tasks.
4.  **Deliver:** Implementation loop (Claim → Code → Test → Complete).
5.  **Archive:** Merge the delta specs into the "current truth" and close the change.

---

## 🗺 Methodology Routing
...

### Example 1: Building a New Feature from Scratch
> "I want to build a new payment gateway. Guide me."
1. Trigger **`product-wizard`** to define the requirements.
2. Trigger **`sdd-speckit`** to create the execution plan and tasks.
3. Use **`agentic-agent task claim`** to start the implementation.

### Example 2: Mapping a Large Brownfield Platform
> "We have a massive legacy monorepo. How do we start using SDD?"
1. Trigger **`platform-contextualizer`** to document the current state.
2. Trigger **`unified-sdd`** to establish the platform-wide engineering constitution.

---

## ❓ FAQ

### Q: Which skill is the most important?
**A:** `agentic-helper`. It is the brain that knows how to use all the other skills. It should be your first point of contact for any workflow questions.

### Q: Can I use multiple skills at once?
**A:** No. Following the **Prime Directive** in `SKILLS.md`, you should trigger **one skill per invocation**. If a task requires multiple skills, split it into sequential steps (e.g., first define the spec with `sdd-speckit`, then implement it with `atdd`).

### Q: What is the difference between Platform and Component scope?
**A:** **Platform** scope deals with multi-team standards, shared context, and cross-cutting architectures. **Component** scope is focused on the implementation of a specific service, library, or feature within that platform.

### Q: How do skills connect to the CLI?
**A:** Skills provide the **logic and strategy**. The `agentic-agent` CLI provides the **physical actions** (creating tasks, claiming branches, generating context files, and validating changes). Always use the CLI to track the work recommended by the skills.
