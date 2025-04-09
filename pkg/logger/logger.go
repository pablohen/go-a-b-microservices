package logger

import (
	"log"
	"os"
)

type Logger interface {
	Info(message string, args ...interface{})
	Error(message string, args ...interface{})
	Debug(message string, args ...interface{})
}

type SimpleLogger struct {
	infoLogger  *log.Logger
	errorLogger *log.Logger
	debugLogger *log.Logger
}

func NewLogger() *SimpleLogger {
	return &SimpleLogger{
		infoLogger:  log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		errorLogger: log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
		debugLogger: log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

func (l *SimpleLogger) Info(message string, args ...interface{}) {
	l.infoLogger.Printf(message, args...)
}

func (l *SimpleLogger) Error(message string, args ...interface{}) {
	l.errorLogger.Printf(message, args...)
}

func (l *SimpleLogger) Debug(message string, args ...interface{}) {
	l.debugLogger.Printf(message, args...)
}
