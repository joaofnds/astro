package models

import (
	"astroapp/habit"
	"astroapp/histogram"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type Show struct {
	habit  habit.Habit
	parent tea.Model
}

func NewShowModel(habit habit.Habit, parent tea.Model) Show {
	return Show{habit, parent}
}

func (m Show) Init() tea.Cmd {
	return nil
}

func (m Show) View() string {
	s := fmt.Sprintf("habits: %s\n", m.habit.Name)
	s += fmt.Sprintf("id: %d\n", m.habit.Id)

	s += "activities:\n"

	s += histogram.Histogram(m.habit)

	s += "\npress 'q' to go back\n"
	return s
}

func (m Show) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m.parent, nil
		case "ctrl+c", "ctrl+d":
			return m, tea.Quit
		}
	}

	return m, nil
}
