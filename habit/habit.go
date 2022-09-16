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
