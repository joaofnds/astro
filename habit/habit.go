package habit

import "time"

type Activity struct {
	Id        int       `json:"id"`
	CreatedAt time.Time `json:"created_at"`
}

type Habit struct {
	Id         int        `json:"id"`
	Name       string     `json:"name"`
	Activities []Activity `json:"activities"`
}

func (h Habit) LatestActivity() time.Time {
	if len(h.Activities) == 0 {
		return time.Time{}
	}

	return h.Activities[len(h.Activities)-1].CreatedAt
}
