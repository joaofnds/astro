package name

import (
	"astro/habit"
	"astro/logger"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

type EditDesc struct {
	habit    *habit.Habit
	name     string
	parent   tea.Model
	textarea textarea.Model
}

func NewEditName(habit *habit.Habit, parent tea.Model) EditDesc {
	ta := textarea.New()
	ta.SetValue(habit.Name)
	ta.Focus()
	ta.SetWidth(80)

	return EditDesc{habit: habit, parent: parent, textarea: ta}
}

func (m EditDesc) Init() tea.Cmd {
	return textarea.Blink
}

func (m EditDesc) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var taCmd tea.Cmd

	m.textarea, taCmd = m.textarea.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m.parent, nil
		case tea.KeyEnter:
			trimmed := strings.TrimSpace(m.textarea.Value())
			if trimmed == "" {
				break
			}
			m.habit.Name = trimmed
			if err := habit.Client.Update(m.habit); err != nil {
				logger.Error.Printf("failed to update habit: %v", err)
			}
			return m.parent, nil
		}
	}

	return m, taCmd
}

func (m EditDesc) View() string {
	return "Check-In Description: \n\n" + m.textarea.View()
}
