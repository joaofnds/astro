package msgs

import (
	"astro/domain"

	tea "charm.land/bubbletea/v2"
)

// --- Navigation Messages ---

// PushScreenMsg tells the root model to push a new screen onto the stack.
type PushScreenMsg struct {
	Screen tea.Model
}

// PopScreenMsg tells the root model to pop the current screen.
// Cmd is an optional command to run after popping (e.g., a data refresh).
type PopScreenMsg struct {
	Cmd tea.Cmd
}

// PushScreen returns a tea.Cmd that produces a PushScreenMsg.
func PushScreen(screen tea.Model) tea.Cmd {
	return func() tea.Msg {
		return PushScreenMsg{Screen: screen}
	}
}

// PopScreen returns a tea.Cmd that produces a PopScreenMsg with no follow-up.
func PopScreen() tea.Cmd {
	return func() tea.Msg {
		return PopScreenMsg{}
	}
}

// PopScreenWith returns a tea.Cmd that pops and runs a follow-up command.
func PopScreenWith(cmd tea.Cmd) tea.Cmd {
	return func() tea.Msg {
		return PopScreenMsg{Cmd: cmd}
	}
}

// --- Async Result Messages ---

// DataLoadedMsg carries the initial data load result.
type DataLoadedMsg struct {
	Habits []*domain.Habit
	Groups []*domain.Group
}

// HabitCreatedMsg carries a newly created habit from the API.
type HabitCreatedMsg struct {
	Habit *domain.Habit
}

// HabitDeletedMsg carries the ID of a deleted habit.
type HabitDeletedMsg struct {
	ID string
}

// HabitUpdatedMsg carries an updated habit fetched from the API.
type HabitUpdatedMsg struct {
	Habit *domain.Habit
}

// CheckInResultMsg carries the habit state after a successful check-in.
type CheckInResultMsg struct {
	Habit *domain.Habit
}

// GroupCreatedMsg carries a newly created group from the API.
type GroupCreatedMsg struct {
	Group *domain.Group
}

// GroupDeletedMsg carries the ID of a deleted group.
type GroupDeletedMsg struct {
	ID string
}

// AddedToGroupMsg confirms a habit was added to a group.
type AddedToGroupMsg struct {
	HabitID string
	GroupID string
}

// RemovedFromGroupMsg confirms a habit was removed from a group.
type RemovedFromGroupMsg struct {
	HabitID string
	GroupID string
}

// ActivityUpdatedMsg confirms an activity description was updated.
type ActivityUpdatedMsg struct {
	HabitID    string
	ActivityID string
	Desc       string
}

// ActivityDeletedMsg confirms an activity was deleted.
type ActivityDeletedMsg struct {
	HabitID    string
	ActivityID string
}

// --- Error Messages ---

// APIErrorMsg is a recoverable error shown as a status message.
// The Op field describes the operation that failed (e.g., "create habit").
type APIErrorMsg struct {
	Err error
	Op  string
}

// FatalErrorMsg causes the application to quit with an error message.
// Used for unrecoverable failures like startup data load errors.
type FatalErrorMsg struct {
	Err error
}

// ClearStatusMsg tells a screen to clear its status/error message.
type ClearStatusMsg struct{}
