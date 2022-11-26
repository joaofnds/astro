package add_to_group

import (
	"astro/habit"
	"astro/msgs"
	"astro/state"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type item struct{ group *habit.Group }

func (i item) Title() string       { return i.group.Name }
func (i item) Description() string { return i.Title() }
func (i item) FilterValue() string { return i.Title() }

func toItems(groups []*habit.Group) []list.Item {
	items := make([]list.Item, len(groups))
	for i, g := range groups {
		items[i] = item{g}
	}
	return items
}

type ChooseGroup struct {
	parent tea.Model
	list   list.Model
	habit  *habit.Habit
}

func NewChooseGroup(parent tea.Model, h *habit.Habit) ChooseGroup {
	l := list.New(toItems(state.Groups()), list.NewDefaultDelegate(), 1, 5)
	l.Title = "Choose a group"
	return ChooseGroup{parent: parent, list: l, habit: h}
}

func (m ChooseGroup) Init() tea.Cmd {
	return nil
}

func (m ChooseGroup) View() string {
	return m.list.View()
}

func (m ChooseGroup) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height)

	case tea.KeyMsg:
		switch {
		case m.list.SettingFilter():
			break

		case len(m.list.VisibleItems()) == 0:
			break

		case msg.Type == tea.KeyEsc:
			return m.parent, nil

		case msg.Type == tea.KeyEnter:
			group := m.list.SelectedItem().(item).group
			state.AddToGroup(*m.habit, *group)
			return m.parent, msgs.UpdateList
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}
