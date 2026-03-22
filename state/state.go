package state

import (
	"astro/api"
	"astro/domain"
	"log"
	"time"
)

var (
	client *api.Client
	groups []*domain.Group
	habits []*domain.Habit
)

// Init sets the API client for all state operations.
func Init(c *api.Client) {
	client = c
}

func GetAll() {
	var err error
	groups, habits, err = client.GroupsAndHabits()
	if err != nil {
		log.Fatal(err)
	}
}

func Habits() []*domain.Habit {
	return habits
}

func Groups() []*domain.Group {
	return groups
}

func At(i int) *domain.Habit {
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

func Get(id string) *domain.Habit {
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

func SetHabit(h *domain.Habit) {
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

func UpdateActivity(h *domain.Habit, activity *domain.Activity) {
	for i, a := range h.Activities {
		if a.ID == activity.ID {
			h.Activities[i] = *activity
		}
	}
}

func DeleteActivity(h *domain.Habit, activity domain.Activity) {
	for i, a := range h.Activities {
		if a.ID == activity.ID {
			h.Activities = append(h.Activities[:i], h.Activities[i+1:]...)
		}
	}
}

func Add(name string) *domain.Habit {
	h, err := client.CreateHabit(name)
	if err != nil {
		log.Fatalf("could not create habit: %s", err)
	}
	GetAll()
	return Get(h.ID)
}

func AddGroup(name string) *domain.Habit {
	g, err := client.CreateGroup(name)
	if err != nil {
		log.Fatalf("could not create group: %s", err)
	}
	GetAll()
	return Get(g.ID)
}

func DeleteGroup(group domain.Group) error {
	err := client.DeleteGroup(group.ID)
	if err != nil {
		log.Fatalf("could not delete group: %s", err)
	}
	GetAll()
	return nil
}

func AddToGroup(h domain.Habit, g domain.Group) {
	err := client.AddToGroup(h.ID, g.ID)
	if err != nil {
		log.Fatalf("could not add to group: %s", err)
	}
	GetAll()
}

func RemoveFromGroup(h domain.Habit, g domain.Group) {
	if err := client.RemoveFromGroup(h.ID, g.ID); err != nil {
		log.Fatalf("could not remove habit from group: %s", err)
	}
	GetAll()
}

func Delete(id string) error {
	i := IndexOf(id)
	if i == -1 {
		return nil
	}

	if err := client.DeleteHabit(id); err != nil {
		return err
	}

	habits = append(habits[:i], habits[i+1:]...)
	return nil
}

func CheckIn(id, desc string, date time.Time) (*domain.Habit, error) {
	return client.CheckIn(api.CheckInDTO{ID: id, Desc: desc, Date: date})
}

func UpdateHabit(h *domain.Habit) error {
	return client.UpdateHabit(h.ID, h.Name)
}

func UpdateHabitActivity(habitID, activityID, desc string) error {
	return client.UpdateActivity(habitID, activityID, desc)
}

func DeleteHabitActivity(habitID, activityID string) error {
	return client.DeleteActivity(habitID, activityID)
}
