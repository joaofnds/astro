package domain

import (
	"fmt"
	"testing"
	"time"
)

// withLocalTZ pins time.Local to zone for the duration of the test so that
// assertions about local-day bucketing are stable regardless of the machine's
// TZ. It does NOT run in parallel with other tests that read time.Local.
func withLocalTZ(t *testing.T, zone *time.Location) {
	t.Helper()
	orig := time.Local
	time.Local = zone
	t.Cleanup(func() { time.Local = orig })
}

func TestHistBucketsLocalDayBoundary(t *testing.T) {
	// São Paulo is UTC-3. Most local hours produce a UTC CreatedAt whose
	// UTC calendar day equals the local calendar day, and the naive
	// DiffInDays(localMidnight, utcInstant) mis-buckets them by one day.
	// Each case places a single activity at a known local hour/day and
	// asserts the bucket index it lands in.
	withLocalTZ(t, time.FixedZone("SP", -3*3600))

	// Use a fixed start rather than date.Today() so the test is stable.
	start := time.Date(2025, 1, 16, 0, 0, 0, 0, time.Local)

	// Activities are stored UTC by the backend; build them by taking a
	// local wall-clock time and converting to UTC.
	at := func(localDay, localHour int) Activity {
		return Activity{
			CreatedAt: time.Date(2025, 1, localDay, localHour, 0, 0, 0, time.Local).UTC(),
		}
	}

	cases := []struct {
		name    string
		act     Activity
		wantBin int // index into a 7-slot histogram starting at Jan 16
	}{
		// All these used to land one bucket too early before the .Local() fix.
		{"today 01:00 local", at(22, 1), 6},
		{"today 14:00 local", at(22, 14), 6},
		// Evening activity: accidentally correct in the broken version too,
		// included so the test fully documents the behavior.
		{"today 22:30 local", at(22, 22), 6},
		{"yesterday 01:00 local", at(21, 1), 5},
		{"yesterday 14:00 local", at(21, 14), 5},
		// Boundary: activity at the exact first instant of the window.
		{"start-day 00:00 local", at(16, 0), 0},
		{"start-day 23:00 local", at(16, 23), 0},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			hist, max := histBuckets(start, []Activity{tc.act}, 7)
			if max != 1 {
				t.Fatalf("max = %d, want 1", max)
			}
			for i, n := range hist {
				want := 0
				if i == tc.wantBin {
					want = 1
				}
				if n != want {
					t.Errorf("hist[%d] = %d, want %d", i, n, want)
				}
			}
		})
	}
}

func TestHistBucketsTokyoLocalDayBoundary(t *testing.T) {
	// Symmetric check for an eastward zone (UTC+9). The broken version
	// mis-bucketed activities in the 00:00-08:59 local window on any
	// non-start day.
	withLocalTZ(t, time.FixedZone("JST", 9*3600))

	start := time.Date(2025, 1, 16, 0, 0, 0, 0, time.Local)

	at := func(localDay, localHour int) Activity {
		return Activity{
			CreatedAt: time.Date(2025, 1, localDay, localHour, 0, 0, 0, time.Local).UTC(),
		}
	}

	cases := []struct {
		name    string
		act     Activity
		wantBin int
	}{
		{"today 05:00 local", at(22, 5), 6},  // previously 5
		{"today 20:00 local", at(22, 20), 6}, // accidentally correct pre-fix
		{"yesterday 05:00 local", at(21, 5), 5},
		{"yesterday 20:00 local", at(21, 20), 5},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			hist, max := histBuckets(start, []Activity{tc.act}, 7)
			if max != 1 {
				t.Fatalf("max = %d, want 1", max)
			}
			for i, n := range hist {
				want := 0
				if i == tc.wantBin {
					want = 1
				}
				if n != want {
					t.Errorf("hist[%d] = %d, want %d", i, n, want)
				}
			}
		})
	}
}

func TestHistBucketsFiltersOutsideWindow(t *testing.T) {
	withLocalTZ(t, time.FixedZone("SP", -3*3600))

	start := time.Date(2025, 1, 16, 0, 0, 0, 0, time.Local)

	before := Activity{CreatedAt: time.Date(2025, 1, 15, 12, 0, 0, 0, time.Local).UTC()}
	inside := Activity{CreatedAt: time.Date(2025, 1, 20, 12, 0, 0, 0, time.Local).UTC()}
	after := Activity{CreatedAt: time.Date(2025, 1, 23, 12, 0, 0, 0, time.Local).UTC()}

	hist, max := histBuckets(start, []Activity{before, inside, after}, 7)
	if max != 1 {
		t.Fatalf("max = %d, want 1", max)
	}
	want := []int{0, 0, 0, 0, 1, 0, 0} // Jan 20 is offset 4 from Jan 16
	for i := range want {
		if hist[i] != want[i] {
			t.Errorf("hist[%d] = %d, want %d", i, hist[i], want[i])
		}
	}
}

func TestHistBucketsAggregatesSameDay(t *testing.T) {
	withLocalTZ(t, time.FixedZone("SP", -3*3600))

	start := time.Date(2025, 1, 16, 0, 0, 0, 0, time.Local)

	// Three activities on Jan 18 local, at different hours that straddle
	// the UTC date boundary (22:00 SP = 01:00 UTC next day).
	activities := []Activity{
		{CreatedAt: time.Date(2025, 1, 18, 2, 0, 0, 0, time.Local).UTC()},
		{CreatedAt: time.Date(2025, 1, 18, 14, 0, 0, 0, time.Local).UTC()},
		{CreatedAt: time.Date(2025, 1, 18, 22, 0, 0, 0, time.Local).UTC()},
	}

	hist, max := histBuckets(start, activities, 7)
	if max != 3 {
		t.Fatalf("max = %d, want 3", max)
	}
	if hist[2] != 3 {
		t.Errorf("hist[2] = %d, want 3 (Jan 18 = offset 2)", hist[2])
	}
	for i := range hist {
		if i == 2 {
			continue
		}
		if hist[i] != 0 {
			t.Errorf("hist[%d] = %d, want 0", i, hist[i])
		}
	}
}

func TestHistogram(t *testing.T) {
	tt := []struct {
		min      int
		max      int
		buckets  int
		expected [][]int
	}{
		{
			min:      0,
			max:      5,
			buckets:  3,
			expected: [][]int{{0}, {1}, {2, 3}, {4, 5}},
		},
		{
			min:      0,
			max:      10,
			buckets:  5,
			expected: [][]int{{0}, {1, 2}, {3, 4}, {5, 6}, {7, 8}, {9, 10}},
		},
		{
			min:      0,
			max:      100,
			buckets:  4,
			expected: [][]int{{0}, seq(1, 25), seq(26, 50), seq(51, 75), seq(76, 100)},
		},
	}

	for _, tc := range tt {
		fit := fitter(tc.min, tc.max, tc.buckets)

		t.Run(fmt.Sprintf("fitter(%d,%d,%d)", tc.min, tc.max, tc.buckets), func(t *testing.T) {
			for expected, nums := range tc.expected {
				for _, num := range nums {
					actual := fit(num)
					if actual != expected {
						t.Errorf(
							"fit(%d) should eq %d, got %d",
							num,
							expected,
							actual,
						)
					}
				}
			}
		})
	}
}

// returns ints in the interval [start, stop]
func seq(start, stop int) []int {
	out := make([]int, stop-start+1)

	for i := range out {
		out[i] = i + start
	}

	return out
}
