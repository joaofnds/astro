package show

import (
	"astro/config"
	"astro/date"
	"astro/habit"
	"astro/histogram"
	"astro/logger"
	"astro/models/textinput"
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
	h := help.New()
	h.Width = config.Width
	return Show{
		habit:    habit,
		parent:   parent,
		selected: selected,
		t:        t,
		help:     h,
		keys:     NewKeymap(),
	}
}

func (m Show) selectedDate() time.Time {
	return m.t.AddDate(0, 0, m.selected)
}

func (m Show) Init() tea.Cmd {
	return nil
}

func (m Show) View() string {
	var s strings.Builder
	s.Grow(11_000)

	s.WriteString(name.Render(fmt.Sprintf("%s - %s streak", m.habit.Name, habit.Streak(m.habit.Activities))) + "\n")
	s.WriteString(histogram.Histogram(m.t, m.habit.Activities, m.selected))
	s.WriteString(habit.ActivitiesOnDate(m.habit.Activities, m.selectedDate()))
	s.WriteString(timeline(m.habit, m.selectedDate()))
	s.WriteString(m.help.View(m.keys))

	return style.Render(s.String())
}

func (m Show) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case textinput.Submit:
		switch msg.Key {
		case "checkin":
			if hab, err := habit.Client.CheckIn(m.habit.ID, msg.Value); err != nil {
				logger.Error.Printf("failed to check: %v", err)
			} else {
				state.SetHabit(hab)
			}
		case "checkin-edit":
			var activity *habit.Activity
			for _, a := range m.habit.Activities {
				if a.ID == msg.ID {
					activity = &a
				}
			}
			if activity == nil {
				break
			}
			activity.Desc = msg.Value
			if err := habit.Client.UpdateActivity(*m.habit, *activity); err != nil {
				logger.Error.Printf("failed to updated activity: %v", err)
			}
			state.UpdateActivity(m.habit, activity)
		}

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.CheckIn):
			hab, err := habit.Client.CheckIn(m.habit.ID, "")
			if err != nil {
				logger.Error.Printf("failed to add activity: %v", err)
			} else {
				state.SetHabit(hab)
			}

		case key.Matches(msg, m.keys.VCheckIn):
			return textinput.New(m, "Check-In Description", "", "checkin", m.habit.ID), nil

		case key.Matches(msg, m.keys.Edit):
			if activity, err := m.habit.LatestActivityOnDate(m.selectedDate()); err == nil {
				return textinput.New(m, "New Description", activity.Desc, "checkin-edit", activity.ID), nil
			}

		case key.Matches(msg, m.keys.Delete):
			if !date.SameDay(m.selectedDate(), date.Today()) {
				break
			}

			activity, err := m.habit.LatestActivityOnDate(m.selectedDate())
			if err != nil {
				break // no activity on date
			}
			if err := habit.Client.DeleteActivity(*m.habit, activity); err != nil {
				logger.Debug.Printf("failed to delete activity: %v", err)
				break
			}
			state.DeleteActivity(m.habit, activity)

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

func timeline(h *habit.Habit, t time.Time) string {
	var s strings.Builder

	s.WriteString("\n")

	for _, a := range h.Activities {
		if !date.SameDay(a.CreatedAt, t) {
			continue
		}

		if a.Desc != "" {
			s.WriteString(a.CreatedAt.Local().Format(config.TimeFormat) + "\n\t" + a.Desc + "\n")
		}
	}

	return s.String()
}
