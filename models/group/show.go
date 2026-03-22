package group

import (
	"astro/api"
	"astro/config"
	"astro/date"
	"astro/domain"
	"astro/models/listitem"
	"astro/models/show"
	"astro/models/textinput"
	"astro/msgs"
	"astro/util"
	"context"
	"fmt"
	"strings"
	"time"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
)

// pendingOp tracks an in-flight optimistic mutation for rollback on failure.
type pendingOp struct {
	op   string
	id   string
	item list.Item
	name string
}

// CreateGroupSubmit is sent when the user submits a new group name from the add input.
// Exported because models/list/list.go handles this message to insert PendingGroupItem.
type CreateGroupSubmit struct {
	Name string
}

type List struct {
	client       *api.Client
	group        *domain.Group
	list         list.Model
	help         help.Model
	km           binds
	t            time.Time
	selected     int
	lastSelected int
	onHist       bool
	width        int
	height       int
	pending      []pendingOp
	cancelOp     context.CancelFunc
}

func NewShow(client *api.Client, g *domain.Group, width, height int) List {
	l := list.New(listitem.HabitsToItems(g.Habits), list.NewDefaultDelegate(), 0, 5)
	l.SetSize(width, height-9)
	l.StatusMessageLifetime = 3 * time.Second

	km := newBinds()
	l.AdditionalShortHelpKeys = km.ToSlice
	l.Title = g.Name

	t, _ := date.TimeFrame()
	selected := date.DiffInDays(t, date.Today()) + config.TimeFrameInDays

	h := help.New()
	h.SetWidth(width)

	return List{
		client:   client,
		t:        t,
		selected: selected,
		group:    g,
		list:     l,
		km:       km,
		help:     h,
		width:    width,
		height:   height,
	}
}

func (m List) Init() tea.Cmd {
	return nil
}

func (m List) View() tea.View {
	activities := m.group.Activities()
	var s strings.Builder
	s.WriteString(domain.Histogram(m.t, activities, m.selected))

	if m.selectedDate().After(date.Today()) {
		s.WriteString("\n")
	} else {
		s.WriteString(domain.ActivitiesOnDateTally(m.group.Habits, m.selectedDate()))
	}

	m.list.Title = domain.Digest(m.group.Name, m.group.Activities())

	s.WriteString("\n")
	s.WriteString(m.list.View())
	s.WriteString("\n")
	s.WriteString(m.help.View(m.km))
	return tea.NewView(s.String())
}

func (m List) selectedDate() time.Time {
	return m.t.AddDate(0, 0, m.selected)
}

func (m List) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case msgs.CheckInResultMsg:
		for i, h := range m.group.Habits {
			if h.ID == msg.Habit.ID {
				m.group.Habits[i] = msg.Habit
				break
			}
		}
		cmds = append(cmds, m.list.SetItems(listitem.HabitsToItems(m.group.Habits)))
		cmds = append(cmds, m.list.NewStatusMessage("Checked in"))

	case msgs.HabitUpdatedMsg:
		// Clear pending for rename on success.
		for i, p := range m.pending {
			if p.op == "update habit" && p.id == msg.Habit.ID {
				m.pending = append(m.pending[:i], m.pending[i+1:]...)
				break
			}
		}
		for i, h := range m.group.Habits {
			if h.ID == msg.Habit.ID {
				m.group.Habits[i] = msg.Habit
				break
			}
		}
		cmds = append(cmds, m.list.SetItems(listitem.HabitsToItems(m.group.Habits)))
		cmds = append(cmds, m.list.NewStatusMessage("Renamed"))

	case msgs.RemovedFromGroupMsg:
		// Already removed optimistically. Clean up pending and group.Habits.
		for i, p := range m.pending {
			if p.op == "remove from group" && p.id == msg.HabitID {
				m.pending = append(m.pending[:i], m.pending[i+1:]...)
				break
			}
		}
		for i, h := range m.group.Habits {
			if h.ID == msg.HabitID {
				m.group.Habits = append(m.group.Habits[:i], m.group.Habits[i+1:]...)
				break
			}
		}

	case msgs.APIErrorMsg:
		errStr := fmt.Sprintf("\u2717 %s: %s", msg.Op, msg.Err)
		for i, p := range m.pending {
			if p.op == msg.Op && (msg.ID == "" || p.id == msg.ID) {
				m.pending = append(m.pending[:i], m.pending[i+1:]...)
				switch p.op {
				case "remove from group":
					if p.item != nil {
						cmds = append(cmds, m.list.InsertItem(len(m.list.Items()), p.item))
						if hi, ok := p.item.(listitem.HabitItem); ok {
							m.group.Habits = append(m.group.Habits, hi.Habit)
						}
						errStr += " \u2014 restored"
					}
				case "update habit":
					// Revert optimistic rename.
					for j, h := range m.group.Habits {
						if h.ID == p.id {
							m.group.Habits[j].Name = p.name
							break
						}
					}
					cmds = append(cmds, m.list.SetItems(listitem.HabitsToItems(m.group.Habits)))
					errStr += " \u2014 restored"
				}
				break
			}
		}
		cmds = append(cmds, m.list.NewStatusMessage(errStr))

	case textinput.Submit:
		switch msg.Key {
		case "habit":
			// Store old state for rollback.
			for _, item := range m.list.Items() {
				if hi, ok := item.(listitem.HabitItem); ok && hi.Habit.ID == msg.ID {
					m.pending = append(m.pending, pendingOp{
						op:   "update habit",
						id:   msg.ID,
						item: item,
						name: hi.Habit.Name,
					})
					// Optimistic rename.
					updated := *hi.Habit
					updated.Name = msg.Value
					for j, h := range m.group.Habits {
						if h.ID == msg.ID {
							m.group.Habits[j] = &updated
							break
						}
					}
					cmds = append(cmds, m.list.SetItems(listitem.HabitsToItems(m.group.Habits)))
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

		case key.Matches(msg, m.km.tab):
			m.onHist = !m.onHist

			if m.onHist {
				m.lastSelected = m.list.Index()
				m.list.Select(-1)
				m.selected -= config.TimeFrameInDays
			} else {
				m.list.Select(m.lastSelected)
				m.selected += config.TimeFrameInDays
			}

		// ClearScreen forces a full sequential redraw on navigation.
		// See comment in show.go for details on the emoji width issue.
		case m.onHist && key.Matches(msg, m.km.left):
			m.selected = util.Max(m.selected-7, 0)
			return m, tea.ClearScreen

		case m.onHist && key.Matches(msg, m.km.right):
			m.selected = util.Min(m.selected+7, config.TimeFrameInDays-1)
			return m, tea.ClearScreen

		case m.onHist && key.Matches(msg, m.km.up) && m.selected > 0:
			m.selected -= 1
			return m, tea.ClearScreen

		case m.onHist && key.Matches(msg, m.km.down) && (m.selected+1) < config.TimeFrameInDays:
			m.selected += 1
			return m, tea.ClearScreen

		case key.Matches(msg, m.km.quit):
			if m.cancelOp != nil {
				m.cancelOp()
				m.cancelOp = nil
			}
			return m, msgs.PopScreen()

		case len(m.list.VisibleItems()) == 0:
			break

		default:
			sel, ok := m.list.SelectedItem().(listitem.HabitItem)
			if !ok {
				break
			}
			selected := sel.Habit

			switch {
			case key.Matches(msg, m.km.checkIn):
				ctx, cancel := context.WithCancel(context.Background())
				m.cancelOp = cancel
				return m, tea.Batch(
					msgs.CheckIn(ctx, m.client, selected.ID, "", time.Now().Local()),
					m.list.NewStatusMessage("Checking in..."),
				)

			case key.Matches(msg, m.km.view):
				if m.cancelOp != nil {
					m.cancelOp()
					m.cancelOp = nil
				}
				return m, msgs.PushScreen(show.NewShow(m.client, selected, m.width))

			case key.Matches(msg, m.km.rename):
				return m, msgs.PushScreen(textinput.New("New Name:", selected.Name, "habit", selected.ID, m.width))

			case key.Matches(msg, m.km.delete):
				idx := m.list.Index()
				item := m.list.SelectedItem()
				m.pending = append(m.pending, pendingOp{
					op:   "remove from group",
					id:   selected.ID,
					item: item,
				})
				m.list.RemoveItem(idx)
				ctx, cancel := context.WithCancel(context.Background())
				m.cancelOp = cancel
				return m, tea.Batch(
					msgs.RemoveFromGroup(ctx, m.client, selected.ID, m.group.ID),
					m.list.NewStatusMessage("Removing "+selected.Name+"..."),
				)
			}
		}
	}

	if !m.onHist {
		var cmd tea.Cmd
		m.list, cmd = m.list.Update(msg)
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}
