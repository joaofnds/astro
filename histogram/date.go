package histogram

import "time"

func TruncateDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func EndOfWeek(t time.Time) time.Time {
	return t.AddDate(0, 0, 7-int(t.Weekday()))
}

func SameDay(t1, t2 time.Time) bool {
	return t1.Day() == t2.Day() &&
		t1.Month() == t2.Month() &&
		t1.Year() == t2.Year()
}
