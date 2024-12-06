package tui

import "github.com/charmbracelet/lipgloss"

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

var helpTextStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#eaeaea"))

var (
	HighlightStyle = lipgloss.Color("#8454fc")
	ErrorStyle     = lipgloss.Color("#ff4e94")
	SuccessStyle   = lipgloss.Color("#9bf3af")
)
