package listitem

import (
	"astro/config"
	"astro/habit"
	"astro/histogram"

	"github.com/charmbracelet/bubbles/list"
)

type GroupItem struct {
	Group      *habit.Group
	activities []habit.Activity
}

func (i GroupItem) FilterValue() string { return i.Title() }

func (i GroupItem) Title() string {
	out := i.Group.Name

	lenHabits := len(i.Group.Habits)
	if lenHabits > 0 {
		out += " ("
		for i, h := range i.Group.Habits {
			out += h.Name
			if i < lenHabits-1 {
				out += ", "
			}
		}
		out += ")"
	}
	return out
}

func (i GroupItem) Description() string {
	return histogram.ShortLineHistogram(i.activities, config.ShortHistSize) + " " + i.lastActivity()
}

func (i GroupItem) lastActivity() string {
	if len(i.activities) == 0 {
		return "no activities"
	}

	return "last activity at " + i.activities[len(i.activities)-1].CreatedAt.Local().Format(config.DateFormat)
}

func newGroupItem(g *habit.Group) GroupItem {
	return GroupItem{Group: g, activities: g.Activities()}
}

func GroupsToItems(groups []*habit.Group) []list.Item {
	items := make([]list.Item, len(groups))
	for i, g := range groups {
		items[i] = newGroupItem(g)
	}
	return items
}
