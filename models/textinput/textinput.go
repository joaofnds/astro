package textinput

import (
	"astro/msgs"
	"strings"

	"charm.land/bubbles/v2/textarea"
	tea "charm.land/bubbletea/v2"
)

type Submit struct {
	Key   string
	ID    string
	Value string
}

type Model struct {
	textarea textarea.Model
	key      string
	id       string
	prompt   string
}

func New(prompt, initialValue, key, id string, width int) Model {
	ta := textarea.New()
	ta.SetValue(initialValue)
	ta.Focus()
	ta.SetWidth(width)

	return Model{textarea: ta, key: key, id: id, prompt: prompt}
}

func (m Model) Init() tea.Cmd {
	return textarea.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, msgs.PopScreen()
		case "enter":
			trimmed := strings.TrimSpace(m.textarea.Value())
			if trimmed == "" {
				break
			}
			return m, func() tea.Msg {
				return msgs.PopScreenMsg{
					Cmd: func() tea.Msg {
						return Submit{Key: m.key, ID: m.id, Value: trimmed}
					},
				}
			}
		}
	}

	var taCmd tea.Cmd
	m.textarea, taCmd = m.textarea.Update(msg)
	return m, taCmd
}

func (m Model) View() tea.View {
	return tea.NewView(m.prompt + "\n" + m.textarea.View())
}
