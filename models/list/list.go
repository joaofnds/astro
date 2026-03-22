package list

import (
	"astro/api"
	"astro/domain"
	"astro/models/add_to_group"
	"astro/models/group"
	"astro/models/listitem"
	"astro/models/show"
	"astro/models/textinput"
	"astro/msgs"
	"strings"
	"time"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
)

type List struct {
	client  *api.Client
	list    list.Model
	help    help.Model
	habitKM habitBinds
	groupKM groupBinds
	groups  []*domain.Group
	width   int
	height  int
}

func NewList(client *api.Client, habits []*domain.Habit, groups []*domain.Group, width, height int) List {
	habitItems := listitem.HabitsToItems(habits)
	groupItems := listitem.GroupsToItems(groups)
	items := make([]list.Item, 0, len(habitItems)+len(groupItems))
	items = append(items, habitItems...)
	items = append(items, groupItems...)

	l := list.New(items, list.NewDefaultDelegate(), width, height)
	l.Title = "Habits"
	l.SetShowHelp(false)
	return List{
		client:  client,
		list:    l,
		help:    help.New(),
		habitKM: NewHabitBinds(),
		groupKM: NewGroupBinds(),
		groups:  groups,
		width:   width,
		height:  height,
	}
}

func (m List) Init() tea.Cmd {
	return nil
}

func (m List) View() tea.View {
	var s strings.Builder
	s.WriteString(m.list.View() + "\n")
	switch m.list.SelectedItem().(type) {
	case listitem.HabitItem:
		s.WriteString(m.help.View(m.habitKM))
	case listitem.GroupItem:
		s.WriteString(m.help.View(m.groupKM))
	}
	return tea.NewView(s.String())
}

func (m List) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case msgs.HabitCreatedMsg:
		cmd := m.list.InsertItem(len(m.list.Items()), listitem.HabitItem{Habit: msg.Habit})
		cmds = append(cmds, cmd, m.list.NewStatusMessage("Added "+msg.Habit.Name))

	case msgs.HabitDeletedMsg:
		for i, item := range m.list.Items() {
			if hi, ok := item.(listitem.HabitItem); ok && hi.Habit.ID == msg.ID {
				m.list.RemoveItem(i)
				break
			}
		}

	case msgs.HabitUpdatedMsg:
		for i, item := range m.list.Items() {
			if hi, ok := item.(listitem.HabitItem); ok && hi.Habit.ID == msg.Habit.ID {
				cmds = append(cmds, m.list.SetItem(i, listitem.HabitItem{Habit: msg.Habit}))
				break
			}
		}

	case msgs.CheckInResultMsg:
		for i, item := range m.list.Items() {
			if hi, ok := item.(listitem.HabitItem); ok && hi.Habit.ID == msg.Habit.ID {
				cmds = append(cmds, m.list.SetItem(i, listitem.HabitItem{Habit: msg.Habit}))
				break
			}
		}

	case msgs.GroupCreatedMsg:
		m.groups = append(m.groups, msg.Group)
		cmd := m.list.InsertItem(len(m.list.Items()), listitem.GroupItem{Group: msg.Group})
		cmds = append(cmds, cmd, m.list.NewStatusMessage("Added "+msg.Group.Name))

	case msgs.GroupDeletedMsg:
		for i, item := range m.list.Items() {
			if gi, ok := item.(listitem.GroupItem); ok && gi.Group.ID == msg.ID {
				m.list.RemoveItem(i)
				break
			}
		}
		for i, g := range m.groups {
			if g.ID == msg.ID {
				m.groups = append(m.groups[:i], m.groups[i+1:]...)
				break
			}
		}

	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.list.SetSize(msg.Width, msg.Height-1)

	case textinput.Submit:
		switch msg.Key {
		case "habit":
			return m, msgs.UpdateHabit(m.client, msg.ID, msg.Value)
		}

	case tea.KeyPressMsg:
		switch {
		case m.list.SettingFilter():
			break

		case key.Matches(msg, m.habitKM.add):
			return m, msgs.PushScreen(newAddInput(m.client))

		case key.Matches(msg, m.habitKM.addGroup):
			return m, msgs.PushScreen(group.NewAddGroup(m.client))

		case len(m.list.VisibleItems()) == 0:
			break

		default:
			switch m.list.SelectedItem().(type) {
			case listitem.HabitItem:
				selected := m.list.SelectedItem().(listitem.HabitItem).Habit
				switch {
				case key.Matches(msg, m.habitKM.view):
					return m, msgs.PushScreen(show.NewShow(m.client, selected, m.width))

				case key.Matches(msg, m.habitKM.rename):
					return m, msgs.PushScreen(textinput.New("New Name:", selected.Name, "habit", selected.ID, m.width))

				case key.Matches(msg, m.habitKM.addToGroup):
					return m, msgs.PushScreen(add_to_group.NewChooseGroup(m.client, selected, m.groups))

				case key.Matches(msg, m.habitKM.delete):
					for i, r := range m.list.Items() {
						if it, ok := r.(listitem.HabitItem); ok && it.Habit.ID == selected.ID {
							m.list.RemoveItem(i)
							break
						}
					}
					return m, tea.Batch(
						msgs.DeleteHabit(m.client, selected.ID),
						m.list.NewStatusMessage("Removed "+selected.Name),
					)

				case key.Matches(msg, m.habitKM.checkIn):
					return m, msgs.CheckIn(m.client, selected.ID, "", time.Now().Local())
				}

			case listitem.GroupItem:
				selected := m.list.SelectedItem().(listitem.GroupItem).Group
				switch {
				case key.Matches(msg, m.groupKM.view):
					return m, msgs.PushScreen(group.NewShow(m.client, selected, m.width, m.height))

				case key.Matches(msg, m.groupKM.delete):
					return m, msgs.DeleteGroup(m.client, selected.ID)
				}
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}
