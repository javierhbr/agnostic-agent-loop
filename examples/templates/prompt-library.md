# Prompt Library

*Copy-paste prompts organized by workflow stage. Adapt `[placeholders]` to your situation.*

---

## Quick Start by Scenario

### Starting from a Vague Idea

```
I have an idea for [feature name]. We currently [current state] and users 
[problem they experience].

Can you help me create a comprehensive PRD using the product-wizard skill?
```

### Starting from Requirements List

```
I have a requirements list for [project name]. Some items conflict and some 
are vague. Can you help me turn this into a structured PRD using the 
product-wizard skill?

Here are the requirements:
[paste your bullet points]
```

### Starting from Existing Process

```
We have a manual [process name] that we want to automate. I've documented 
the current workflow. Can you help me create a PRD using the product-wizard 
skill that preserves our business logic?

Current process:
[paste process description]
```

### Starting from Codebase

```
Our [system name] code is complex and needs refactoring. I've described the 
current system and what we want to change. Can you help create a PRD for 
this refactoring project using the product-wizard skill?

Current system:
[paste system description]

Desired changes:
[paste what needs to change]
```

---

## Discovery Phase Prompts

### Answering Problem Questions

```
Problem: [specific user pain point]

Evidence: [support ticket count / user feedback / metrics]

Impact: [time wasted / money lost / churn risk / competitive pressure]

Why now: [what changed / deadline / opportunity]
```

###Answering User Questions

```
We have [number] users in [number] segments:

1. [User Type 1] - [their role and needs]
2. [User Type 2] - [their role and needs]
3. [User Type 3] - [their role and needs]

Our primary focus is [User Type X] because [reason].
```

### Answering Success Metrics Questions

```
Success metrics:

1. [Metric name]: [baseline] → [target] within [timeframe]
   Measured by: [analytics tool / dashboard / manual tracking]

2. [Metric name]: [baseline] → [target] within [timeframe]
   Measured by: [analytics tool / dashboard / manual tracking]

We'll consider this successful if we hit [primary metric] even if others lag.
```

### Answering Constraints Questions

```
Constraints:

Timeline: [delivery date], [reason for deadline]
Budget: $[amount] or [team size] engineers for [duration]
Technical: [must integrate with X, must use Y, cannot use Z]
Compliance: [GDPR / HIPAA / SOC2 / other requirements]
```

### Answering Scope Questions

```
In scope (must-have for MVP):
- [Feature 1]
- [Feature 2]
- [Feature 3]

Out of scope (Phase 2 or never):
- [Feature A]
- [Feature B]
- [Feature C]

If we had to cut scope, we'd drop [feature] first.
```

---

## Refinement Prompts

### Clarifying Vague Requirements

```
When you say "[vague term like 'fast' or 'easy']", let me be more specific:
- [Concrete measurement]
- [Acceptable threshold]
- [How we'll test it]
```

### Resolving Conflicts

```
I see items [#X] and [#Y] conflict. Here's how to resolve it:
[Your decision and reasoning]

Let's prioritize [option] because [business justification].
```

### Adding Missing Context

```
For additional context:
- Company size: [employees] employees, [users] users
- Industry: [industry]
- Current tech stack: [technologies]
- Integration requirements: [systems we use]
```

### Narrowing Scope

```
That scope is too large for our [timeline/budget]. Let's focus only on:
- [Core feature 1]
- [Core feature 2]

Everything else moves to Phase 2. Our MVP success criteria is [metric].
```

---

## Development Planning Prompts

### Creating Dev Plan from PRD

```
Now that we have a PRD at [file path], can you create a development plan 
using the dev-plans skill? Please break this into phases with clear review 
checkpoints.
```

### Requesting Specific Phase Breakdown

```
For the development plan, please organize into [number] phases:
- Phase 1: [focus area]
- Phase 2: [focus area]
- Phase 3: [focus area]

Each phase should have a review checkpoint before proceeding.
```

### Adding Technical Details to Plan

```
For technical context in the dev plan:
- Frontend: [framework/language]
- Backend: [framework/language]
- Database: [database type]
- Deployment: [platform]
- Testing: [test framework]
```

---

## Acceptance Testing Prompts

### Generating Tests for a Task

```
For task [TASK-ID], please generate acceptance tests using the atdd skill. 
Show the RED-GREEN-REFACTOR cycle with failing tests first.
```

### Requesting Specific Test Coverage

```
When generating acceptance tests, please ensure coverage for:
- Happy path scenarios
- Edge cases: [specific edge cases]
- Error handling: [specific errors]
- Cross-browser compatibility: [browsers to test]
```

### Understanding Test Output

```
Can you explain these acceptance tests in non-technical terms? Specifically:
- What does each test verify?
- Why is this test important?
- How does it map to the acceptance criteria?
```

---

## Iteration & Refinement Prompts

### Requesting Changes to PRD

```
The [section name] section needs refinement. Specifically:
- [Change 1]
- [Change 2]

Can you update the PRD with these changes?
```

### Adding User Stories

```
Please add a user story for:
**As a** [user type]
**I want to** [action]
**So that** [benefit]

With acceptance criteria:
- [Criterion 1]
- [Criterion 2]
```

### Adjusting Success Metrics

```
Let's revise the success metrics:
- Remove: [metric that's not actionable]
- Add: [new metric] with target [target value]
- Change: [existing metric] target from [old] to [new]
```

---

## Troubleshooting Prompts

### AI Doesn't Understand Domain

```
Let me provide more context about [domain/industry]:
[Explain industry-specific terms, workflows, or constraints]

In our industry, [term] means [definition].
```

### PRD is Too Generic

```
This feels generic. Let me be more specific about our users:

[Detailed persona]: [age, role, goals, pain points, daily workflow]

They specifically struggle with [concrete problem] which costs them [time/money].
```

### Missing Technical Details

```
For technical context:
- Our platform: [web/mobile/desktop]
- Current architecture: [monolith/microservices/serverless]
- APIs we use: [third-party APIs]
- Data we handle: [types of data, volumes, sensitivity]
```

### Scope is Too Broad

```
This is too ambitious for [timeline]. Let's use the MoSCoW method:

Must have: [features]
Should have: [features]
Could have: [features]
Won't have (this release): [features]
```

---

## Advanced Prompts

### Requesting Alternatives

```
Before we commit to this approach, can you show 2-3 alternative solutions 
to [problem] with trade-offs for each?

Consider:
- Build vs. buy
- Simple MVP vs. comprehensive solution
- Different technical approaches
```

### Risk Analysis

```
What are the top 5 risks for this project? For each risk:
- Likelihood (high/medium/low)
- Impact (high/medium/low)
- Mitigation strategy
```

### Dependency Mapping

```
What are the dependencies for this project?
- Internal systems we need to integrate with
- External APIs or services
- Team dependencies (design, infrastructure, etc.)
- Data requirements
```

---

## Integration with Agnostic Agent

### Creating OpenSpec from PRD

```bash
# After PRD is created, use CLI to generate tasks
agentic-agent openspec init "<feature-name>" --from .agentic/spec/prd-<feature>.md
```

### Building Context for Task

```bash
# Before starting work on a task
agentic-agent context build --task <TASK-ID>
```

### Running Validation

```bash
# Before completing work
agentic-agent validate
```

---

## Tips for Effective Prompting

### ✅ Do

- **Be specific:** "200ms response time" not "fast"
- **Provide context:** Company size, user count, current state
- **Use evidence:** "30% of tickets" not "users complain"
- **Name the skill:** "using the product-wizard skill"
- **Ask follow-ups:** If unclear, probe deeper

### ❌ Don't

- **Be vague:** "Better notifications" (better how?)
- **Assume knowledge:** Explain your domain
- **Skip the why:** Always explain the problem
- **Ignore constraints:** Unlimited scope = unrealistic PRDs
- **Rush discovery:** Questions exist for a reason

---

## Example Full Conversation Flow

1. **You:** [Initial prompt with context]
2. **AI:** [Asks 4-6 discovery questions]
3. **You:** [Answer each with specifics]
4. **AI:** [Presents draft PRD sections]
5. **You:** [Request refinements]
6. **AI:** [Updates PRD]
7. **You:** "Save PRD to `.agentic/spec/prd-[name].md`"
8. **AI:** [Confirms save]
9. **You:** "Create development plan with dev-plans skill"
10. **AI:** [Creates phased plan]
11. **You:** "Generate acceptance tests for task [ID] with atdd skill"
12. **AI:** [Creates tests]

**Typical duration:** 30-60 minutes for complete PRD + plan + tests

---

## Prompt Patterns by Use Case

### For Product Managers

Focus on: Problem, users, success metrics, market opportunity

```
We're losing deals because [competitor feature]. Our sales team reports 
[specific feedback]. I want to create a PRD for [feature] that addresses 
this gap and helps us win [specific customer segment].
```

### For Business Analysts

Focus on: Process improvement, efficiency gains, ROI

```
Our current [process] takes [time] and costs $[amount] per [unit]. We want 
to automate [specific steps] to reduce to [target time/cost]. Here's the 
current workflow: [description]
```

### For Engineering Managers

Focus on: Technical debt, refactoring, system improvements

```
Our [system] has [technical issue] causing [business impact]. The codebase 
is [state]. We need to [refactor/rebuild] while maintaining [requirements]. 
Current architecture: [description]
```

---

**Remember:** Good prompts = good PRDs. Spend time crafting clear, specific prompts and the AI will produce clear, specific specifications.

---

[← Back to Main README](../README.md) | [View Scenarios](../README.md#choose-your-scenario)
