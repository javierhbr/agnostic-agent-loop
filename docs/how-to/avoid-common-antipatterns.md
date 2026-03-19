# Anti-Patterns to Avoid

When working with the Agentic Agent CLI and the SDD (Spec-Driven Development) workflow, there are a few common pitfalls that can derail your progress, clutter your context, or lead to misaligned code. 

Here are the most common anti-patterns and how to avoid them.

## 1. The "Monolith" Anti-pattern

**The Anti-pattern:** Writing a massive prompt or creating a single huge task (e.g., "Build an entire e-commerce backend") and sending it directly to `agentic-agent deliver`.

**Why it's bad:**
- Exceeds the AI's reliable context window and cognitive limits.
- If the AI fails halfway through, you have a massive, tangled mess of partially written code to clean up.
- Validation is impossible because the scope is too broad.

**The Solution:** Use `task decompose`. Break large features into atomic units of work. Ensure each task has specific, testable acceptance criteria before running `deliver`.

---

## 2. The "Orphaned Spec" Anti-pattern

**The Anti-pattern:** You spot a bug or a needed feature, and you jump straight into the code to fix it, bypassing the `.agentic/spec/` directories and the `agentic-agent specify` command.

**Why it's bad:**
- Your Markdown specifications become immediately outdated.
- The next time the agent runs a task relying on that spec, it will use the obsolete knowledge, potentially overwriting your manual fix or writing incompatible code.
- You lose the "Source of Truth" guarantee.

**The Solution:** Always update the Spec first. Run `agentic-agent specify` or update the markdown manually, then run `agentic-agent route` to validate the spec, *then* execute the work.

---

## 3. The "Premature Execution" Anti-pattern

**The Anti-pattern:** Generating a task and immediately running `agentic-agent deliver` without verifying the context or creating an implementation plan.

**Why it's bad:**
- The agent might start writing code using the wrong architectural patterns because it hasn't mapped the local directory context.
- You end up burning tokens on an approach you would have vetoed if you had seen a plan first.

**The Solution:** Always run `agentic-agent plan` before `deliver`. Let the external planner (like SpecKit or BMAD) generate the execution strategy. Review the plan. Only when the plan is solid should you run `deliver`.

---

## 4. The "Global Context Bloat" Anti-pattern

**The Anti-pattern:** Shoving every single technical detail, library version, code snippet, and design choice into `.agentic/context/global-context.md`.

**Why it's bad:**
- This file is injected into *every single prompt*.
- If it grows too large, it wastes tokens on every turn.
- It dilutes the AI's attention, making it miss critical local instructions.

**The Solution:** Respect the 3-Tier Layered Context Model:
1. **Global (`global-context.md`)**: Only universal project rules and core architecture.
2. **Track/Feature (`decisions.md`)**: Specific to the current initiative.
3. **Task/Local (`assumptions.md` / `rolling-summary.md`)**: Ephemeral context that is cleared when the task is done.

Keep the global context lean and ruthless.

---

## 5. The "Silent Failure" Anti-pattern

**The Anti-pattern:** The agent encounters a compilation error or test failure during `deliver`, and instead of fixing it, it just moves on to the next file or marks the task as "done".

**Why it's bad:**
- You accumulate technical debt instantly.
- Downstream tasks will fail because the foundational code is broken.

**The Solution:** Ensure you are utilizing the `agentic-agent sync` and `agentic-agent validate` commands. The framework relies on strict Validation Gates. If a unit test or behavioral integrity check fails, the task should *never* reach the "done" state.

---

## 6. The "No Platform Alignment" Anti-pattern

**The Anti-pattern:** Starting implementation work in a component repo without checking platform-ref.yaml or confirming alignment with the platform version.

**Why it's bad:**

- Component work may drift from platform truth without anyone noticing.
- Impact tiers from dependency-map.md are ignored, leading to uncoordinated changes across tier 1 dependencies.
- Verification cannot confirm platform compliance if the alignment was never established.

**The Solution:** Always populate platform-ref.yaml in the Assess phase. Pin the platform version, list affected contracts and capabilities, and populate impact tiers from dependency-map.md. Read impact tiers before writing design.md — tier 1 entries become hard constraints.

---

## 7. The "Skipping Assess" Anti-pattern

**The Anti-pattern:** Jumping directly from a request to writing specs (Specify phase) without going through the Assess phase to classify the change and open a change package.

**Why it's bad:**

- Without classification, you don't know if the change is component-only or shared, and may miss platform-level coordination.
- No change package means no canonical execution container — artifacts are scattered.
- Ownership is not verified against component-ownership files, leading to wrong teams owning the wrong changes.

**The Solution:** Always run the Assess phase. Use BMAD to classify (greenfield/brownfield, size, impact). Read component-ownership files. Open a change package with platform-ref.yaml and jira-traceability.yaml before any spec work begins.

---

## 8. The "BMAD in Component Repos" Anti-pattern

**The Anti-pattern:** Using BMAD or Speckit tools inside component repositories instead of keeping them OpenSpec-only.

**Why it's bad:**

- Creates competing sources of intent between platform-side and component-side artifacts.
- Violates the core boundary rule (ADR-011): component repos use OpenSpec ONLY.
- Makes traceability harder because artifacts from different tools have different formats and assumptions.

**The Solution:** Follow the core boundary rule strictly. Platform-side work may use BMAD + OpenSpec + Speckit. Once work enters a component repo, use OpenSpec only for proposal.md, delta specs, design.md, tasks.md, and all delivery artifacts.
