package habit

import (
	"astro/config"
	"astro/date"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"
)

type Activity struct {
	Id        string    `json:"id"`
	Desc      string    `json:"description"`
	CreatedAt time.Time `json:"created_at"`
}

type Habit struct {
	ID         string     `json:"id"`
	Name       string     `json:"name"`
	UserID     string     `json:"user_id"`
	Activities []Activity `json:"activities"`
}

type Group struct {
	ID     string   `json:"id"`
	Name   string   `json:"name"`
	Habits []*Habit `json:"habits"`
}

func (h Habit) LatestActivity() time.Time {
	if len(h.Activities) == 0 {
		return time.Time{}
	}

	return h.Activities[len(h.Activities)-1].CreatedAt
}

func (h Habit) LatestActivityOnDate(time time.Time) (Activity, error) {
	for i := len(h.Activities) - 1; i >= 0; i-- {
		if date.SameDay(time, h.Activities[i].CreatedAt) {
			return h.Activities[i], nil
		}
	}

	return Activity{}, errors.New("no activity on date")
}

func sortActivities(activities []Activity) {
	sort.SliceStable(activities, func(i, j int) bool {
		return activities[i].CreatedAt.Before(activities[j].CreatedAt)
	})
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

func ActivitiesOnDate(activities []Activity, t time.Time) string {
	var count int
	for _, a := range activities {
		if date.SameDay(a.CreatedAt, t) {
			count++
		}
	}
	w := "activities"
	if count == 1 {
		w = "activity"
	}
	return fmt.Sprintf("%d %s on %s\n", count, w, t.Format(config.DateFormat))
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
