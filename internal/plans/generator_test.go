package plans

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateFromSpec_WithItems(t *testing.T) {
	dir := t.TempDir()
	specPath := filepath.Join(dir, "spec.md")
	planPath := filepath.Join(dir, "plan.md")

	spec := `# Specification: Auth

## Requirements

- [ ] Implement login endpoint
- [ ] Implement logout endpoint
- [ ] Add session management

## Acceptance Criteria

- [ ] Users can log in with email
- [ ] Sessions expire after 1 hour
`
	require.NoError(t, os.WriteFile(specPath, []byte(spec), 0644))

	err := GenerateFromSpec(specPath, planPath, "Auth System")
	require.NoError(t, err)

	data, err := os.ReadFile(planPath)
	require.NoError(t, err)
	content := string(data)

	assert.Contains(t, content, "# Plan: Auth System")
	assert.Contains(t, content, "## Phase 1: Setup")
	assert.Contains(t, content, "## Phase 2: Implementation")
	assert.Contains(t, content, "Implement login endpoint")
	assert.Contains(t, content, "Implement logout endpoint")
	assert.Contains(t, content, "Add session management")
	assert.Contains(t, content, "## Phase 3: Validation")
	assert.Contains(t, content, "Users can log in with email")
	assert.Contains(t, content, "Sessions expire after 1 hour")
}

func TestGenerateFromSpec_PlaceholdersSkipped(t *testing.T) {
	dir := t.TempDir()
	specPath := filepath.Join(dir, "spec.md")
	planPath := filepath.Join(dir, "plan.md")

	spec := `# Specification: Test

## Requirements

- [ ] Requirement 1
- [ ] Requirement 2

## Acceptance Criteria

- [ ] Criterion 1
- [ ] Criterion 2
`
	require.NoError(t, os.WriteFile(specPath, []byte(spec), 0644))

	err := GenerateFromSpec(specPath, planPath, "Test")
	require.NoError(t, err)

	data, err := os.ReadFile(planPath)
	require.NoError(t, err)
	content := string(data)

	// Placeholders should be skipped, fallback items used
	assert.NotContains(t, content, "Requirement 1")
	assert.Contains(t, content, "Core implementation")
}

func TestGenerateFromSpec_Empty(t *testing.T) {
	dir := t.TempDir()
	specPath := filepath.Join(dir, "spec.md")
	planPath := filepath.Join(dir, "plan.md")

	spec := `# Specification: Empty
`
	require.NoError(t, os.WriteFile(specPath, []byte(spec), 0644))

	err := GenerateFromSpec(specPath, planPath, "Empty")
	require.NoError(t, err)

	data, err := os.ReadFile(planPath)
	require.NoError(t, err)
	content := string(data)

	// Should have default fallback tasks
	assert.Contains(t, content, "Core implementation")
	assert.Contains(t, content, "Write tests")
}

func TestExtractSpecItems(t *testing.T) {
	dir := t.TempDir()
	specPath := filepath.Join(dir, "spec.md")

	spec := `# Spec

## Requirements

- [ ] Build API
- [x] Design schema

## Acceptance Criteria

- [ ] API returns 200
- [ ] Schema validates
`
	require.NoError(t, os.WriteFile(specPath, []byte(spec), 0644))

	reqs, acc, err := extractSpecItems(specPath)
	require.NoError(t, err)

	assert.Equal(t, []string{"Build API", "Design schema"}, reqs)
	assert.Equal(t, []string{"API returns 200", "Schema validates"}, acc)
}

func TestIsPlaceholderItem(t *testing.T) {
	assert.True(t, isPlaceholderItem("Requirement 1"))
	assert.True(t, isPlaceholderItem("Criterion 2"))
	assert.False(t, isPlaceholderItem("Build the authentication endpoint"))
}

func TestGenerateFromSpec_ParseableByParser(t *testing.T) {
	dir := t.TempDir()
	specPath := filepath.Join(dir, "spec.md")
	planPath := filepath.Join(dir, "plan.md")

	spec := `# Specification: Auth

## Requirements

- [ ] Login endpoint
- [ ] Token refresh

## Acceptance Criteria

- [ ] Returns valid JWT
`
	require.NoError(t, os.WriteFile(specPath, []byte(spec), 0644))
	require.NoError(t, GenerateFromSpec(specPath, planPath, "Auth"))

	// The generated plan should be parseable by the existing plan parser
	plan, err := ParseFile(planPath)
	require.NoError(t, err)

	assert.Equal(t, "Auth", plan.Title)
	assert.True(t, len(plan.Phases) >= 3, "expected at least 3 phases, got %d", len(plan.Phases))

	// All tasks should be pending
	for _, phase := range plan.Phases {
		for _, task := range phase.Tasks {
			assert.Equal(t, PlanTaskPending, task.Status)
		}
	}

	// Check specific items are present
	var allTitles []string
	for _, phase := range plan.Phases {
		for _, task := range phase.Tasks {
			allTitles = append(allTitles, task.Title)
		}
	}
	joined := strings.Join(allTitles, "|")
	assert.Contains(t, joined, "Login endpoint")
	assert.Contains(t, joined, "Token refresh")
	assert.Contains(t, joined, "Returns valid JWT")
}
