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
	"context"
	"fmt"
	"strings"
	"time"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
)

// pendingOp stores undo information for one in-flight API call.
type pendingOp struct {
	op   string    // operation name, matches APIErrorMsg.Op
	id   string    // entity ID for correlation
	item list.Item // the removed/modified item for rollback (nil for creates)
	name string    // old name for rename rollback
}

// errorQueue holds error messages for sequential display.
type errorQueue struct {
	msgs []string
}

func (q *errorQueue) push(msg string) {
	q.msgs = append(q.msgs, msg)
}

func (q *errorQueue) pop() (string, bool) {
	if len(q.msgs) == 0 {
		return "", false
	}
	msg := q.msgs[0]
	q.msgs = q.msgs[1:]
	return msg, true
}

// createHabitSubmit is sent when the user submits a new habit name.
type createHabitSubmit struct {
	Name string
}

type List struct {
	client   *api.Client
	list     list.Model
	help     help.Model
	habitKM  habitBinds
	groupKM  groupBinds
	groups   []*domain.Group
	width    int
	height   int
	pending  []pendingOp
	errQueue errorQueue
	cancelOp context.CancelFunc
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
	l.StatusMessageLifetime = 3 * time.Second

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
	default:
		// nil or pending item -- no contextual help
	}
	return tea.NewView(s.String())
}

func (m List) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case createHabitSubmit:
		cmd := m.list.InsertItem(len(m.list.Items()), listitem.PendingHabitItem{Name: msg.Name})
		m.pending = append(m.pending, pendingOp{op: "create habit"})
		ctx, cancel := context.WithCancel(context.Background())
		m.cancelOp = cancel
		return m, tea.Batch(
			cmd,
			msgs.CreateHabit(ctx, m.client, msg.Name),
			m.list.NewStatusMessage("Creating "+msg.Name+"..."),
		)

	case group.CreateGroupSubmit:
		cmd := m.list.InsertItem(len(m.list.Items()), listitem.PendingGroupItem{Name: msg.Name})
		m.pending = append(m.pending, pendingOp{op: "create group"})
		ctx, cancel := context.WithCancel(context.Background())
		m.cancelOp = cancel
		return m, tea.Batch(
			cmd,
			msgs.CreateGroup(ctx, m.client, msg.Name),
			m.list.NewStatusMessage("Creating "+msg.Name+"..."),
		)

	case msgs.HabitCreatedMsg:
		for i, item := range m.list.Items() {
			if _, ok := item.(listitem.PendingHabitItem); ok {
				m.list.RemoveItem(i)
				break
			}
		}
		for i, p := range m.pending {
			if p.op == "create habit" {
				m.pending = append(m.pending[:i], m.pending[i+1:]...)
				break
			}
		}
		cmd := m.list.InsertItem(len(m.list.Items()), listitem.HabitItem{Habit: msg.Habit})
		cmds = append(cmds, cmd, m.list.NewStatusMessage("Added "+msg.Habit.Name))

	case msgs.HabitDeletedMsg:
		for i, p := range m.pending {
			if p.op == "delete habit" && p.id == msg.ID {
				m.pending = append(m.pending[:i], m.pending[i+1:]...)
				break
			}
		}
		// Remove from list if still present (non-optimistic path).
		for i, item := range m.list.Items() {
			if hi, ok := item.(listitem.HabitItem); ok && hi.Habit.ID == msg.ID {
				m.list.RemoveItem(i)
				break
			}
		}

	case msgs.HabitUpdatedMsg:
		for i, p := range m.pending {
			if p.op == "update habit" && p.id == msg.Habit.ID {
				m.pending = append(m.pending[:i], m.pending[i+1:]...)
				break
			}
		}
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
		cmds = append(cmds, m.list.NewStatusMessage("Checked in"))

	case msgs.GroupCreatedMsg:
		for i, item := range m.list.Items() {
			if _, ok := item.(listitem.PendingGroupItem); ok {
				m.list.RemoveItem(i)
				break
			}
		}
		for i, p := range m.pending {
			if p.op == "create group" {
				m.pending = append(m.pending[:i], m.pending[i+1:]...)
				break
			}
		}
		m.groups = append(m.groups, msg.Group)
		cmd := m.list.InsertItem(len(m.list.Items()), listitem.GroupItem{Group: msg.Group})
		cmds = append(cmds, cmd, m.list.NewStatusMessage("Added "+msg.Group.Name))

	case msgs.GroupDeletedMsg:
		for i, p := range m.pending {
			if p.op == "delete group" && p.id == msg.ID {
				m.pending = append(m.pending[:i], m.pending[i+1:]...)
				break
			}
		}
		// Remove from list if still present (non-optimistic path).
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

	case msgs.APIErrorMsg:
		errStr := fmt.Sprintf("\u2717 %s: %s", msg.Op, msg.Err)
		for i, p := range m.pending {
			if p.op == msg.Op && (msg.ID == "" || p.id == msg.ID) {
				m.pending = append(m.pending[:i], m.pending[i+1:]...)
				switch p.op {
				case "delete habit", "delete group":
					cmds = append(cmds, m.list.InsertItem(len(m.list.Items()), p.item))
					if p.op == "delete group" {
						if gi, ok := p.item.(listitem.GroupItem); ok {
							m.groups = append(m.groups, gi.Group)
						}
					}
					errStr += " \u2014 restored"
				case "update habit":
					if p.item != nil {
						for j, item := range m.list.Items() {
							if hi, ok := item.(listitem.HabitItem); ok && hi.Habit.ID == p.id {
								cmds = append(cmds, m.list.SetItem(j, p.item))
								break
							}
						}
					}
					errStr += " \u2014 restored"
				case "create habit":
					for j, item := range m.list.Items() {
						if _, ok := item.(listitem.PendingHabitItem); ok {
							m.list.RemoveItem(j)
							break
						}
					}
					errStr += " \u2014 removed"
				case "create group":
					for j, item := range m.list.Items() {
						if _, ok := item.(listitem.PendingGroupItem); ok {
							m.list.RemoveItem(j)
							break
						}
					}
					errStr += " \u2014 removed"
				}
				break
			}
		}
		m.errQueue.push(errStr)
		if len(m.errQueue.msgs) == 1 {
			return m, tea.Batch(append(cmds,
				m.list.NewStatusMessage(errStr),
				msgs.ClearStatusAfter(3*time.Second),
			)...)
		}
		return m, tea.Batch(cmds...)

	case msgs.ClearStatusMsg:
		if errMsg, ok := m.errQueue.pop(); ok {
			return m, tea.Batch(
				m.list.NewStatusMessage(errMsg),
				msgs.ClearStatusAfter(3*time.Second),
			)
		}

	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.list.SetSize(msg.Width, msg.Height-1)

	case textinput.Submit:
		switch msg.Key {
		case "habit":
			for _, item := range m.list.Items() {
				if hi, ok := item.(listitem.HabitItem); ok && hi.Habit.ID == msg.ID {
					m.pending = append(m.pending, pendingOp{
						op:   "update habit",
						id:   msg.ID,
						item: item,
						name: hi.Habit.Name,
					})
					updated := *hi.Habit
					updated.Name = msg.Value
					for i, it := range m.list.Items() {
						if h, ok := it.(listitem.HabitItem); ok && h.Habit.ID == msg.ID {
							cmds = append(cmds, m.list.SetItem(i, listitem.HabitItem{Habit: &updated}))
							break
						}
					}
					break
				}
			}
			ctx, cancel := context.WithCancel(context.Background())
			m.cancelOp = cancel
			return m, tea.Batch(
				append(cmds,
					msgs.UpdateHabit(ctx, m.client, msg.ID, msg.Value),
					m.list.NewStatusMessage("Renaming..."),
				)...,
			)
		}

	case tea.KeyPressMsg:
		switch {
		case m.list.SettingFilter():
			break

		case key.Matches(msg, m.habitKM.add):
			return m, msgs.PushScreen(newAddInput(m.client))

		case key.Matches(msg, m.habitKM.addGroup):
			return m, msgs.PushScreen(group.NewAddGroup())

		case len(m.list.VisibleItems()) == 0:
			break

		default:
			switch sel := m.list.SelectedItem().(type) {
			case listitem.HabitItem:
				selected := sel.Habit
				switch {
				case key.Matches(msg, m.habitKM.view):
					if m.cancelOp != nil {
						m.cancelOp()
						m.cancelOp = nil
					}
					return m, msgs.PushScreen(show.NewShow(m.client, selected, m.width))

				case key.Matches(msg, m.habitKM.rename):
					return m, msgs.PushScreen(textinput.New("New Name:", selected.Name, "habit", selected.ID, m.width))

				case key.Matches(msg, m.habitKM.addToGroup):
					if m.cancelOp != nil {
						m.cancelOp()
						m.cancelOp = nil
					}
					return m, msgs.PushScreen(add_to_group.NewChooseGroup(m.client, selected, m.groups))

				case key.Matches(msg, m.habitKM.delete):
					idx := m.list.Index()
					item := m.list.SelectedItem()
					m.pending = append(m.pending, pendingOp{
						op:   "delete habit",
						id:   selected.ID,
						item: item,
					})
					m.list.RemoveItem(idx)
					ctx, cancel := context.WithCancel(context.Background())
					m.cancelOp = cancel
					return m, tea.Batch(
						msgs.DeleteHabit(ctx, m.client, selected.ID),
						m.list.NewStatusMessage("Deleting "+selected.Name+"..."),
					)

				case key.Matches(msg, m.habitKM.checkIn):
					ctx, cancel := context.WithCancel(context.Background())
					m.cancelOp = cancel
					return m, tea.Batch(
						msgs.CheckIn(ctx, m.client, selected.ID, "", time.Now().Local()),
						m.list.NewStatusMessage("Checking in..."),
					)
				}

			case listitem.GroupItem:
				selected := sel.Group
				switch {
				case key.Matches(msg, m.groupKM.view):
					if m.cancelOp != nil {
						m.cancelOp()
						m.cancelOp = nil
					}
					return m, msgs.PushScreen(group.NewShow(m.client, selected, m.width, m.height))

				case key.Matches(msg, m.groupKM.delete):
					idx := m.list.Index()
					item := m.list.SelectedItem()
					m.pending = append(m.pending, pendingOp{
						op:   "delete group",
						id:   selected.ID,
						item: item,
					})
					m.list.RemoveItem(idx)
					ctx, cancel := context.WithCancel(context.Background())
					m.cancelOp = cancel
					return m, tea.Batch(
						msgs.DeleteGroup(ctx, m.client, selected.ID),
						m.list.NewStatusMessage("Deleting "+selected.Name+"..."),
					)
				}
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}
