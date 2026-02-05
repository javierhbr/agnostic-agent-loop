package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/javierbenavides/agentic-agent/internal/ui/styles"
)

// SelectOption represents an option in a select list
type SelectOption struct {
	title       string
	description string
	value       string
}

func (i SelectOption) FilterValue() string { return i.title }
func (i SelectOption) Title() string       { return i.title }
func (i SelectOption) Description() string { return i.description }
func (i SelectOption) Value() string       { return i.value }

// NewSelectOption creates a new select option
func NewSelectOption(title, description, value string) SelectOption {
	return SelectOption{
		title:       title,
		description: description,
		value:       value,
	}
}

// SelectList is a selectable list component
type SelectList struct {
	list  list.Model
	Label string
}

// NewSelectList creates a new select list
func NewSelectList(label string, options []SelectOption, width, height int) SelectList {
	items := make([]list.Item, len(options))
	for i, opt := range options {
		items[i] = opt
	}

	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = lipgloss.NewStyle().
		Foreground(styles.Secondary).
		Bold(true).
		Border(lipgloss.NormalBorder(), false, false, false, true).
		BorderForeground(styles.Secondary).
		Padding(0, 0, 0, 1)

	delegate.Styles.SelectedDesc = lipgloss.NewStyle().
		Foreground(styles.Subtle).
		Border(lipgloss.NormalBorder(), false, false, false, true).
		BorderForeground(styles.Secondary).
		Padding(0, 0, 0, 1)

	l := list.New(items, delegate, width, height)
	l.Title = label
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = styles.TitleStyle

	return SelectList{
		list:  l,
		Label: label,
	}
}

// Update handles list messages
func (sl *SelectList) Update(msg tea.Msg) (SelectList, tea.Cmd) {
	var cmd tea.Cmd
	sl.list, cmd = sl.list.Update(msg)
	return *sl, cmd
}

// View renders the select list
func (sl SelectList) View() string {
	return sl.list.View()
}

// SelectedValue returns the selected option value
func (sl SelectList) SelectedValue() string {
	selected := sl.list.SelectedItem()
	if selected == nil {
		return ""
	}
	return selected.(SelectOption).Value()
}

// SelectedOption returns the selected option
func (sl SelectList) SelectedOption() *SelectOption {
	selected := sl.list.SelectedItem()
	if selected == nil {
		return nil
	}
	opt := selected.(SelectOption)
	return &opt
}

// SimpleSelect is a simpler select component without the full list UI
type SimpleSelect struct {
	Label       string
	Options     []SelectOption
	SelectedIdx int
}

// NewSimpleSelect creates a simple select component
func NewSimpleSelect(label string, options []SelectOption) SimpleSelect {
	return SimpleSelect{
		Label:       label,
		Options:     options,
		SelectedIdx: 0,
	}
}

// Update handles key presses for simple select
func (ss *SimpleSelect) Update(msg tea.Msg) SimpleSelect {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if ss.SelectedIdx > 0 {
				ss.SelectedIdx--
			}
		case "down", "j":
			if ss.SelectedIdx < len(ss.Options)-1 {
				ss.SelectedIdx++
			}
		}
	}
	return *ss
}

// View renders the simple select
func (ss SimpleSelect) View() string {
	var b strings.Builder

	b.WriteString(styles.BoldStyle.Render(ss.Label) + "\n\n")

	for i, opt := range ss.Options {
		cursor := "  "
		optStyle := styles.ListItemStyle
		if i == ss.SelectedIdx {
			cursor = styles.IconArrow + " "
			optStyle = styles.SelectedItemStyle
		}

		b.WriteString(fmt.Sprintf("%s%s\n", cursor, optStyle.Render(opt.Title())))
		if opt.Description() != "" {
			desc := styles.MutedStyle.Render("  " + opt.Description())
			b.WriteString(fmt.Sprintf("   %s\n", desc))
		}
	}

	return b.String()
}

// SelectedValue returns the selected value
func (ss SimpleSelect) SelectedValue() string {
	if ss.SelectedIdx >= 0 && ss.SelectedIdx < len(ss.Options) {
		return ss.Options[ss.SelectedIdx].Value()
	}
	return ""
}

// SelectedOption returns the selected option
func (ss SimpleSelect) SelectedOption() *SelectOption {
	if ss.SelectedIdx >= 0 && ss.SelectedIdx < len(ss.Options) {
		return &ss.Options[ss.SelectedIdx]
	}
	return nil
}
