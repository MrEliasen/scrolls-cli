package tui

import "github.com/charmbracelet/lipgloss"

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

var helpTextStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#eaeaea"))

var (
	HighlightStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#8454fc"))
	ErrorStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff4e94"))
	SuccessStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#9bf3af"))
)
