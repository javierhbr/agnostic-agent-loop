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

// MultiSelect is a select component that allows toggling multiple options
type MultiSelect struct {
	Label      string
	Options    []SelectOption
	CursorIdx  int
	Selected   map[int]bool
	MaxVisible int // max items to show at once; 0 = show all
	offset     int // first visible item index
}

// NewMultiSelect creates a multi-select component
func NewMultiSelect(label string, options []SelectOption) MultiSelect {
	return MultiSelect{
		Label:     label,
		Options:   options,
		CursorIdx: 0,
		Selected:  make(map[int]bool),
	}
}

// Update handles key presses for multi-select
func (ms *MultiSelect) Update(msg tea.Msg) MultiSelect {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if ms.CursorIdx > 0 {
				ms.CursorIdx--
				ms.ensureCursorVisible()
			}
		case "down", "j":
			if ms.CursorIdx < len(ms.Options)-1 {
				ms.CursorIdx++
				ms.ensureCursorVisible()
			}
		case " ":
			ms.Selected[ms.CursorIdx] = !ms.Selected[ms.CursorIdx]
		}
	}
	return *ms
}

// SetMaxVisible sets the maximum number of visible items based on available height.
// Each item uses 2 lines (title + description). Pass the total available lines.
func (ms *MultiSelect) SetMaxVisible(availableLines int) {
	// Each option takes 2 lines (title + description), plus 1 for the label
	if availableLines > 0 {
		ms.MaxVisible = max((availableLines-2)/2, 1) // subtract label + blank line
	}
}

func (ms *MultiSelect) ensureCursorVisible() {
	if ms.MaxVisible <= 0 || ms.MaxVisible >= len(ms.Options) {
		return
	}
	if ms.CursorIdx < ms.offset {
		ms.offset = ms.CursorIdx
	}
	if ms.CursorIdx >= ms.offset+ms.MaxVisible {
		ms.offset = ms.CursorIdx - ms.MaxVisible + 1
	}
}

// View renders the multi-select
func (ms MultiSelect) View() string {
	var b strings.Builder

	b.WriteString(styles.BoldStyle.Render(ms.Label) + "\n\n")

	start := ms.offset
	end := len(ms.Options)
	if ms.MaxVisible > 0 && ms.MaxVisible < len(ms.Options) {
		end = min(start+ms.MaxVisible, len(ms.Options))
		if start > 0 {
			b.WriteString(styles.MutedStyle.Render("  ↑ more") + "\n")
		}
	}

	for i := start; i < end; i++ {
		opt := ms.Options[i]
		cursor := "  "
		if i == ms.CursorIdx {
			cursor = styles.IconArrow + " "
		}

		check := "[ ]"
		if ms.Selected[i] {
			check = "[" + styles.IconCheckmark + "]"
		}

		optStyle := styles.ListItemStyle
		if i == ms.CursorIdx {
			optStyle = styles.SelectedItemStyle
		}

		fmt.Fprintf(&b, "%s%s %s\n", cursor, check, optStyle.Render(opt.Title()))
		if opt.Description() != "" {
			desc := styles.MutedStyle.Render("     " + opt.Description())
			fmt.Fprintf(&b, "   %s\n", desc)
		}
	}

	if ms.MaxVisible > 0 && end < len(ms.Options) {
		b.WriteString(styles.MutedStyle.Render("  ↓ more") + "\n")
	}

	return b.String()
}

// SelectedValues returns the values of all selected options
func (ms MultiSelect) SelectedValues() []string {
	var values []string
	for i, opt := range ms.Options {
		if ms.Selected[i] {
			values = append(values, opt.Value())
		}
	}
	return values
}

// HasSelection returns true if at least one option is selected
func (ms MultiSelect) HasSelection() bool {
	for _, v := range ms.Selected {
		if v {
			return true
		}
	}
	return false
}
