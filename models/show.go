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
	t        time.Time
}

func NewShowModel(habit habit.Habit, parent tea.Model) Show {
	t, _ := histogram.TimeFrame()
	return Show{habit, parent, histogram.TimeFrameInDays - 1, t}
}

func (m Show) Init() tea.Cmd {
	return nil
}

func (m Show) View() string {
	s := new(strings.Builder)

	fmt.Fprintf(s, "Habit: %s\n", m.habit.Name)
	s.WriteString(histogram.Histogram(m.t, m.habit, m.selected))
	s.WriteString(activitiesOnDate(m.habit, m.t.AddDate(0, 0, m.selected)))

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

func activitiesOnDate(h habit.Habit, t time.Time) string {
	var count int
	for _, a := range h.Activites {
		if histogram.SameDay(a.CreatedAt, t) {
			count++
		}
	}
	w := "activities"
	if count == 1 {
		w = "activity"
	}
	return fmt.Sprintf("%d %s on %s\n", count, w, t.Format("Jan 02, 2006"))
}
