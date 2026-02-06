package orchestrator

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Archiver handles archiving of progress and task data when switching branches
type Archiver struct {
	archiveDir     string
	lastBranchFile string
	progressText   string
	progressYAML   string
	tasksDir       string
}

// NewArchiver creates a new archiver with the specified paths
func NewArchiver(archiveDir, progressText, progressYAML, tasksDir string) *Archiver {
	return &Archiver{
		archiveDir:     archiveDir,
		lastBranchFile: filepath.Join(filepath.Dir(archiveDir), ".last-branch"),
		progressText:   progressText,
		progressYAML:   progressYAML,
		tasksDir:       tasksDir,
	}
}

// ArchiveIfBranchChanged archives the previous run if the branch has changed
func (a *Archiver) ArchiveIfBranchChanged(currentBranch string) error {
	lastBranch, err := a.readLastBranch()
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to read last branch: %w", err)
	}

	// If this is the first run or branch hasn't changed, just update the file
	if lastBranch == "" || lastBranch == currentBranch {
		return a.saveLastBranch(currentBranch)
	}

	// Branch changed - archive the previous run
	if err := a.archivePreviousRun(lastBranch); err != nil {
		return fmt.Errorf("failed to archive previous run: %w", err)
	}

	// Reset progress files for new run
	if err := a.resetProgressFiles(); err != nil {
		return fmt.Errorf("failed to reset progress files: %w", err)
	}

	// Save current branch
	return a.saveLastBranch(currentBranch)
}

// archivePreviousRun creates an archive of the previous run's data
func (a *Archiver) archivePreviousRun(branchName string) error {
	// Create archive folder: archive/YYYY-MM-DD-branch-name/
	timestamp := time.Now().Format("2006-01-02")
	// Clean branch name (remove "ralph/" prefix if present, sanitize)
	cleanBranch := filepath.Base(branchName)
	archivePath := filepath.Join(a.archiveDir, fmt.Sprintf("%s-%s", timestamp, cleanBranch))

	// Create archive directory
	if err := os.MkdirAll(archivePath, 0755); err != nil {
		return fmt.Errorf("failed to create archive directory: %w", err)
	}

	// Copy progress files if they exist
	if err := a.copyFileIfExists(a.progressText, filepath.Join(archivePath, "progress.txt")); err != nil {
		return err
	}
	if err := a.copyFileIfExists(a.progressYAML, filepath.Join(archivePath, "progress.yaml")); err != nil {
		return err
	}

	// Copy completed tasks
	doneTasks := filepath.Join(a.tasksDir, "done.yaml")
	if err := a.copyFileIfExists(doneTasks, filepath.Join(archivePath, "done.yaml")); err != nil {
		return err
	}

	fmt.Printf("Archived previous run (%s) to: %s\n", branchName, archivePath)
	return nil
}

// resetProgressFiles creates fresh progress files for a new run
func (a *Archiver) resetProgressFiles() error {
	// Reset progress.txt with new header
	header := fmt.Sprintf("# Progress Log\n\nStarted: %s\n\n---\n", time.Now().Format("2006-01-02 15:04:05"))
	if err := os.WriteFile(a.progressText, []byte(header), 0644); err != nil {
		return fmt.Errorf("failed to reset progress.txt: %w", err)
	}

	// Remove progress.yaml (will be recreated on first entry)
	if err := os.Remove(a.progressYAML); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove progress.yaml: %w", err)
	}

	return nil
}

// copyFileIfExists copies a file from src to dst if src exists
func (a *Archiver) copyFileIfExists(src, dst string) error {
	if _, err := os.Stat(src); os.IsNotExist(err) {
		return nil // File doesn't exist, skip
	}

	data, err := os.ReadFile(src)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", src, err)
	}

	if err := os.WriteFile(dst, data, 0644); err != nil {
		return fmt.Errorf("failed to write %s: %w", dst, err)
	}

	return nil
}

// readLastBranch reads the last branch name from the tracking file
func (a *Archiver) readLastBranch() (string, error) {
	data, err := os.ReadFile(a.lastBranchFile)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// saveLastBranch saves the current branch name to the tracking file
func (a *Archiver) saveLastBranch(branchName string) error {
	// Ensure directory exists
	dir := filepath.Dir(a.lastBranchFile)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	return os.WriteFile(a.lastBranchFile, []byte(branchName), 0644)
}

// GetLastBranch returns the last branch name (for external use)
func (a *Archiver) GetLastBranch() (string, error) {
	return a.readLastBranch()
}

// ListArchives returns a list of archived runs
func (a *Archiver) ListArchives() ([]string, error) {
	if _, err := os.Stat(a.archiveDir); os.IsNotExist(err) {
		return nil, nil
	}

	entries, err := os.ReadDir(a.archiveDir)
	if err != nil {
		return nil, err
	}

	var archives []string
	for _, entry := range entries {
		if entry.IsDir() {
			archives = append(archives, entry.Name())
		}
	}

	return archives, nil
}
