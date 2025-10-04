// Package applog helps to retain the logs of the analysis done for verification
package applog

import (
	"io"
	"log"
	"os"
)

// Init initializes the multi writer to write to both stdout and app.log file
func Init() *os.File {
	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)
	return logFile
}
