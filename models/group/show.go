package group

import (
	"astro/api"
	"astro/config"
	"astro/date"
	"astro/domain"
	"astro/models/listitem"
	"astro/models/show"
	"astro/models/textinput"
	"astro/msgs"
	"astro/util"
	"strings"
	"time"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
)

type List struct {
	client       *api.Client
	group        *domain.Group
	list         list.Model
	km           binds
	t            time.Time
	selected     int
	lastSelected int
	onHist       bool
	width        int
	height       int
}

func NewShow(client *api.Client, g *domain.Group, width, height int) List {
	l := list.New(listitem.HabitsToItems(g.Habits), list.NewDefaultDelegate(), 0, 5)
	l.SetSize(width, height-9)

	km := newBinds()
	l.AdditionalShortHelpKeys = km.ToSlice
	l.Title = g.Name

	t, _ := date.TimeFrame()
	selected := date.DiffInDays(t, date.Today()) + config.TimeFrameInDays

	return List{
		client:   client,
		t:        t,
		selected: selected,
		group:    g,
		list:     l,
		km:       km,
		width:    width,
		height:   height,
	}
}

func (m List) Init() tea.Cmd {
	return nil
}

func (m List) View() tea.View {
	activities := m.group.Activities()
	var s strings.Builder
	s.WriteString(domain.Histogram(m.t, activities, m.selected))

	if m.selectedDate().After(date.Today()) {
		s.WriteString("\n")
	} else {
		s.WriteString(domain.ActivitiesOnDateTally(m.group.Habits, m.selectedDate()))
	}

	m.list.Title = domain.Digest(m.group.Name, m.group.Activities())

	s.WriteString("\n")
	s.WriteString(m.list.View())
	return tea.NewView(s.String())
}

func (m List) selectedDate() time.Time {
	return m.t.AddDate(0, 0, m.selected)
}

func (m List) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case msgs.CheckInResultMsg:
		for i, h := range m.group.Habits {
			if h.ID == msg.Habit.ID {
				m.group.Habits[i] = msg.Habit
				break
			}
		}
		cmds = append(cmds, m.list.SetItems(listitem.HabitsToItems(m.group.Habits)))

	case msgs.HabitUpdatedMsg:
		for i, h := range m.group.Habits {
			if h.ID == msg.Habit.ID {
				m.group.Habits[i] = msg.Habit
				break
			}
		}
		cmds = append(cmds, m.list.SetItems(listitem.HabitsToItems(m.group.Habits)))

	case textinput.Submit:
		switch msg.Key {
		case "habit":
			return m, msgs.UpdateHabit(m.client, msg.ID, msg.Value)
		}

	case tea.KeyPressMsg:
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

		// ClearScreen forces a full sequential redraw on navigation.
		// See comment in show.go for details on the emoji width issue.
		case m.onHist && key.Matches(msg, m.km.left):
			m.selected = util.Max(m.selected-7, 0)
			return m, tea.ClearScreen

		case m.onHist && key.Matches(msg, m.km.right):
			m.selected = util.Min(m.selected+7, config.TimeFrameInDays-1)
			return m, tea.ClearScreen

		case m.onHist && key.Matches(msg, m.km.up) && m.selected > 0:
			m.selected -= 1
			return m, tea.ClearScreen

		case m.onHist && key.Matches(msg, m.km.down) && (m.selected+1) < config.TimeFrameInDays:
			m.selected += 1
			return m, tea.ClearScreen

		case key.Matches(msg, m.km.quit):
			return m, msgs.PopScreen()

		case len(m.list.VisibleItems()) == 0:
			break

		case key.Matches(msg, m.km.checkIn):
			selected := m.list.SelectedItem().(listitem.HabitItem).Habit
			return m, msgs.CheckIn(m.client, selected.ID, "", time.Now().Local())

		case key.Matches(msg, m.km.view):
			selected := m.list.SelectedItem().(listitem.HabitItem).Habit
			return m, msgs.PushScreen(show.NewShow(m.client, selected, m.width))

		case key.Matches(msg, m.km.rename):
			selected := m.list.SelectedItem().(listitem.HabitItem).Habit
			return m, msgs.PushScreen(textinput.New("New Name:", selected.Name, "habit", selected.ID, m.width))

		case key.Matches(msg, m.km.delete):
			selected := m.list.SelectedItem().(listitem.HabitItem).Habit
			for i, r := range m.list.Items() {
				if it, ok := r.(listitem.HabitItem); ok && it.Habit.ID == selected.ID {
					m.list.RemoveItem(i)
					break
				}
			}
			return m, tea.Batch(
				msgs.RemoveFromGroup(m.client, selected.ID, m.group.ID),
				m.list.NewStatusMessage("Removed "+selected.Name),
			)
		}
	}

	if !m.onHist {
		var cmd tea.Cmd
		m.list, cmd = m.list.Update(msg)
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}
