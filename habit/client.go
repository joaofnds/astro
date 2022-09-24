package habit

import (
	"encoding/json"
	"io"
	"sort"
)

var Client *client

func init() {
	Client = NewClient()
}

type client struct {
	api *API
}

func NewClient() *client {
	return &client{NewAPI()}
}

func (d *client) List() ([]*Habit, error) {
	res, err := d.api.List()
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	str, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	habits := []*Habit{}
	err = json.Unmarshal(str, &habits)
	if err != nil {
		return nil, err
	}

	sort.SliceStable(habits, func(i, j int) bool {
		return habits[i].Name < habits[j].Name
	})

	for _, h := range habits {
		sortActivities(h)
	}

	return habits, err
}

func (d *client) Create(name string) (*Habit, error) {
	var h Habit
	res, err := d.api.Create(name)
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

	sortActivities(&h)
	return &h, nil
}

func (d *client) Delete(name string) error {
	_, err := d.api.Delete(name)
	return err
}

func (d *client) Get(name string) (*Habit, error) {
	h := Habit{}

	res, err := d.api.Get(name)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(b, &h)
	sortActivities(&h)
	return &h, err
}

func (d *client) CheckIn(name string) (*Habit, error) {
	_, err := d.api.AddActivity(name)
	if err != nil {
		return nil, err
	}

	h, err := d.Get(name)
	return h, err
}
