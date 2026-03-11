# Agentic Agent Documentation

Welcome to the Agentic Agent framework documentation! We follow the **Diátaxis** documentation framework, dividing our documentation into four distinct quadrants based on what you are trying to achieve.

## 🎓 Tutorials (Learning)
*Step-by-step guides for beginners to achieve a basic level of competence.*

- [01: Getting Started](tutorials/01-getting-started.md) - Complete CLI tutorial and setup guide.
- [02: Worktree PR Workflow](tutorials/02-worktree-pr-workflow.md) - How to safely isolate tasks using worktrees.

## 🛠️ How-To Guides (Problem Solving)
*Step-by-step guides to answer specific questions or solve specific problems.*

- [How to Avoid Common Anti-Patterns](how-to/avoid-common-antipatterns.md) - *Start here to prevent costly mistakes.*
- [How to Use the SDD Process](how-to/use-sdd-process.md) - End-to-end Spec-Driven Development workflow.
- [How to Maintain Layered Context](how-to/maintain-layered-context.md) - Keeping your Agent's memory clean.
- [How to Use Checkpoints and Resume](how-to/use-checkpoints-and-resume.md) - Pause and resume agentic runs.

## 📖 Reference (Information)
*Technical descriptions, CLI commands, and objective information.*

- [CLI Reference: Spec-Driven Development](reference/spec-driven-development-cli.md) - All `specify` and `deliver` commands.
- [Layered Context Model](reference/layered-context-model.md) - The 3-tier memory schema definitions.
- [SDD Skills Directory](reference/sdd-skills.md) - Full list of available agent skills.
- [Agentic Helper Agent](reference/agentic-helper-agent.md) - Internal agent documentation.
- [Worktree PR Architecture](reference/worktree-pr-architecture.md) - Deep dive into gitops.
- [Integrations Reference](reference/integrations.md) - Superpowers, Ralph PDR, AgentSkills.

## 🧠 Explanation (Understanding)
*Clarifies and illuminates topics, providing background and conceptual context.*

- [SDD Methodology](explanation/sdd-methodology.md) - Why we build spec-first.
- [Layered Context Theory](explanation/layered-context-theory.md) - The psychology behind limiting agent context.
- [Autopilot vs Ralph Loop](explanation/autopilot-vs-ralph.md) - Comparing iterative paradigms.
- [Multi-Agent Execution](explanation/multi-agent-execution.md) - How multiple agents orchestrate.

---

### Internal Architecture & Design
*(For framework contributors)*
- [Implementation Plan](architecture/plan.md)
- [Foundational Framework (TODO-01)](architecture/TODO-01.md)
- [Advanced Features (TODO-02)](architecture/TODO-02.md)
- [TOON Integration (TODO-03)](architecture/TODO-03.md)
- [Task Templates Design](architecture/TASK_TEMPLATES_FEATURE.md)

### BDD Testing
- [BDD Framework Overview](bdd/BDD_IMPLEMENTATION_SUMMARY.md)
- [BDD Quick Reference](bdd/BDD_QUICK_REFERENCE.md)
- [BDD Guide](bdd/BDD_GUIDE.md)