package listitem_test

import (
	"astro/domain"
	"astro/models/listitem"
	"testing"
	"time"
)

func TestGroupItemReflectsLiveGroupState(t *testing.T) {
	habit := &domain.Habit{ID: "h1", Name: "read"}
	group := &domain.Group{ID: "g1", Name: "morning", Habits: []*domain.Habit{habit}}
	item := listitem.GroupItem{Group: group}

	before := item.Description()

	habit.Activities = append(habit.Activities, domain.Activity{ID: "a1", CreatedAt: time.Now()})
	after := item.Description()

	if before == after {
		t.Fatalf("description did not change after adding an activity; got %q both times", before)
	}
}
