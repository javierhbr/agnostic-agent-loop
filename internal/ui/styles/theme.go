package styles

import "github.com/charmbracelet/lipgloss"

// Color palette
var (
	Primary   = lipgloss.Color("#7D56F4") // Purple
	Secondary = lipgloss.Color("#00D9FF") // Cyan
	Success   = lipgloss.Color("#04B575") // Green
	Warning   = lipgloss.Color("#FFAA00") // Orange
	Error     = lipgloss.Color("#FF4672") // Red
	Muted     = lipgloss.Color("#6C757D") // Gray
	Text      = lipgloss.Color("#FFFFFF") // White
	Subtle    = lipgloss.Color("#A0A0A0") // Light gray
)

// Text styles
var (
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(Primary).
			MarginBottom(1)

	SubtitleStyle = lipgloss.NewStyle().
			Foreground(Secondary).
			Italic(true)

	HelpStyle = lipgloss.NewStyle().
			Foreground(Muted).
			Italic(true)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(Error).
			Bold(true)

	SuccessStyle = lipgloss.NewStyle().
			Foreground(Success).
			Bold(true)

	WarningStyle = lipgloss.NewStyle().
			Foreground(Warning).
			Bold(true)

	MutedStyle = lipgloss.NewStyle().
			Foreground(Muted)

	BoldStyle = lipgloss.NewStyle().
			Bold(true)
)

// Component styles
var (
	InputStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Primary).
			Padding(0, 1)

	InputFocusedStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(Secondary).
				Padding(0, 1)

	InputErrorStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Error).
			Padding(0, 1)

	BoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Primary).
			Padding(1, 2)

	CardStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Secondary).
			Padding(1, 2).
			MarginTop(1).
			MarginBottom(1)

	ListItemStyle = lipgloss.NewStyle().
			PaddingLeft(2)

	SelectedItemStyle = lipgloss.NewStyle().
				Foreground(Secondary).
				Bold(true).
				PaddingLeft(0)

	ButtonStyle = lipgloss.NewStyle().
			Foreground(Text).
			Background(Primary).
			Padding(0, 3).
			MarginRight(2)

	ButtonActiveStyle = lipgloss.NewStyle().
				Foreground(Text).
				Background(Secondary).
				Padding(0, 3).
				MarginRight(2)

	TabStyle = lipgloss.NewStyle().
			Foreground(Muted).
			Padding(0, 2)

	SelectedTabStyle = lipgloss.NewStyle().
				Foreground(Secondary).
				Bold(true).
				Padding(0, 2).
				BorderBottom(true).
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(Secondary)
)

// Layout styles
var (
	ContainerStyle = lipgloss.NewStyle().
			Padding(1, 2)

	HeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(Primary).
			BorderBottom(true).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(Primary).
			Padding(0, 1, 1, 1)

	FooterStyle = lipgloss.NewStyle().
			Foreground(Muted).
			BorderTop(true).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(Muted).
			Padding(1, 1, 0, 1)
)

// Icons and symbols
const (
	IconCheckmark = "✓"
	IconCross     = "✗"
	IconArrow     = "→"
	IconSpinner   = "⠋"
	IconPending   = "○"
	IconProgress  = "◐"
	IconBullet    = "•"
	IconPrompt    = "›"
)

// RenderTitle renders a styled title
func RenderTitle(title string) string {
	return TitleStyle.Render(title)
}

// RenderSubtitle renders a styled subtitle
func RenderSubtitle(subtitle string) string {
	return SubtitleStyle.Render(subtitle)
}

// RenderHelp renders help text
func RenderHelp(help string) string {
	return HelpStyle.Render(help)
}

// RenderError renders an error message
func RenderError(msg string) string {
	return ErrorStyle.Render(IconCross + " " + msg)
}

// RenderSuccess renders a success message
func RenderSuccess(msg string) string {
	return SuccessStyle.Render(IconCheckmark + " " + msg)
}

// RenderWarning renders a warning message
func RenderWarning(msg string) string {
	return WarningStyle.Render("! " + msg)
}
