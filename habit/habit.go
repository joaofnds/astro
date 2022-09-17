package habit

import "time"

type Activity struct {
	Id        int       `json:"id"`
	CreatedAt time.Time `json:"created_at"`
}

type Habit struct {
	Id        int        `json:"id"`
	Name      string     `json:"name"`
	Activites []Activity `json:"activities"`
}

func (h Habit) LatestActivity() time.Time {
	if len(h.Activites) == 0 {
		return time.Time{}
	}

	return h.Activites[len(h.Activites)-1].CreatedAt
}
