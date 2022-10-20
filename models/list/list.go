package list

import (
	"astro/config"
	"astro/habit"
	"astro/histogram"
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

	return "last activity at " + i.habit.LatestActivity().Format(config.TimeFormat)
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
	list := list.New(toItems(state.Habits()), list.NewDefaultDelegate(), 0, 5)
	list.Title = "Habits"
	list.AdditionalShortHelpKeys = km.ToSlice
	return List{list, km}
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
		if m.list.SettingFilter() {
			break
		}
		switch {
		case key.Matches(msg, m.km.delete):
			h := state.At(m.list.Index())
			if err := state.Delete(h.ID); err != nil {
				panic(err)
			}
			m.list.RemoveItem(m.list.Index())
			m.list.Select(util.Min(m.list.Index(), len(state.Habits())-1))
			return m, m.list.NewStatusMessage("Removed " + h.Name)

		case key.Matches(msg, m.km.add):
			return newAddInput(m), nil

		case key.Matches(msg, m.km.view):
			habit := state.At(m.list.Index())
			return show.NewShow(habit, m), nil
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}
