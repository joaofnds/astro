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
	data := []*Habit{}

	res, err := d.api.List()
	if err != nil {
		return data, err
	}
	defer res.Body.Close()

	str, err := io.ReadAll(res.Body)
	if err != nil {
		return data, err
	}

	err = json.Unmarshal(str, &data)
	if err != nil {
		return data, err
	}

	sort.SliceStable(data, func(i, j int) bool {
		return data[i].Name < data[j].Name
	})

	for _, h := range data {
		sort.SliceStable(h.Activities, func(i, j int) bool {
			return h.Activities[i].CreatedAt.Before(h.Activities[j].CreatedAt)
		})
	}

	return data, err
}

func (d *client) Create(name string) error {
	_, err := d.api.Create(name)
	return err
}

func (d *client) Get(name string) (*Habit, error) {
	data := Habit{}

	res, err := d.api.Get(name)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	str, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(str, &data)

	return &data, err
}

func (d *client) CheckIn(name string) (*Habit, error) {
	_, err := d.api.AddActivity(name)
	if err != nil {
		return nil, err
	}

	h, err := d.Get(name)
	return h, err
}
