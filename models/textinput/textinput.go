package textinput

import (
	"astro/config"
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
	parent   tea.Model
	textarea textarea.Model
	key      string
	id       string
	prompt   string
}

func New(parent tea.Model, prompt, initialValue, key, id string) Model {
	ta := textarea.New()
	ta.SetValue(initialValue)
	ta.Focus()
	ta.SetWidth(config.Width)

	return Model{parent: parent, textarea: ta, key: key, id: id, prompt: prompt}
}

func (m Model) Init() tea.Cmd {
	return textarea.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m.parent, nil
		case "enter":
			trimmed := strings.TrimSpace(m.textarea.Value())
			if trimmed == "" {
				break
			}
			return m.parent, func() tea.Msg { return Submit{Key: m.key, ID: m.id, Value: trimmed} }
		}
	}

	var taCmd tea.Cmd
	m.textarea, taCmd = m.textarea.Update(msg)
	return m, taCmd
}

func (m Model) View() tea.View {
	v := tea.NewView(m.prompt + "\n" + m.textarea.View())
	v.AltScreen = true
	return v
}
