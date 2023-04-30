package list

import (
	"astro/config"
	"astro/habit"
	"astro/logger"
	"astro/models/add_to_group"
	"astro/models/group"
	"astro/models/listitem"
	"astro/models/show"
	"astro/models/textinput"
	"astro/msgs"
	"astro/state"
	"astro/util"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type List struct {
	list    list.Model
	help    help.Model
	habitKM habitBinds
	groupKM groupBinds
}

func NewList() List {
	l := list.New(items(), list.NewDefaultDelegate(), 0, 5)
	l.Title = "Habits"
	l.SetShowHelp(false)
	return List{
		list:    l,
		help:    help.New(),
		habitKM: NewHabitBinds(),
		groupKM: NewGroupBinds(),
	}
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
	var s strings.Builder
	s.WriteString(m.list.View() + "\n")
	switch m.list.SelectedItem().(type) {
	case listitem.HabitItem:
		s.WriteString(m.help.View(m.habitKM))
	case listitem.GroupItem:
		s.WriteString(m.help.View(m.groupKM))
	}
	return s.String()
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
		m.list.SetSize(msg.Width, msg.Height-1)

	case textinput.Submit:
		switch msg.Key {
		case "habit":
			hab := state.Get(msg.ID)
			hab.Name = msg.Value
			if err := habit.Client.Update(hab); err != nil {
				logger.Error.Printf("failed to update habit: %v", err)
			}
			cmds = append(cmds, msgs.UpdateList)
		}

	case tea.KeyMsg:
		switch {
		case m.list.SettingFilter():
			break

		case key.Matches(msg, m.habitKM.add):
			return newAddInput(m), nil

		case key.Matches(msg, m.habitKM.addGroup):
			return group.NewAddGroup(m), nil

		case len(m.list.VisibleItems()) == 0:
			break

		default:
			switch m.list.SelectedItem().(type) {
			case listitem.HabitItem:
				selected := m.list.SelectedItem().(listitem.HabitItem).Habit
				switch {
				case key.Matches(msg, m.habitKM.view):
					return show.NewShow(selected, m), nil

				case key.Matches(msg, m.habitKM.rename):
					return textinput.New(m, "New Name:", selected.Name, "habit", selected.ID), nil

				case key.Matches(msg, m.habitKM.addToGroup):
					selected := m.list.SelectedItem().(listitem.HabitItem).Habit
					return add_to_group.NewChooseGroup(m, selected), nil

				case key.Matches(msg, m.habitKM.delete):
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

				case key.Matches(msg, m.habitKM.checkIn):
					selected := m.list.SelectedItem().(listitem.HabitItem).Habit
					dto := habit.CheckInDTO{ID: selected.ID, Desc: "", Date: time.Now().Local()}
					hab, err := habit.Client.CheckIn(dto)
					if err != nil {
						logger.Error.Printf("failed to add activity: %v", err)
					} else {
						state.SetHabit(hab)
					}
				}

			case listitem.GroupItem:
				selected := m.list.SelectedItem().(listitem.GroupItem).Group
				switch {
				case key.Matches(msg, m.groupKM.view):
					return group.NewShow(selected, m), nil

				case key.Matches(msg, m.groupKM.delete):
					if err := state.DeleteGroup(*selected); err != nil {
						logger.Error.Printf("failed to delete group: %v", err)
					}
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
