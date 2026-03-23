package date_test

import (
	"astro/config"
	"astro/date"
	"testing"
	"time"
)

func TestTruncateDay(t *testing.T) {
	tt := []struct {
		name string
		in   time.Time
		want time.Time
	}{
		{
			name: "zeroes hours/minutes/seconds/nanoseconds",
			in:   time.Date(2025, 3, 15, 14, 30, 45, 123456789, time.UTC),
			want: time.Date(2025, 3, 15, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "already truncated returns identical value",
			in:   time.Date(2025, 3, 15, 0, 0, 0, 0, time.UTC),
			want: time.Date(2025, 3, 15, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "preserves location",
			in:   time.Date(2025, 6, 1, 23, 59, 59, 0, time.FixedZone("EST", -5*3600)),
			want: time.Date(2025, 6, 1, 0, 0, 0, 0, time.FixedZone("EST", -5*3600)),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			got := date.TruncateDay(tc.in)
			if !got.Equal(tc.want) {
				t.Errorf("TruncateDay(%v) = %v, want %v", tc.in, got, tc.want)
			}
			if got.Location().String() != tc.in.Location().String() {
				t.Errorf("TruncateDay(%v) location = %v, want %v", tc.in, got.Location(), tc.in.Location())
			}
		})
	}
}

func TestDiffInDays(t *testing.T) {
	base := time.Date(2025, 3, 15, 12, 0, 0, 0, time.UTC)

	tt := []struct {
		name string
		t1   time.Time
		t2   time.Time
		want int
	}{
		{
			name: "same time returns 0",
			t1:   base,
			t2:   base,
			want: 0,
		},
		{
			name: "1 day apart returns 1",
			t1:   base.AddDate(0, 0, 1),
			t2:   base,
			want: 1,
		},
		{
			name: "7 days apart returns 7",
			t1:   base.AddDate(0, 0, 7),
			t2:   base,
			want: 7,
		},
		{
			name: "negative difference returns positive value",
			t1:   base,
			t2:   base.AddDate(0, 0, 5),
			want: 5,
		},
		{
			name: "different hours same day returns 0",
			t1:   time.Date(2025, 3, 15, 1, 0, 0, 0, time.UTC),
			t2:   time.Date(2025, 3, 15, 23, 0, 0, 0, time.UTC),
			want: 0,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			got := date.DiffInDays(tc.t1, tc.t2)
			if got != tc.want {
				t.Errorf("DiffInDays(%v, %v) = %d, want %d", tc.t1, tc.t2, got, tc.want)
			}
		})
	}
}

func TestSameDay(t *testing.T) {
	tt := []struct {
		name string
		t1   time.Time
		t2   time.Time
		want bool
	}{
		{
			name: "identical times return true",
			t1:   time.Date(2025, 3, 15, 14, 30, 0, 0, time.Local),
			t2:   time.Date(2025, 3, 15, 14, 30, 0, 0, time.Local),
			want: true,
		},
		{
			name: "same date different hours return true",
			t1:   time.Date(2025, 3, 15, 1, 0, 0, 0, time.Local),
			t2:   time.Date(2025, 3, 15, 23, 59, 0, 0, time.Local),
			want: true,
		},
		{
			name: "different dates return false",
			t1:   time.Date(2025, 3, 15, 12, 0, 0, 0, time.Local),
			t2:   time.Date(2025, 3, 16, 12, 0, 0, 0, time.Local),
			want: false,
		},
		{
			name: "midnight boundary returns false",
			t1:   time.Date(2025, 3, 15, 23, 59, 59, 0, time.Local),
			t2:   time.Date(2025, 3, 16, 0, 0, 0, 0, time.Local),
			want: false,
		},
		{
			name: "different timezones same local date return true",
			t1:   time.Date(2025, 3, 15, 20, 0, 0, 0, time.UTC),
			t2:   time.Date(2025, 3, 15, 10, 0, 0, 0, time.FixedZone("EST", -5*3600)),
			want: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			got := date.SameDay(tc.t1, tc.t2)
			if got != tc.want {
				t.Errorf("SameDay(%v, %v) = %v, want %v", tc.t1, tc.t2, got, tc.want)
			}
		})
	}
}

func TestEndOfWeek(t *testing.T) {
	tt := []struct {
		name string
		in   time.Time
		want time.Time
	}{
		{
			name: "Sunday adds 7 days",
			in:   time.Date(2025, 3, 16, 0, 0, 0, 0, time.UTC), // Sunday
			want: time.Date(2025, 3, 23, 0, 0, 0, 0, time.UTC), // next Sunday
		},
		{
			name: "Monday returns following Sunday",
			in:   time.Date(2025, 3, 17, 0, 0, 0, 0, time.UTC), // Monday
			want: time.Date(2025, 3, 23, 0, 0, 0, 0, time.UTC), // Sunday
		},
		{
			name: "Saturday returns next day Sunday",
			in:   time.Date(2025, 3, 22, 0, 0, 0, 0, time.UTC), // Saturday
			want: time.Date(2025, 3, 23, 0, 0, 0, 0, time.UTC), // Sunday
		},
		{
			name: "Wednesday returns following Sunday",
			in:   time.Date(2025, 3, 19, 0, 0, 0, 0, time.UTC), // Wednesday
			want: time.Date(2025, 3, 23, 0, 0, 0, 0, time.UTC), // Sunday
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			got := date.EndOfWeek(tc.in)
			if !got.Equal(tc.want) {
				t.Errorf("EndOfWeek(%v) = %v (weekday %v), want %v (weekday %v)",
					tc.in, got, got.Weekday(), tc.want, tc.want.Weekday())
			}
		})
	}
}

func TestToday(t *testing.T) {
	got := date.Today()

	if got.Hour() != 0 || got.Minute() != 0 || got.Second() != 0 || got.Nanosecond() != 0 {
		t.Errorf("Today() has non-zero time components: %v", got)
	}

	if got.Location() != time.Now().Local().Location() {
		t.Errorf("Today() location = %v, want Local (%v)", got.Location(), time.Now().Local().Location())
	}

	now := time.Now().Local()
	if got.Year() != now.Year() || got.Month() != now.Month() || got.Day() != now.Day() {
		t.Errorf("Today() date = %v-%v-%v, want %v-%v-%v",
			got.Year(), got.Month(), got.Day(),
			now.Year(), now.Month(), now.Day())
	}
}

func TestTimeFrame(t *testing.T) {
	beg, end := date.TimeFrame()

	if end.Weekday() != time.Sunday {
		t.Errorf("TimeFrame() end weekday = %v, want Sunday", end.Weekday())
	}

	span := date.DiffInDays(end, beg)
	if span != config.TimeFrameInDays {
		t.Errorf("TimeFrame() span = %d days, want %d", span, config.TimeFrameInDays)
	}

	if end.Hour() != 0 || end.Minute() != 0 || end.Second() != 0 {
		t.Errorf("TimeFrame() end has non-zero time: %v", end)
	}

	if beg.Hour() != 0 || beg.Minute() != 0 || beg.Second() != 0 {
		t.Errorf("TimeFrame() beg has non-zero time: %v", beg)
	}
}

func TestCombineDateWithTime(t *testing.T) {
	tt := []struct {
		name     string
		baseDate time.Time
		baseTime time.Time
		want     time.Time
	}{
		{
			name:     "takes date from first arg and time from second",
			baseDate: time.Date(2025, 3, 15, 0, 0, 0, 0, time.UTC),
			baseTime: time.Date(2000, 1, 1, 14, 30, 45, 0, time.UTC),
			want:     time.Date(2025, 3, 15, 14, 30, 45, 0, time.UTC),
		},
		{
			name:     "preserves nanoseconds from time arg",
			baseDate: time.Date(2025, 6, 20, 0, 0, 0, 0, time.UTC),
			baseTime: time.Date(2000, 1, 1, 8, 15, 0, 123456789, time.UTC),
			want:     time.Date(2025, 6, 20, 8, 15, 0, 123456789, time.UTC),
		},
		{
			name:     "uses location from time arg",
			baseDate: time.Date(2025, 3, 15, 0, 0, 0, 0, time.UTC),
			baseTime: time.Date(2000, 1, 1, 10, 0, 0, 0, time.FixedZone("EST", -5*3600)),
			want:     time.Date(2025, 3, 15, 10, 0, 0, 0, time.FixedZone("EST", -5*3600)),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			got := date.CombineDateWithTime(tc.baseDate, tc.baseTime)
			if !got.Equal(tc.want) {
				t.Errorf("CombineDateWithTime(%v, %v) = %v, want %v",
					tc.baseDate, tc.baseTime, got, tc.want)
			}
			if got.Location().String() != tc.baseTime.Location().String() {
				t.Errorf("CombineDateWithTime location = %v, want %v",
					got.Location(), tc.baseTime.Location())
			}
		})
	}
}
