package histogram

import (
	"astroapp/habit"
	"math"
	"strings"
	"time"
)

const (
	TimeFrameInDays = 52 * 7

	graphic = "⬛"

	selectedGraphic    = "⭕"
	selectedBackground = "\033[48;2;0;0;0m"

	// \033[38;2;<r>;<g>;<b>m
	color0 = "\033[38;2;235;237;240m"
	color1 = "\033[38;2;155;233;168m"
	color2 = "\033[38;2;64;196;99m"
	color3 = "\033[38;2;48;161;78m"
	color4 = "\033[38;2;33;110;57m"

	resetStyle = "\033[m"
)

var colors = []string{color0, color1, color2, color3, color4}

func fitter(min, max, buckets int) func(n int) int {
	i := ((float64(buckets) - float64(min)) / (float64(max) - float64(min)))
	return func(n int) int {
		return int(math.Floor(i * float64(n)))
	}
}

func Histogram(t time.Time, h habit.Habit, selected int) string {
	hist := make([]int, TimeFrameInDays)
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
	s.Grow(TimeFrameInDays*len(colors[0]+graphic+resetStyle) + 52)
	for weekday := 0; weekday < 7; weekday++ {
		switch weekday {
		case 1:
			s.WriteString("mon ")
		case 3:
			s.WriteString("wed ")
		case 5:
			s.WriteString("fri ")
		default:
			s.WriteString("    ")
		}

		for week := 0; week < 52; week++ {
			day := weekday + week*7
			g := graphic
			if day == selected {
				g = selectedGraphic
			}
			s.WriteString(colors[fit(hist[day])] + g + resetStyle)
		}
		s.WriteString("\n")
	}
	return s.String()
}
