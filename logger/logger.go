package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

var (
	Debug *log.Logger
	Error *log.Logger
	flags = log.Ldate | log.Ltime | log.Lshortfile | log.Lmsgprefix
)

func Init() error {
	f, err := os.OpenFile("log.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}

	Debug = log.New(f, "[DEBUG] ", flags)
	Error = log.New(f, "[ERROR] ", flags)

	return nil
}

func DebugTime(label string, t1 time.Time) {
	elapsed := time.Since(t1)
	Debug.Printf("%s took: %s", label, elapsed)
}
