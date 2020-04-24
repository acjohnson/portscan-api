package logger

import (
	"fmt"
	"log"
	"os"
)

var f *os.File

func Load() (*log.Logger, error) {
	f, err := os.OpenFile("/tmp/portscan-api.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		fmt.Printf("error opening file: %v", err)
		os.Exit(1)
	}

	logger := log.New(f, "applog: ", log.Lshortfile|log.LstdFlags)
	return logger, err
}
