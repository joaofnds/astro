package app_test

import (
	"astro/app"
	"astro/domain"
	"testing"
)

func habit(id, name string) *domain.Habit {
	return &domain.Habit{ID: id, Name: name}
}

func habitWithActivities(id, name string, activities ...domain.Activity) *domain.Habit {
	return &domain.Habit{ID: id, Name: name, Activities: activities}
}

func group(id, name string, habits ...*domain.Habit) *domain.Group {
	return &domain.Group{ID: id, Name: name, Habits: habits}
}

func activity(id, desc string) domain.Activity {
	return domain.Activity{ID: id, Desc: desc}
}

func TestNewAppState(t *testing.T) {
	s := app.NewAppState()
	if s.Habits() != nil {
		t.Fatalf("expected nil habits, got %v", s.Habits())
	}
	if s.Groups() != nil {
		t.Fatalf("expected nil groups, got %v", s.Groups())
	}
}

func TestSetAll(t *testing.T) {
	s := app.NewAppState()
	habits := []*domain.Habit{habit("h1", "A"), habit("h2", "B")}
	groups := []*domain.Group{group("g1", "G")}

	s.SetAll(habits, groups)

	if got := len(s.Habits()); got != 2 {
		t.Fatalf("expected 2 habits, got %d", got)
	}
	if got := len(s.Groups()); got != 1 {
		t.Fatalf("expected 1 group, got %d", got)
	}
}

func TestHabitByID_TopLevel(t *testing.T) {
	s := app.NewAppState()
	h := habit("h1", "Test")
	s.SetAll([]*domain.Habit{h}, nil)

	got := s.HabitByID("h1")
	if got == nil {
		t.Fatal("expected habit, got nil")
	}
	if got.Name != "Test" {
		t.Fatalf("expected name 'Test', got %q", got.Name)
	}
}

func TestHabitByID_InGroup(t *testing.T) {
	s := app.NewAppState()
	h := habit("h2", "Grouped")
	g := group("g1", "G", h)
	s.SetAll(nil, []*domain.Group{g})

	got := s.HabitByID("h2")
	if got == nil {
		t.Fatal("expected habit from group, got nil")
	}
	if got.Name != "Grouped" {
		t.Fatalf("expected name 'Grouped', got %q", got.Name)
	}
}

func TestHabitByID_NotFound(t *testing.T) {
	s := app.NewAppState()
	s.SetAll([]*domain.Habit{habit("h1", "A")}, nil)

	got := s.HabitByID("nonexistent")
	if got != nil {
		t.Fatalf("expected nil, got %v", got)
	}
}

func TestAddHabit(t *testing.T) {
	s := app.NewAppState()
	s.SetAll(nil, nil)
	s.AddHabit(habit("h1", "New"))

	habits := s.Habits()
	if len(habits) != 1 {
		t.Fatalf("expected 1 habit, got %d", len(habits))
	}
	if habits[0].ID != "h1" {
		t.Fatalf("expected ID 'h1', got %q", habits[0].ID)
	}
}

func TestRemoveHabit(t *testing.T) {
	s := app.NewAppState()
	h := habit("h1", "A")
	s.SetAll([]*domain.Habit{h, habit("h2", "B")}, nil)

	s.RemoveHabit("h1")

	habits := s.Habits()
	if len(habits) != 1 {
		t.Fatalf("expected 1 habit after removal, got %d", len(habits))
	}
	if habits[0].ID != "h2" {
		t.Fatalf("expected remaining habit 'h2', got %q", habits[0].ID)
	}
}

func TestRemoveHabit_AlsoRemovesFromGroups(t *testing.T) {
	s := app.NewAppState()
	h1 := habit("h1", "A")
	h2 := habit("h2", "B")
	g := group("g1", "G", h1, h2)
	s.SetAll([]*domain.Habit{h1, h2}, []*domain.Group{g})

	s.RemoveHabit("h1")

	groups := s.Groups()
	if len(groups[0].Habits) != 1 {
		t.Fatalf("expected 1 habit in group after removal, got %d", len(groups[0].Habits))
	}
	if groups[0].Habits[0].ID != "h2" {
		t.Fatalf("expected remaining group habit 'h2', got %q", groups[0].Habits[0].ID)
	}
}

func TestMergeHabit_TopLevel(t *testing.T) {
	s := app.NewAppState()
	s.SetAll([]*domain.Habit{habit("h1", "Old")}, nil)

	updated := habit("h1", "New")
	s.MergeHabit(updated)

	h := s.HabitByID("h1")
	if h.Name != "New" {
		t.Fatalf("expected name 'New', got %q", h.Name)
	}
}

func TestMergeHabit_InGroups(t *testing.T) {
	s := app.NewAppState()
	h := habit("h1", "Old")
	g := group("g1", "G", h)
	s.SetAll([]*domain.Habit{h}, []*domain.Group{g})

	updated := habit("h1", "New")
	s.MergeHabit(updated)

	// Verify top-level updated
	if s.Habits()[0].Name != "New" {
		t.Fatalf("expected top-level name 'New', got %q", s.Habits()[0].Name)
	}
	// Verify group member updated
	if s.Groups()[0].Habits[0].Name != "New" {
		t.Fatalf("expected group habit name 'New', got %q", s.Groups()[0].Habits[0].Name)
	}
}

func TestAddGroup(t *testing.T) {
	s := app.NewAppState()
	s.SetAll(nil, nil)
	s.AddGroup(group("g1", "New Group"))

	groups := s.Groups()
	if len(groups) != 1 {
		t.Fatalf("expected 1 group, got %d", len(groups))
	}
	if groups[0].Name != "New Group" {
		t.Fatalf("expected name 'New Group', got %q", groups[0].Name)
	}
}

func TestRemoveGroup(t *testing.T) {
	s := app.NewAppState()
	s.SetAll(nil, []*domain.Group{group("g1", "A"), group("g2", "B")})

	s.RemoveGroup("g1")

	groups := s.Groups()
	if len(groups) != 1 {
		t.Fatalf("expected 1 group after removal, got %d", len(groups))
	}
	if groups[0].ID != "g2" {
		t.Fatalf("expected remaining group 'g2', got %q", groups[0].ID)
	}
}

func TestUpdateHabitActivity(t *testing.T) {
	s := app.NewAppState()
	h := habitWithActivities("h1", "H", activity("a1", "old"), activity("a2", "keep"))
	s.SetAll([]*domain.Habit{h}, nil)

	updated := domain.Activity{ID: "a1", Desc: "new"}
	s.UpdateHabitActivity(h, &updated)

	if h.Activities[0].Desc != "new" {
		t.Fatalf("expected desc 'new', got %q", h.Activities[0].Desc)
	}
	if h.Activities[1].Desc != "keep" {
		t.Fatalf("expected desc 'keep' unchanged, got %q", h.Activities[1].Desc)
	}
}

func TestDeleteHabitActivity(t *testing.T) {
	s := app.NewAppState()
	h := habitWithActivities("h1", "H", activity("a1", "remove"), activity("a2", "keep"))
	s.SetAll([]*domain.Habit{h}, nil)

	s.DeleteHabitActivity(h, "a1")

	if len(h.Activities) != 1 {
		t.Fatalf("expected 1 activity after deletion, got %d", len(h.Activities))
	}
	if h.Activities[0].ID != "a2" {
		t.Fatalf("expected remaining activity 'a2', got %q", h.Activities[0].ID)
	}
}

func TestRemoveHabit_NotFound(t *testing.T) {
	s := app.NewAppState()
	s.SetAll([]*domain.Habit{habit("h1", "A")}, nil)

	// Should not panic when removing a nonexistent habit.
	s.RemoveHabit("nonexistent")

	if len(s.Habits()) != 1 {
		t.Fatal("habits should be unchanged")
	}
}

func TestRemoveGroup_NotFound(t *testing.T) {
	s := app.NewAppState()
	s.SetAll(nil, []*domain.Group{group("g1", "A")})

	// Should not panic when removing a nonexistent group.
	s.RemoveGroup("nonexistent")

	if len(s.Groups()) != 1 {
		t.Fatal("groups should be unchanged")
	}
}

func TestMergeHabit_NotFound(t *testing.T) {
	s := app.NewAppState()
	s.SetAll([]*domain.Habit{habit("h1", "A")}, nil)

	// Merging a habit that doesn't exist should be a no-op.
	s.MergeHabit(habit("nonexistent", "X"))

	if s.Habits()[0].Name != "A" {
		t.Fatal("existing habit should be unchanged")
	}
}
