package api

import (
	"astro/domain"
	"fmt"
	"net/url"
	"strings"
	"time"
)

// CheckInDTO carries parameters for a check-in operation.
type CheckInDTO struct {
	ID   string
	Desc string
	Date time.Time
}

func (c *Client) ListHabits() ([]*domain.Habit, error) {
	var habits []*domain.Habit
	if err := c.doRequest("GET", "/habits", nil, &habits); err != nil {
		return nil, err
	}
	domain.SortHabits(habits)
	return habits, nil
}

func (c *Client) CreateHabit(name string) (*domain.Habit, error) {
	var h domain.Habit
	if err := c.doRequest("POST", "/habits?name="+url.QueryEscape(name), nil, &h); err != nil {
		return nil, err
	}
	domain.SortActivities(h.Activities)
	return &h, nil
}

func (c *Client) GetHabit(id string) (*domain.Habit, error) {
	var h domain.Habit
	if err := c.doRequest("GET", "/habits/"+id, nil, &h); err != nil {
		return nil, err
	}
	domain.SortActivities(h.Activities)
	return &h, nil
}

func (c *Client) UpdateHabit(id, name string) error {
	body := fmt.Sprintf(`{"name":%q}`, name)
	return c.doRequest("PATCH", "/habits/"+id, strings.NewReader(body), nil)
}

func (c *Client) DeleteHabit(id string) error {
	return c.doRequest("DELETE", "/habits/"+id, nil, nil)
}

func (c *Client) AddActivity(id, desc string, date time.Time) error {
	body := fmt.Sprintf(`{"description":%q,"date":%q}`, desc, date.UTC().Format(time.RFC3339))
	return c.doRequest("POST", "/habits/"+id, strings.NewReader(body), nil)
}

func (c *Client) UpdateActivity(habitID, activityID, desc string) error {
	body := fmt.Sprintf(`{"description":%q}`, desc)
	return c.doRequest("PATCH", "/habits/"+habitID+"/"+activityID, strings.NewReader(body), nil)
}

func (c *Client) DeleteActivity(habitID, activityID string) error {
	return c.doRequest("DELETE", "/habits/"+habitID+"/"+activityID, nil, nil)
}

func (c *Client) CheckIn(dto CheckInDTO) (*domain.Habit, error) {
	if err := c.AddActivity(dto.ID, dto.Desc, dto.Date); err != nil {
		return nil, err
	}
	return c.GetHabit(dto.ID)
}
