---
name: tdd
description: TDD orchestrator enforcing red-green-refactor discipline with phased workflow, validation gates, and multi-agent coordination.
---

# skill:tdd

## Does exactly this

Enforces strict RED-GREEN-REFACTOR test-driven development discipline with phased workflow and validation gates.

---

## When to use

- Implementing new features or fixing bugs
- You need strict red-green-refactor cycle enforcement
- Coordinating multi-phase TDD workflows

---

## Do not use this skill when

- The task is exploratory or spike work (TDD after spike)
- UI layout-only changes with no testable logic
- The task is unrelated to writing or modifying code

---

## The TDD Cycle — Three Laws

```
RED    → Write a failing test that defines expected behaviour
         ↓
GREEN  → Write the minimum code to make it pass
         ↓
REFACTOR → Improve code quality while keeping tests green
         ↓
        Repeat...
```

**The Three Laws of TDD:**
1. Write production code only to make a failing test pass
2. Write only enough test to demonstrate failure
3. Write only enough code to make the test pass

---

## Phased Workflow (6 Phases)

**Phase 1: Test Specification and Design**
- Requirements Analysis
- Test Architecture (fixtures, mocks, test data)

**Phase 2: RED — Write Failing Tests**
- Write unit tests covering happy paths, edge cases, error scenarios (DO NOT implement)
- Verify all tests fail for the right reasons
- **Gate:** Do not proceed until all tests fail appropriately

**Phase 3: GREEN — Make Tests Pass**
- Write minimum code to pass tests (no extra features)
- Verify all tests pass
- **Gate:** All tests must pass before proceeding

**Phase 4: REFACTOR — Improve Code Quality**
- Apply SOLID principles, remove duplication, improve naming
- Refactor tests: remove duplication, improve names, extract fixtures
- Run tests after each change

**Phase 5: Integration Tests**
- Write integration tests (failing first)
- Make integration tests pass

**Phase 6: Continuous Improvement**
- Add performance, stress, boundary, and error recovery tests
- Final review: verify TDD process, check code and test quality

---

## Validation Checkpoints

**RED Phase:**
- [ ] All tests written before implementation
- [ ] All tests fail with meaningful error messages
- [ ] Failures are due to missing implementation, not syntax
- [ ] No test passes accidentally

**GREEN Phase:**
- [ ] All tests pass
- [ ] No extra code beyond test requirements
- [ ] Coverage ≥ 80% line, 75% branch
- [ ] No test was modified to make it pass

**REFACTOR Phase:**
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

## Anti-Patterns to Avoid

| Don't | Do |
|-------|-----|
| Write implementation before tests | Watch test fail first |
| Write tests that already pass | Ensure tests fail initially |
| Skip the refactor phase | Refactoring is NOT optional |
| Modify tests to make them pass | Fix implementation instead |
| Write tests after implementation | Tests are the specification |
| Multiple behaviours per test | One assertion per test |

---

## If you need more detail

→ `resources/red-green-refactor.md` — Phase-specific patterns, test design patterns, and strategies per phase
→ `resources/implementation-playbook.md` — Detailed examples, worked walkthroughs, and implementation scenarios
