package histogram

import (
	"astroapp/habit"
	"math"
	"strings"
	"time"
)

const (
	timeFrameInDays = 52 * 7

	graphic = "⬢"

	// \033[38;2;<r>;<g>;<b>m
	color0 = "\033[38;2;235;237;240m"
	color1 = "\033[38;2;155;233;168m"
	color2 = "\033[38;2;64;196;99m"
	color3 = "\033[38;2;48;161;78m"
	color4 = "\033[38;2;33;110;57m"

	resetStyle = "\033[0m"
)

var colors = []string{color0, color1, color2, color3, color4}

func fitter(min, max, buckets int) func(n int) int {
	return func(n int) int {
		return int(math.Floor(((float64(buckets) - float64(min)) / (float64(max) - float64(min))) * float64(n)))
	}
}

func Histogram(h habit.Habit) string {
	min, max := 0, 0
	hist := make([]int, timeFrameInDays)
	t := endOfWeek(truncateDay(time.Now())).AddDate(0, 0, -timeFrameInDays)
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
	for i := 0; i < 7; i++ {
		for j := 0; j < 52; j++ {
			n := hist[i+j*7]
			s.WriteString(colors[fit(n)] + graphic + resetStyle)
		}
		s.WriteString("\n")
	}
	return s.String()
}
