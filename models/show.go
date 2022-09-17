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

	fmt.Fprintf(s, "Habit: %s\n", m.habit.Name)

	s.WriteString(histogram.Histogram(m.habit, m.selected))

	var count int
	for _, a := range m.habit.Activites {
		if histogram.SameDay(a.CreatedAt, selectedDate) {
			count++
		}
	}
	w := "activities"
	if count == 1 {
		w = "activity"
	}
	fmt.Fprintf(s, "%d %s on %s\n", count, w, selectedDate.Format("Jan 02, 2006"))

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
