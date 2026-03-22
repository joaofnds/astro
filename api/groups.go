package api

import (
	"astro/domain"
	"context"
	"fmt"
	"strings"
)

// GroupsAndHabitsPayload is the response shape for the /groups endpoint.
type GroupsAndHabitsPayload struct {
	Groups []*domain.Group `json:"groups"`
	Habits []*domain.Habit `json:"habits"`
}

func (c *Client) GroupsAndHabits(ctx context.Context) ([]*domain.Group, []*domain.Habit, error) {
	var data GroupsAndHabitsPayload
	if err := c.doRequest(ctx, "GET", "/groups", nil, &data); err != nil {
		return nil, nil, err
	}
	domain.SortHabits(data.Habits)
	domain.SortGroups(data.Groups)
	return data.Groups, data.Habits, nil
}

func (c *Client) CreateGroup(ctx context.Context, name string) (*domain.Group, error) {
	body := fmt.Sprintf(`{"name":%q}`, name)
	var group domain.Group
	if err := c.doRequest(ctx, "POST", "/groups", strings.NewReader(body), &group); err != nil {
		return nil, err
	}
	return &group, nil
}

func (c *Client) AddToGroup(ctx context.Context, habitID, groupID string) error {
	return c.doRequest(ctx, "POST", "/groups/"+groupID+"/"+habitID, nil, nil)
}

func (c *Client) RemoveFromGroup(ctx context.Context, habitID, groupID string) error {
	return c.doRequest(ctx, "DELETE", "/groups/"+groupID+"/"+habitID, nil, nil)
}

func (c *Client) DeleteGroup(ctx context.Context, groupID string) error {
	return c.doRequest(ctx, "DELETE", "/groups/"+groupID, nil, nil)
}
