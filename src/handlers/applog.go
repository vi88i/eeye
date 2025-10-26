// Package handlers provides handlers for app logging, interrupt etc.
package handlers

import (
	"io"
	"log"
	"os"
)

// AppLog creates a instance to start capturing the app logs
type AppLog struct {
	handle *os.File
}

// Init initializes the log output based on verbose flag
// If verbose is true, writes to both stdout and app.log file
// If verbose is false, writes only to app.log file
func (a *AppLog) Init(verbose bool) {
	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	if verbose {
		mw := io.MultiWriter(os.Stdout, logFile)
		log.SetOutput(mw)
	} else {
		log.SetOutput(logFile)
	}
	a.handle = logFile
}

// Close closes the multi-writer
func (a *AppLog) Close() {
	if a.handle != nil {
		_ = a.handle.Close()
	}
}

// GetAppLog creates an instance of AppLog handler
func GetAppLog(verbose bool) *AppLog {
	a := AppLog{}
	a.Init(verbose)
	return &a
}
