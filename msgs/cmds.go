package msgs

import (
	"astro/api"
	"context"
	"time"

	tea "charm.land/bubbletea/v2"
)

// LoadAll fetches all groups and habits. Returns FatalErrorMsg on failure
// because startup data is required for the application to function.
// If the context is cancelled, the result is silently discarded.
func LoadAll(ctx context.Context, client *api.Client) tea.Cmd {
	return func() tea.Msg {
		groups, habits, err := client.GroupsAndHabits(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			return FatalErrorMsg{Err: err}
		}
		return DataLoadedMsg{Habits: habits, Groups: groups}
	}
}

// CreateHabit creates a new habit via the API.
func CreateHabit(ctx context.Context, client *api.Client, name string) tea.Cmd {
	return func() tea.Msg {
		h, err := client.CreateHabit(ctx, name)
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			return APIErrorMsg{Err: err, Op: "create habit"}
		}
		return HabitCreatedMsg{Habit: h}
	}
}

// DeleteHabit deletes a habit by ID via the API.
func DeleteHabit(ctx context.Context, client *api.Client, id string) tea.Cmd {
	return func() tea.Msg {
		if err := client.DeleteHabit(ctx, id); err != nil {
			if ctx.Err() != nil {
				return nil
			}
			return APIErrorMsg{Err: err, Op: "delete habit", ID: id}
		}
		return HabitDeletedMsg{ID: id}
	}
}

// UpdateHabit updates a habit's name and fetches the updated entity.
// The API's UpdateHabit returns only an error, so we follow up with
// GetHabit to populate HabitUpdatedMsg with the full habit.
func UpdateHabit(ctx context.Context, client *api.Client, id, name string) tea.Cmd {
	return func() tea.Msg {
		if err := client.UpdateHabit(ctx, id, name); err != nil {
			if ctx.Err() != nil {
				return nil
			}
			return APIErrorMsg{Err: err, Op: "update habit", ID: id}
		}
		h, err := client.GetHabit(ctx, id)
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			return APIErrorMsg{Err: err, Op: "update habit", ID: id}
		}
		return HabitUpdatedMsg{Habit: h}
	}
}

// CheckIn performs a check-in (add activity + fetch updated habit).
func CheckIn(ctx context.Context, client *api.Client, id, desc string, date time.Time) tea.Cmd {
	return func() tea.Msg {
		h, err := client.CheckIn(ctx, api.CheckInDTO{ID: id, Desc: desc, Date: date})
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			return APIErrorMsg{Err: err, Op: "check in", ID: id}
		}
		return CheckInResultMsg{Habit: h}
	}
}

// CreateGroup creates a new group via the API.
func CreateGroup(ctx context.Context, client *api.Client, name string) tea.Cmd {
	return func() tea.Msg {
		g, err := client.CreateGroup(ctx, name)
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			return APIErrorMsg{Err: err, Op: "create group"}
		}
		return GroupCreatedMsg{Group: g}
	}
}

// DeleteGroup deletes a group by ID via the API.
func DeleteGroup(ctx context.Context, client *api.Client, id string) tea.Cmd {
	return func() tea.Msg {
		if err := client.DeleteGroup(ctx, id); err != nil {
			if ctx.Err() != nil {
				return nil
			}
			return APIErrorMsg{Err: err, Op: "delete group", ID: id}
		}
		return GroupDeletedMsg{ID: id}
	}
}

// AddToGroup adds a habit to a group via the API.
func AddToGroup(ctx context.Context, client *api.Client, habitID, groupID string) tea.Cmd {
	return func() tea.Msg {
		if err := client.AddToGroup(ctx, habitID, groupID); err != nil {
			if ctx.Err() != nil {
				return nil
			}
			return APIErrorMsg{Err: err, Op: "add to group", ID: habitID}
		}
		return AddedToGroupMsg{HabitID: habitID, GroupID: groupID}
	}
}

// RemoveFromGroup removes a habit from a group via the API.
func RemoveFromGroup(ctx context.Context, client *api.Client, habitID, groupID string) tea.Cmd {
	return func() tea.Msg {
		if err := client.RemoveFromGroup(ctx, habitID, groupID); err != nil {
			if ctx.Err() != nil {
				return nil
			}
			return APIErrorMsg{Err: err, Op: "remove from group", ID: habitID}
		}
		return RemovedFromGroupMsg{HabitID: habitID, GroupID: groupID}
	}
}

// UpdateActivity updates an activity's description via the API.
func UpdateActivity(ctx context.Context, client *api.Client, habitID, activityID, desc string) tea.Cmd {
	return func() tea.Msg {
		if err := client.UpdateActivity(ctx, habitID, activityID, desc); err != nil {
			if ctx.Err() != nil {
				return nil
			}
			return APIErrorMsg{Err: err, Op: "update activity", ID: activityID}
		}
		return ActivityUpdatedMsg{HabitID: habitID, ActivityID: activityID, Desc: desc}
	}
}

// DeleteActivity deletes an activity via the API.
func DeleteActivity(ctx context.Context, client *api.Client, habitID, activityID string) tea.Cmd {
	return func() tea.Msg {
		if err := client.DeleteActivity(ctx, habitID, activityID); err != nil {
			if ctx.Err() != nil {
				return nil
			}
			return APIErrorMsg{Err: err, Op: "delete activity", ID: activityID}
		}
		return ActivityDeletedMsg{HabitID: habitID, ActivityID: activityID}
	}
}

// ClearStatusAfter returns a tea.Cmd that produces ClearStatusMsg after
// the given duration. Uses tea.Tick for integration with the Bubbletea
// event loop.
func ClearStatusAfter(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(time.Time) tea.Msg {
		return ClearStatusMsg{}
	})
}
