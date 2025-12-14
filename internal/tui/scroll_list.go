package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mreliasen/scrolls-cli/internal/library"
)

type listResult struct {
	value  *library.Scroll
	cancel bool
}

type listModel struct {
	table     table.Model
	list      []*library.Scroll
	selection *listResult
}

func (m listModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m listModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			m.selection.cancel = true
			return m, tea.Quit

		case "enter":
			s := m.table.Cursor()
			m.selection.value = m.list[s]
			return m, tea.Quit
		}
	}

	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m listModel) View() string {
	return "\n" +
		"Heres your scrolls!" +
		"\n" +
		baseStyle.Render(
			m.table.View(),
		) +
		helpTextStyle.Render("\n  Enter:  Manage selected scroll\n  ↑↓/kj:  Move up/down the list\n  CTRL+C: Exit without selection")
}

func NewScrollList(scrolls []*library.Scroll) (*library.Scroll, bool) {
	if len(scrolls) == 0 {
		fmt.Println("No scrolls found.")
		return nil, true
	}

	columns := []table.Column{
		{Title: "Name", Width: 20},
		{Title: "Type", Width: 10},
	}

	rows := []table.Row{}

	for _, f := range scrolls {
		rows = append(rows, table.Row{
			f.Name(),
			f.Type(),
		})
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(25),
	)

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

	listRes := &listResult{
		cancel: false,
		value:  nil,
	}
	m := listModel{t, scrolls, listRes}

	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		return nil, true
	}

	return listRes.value, listRes.cancel
}
