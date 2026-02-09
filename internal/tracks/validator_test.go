package tracks

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateSpecContent_Complete(t *testing.T) {
	content := `# Specification: Auth

## Purpose

Allow users to securely log in.

## Constraints

Must use existing JWT library.

## Success Criteria

- [x] Users can register

## Alternatives Considered

### Approach A
We considered session-based auth but chose JWT.

## Design

### Architecture
REST API with middleware.

## Requirements

- [ ] Login endpoint
- [ ] Logout endpoint

## Acceptance Criteria

- [ ] Users can log in
- [ ] Tokens expire after 1h
`

	report := ValidateSpecContent(content)
	assert.True(t, report.Complete)
	assert.Empty(t, report.Missing)
	assert.Contains(t, report.Present, "purpose")
	assert.Contains(t, report.Present, "constraints")
	assert.Contains(t, report.Present, "success")
	assert.Contains(t, report.Present, "design")
	assert.Contains(t, report.Present, "requirements")
	assert.Contains(t, report.Present, "acceptance")
}

func TestValidateSpecContent_Incomplete(t *testing.T) {
	content := `# Specification: Auth

## Purpose

Allow users to log in.

## Constraints

<!-- What must we work within? -->

## Success Criteria

<!-- What does "done" look like? -->

## Design

<!-- High-level structure -->
`

	report := ValidateSpecContent(content)
	assert.False(t, report.Complete)
	assert.Contains(t, report.Present, "purpose")
	assert.Contains(t, report.Missing, "constraints")
	assert.Contains(t, report.Missing, "success")
	assert.Contains(t, report.Missing, "design")
	assert.Contains(t, report.Missing, "requirements")
	assert.Contains(t, report.Missing, "acceptance")
}

func TestValidateSpecContent_PlaceholderText(t *testing.T) {
	content := `# Specification: Auth

## Purpose

Why are we building this?

## Requirements

- [ ] Requirement 1
- [ ] Requirement 2
`

	report := ValidateSpecContent(content)
	assert.False(t, report.Complete)
	assert.Contains(t, report.Missing, "purpose")
	assert.Contains(t, report.Missing, "requirements")
}

func TestValidateSpecContent_AlternativesOptional(t *testing.T) {
	content := `# Specification: Auth

## Purpose

Real purpose here.

## Constraints

Real constraints here.

## Success Criteria

- [ ] Real criterion

## Design

Real design content.

## Requirements

- [ ] Real requirement

## Acceptance Criteria

- [ ] Real acceptance criterion
`

	report := ValidateSpecContent(content)
	assert.True(t, report.Complete)
	// alternatives is a warning, not a blocker
	assert.Contains(t, report.Warnings, "alternatives")
}
