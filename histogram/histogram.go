package histogram

import (
	"astro/config"
	"astro/date"
	"astro/habit"
	"math"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

var (
	colors = []lipgloss.Style{
		lipgloss.NewStyle().Foreground(lipgloss.Color("#ebedf0")),
		lipgloss.NewStyle().Foreground(lipgloss.Color("#631c03")),
		lipgloss.NewStyle().Foreground(lipgloss.Color("#bd561d")),
		lipgloss.NewStyle().Foreground(lipgloss.Color("#fa7a18")),
		lipgloss.NewStyle().Foreground(lipgloss.Color("#fddf68")),
	}
)

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

func ShortLineHistogram(h habit.Habit, days int) string {
	hist := make([]int, days)
	start := date.Today().AddDate(0, 0, 1-days)

	min, max := 0, 0
	for i := len(h.Activities) - 1; i >= 0; i-- {
		a := h.Activities[i]

		if start.After(a.CreatedAt) {
			break
		}

		diffInDays := date.DiffInDays(start, a.CreatedAt)

		hist[diffInDays]++

		if hist[diffInDays] > max {
			max = hist[diffInDays]
		}

		if hist[diffInDays] < min {
			min = hist[diffInDays]
		}
	}

	fit := fitter(min, max, len(colors)-1)

	var s strings.Builder
	s.Grow(len(hist) * 30)
	for _, day := range hist {
		s.WriteString(colors[fit(day)].Render(config.Graphic))
	}
	return s.String()
}

func Histogram(t time.Time, h habit.Habit, selected int) string {
	hist := make([]int, config.TimeFrameInDays)
	min, max := 0, 0
	for _, a := range h.Activities {
		diffInDays := date.DiffInDays(t, a.CreatedAt)
		if diffInDays >= 0 {
			hist[diffInDays]++

			if hist[diffInDays] > max {
				max = hist[diffInDays]
			}

			if hist[diffInDays] < min {
				min = hist[diffInDays]
			}
		}
	}

	fit := fitter(min, max, len(colors)-1)

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
