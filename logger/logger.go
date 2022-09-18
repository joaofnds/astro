package logger

import (
	"log"
	"os"
)

var (
	Debug *log.Logger
	Error *log.Logger
	flags = log.Ldate | log.Ltime | log.Lshortfile | log.Lmsgprefix
)

func init() {
	f, err := os.OpenFile("log.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("failed to open log file: %v", err)
	}

	Debug = log.New(f, "[DEBUG] ", flags)
	Error = log.New(f, "[ERROR] ", flags)
}
