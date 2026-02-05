# Phase 3: File Pickers & Interactive Init - Complete ✅

## Overview

Phase 3 adds file/directory picking capabilities and an interactive init wizard, making project setup and file selection intuitive for junior developers.

## What Was Built

### 1. File Picker Component (`components/filepicker.go`)

A fully-featured file and directory browser:

**Features:**
- Browse directories with keyboard navigation (↑/↓/j/k)
- Enter directories or select files
- Multi-select mode with Space bar
- Show/hide hidden files (press 'h')
- Visual indicators (📁 folders, 📄 files, ↑ parent)
- Selection checkmarks (✓)
- Scroll support for long lists
- Relative path display from root

**Modes:**
- **DirsOnly** - Only show directories (for scope selection)
- **MultiSelect** - Select multiple files/dirs (with Space)
- **Single Select** - Select one item (with Enter)

**Usage:**
```go
picker := components.NewFilePicker("Select Files", ".", false, true)
// Update in your wizard's Update()
picker, cmd = picker.Update(msg)
// Get selected paths
selected := picker.GetSelected()
```

### 2. Interactive Init Wizard (`models/init.go`)

A complete 7-step initialization wizard:

1. **Welcome** - Introduction to what will be created
2. **Project Name** - Validated input with real-time feedback
3. **AI Model** - Select from Claude/GPT-4 models with descriptions
4. **Validators** - Choose validation rule level (all/essential/none)
5. **Preview** - Review configuration and directory structure
6. **Initializing** - Animated spinner during setup
7. **Complete** - Success message with next steps

**Features:**
- Step-by-step guidance
- Model selection with descriptions
- Validator configuration explained
- Directory structure preview
- Configuration summary
- Error handling
- Next steps on completion

### 3. Command Integrations

#### Enhanced Init Command (`cmd/agentic-agent/init.go`)

```bash
# Interactive mode (auto-detect)
agentic-agent init
→ Launches wizard

# Flag mode (traditional)
agentic-agent init --name "My Project"
→ Uses flags

# Force non-interactive
agentic-agent init --no-interactive --name "Project"
→ Forces flag mode
```

## User Experience

### Interactive Init Flow

```
┌─ Initialize Agentic Agent Project ────┐
│                                        │
│  This wizard will guide you through:  │
│                                        │
│    • Project configuration             │
│    • Directory structure (.agentic/)   │
│    • Specification templates           │
│    • Task management files             │
│    • AI model configuration            │
│                                        │
│  Press Enter to continue               │
└────────────────────────────────────────┘

        ↓ [Project name]

┌─ Project Configuration ───────────────┐
│                                        │
│  Project Name                          │
│  ┌──────────────────────────────────┐ │
│  │ my-awesome-project               │ │
│  └──────────────────────────────────┘ │
│                                        │
│  Enter to continue • Esc to cancel     │
└────────────────────────────────────────┘

        ↓ [AI Model selection]

┌─ AI Model Selection ──────────────────┐
│                                        │
│  AI Model                              │
│                                        │
│  → Claude 3.5 Sonnet (Recommended)     │
│      Latest Claude model - excellent   │
│      balance of speed and capability   │
│                                        │
│    GPT-4 Turbo                         │
│      OpenAI's powerful model           │
│                                        │
│  ↑/↓ navigate • Enter continue         │
└────────────────────────────────────────┘

        ↓ [Validator selection]

┌─ Validation Rules ────────────────────┐
│                                        │
│  Validation Rules                      │
│                                        │
│  → All Validators (Recommended)        │
│      Enable all quality checks         │
│                                        │
│    Essential Only                      │
│      Only critical validators          │
│                                        │
│  Validators enforce best practices:    │
│    • Context files in directories      │
│    • Task scope enforcement            │
│    • Task size limits                  │
│                                        │
│  ↑/↓ navigate • Enter continue         │
└────────────────────────────────────────┘

        ↓ [Preview]

┌─ Confirm Configuration ───────────────┐
│                                        │
│  ┌────────────────────────────────┐   │
│  │ Project Name: my-awesome-proj  │   │
│  │ AI Model: Claude 3.5 Sonnet    │   │
│  │ Validators: All Validators     │   │
│  └────────────────────────────────┘   │
│                                        │
│  Directory structure to be created:    │
│                                        │
│  .agentic/                             │
│  ├── spec/           # Specifications  │
│  ├── context/        # Context summaries│
│  ├── tasks/          # Task management │
│  └── agent-rules/    # Tool configs    │
│  agnostic-agent.yaml # Project config  │
│                                        │
│  Press Enter to initialize             │
└────────────────────────────────────────┘

        ↓ [Initializing...]

┌─ Initializing Project ────────────────┐
│                                        │
│  ⠋ Initializing project...             │
│                                        │
│  Creating directory structure...       │
└────────────────────────────────────────┘

        ↓ [Success!]

┌─ Complete ────────────────────────────┐
│                                        │
│  ✓ Project initialized successfully!   │
│                                        │
│  Next steps:                           │
│                                        │
│  1. Create a task:                     │
│     agentic-agent task create          │
│                                        │
│  2. Review specs:                      │
│     ls .agentic/spec/                  │
│                                        │
│  3. Start working:                     │
│     agentic-agent work                 │
│                                        │
│  Press Enter to exit                   │
└────────────────────────────────────────┘
```

### File Picker Example

```
┌─ Select Scope Directories ────────────┐
│                                        │
│  Select Scope (files/directories)      │
│  Current: /                            │
│                                        │
│  → ✓ 📁 src/                           │
│      📁 tests/                         │
│      📁 docs/                          │
│      📄 README.md                      │
│      📄 go.mod                         │
│                                        │
│  2 selected • ↑/↓ navigate • Space toggle│
└────────────────────────────────────────┘
```

## Technical Details

### File Picker Architecture

**State Management:**
```go
type FilePicker struct {
    CurrentDir   string              // Current directory path
    Entries      []FileEntry         // Files/dirs in current dir
    CursorPos    int                 // Selected item index
    Selected     map[string]bool     // Multi-select state
    ShowHidden   bool                // Show/hide dotfiles
    DirsOnly     bool                // Filter files
    MultiSelect  bool                // Enable multi-selection
}
```

**Navigation:**
- `↑`/`k` - Move cursor up
- `↓`/`j` - Move cursor down
- `Enter` - Enter directory or select file
- `Space` - Toggle selection (multi-select mode)
- `h` - Toggle hidden files
- `Esc` - Cancel/back

**Features:**
- Automatic sorting (directories first, then alphabetical)
- Parent directory (`..`) navigation
- Scrolling for long lists (10 items visible)
- Visual indicators and icons
- Selection count display

### Init Wizard State Machine

```go
type InitWizardStep int

const (
    InitStepWelcome
    InitStepProjectName
    InitStepModel
    InitStepValidators
    InitStepPreview
    InitStepInitializing
    InitStepComplete
)
```

Each step validates before advancing and can be cancelled with Esc.

## Files Created/Modified

### New Files
- `internal/ui/components/filepicker.go` - File/directory picker
- `internal/ui/models/init.go` - Interactive init wizard

### Modified Files
- `cmd/agentic-agent/init.go` - Added interactive mode support

## Keyboard Shortcuts

### File Picker
- `↑`/`↓` or `k`/`j` - Navigate
- `Enter` - Open directory / Select file
- `Space` - Toggle selection (multi-select)
- `h` - Toggle hidden files
- `Esc` - Cancel

### Init Wizard
- `Enter` - Advance to next step
- `↑`/`↓` - Navigate selections
- `Esc` - Cancel wizard
- `Ctrl+C` - Force quit

## Testing

### Build and Test

```bash
# Build
go build -o agentic-agent ./cmd/agentic-agent

# Test interactive init
./agentic-agent init

# Test flag mode still works
./agentic-agent init --name "Test Project"

# View help
./agentic-agent init --help

# Test task creation (still works from Phase 2)
./agentic-agent task create
```

### Validation Checklist

Phase 3:
- ✅ Init wizard launches without flags
- ✅ Project name validation works
- ✅ AI model selection works
- ✅ Validator selection works
- ✅ Preview shows configuration
- ✅ Project initializes successfully
- ✅ Success screen shows next steps
- ✅ Flag mode still works identically
- ✅ `--no-interactive` forces flag mode

File Picker (component ready, integration in Phase 4):
- ✅ Directory navigation works
- ✅ File listing and sorting works
- ✅ Multi-select with Space works
- ✅ Hidden files toggle works
- ✅ Selection count displays
- ✅ Icons and indicators show correctly

## Integration Points for Phase 4

The file picker component is ready to be integrated into:

1. **Task Creation Wizard** - For spec refs, scope, and outputs selection
2. **Context Generate** - For directory selection
3. **Work Command** - For scope review

Example integration:
```go
// In task creation wizard
case TaskStepSelectScope:
    if !m.filePicker.Initialized {
        m.filePicker = components.NewFilePicker("Select Scope", ".", true, true)
    }
    m.filePicker, cmd = m.filePicker.Update(msg)

    // On Enter (when done selecting)
    if keyMsg.String() == "enter" && m.filePicker.HasSelection() {
        selectedPaths := m.filePicker.GetSelected()
        // Use selectedPaths...
        m.step = NextStep
    }
```

## Backward Compatibility

**100% Maintained:**
- All existing commands work identically
- Flag-based init unchanged
- Scripts and CI/CD unaffected
- Only bare `init` triggers wizard

## What's Next

### Phase 4: Work Command & Task Selection (Planned)
- Integrate file picker into task creation
- Task selection menu with tabs (Backlog/In Progress/Done)
- Complete work workflow wizard
- Progress tracking
- Validation results viewer

### Phase 5: Polish & Advanced (Planned)
- Search/filter in file picker
- Bulk operations
- Task templates
- Help system (`?` key)
- Back/undo navigation
- Performance optimization

## Success Metrics

✅ **Phase 3 Complete:**
- File picker component built and tested
- Interactive init wizard working
- AI model and validator selection
- Configuration preview
- Full backward compatibility
- Documentation updated

✅ **Component Reusability:**
- File picker ready for multiple integrations
- Consistent with existing UI patterns
- Well-documented for future use

## Demo Commands

```bash
# Try the interactive init wizard
./agentic-agent init

# When prompted:
# 1. Press Enter on welcome
# 2. Enter project name: "my-test-project"
# 3. Select AI model (default: Claude 3.5 Sonnet)
# 4. Select validators (default: All)
# 5. Review preview
# 6. Press Enter to initialize
# 7. See success message

# Verify project was created
ls -la .agentic/

# Try flag mode (still works)
rm -rf .agentic agnostic-agent.yaml
./agentic-agent init --name "Another Project"

# Try task creation (Phase 2 still works)
./agentic-agent task create
```

## File Picker Demo (Component Ready)

The file picker is built and ready for integration. Here's how it will work when integrated:

```bash
# Future: Task creation with file picker
./agentic-agent task create

# Step through wizard...
# When asked "Add scope?"
# → Select "Yes"
#
# File picker appears:
# ┌─ Select Scope ─────────────────┐
# │ Current: /                     │
# │  → 📁 src/                     │
# │    📁 internal/                │
# │    📁 cmd/                     │
# │  ↑/↓ nav • Enter open • Space  │
# └────────────────────────────────┘
#
# Navigate with arrows, select with Space
# Selected directories get ✓ checkmark
# Press Enter when done
```

## Conclusion

Phase 3 successfully delivers:

1. **Powerful File Picker** - Navigate, select files/dirs, multi-select
2. **Interactive Init** - Complete project setup wizard
3. **Better Onboarding** - AI model and validator configuration
4. **Foundation for Phase 4** - File picker ready for integration

Junior developers can now:
- Set up projects through guided wizards
- Choose AI models with descriptions
- Configure validators with explanations
- Preview before committing changes
- Get clear next steps

The file picker component is production-ready and waiting to be integrated into task creation and context workflows in Phase 4!

---

**Status:** Phase 3 Complete ✅
**Next:** Phase 4 - Work Command & Complete Workflows
**Built with:** Bubble Tea, Bubbles, Lipgloss, Go 1.22+
