package habit

import (
	"sort"
	"time"
)

type Activity struct {
	Id        string    `json:"id"`
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

func sortActivities(h *Habit) {
	sort.SliceStable(h.Activities, func(i, j int) bool {
		return h.Activities[i].CreatedAt.Before(h.Activities[j].CreatedAt)
	})
}
