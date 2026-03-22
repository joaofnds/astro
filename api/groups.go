package api

import (
	"astro/domain"
	"fmt"
	"strings"
)

// GroupsAndHabitsPayload is the response shape for the /groups endpoint.
type GroupsAndHabitsPayload struct {
	Groups []*domain.Group `json:"groups"`
	Habits []*domain.Habit `json:"habits"`
}

func (c *Client) GroupsAndHabits() ([]*domain.Group, []*domain.Habit, error) {
	var data GroupsAndHabitsPayload
	if err := c.doRequest("GET", "/groups", nil, &data); err != nil {
		return nil, nil, err
	}
	domain.SortHabits(data.Habits)
	domain.SortGroups(data.Groups)
	return data.Groups, data.Habits, nil
}

func (c *Client) CreateGroup(name string) (*domain.Group, error) {
	body := fmt.Sprintf(`{"name":%q}`, name)
	var group domain.Group
	if err := c.doRequest("POST", "/groups", strings.NewReader(body), &group); err != nil {
		return nil, err
	}
	return &group, nil
}

func (c *Client) AddToGroup(habitID, groupID string) error {
	return c.doRequest("POST", "/groups/"+groupID+"/"+habitID, nil, nil)
}

func (c *Client) RemoveFromGroup(habitID, groupID string) error {
	return c.doRequest("DELETE", "/groups/"+groupID+"/"+habitID, nil, nil)
}

func (c *Client) DeleteGroup(groupID string) error {
	return c.doRequest("DELETE", "/groups/"+groupID, nil, nil)
}
