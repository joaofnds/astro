package listitem

import (
	"astro/config"
	"astro/domain"

	"charm.land/bubbles/v2/list"
)

type GroupItem struct {
	Group *domain.Group
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
	activities := i.Group.Activities()
	return domain.ShortLineHistogram(activities, config.ShortHistSize) + " " + lastActivityLine(activities)
}

func lastActivityLine(activities []domain.Activity) string {
	if len(activities) == 0 {
		return "no activities"
	}
	return "last activity at " + activities[len(activities)-1].CreatedAt.Local().Format(config.DateFormat)
}

type PendingGroupItem struct {
	Name string
}

func (i PendingGroupItem) Title() string       { return pendingStyle.Render(i.Name) }
func (i PendingGroupItem) Description() string { return pendingStyle.Render("creating...") }
func (i PendingGroupItem) FilterValue() string { return i.Name }

func GroupsToItems(groups []*domain.Group) []list.Item {
	items := make([]list.Item, len(groups))
	for i, g := range groups {
		items[i] = GroupItem{Group: g}
	}
	return items
}
