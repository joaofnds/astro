package desc

import (
	"astro/habit"
	"astro/logger"
	"astro/state"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

type EditDesc struct {
	habit       *habit.Habit
	activity    *habit.Activity
	parent      tea.Model
	description string
	textarea    textarea.Model
}

func NewEditEditDesc(habit *habit.Habit, activity *habit.Activity, parent tea.Model) EditDesc {
	ta := textarea.New()
	ta.SetValue(strings.TrimSpace(activity.Desc))
	ta.Focus()
	ta.SetWidth(80)

	return EditDesc{habit: habit, activity: activity, parent: parent, textarea: ta}
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
			m.activity.Desc = m.textarea.Value()
			if err := habit.Client.UpdateActivity(*m.habit, *m.activity); err != nil {
				logger.Error.Printf("failed to update activity: %v", err)
			}
			state.UpdateActivity(m.habit, m.activity)
			return m.parent, nil
		}
	}

	return m, taCmd
}

func (m EditDesc) View() string {
	return "Check-In Description: \n\n" + m.textarea.View()
}
