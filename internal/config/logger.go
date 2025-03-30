package config

import (
	"fmt"
	"log"
	"os"
	"path"
	"sort"
	"time"
)

type LogInstance struct {
	Log     *log.Logger
	LogFile *os.File
	fields  map[string]interface{}
}

func (l *LogInstance) Close() error {
	return l.LogFile.Close()
}

func CreateLogger() *LogInstance {
	// Create a logger that writes to a file
	logDir := "./logs"
	currDateTime := time.Now()
	logName := fmt.Sprintf("%s/%d%02d%02d_app.log", logDir, currDateTime.Year(), currDateTime.Month(), currDateTime.Day())

	// Create directory if it doesn't exist
	if err := os.MkdirAll(logDir, 0755); err != nil {
		log.Fatalf("Failed to create logs directory: %v", err)
	}

	// Create log file to write to
	file, err := os.OpenFile(logName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	// Limit number of log files to 30
	err = cleanLogDir(logDir)
	if err != nil {
		file.Close()
		log.Fatalf("Failed to clean log directory: %v", err)
	}

	logger := log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	logInstance := LogInstance{
		Log:     logger,
		LogFile: file,
	}

	return &logInstance
}

const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[37m"
)

func (s *State) LogInfo(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	fmt.Println(msg)
	s.Logger.Log.Println("[INFO]", msg)
}

func (s *State) LogDebug(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	s.Logger.Log.Println("[DEBUG]", msg)
}

func (s *State) LogError(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	fmt.Println(msg)
	s.Logger.Log.Println("[ERROR]", msg)
}

func cleanLogDir(dir string) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	if numfiles := len(files); numfiles > 30 {
		// Sort files by modification time (oldest first)
		sort.Slice(files, func(i, j int) bool {
			infoI, _ := files[i].Info()
			infoJ, _ := files[j].Info()
			return infoI.ModTime().Before(infoJ.ModTime())
		})

		// Remove excess oldest files
		for i := 0; i < numfiles-30; i++ {
			err = os.Remove(path.Join(dir, files[i].Name()))
			if err != nil {
				return err
			}
		}
	}
	return nil
}
