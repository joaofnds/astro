package add_to_group

import (
	"astro/domain"
	"astro/msgs"
	"astro/state"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
)

type item struct{ group *domain.Group }

func (i item) Title() string       { return i.group.Name }
func (i item) Description() string { return i.Title() }
func (i item) FilterValue() string { return i.Title() }

func toItems(groups []*domain.Group) []list.Item {
	items := make([]list.Item, len(groups))
	for i, g := range groups {
		items[i] = item{g}
	}
	return items
}

type ChooseGroup struct {
	parent tea.Model
	list   list.Model
	habit  *domain.Habit
}

func NewChooseGroup(parent tea.Model, h *domain.Habit) ChooseGroup {
	l := list.New(toItems(state.Groups()), list.NewDefaultDelegate(), 1, 5)
	l.Title = "Choose a group"
	return ChooseGroup{parent: parent, list: l, habit: h}
}

func (m ChooseGroup) Init() tea.Cmd {
	return nil
}

func (m ChooseGroup) View() tea.View {
	v := tea.NewView(m.list.View())
	v.AltScreen = true
	return v
}

func (m ChooseGroup) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height)

	case tea.KeyPressMsg:
		switch {
		case m.list.SettingFilter():
			break

		case len(m.list.VisibleItems()) == 0:
			break

		case msg.String() == "esc":
			return m.parent, nil

		case msg.String() == "enter":
			group := m.list.SelectedItem().(item).group
			state.AddToGroup(*m.habit, *group)
			return m.parent, msgs.UpdateList
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}
