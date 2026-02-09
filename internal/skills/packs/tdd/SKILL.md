---
name: tdd
description: TDD orchestrator enforcing red-green-refactor discipline with phased workflow, validation gates, and multi-agent coordination.
metadata:
  model: opus
---

## Use this skill when

- Implementing new features or fixing bugs using test-driven development
- You need strict red-green-refactor cycle enforcement
- Coordinating multi-phase TDD workflows

## Do not use this skill when

- The task is exploratory or spike work (TDD after spike)
- UI layout-only changes with no testable logic
- The task is unrelated to writing or modifying code

## Instructions

Follow strict RED-GREEN-REFACTOR discipline. Each phase must complete before the next begins. Reference `resources/red-green-refactor.md` for phase-specific patterns and `resources/implementation-playbook.md` for detailed examples.

---

## Quick Reference: The TDD Cycle

```
RED    → Write a failing test that defines expected behavior
         ↓
GREEN  → Write the minimum code to make it pass
         ↓
REFACTOR → Improve code quality while keeping tests green
         ↓
        Repeat...
```

### Three Laws of TDD

1. Write production code only to make a failing test pass
2. Write only enough test to demonstrate failure
3. Write only enough code to make the test pass

### Test Pattern (AAA)

| Step | Purpose |
|------|---------|
| **Arrange** | Set up test data and dependencies |
| **Act** | Execute the code under test |
| **Assert** | Verify expected outcome |

---

## Phased Workflow

### Phase 1: Test Specification and Design

1. **Requirements Analysis** — Analyze requirements, define acceptance criteria, identify edge cases, create test scenarios.
2. **Test Architecture** — Design test structure, fixtures, mocks, and test data strategy.

### Phase 2: RED — Write Failing Tests

3. **Write Unit Tests (Failing)** — Write failing tests covering happy paths, edge cases, and error scenarios. DO NOT implement production code.
4. **Verify Test Failure** — Confirm all tests fail for the right reasons (missing implementation, not syntax errors). No false positives.

**Gate:** Do not proceed until all tests fail appropriately.

### Phase 3: GREEN — Make Tests Pass

5. **Minimal Implementation** — Write the minimum code to make tests pass. No extra features or optimizations.
6. **Verify Test Success** — Run all tests, confirm they pass, check coverage metrics.

**Gate:** All tests must pass before proceeding.

### Phase 4: REFACTOR — Improve Code Quality

7. **Code Refactoring** — Apply SOLID principles, remove duplication, improve naming. Run tests after each change.
8. **Test Refactoring** — Remove test duplication, improve names, extract fixtures. Coverage must be maintained.

### Phase 5: Integration Tests

9. **Write Integration Tests (Failing First)** — Test component interactions, API contracts, data flow.
10. **Implement Integration** — Make integration tests pass.

### Phase 6: Continuous Improvement

11. **Performance and Edge Case Tests** — Add stress tests, boundary tests, error recovery tests.
12. **Final Review** — Verify TDD process was followed, check code and test quality.

---

## Validation Checkpoints

### RED Phase
- [ ] All tests written before implementation
- [ ] All tests fail with meaningful error messages
- [ ] Failures are due to missing implementation
- [ ] No test passes accidentally

### GREEN Phase
- [ ] All tests pass
- [ ] No extra code beyond test requirements
- [ ] Coverage meets thresholds (80% line, 75% branch)
- [ ] No test was modified to make it pass

### REFACTOR Phase
- [ ] All tests still pass after refactoring
- [ ] Code complexity reduced
- [ ] Duplication eliminated
- [ ] Test readability improved

---

## Configuration Thresholds

| Metric | Threshold |
|--------|-----------|
| Line coverage | 80% |
| Branch coverage | 75% |
| Critical path coverage | 100% |
| Cyclomatic complexity | < 10 |
| Method length | < 20 lines |
| Class length | < 200 lines |

---

## Modes

### Incremental (default)
1. Write ONE failing test
2. Make ONLY that test pass
3. Refactor if needed
4. Repeat for next test

### Suite
1. Write ALL tests for a feature (failing)
2. Implement code to pass ALL tests
3. Refactor entire module
4. Add integration tests

---

## Anti-Patterns to Avoid

| Don't | Do |
|-------|-----|
| Write implementation before tests | Watch test fail first |
| Write tests that already pass | Ensure tests fail initially |
| Skip the refactor phase | Refactoring is NOT optional |
| Modify tests to make them pass | Fix implementation instead |
| Write tests after implementation | Tests are the specification |
| Multiple behaviors per test | One assertion per test |

---

## Failure Recovery

If TDD discipline is broken:
1. STOP immediately
2. Identify which phase was violated
3. Rollback to last valid state
4. Resume from correct phase
5. Document lesson learned

---

## When to Use TDD

| Scenario | TDD Value |
|----------|-----------|
| New feature | High |
| Bug fix | High (write test first) |
| Complex logic | High |
| Exploratory/spike | Low (spike, then TDD) |
| UI layout | Low |
