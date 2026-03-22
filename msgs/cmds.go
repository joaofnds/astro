package msgs

import (
	"astro/api"
	"time"

	tea "charm.land/bubbletea/v2"
)

// LoadAll fetches all groups and habits. Returns FatalErrorMsg on failure
// because startup data is required for the application to function.
func LoadAll(client *api.Client) tea.Cmd {
	return func() tea.Msg {
		groups, habits, err := client.GroupsAndHabits()
		if err != nil {
			return FatalErrorMsg{Err: err}
		}
		return DataLoadedMsg{Habits: habits, Groups: groups}
	}
}

// CreateHabit creates a new habit via the API.
func CreateHabit(client *api.Client, name string) tea.Cmd {
	return func() tea.Msg {
		h, err := client.CreateHabit(name)
		if err != nil {
			return APIErrorMsg{Err: err, Op: "create habit"}
		}
		return HabitCreatedMsg{Habit: h}
	}
}

// DeleteHabit deletes a habit by ID via the API.
func DeleteHabit(client *api.Client, id string) tea.Cmd {
	return func() tea.Msg {
		if err := client.DeleteHabit(id); err != nil {
			return APIErrorMsg{Err: err, Op: "delete habit"}
		}
		return HabitDeletedMsg{ID: id}
	}
}

// UpdateHabit updates a habit's name and fetches the updated entity.
// The API's UpdateHabit returns only an error, so we follow up with
// GetHabit to populate HabitUpdatedMsg with the full habit.
func UpdateHabit(client *api.Client, id, name string) tea.Cmd {
	return func() tea.Msg {
		if err := client.UpdateHabit(id, name); err != nil {
			return APIErrorMsg{Err: err, Op: "update habit"}
		}
		h, err := client.GetHabit(id)
		if err != nil {
			return APIErrorMsg{Err: err, Op: "update habit"}
		}
		return HabitUpdatedMsg{Habit: h}
	}
}

// CheckIn performs a check-in (add activity + fetch updated habit).
func CheckIn(client *api.Client, id, desc string, date time.Time) tea.Cmd {
	return func() tea.Msg {
		h, err := client.CheckIn(api.CheckInDTO{ID: id, Desc: desc, Date: date})
		if err != nil {
			return APIErrorMsg{Err: err, Op: "check in"}
		}
		return CheckInResultMsg{Habit: h}
	}
}

// CreateGroup creates a new group via the API.
func CreateGroup(client *api.Client, name string) tea.Cmd {
	return func() tea.Msg {
		g, err := client.CreateGroup(name)
		if err != nil {
			return APIErrorMsg{Err: err, Op: "create group"}
		}
		return GroupCreatedMsg{Group: g}
	}
}

// DeleteGroup deletes a group by ID via the API.
func DeleteGroup(client *api.Client, id string) tea.Cmd {
	return func() tea.Msg {
		if err := client.DeleteGroup(id); err != nil {
			return APIErrorMsg{Err: err, Op: "delete group"}
		}
		return GroupDeletedMsg{ID: id}
	}
}

// AddToGroup adds a habit to a group via the API.
func AddToGroup(client *api.Client, habitID, groupID string) tea.Cmd {
	return func() tea.Msg {
		if err := client.AddToGroup(habitID, groupID); err != nil {
			return APIErrorMsg{Err: err, Op: "add to group"}
		}
		return AddedToGroupMsg{HabitID: habitID, GroupID: groupID}
	}
}

// RemoveFromGroup removes a habit from a group via the API.
func RemoveFromGroup(client *api.Client, habitID, groupID string) tea.Cmd {
	return func() tea.Msg {
		if err := client.RemoveFromGroup(habitID, groupID); err != nil {
			return APIErrorMsg{Err: err, Op: "remove from group"}
		}
		return RemovedFromGroupMsg{HabitID: habitID, GroupID: groupID}
	}
}

// UpdateActivity updates an activity's description via the API.
func UpdateActivity(client *api.Client, habitID, activityID, desc string) tea.Cmd {
	return func() tea.Msg {
		if err := client.UpdateActivity(habitID, activityID, desc); err != nil {
			return APIErrorMsg{Err: err, Op: "update activity"}
		}
		return ActivityUpdatedMsg{HabitID: habitID, ActivityID: activityID, Desc: desc}
	}
}

// DeleteActivity deletes an activity via the API.
func DeleteActivity(client *api.Client, habitID, activityID string) tea.Cmd {
	return func() tea.Msg {
		if err := client.DeleteActivity(habitID, activityID); err != nil {
			return APIErrorMsg{Err: err, Op: "delete activity"}
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
