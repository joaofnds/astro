package state

import (
	"astroapp/habit"
	"log"
)

var habits []*habit.Habit

func init() {
	var err error
	habits, err = habit.Client.List()
	if err != nil {
		log.Fatal(err)
	}
}

func Habits() []*habit.Habit {
	return habits
}

func At(i int) *habit.Habit {
	return habits[i]
}

func SetHabit(h *habit.Habit) {
	for i := range habits {
		if habits[i].Name == h.Name {
			*habits[i] = *h
		}
	}
}
