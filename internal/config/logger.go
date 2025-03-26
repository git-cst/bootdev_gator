package config

import (
	"fmt"
	"log"
	"os"
	"time"
)

func CreateLogger() *log.Logger {
	// Create a logger that writes to a file
	logDir := "./logs"
	currDateTime := time.Now()
	logName := fmt.Sprintf("%s/%d%02d%02d_app.log", logDir, currDateTime.Year(), currDateTime.Month(), currDateTime.Day())

	// Create directory if it doesn't exist
	if err := os.MkdirAll(logDir, 0755); err != nil {
		log.Fatal("Failed to create logs directory")
	}

	file, err := os.OpenFile(logName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open log file")
	}

	logger := log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)

	return logger
}

func (s *State) LogInfo(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	fmt.Println(msg)
	s.Logger.Println("[INFO]", msg)
}

func (s *State) LogDebug(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	s.Logger.Println("[DEBUG]", msg)
}

func (s *State) LogError(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	fmt.Println(msg)
	s.Logger.Println("[ERROR]", msg)
}
