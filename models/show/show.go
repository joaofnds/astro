package show

import (
	"astro/api"
	"astro/config"
	"astro/date"
	"astro/domain"
	"astro/models/textinput"
	"astro/msgs"
	"astro/util"
	"context"
	"fmt"
	"strings"
	"time"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

var (
	style = lipgloss.NewStyle().Padding(0, 2)
	name  = lipgloss.NewStyle().
		Background(lipgloss.Color("#5F5FD7")).
		Foreground(lipgloss.Color("#FFFFD7")).
		Padding(0, 1)
)

type Show struct {
	client        *api.Client
	habit         *domain.Habit
	habitSnapshot *domain.Habit // Pre-mutation snapshot for rollback
	selected      int
	t             time.Time
	help          help.Model
	keys          keymap
	width         int
	statusMsg     string
	cancelOp      context.CancelFunc
}

// snapshotHabit returns a deep copy of h suitable for rollback.
// The Activities slice is copied to avoid shared backing arrays.
func snapshotHabit(h *domain.Habit) *domain.Habit {
	cp := *h
	cp.Activities = make([]domain.Activity, len(h.Activities))
	copy(cp.Activities, h.Activities)
	return &cp
}

func NewShow(client *api.Client, habit *domain.Habit, width int) Show {
	t, _ := date.TimeFrame()
	selected := date.DiffInDays(t, date.Today())
	h := help.New()
	h.SetWidth(width)
	return Show{
		client:   client,
		habit:    habit,
		selected: selected,
		t:        t,
		help:     h,
		keys:     NewKeymap(),
		width:    width,
	}
}

func (m Show) selectedDate() time.Time {
	return m.t.AddDate(0, 0, m.selected)
}

func (m Show) Init() tea.Cmd {
	return nil
}

func (m Show) View() tea.View {
	var s strings.Builder
	s.Grow(11_000)

	s.WriteString(name.Render(domain.Digest(m.habit.Name, m.habit.Activities)) + "\n")
	s.WriteString(domain.Histogram(m.t, m.habit.Activities, m.selected))
	s.WriteString(domain.ActivitiesOnDate(m.habit.Activities, m.selectedDate()))
	s.WriteString(timeline(m.habit, m.selectedDate()))
	if m.statusMsg != "" {
		s.WriteString("\n" + m.statusMsg + "\n")
	}
	s.WriteString(m.help.View(m.keys))

	return tea.NewView(style.Render(s.String()))
}

func (m Show) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case msgs.APIErrorMsg:
		errStr := fmt.Sprintf("\u2717 %s: %s", msg.Op, msg.Err)
		if msg.Op == "check in" && m.habitSnapshot != nil {
			m.habit = m.habitSnapshot
			m.habitSnapshot = nil
			errStr += " \u2014 restored"
		}
		m.statusMsg = errStr
		return m, msgs.ClearStatusAfter(3 * time.Second)

	case msgs.ClearStatusMsg:
		m.statusMsg = ""
		return m, nil

	case textinput.Submit:
		switch msg.Key {
		case "checkin":
			m.habitSnapshot = snapshotHabit(m.habit)
			checkInTime := m.checkInDate()
			m.habit.Activities = append(m.habit.Activities, domain.Activity{
				ID:        "pending",
				Desc:      msg.Value,
				CreatedAt: checkInTime,
			})
			m.statusMsg = "Checking in..."
			ctx, cancel := context.WithCancel(context.Background())
			m.cancelOp = cancel
			return m, msgs.CheckIn(ctx, m.client, m.habit.ID, msg.Value, checkInTime)
		case "checkin-edit":
			m.statusMsg = "Updating..."
			ctx, cancel := context.WithCancel(context.Background())
			m.cancelOp = cancel
			return m, msgs.UpdateActivity(ctx, m.client, m.habit.ID, msg.ID, msg.Value)
		}

	case msgs.CheckInResultMsg:
		if msg.Habit.ID == m.habit.ID {
			m.habit = msg.Habit
			m.habitSnapshot = nil
			m.statusMsg = "Checked in"
			return m, msgs.ClearStatusAfter(3 * time.Second)
		}

	case msgs.ActivityUpdatedMsg:
		if msg.HabitID == m.habit.ID {
			for i, a := range m.habit.Activities {
				if a.ID == msg.ActivityID {
					m.habit.Activities[i].Desc = msg.Desc
				}
			}
			m.statusMsg = "Updated"
			return m, msgs.ClearStatusAfter(3 * time.Second)
		}

	case msgs.ActivityDeletedMsg:
		if msg.HabitID == m.habit.ID {
			for i, a := range m.habit.Activities {
				if a.ID == msg.ActivityID {
					m.habit.Activities = append(m.habit.Activities[:i], m.habit.Activities[i+1:]...)
					break
				}
			}
			m.statusMsg = "Deleted"
			return m, msgs.ClearStatusAfter(3 * time.Second)
		}

	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, m.keys.CheckIn):
			if m.selectedDate().After(time.Now()) {
				break
			}
			m.habitSnapshot = snapshotHabit(m.habit)
			checkInTime := m.checkInDate()
			m.habit.Activities = append(m.habit.Activities, domain.Activity{
				ID:        "pending",
				Desc:      "",
				CreatedAt: checkInTime,
			})
			m.statusMsg = "Checking in..."
			ctx, cancel := context.WithCancel(context.Background())
			m.cancelOp = cancel
			return m, msgs.CheckIn(ctx, m.client, m.habit.ID, "", checkInTime)

		case key.Matches(msg, m.keys.VCheckIn):
			if m.selectedDate().After(time.Now()) {
				break
			}
			return m, msgs.PushScreen(textinput.New("Check-In Description", "", "checkin", m.habit.ID, m.width))

		case key.Matches(msg, m.keys.Edit):
			if activity, err := m.habit.LatestActivityOnDate(m.selectedDate()); err == nil && activity.ID != "pending" {
				return m, msgs.PushScreen(textinput.New("New Description", activity.Desc, "checkin-edit", activity.ID, m.width))
			}

		case key.Matches(msg, m.keys.Delete):
			activity, err := m.habit.LatestActivityOnDate(m.selectedDate())
			if err != nil || activity.ID == "pending" {
				break
			}
			m.statusMsg = "Deleting activity..."
			ctx, cancel := context.WithCancel(context.Background())
			m.cancelOp = cancel
			return m, msgs.DeleteActivity(ctx, m.client, m.habit.ID, activity.ID)

		// ClearScreen forces a full sequential redraw on navigation.
		// The v2 renderer's differential updates miscalculate cursor
		// positions for emoji characters (see original comment).
		case key.Matches(msg, m.keys.Up):
			m.selected = util.Max(m.selected-1, 0)
			return m, tea.ClearScreen

		case key.Matches(msg, m.keys.Down):
			m.selected = util.Min(m.selected+1, config.TimeFrameInDays-1)
			return m, tea.ClearScreen

		case key.Matches(msg, m.keys.Left):
			m.selected = util.Max(m.selected-7, 0)
			return m, tea.ClearScreen

		case key.Matches(msg, m.keys.Right):
			m.selected = util.Min(m.selected+7, config.TimeFrameInDays-1)
			return m, tea.ClearScreen

		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll

		case key.Matches(msg, m.keys.Quit):
			if m.cancelOp != nil {
				m.cancelOp()
				m.cancelOp = nil
			}
			m.habitSnapshot = nil
			return m, msgs.PopScreen()
		}
	}

	return m, nil
}

func (m Show) checkInDate() time.Time {
	return date.CombineDateWithTime(m.selectedDate(), time.Now().Local())
}

func timeline(h *domain.Habit, t time.Time) string {
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
