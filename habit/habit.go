package habit

import (
	"astro/date"
	"errors"
	"sort"
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

func sortActivities(h *Habit) {
	sort.SliceStable(h.Activities, func(i, j int) bool {
		return h.Activities[i].CreatedAt.Before(h.Activities[j].CreatedAt)
	})
}
