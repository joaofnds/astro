package listitem

import (
	"astro/config"
	"astro/domain"

	"charm.land/bubbles/v2/list"
)

type GroupItem struct {
	Group      *domain.Group
	activities []domain.Activity
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
	return domain.ShortLineHistogram(i.activities, config.ShortHistSize) + " " + i.lastActivity()
}

func (i GroupItem) lastActivity() string {
	if len(i.activities) == 0 {
		return "no activities"
	}

	return "last activity at " + i.activities[len(i.activities)-1].CreatedAt.Local().Format(config.DateFormat)
}

// PendingGroupItem is a placeholder shown during optimistic group create.
// It renders with dimmed/italic styling until the API confirms creation.
type PendingGroupItem struct {
	Name string
}

func (i PendingGroupItem) Title() string       { return pendingStyle.Render(i.Name) }
func (i PendingGroupItem) Description() string { return pendingStyle.Render("creating...") }
func (i PendingGroupItem) FilterValue() string { return i.Name }

func newGroupItem(g *domain.Group) GroupItem {
	return GroupItem{Group: g, activities: g.Activities()}
}

func GroupsToItems(groups []*domain.Group) []list.Item {
	items := make([]list.Item, len(groups))
	for i, g := range groups {
		items[i] = newGroupItem(g)
	}
	return items
}
