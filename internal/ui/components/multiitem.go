package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/javierbenavides/agentic-agent/internal/ui/styles"
)

// MultiItemEditor allows editing a list of strings
type MultiItemEditor struct {
	Label       string
	Items       []string
	CurrentItem textinput.Model
	EditIndex   int  // -1 means adding new
	IsEditing   bool
}

// NewMultiItemEditor creates a new multi-item editor
func NewMultiItemEditor(label string) MultiItemEditor {
	ti := textinput.New()
	ti.Placeholder = "Enter item and press Enter to add"
	ti.CharLimit = 200
	ti.Width = 60

	return MultiItemEditor{
		Label:       label,
		Items:       []string{},
		CurrentItem: ti,
		EditIndex:   -1,
		IsEditing:   false,
	}
}

// Update handles multi-item editor messages
func (m *MultiItemEditor) Update(msg tea.Msg) (MultiItemEditor, tea.Cmd) {
	if !m.IsEditing {
		return *m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			// Add or update item
			value := strings.TrimSpace(m.CurrentItem.Value())
			if value != "" {
				if m.EditIndex == -1 {
					// Adding new item
					m.Items = append(m.Items, value)
				} else {
					// Updating existing item
					m.Items[m.EditIndex] = value
				}
				m.CurrentItem.SetValue("")
				m.EditIndex = -1
			}
			return *m, nil

		case "esc":
			// Cancel editing
			m.CurrentItem.SetValue("")
			m.EditIndex = -1
			return *m, nil
		}
	}

	var cmd tea.Cmd
	m.CurrentItem, cmd = m.CurrentItem.Update(msg)
	return *m, cmd
}

// View renders the multi-item editor
func (m MultiItemEditor) View() string {
	var b strings.Builder

	// Label
	b.WriteString(styles.BoldStyle.Render(m.Label))
	b.WriteString(styles.MutedStyle.Render(" (optional)") + "\n\n")

	// Show existing items
	if len(m.Items) > 0 {
		for i, item := range m.Items {
			bullet := styles.IconBullet
			b.WriteString(fmt.Sprintf("  %s %s\n", bullet, item))
			_ = i // Could add edit/delete controls here
		}
		b.WriteString("\n")
	}

	// Show input if editing
	if m.IsEditing {
		if m.EditIndex == -1 {
			b.WriteString(styles.MutedStyle.Render("  Add item:") + "\n")
		} else {
			b.WriteString(styles.MutedStyle.Render(fmt.Sprintf("  Edit item %d:", m.EditIndex+1)) + "\n")
		}
		b.WriteString("  " + m.CurrentItem.View() + "\n")
		b.WriteString(styles.HelpStyle.Render("  Enter to add • Esc to cancel") + "\n")
	} else {
		b.WriteString(styles.HelpStyle.Render(fmt.Sprintf("  %d items • Press 'a' to add • Enter to continue", len(m.Items))) + "\n")
	}

	return b.String()
}

// StartEditing begins editing mode
func (m *MultiItemEditor) StartEditing() tea.Cmd {
	m.IsEditing = true
	m.EditIndex = -1
	return m.CurrentItem.Focus()
}

// StopEditing exits editing mode
func (m *MultiItemEditor) StopEditing() {
	m.IsEditing = false
	m.CurrentItem.Blur()
	m.CurrentItem.SetValue("")
	m.EditIndex = -1
}

// RemoveItem removes an item by index
func (m *MultiItemEditor) RemoveItem(index int) {
	if index >= 0 && index < len(m.Items) {
		m.Items = append(m.Items[:index], m.Items[index+1:]...)
	}
}

// GetItems returns all items
func (m MultiItemEditor) GetItems() []string {
	return m.Items
}

// HasItems returns true if there are items
func (m MultiItemEditor) HasItems() bool {
	return len(m.Items) > 0
}
