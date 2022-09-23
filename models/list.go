package models

import (
	"astroapp/config"
	"astroapp/habit"
	"astroapp/state"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type item struct{ habit *habit.Habit }

func (i item) Title() string { return i.habit.Name }
func (i item) Description() string {
	if len(i.habit.Activities) == 0 {
		return "no activities"
	}

	return "latest activity on " + i.habit.LatestActivity().Format(config.TimeFormat)
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
}

func NewList() List {
	list := list.New(toItems(state.Habits()), list.NewDefaultDelegate(), 0, 5)
	list.Title = "Habits"
	return List{list}
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
		switch msg.String() {
		case "enter":
			habit := state.At(m.list.Index())
			return NewShow(habit, m), nil
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}
