package domain

import (
	"astro/config"
	"astro/date"
	"fmt"
	"sort"
	"strings"
	"time"
)

type Group struct {
	ID     string   `json:"id"`
	Name   string   `json:"name"`
	Habits []*Habit `json:"habits"`
}

func SortGroups(groups []*Group) {
	sort.SliceStable(groups, func(i, j int) bool {
		return groups[i].Name < groups[j].Name
	})

	for _, g := range groups {
		SortHabits(g.Habits)
	}
}

func (g *Group) Activities() []Activity {
	var size int
	for _, h := range g.Habits {
		size += len(h.Activities)
	}
	activities := make([]Activity, 0, size)
	for _, h := range g.Habits {
		activities = append(activities, h.Activities...)
	}
	sort.SliceStable(activities, func(i, j int) bool {
		return activities[i].CreatedAt.Before(activities[j].CreatedAt)
	})
	return activities
}

func ActivitiesOnDateTally(habits []*Habit, t time.Time) string {
	habitCount := map[string]int{}
	var total int

	for _, h := range habits {
		for _, a := range h.Activities {
			if date.SameDay(a.CreatedAt, t) {
				habitCount[h.Name]++
				total++
			}
		}
	}

	act := "activities"
	if total == 1 {
		act = "activity"
	}

	counts := make([]string, 0, total)
	for name, count := range habitCount {
		counts = append(counts, fmt.Sprintf("%s:%d", name, count))
	}
	var countsStr string
	if total > 0 {
		countsStr = " (" + strings.Join(counts, ", ") + ")"
	}
	return fmt.Sprintf(
		"%d %s on %s%s\n",
		total,
		act,
		t.Format(config.DateFormat),
		countsStr,
	)
}
