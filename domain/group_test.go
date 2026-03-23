package domain

import (
	"astro/config"
	"astro/date"
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestSortGroups(t *testing.T) {
	t.Run("sorts groups by name", func(t *testing.T) {
		groups := []*Group{
			{Name: "zzz"},
			{Name: "aaa"},
			{Name: "mmm"},
		}
		SortGroups(groups)
		for i := 0; i < len(groups)-1; i++ {
			if groups[i].Name >= groups[i+1].Name {
				t.Errorf("groups not sorted: %q >= %q at index %d", groups[i].Name, groups[i+1].Name, i)
			}
		}
	})

	t.Run("sorts habits within each group", func(t *testing.T) {
		groups := []*Group{
			{
				Name: "fitness",
				Habits: []*Habit{
					{Name: "yoga"},
					{Name: "abs"},
					{Name: "run"},
				},
			},
		}
		SortGroups(groups)
		for i := 0; i < len(groups[0].Habits)-1; i++ {
			if groups[0].Habits[i].Name >= groups[0].Habits[i+1].Name {
				t.Errorf("habits not sorted: %q >= %q at index %d",
					groups[0].Habits[i].Name, groups[0].Habits[i+1].Name, i)
			}
		}
	})

	t.Run("empty slice is no-op", func(t *testing.T) {
		SortGroups(nil)
		// no panic = pass
	})
}

func TestGroupActivities(t *testing.T) {
	base := time.Date(2025, 3, 15, 10, 0, 0, 0, time.Local)

	t.Run("empty group returns empty slice", func(t *testing.T) {
		g := &Group{}
		got := g.Activities()
		if len(got) != 0 {
			t.Errorf("Activities() length = %d, want 0", len(got))
		}
	})

	t.Run("single habit returns its activities", func(t *testing.T) {
		g := &Group{
			Habits: []*Habit{
				{
					Activities: []Activity{
						{ID: "a1", CreatedAt: base},
						{ID: "a2", CreatedAt: base.Add(1 * time.Hour)},
					},
				},
			},
		}
		got := g.Activities()
		if len(got) != 2 {
			t.Fatalf("Activities() length = %d, want 2", len(got))
		}
		if got[0].ID != "a1" || got[1].ID != "a2" {
			t.Errorf("Activities() = [%s, %s], want [a1, a2]", got[0].ID, got[1].ID)
		}
	})

	t.Run("multiple habits returns sorted by CreatedAt", func(t *testing.T) {
		g := &Group{
			Habits: []*Habit{
				{
					Activities: []Activity{
						{ID: "h1-a1", CreatedAt: base.Add(2 * time.Hour)},
					},
				},
				{
					Activities: []Activity{
						{ID: "h2-a1", CreatedAt: base},
						{ID: "h2-a2", CreatedAt: base.Add(3 * time.Hour)},
					},
				},
			},
		}
		got := g.Activities()
		if len(got) != 3 {
			t.Fatalf("Activities() length = %d, want 3", len(got))
		}
		for i := 0; i < len(got)-1; i++ {
			if !got[i].CreatedAt.Before(got[i+1].CreatedAt) {
				t.Errorf("activities not sorted: %v >= %v at index %d",
					got[i].CreatedAt, got[i+1].CreatedAt, i)
			}
		}
	})
}

func TestActivitiesOnDateTally(t *testing.T) {
	ref := time.Date(2025, 3, 15, 0, 0, 0, 0, time.Local)
	formatted := ref.Format(config.DateFormat)

	t.Run("0 activities", func(t *testing.T) {
		habits := []*Habit{
			{Name: "run", Activities: nil},
		}
		got := ActivitiesOnDateTally(habits, ref)
		want := fmt.Sprintf("0 activities on %s\n", formatted)
		if got != want {
			t.Errorf("ActivitiesOnDateTally() = %q, want %q", got, want)
		}
	})

	t.Run("1 activity singular with tally", func(t *testing.T) {
		habits := []*Habit{
			{
				Name: "run",
				Activities: []Activity{
					{CreatedAt: ref.Add(10 * time.Hour)},
				},
			},
		}
		got := ActivitiesOnDateTally(habits, ref)
		want := fmt.Sprintf("1 activity on %s (run:1)\n", formatted)
		if got != want {
			t.Errorf("ActivitiesOnDateTally() = %q, want %q", got, want)
		}
	})

	t.Run("multiple habits with tally", func(t *testing.T) {
		habits := []*Habit{
			{
				Name: "run",
				Activities: []Activity{
					{CreatedAt: ref.Add(8 * time.Hour)},
					{CreatedAt: ref.Add(18 * time.Hour)},
				},
			},
			{
				Name: "read",
				Activities: []Activity{
					{CreatedAt: ref.Add(12 * time.Hour)},
				},
			},
		}
		got := ActivitiesOnDateTally(habits, ref)

		// Map iteration order is non-deterministic, so check components
		if !strings.HasPrefix(got, fmt.Sprintf("3 activities on %s (", formatted)) {
			t.Errorf("ActivitiesOnDateTally() prefix = %q, want %q prefix",
				got, fmt.Sprintf("3 activities on %s (", formatted))
		}
		if !strings.Contains(got, "run:2") {
			t.Errorf("ActivitiesOnDateTally() missing run:2 tally in %q", got)
		}
		if !strings.Contains(got, "read:1") {
			t.Errorf("ActivitiesOnDateTally() missing read:1 tally in %q", got)
		}
	})

	t.Run("activities on different dates are excluded", func(t *testing.T) {
		habits := []*Habit{
			{
				Name: "run",
				Activities: []Activity{
					{CreatedAt: ref.AddDate(0, 0, -1)}, // different date
					{CreatedAt: ref.Add(10 * time.Hour)},
				},
			},
		}
		got := ActivitiesOnDateTally(habits, ref)
		want := fmt.Sprintf("1 activity on %s (run:1)\n", formatted)
		if got != want {
			t.Errorf("ActivitiesOnDateTally() = %q, want %q", got, want)
		}
	})
}

// TestActivitiesOnDateTallyUsesToday verifies the function works with the
// current date, matching the usage pattern in the real application where
// date.Today() is passed.
func TestActivitiesOnDateTallyUsesToday(t *testing.T) {
	today := date.Today()
	habits := []*Habit{
		{Name: "meditate"},
	}
	got := ActivitiesOnDateTally(habits, today)
	want := fmt.Sprintf("0 activities on %s\n", today.Format(config.DateFormat))
	if got != want {
		t.Errorf("ActivitiesOnDateTally(today) = %q, want %q", got, want)
	}
}
