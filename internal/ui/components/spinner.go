package components

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/javierbenavides/agentic-agent/internal/ui/styles"
)

// Spinner is a loading spinner component
type Spinner struct {
	spinner spinner.Model
	Message string
}

// NewSpinner creates a new spinner
func NewSpinner(message string) Spinner {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = styles.SuccessStyle.Copy()

	return Spinner{
		spinner: s,
		Message: message,
	}
}

// Init initializes the spinner
func (s Spinner) Init() tea.Cmd {
	return s.spinner.Tick
}

// Update handles spinner messages
func (s *Spinner) Update(msg tea.Msg) (Spinner, tea.Cmd) {
	var cmd tea.Cmd
	s.spinner, cmd = s.spinner.Update(msg)
	return *s, cmd
}

// View renders the spinner
func (s Spinner) View() string {
	return s.spinner.View() + " " + s.Message
}

// SetMessage updates the spinner message
func (s *Spinner) SetMessage(message string) {
	s.Message = message
}
