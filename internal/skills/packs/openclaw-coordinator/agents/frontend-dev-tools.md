# FrontendDev — Tools & CLI Reference

## Task & Context Management

- `agentic-agent task list` — view backlog and in-progress tasks
- `agentic-agent task claim <ID>` — claim task (records branch + timestamp)
- `agentic-agent context build --task <ID>` — load context bundle (includes openspec + api_contracts + tech-stack)
- `agentic-agent validate` — run all quality validators before completing
- `agentic-agent task complete <ID>` — mark done (auto-captures commits)
- `agentic-agent openspec check <spec-id>` — verify spec readiness before claiming

## Framework & Build Tools

- **React/Next.js:**
  - `npm run dev` — start dev server
  - `npm run build` — production build
  - `npm test` — run tests (Jest)
  - `npm run lint` — ESLint + Prettier
  - `npx tsc --noEmit` — type checking
- **Vue/Vite:**
  - `npm run dev` — dev server
  - `npm run build` — build
  - `npm test` — Vitest
  - `npm run lint` — ESLint + Prettier
- **E2E Testing:**
  - `npx playwright test` — run all Playwright tests
  - `npx cypress run` — run Cypress tests

## Linting & Code Quality

- `eslint src/` — check for JS/TS errors
- `stylelint src/**/*.css` — check CSS
- `prettier --check src/` — check formatting
- `tsc --noEmit` — TypeScript type checking

## agent-browser — Visual Testing & Accessibility

**Installation:**
```bash
npm install -g agent-browser
agent-browser install
```

**Key Commands:**

Navigation & Interaction:
- `agent-browser open <url>` — open in browser
- `agent-browser click "<selector or aria label>"` — click element
- `agent-browser type "<selector>" "<text>"` — type into input
- `agent-browser fill-form "<form selector>" '{"field": "value"}'` — fill form
- `agent-browser hover "<selector>"` — hover element
- `agent-browser scroll` — scroll to element or by amount
- `agent-browser drag "<from>" "<to>"` — drag and drop

Inspection & Verification:
- `agent-browser text "<selector>"` — extract text content
- `agent-browser html "<selector>"` — get HTML
- `agent-browser accessibility-tree` — print accessibility tree (verify ARIA labels, roles)
- `agent-browser keyboard-nav` — test keyboard navigation order
- `agent-browser get-styles "<selector>"` — check computed styles (colors, fonts, spacing)
- `agent-browser contrast-check` — verify color contrast (WCAG AA/AAA)

Evidence & Debugging:
- `agent-browser screenshot "<url>" --output evidence/` — save screenshot for AC verification
- `agent-browser screenshot "<url>" --fullpage` — full-page screenshot
- `agent-browser pdf "<url>" --output evidence/page.pdf` — save as PDF
- `agent-browser console-logs` — retrieve console messages
- `agent-browser execute "<js code>"` — run JavaScript and return result

Browser Control:
- `agent-browser viewport --width 1920 --height 1080` — set viewport size
- `agent-browser viewport --device "iPhone 12"` — mobile emulation
- `agent-browser cookie set "<name>" "<value>"` — set cookie
- `agent-browser localStorage set "<key>" "<value>"` — set localStorage

**Workflow Example:**
```bash
# Start dev server
npm run dev

# Verify component renders
agent-browser open http://localhost:3000/components/button
agent-browser screenshot http://localhost:3000/components/button --output evidence/button-default.png

# Test interaction
agent-browser click "button[aria-label='Submit']"
agent-browser screenshot http://localhost:3000/components/button --output evidence/button-clicked.png

# Check accessibility
agent-browser accessibility-tree
agent-browser contrast-check
agent-browser keyboard-nav
```

## Key Paths

- `.agentic/spec/` — all openspec proposals with acceptance criteria
- `.agentic/contracts/` — API contracts (read before writing HTTP calls)
- `.agentic/context/` — global-context.md (design tokens, conventions), tech-stack.md (frameworks, versions)
- `.agentic/coordination/` — announcements.yaml, reservations.yaml
- `src/components/` — component source files
- `src/pages/` — page components
- `public/` — static assets

## API Contract Reference

Before writing any HTTP call:

1. Read the contract from `api_contracts[].path` in your context bundle
2. Contract is an OpenAPI spec (YAML or JSON) with:
   - `paths:` — endpoints (GET /api/v1/users, POST /api/v1/auth/login, etc.)
   - `components.schemas:` — data models
   - `securitySchemes:` — authentication (Bearer token, API key, etc.)
3. Never assume additional fields or endpoints not in the contract
4. If contract is incomplete: report `contract-deviation` to TechLead, don't invent

## Announcements Format

When announcing task completion to TechLead:

```yaml
- from_agent: frontend-dev
  to_agent: tech-lead
  project_id: proj-001
  status: complete
  summary: "Auth UI components implemented (login, register, forgot-password)"
  data:
    task_id: TASK-043
    branch: feature/auth-ui-v1
    components_created:
      - src/components/LoginForm.tsx
      - src/components/RegisterForm.tsx
      - src/components/ForgotPasswordForm.tsx
    test_results:
      unit_tests: 18 passed
      e2e_tests: 5 passed
    ac_coverage:
      - "Login form displays email/password fields" ✅
      - "Submit button disabled until both fields filled" ✅
      - "Login form submits POST to /auth/login" ✅
    screenshot_evidence:
      - evidence/login-form-default.png
      - evidence/login-form-filled.png
      - evidence/login-form-error.png
    accessibility_verified: true
    commits: [abc123, def456]

# Or if you find an API mismatch:
- from_agent: frontend-dev
  to_agent: tech-lead
  project_id: proj-001
  status: contract-deviation
  summary: "POST /auth/login returns 400 instead of spec's 401 for invalid credentials"
  data:
    task_id: TASK-043
    spec_ref: .agentic/contracts/auth-api.yaml
    endpoint: POST /api/v1/auth/login
    expected: "401 Unauthorized on invalid credentials"
    actual: "400 Bad Request"
    reproduction: "curl -X POST http://localhost:8080/api/v1/auth/login -d '{\"email\": \"test\", \"password\": \"wrong\"}'"
    severity: blocking
```

## Git Workflow

- Create feature branch: `git checkout -b feature/task-id-short-name` (TechLead handles main)
- Commit messages: `feat: <component name>: <AC description>` or `fix: <component>: <bug description>`
- Push to feature branch: `git push origin feature/<branch>`
- Never force-push or merge directly to main
