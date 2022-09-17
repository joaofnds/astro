package models

import (
	"astroapp/habit"
	"astroapp/histogram"
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type Show struct {
	habit    habit.Habit
	parent   tea.Model
	selected int
}

func NewShowModel(habit habit.Habit, parent tea.Model) Show {
	return Show{habit, parent, histogram.TimeFrameInDays - 1}
}

func (m Show) Init() tea.Cmd {
	return nil
}

func (m Show) View() string {
	t := histogram.EndOfWeek(histogram.TruncateDay(time.Now())).AddDate(0, 0, -histogram.TimeFrameInDays)
	selectedDate := t.AddDate(0, 0, m.selected)
	s := new(strings.Builder)

	fmt.Fprintf(s, "habits: %s\n", m.habit.Name)
	fmt.Fprintf(s, "id: %d\n", m.habit.Id)

	fmt.Fprintf(s, "activities:  %s\n", selectedDate.Format("2006-01-02"))

	s.WriteString(histogram.Histogram(m.habit, m.selected))

	s.WriteString("activities on this day\n")
	for _, a := range m.habit.Activites {
		if histogram.SameDay(a.CreatedAt, selectedDate) {
			fmt.Fprintf(s, " - %s\n", a.CreatedAt)
		}
	}

	s.WriteString("\npress 'q' to go back\n")
	return s.String()
}

func (m Show) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "l": // right
			if m.selected+7 < histogram.TimeFrameInDays {
				m.selected += 7
			}
		case "h": // left
			if m.selected-7 >= 0 {
				m.selected -= 7
			}
		case "j": // down
			if m.selected+1 < histogram.TimeFrameInDays {
				m.selected++
			}
		case "k": // up
			if m.selected > 0 {
				m.selected--
			}
		case "q":
			return m.parent, nil
		case "ctrl+c", "ctrl+d":
			return m, tea.Quit
		}
	}

	return m, nil
}
