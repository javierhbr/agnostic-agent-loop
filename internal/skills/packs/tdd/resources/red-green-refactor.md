# RED-GREEN-REFACTOR Phase Guide

Detailed instructions for each TDD phase. Referenced by the main TDD skill.

---

## RED Phase: Write Failing Tests

Generate comprehensive failing tests that define expected behavior.

### Core Requirements

1. **Test Structure**
   - Framework-appropriate setup (Jest/pytest/JUnit/Go/RSpec)
   - Arrange-Act-Assert pattern
   - `should_X_when_Y` naming convention
   - Isolated fixtures with no interdependencies

2. **Behavior Coverage**
   - Happy path scenarios
   - Edge cases (empty, null, boundary values)
   - Error handling and exceptions
   - Concurrent access (if applicable)

3. **Failure Verification**
   - Tests MUST fail when run
   - Failures for RIGHT reasons (not syntax/import errors)
   - Meaningful diagnostic error messages
   - No cascading failures

4. **Test Categories**
   - Unit: Isolated component behavior
   - Integration: Component interaction
   - Contract: API/interface contracts
   - Property: Mathematical invariants

### Framework Patterns

**JavaScript/TypeScript (Jest/Vitest)**
- Mock dependencies with `vi.fn()` or `jest.fn()`
- Use `@testing-library` for React components
- Property tests with `fast-check`

**Python (pytest)**
- Fixtures with appropriate scopes
- Parametrize for multiple test cases
- Hypothesis for property-based tests

**Go**
- Table-driven tests with subtests
- `t.Parallel()` for parallel execution
- Use `testify/assert` for cleaner assertions

**Ruby (RSpec)**
- `let` for lazy loading, `let!` for eager
- Contexts for different scenarios
- Shared examples for common behavior

### Edge Case Categories

- **Null/Empty**: undefined, null, empty string/array/object
- **Boundaries**: min/max values, single element, capacity limits
- **Special Cases**: Unicode, whitespace, special characters
- **State**: Invalid transitions, concurrent modifications
- **Errors**: Network failures, timeouts, permissions

### Quality Checklist

- [ ] Readable test names documenting intent
- [ ] One behavior per test
- [ ] No implementation leakage
- [ ] Meaningful test data (not 'foo'/'bar')
- [ ] Tests serve as living documentation

---

## GREEN Phase: Minimal Implementation

Implement the minimum code to make failing tests pass.

### Implementation Strategy

1. **Fake It** — Return hard-coded values when appropriate
2. **Obvious Implementation** — When solution is trivial and clear
3. **Triangulation** — Generalize only when multiple tests require it

### Rules

- Write the minimal code that could possibly work
- Don't add functionality not required by tests
- Use simple data structures initially
- Defer architectural decisions until refactor phase
- Don't add error handling unless tests require it
- Run tests after each change

### Progressive Implementation

| Stage | Approach |
|-------|----------|
| First test | Hard-coded return (fake it) |
| Second test | Simple if/else or pattern |
| Third+ tests | Generalize into real logic |

### Anti-Patterns

- Gold plating or adding unrequested features
- Implementing design patterns prematurely
- Performance optimizations without metrics
- Adding tests during green phase
- Refactoring during implementation

### When to Fake vs Go Real

| Fake It | Go Real |
|---------|---------|
| First test for new feature | Second/third test reveals pattern |
| Complex external dependencies | Implementation is obvious and simple |
| Approach still uncertain | Faking would be more complex than real code |

### Framework Progression Examples

**React:** Props-only → Hooks → Context
**Django:** Simple function → Class view → Generic view
**Express:** Inline logic → Extract middleware → Service layer
**Go:** Simple function → Interface → Full implementation

### Success Criteria

- All tests pass (green)
- No extra functionality beyond test requirements
- Code is readable even if not optimal
- No broken existing functionality
- Clear path to refactoring identified

---

## REFACTOR Phase: Improve Code Quality

Refactor with confidence using the test safety net.

### Core Process

**1. Pre-Assessment**
- Run tests to establish green baseline
- Analyze code smells and test coverage
- Document current performance metrics

**2. Code Smell Detection**

| Smell | Refactoring |
|-------|-------------|
| Duplicated code | Extract methods/classes |
| Long methods | Decompose into focused functions |
| Large classes | Split responsibilities |
| Long parameter lists | Parameter objects |
| Feature Envy | Move methods to appropriate classes |
| Primitive Obsession | Value objects |
| Switch statements | Polymorphism |

**3. SOLID Principles**
- **S**ingle Responsibility: One reason to change
- **O**pen/Closed: Open for extension, closed for modification
- **L**iskov Substitution: Subtypes substitutable
- **I**nterface Segregation: Small, focused interfaces
- **D**ependency Inversion: Depend on abstractions

**4. Refactoring Techniques**
- Extract Method/Variable/Interface
- Inline unnecessary indirection
- Rename for clarity
- Move Method/Field to appropriate classes
- Replace Magic Numbers with constants
- Replace Conditional with Polymorphism

**5. Performance Optimization**
- Profile to identify bottlenecks first
- Optimize algorithms and data structures
- Implement caching where beneficial
- Always measure before and after

**6. Incremental Steps**
- Make small, atomic changes
- Run tests after each modification
- Commit after each successful refactoring
- Keep refactoring separate from behavior changes

### Safety Checklist

Before committing:
- [ ] All tests pass (100% green)
- [ ] No functionality regression
- [ ] Performance metrics acceptable
- [ ] Code coverage maintained or improved
- [ ] Documentation updated

### Recovery Protocol

If tests fail during refactoring:
1. Immediately revert last change
2. Identify breaking refactoring
3. Apply smaller incremental changes
4. Use version control for safe experimentation

### Advanced Patterns

- **Strangler Fig**: Gradual legacy replacement
- **Branch by Abstraction**: Large-scale changes safely
- **Parallel Change**: Expand-contract pattern
- **Mikado Method**: Navigate dependency graphs
