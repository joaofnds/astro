package group

import (
	"astro/msgs"
	"strings"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
)

type (
	errMsg error
)

type model struct {
	input textinput.Model
	err   error
}

func NewAddGroup() model {
	input := textinput.New()
	input.Placeholder = "Health"
	input.Focus()
	input.CharLimit = 50
	input.SetWidth(20)

	return model{input: input, err: nil}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "esc", "ctrl+c":
			return m, msgs.PopScreen()
		case "enter":
			trimmed := strings.TrimSpace(m.input.Value())
			if trimmed == "" {
				break
			}
			return m, func() tea.Msg {
				return msgs.PopScreenMsg{
					Cmd: func() tea.Msg {
						return CreateGroupSubmit{Name: trimmed}
					},
				}
			}
		}

	case errMsg:
		m.err = msg
		return m, nil
	}

	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m model) View() tea.View {
	return tea.NewView("What is the name of your new group?\n" + m.input.View() + "\n\n(esc to quit)")
}
