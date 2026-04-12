package domain

import (
	"astro/config"
	"astro/date"
	"math"
	"strings"
	"time"

	"charm.land/lipgloss/v2"
)

var colors = []lipgloss.Style{
	lipgloss.NewStyle().Foreground(lipgloss.Color("#ebedf0")),
	lipgloss.NewStyle().Foreground(lipgloss.Color("#39d353")),
	lipgloss.NewStyle().Foreground(lipgloss.Color("#26a641")),
	lipgloss.NewStyle().Foreground(lipgloss.Color("#006d32")),
	lipgloss.NewStyle().Foreground(lipgloss.Color("#0e4429")),
}

func fitter(min, max, buckets int) func(n int) int {
	if min == 0 && max == 0 {
		return func(_ int) int {
			return 0
		}
	}
	i := float64(buckets-min) / float64(max-min)
	return func(n int) int {
		return int(math.Ceil(i * float64(n)))
	}
}

// histBuckets counts activities into `size` day buckets indexed by local
// calendar-day offset from start. hist[i] is the number of activities whose
// local day equals start's local day + i. Activities outside
// [start, start+size) are ignored. `start` must already be at local midnight.
//
// Each activity's CreatedAt is converted to local time before the day-offset
// is computed. The backend stores CreatedAt in UTC, but date.DiffInDays
// truncates each argument in its own Location — mixing a local-midnight start
// with a UTC CreatedAt produces off-by-one errors at the local-midnight
// boundary, which is exactly the activity-graph bug this fix addresses.
func histBuckets(start time.Time, activities []Activity, size int) ([]int, int) {
	hist := make([]int, size)
	max := 0
	for _, a := range activities {
		if a.CreatedAt.Before(start) {
			continue
		}
		diff := date.DiffInDays(start, a.CreatedAt.Local())
		if diff >= size {
			continue
		}
		hist[diff]++
		if hist[diff] > max {
			max = hist[diff]
		}
	}
	return hist, max
}

func ShortLineHistogram(activities []Activity, days int) string {
	start := date.Today().AddDate(0, 0, 1-days)
	hist, max := histBuckets(start, activities, days)
	fit := fitter(0, max, len(colors)-1)

	var s strings.Builder
	s.Grow(len(hist) * 30)
	for _, day := range hist {
		s.WriteString(colors[fit(day)].Render(config.Graphic))
	}
	return s.String()
}

func Histogram(t time.Time, activities []Activity, selected int) string {
	hist, max := histBuckets(t, activities, config.TimeFrameInDays)
	fit := fitter(0, max, len(colors)-1)

	var s strings.Builder
	s.Grow(config.TimeFrameInDays * 30)
	for weekday := 0; weekday < 7; weekday++ {
		switch weekday {
		case 1:
			s.WriteString("Mon ")
		case 3:
			s.WriteString("Wed ")
		case 5:
			s.WriteString("Fri ")
		default:
			s.WriteString("    ")
		}

		for week := 0; week < 52; week++ {
			day := weekday + week*7
			g := config.Graphic
			if day == selected {
				g = config.SelectedGraphic
			}
			s.WriteString(colors[fit(hist[day])].Render(g))
		}
		s.WriteString("\n")
	}
	return s.String()
}
