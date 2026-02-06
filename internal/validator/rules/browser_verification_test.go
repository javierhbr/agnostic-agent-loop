package rules

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/javierbenavides/agentic-agent/internal/validator"
	"github.com/javierbenavides/agentic-agent/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

// taskList is a helper struct for YAML marshaling in tests
type taskList struct {
	Tasks []models.Task `yaml:"tasks"`
}

func TestBrowserVerificationRule_Name(t *testing.T) {
	rule := &BrowserVerificationRule{}
	assert.Equal(t, "browser-verification", rule.Name())
}

func TestBrowserVerificationRule_Validate_NoInProgressTasks(t *testing.T) {
	// Setup temp directory
	tmpDir := t.TempDir()
	tasksDir := filepath.Join(tmpDir, ".agentic", "tasks")
	err := os.MkdirAll(tasksDir, 0755)
	require.NoError(t, err)

	// Create empty in-progress.yaml
	inProgressPath := filepath.Join(tasksDir, "in-progress.yaml")
	emptyList := taskList{Tasks: []models.Task{}}
	data, err := yaml.Marshal(emptyList)
	require.NoError(t, err)
	err = os.WriteFile(inProgressPath, data, 0644)
	require.NoError(t, err)

	// Run validation
	rule := &BrowserVerificationRule{}
	ctx := &validator.ValidationContext{
		ProjectRoot: tmpDir,
	}

	result, err := rule.Validate(ctx)
	require.NoError(t, err)
	assert.Equal(t, "PASS", result.Status)
	assert.Len(t, result.Errors, 0)
}

func TestBrowserVerificationRule_Validate_UITaskWithBrowserVerification(t *testing.T) {
	tmpDir := t.TempDir()
	tasksDir := filepath.Join(tmpDir, ".agentic", "tasks")
	err := os.MkdirAll(tasksDir, 0755)
	require.NoError(t, err)

	// Create task with UI file and browser verification
	task := models.Task{
		ID:          "US-001",
		Title:       "Add login button",
		Description: "Create login button component",
		Status:      "in-progress",
		Scope:       []string{"src/components/LoginButton.tsx"},
		Acceptance: []string{
			"Button renders correctly",
			"Typecheck passes",
			"Verify in browser using dev-browser skill",
		},
	}

	inProgressPath := filepath.Join(tasksDir, "in-progress.yaml")
	taskList := taskList{Tasks: []models.Task{task}}
	data, err := yaml.Marshal(taskList)
	require.NoError(t, err)
	err = os.WriteFile(inProgressPath, data, 0644)
	require.NoError(t, err)

	// Run validation
	rule := &BrowserVerificationRule{}
	ctx := &validator.ValidationContext{
		ProjectRoot: tmpDir,
	}

	result, err := rule.Validate(ctx)
	require.NoError(t, err)
	assert.Equal(t, "PASS", result.Status)
	assert.Len(t, result.Errors, 0)
}

func TestBrowserVerificationRule_Validate_UITaskWithoutBrowserVerification(t *testing.T) {
	tmpDir := t.TempDir()
	tasksDir := filepath.Join(tmpDir, ".agentic", "tasks")
	err := os.MkdirAll(tasksDir, 0755)
	require.NoError(t, err)

	// Create task with UI file but NO browser verification
	task := models.Task{
		ID:          "US-002",
		Title:       "Add dashboard component",
		Description: "Create dashboard UI",
		Status:      "in-progress",
		Scope:       []string{"src/components/Dashboard.jsx"},
		Acceptance: []string{
			"Component created",
			"Typecheck passes",
			// Missing browser verification
		},
	}

	inProgressPath := filepath.Join(tasksDir, "in-progress.yaml")
	taskList := taskList{Tasks: []models.Task{task}}
	data, err := yaml.Marshal(taskList)
	require.NoError(t, err)
	err = os.WriteFile(inProgressPath, data, 0644)
	require.NoError(t, err)

	// Run validation
	rule := &BrowserVerificationRule{}
	ctx := &validator.ValidationContext{
		ProjectRoot: tmpDir,
	}

	result, err := rule.Validate(ctx)
	require.NoError(t, err)
	assert.Equal(t, "FAIL", result.Status)
	require.Len(t, result.Errors, 1)
	assert.Contains(t, result.Errors[0], "US-002")
	assert.Contains(t, result.Errors[0], "Add dashboard component")
	assert.Contains(t, result.Errors[0], "browser verification")
}

func TestBrowserVerificationRule_Validate_NonUITask(t *testing.T) {
	tmpDir := t.TempDir()
	tasksDir := filepath.Join(tmpDir, ".agentic", "tasks")
	err := os.MkdirAll(tasksDir, 0755)
	require.NoError(t, err)

	// Create task with backend files (no UI)
	task := models.Task{
		ID:          "US-003",
		Title:       "Add database migration",
		Description: "Create user table migration",
		Status:      "in-progress",
		Scope:       []string{"internal/database/migrations/001_users.sql"},
		Acceptance: []string{
			"Migration runs successfully",
			"Table created with correct schema",
		},
	}

	inProgressPath := filepath.Join(tasksDir, "in-progress.yaml")
	taskList := taskList{Tasks: []models.Task{task}}
	data, err := yaml.Marshal(taskList)
	require.NoError(t, err)
	err = os.WriteFile(inProgressPath, data, 0644)
	require.NoError(t, err)

	// Run validation
	rule := &BrowserVerificationRule{}
	ctx := &validator.ValidationContext{
		ProjectRoot: tmpDir,
	}

	result, err := rule.Validate(ctx)
	require.NoError(t, err)
	assert.Equal(t, "PASS", result.Status)
	assert.Len(t, result.Errors, 0)
}

func TestBrowserVerificationRule_Validate_MultipleUIFiles(t *testing.T) {
	tmpDir := t.TempDir()
	tasksDir := filepath.Join(tmpDir, ".agentic", "tasks")
	err := os.MkdirAll(tasksDir, 0755)
	require.NoError(t, err)

	// Create task with multiple UI files but no browser verification
	task := models.Task{
		ID:     "US-004",
		Title:  "Add profile page",
		Status: "in-progress",
		Scope: []string{
			"src/pages/Profile.vue",
			"src/components/ProfileCard.vue",
			"src/styles/profile.css",
		},
		Acceptance: []string{
			"Profile page renders",
			"Card shows user info",
		},
	}

	inProgressPath := filepath.Join(tasksDir, "in-progress.yaml")
	taskList := taskList{Tasks: []models.Task{task}}
	data, err := yaml.Marshal(taskList)
	require.NoError(t, err)
	err = os.WriteFile(inProgressPath, data, 0644)
	require.NoError(t, err)

	// Run validation
	rule := &BrowserVerificationRule{}
	ctx := &validator.ValidationContext{
		ProjectRoot: tmpDir,
	}

	result, err := rule.Validate(ctx)
	require.NoError(t, err)
	assert.Equal(t, "FAIL", result.Status)
	require.Len(t, result.Errors, 1)
	assert.Contains(t, result.Errors[0], "US-004")
}

func TestBrowserVerificationRule_Validate_MixedUIAndNonUI(t *testing.T) {
	tmpDir := t.TempDir()
	tasksDir := filepath.Join(tmpDir, ".agentic", "tasks")
	err := os.MkdirAll(tasksDir, 0755)
	require.NoError(t, err)

	// Create task with both UI and non-UI files, no browser verification
	task := models.Task{
		ID:     "US-005",
		Title:  "Add user feature",
		Status: "in-progress",
		Scope: []string{
			"internal/services/user_service.go",
			"src/components/UserList.tsx",
			"internal/repositories/user_repo.go",
		},
		Acceptance: []string{
			"Service implements CRUD",
			"Component displays users",
		},
	}

	inProgressPath := filepath.Join(tasksDir, "in-progress.yaml")
	taskList := taskList{Tasks: []models.Task{task}}
	data, err := yaml.Marshal(taskList)
	require.NoError(t, err)
	err = os.WriteFile(inProgressPath, data, 0644)
	require.NoError(t, err)

	// Run validation
	rule := &BrowserVerificationRule{}
	ctx := &validator.ValidationContext{
		ProjectRoot: tmpDir,
	}

	result, err := rule.Validate(ctx)
	require.NoError(t, err)
	// Should FAIL because it has UI file (UserList.tsx) without browser verification
	assert.Equal(t, "FAIL", result.Status)
	require.Len(t, result.Errors, 1)
	assert.Contains(t, result.Errors[0], "US-005")
	assert.Contains(t, result.Errors[0], "Add user feature")
}

func TestBrowserVerificationRule_isUIFile_Extensions(t *testing.T) {
	rule := &BrowserVerificationRule{}

	testCases := []struct {
		path     string
		expected bool
	}{
		// UI files
		{"components/Button.tsx", true},
		{"pages/Home.jsx", true},
		{"views/Dashboard.vue", true},
		{"components/Card.svelte", true},
		{"templates/index.html", true},
		{"styles/main.css", true},
		{"styles/theme.scss", true},
		{"styles/variables.sass", true},
		{"styles/mixins.less", true},

		// Non-UI files
		{"services/user_service.go", false},
		{"repositories/user_repo.go", false},
		{"models/user.go", false},
		{"config/app.yaml", false},
		{"scripts/migrate.sql", false},
		{"README.md", false},
		{"tests/user_test.go", false},
	}

	for _, tc := range testCases {
		t.Run(tc.path, func(t *testing.T) {
			result := rule.isUIFile(tc.path)
			assert.Equal(t, tc.expected, result, "File: %s", tc.path)
		})
	}
}

func TestBrowserVerificationRule_isUIFile_Directories(t *testing.T) {
	rule := &BrowserVerificationRule{}

	testCases := []struct {
		path     string
		expected bool
	}{
		// UI directories with UI file extensions (.ts/.js only for directories)
		{"src/components/Button.ts", true},
		{"internal/ui/dashboard/widget.js", true},
		{"app/views/home/index.ts", true},
		{"frontend/pages/about.js", true},
		{"web/layouts/main.ts", true},

		// .go files are NOT UI files, even in UI directories
		{"src/components/Button.go", false},
		{"internal/ui/widget.go", false},

		// Non-UI directories
		{"internal/services/user.go", false},
		{"pkg/utils/helper.go", false},
		{"cmd/app/main.go", false},
	}

	for _, tc := range testCases {
		t.Run(tc.path, func(t *testing.T) {
			result := rule.isUIFile(tc.path)
			assert.Equal(t, tc.expected, result, "File: %s", tc.path)
		})
	}
}

func TestBrowserVerificationRule_hasBrowserVerification(t *testing.T) {
	rule := &BrowserVerificationRule{}

	testCases := []struct {
		name     string
		criteria []string
		expected bool
	}{
		{
			name: "Has verify in browser",
			criteria: []string{
				"Component renders",
				"Verify in browser",
				"Tests pass",
			},
			expected: true,
		},
		{
			name: "Has browser verification",
			criteria: []string{
				"Feature works",
				"Browser verification completed",
			},
			expected: true,
		},
		{
			name: "Has visual verification",
			criteria: []string{
				"Visual verification done",
			},
			expected: true,
		},
		{
			name: "Has test in browser",
			criteria: []string{
				"Test in browser",
			},
			expected: true,
		},
		{
			name: "Case insensitive",
			criteria: []string{
				"VERIFY IN BROWSER",
			},
			expected: true,
		},
		{
			name: "No browser verification",
			criteria: []string{
				"Component created",
				"Tests pass",
				"Typecheck passes",
			},
			expected: false,
		},
		{
			name:     "Empty criteria",
			criteria: []string{},
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := rule.hasBrowserVerification(tc.criteria)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestBrowserVerificationRule_Validate_MultipleTasks(t *testing.T) {
	tmpDir := t.TempDir()
	tasksDir := filepath.Join(tmpDir, ".agentic", "tasks")
	err := os.MkdirAll(tasksDir, 0755)
	require.NoError(t, err)

	// Create multiple tasks - mix of valid and invalid
	tasks := []models.Task{
		{
			ID:     "US-001",
			Title:  "Backend task",
			Status: "in-progress",
			Scope:  []string{"internal/services/auth.go"},
			Acceptance: []string{
				"Service works",
			},
		},
		{
			ID:     "US-002",
			Title:  "UI task with verification",
			Status: "in-progress",
			Scope:  []string{"src/components/Login.tsx"},
			Acceptance: []string{
				"Component renders",
				"Verify in browser",
			},
		},
		{
			ID:     "US-003",
			Title:  "UI task without verification",
			Status: "in-progress",
			Scope:  []string{"src/pages/Dashboard.jsx"},
			Acceptance: []string{
				"Page created",
			},
		},
		{
			ID:     "US-004",
			Title:  "Another UI task without verification",
			Status: "in-progress",
			Scope:  []string{"src/components/Header.vue"},
			Acceptance: []string{
				"Header displays correctly",
			},
		},
	}

	inProgressPath := filepath.Join(tasksDir, "in-progress.yaml")
	taskList := taskList{Tasks: tasks}
	data, err := yaml.Marshal(taskList)
	require.NoError(t, err)
	err = os.WriteFile(inProgressPath, data, 0644)
	require.NoError(t, err)

	// Run validation
	rule := &BrowserVerificationRule{}
	ctx := &validator.ValidationContext{
		ProjectRoot: tmpDir,
	}

	result, err := rule.Validate(ctx)
	require.NoError(t, err)
	assert.Equal(t, "FAIL", result.Status)
	// Should have 2 failures: US-003 and US-004
	require.Len(t, result.Errors, 2)
	assert.Contains(t, result.Errors[0], "US-003")
	assert.Contains(t, result.Errors[1], "US-004")
}
