package listitem

import (
	"astro/config"
	"astro/domain"

	"charm.land/bubbles/v2/list"
	"charm.land/lipgloss/v2"
)

type HabitItem struct{ Habit *domain.Habit }

func (i HabitItem) FilterValue() string { return i.Title() }

func (i HabitItem) Title() string { return i.Habit.Name }

func (i HabitItem) Description() string {
	return domain.ShortLineHistogram(i.Habit.Activities, config.ShortHistSize) + " " + i.lastActivity()
}

func (i HabitItem) lastActivity() string {
	if len(i.Habit.Activities) == 0 {
		return "no activities"
	}

	return "last activity at " + i.Habit.LatestActivity().Local().Format(config.DateFormat)
}

// pendingStyle renders text as dimmed and italic, used for optimistic
// create placeholders that haven't been confirmed by the API yet.
var pendingStyle = lipgloss.NewStyle().Faint(true).Italic(true)

// PendingHabitItem is a placeholder shown during optimistic create.
// It renders with dimmed/italic styling until the API confirms creation.
type PendingHabitItem struct {
	Name string
}

func (i PendingHabitItem) Title() string       { return pendingStyle.Render(i.Name) }
func (i PendingHabitItem) Description() string { return pendingStyle.Render("creating...") }
func (i PendingHabitItem) FilterValue() string { return i.Name }

func HabitsToItems(habits []*domain.Habit) []list.Item {
	items := make([]list.Item, len(habits))
	for i, h := range habits {
		items[i] = HabitItem{h}
	}
	return items
}
