package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mreliasen/scrolls-cli/internal/file_types"
)

type selectionResult struct {
	value  string
	cancel bool
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
			m.selection.cancel = true
			return m, tea.Quit

		case "enter":
			s := m.table.SelectedRow()
			if len(s) > 0 {
				m.selection.value = s[0]
				return m, tea.Quit
			}

			r := m.table.Rows()
			if len(r) == 1 {
				m.table.SetCursor(0)
				m.selection.value = r[0][0]
				return m, tea.Quit
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
		if strings.Contains(r[0], term) || strings.Contains(r[1], term) {
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
		helpTextStyle.Render("\n  Enter:  Choose selection & save\n  TAB:    Switch between search and list\n  ↑↓/kj:  Move up/down the list\n  CTRL+C: Exit without selection") +
		"\n" +
		baseStyle.Render(
			m.table.View(),
		)
}

func NewSelector(initial string) (string, bool) {
	columns := []table.Column{
		{Title: "File Type", Width: 25},
		{Title: "Extension", Width: 10},
	}

	rows := []table.Row{}

	i := 0
	tSelection := -1
	for bin, cfg := range file_types.ExecList {
		if initial == bin {
			rows = append([]table.Row{{bin, cfg.Ext}}, rows...)
			tSelection = 0
		} else {
			rows = append(rows, table.Row{bin, cfg.Ext})
		}
		i++
	}

	search := textinput.New()
	search.Placeholder = "Search"
	search.CharLimit = 20
	search.Width = 20

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(false),
		table.WithHeight(10),
	)

	if tSelection > -1 {
		t.Focus()
		t.SetCursor(tSelection)
	} else {
		search.Focus()
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

	selection := &selectionResult{
		cancel: false,
		value:  initial,
	}
	m := model{t, search, rows, selection}

	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		return "", true
	}

	return selection.value, selection.cancel
}
