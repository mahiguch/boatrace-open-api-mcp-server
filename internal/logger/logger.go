package logger

import (
	"fmt"
	"log"
	"os"
)

// Logger provides structured logging
type Logger struct {
	infoLog  *log.Logger
	errorLog *log.Logger
}

// New creates a new logger
func New() *Logger {
	return &Logger{
		infoLog:  log.New(os.Stdout, "[INFO] ", log.LstdFlags),
		errorLog: log.New(os.Stderr, "[ERROR] ", log.LstdFlags),
	}
}

// Info logs an info message with key-value pairs
func (l *Logger) Info(msg string, fields map[string]interface{}) {
	output := msg
	for k, v := range fields {
		output += " " + k + "=" + toString(v)
	}
	l.infoLog.Println(output)
}

// Error logs an error message with context
func (l *Logger) Error(msg string, err error, fields map[string]interface{}) {
	output := msg + ": " + err.Error()
	for k, v := range fields {
		output += " " + k + "=" + toString(v)
	}
	l.errorLog.Println(output)
}

// toString converts a value to string for logging
func toString(v interface{}) string {
	switch val := v.(type) {
	case string:
		return val
	case int, int64:
		return fmt.Sprintf("%d", val)
	case float64:
		return fmt.Sprintf("%f", val)
	case bool:
		return fmt.Sprintf("%v", val)
	default:
		return fmt.Sprintf("%v", v)
	}
}
