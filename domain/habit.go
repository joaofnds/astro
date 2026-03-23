package domain

import (
	"astro/config"
	"astro/date"
	"errors"
	"fmt"
	"sort"
	"time"
)

type Activity struct {
	ID        string    `json:"id"`
	Desc      string    `json:"description"`
	CreatedAt time.Time `json:"created_at"`
}

type Habit struct {
	ID         string     `json:"id"`
	Name       string     `json:"name"`
	UserID     string     `json:"user_id"`
	Activities []Activity `json:"activities"`
}

func (h Habit) LatestActivity() time.Time {
	if len(h.Activities) == 0 {
		return time.Time{}
	}

	return h.Activities[len(h.Activities)-1].CreatedAt
}

func Streak(activities []Activity) string {
	streak := CurrentStreakDays(activities)
	if streak == 1 {
		return "1 day"
	}
	return fmt.Sprintf("%d days", streak)
}

func CurrentStreakDays(activities []Activity) int {
	if len(activities) == 0 {
		return 0
	}

	var streak int
	cur := date.Today()

	i := len(activities) - 1
	if !date.SameDay(activities[i].CreatedAt, cur) {
		cur = cur.AddDate(0, 0, -1)
	}

	for i >= 0 {
		if !date.SameDay(activities[i].CreatedAt, cur) {
			break
		}

		for date.SameDay(activities[i].CreatedAt, cur) && i > 0 {
			i--
		}

		streak++
		cur = cur.AddDate(0, 0, -1)
	}

	return streak
}

func Momentum(activities []Activity) int {
	if len(activities) == 0 {
		return 0
	}

	var momentum int

	idx := 0
	day := date.TruncateDay(activities[idx].CreatedAt)

	for idx < len(activities) && !day.After(date.Today()) {

		if date.SameDay(activities[idx].CreatedAt, day) {
			for idx < len(activities) && date.SameDay(activities[idx].CreatedAt, day) {
				momentum++
				idx++
			}
		} else if momentum > 0 {
			momentum--
		}

		day = day.AddDate(0, 0, 1)
	}

	return momentum
}

func Digest(name string, activities []Activity) string {
	return fmt.Sprintf(
		"%s - streak: %s, momentum: %d",
		name,
		Streak(activities),
		Momentum(activities),
	)
}

func (h Habit) LatestActivityOnDate(time time.Time) (Activity, error) {
	for i := len(h.Activities) - 1; i >= 0; i-- {
		if date.SameDay(time, h.Activities[i].CreatedAt) {
			return h.Activities[i], nil
		}
	}

	return Activity{}, errors.New("no activity on date")
}

func SortHabits(habits []*Habit) {
	sort.SliceStable(habits, func(i, j int) bool {
		return habits[i].Name < habits[j].Name
	})

	for _, h := range habits {
		SortActivities(h.Activities)
	}
}

func SortActivities(activities []Activity) {
	sort.SliceStable(activities, func(i, j int) bool {
		return activities[i].CreatedAt.Before(activities[j].CreatedAt)
	})
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
