package list

import (
	"astro/msgs"
	"strings"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
)

type (
	errMsg error
)

type addModel struct {
	input textinput.Model
	err   error
}

func newAddInput() addModel {
	input := textinput.New()
	input.Placeholder = "Read"
	input.Focus()
	input.CharLimit = 50
	input.SetWidth(20)

	return addModel{input: input}
}

func (m addModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m addModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
						return createHabitSubmit{Name: trimmed}
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

func (m addModel) View() tea.View {
	return tea.NewView("What is the name of your new habit?\n" + m.input.View() + "\n\n(esc to quit)")
}
