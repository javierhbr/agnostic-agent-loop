# Interactive UI Package

This package provides the interactive terminal UI components for the Agentic Agent CLI, powered by Bubble Tea.

## Structure

```
internal/ui/
├── components/         # Reusable UI components
│   ├── input.go       # Validated text input
│   ├── select.go      # Selection lists
│   └── spinner.go     # Loading spinners
├── styles/
│   └── theme.go       # Lipgloss styling theme
├── models/
│   └── start.go       # Wizard models (Bubble Tea)
└── helpers/
    └── mode.go        # Interactive mode detection
```

## Components

### ValidatedInput

A text input component with real-time validation.

```go
import "github.com/javierbenavides/agentic-agent/internal/ui/components"

// Create a validated input
input := components.NewValidatedInput(
    "Project Name",
    "my-project",
    func(s string) error {
        if len(s) == 0 {
            return fmt.Errorf("cannot be empty")
        }
        return nil
    },
)

// In your Update function
input, cmd = input.Update(msg)

// In your View function
view := input.View()

// Check validity
if input.IsValid() {
    value := input.Value()
}
```

### SimpleSelect

A keyboard-driven selection list.

```go
// Create options
options := []components.SelectOption{
    components.NewSelectOption(
        "Option 1",
        "Description of option 1",
        "value1",
    ),
    components.NewSelectOption(
        "Option 2",
        "Description of option 2",
        "value2",
    ),
}

// Create select
sel := components.NewSimpleSelect("Choose an option", options)

// In your Update function
sel = sel.Update(msg)

// Get selected value
value := sel.SelectedValue()
```

### Spinner

An animated loading indicator.

```go
spinner := components.NewSpinner("Loading...")

// In Init
return spinner.Init()

// In Update
spinner, cmd = spinner.Update(msg)

// In View
view := spinner.View()

// Update message
spinner.SetMessage("Processing...")
```

## Styling

All styles are defined in `styles/theme.go` for consistency.

### Colors

```go
import "github.com/javierbenavides/agentic-agent/internal/ui/styles"

// Available colors
styles.Primary    // Purple
styles.Secondary  // Cyan
styles.Success    // Green
styles.Warning    // Orange
styles.Error      // Red
styles.Muted      // Gray
```

### Pre-built Styles

```go
// Text styles
styles.TitleStyle.Render("My Title")
styles.SubtitleStyle.Render("Subtitle")
styles.HelpStyle.Render("Press Enter to continue")
styles.ErrorStyle.Render("Error message")
styles.SuccessStyle.Render("Success!")

// Component styles
styles.InputStyle        // Input boxes
styles.CardStyle         // Content cards
styles.BoxStyle          // Generic boxes
styles.ListItemStyle     // List items
styles.SelectedItemStyle // Selected items

// Helper functions
styles.RenderError("Something went wrong")
styles.RenderSuccess("Operation complete")
styles.RenderWarning("Be careful")
```

### Icons

```go
styles.IconCheckmark  // ✓
styles.IconCross      // ✗
styles.IconArrow      // →
styles.IconSpinner    // ⠋
styles.IconPending    // ○
styles.IconProgress   // ◐
styles.IconBullet     // •
styles.IconPrompt     // ›
```

## Mode Detection

Use the mode detection helper to determine if commands should run in interactive mode.

```go
import "github.com/javierbenavides/agentic-agent/internal/ui/helpers"

// In your command's Run function
if helpers.ShouldUseInteractiveMode(cmd) {
    // Launch interactive UI
    runInteractiveMode()
} else {
    // Use traditional flag-based mode
    runFlagMode(cmd, args)
}
```

**Auto-Detection Logic:**
- Interactive mode is enabled when:
  - No command-specific flags are provided
  - Input is a TTY (not piped or redirected)
  - `--no-interactive` flag is not set
  - Not running in CI environment

## Creating a New Wizard

Follow the Bubble Tea pattern to create new interactive wizards:

### 1. Define Your Model

```go
package models

import (
    tea "github.com/charmbracelet/bubbletea"
    "github.com/javierbenavides/agentic-agent/internal/ui/components"
    "github.com/javierbenavides/agentic-agent/internal/ui/styles"
)

type MyWizardStep int

const (
    StepWelcome MyWizardStep = iota
    StepInput
    StepConfirm
    StepComplete
)

type MyWizardModel struct {
    step      MyWizardStep
    input     components.ValidatedInput
    quitting  bool
}

func NewMyWizardModel() MyWizardModel {
    return MyWizardModel{
        step:  StepWelcome,
        input: components.NewValidatedInput("Label", "placeholder", validator),
    }
}
```

### 2. Implement Init, Update, View

```go
func (m MyWizardModel) Init() tea.Cmd {
    return m.input.Focus()
}

func (m MyWizardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "ctrl+c", "esc":
            m.quitting = true
            return m, tea.Quit
        case "enter":
            return m.handleEnter()
        }
    }

    // Handle step-specific updates
    switch m.step {
    case StepInput:
        var cmd tea.Cmd
        m.input, cmd = m.input.Update(msg)
        return m, cmd
    }

    return m, nil
}

func (m MyWizardModel) View() string {
    switch m.step {
    case StepWelcome:
        return m.renderWelcome()
    case StepInput:
        return m.renderInput()
    // ... other steps
    }
    return ""
}
```

### 3. Implement Render Functions

```go
func (m MyWizardModel) renderWelcome() string {
    var b strings.Builder

    b.WriteString(styles.TitleStyle.Render("Welcome!") + "\n\n")
    b.WriteString("Description of what this wizard does\n\n")
    b.WriteString(styles.HelpStyle.Render("Press Enter to continue") + "\n")

    return styles.ContainerStyle.Render(b.String())
}

func (m MyWizardModel) renderInput() string {
    var b strings.Builder

    b.WriteString(styles.TitleStyle.Render("Step Title") + "\n\n")
    b.WriteString(m.input.View() + "\n")
    b.WriteString(styles.HelpStyle.Render("Enter to continue • Esc to cancel") + "\n")

    return styles.ContainerStyle.Render(b.String())
}
```

### 4. Create the Command

```go
package main

import (
    tea "github.com/charmbracelet/bubbletea"
    "github.com/javierbenavides/agentic-agent/internal/ui/models"
    "github.com/spf13/cobra"
)

var myCmd = &cobra.Command{
    Use:   "mycommand",
    Short: "Description",
    Run: func(cmd *cobra.Command, args []string) {
        model := models.NewMyWizardModel()
        p := tea.NewProgram(model)

        if _, err := p.Run(); err != nil {
            fmt.Printf("Error: %v\n", err)
            os.Exit(1)
        }
    },
}
```

## Best Practices

### 1. Consistent Styling

Always use the pre-defined styles from `styles/theme.go`:

```go
// Good
styles.TitleStyle.Render("My Title")

// Avoid
lipgloss.NewStyle().Bold(true).Render("My Title")
```

### 2. Validation Feedback

Provide immediate validation feedback:

```go
validator := func(s string) error {
    if len(s) < 3 {
        return fmt.Errorf("must be at least 3 characters")
    }
    return nil
}
```

### 3. Help Text

Always show keyboard shortcuts:

```go
styles.HelpStyle.Render("↑/↓ Navigate • Enter Select • Esc Cancel")
```

### 4. Error Handling

Use styled error messages:

```go
if err != nil {
    return styles.RenderError("Operation failed: " + err.Error())
}
```

### 5. Confirmation Steps

Always confirm before destructive operations:

```go
const (
    StepInput
    StepConfirm    // Show what will happen
    StepExecute    // Actually do it
    StepComplete
)
```

## Testing Interactive UIs

### Manual Testing

```bash
# Build and run
go build -o agentic-agent ./cmd/agentic-agent
./agentic-agent mycommand

# Test keyboard navigation
# - Arrow keys for navigation
# - Enter to proceed
# - Esc to cancel
# - Ctrl+C to force quit
```

### Test Checklist

- [ ] Welcome screen displays correctly
- [ ] Input validation works (try invalid inputs)
- [ ] Navigation keys work (↑/↓/Enter/Esc)
- [ ] Error messages are clear
- [ ] Success state shows next steps
- [ ] Terminal resizing doesn't break layout
- [ ] Works in different terminals (Terminal.app, iTerm2, etc.)

## Examples

See `models/start.go` for a complete example of a multi-step wizard implementing:
- Welcome screen
- Validated input
- Selection menu
- Confirmation
- Async operation with spinner
- Success/error states

## Dependencies

- **Bubble Tea** - The Elm Architecture for terminal UIs
- **Bubbles** - Pre-built Bubble Tea components
- **Lipgloss** - Styling and layout
- **golang.org/x/term** - Terminal detection

## Resources

- [Bubble Tea Documentation](https://github.com/charmbracelet/bubbletea)
- [Bubbles Components](https://github.com/charmbracelet/bubbles)
- [Lipgloss Styling](https://github.com/charmbracelet/lipgloss)
- [The Elm Architecture](https://guide.elm-lang.org/architecture/)
