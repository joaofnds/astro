package add_to_group

import (
	"astro/api"
	"astro/domain"
	"astro/msgs"
	"context"

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
	client *api.Client
	list   list.Model
	habit  *domain.Habit
}

func NewChooseGroup(client *api.Client, h *domain.Habit, groups []*domain.Group) ChooseGroup {
	l := list.New(toItems(groups), list.NewDefaultDelegate(), 1, 5)
	l.Title = "Choose a group"
	return ChooseGroup{client: client, list: l, habit: h}
}

func (m ChooseGroup) Init() tea.Cmd {
	return nil
}

func (m ChooseGroup) View() tea.View {
	return tea.NewView(m.list.View())
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
			return m, msgs.PopScreen()

		case msg.String() == "enter":
			sel, ok := m.list.SelectedItem().(item)
			if !ok {
				break
			}
			return m, func() tea.Msg {
				return msgs.PopScreenMsg{
					Cmd: msgs.AddToGroup(context.Background(), m.client, m.habit.ID, sel.group.ID),
				}
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}
