package api

import (
	"astro/domain"
	"context"
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

func (c *Client) ListHabits(ctx context.Context) ([]*domain.Habit, error) {
	var habits []*domain.Habit
	if err := c.doRequest(ctx, "GET", "/habits", nil, &habits); err != nil {
		return nil, err
	}
	domain.SortHabits(habits)
	return habits, nil
}

func (c *Client) CreateHabit(ctx context.Context, name string) (*domain.Habit, error) {
	var h domain.Habit
	if err := c.doRequest(ctx, "POST", "/habits?name="+url.QueryEscape(name), nil, &h); err != nil {
		return nil, err
	}
	domain.SortActivities(h.Activities)
	return &h, nil
}

func (c *Client) GetHabit(ctx context.Context, id string) (*domain.Habit, error) {
	var h domain.Habit
	if err := c.doRequest(ctx, "GET", "/habits/"+id, nil, &h); err != nil {
		return nil, err
	}
	domain.SortActivities(h.Activities)
	return &h, nil
}

func (c *Client) UpdateHabit(ctx context.Context, id, name string) error {
	body := fmt.Sprintf(`{"name":%q}`, name)
	return c.doRequest(ctx, "PATCH", "/habits/"+id, strings.NewReader(body), nil)
}

func (c *Client) DeleteHabit(ctx context.Context, id string) error {
	return c.doRequest(ctx, "DELETE", "/habits/"+id, nil, nil)
}

func (c *Client) AddActivity(ctx context.Context, id, desc string, date time.Time) error {
	body := fmt.Sprintf(`{"description":%q,"date":%q}`, desc, date.UTC().Format(time.RFC3339))
	return c.doRequest(ctx, "POST", "/habits/"+id, strings.NewReader(body), nil)
}

func (c *Client) UpdateActivity(ctx context.Context, habitID, activityID, desc string) error {
	body := fmt.Sprintf(`{"description":%q}`, desc)
	return c.doRequest(ctx, "PATCH", "/habits/"+habitID+"/"+activityID, strings.NewReader(body), nil)
}

func (c *Client) DeleteActivity(ctx context.Context, habitID, activityID string) error {
	return c.doRequest(ctx, "DELETE", "/habits/"+habitID+"/"+activityID, nil, nil)
}

func (c *Client) CheckIn(ctx context.Context, dto CheckInDTO) (*domain.Habit, error) {
	if err := c.AddActivity(ctx, dto.ID, dto.Desc, dto.Date); err != nil {
		return nil, err
	}
	return c.GetHabit(ctx, dto.ID)
}
