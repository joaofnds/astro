package date_test

import (
	"astro/date"
	"testing"
	"time"
)

func TestDiffInDays(t *testing.T) {
	tests := []struct {
		name     string
		t1       time.Time
		t2       time.Time
		expected int
	}{
		{
			name:     "Same day",
			t1:       time.Date(2025, 1, 20, 0, 0, 0, 0, time.UTC),
			t2:       time.Date(2025, 1, 20, 23, 59, 59, 0, time.UTC),
			expected: 0,
		},
		{
			name:     "Next day, small hour difference",
			t1:       time.Date(2025, 1, 20, 23, 0, 0, 0, time.UTC),
			t2:       time.Date(2025, 1, 21, 1, 0, 0, 0, time.UTC),
			expected: 1,
		},
		{
			name:     "Multiple days difference",
			t1:       time.Date(2025, 1, 20, 0, 0, 0, 0, time.UTC),
			t2:       time.Date(2025, 1, 25, 0, 0, 0, 0, time.UTC),
			expected: 5,
		},
		{
			name:     "Previous day",
			t1:       time.Date(2025, 1, 21, 0, 0, 0, 0, time.UTC),
			t2:       time.Date(2025, 1, 20, 0, 0, 0, 0, time.UTC),
			expected: 1,
		},
		{
			name:     "Different months",
			t1:       time.Date(2025, 1, 31, 0, 0, 0, 0, time.UTC),
			t2:       time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC),
			expected: 1,
		},
		{
			name:     "Different years",
			t1:       time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC),
			t2:       time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: 1,
		},
		{
			name:     "365 days (1 year) difference",
			t1:       time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			t2:       time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: 365,
		},
		{
			name:     "Leap year crossing, 366 days",
			t1:       time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC),
			t2:       time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
			expected: 366, // 2024 is a leap year
		},
		{
			name:     "Negative order, 200 days difference",
			t1:       time.Date(2025, 1, 17, 0, 0, 0, 0, time.UTC),
			t2:       time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC),
			expected: 200,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := date.DiffInDays(test.t1, test.t2)
			if result != test.expected {
				t.Errorf("DiffInDays(%v, %v) = %d; expected %d", test.t1, test.t2, result, test.expected)
			}
		})
	}
}

func TestIncremental(t *testing.T) {
	now := time.Now()
	t0 := time.Date(2025, 1, 1, now.Hour(), now.Minute(), now.Second(), now.Nanosecond(), now.Location())
	for i := 0; i <= 500; i++ {
		t1 := t0.Add(time.Duration(i*24) * time.Hour)
		if date.DiffInDays(t0, t1) != i {
			t.Errorf("DiffInDays(%v, %v) != %d", t0, t1, i)
		}
	}
}
