# Overview and Philosophy

## What Spec Kit is
Spec Kit is an open-source toolkit for Spec-Driven Development.
Its core idea is that specifications are not throwaway documentation; they become executable artifacts that drive implementation.

## Core philosophy
- define intent before implementation
- focus on the what and why first
- refine through multiple steps instead of one-shot prompting
- use project guardrails so agents make better decisions

## Development modes
### Greenfield
Use for 0-to-1 work where the feature or product starts from a blank slate.

### Brownfield
Use for iterative enhancement, modernization, or feature additions in an existing codebase.

## Why this matters for Codex skills
A good Codex skill should not only explain commands. It should teach the model:
- when to use each artifact
- how to separate concerns between artifacts
- how to preserve intent as work moves from spec to implementation
