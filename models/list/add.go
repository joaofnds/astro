package list

import (
	"astro/api"
	"astro/msgs"
	"context"
	"strings"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
)

type (
	errMsg error
)

type model struct {
	client *api.Client
	input  textinput.Model
	err    error
}

func newAddInput(client *api.Client) model {
	input := textinput.New()
	input.Placeholder = "Read"
	input.Focus()
	input.CharLimit = 50
	input.SetWidth(20)

	return model{client: client, input: input, err: nil}
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
					Cmd: msgs.CreateHabit(context.Background(), m.client, trimmed),
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
	return tea.NewView("What is the name of your new habit?\n" + m.input.View() + "\n\n(esc to quit)")
}
