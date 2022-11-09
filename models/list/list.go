package list

import (
	"astro/config"
	"astro/habit"
	"astro/histogram"
	"astro/logger"
	"astro/models/name"
	"astro/models/show"
	"astro/state"
	"astro/util"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type item struct{ habit *habit.Habit }

func (i item) Title() string { return i.habit.Name }
func (i item) Description() string {
	return histogram.ShortLineHistogram(*i.habit, config.ShortHistSize) + " " + i.lastActivity()
}

func (i item) lastActivity() string {
	if len(i.habit.Activities) == 0 {
		return "no activities"
	}

	return "last activity at " + i.habit.LatestActivity().Local().Format(config.DateFormat)
}

func (i item) FilterValue() string { return i.habit.Name }
func toItems(habits []*habit.Habit) []list.Item {
	items := make([]list.Item, len(habits))
	for i, h := range habits {
		items[i] = item{h}
	}
	return items
}

type List struct {
	list list.Model
	km   keymap
}

func NewList() List {
	km := NewKeymap()
	l := list.New(toItems(state.Habits()), list.NewDefaultDelegate(), 0, 5)
	l.Title = "Habits"
	l.AdditionalShortHelpKeys = km.ToSlice
	return List{l, km}
}

func (m List) Init() tea.Cmd {
	return nil
}

func (m List) View() string {
	return m.list.View()
}

func (m List) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height)

	case tea.KeyMsg:
		switch {
		case m.list.SettingFilter():
			break

		case key.Matches(msg, m.km.add):
			return newAddInput(m), nil

		case len(m.list.VisibleItems()) == 0:
			break

		case key.Matches(msg, m.km.checkIn):
			selected := m.list.SelectedItem().(item).habit
			hab, err := habit.Client.CheckIn(selected.ID, "")
			if err != nil {
				logger.Error.Printf("failed to add activity: %v", err)
			} else {
				state.SetHabit(hab)
			}

		case key.Matches(msg, m.km.view):
			selected := m.list.SelectedItem().(item).habit
			return show.NewShow(selected, m), nil

		case key.Matches(msg, m.km.rename):
			selected := m.list.SelectedItem().(item).habit
			return name.NewEditName(selected, m), nil

		case key.Matches(msg, m.km.delete):
			selected := m.list.SelectedItem().(item).habit
			for i, r := range m.list.Items() {
				if it, ok := r.(item); ok && it.habit.ID == selected.ID {
					m.list.RemoveItem(i)
				}
			}
			if err := state.Delete(selected.ID); err != nil {
				panic(err)
			}
			m.list.SetFilteringEnabled(false)
			m.list.Select(util.Min(m.list.Index(), len(state.Habits())-1))
			return m, m.list.NewStatusMessage("Removed " + selected.Name)
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}
