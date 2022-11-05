package list

import (
	"astro/state"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type (
	errMsg  error
)

type model struct {
	parent List
	input  textinput.Model
	err    error
}

func newAddInput(parent List) model {
	input := textinput.New()
	input.Placeholder = "Read"
	input.Focus()
	input.CharLimit = 50
	input.Width = 20

	return model{parent: parent, input: input, err: nil}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc, tea.KeyCtrlC:
			return m.parent, nil
		case tea.KeyEnter:
			h := state.Add(m.input.Value())
			m.parent.list.SetItems(toItems(state.Habits()))
			m.parent.list.ResetSelected()
			return m.parent, m.parent.list.NewStatusMessage("Added " + h.Name)
		}

	case errMsg:
		m.err = msg
		return m, nil
	}

	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return "What is the name of your new habit?\n" + m.input.View() + "\n\n(esc to quit)"
}
