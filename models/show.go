package models

import (
	"astroapp/habit"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type Show struct {
	habit  habit.Habit
	parent tea.Model
	cursor int
}

func NewShowModel(habit habit.Habit, parent tea.Model) Show {
	return Show{habit, parent, 0}
}

func (m Show) Init() tea.Cmd {
	return nil
}

func (m Show) View() string {
	s := fmt.Sprintf("habits: %s\n", m.habit.Name)
	s += fmt.Sprintf("id: %d\n", m.habit.Id)

	s += "activities:\n"
	for i, activity := range m.habit.Activites {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		s += fmt.Sprintf("%s %s\n", cursor, activity.CreatedAt.Format("06/01/02 15:04"))
	}

	s += "\npress 'q' to go back\n"
	return s
}

func (m Show) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "j":
			if m.cursor < len(m.habit.Activites)-1 {
				m.cursor++
			}
		case "q":
			return m.parent, nil
		case "ctrl+c", "ctrl+d":
			return m, tea.Quit
		}
	}

	return m, nil
}
