package models

import (
	"astroapp/config"
	"astroapp/date"
	"astroapp/habit"
	"astroapp/histogram"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	style = lipgloss.NewStyle().Padding(0, 2)
	name  = lipgloss.NewStyle().
		Background(lipgloss.Color("#5F5FD7")).
		Foreground(lipgloss.Color("#FFFFD7")).
		Padding(0, 1)
	keymap = KeyMap{
		Quit: key.NewBinding(
			key.WithKeys("q"),
			key.WithHelp("q", "quit"),
		),
		Up: key.NewBinding(
			key.WithKeys("k"),
			key.WithHelp("k", "+day"),
		),
		Down: key.NewBinding(
			key.WithKeys("j"),
			key.WithHelp("j", "-day"),
		),
		Left: key.NewBinding(
			key.WithKeys("h"),
			key.WithHelp("h", "-week"),
		),
		Right: key.NewBinding(
			key.WithKeys("l"),
			key.WithHelp("l", "+week"),
		),
	}
)

type KeyMap struct {
	Quit  key.Binding
	Up    key.Binding
	Down  key.Binding
	Left  key.Binding
	Right key.Binding
}

type Show struct {
	habit    habit.Habit
	parent   tea.Model
	selected int
	t        time.Time
}

func NewShow(habit habit.Habit, parent tea.Model) Show {
	t, _ := date.TimeFrame()
	selected := date.DiffInDays(t, date.Today())
	return Show{habit, parent, selected, t}
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
		switch {
		case key.Matches(msg, keymap.Up):
			if m.selected > 0 {
				m.selected--
			}
		case key.Matches(msg, keymap.Down):
			if m.selected+1 < config.TimeFrameInDays {
				m.selected++
			}
		case key.Matches(msg, keymap.Left):
			if m.selected-7 >= 0 {
				m.selected -= 7
			}
		case key.Matches(msg, keymap.Right):
			if m.selected+7 < config.TimeFrameInDays {
				m.selected += 7
			}
		case key.Matches(msg, keymap.Quit):
			if m.parent == nil {
				return m, tea.Quit
			}
			return m.parent, nil
		}
	}

	return m, nil
}

func activitiesOnDate(h habit.Habit, t time.Time) string {
	var count int
	for _, a := range h.Activites {
		if date.SameDay(a.CreatedAt, t) {
			count++
		}
	}
	w := "activities"
	if count == 1 {
		w = "activity"
	}
	return fmt.Sprintf("%d %s on %s\n", count, w, t.Format(config.TimeFormat))
}
