package models

import (
	"astroapp/config"
	"astroapp/habit"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle()

type item habit.Habit

func (i item) Title() string { return i.Name }
func (i item) Description() string {
	if len(i.Activities) == 0 {
		return "no activities"
	}

	return "latest activity on " + i.ToHabit().LatestActivity().Format(config.TimeFormat)
}
func (i item) FilterValue() string  { return i.Name }
func (i item) ToHabit() habit.Habit { return habit.Habit(i) }

type List struct {
	list list.Model
}

func NewList(habits []habit.Habit) List {
	items := make([]list.Item, len(habits))
	for i, h := range habits {
		items[i] = item(h)
	}
	list := list.New(items, list.NewDefaultDelegate(), 0, 5)
	list.Title = "Habits"
	return List{list}
}

func (m List) Init() tea.Cmd {
	return nil
}

func (m List) View() string {
	return docStyle.Render(m.list.View())
}

func (m List) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			return NewShow(m.list.SelectedItem().(item).ToHabit(), m), nil
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}
