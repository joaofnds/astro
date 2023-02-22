package group

import (
	"astro/config"
	"astro/date"
	"astro/habit"
	"astro/histogram"
	"astro/logger"
	"astro/models/listitem"
	"astro/models/show"
	"astro/models/textinput"
	"astro/msgs"
	"astro/state"
	"astro/util"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type List struct {
	group        *habit.Group
	parent       tea.Model
	list         list.Model
	km           binds
	t            time.Time
	selected     int
	lastSelected int
	onHist       bool
}

func NewShow(g *habit.Group, parent tea.Model) List {
	l := list.New(listitem.HabitsToItems(g.Habits), list.NewDefaultDelegate(), 0, 5)
	l.SetSize(config.Width, config.Height-9)

	km := newBinds()
	l.AdditionalShortHelpKeys = km.ToSlice
	l.Title = g.Name

	t, _ := date.TimeFrame()
	selected := date.DiffInDays(t, date.Today()) + config.TimeFrameInDays

	return List{t: t, selected: selected, group: g, parent: parent, list: l, km: km}
}

func (m List) Init() tea.Cmd {
	return nil
}

func (m List) View() string {
	activities := m.group.Activities()
	var s strings.Builder
	s.WriteString(histogram.Histogram(m.t, activities, m.selected))

	if m.selectedDate().After(date.Today()) {
		s.WriteString("\n")
	} else {
		s.WriteString(habit.ActivitiesOnDateTally(m.group.Habits, m.selectedDate()))
	}

	m.list.Title = fmt.Sprintf("%s - %s streak", m.group.Name, habit.Streak(m.group.Activities()))

	s.WriteString("\n")
	s.WriteString(m.list.View())
	return s.String()
}

func (m List) selectedDate() time.Time {
	return m.t.AddDate(0, 0, m.selected)
}

func (m List) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case msgs.Msg:
		switch msg {
		case msgs.MsgUpdateList:
			cmds = append(cmds, m.list.SetItems(listitem.HabitsToItems(m.group.Habits)))
		}
	case textinput.Submit:
		switch msg.Key {
		case "habit":
			hab := state.Get(msg.ID)
			hab.Name = msg.Value
			if err := habit.Client.Update(hab); err != nil {
				logger.Error.Printf("failed to update habit: %v", err)
			}
			cmds = append(cmds, msgs.UpdateList)
		}

	case tea.KeyMsg:
		switch {
		case m.list.SettingFilter():
			break

		case key.Matches(msg, m.km.tab):
			m.onHist = !m.onHist

			if m.onHist {
				m.lastSelected = m.list.Index()
				m.list.Select(-1)
				m.selected -= config.TimeFrameInDays
			} else {
				m.list.Select(m.lastSelected)
				m.selected += config.TimeFrameInDays
			}

		case m.onHist && key.Matches(msg, m.km.left):
			m.selected = util.Max(m.selected-7, 0)

		case m.onHist && key.Matches(msg, m.km.right):
			m.selected = util.Min(m.selected+7, config.TimeFrameInDays-1)

		case m.onHist && key.Matches(msg, m.km.up) && m.selected > 0:
			m.selected -= 1

		case m.onHist && key.Matches(msg, m.km.down) && (m.selected+1) < config.TimeFrameInDays:
			m.selected += 1

		case key.Matches(msg, m.km.quit):
			return m.parent, msgs.UpdateList

		case len(m.list.VisibleItems()) == 0:
			break

		case key.Matches(msg, m.km.checkIn):
			selected := m.list.SelectedItem().(listitem.HabitItem).Habit
			dto := habit.CheckInDTO{ID: selected.ID, Desc: "", Date: time.Now().Local()}
			hab, err := habit.Client.CheckIn(dto)
			if err != nil {
				logger.Error.Printf("failed to add activity: %v", err)
			} else {
				state.SetHabit(hab)
			}

		case key.Matches(msg, m.km.view):
			selected := m.list.SelectedItem().(listitem.HabitItem).Habit
			return show.NewShow(selected, m), nil

		case key.Matches(msg, m.km.rename):
			selected := m.list.SelectedItem().(listitem.HabitItem).Habit
			return textinput.New(m, "New Name:", selected.Name, "habit", selected.ID), nil

		case key.Matches(msg, m.km.delete):
			selected := m.list.SelectedItem().(listitem.HabitItem).Habit
			for i, r := range m.list.Items() {
				if it, ok := r.(listitem.HabitItem); ok && it.Habit.ID == selected.ID {
					m.list.RemoveItem(i)
				}
			}
			state.RemoveFromGroup(*selected, *m.group)
			m.list.SetFilteringEnabled(false)
			m.list.Select(util.Min(m.list.Index(), len(state.Habits())-1))
			return m, m.list.NewStatusMessage("Removed " + selected.Name)
		}
	}

	if !m.onHist {
		var cmd tea.Cmd
		m.list, cmd = m.list.Update(msg)
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}
