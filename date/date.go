package date

import (
	"astro/config"
	"time"
)

func Today() time.Time {
	return TruncateDay(time.Now().Local())
}

func TruncateDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func DiffInDays(t1, t2 time.Time) int {
	days := int(TruncateDay(t1).Sub(TruncateDay(t2)).Hours() / 24)

	if days < 0 {
		return -days
	}
	return days
}

func EndOfWeek(t time.Time) time.Time {
	return t.AddDate(0, 0, 7-int(t.Weekday()))
}

func SameDay(t1, t2 time.Time) bool {
	t1, t2 = t1.Local(), t2.Local()
	return t1.Day() == t2.Day() &&
		t1.Month() == t2.Month() &&
		t1.Year() == t2.Year()
}

func TimeFrame() (time.Time, time.Time) {
	end := EndOfWeek(Today())
	beg := end.AddDate(0, 0, -config.TimeFrameInDays)
	return beg, end
}

func CombineDateWithTime(baseDate time.Time, baseTime time.Time) time.Time {
	return time.Date(
		baseDate.Year(),
		baseDate.Month(),
		baseDate.Day(),

		baseTime.Hour(),
		baseTime.Minute(),
		baseTime.Second(),
		baseTime.Nanosecond(),
		baseTime.Location(),
	)
}
