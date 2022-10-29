package state

import (
	"astro/habit"
	"log"
)

var habits []*habit.Habit

func GetAll() error {
	var err error
	habits, err = habit.Client.List()
	return err
}

func Habits() []*habit.Habit {
	return habits
}

func At(i int) *habit.Habit {
	return habits[i]
}

func IndexOf(id string) int {
	for i, h := range habits {
		if h.ID == id {
			return i
		}
	}
	return -1
}

func Get(id string) *habit.Habit {
	for _, h := range habits {
		if h.ID == id {
			return h
		}
	}
	return nil
}

func SetHabit(h *habit.Habit) {
	for i := range habits {
		if habits[i].ID == h.ID {
			*habits[i] = *h
		}
	}
}

func UpdateActivity(h *habit.Habit, activity *habit.Activity) {
	for i, a := range h.Activities {
		if a.Id == activity.Id {
			h.Activities[i] = *activity
		}
	}
}

func DeleteActivity(h *habit.Habit, activity habit.Activity) {
	for i, a := range h.Activities {
		if a.Id == activity.Id {
			h.Activities = append(h.Activities[:i], h.Activities[i+1:]...)
		}
	}
}

func Add(name string) *habit.Habit {
	h, err := habit.Client.Create(name)
	if err != nil {
		log.Fatalf("could not create habit: %s", err)
	}
	GetAll()
	return Get(h.ID)
}

func Delete(id string) error {
	i := IndexOf(id)
	if i == -1 {
		return nil
	}

	if err := habit.Client.Delete(id); err != nil {
		return err
	}

	habits = append(habits[:i], habits[i+1:]...)
	return nil
}
