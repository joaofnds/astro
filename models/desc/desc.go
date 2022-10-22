package desc

import (
	"astro/habit"
	"astro/logger"
	"astro/state"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

type Desc struct {
	habit       *habit.Habit
	parent      tea.Model
	description string
	textarea    textarea.Model
}

func NewDesc(habit *habit.Habit, parent tea.Model) Desc {
	ta := textarea.New()
	ta.Placeholder = "Check-in description"
	ta.Focus()
	ta.SetWidth(80)

	return Desc{habit: habit, parent: parent, textarea: ta}
}

func (m Desc) Init() tea.Cmd {
	return textarea.Blink
}

func (m Desc) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		taCmd tea.Cmd
	)

	m.textarea, taCmd = m.textarea.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m.parent, nil
		case tea.KeyEnter:
			hab, err := habit.Client.CheckIn(m.habit.ID, m.textarea.Value())
			if err != nil {
				logger.Error.Printf("failed to add activity: %v", err)
			} else {
				state.SetHabit(hab)
			}
			return m.parent, nil
		}
	}

	return m, taCmd
}

func (m Desc) View() string {
	return "Check-In Description: \n\n" + m.textarea.View()
}
