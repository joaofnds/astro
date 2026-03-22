package listitem

import (
	"astro/config"
	"astro/domain"

	"charm.land/bubbles/v2/list"
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

func HabitsToItems(habits []*domain.Habit) []list.Item {
	items := make([]list.Item, len(habits))
	for i, h := range habits {
		items[i] = HabitItem{h}
	}
	return items
}
