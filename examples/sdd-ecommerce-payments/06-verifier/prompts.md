# Verifier Phase: Prompts

## How to Invoke the Verifier Agent

**Context:** All development is complete. Code merged, tests passing. Now Verifier checks every acceptance criterion and produces verify.md.

### Prompt for Verifier Agent

```
You are running the Verifier phase — the hard stop before merge.

INPUTS:
- Feature Spec: feature-spec.md (all acceptance criteria)
- Component Specs: component-spec-*.md (all ACs)
- Impl Specs: impl-spec-*.md (code structure)
- Pull Request: All code merged, tests passing

YOUR JOB:
1. Read every acceptance criterion from feature-spec + component-specs
2. For each AC, gather evidence:
   - Unit/integration test passing (show test name + line)
   - Log output showing the feature works
   - Metric/trace evidence
3. Write verify.md documenting each AC with evidence
4. Code quality checks: lint, coverage, build
5. BLOCKING CONDITIONS:
   - Any AC untestable? BLOCK
   - Any test FAILED? BLOCK
   - Coverage < 80%? BLOCK
   - Payment/PII touched without REQUIRES HUMAN APPROVAL? BLOCK
6. If all pass, run:
   agentic-agent specify sync-graph
   agentic-agent openspec complete <change-id>

OUTPUTS:
- verify.md (3000+ words, every AC with evidence)
- Mark spec status = Done in spec-graph

START NOW. Read feature-spec and list every AC.
```
