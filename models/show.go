package models

import (
	"astroapp/config"
	"astroapp/habit"
	"astroapp/histogram"
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	style = lipgloss.NewStyle().Padding(0, 2)
	name  = lipgloss.NewStyle().
		Background(lipgloss.Color("#5F5FD7")).
		Foreground(lipgloss.Color("#FFFFD7")).
		Padding(0, 1)
)

type Show struct {
	habit    habit.Habit
	parent   tea.Model
	selected int
	t        time.Time
}

func NewShow(habit habit.Habit, parent tea.Model) Show {
	t, _ := histogram.TimeFrame()
	return Show{habit, parent, config.TimeFrameInDays - 1, t}
}

func (m Show) Init() tea.Cmd {
	return nil
}

func (m Show) View() string {
	s := new(strings.Builder)

	s.WriteString(name.Render(m.habit.Name) + "\n")
	s.WriteString(histogram.Histogram(m.t, m.habit, m.selected))
	s.WriteString(activitiesOnDate(m.habit, m.t.AddDate(0, 0, m.selected)))

	return style.Render(s.String())
}

func (m Show) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "l": // right
			if m.selected+7 < config.TimeFrameInDays {
				m.selected += 7
			}
		case "h": // left
			if m.selected-7 >= 0 {
				m.selected -= 7
			}
		case "j": // down
			if m.selected+1 < config.TimeFrameInDays {
				m.selected++
			}
		case "k": // up
			if m.selected > 0 {
				m.selected--
			}
		case "q":
			if m.parent == nil {
				return m, tea.Quit
			}
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
	return fmt.Sprintf("%d %s on %s\n", count, w, t.Format(config.TimeFormat))
}
