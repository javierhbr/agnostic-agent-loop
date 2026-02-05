package components

import (
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/javierbenavides/agentic-agent/internal/ui/styles"
)

// TextArea is a multi-line text input component
type TextArea struct {
	textarea textarea.Model
	Label    string
	Optional bool
}

// NewTextArea creates a new text area
func NewTextArea(label string, placeholder string, optional bool) TextArea {
	ta := textarea.New()
	ta.Placeholder = placeholder
	ta.CharLimit = 1000
	ta.SetWidth(60)
	ta.SetHeight(5)
	ta.ShowLineNumbers = false

	return TextArea{
		textarea: ta,
		Label:    label,
		Optional: optional,
	}
}

// Update handles textarea messages
func (t *TextArea) Update(msg tea.Msg) (TextArea, tea.Cmd) {
	var cmd tea.Cmd
	t.textarea, cmd = t.textarea.Update(msg)
	return *t, cmd
}

// View renders the text area
func (t TextArea) View() string {
	label := styles.BoldStyle.Render(t.Label)
	if t.Optional {
		label += styles.MutedStyle.Render(" (optional)")
	}

	view := label + "\n"
	view += t.textarea.View() + "\n"

	return view
}

// Value returns the current textarea value
func (t TextArea) Value() string {
	return t.textarea.Value()
}

// SetValue sets the textarea value
func (t *TextArea) SetValue(value string) {
	t.textarea.SetValue(value)
}

// Focus focuses the textarea
func (t *TextArea) Focus() tea.Cmd {
	return t.textarea.Focus()
}

// Blur blurs the textarea
func (t *TextArea) Blur() {
	t.textarea.Blur()
}

// Focused returns true if focused
func (t TextArea) Focused() bool {
	return t.textarea.Focused()
}
