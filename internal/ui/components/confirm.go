package components

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/javierbenavides/agentic-agent/internal/ui/styles"
)

// Confirm is a yes/no confirmation component
type Confirm struct {
	Label    string
	Selected bool // true = yes, false = no
}

// NewConfirm creates a new confirmation component
func NewConfirm(label string, defaultYes bool) Confirm {
	return Confirm{
		Label:    label,
		Selected: defaultYes,
	}
}

// Update handles confirm messages
func (c *Confirm) Update(msg tea.Msg) Confirm {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "left", "right", "h", "l":
			c.Selected = !c.Selected
		case "y":
			c.Selected = true
		case "n":
			c.Selected = false
		}
	}
	return *c
}

// View renders the confirmation
func (c Confirm) View() string {
	yesStyle := styles.ListItemStyle
	noStyle := styles.ListItemStyle

	if c.Selected {
		yesStyle = styles.SelectedItemStyle
	} else {
		noStyle = styles.SelectedItemStyle
	}

	yes := yesStyle.Render("[Yes]")
	no := noStyle.Render("[No]")

	cursor := ""
	if c.Selected {
		cursor = styles.IconArrow + " "
	}

	view := styles.BoldStyle.Render(c.Label) + "\n\n"

	if c.Selected {
		view += "  " + cursor + yes + "  " + no + "\n"
	} else {
		view += "  " + yes + "  " + cursor + no + "\n"
	}

	view += "\n" + styles.HelpStyle.Render("←/→ to toggle • y/n for yes/no • Enter to confirm")

	return view
}

// IsYes returns true if yes is selected
func (c Confirm) IsYes() bool {
	return c.Selected
}
