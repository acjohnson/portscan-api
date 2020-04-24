package logger

import (
	"fmt"
	"log"
	"os"
)

var f *os.File

func Load(log_level string) (*log.Logger, error) {
	f, err := os.OpenFile("/tmp/portscan-api.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		fmt.Printf("error opening file: %v", err)
		os.Exit(1)
	}

	prefix_str := "info: "

	if log_level == "DEBUG" {
		prefix_str = "debug: "
	}
	if log_level == "WARN" {
		prefix_str = "warn: "
	}

	logger := log.New(f, prefix_str, log.Lshortfile|log.LstdFlags)
	return logger, err
}
