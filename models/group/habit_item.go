package group

import (
	"astro/config"
	"astro/habit"
	"astro/histogram"

	"github.com/charmbracelet/bubbles/list"
)

type habitItem struct{ habit *habit.Habit }

func (i habitItem) FilterValue() string { return i.Title() }

func (i habitItem) Title() string { return i.habit.Name }

func (i habitItem) Description() string {
	return histogram.ShortLineHistogram(i.habit.Activities, config.ShortHistSize) + " " + i.lastActivity()
}

func (i habitItem) lastActivity() string {
	if len(i.habit.Activities) == 0 {
		return "no activities"
	}

	return "last activity at " + i.habit.LatestActivity().Local().Format(config.DateFormat)
}

func toItems(habits []*habit.Habit) []list.Item {
	items := make([]list.Item, len(habits))
	for i, h := range habits {
		items[i] = toItem(h)
	}
	return items
}

func toItem(h *habit.Habit) list.Item {
	return habitItem{h}
}
