package group

import (
	"astro/config"
	"astro/date"
	"astro/habit"
	"astro/histogram"
	"astro/logger"
	"astro/models/name"
	"astro/models/show"
	"astro/msgs"
	"astro/state"
	"astro/util"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type List struct {
	group    *habit.Group
	parent   tea.Model
	list     list.Model
	km       binds
	t        time.Time
	selected int
	onHist   bool
}

func NewShow(g *habit.Group, parent tea.Model) List {
	l := list.New(toItems(g.Habits), list.NewDefaultDelegate(), 0, 5)
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

	s.WriteString("\n")
	s.WriteString(m.list.View())
	return s.String()
}

func (m List) selectedDate() time.Time {
	return m.t.AddDate(0, 0, m.selected)
}

func (m List) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case m.list.SettingFilter():
			break

		case m.onHist && key.Matches(msg, m.km.left):
			m.selected = util.Max(m.selected-7, 0)

		case m.onHist && key.Matches(msg, m.km.right):
			m.selected = util.Min(m.selected+7, config.TimeFrameInDays-1)

		case key.Matches(msg, m.km.up):
			if m.onHist {
				if m.selected > 0 {
					m.selected -= 1
				}
			}

			if m.list.Index() == 0 {
				m.onHist = true
				m.selected -= config.TimeFrameInDays
				m.list.Select(-1)
			}

		case key.Matches(msg, m.km.down):
			if m.onHist {
				if (m.selected+1)%7 == 0 {
					m.onHist = false
					m.selected += config.TimeFrameInDays
				} else {
					m.selected += 1
				}
			}

		case key.Matches(msg, m.km.quit):
			return m.parent, msgs.UpdateList

		case len(m.list.VisibleItems()) == 0:
			break

		case key.Matches(msg, m.km.checkIn):
			selected := m.list.SelectedItem().(habitItem).habit
			hab, err := habit.Client.CheckIn(selected.ID, "")
			if err != nil {
				logger.Error.Printf("failed to add activity: %v", err)
			} else {
				state.SetHabit(hab)
			}

		case key.Matches(msg, m.km.view):
			selected := m.list.SelectedItem().(habitItem).habit
			return show.NewShow(selected, m), nil

		case key.Matches(msg, m.km.rename):
			selected := m.list.SelectedItem().(habitItem).habit
			return name.NewEditName(selected, m), nil

		case key.Matches(msg, m.km.delete):
			selected := m.list.SelectedItem().(habitItem).habit
			for i, r := range m.list.Items() {
				if it, ok := r.(habitItem); ok && it.habit.ID == selected.ID {
					m.list.RemoveItem(i)
				}
			}
			state.RemoveFromGroup(*selected, *m.group)
			m.list.SetFilteringEnabled(false)
			m.list.Select(util.Min(m.list.Index(), len(state.Habits())-1))
			return m, m.list.NewStatusMessage("Removed " + selected.Name)
		}
	}

	var cmd tea.Cmd
	if !m.onHist {
		m.list, cmd = m.list.Update(msg)
	}
	return m, cmd
}
