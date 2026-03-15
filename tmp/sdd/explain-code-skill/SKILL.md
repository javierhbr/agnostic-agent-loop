---
name: explain-code
description: Explains code with visual diagrams and analogies. Use when explaining how code works, teaching about a codebase, or when the user asks "how does this work?"
---

# Explain Code

Use this skill when the task is not to change code first, but to make code,
architecture, or runtime behavior understandable to humans.

This skill is a support skill. It works well with BMAD, OpenSpec, and Speckit
when teams need to explain existing code, planned changes, pull requests, or
implementation behavior across roles.

## Required explanation format

When explaining code, always include:

1. an analogy from everyday life
2. an ASCII diagram that shows flow, structure, or relationships
3. a step-by-step walkthrough of what happens
4. one gotcha, misconception, or common failure point

## Operating rules

- keep explanations conversational and concrete
- tie the explanation to the actual code or artifact being discussed
- prefer ASCII diagrams so the explanation stays portable in docs and chat
- use more than one analogy only when the first analogy is not enough
- explain cause and effect, not just definitions
- call out where data enters, changes, and leaves the flow

## Recommended output shape

### Analogy

Give one short analogy that matches the code's actual job.

### Diagram

Draw a small ASCII diagram that shows the important path.

### Walkthrough

Explain the flow in order:

1. where the input starts
2. which component or function handles it next
3. what state or data changes happen
4. what output, side effect, or decision comes out

### Gotcha

Call out one thing that reviewers, new developers, or product partners often
miss.

## How this skill supports the unified SDD workflow

- `Platform`
  - explain current platform structure, dependencies, and existing constraints
- `Route`
  - explain the current code path or blast radius before sizing a change
- `Specify`
  - explain existing behavior so teams can compare current vs expected behavior
- `Plan`
  - explain architecture choices, interfaces, and affected components
- `Deliver`
  - explain pull requests, code paths, and implementation gotchas during review

## Prompt pattern

Use prompts like:

- "Using the explain-code skill, explain how this code path works with an analogy, an ASCII diagram, a step-by-step walkthrough, and one gotcha."
- "Using the explain-code skill, explain how this pull request changes the flow, and call out one risk reviewers should pay attention to."
