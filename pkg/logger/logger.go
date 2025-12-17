package logger

import (
	"log"
	"os"
)

// Logger interface
type Logger interface {
	Info(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Debug(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
}

// SimpleLogger is a simple logger implementation
type SimpleLogger struct {
	infoLog  *log.Logger
	errorLog *log.Logger
	debugLog *log.Logger
	warnLog  *log.Logger
}

// NewLogger creates a new logger
func NewLogger() Logger {
	return &SimpleLogger{
		infoLog:  log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		errorLog: log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
		debugLog: log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile),
		warnLog:  log.New(os.Stdout, "WARN: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

// Info logs an info message
func (l *SimpleLogger) Info(msg string, args ...interface{}) {
	l.infoLog.Printf(msg, args...)
}

// Error logs an error message
func (l *SimpleLogger) Error(msg string, args ...interface{}) {
	l.errorLog.Printf(msg, args...)
}

// Debug logs a debug message
func (l *SimpleLogger) Debug(msg string, args ...interface{}) {
	l.debugLog.Printf(msg, args...)
}

// Warn logs a warning message
func (l *SimpleLogger) Warn(msg string, args ...interface{}) {
	l.warnLog.Printf(msg, args...)
}
