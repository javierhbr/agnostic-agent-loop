# QADev — Testing Tools & CLI Reference

## Task & Context Management

- `agentic-agent task list` — view backlog and in-progress tasks
- `agentic-agent task claim <ID>` — claim task (records branch + timestamp)
- `agentic-agent context build --task <ID>` — load context bundle (includes sdd-openspec + tech-stack)
- `agentic-agent validate` — run all quality validators before completing
- `agentic-agent task complete <ID>` — mark done (auto-captures commits)
- `agentic-agent sdd-openspec check <spec-id>` — verify all ACs have test coverage before scoring

## Language Test Runners & Coverage

**Go:**
- `go test -cover ./...` — run all tests with coverage
- `go tool cover -html=coverage.out -o coverage/report.html` — generate HTML coverage report
- `go test -v ./...` — verbose test output

**JavaScript/TypeScript (Jest):**
- `npm test` — run all tests
- `npm test -- --coverage` — run with coverage
- `npm test -- --testPathPattern=auth` — run specific test file

**Python (pytest):**
- `pytest` — run all tests
- `pytest --cov=src/` — run with coverage
- `pytest --cov=src/ --cov-report=html` — HTML coverage report
- `coverage run -m pytest && coverage report` — detailed coverage

**Flutter/Dart:**
- `flutter test` — run all unit tests
- `flutter test --coverage` — generate coverage data
- `coverage format --lcov --in=coverage/lcov.info --out=coverage/report.html` — HTML coverage

## E2E & Integration Testing

**Playwright (React/Vue/Next.js):**
- `npx playwright test` — run all E2E tests
- `npx playwright test --headed` — run with visible browser
- `npx playwright test --debug` — debug mode
- `npx playwright test --project=chromium` — specific browser
- Screenshots: `await page.screenshot({ path: 'evidence/screenshot.png' })`

**Cypress:**
- `npx cypress run` — run all tests
- `npx cypress open` — interactive test runner
- Screenshots: `cy.screenshot('auth-form')`

**Flutter integration tests:**
- `flutter test integration_test/` — run integration tests
- `flutter test --coverage` — coverage from integration tests

## Browser & Accessibility Testing (agent-browser)

**Installation:**
```bash
npm install -g agent-browser
agent-browser install
```

**Visual Verification for Frontend:**
- Start dev server: `npm run dev`
- `agent-browser open http://localhost:3000`
- `agent-browser screenshot <url> --output evidence/screenshot.png` — capture state for evidence
- `agent-browser accessibility-tree` — verify ARIA labels and roles
- `agent-browser keyboard-nav` — test tab order and keyboard interaction
- `agent-browser get-styles "<selector>"` — verify colors, fonts match design tokens
- `agent-browser contrast-check` — verify WCAG AA/AAA color contrast

**Widget Verification for Flutter:**
- Run app: `flutter run`
- Use MCP widget tree inspection (if available) to verify layout
- Screenshots via device emulator or `flutter screenshot` command

## Security & Quality Scanning

**Backend:**
- `go vet ./...` — Go linter
- `gosec ./...` — Go security scanner
- `npm audit` — Node.js dependency vulnerabilities
- `trivy fs .` — scan for known vulnerabilities

**Frontend:**
- `npm audit` — dependency vulnerabilities
- `eslint src/` — code quality
- `tsc --noEmit` — TypeScript type checking

**Mobile (Flutter):**
- `flutter analyze` — Dart analyzer (must be zero warnings)
- `dart pub audit` — check Dart package vulnerabilities

## Performance Testing

- **Lighthouse:** `npm install -g lighthouse && lighthouse <url>` (React/Web)
- **Bundle size:** `npm run build && npm install -g bundlesize && bundlesize`
- **Load testing:** `npm install -g k6 && k6 run load-test.js`

## Key Paths

- `.agentic/spec/` — sdd-openspec proposals with acceptance criteria
- `.agentic/context/` — tech-stack.md (test framework, coverage tool), global-context.md
- `.agentic/coordination/` — announcements.yaml, reservations.yaml
- Coverage reports: `coverage/`, `coverage.out`, `coverage/lcov.info`
- Test results: `test-results.json` (if configured)
- Screenshots/evidence: `evidence/` directory

## 10-Point QA Rubric Checklist

Before announcing QA results, verify:

- [ ] **1. AC Coverage** — all sdd-openspec ACs have passing tests
- [ ] **2. Unit Tests** — core functions tested, ≥80% line coverage
- [ ] **3. Integration** — API contract matches spec, service boundaries tested, no deviations
- [ ] **4. E2E/Visual** — at least one happy-path E2E test with screenshots
- [ ] **5. Edge Cases** — ≥3 error/boundary condition tests
- [ ] **6. Performance** — response time/bundle size/load time not regressed
- [ ] **7. Security** — auth, input validation, injection prevention (if applicable)
- [ ] **8. Accessibility** — ARIA labels, keyboard nav, contrast (frontend/mobile) or N/A (backend)
- [ ] **9. Regression** — all existing tests still pass
- [ ] **10. Documentation** — test descriptions explain intent

**Score: ___ / 10**

- **8+** = QA pass → announce `qa-complete`
- **<8** = Request fixes → announce `qa-fix-requested` with failed points listed

## Announcements Format

When announcing QA results to TechLead:

```yaml
# QA Complete (score ≥ 8/10)
- from_agent: qa-dev
  to_agent: tech-lead
  project_id: proj-001
  status: qa-complete
  summary: "Auth API + UI features passed QA — 9/10"
  data:
    task_id: TASK-042
    qa_score: 9
    rubric:
      ac_coverage: pass
      unit_tests: pass (87% coverage)
      integration: pass (contract verified)
      e2e_visual: pass (3 flows tested)
      edge_cases: pass (5 tested)
      performance: pass
      security: pass (JWT validation, input sanitization)
      accessibility: n/a (backend task)
      regression: pass (all existing tests)
      documentation: partial (8/10 tests documented)
    test_counts:
      unit_tests: 34 passed
      integration_tests: 8 passed
      e2e_tests: 3 passed
    evidence_paths:
      - coverage/report.html (87% line coverage)
      - evidence/auth-login.png
      - evidence/auth-register.png
    failed_points: []
    commits: [abc123, def456]

# QA Fix Requested (score < 8/10)
- from_agent: qa-dev
  to_agent: backend-dev
  project_id: proj-001
  status: qa-fix-requested
  summary: "Auth API needs security fixes before QA pass (4/10)"
  data:
    task_id: TASK-042
    qa_score: 4
    failed_points:
      - "AC Coverage: POST /auth/login missing happy-path test"
      - "Integration: JWT token validation not tested"
      - "Security: Password hashing not verified in tests"
      - "Regression: 2 existing tests now failing"
      - "Documentation: Test descriptions unclear"
    needed_evidence:
      - "Unit test for JWT refresh token rotation"
      - "Integration test verifying bcrypt password hashing"
      - "Fix failing regression tests in existing suite"
      - "Add doc strings to test functions"
    ac_not_covered:
      - "POST /auth/refresh returns 200 with new token pair"
      - "Invalid password rejected with 401"
```

## Git Workflow

- Create feature branch: `git checkout -b feature/qa-tests-task-id` (TechLead handles main)
- Commit messages: `test: <component>: <AC description>` or `fix: <test name>: <what was wrong>`
- Push to feature branch: `git push origin feature/qa-tests-<branch>`
- Never force-push or merge directly to main
