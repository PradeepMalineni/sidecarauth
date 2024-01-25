package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

const logDir = "/path/to/log/directory"
const timestampFormat = "2006-01-02_15-04-05"

func InitLogger() {
	// Generate the log file name with the current timestamp
	currentTime := time.Now()
	applogFileName := fmt.Sprintf("%s/app_%s.log", logDir, currentTime.Format(timestampFormat))

	// Open or create a log file
	logFile, err := os.OpenFile(applogFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("[%s]: Error opening log file: %s", time.Now().Format(timestampFormat), err)
	}
	defer logFile.Close()

	// Set log output to both console and the log file
	log.SetOutput(io.MultiWriter(os.Stdout, logFile))
	log.Print("SideCarAuthSvcs Initializing")
}
