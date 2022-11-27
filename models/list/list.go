package list

import (
	"astro/config"
	"astro/habit"
	"astro/logger"
	"astro/models/add_to_group"
	"astro/models/group"
	"astro/models/listitem"
	"astro/models/name"
	"astro/models/show"
	"astro/msgs"
	"astro/state"
	"astro/util"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type List struct {
	list list.Model
	km   habitBinds
}

func NewList() List {
	l := list.New(items(), list.NewDefaultDelegate(), 0, 5)
	km := NewHabitBinds()
	l.Title = "Habits"
	l.AdditionalShortHelpKeys = km.ToSlice
	return List{list: l, km: km}
}

func items() []list.Item {
	habits, groups := listitem.HabitsToItems(state.Habits()), listitem.GroupsToItems(state.Groups())
	items := make([]list.Item, 0, len(habits)+len(groups))
	items = append(items, habits...)
	items = append(items, groups...)
	return items
}

func (m List) Init() tea.Cmd {
	return nil
}

func (m List) View() string {
	return m.list.View()
}

func (m List) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case msgs.Msg:
		switch msg {
		case msgs.MsgUpdateList:
			cmds = append(cmds, m.list.SetItems(items()))
		}

	case tea.WindowSizeMsg:
		config.Width, config.Height = msg.Width, msg.Height
		m.list.SetSize(msg.Width, msg.Height)

	case tea.KeyMsg:
		switch {
		case m.list.SettingFilter():
			break

		case key.Matches(msg, m.km.add):
			return newAddInput(m), nil

		case key.Matches(msg, m.km.addGroup):
			return group.NewAddGroup(m), nil

		case len(m.list.VisibleItems()) == 0:
			break

		default:
			switch m.list.SelectedItem().(type) {
			case listitem.HabitItem:
				selected := m.list.SelectedItem().(listitem.HabitItem).Habit
				switch {
				case key.Matches(msg, m.km.view):
					return show.NewShow(selected, m), nil

				case key.Matches(msg, m.km.rename):
					return name.NewEditName(selected, m), nil

				case key.Matches(msg, m.km.addToGroup):
					selected := m.list.SelectedItem().(listitem.HabitItem).Habit
					return add_to_group.NewChooseGroup(m, selected), nil

				case key.Matches(msg, m.km.delete):
					for i, r := range m.list.Items() {
						if it, ok := r.(listitem.HabitItem); ok && it.Habit.ID == selected.ID {
							m.list.RemoveItem(i)
						}
					}
					if err := state.Delete(selected.ID); err != nil {
						panic(err)
					}
					m.list.SetFilteringEnabled(false)
					m.list.Select(util.Min(m.list.Index(), len(state.Habits())-1))
					return m, m.list.NewStatusMessage("Removed " + selected.Name)

				case key.Matches(msg, m.km.checkIn):
					selected := m.list.SelectedItem().(listitem.HabitItem).Habit
					hab, err := habit.Client.CheckIn(selected.ID, "")
					if err != nil {
						logger.Error.Printf("failed to add activity: %v", err)
					} else {
						state.SetHabit(hab)
					}
				}
			case listitem.GroupItem:
				selected := m.list.SelectedItem().(listitem.GroupItem).Group
				switch {
				case key.Matches(msg, m.km.view):
					return group.NewShow(selected, m), nil

				case key.Matches(msg, m.km.delete):
					state.DeleteGroup(*selected)
					cmds = append(cmds, msgs.UpdateList)
				}
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}
