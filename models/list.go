package models

import (
	"astroapp/habit"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type List struct {
	habits []habit.Habit
	cursor int
}

func NewList(habits []habit.Habit) List {
	return List{habits, 0}
}

func (m List) Init() tea.Cmd {
	return nil
}

func (m List) View() string {
	s := "habits:\n"

	for i, h := range m.habits {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		s += fmt.Sprintf("%s %s\n", cursor, h.Name)
	}

	return s
}

func (m List) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "ctrl+d", "q":
			return m, tea.Quit
		case "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "j":
			if m.cursor < len(m.habits)-1 {
				m.cursor++
			}
		case "enter", " ":
			return NewShowModel(m.habits[m.cursor], m), nil
		}
	}

	return m, nil
}
