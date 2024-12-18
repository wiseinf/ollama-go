package ollama

import (
	"fmt"
	"log"
	"os"
)

// Logger interface
type Logger interface {
	Debug(format string, v ...interface{})
	Info(format string, v ...interface{})
	Error(format string, v ...interface{})
}

// default logger
type defaultLogger struct {
	debug *log.Logger
	info  *log.Logger
	error *log.Logger
}

func newDefaultLogger() *defaultLogger {
	return &defaultLogger{
		debug: log.New(os.Stdout, "[DEBUG] ", log.LstdFlags),
		info:  log.New(os.Stdout, "[INFO] ", log.LstdFlags),
		error: log.New(os.Stderr, "[ERROR] ", log.LstdFlags),
	}
}

func (l *defaultLogger) Debug(format string, v ...interface{}) {
	l.debug.Output(2, fmt.Sprintf(format, v...))
}

func (l *defaultLogger) Info(format string, v ...interface{}) {
	l.info.Output(2, fmt.Sprintf(format, v...))
}

func (l *defaultLogger) Error(format string, v ...interface{}) {
	l.error.Output(2, fmt.Sprintf(format, v...))
}
