package state

import (
	"astro/habit"
	"log"
)

var (
	groups []*habit.Group
	habits []*habit.Habit
)

func GetAll() {
	var err error
	groups, habits, err = habit.Client.GroupsAndHabits()
	if err != nil {
		log.Fatal(err)
	}
}

func Habits() []*habit.Habit {
	return habits
}

func Groups() []*habit.Group {
	return groups
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
	for _, g := range groups {
		for _, h := range g.Habits {
			if h.ID == id {
				return h
			}
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
	for _, g := range groups {
		for i := range g.Habits {
			if g.Habits[i].ID == h.ID {
				*g.Habits[i] = *h
			}
		}
	}
}

func UpdateActivity(h *habit.Habit, activity *habit.Activity) {
	for i, a := range h.Activities {
		if a.ID == activity.ID {
			h.Activities[i] = *activity
		}
	}
}

func DeleteActivity(h *habit.Habit, activity habit.Activity) {
	for i, a := range h.Activities {
		if a.ID == activity.ID {
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

func AddGroup(name string) *habit.Habit {
	h, err := habit.Client.CreateGroup(name)
	if err != nil {
		log.Fatalf("could not create group: %s", err)
	}
	GetAll()
	return Get(h.ID)
}

func DeleteGroup(group habit.Group) error {
	err := habit.Client.DeleteGroup(group)
	if err != nil {
		log.Fatalf("could not delete group: %s", err)
	}
	GetAll()
	return nil
}

func AddToGroup(h habit.Habit, g habit.Group) {
	err := habit.Client.AddToGroup(h, g)
	if err != nil {
		log.Fatalf("could not create group: %s", err)
	}
	GetAll()
}

func RemoveFromGroup(h habit.Habit, g habit.Group) {
	if err := habit.Client.RemoveFromGroup(h, g); err != nil {
		log.Fatalf("could not remove habti from group: %s", err)
	}
	GetAll()
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
