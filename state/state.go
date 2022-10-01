package state

import (
	"astroapp/habit"
	"log"
)

var habits []*habit.Habit

func init() {
	GetAll()
}

func GetAll() {
	var err error
	habits, err = habit.Client.List()
	if err != nil {
		log.Fatalf("could not list habits: %s", err)
	}
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
