package components

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/javierbenavides/agentic-agent/internal/ui/styles"
)

// ValidatedInput is a text input with validation
type ValidatedInput struct {
	Input     textinput.Model
	Validator func(string) error
	Error     string
	Label     string
}

// NewValidatedInput creates a new validated input
func NewValidatedInput(label, placeholder string, validator func(string) error) ValidatedInput {
	input := textinput.New()
	input.Placeholder = placeholder
	input.CharLimit = 200
	input.Width = 50
	input.Focus()

	return ValidatedInput{
		Input:     input,
		Validator: validator,
		Label:     label,
	}
}

// Update handles input messages and validation
func (vi *ValidatedInput) Update(msg tea.Msg) (ValidatedInput, tea.Cmd) {
	var cmd tea.Cmd
	vi.Input, cmd = vi.Input.Update(msg)

	// Validate on change if validator is provided
	if vi.Validator != nil {
		if err := vi.Validator(vi.Input.Value()); err != nil {
			vi.Error = err.Error()
		} else {
			vi.Error = ""
		}
	}

	return *vi, cmd
}

// View renders the validated input
func (vi ValidatedInput) View() string {
	view := styles.BoldStyle.Render(vi.Label) + "\n"
	view += vi.Input.View() + "\n"

	if vi.Error != "" {
		view += styles.RenderError(vi.Error) + "\n"
	}

	return view
}

// Value returns the current input value
func (vi ValidatedInput) Value() string {
	return vi.Input.Value()
}

// IsValid returns true if the input is valid
func (vi ValidatedInput) IsValid() bool {
	return vi.Error == "" && vi.Input.Value() != ""
}

// SetValue sets the input value
func (vi *ValidatedInput) SetValue(value string) {
	vi.Input.SetValue(value)
	if vi.Validator != nil {
		if err := vi.Validator(value); err != nil {
			vi.Error = err.Error()
		} else {
			vi.Error = ""
		}
	}
}

// Focus focuses the input
func (vi *ValidatedInput) Focus() tea.Cmd {
	return vi.Input.Focus()
}

// Blur blurs the input
func (vi *ValidatedInput) Blur() {
	vi.Input.Blur()
}
