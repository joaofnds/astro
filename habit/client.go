package habit

import (
	"encoding/json"
	"io"
	"sort"
)

type Client struct {
	api *API
}

func NewClient() *Client {
	return &Client{NewAPI()}
}

func (d *Client) List() ([]Habit, error) {
	data := []Habit{}

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

	for _, h := range data {
		sort.SliceStable(h.Activities, func(i, j int) bool {
			return h.Activities[i].CreatedAt.Before(h.Activities[j].CreatedAt)
		})
	}

	return data, err
}

func (d *Client) Create(name string) error {
	_, err := d.api.Create(name)
	return err
}

func (d *Client) Get(name string) (Habit, error) {
	data := Habit{}

	res, err := d.api.Get(name)
	if err != nil {
		return data, err
	}
	defer res.Body.Close()

	str, err := io.ReadAll(res.Body)
	if err != nil {
		return data, err
	}

	err = json.Unmarshal(str, &data)

	return data, err
}
