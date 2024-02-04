// logger/logger.go
package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

const logDir = "logs" // Specify your desired log directory

func InitLogger() {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		log.Fatalf("Error creating log directory: %s", err)
	}
	formattedTime := "2006-01-02"
	logFileName := fmt.Sprintf("%s/app_%s.log", logDir, time.Now().Format(formattedTime))

	//logFileName := filepath.Join(logDir, "app.log")
	logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error opening log file: %s", err)
	}

	log.SetOutput(logFile)
}

// Log logs a message with a timestamp
func Log(message string) {
	log.Printf("%s", message)
}

// LogF logs a variable name and its value with a timestamp
func LogF(variableName string, variableValue interface{}) {
	log.Printf("%s: %v", variableName, variableValue)
}
