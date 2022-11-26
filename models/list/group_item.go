package list

import (
	"astro/config"
	"astro/habit"
	"astro/histogram"

	"github.com/charmbracelet/bubbles/list"
)

type groupItem struct {
	group      *habit.Group
	activities []habit.Activity
}

func (i groupItem) FilterValue() string { return i.Title() }

func (i groupItem) Title() string {
	out := i.group.Name

	lenHabits := len(i.group.Habits)
	if lenHabits > 0 {
		out += " ("
		for i, h := range i.group.Habits {
			out += h.Name
			if i < lenHabits-1 {
				out += ", "
			}
		}
		out += ")"
	}
	return out
}

func (i groupItem) Description() string {
	return histogram.ShortLineHistogram(i.activities, config.ShortHistSize) + " " + i.lastActivity()
}

func (i groupItem) lastActivity() string {
	if len(i.activities) == 0 {
		return "no activities"
	}

	return "last activity at " + i.activities[len(i.activities)-1].CreatedAt.Local().Format(config.DateFormat)
}

func newGroupItem(g *habit.Group) groupItem {
	return groupItem{group: g, activities: g.Activities()}
}

func groupsToItems(groups []*habit.Group) []list.Item {
	items := make([]list.Item, len(groups))
	for i, g := range groups {
		items[i] = newGroupItem(g)
	}
	return items
}
