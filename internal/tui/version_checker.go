package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type spinnerValue struct {
	value string
	err   error
	done  bool
	ch    chan int
}

type SpinnerModel struct {
	spinner spinner.Model
	label   string
	value   *spinnerValue
	err     error
}

func spinnerModel(label string, val *spinnerValue) SpinnerModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = HighlightStyle

	return SpinnerModel{
		spinner: s,
		label:   label,
		value:   val,
	}
}

func (m SpinnerModel) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m SpinnerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	select {
	case <-m.value.ch:
		m.value.done = true
		return m, tea.Quit
	default:
		// continue
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		default:
			return m, nil
		}

	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
}

func (m SpinnerModel) View() string {
	if m.err != nil {
		return m.err.Error()
	}

	if m.value.done {
		return "\n\033[F\033[K"
	}

	return fmt.Sprintf("%s %s Press q to quit\n", m.spinner.View(), m.label)
}

func NewSpinner(label string, while func() (string, error)) (string, error) {
	// setup graceful shutdown handler
	done := make(chan int, 1)
	defer close(done)

	value := &spinnerValue{
		ch:   done,
		done: false,
	}

	go func() {
		val, err := while()
		value.value = val
		value.err = err
		done <- 1
	}()

	p := tea.NewProgram(spinnerModel(label, value))
	if _, err := p.Run(); err != nil {
		return "", err
	}

	return value.value, value.err
}
