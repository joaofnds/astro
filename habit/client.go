package habit

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
)

var Client *client

func InitClient(token string) {
	Client = NewClient(token)
}

type client struct {
	api   *API
	token string
}

func NewClient(token string) *client {
	return &client{NewAPI(), token}
}

func (d *client) List() ([]*Habit, error) {
	res, err := d.api.List(d.token)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	str, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var habits []*Habit
	err = json.Unmarshal(str, &habits)
	if err != nil {
		return nil, err
	}

	sort.SliceStable(habits, func(i, j int) bool {
		return habits[i].Name < habits[j].Name
	})

	for _, h := range habits {
		sortActivities(h.Activities)
	}

	return habits, err
}

func (d *client) Create(name string) (*Habit, error) {
	var h Habit
	res, err := d.api.Create(d.token, name)
	if err != nil {
		return &h, err
	}
	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return &h, err
	}

	err = json.Unmarshal(b, &h)
	if err != nil {
		return &h, err
	}

	sortActivities(h.Activities)
	return &h, nil
}

func (d *client) Update(habit *Habit) error {
	_, err := d.api.Update(d.token, habit.ID, habit.Name)
	return err
}

func (d *client) Delete(id string) error {
	_, err := d.api.Delete(d.token, id)
	return err
}

func (d *client) UpdateActivity(habit Habit, activity Activity) error {
	_, err := d.api.UpdateActivity(d.token, habit.ID, activity.ID, activity.Desc)
	return err
}

func (d *client) DeleteActivity(habit Habit, activity Activity) error {
	_, err := d.api.DeleteActivity(d.token, habit.ID, activity.ID)
	return err
}

func (d *client) CreateGroup(name string) (Group, error) {
	var group Group

	res, err := d.api.CreateGroup(d.token, name)
	if err != nil {
		return group, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		return group, fmt.Errorf("failed to create group (code %d != %d)", res.StatusCode, http.StatusCreated)
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return group, fmt.Errorf("failed to create group: %s", err)
	}

	return group, json.Unmarshal(b, &group)
}

func (d *client) AddToGroup(habit Habit, group Group) error {
	res, err := d.api.AddToGroup(d.token, habit.ID, group.ID)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to create group (code %d != %d)", res.StatusCode, http.StatusCreated)
	}

	return nil
}

func (d *client) RemoveFromGroup(habit Habit, group Group) error {
	res, err := d.api.RemoveFromGroup(d.token, habit.ID, group.ID)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to remove from group (code %d != %d)", res.StatusCode, http.StatusOK)
	}

	return nil
}

func (d *client) DeleteGroup(group Group) error {
	res, err := d.api.DeleteGroup(d.token, group.ID)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to delete group (code %d != %d)", res.StatusCode, http.StatusOK)
	}

	return nil
}

func (d *client) GroupsAndHabits() ([]*Group, []*Habit, error) {
	res, err := d.api.GroupsAndHabits(d.token)
	if err != nil {
		return nil, nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("failed to create group (code %d != %d)", res.StatusCode, http.StatusOK)
	}

	str, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, nil, err
	}

	data := GroupsAndHabitsPayload{}

	if err := json.Unmarshal(str, &data); err != nil {
		return nil, nil, err
	}

	sortHabits(data.Habits)
	sortGroups(data.Groups)

	return data.Groups, data.Habits, nil
}

func (d *client) Get(id string) (*Habit, error) {
	h := Habit{}

	res, err := d.api.Get(d.token, id)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(b, &h)
	sortActivities(h.Activities)
	return &h, err
}

func (d *client) CheckIn(id, desc string) (*Habit, error) {
	_, err := d.api.AddActivity(d.token, id, desc)
	if err != nil {
		return nil, err
	}

	h, err := d.Get(id)
	return h, err
}

type GroupsAndHabitsPayload struct {
	Groups []*Group `json:"groups"`
	Habits []*Habit `json:"habits"`
}
