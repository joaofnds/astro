package histogram

import (
	"astroapp/config"
	"astroapp/habit"
	"math"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

var (
	colors = []lipgloss.Style{
		lipgloss.NewStyle().Foreground(lipgloss.Color("#ebedf0")),
		lipgloss.NewStyle().Foreground(lipgloss.Color("#39d353")),
		lipgloss.NewStyle().Foreground(lipgloss.Color("#26a641")),
		lipgloss.NewStyle().Foreground(lipgloss.Color("#006d32")),
		lipgloss.NewStyle().Foreground(lipgloss.Color("#0e4429")),
	}
)

func fitter(min, max, buckets int) func(n int) int {
	i := ((float64(buckets) - float64(min)) / (float64(max) - float64(min)))
	return func(n int) int {
		return int(math.Floor(i * float64(n)))
	}
}

func Histogram(t time.Time, h habit.Habit, selected int) string {
	hist := make([]int, config.TimeFrameInDays)
	min, max := 0, 0
	for _, a := range h.Activites {
		diffInDays := int(a.CreatedAt.Sub(t).Hours() / 24)
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
	s.Grow(config.TimeFrameInDays*10 + 52)
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
