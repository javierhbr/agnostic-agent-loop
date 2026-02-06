package functional

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/javierbenavides/agentic-agent/pkg/models"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

// TaskList represents a collection of tasks in a YAML file.
type TaskList struct {
	Tasks []models.Task `yaml:"tasks"`
}

// SetupTestProject creates an isolated test environment with the following:
// - Temporary directory
// - Git repository initialization (required for some validations)
// - Cleanup function registered with t.Cleanup()
// - Working directory changed to test directory
//
// Returns the test directory path and cleanup function.
func SetupTestProject(t *testing.T) string {
	t.Helper()

	// Create temporary test directory
	tmpDir, err := os.MkdirTemp("", "agentic-functional-test-*")
	require.NoError(t, err, "Failed to create temp directory")

	// Register cleanup to remove test directory
	t.Cleanup(func() {
		os.RemoveAll(tmpDir)
	})

	// Save original working directory
	originalDir, err := os.Getwd()
	require.NoError(t, err, "Failed to get current directory")

	// Change to test directory
	err = os.Chdir(tmpDir)
	require.NoError(t, err, "Failed to change to test directory")

	// Register cleanup to restore original directory
	t.Cleanup(func() {
		os.Chdir(originalDir)
	})

	// Initialize git repository (required for some validation rules)
	cmd := exec.Command("git", "init")
	cmd.Dir = tmpDir
	err = cmd.Run()
	require.NoError(t, err, "Failed to initialize git repository")

	// Configure git for the test (required for commits)
	configCmds := [][]string{
		{"git", "config", "user.email", "test@example.com"},
		{"git", "config", "user.name", "Test User"},
	}
	for _, cmdArgs := range configCmds {
		cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
		cmd.Dir = tmpDir
		err = cmd.Run()
		require.NoError(t, err, "Failed to configure git: %v", cmdArgs)
	}

	return tmpDir
}

// ExecuteCommandArgs executes a Cobra command with the given arguments and captures output.
// This is useful for testing commands that don't call os.Exit().
//
// Returns the captured stdout/stderr combined and any error.
func ExecuteCommandArgs(root *cobra.Command, args ...string) (string, error) {
	// Save original stdout/stderr
	oldStdout := os.Stdout
	oldStderr := os.Stderr
	defer func() {
		os.Stdout = oldStdout
		os.Stderr = oldStderr
	}()

	// Create pipes to capture output
	r, w, err := os.Pipe()
	if err != nil {
		return "", err
	}

	os.Stdout = w
	os.Stderr = w

	// Create a channel to capture output
	outC := make(chan string)
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	// Set command args and execute
	root.SetArgs(args)
	err = root.Execute()

	// Close writer and wait for output
	w.Close()
	output := <-outC

	return output, err
}

// ExecuteCommandInDir runs a command in a specific directory using subprocess execution.
// This is necessary for commands that call os.Exit(), as they cannot be tested in-process.
//
// Returns stdout/stderr combined, exit code, and any error.
func ExecuteCommandInDir(binaryPath, dir string, args ...string) (string, int, error) {
	cmd := exec.Command(binaryPath, args...)
	cmd.Dir = dir

	// Capture combined output
	output, err := cmd.CombinedOutput()

	// Get exit code
	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		}
	}

	return string(output), exitCode, err
}

// CaptureOutput captures stdout and stderr during function execution.
// Returns stdout and stderr as separate strings.
func CaptureOutput(f func()) (stdout, stderr string) {
	// Save original stdout/stderr
	oldStdout := os.Stdout
	oldStderr := os.Stderr
	defer func() {
		os.Stdout = oldStdout
		os.Stderr = oldStderr
	}()

	// Create pipes for stdout
	rOut, wOut, err := os.Pipe()
	if err != nil {
		panic(err)
	}
	os.Stdout = wOut

	// Create pipes for stderr
	rErr, wErr, err := os.Pipe()
	if err != nil {
		panic(err)
	}
	os.Stderr = wErr

	// Capture stdout
	outC := make(chan string)
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, rOut)
		outC <- buf.String()
	}()

	// Capture stderr
	errC := make(chan string)
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, rErr)
		errC <- buf.String()
	}()

	// Execute function
	f()

	// Close writers and collect output
	wOut.Close()
	wErr.Close()
	stdout = <-outC
	stderr = <-errC

	return
}

// VerifyProjectStructure verifies that the standard agentic project structure exists.
func VerifyProjectStructure(t *testing.T, projectDir string) {
	t.Helper()

	requiredDirs := []string{
		".agentic",
		".agentic/tasks",
		".agentic/context",
		".agentic/spec",
		".agentic/agent-rules",
	}

	for _, dir := range requiredDirs {
		fullPath := filepath.Join(projectDir, dir)
		assert.DirExists(t, fullPath, "Required directory should exist: %s", dir)
	}

	requiredFiles := []string{
		"agnostic-agent.yaml",
		".agentic/tasks/backlog.yaml",
		".agentic/tasks/in-progress.yaml",
		".agentic/tasks/done.yaml",
	}

	for _, file := range requiredFiles {
		fullPath := filepath.Join(projectDir, file)
		assert.FileExists(t, fullPath, "Required file should exist: %s", file)
	}
}

// VerifyTaskFile parses and verifies a task YAML file.
// Returns the parsed TaskList for further assertions.
func VerifyTaskFile(t *testing.T, filePath string) *TaskList {
	t.Helper()

	// Verify file exists
	assert.FileExists(t, filePath, "Task file should exist")

	// Read file
	data, err := os.ReadFile(filePath)
	require.NoError(t, err, "Failed to read task file")

	// Parse YAML
	var taskList TaskList
	err = yaml.Unmarshal(data, &taskList)
	require.NoError(t, err, "Failed to parse task file YAML")

	return &taskList
}

// VerifyTaskInFile verifies that a task with the given ID exists in the task file.
// Returns the found task for further assertions.
func VerifyTaskInFile(t *testing.T, filePath, taskID string) *models.Task {
	t.Helper()

	taskFile := VerifyTaskFile(t, filePath)

	// Find task with matching ID
	for _, task := range taskFile.Tasks {
		if task.ID == taskID {
			return &task
		}
	}

	t.Fatalf("Task %s not found in %s", taskID, filePath)
	return nil
}

// VerifyContextFile verifies that a context.md file exists and contains expected sections.
func VerifyContextFile(t *testing.T, filePath string) {
	t.Helper()

	// Verify file exists
	assert.FileExists(t, filePath, "Context file should exist")

	// Read content
	content, err := os.ReadFile(filePath)
	require.NoError(t, err, "Failed to read context file")

	contentStr := string(content)

	// Verify it's a markdown file with expected structure
	// Context files should have headings
	assert.Contains(t, contentStr, "#", "Context file should contain markdown headings")
}

// BuildBinary builds the agentic-agent binary for testing.
// Returns the path to the built binary.
func BuildBinary(t *testing.T) string {
	t.Helper()

	// Create temporary directory for binary
	binDir, err := os.MkdirTemp("", "agentic-binary-*")
	require.NoError(t, err, "Failed to create binary directory")

	t.Cleanup(func() {
		os.RemoveAll(binDir)
	})

	binPath := filepath.Join(binDir, "agentic-agent")

	// Build the binary
	cmd := exec.Command("go", "build", "-o", binPath, "./cmd/agentic-agent")
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "Failed to build binary: %s", string(output))

	return binPath
}

// AssertOutputContains is a helper that asserts the output contains all expected strings.
func AssertOutputContains(t *testing.T, output string, expectedStrings ...string) {
	t.Helper()

	for _, expected := range expectedStrings {
		assert.Contains(t, output, expected, "Output should contain: %s", expected)
	}
}

// AssertOutputNotContains is a helper that asserts the output does not contain any of the strings.
func AssertOutputNotContains(t *testing.T, output string, unexpectedStrings ...string) {
	t.Helper()

	for _, unexpected := range unexpectedStrings {
		assert.NotContains(t, output, unexpected, "Output should not contain: %s", unexpected)
	}
}

// CreateTestFile creates a file with the given content in the test directory.
func CreateTestFile(t *testing.T, filePath, content string) {
	t.Helper()

	// Create parent directories if needed
	dir := filepath.Dir(filePath)
	if dir != "." && dir != "" {
		err := os.MkdirAll(dir, 0755)
		require.NoError(t, err, "Failed to create directories for %s", filePath)
	}

	// Write file
	err := os.WriteFile(filePath, []byte(content), 0644)
	require.NoError(t, err, "Failed to write file %s", filePath)
}

// ParseTaskID extracts a task ID from command output.
// Looks for patterns like "TASK-1234567890" or "Created task: TASK-XXX".
func ParseTaskID(output string) (string, error) {
	// Look for TASK-XXXXXXXXXX pattern
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.Contains(line, "TASK-") {
			// Extract the task ID
			parts := strings.Fields(line)
			for _, part := range parts {
				if strings.HasPrefix(part, "TASK-") {
					// Clean up any trailing punctuation or ANSI codes
					taskID := strings.TrimSpace(part)
					taskID = strings.Trim(taskID, ".:,")
					return taskID, nil
				}
			}
		}
	}

	return "", fmt.Errorf("no task ID found in output")
}
