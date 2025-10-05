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

// Init initializes the multi writer to write to both stdout and app.log file
func (a *AppLog) Init() {
	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)
	a.handle = logFile
}

// Close closes the multi-writer
func (a *AppLog) Close() {
	if a.handle != nil {
		_ = a.handle.Close()
	}
}

// GetAppLog creates an instance of AppLog handler
func GetAppLog() *AppLog {
	a := AppLog{}
	a.Init()
	return &a
}
