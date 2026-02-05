package components

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/javierbenavides/agentic-agent/internal/ui/styles"
)

// FileType represents the type of file system entry
type FileType int

const (
	FileTypeDir FileType = iota
	FileTypeFile
)

// FileEntry represents a file or directory
type FileEntry struct {
	Name     string
	Path     string
	Type     FileType
	IsHidden bool
	Selected bool
}

// FilePicker is a file/directory picker component
type FilePicker struct {
	Label        string
	RootDir      string
	CurrentDir   string
	Entries      []FileEntry
	CursorPos    int
	ShowHidden   bool
	DirsOnly     bool
	MultiSelect  bool
	Selected     map[string]bool
	Height       int
	Offset       int
}

// NewFilePicker creates a new file picker
func NewFilePicker(label, rootDir string, dirsOnly, multiSelect bool) FilePicker {
	fp := FilePicker{
		Label:       label,
		RootDir:     rootDir,
		CurrentDir:  rootDir,
		ShowHidden:  false,
		DirsOnly:    dirsOnly,
		MultiSelect: multiSelect,
		Selected:    make(map[string]bool),
		Height:      10,
		Offset:      0,
	}

	fp.loadDirectory()
	return fp
}

// loadDirectory loads the current directory contents
func (fp *FilePicker) loadDirectory() error {
	entries, err := os.ReadDir(fp.CurrentDir)
	if err != nil {
		return err
	}

	fp.Entries = []FileEntry{}

	// Add parent directory entry if not at root
	if fp.CurrentDir != fp.RootDir {
		fp.Entries = append(fp.Entries, FileEntry{
			Name: "..",
			Path: filepath.Dir(fp.CurrentDir),
			Type: FileTypeDir,
		})
	}

	// Add all entries
	for _, entry := range entries {
		name := entry.Name()

		// Skip hidden files if not showing them
		if !fp.ShowHidden && strings.HasPrefix(name, ".") {
			continue
		}

		// Skip files if dirs only
		if fp.DirsOnly && !entry.IsDir() {
			continue
		}

		path := filepath.Join(fp.CurrentDir, name)
		fileType := FileTypeFile
		if entry.IsDir() {
			fileType = FileTypeDir
		}

		fp.Entries = append(fp.Entries, FileEntry{
			Name:     name,
			Path:     path,
			Type:     fileType,
			IsHidden: strings.HasPrefix(name, "."),
			Selected: fp.Selected[path],
		})
	}

	// Sort: directories first, then alphabetically
	sort.Slice(fp.Entries, func(i, j int) bool {
		if fp.Entries[i].Name == ".." {
			return true
		}
		if fp.Entries[j].Name == ".." {
			return false
		}
		if fp.Entries[i].Type != fp.Entries[j].Type {
			return fp.Entries[i].Type == FileTypeDir
		}
		return fp.Entries[i].Name < fp.Entries[j].Name
	})

	// Reset cursor if out of bounds
	if fp.CursorPos >= len(fp.Entries) {
		fp.CursorPos = len(fp.Entries) - 1
	}
	if fp.CursorPos < 0 {
		fp.CursorPos = 0
	}

	return nil
}

// Update handles file picker messages
func (fp *FilePicker) Update(msg tea.Msg) (FilePicker, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if fp.CursorPos > 0 {
				fp.CursorPos--
				if fp.CursorPos < fp.Offset {
					fp.Offset = fp.CursorPos
				}
			}

		case "down", "j":
			if fp.CursorPos < len(fp.Entries)-1 {
				fp.CursorPos++
				if fp.CursorPos >= fp.Offset+fp.Height {
					fp.Offset = fp.CursorPos - fp.Height + 1
				}
			}

		case "enter":
			// Enter directory or select file
			if fp.CursorPos < len(fp.Entries) {
				entry := fp.Entries[fp.CursorPos]
				if entry.Type == FileTypeDir {
					fp.CurrentDir = entry.Path
					fp.CursorPos = 0
					fp.Offset = 0
					fp.loadDirectory()
				} else if fp.MultiSelect {
					// Toggle selection
					fp.Selected[entry.Path] = !fp.Selected[entry.Path]
					fp.loadDirectory()
				}
			}

		case " ":
			// Space to toggle selection in multi-select mode
			if fp.MultiSelect && fp.CursorPos < len(fp.Entries) {
				entry := fp.Entries[fp.CursorPos]
				if entry.Name != ".." {
					fp.Selected[entry.Path] = !fp.Selected[entry.Path]
					fp.loadDirectory()
				}
			}

		case "h":
			// Toggle hidden files
			fp.ShowHidden = !fp.ShowHidden
			fp.loadDirectory()
		}
	}

	return *fp, nil
}

// View renders the file picker
func (fp FilePicker) View() string {
	var b strings.Builder

	// Label
	b.WriteString(styles.BoldStyle.Render(fp.Label) + "\n")

	// Current directory
	relPath, _ := filepath.Rel(fp.RootDir, fp.CurrentDir)
	if relPath == "." {
		relPath = "/"
	}
	b.WriteString(styles.MutedStyle.Render(fmt.Sprintf("Current: %s", relPath)) + "\n\n")

	// File list
	visibleEntries := fp.Entries[fp.Offset:]
	if len(visibleEntries) > fp.Height {
		visibleEntries = visibleEntries[:fp.Height]
	}

	for i, entry := range visibleEntries {
		actualIndex := i + fp.Offset
		cursor := "  "
		if actualIndex == fp.CursorPos {
			cursor = styles.IconArrow + " "
		}

		// Selection indicator
		selectIndicator := " "
		if fp.MultiSelect && fp.Selected[entry.Path] {
			selectIndicator = "✓"
		}

		// Icon
		icon := "  "
		if entry.Type == FileTypeDir {
			if entry.Name == ".." {
				icon = "↑ "
			} else {
				icon = "📁"
			}
		} else {
			icon = "📄"
		}

		// Style
		style := styles.ListItemStyle
		if actualIndex == fp.CursorPos {
			style = styles.SelectedItemStyle
		}
		if entry.IsHidden {
			style = styles.MutedStyle
		}

		name := entry.Name
		if entry.Type == FileTypeDir && entry.Name != ".." {
			name += "/"
		}

		line := fmt.Sprintf("%s%s %s %s", cursor, selectIndicator, icon, name)
		b.WriteString(style.Render(line) + "\n")
	}

	// Help text
	b.WriteString("\n")
	helpParts := []string{"↑/↓ navigate", "Enter open/select"}
	if fp.MultiSelect {
		helpParts = append(helpParts, "Space toggle")
	}
	helpParts = append(helpParts, "h show hidden", "Esc back")

	selectedCount := len(fp.Selected)
	if selectedCount > 0 {
		b.WriteString(styles.SuccessStyle.Render(fmt.Sprintf("%d selected • ", selectedCount)))
	}
	b.WriteString(styles.HelpStyle.Render(strings.Join(helpParts, " • ")))

	return b.String()
}

// GetSelected returns all selected paths
func (fp FilePicker) GetSelected() []string {
	paths := []string{}
	for path, selected := range fp.Selected {
		if selected {
			paths = append(paths, path)
		}
	}
	sort.Strings(paths)
	return paths
}

// GetCurrentSelection returns the currently highlighted entry
func (fp FilePicker) GetCurrentSelection() string {
	if fp.CursorPos >= 0 && fp.CursorPos < len(fp.Entries) {
		entry := fp.Entries[fp.CursorPos]
		if entry.Name != ".." {
			return entry.Path
		}
	}
	return ""
}

// HasSelection returns true if any paths are selected
func (fp FilePicker) HasSelection() bool {
	return len(fp.Selected) > 0
}
