package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mreliasen/scrolls-cli/internal/scrolls/file_handler"
)

type selectionResult struct {
	value string
}

type model struct {
	table     table.Model
	input     textinput.Model
	rows      []table.Row
	selection *selectionResult
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "shift+tab":
			if m.table.Focused() {
				m.table.Blur()
				m.input.Focus()
			} else {
				m.input.Blur()
				m.table.Focus()
			}

		case "ctrl+c":
			return m, tea.Quit

		case "enter":
			if m.table.Focused() {
				s := m.table.SelectedRow()
				if len(s) > 0 {
					m.selection.value = s[0]
					return m, tea.Quit
				}
			}
		}
	}

	if m.table.Focused() {
		m.table, cmd = m.table.Update(msg)
	} else {
		m.input, cmd = m.input.Update(msg)

		m.table.SetRows(filterRows(m.rows, m.input.Value()))
		tModel, tCmd := m.table.Update(msg)
		m.table = tModel

		tea.Batch(cmd, tCmd)
	}

	return m, cmd
}

func filterRows(rows []table.Row, term string) []table.Row {
	filterd := []table.Row{}

	for _, r := range rows {
		if strings.Contains(r[0], term) {
			filterd = append(filterd, r)
		}
	}

	return filterd
}

func (m model) View() string {
	return "\n" +
		"Select the type associated with your scroll" +
		"\n" +
		m.input.View() +
		"\n" +
		baseStyle.Render(
			m.table.View(),
		) + helpTextStyle.Render(
		"\n  ↑↓/kj:  Move up/down the list\n  Enter:  Choose selection & save\n  TAB:    Switch between search and list\n  CTRL+C: Exit without selection",
	)
}

func NewSelector(initial string) string {
	columns := []table.Column{
		{Title: "File Type", Width: 25},
	}

	rows := []table.Row{}

	i := 0
	tSelection := -1
	for bin := range file_handler.ExecList {
		rows = append(rows, table.Row{bin})

		if initial == bin {
			tSelection = i
		}
		i++
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(false),
		table.WithHeight(10),
	)

	if tSelection > -1 {
		t.SetCursor(tSelection)
	}

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	search := textinput.New()
	search.Placeholder = "Search"
	search.Focus()
	search.CharLimit = 20
	search.Width = 20

	selection := &selectionResult{
		value: initial,
	}
	m := model{t, search, rows, selection}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		return ""
	}

	return selection.value
}
