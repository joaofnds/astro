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
		log.Fatal(err)
	}
}

func Habits() []*habit.Habit {
	return habits
}

func At(i int) *habit.Habit {
	return habits[i]
}

func IndexOf(name string) int {
	for i, h := range habits {
		if h.Name == name {
			return i
		}
	}
	return -1
}

func Get(name string) *habit.Habit {
	for _, h := range habits {
		if h.Name == name {
			return h
		}
	}
	return nil
}

func SetHabit(h *habit.Habit) {
	for i := range habits {
		if habits[i].Name == h.Name {
			*habits[i] = *h
		}
	}
}

func Add(name string) *habit.Habit {
	_, err := habit.Client.Create(name)
	if err != nil {
		log.Fatal(err)
	}
	GetAll()
	return Get(name)
}

func Delete(name string) error {
	i := IndexOf(name)
	if i == -1 {
		return nil
	}

	if err := habit.Client.Delete(name); err != nil {
		return err
	}

	habits = append(habits[:i], habits[i+1:]...)
	return nil
}
