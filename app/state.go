package app

import (
	"astro/domain"
)

// AppState holds all application data. The root model is the sole owner.
// All fields are unexported; access is via methods only.
// No *api.Client field -- command factories in msgs/cmds.go take the
// client as a parameter, keeping AppState purely about data.
type AppState struct {
	habits []*domain.Habit
	groups []*domain.Group
}

// NewAppState creates an empty AppState.
func NewAppState() AppState {
	return AppState{}
}

// --- Accessor methods (value receiver) ---

// Habits returns the current habit list.
func (s AppState) Habits() []*domain.Habit {
	return s.habits
}

// Groups returns the current group list.
func (s AppState) Groups() []*domain.Group {
	return s.groups
}

// HabitByID returns the habit with the given ID, searching both the
// top-level list and within groups. Returns nil if not found.
func (s AppState) HabitByID(id string) *domain.Habit {
	for _, h := range s.habits {
		if h.ID == id {
			return h
		}
	}
	for _, g := range s.groups {
		for _, h := range g.Habits {
			if h.ID == id {
				return h
			}
		}
	}
	return nil
}

// --- Mutation methods (pointer receiver) ---

// SetAll replaces all habits and groups. Used after initial data load
// or a full refresh.
func (s *AppState) SetAll(habits []*domain.Habit, groups []*domain.Group) {
	s.habits = habits
	s.groups = groups
}

// AddHabit appends a habit to the top-level list.
func (s *AppState) AddHabit(h *domain.Habit) {
	s.habits = append(s.habits, h)
}

// RemoveHabit removes a habit by ID from the top-level list and from
// every group's Habits slice.
func (s *AppState) RemoveHabit(id string) {
	s.habits = removeHabitFromSlice(s.habits, id)
	for _, g := range s.groups {
		g.Habits = removeHabitFromSlice(g.Habits, id)
	}
}

// MergeHabit replaces an existing habit (matched by ID) in the
// top-level list and within every group. If the ID is not found,
// this is a no-op.
func (s *AppState) MergeHabit(updated *domain.Habit) {
	for i, h := range s.habits {
		if h.ID == updated.ID {
			s.habits[i] = updated
			break
		}
	}
	for _, g := range s.groups {
		for i, h := range g.Habits {
			if h.ID == updated.ID {
				g.Habits[i] = updated
				break
			}
		}
	}
}

// AddGroup appends a group to the group list.
func (s *AppState) AddGroup(g *domain.Group) {
	s.groups = append(s.groups, g)
}

// RemoveGroup removes a group by ID.
func (s *AppState) RemoveGroup(id string) {
	for i, g := range s.groups {
		if g.ID == id {
			s.groups = append(s.groups[:i], s.groups[i+1:]...)
			return
		}
	}
}

// UpdateHabitActivity finds the activity by ID within the habit and
// updates its description.
func (s *AppState) UpdateHabitActivity(habit *domain.Habit, activity *domain.Activity) {
	for i := range habit.Activities {
		if habit.Activities[i].ID == activity.ID {
			habit.Activities[i].Desc = activity.Desc
			return
		}
	}
}

// DeleteHabitActivity removes the activity with the given ID from the
// habit's Activities slice.
func (s *AppState) DeleteHabitActivity(habit *domain.Habit, activityID string) {
	for i := range habit.Activities {
		if habit.Activities[i].ID == activityID {
			habit.Activities = append(habit.Activities[:i], habit.Activities[i+1:]...)
			return
		}
	}
}

// removeHabitFromSlice returns a new slice with the habit matching id removed.
func removeHabitFromSlice(habits []*domain.Habit, id string) []*domain.Habit {
	for i, h := range habits {
		if h.ID == id {
			return append(habits[:i], habits[i+1:]...)
		}
	}
	return habits
}
