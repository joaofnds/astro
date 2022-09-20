package logger

import "time"

func DebugTime(label string, t1 time.Time) {
	elapsed := time.Since(t1)
	Debug.Printf("%s took: %s", label, elapsed)
}
