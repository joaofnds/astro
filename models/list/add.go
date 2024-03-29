package list

import (
	"astro/msgs"
	"astro/state"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type (
	errMsg error
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
			trimmed := strings.TrimSpace(m.input.Value())
			if trimmed == "" {
				break
			}
			h := state.Add(trimmed)
			return m.parent, tea.Batch(
				m.parent.list.NewStatusMessage("Added "+h.Name),
				msgs.UpdateList,
			)
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
