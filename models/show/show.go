package show

import (
	"astro/config"
	"astro/date"
	"astro/habit"
	"astro/histogram"
	"astro/logger"
	"astro/state"
	"astro/util"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
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
)

type Show struct {
	habit    *habit.Habit
	parent   tea.Model
	selected int
	t        time.Time
	help     help.Model
	keys     keymap
}

func NewShow(habit *habit.Habit, parent tea.Model) Show {
	t, _ := date.TimeFrame()
	selected := date.DiffInDays(t, date.Today())
	return Show{
		habit:    habit,
		parent:   parent,
		selected: selected,
		t:        t,
		help:     help.New(),
		keys:     NewKeymap(),
	}
}

func (m Show) Init() tea.Cmd {
	return nil
}

func (m Show) View() string {
	s := new(strings.Builder)
	s.Grow(11_000)

	s.WriteString(name.Render(m.habit.Name) + "\n")
	s.WriteString(histogram.Histogram(m.t, *m.habit, m.selected))
	s.WriteString(activitiesOnDate(m.habit, m.t.AddDate(0, 0, m.selected)))
	s.WriteString(m.help.View(m.keys))

	return style.Render(s.String())
}

func (m Show) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.help.Width = msg.Width
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.CheckIn):
			habit, err := habit.Client.CheckIn(m.habit.ID)
			if err != nil {
				logger.Error.Printf("failed to add activity: %v", err)
			} else {
				state.SetHabit(habit)
			}
		case key.Matches(msg, m.keys.Up):
			m.selected = util.Max(m.selected-1, 0)
		case key.Matches(msg, m.keys.Down):
			m.selected = util.Min(m.selected+1, config.TimeFrameInDays-1)
		case key.Matches(msg, m.keys.Left):
			m.selected = util.Max(m.selected-7, 0)
		case key.Matches(msg, m.keys.Right):
			m.selected = util.Min(m.selected+7, config.TimeFrameInDays-1)
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll // FIX: only works after resizing
		case key.Matches(msg, m.keys.Quit):
			if m.parent == nil {
				return m, tea.Quit
			}
			return m.parent, nil
		}
	}

	return m, nil
}

func activitiesOnDate(h *habit.Habit, t time.Time) string {
	var count int
	for _, a := range h.Activities {
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
