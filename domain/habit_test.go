package domain

import (
	"astro/config"
	"astro/date"
	"fmt"
	"testing"
	"time"
)

func TestMomentum(t *testing.T) {
	activities := []Activity{
		{CreatedAt: date.Today().AddDate(0, 0, -9)}, // +1
		//                                     skipped: -1 = 0
		{CreatedAt: date.Today().AddDate(0, 0, -7)}, // +1 = 1
		//                                     skipped: -1 = 0
		//                                     skipped: -1 = 0  <- min is 0
		{CreatedAt: date.Today().AddDate(0, 0, -4)}, // +1 = 1
		{CreatedAt: date.Today().AddDate(0, 0, -4)}, // +1 = 2
		{CreatedAt: date.Today().AddDate(0, 0, -3)}, // +1 = 3
		//                                     skipped: -1 = 2
		{CreatedAt: date.Today().AddDate(0, 0, -1)}, // +1 = 3
	}

	expected := 3
	got := Momentum(activities)

	if got != expected {
		t.Errorf("expected %d, got %d", expected, got)
	}
}

func TestMomentumWithNoActivities(t *testing.T) {
	expected := 0
	got := Momentum([]Activity{})

	if got != expected {
		t.Errorf("expected %d, got %d", expected, got)
	}
}

func TestDigest(t *testing.T) {
	activities := []Activity{
		{CreatedAt: date.Today().AddDate(0, 0, -5)},
		{CreatedAt: date.Today().AddDate(0, 0, -4)},
		{CreatedAt: date.Today().AddDate(0, 0, -2)},
		{CreatedAt: date.Today().AddDate(0, 0, -1)},
	}
	got := Digest("run", activities)
	want := "run - streak: 2 days, momentum: 3"

	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestStreak(t *testing.T) {
	tt := []struct {
		name       string
		activities []Activity
		want       string
	}{
		{
			name:       "0 activities returns 0 days",
			activities: nil,
			want:       "0 days",
		},
		{
			name: "1-day streak returns singular",
			activities: []Activity{
				{CreatedAt: date.Today()},
			},
			want: "1 day",
		},
		{
			name: "multi-day streak returns plural",
			activities: []Activity{
				{CreatedAt: date.Today().AddDate(0, 0, -2)},
				{CreatedAt: date.Today().AddDate(0, 0, -1)},
				{CreatedAt: date.Today()},
			},
			want: "3 days",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			got := Streak(tc.activities)
			if got != tc.want {
				t.Errorf("Streak() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestCurrentStreakDays(t *testing.T) {
	tt := []struct {
		name       string
		activities []Activity
		want       int
	}{
		{
			name:       "0 activities returns 0",
			activities: nil,
			want:       0,
		},
		{
			name: "activity today returns 1",
			activities: []Activity{
				{CreatedAt: date.Today()},
			},
			want: 1,
		},
		{
			name: "activity yesterday not today returns 1",
			activities: []Activity{
				{CreatedAt: date.Today().AddDate(0, 0, -1)},
			},
			want: 1,
		},
		{
			name: "consecutive days returns correct count",
			activities: []Activity{
				{CreatedAt: date.Today().AddDate(0, 0, -3)},
				{CreatedAt: date.Today().AddDate(0, 0, -2)},
				{CreatedAt: date.Today().AddDate(0, 0, -1)},
				{CreatedAt: date.Today()},
			},
			want: 4,
		},
		{
			name: "gap breaks streak",
			activities: []Activity{
				{CreatedAt: date.Today().AddDate(0, 0, -5)},
				{CreatedAt: date.Today().AddDate(0, 0, -4)},
				// gap on -3
				{CreatedAt: date.Today().AddDate(0, 0, -1)},
				{CreatedAt: date.Today()},
			},
			want: 2,
		},
		{
			name: "multiple activities on same day count as 1",
			activities: []Activity{
				{CreatedAt: date.Today().AddDate(0, 0, -1)},
				{CreatedAt: date.Today()},
				{CreatedAt: date.Today().Add(1 * time.Hour)},
				{CreatedAt: date.Today().Add(2 * time.Hour)},
			},
			want: 2,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			got := CurrentStreakDays(tc.activities)
			if got != tc.want {
				t.Errorf("CurrentStreakDays() = %d, want %d", got, tc.want)
			}
		})
	}
}

func TestLatestActivity(t *testing.T) {
	t.Run("empty activities returns zero time", func(t *testing.T) {
		h := Habit{}
		got := h.LatestActivity()
		if !got.IsZero() {
			t.Errorf("LatestActivity() = %v, want zero time", got)
		}
	})

	t.Run("returns last activity CreatedAt", func(t *testing.T) {
		ts := time.Date(2025, 3, 15, 10, 0, 0, 0, time.Local)
		h := Habit{
			Activities: []Activity{
				{CreatedAt: ts.AddDate(0, 0, -2)},
				{CreatedAt: ts.AddDate(0, 0, -1)},
				{CreatedAt: ts},
			},
		}
		got := h.LatestActivity()
		if !got.Equal(ts) {
			t.Errorf("LatestActivity() = %v, want %v", got, ts)
		}
	})
}

func TestLatestActivityOnDate(t *testing.T) {
	ref := time.Date(2025, 3, 15, 0, 0, 0, 0, time.Local)

	t.Run("no activity on date returns error", func(t *testing.T) {
		h := Habit{
			Activities: []Activity{
				{ID: "a1", CreatedAt: ref.AddDate(0, 0, -1)},
			},
		}
		_, err := h.LatestActivityOnDate(ref)
		if err == nil {
			t.Error("expected error for missing date, got nil")
		}
	})

	t.Run("returns latest activity matching date", func(t *testing.T) {
		h := Habit{
			Activities: []Activity{
				{ID: "a1", CreatedAt: ref.Add(8 * time.Hour)},
				{ID: "a2", CreatedAt: ref.Add(10 * time.Hour)},
				{ID: "a3", CreatedAt: ref.Add(14 * time.Hour)},
			},
		}
		got, err := h.LatestActivityOnDate(ref)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.ID != "a3" {
			t.Errorf("LatestActivityOnDate() ID = %q, want %q", got.ID, "a3")
		}
	})

	t.Run("multiple activities returns last via reverse iteration", func(t *testing.T) {
		h := Habit{
			Activities: []Activity{
				{ID: "old", CreatedAt: ref.AddDate(0, 0, -1)},
				{ID: "first", CreatedAt: ref.Add(1 * time.Hour)},
				{ID: "second", CreatedAt: ref.Add(5 * time.Hour)},
			},
		}
		got, err := h.LatestActivityOnDate(ref)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.ID != "second" {
			t.Errorf("LatestActivityOnDate() ID = %q, want %q", got.ID, "second")
		}
	})
}

func TestSortHabits(t *testing.T) {
	t.Run("sorts by name alphabetically", func(t *testing.T) {
		habits := []*Habit{
			{Name: "yoga"},
			{Name: "abc"},
			{Name: "meditation"},
		}
		SortHabits(habits)
		for i := 0; i < len(habits)-1; i++ {
			if habits[i].Name >= habits[i+1].Name {
				t.Errorf("habits not sorted: %q >= %q at index %d", habits[i].Name, habits[i+1].Name, i)
			}
		}
	})

	t.Run("also sorts activities within each habit", func(t *testing.T) {
		later := time.Date(2025, 3, 15, 12, 0, 0, 0, time.Local)
		earlier := time.Date(2025, 3, 15, 8, 0, 0, 0, time.Local)
		habits := []*Habit{
			{
				Name: "run",
				Activities: []Activity{
					{CreatedAt: later},
					{CreatedAt: earlier},
				},
			},
		}
		SortHabits(habits)
		if !habits[0].Activities[0].CreatedAt.Before(habits[0].Activities[1].CreatedAt) {
			t.Error("activities within habit not sorted by CreatedAt")
		}
	})
}

func TestSortActivities(t *testing.T) {
	t.Run("sorts by CreatedAt ascending", func(t *testing.T) {
		base := time.Date(2025, 3, 15, 10, 0, 0, 0, time.Local)
		activities := []Activity{
			{ID: "c", CreatedAt: base.Add(2 * time.Hour)},
			{ID: "a", CreatedAt: base},
			{ID: "b", CreatedAt: base.Add(1 * time.Hour)},
		}
		SortActivities(activities)
		for i := 0; i < len(activities)-1; i++ {
			if !activities[i].CreatedAt.Before(activities[i+1].CreatedAt) {
				t.Errorf("activities not sorted: %v >= %v at index %d",
					activities[i].CreatedAt, activities[i+1].CreatedAt, i)
			}
		}
	})

	t.Run("empty slice is no-op", func(t *testing.T) {
		SortActivities(nil)
		// no panic = pass
	})
}

func TestActivitiesOnDate(t *testing.T) {
	ref := time.Date(2025, 3, 15, 0, 0, 0, 0, time.Local)
	formatted := ref.Format(config.DateFormat)

	tt := []struct {
		name       string
		activities []Activity
		want       string
	}{
		{
			name:       "0 activities",
			activities: nil,
			want:       fmt.Sprintf("0 activities on %s\n", formatted),
		},
		{
			name: "1 activity singular",
			activities: []Activity{
				{CreatedAt: ref.Add(10 * time.Hour)},
			},
			want: fmt.Sprintf("1 activity on %s\n", formatted),
		},
		{
			name: "multiple activities",
			activities: []Activity{
				{CreatedAt: ref.Add(8 * time.Hour)},
				{CreatedAt: ref.Add(12 * time.Hour)},
				{CreatedAt: ref.Add(18 * time.Hour)},
			},
			want: fmt.Sprintf("3 activities on %s\n", formatted),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			got := ActivitiesOnDate(tc.activities, ref)
			if got != tc.want {
				t.Errorf("ActivitiesOnDate() = %q, want %q", got, tc.want)
			}
		})
	}
}
