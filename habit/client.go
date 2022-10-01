package habit

import (
	"astro/config"
	"encoding/json"
	"io"
	"sort"
)

var Client *client

func init() {
	Client = NewClient(string(config.Token))
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

	sortActivities(&h)
	return &h, nil
}

func (d *client) Delete(name string) error {
	_, err := d.api.Delete(d.token, name)
	return err
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
	sortActivities(&h)
	return &h, err
}

func (d *client) CheckIn(id string) (*Habit, error) {
	_, err := d.api.AddActivity(d.token, id)
	if err != nil {
		return nil, err
	}

	h, err := d.Get(id)
	return h, err
}
