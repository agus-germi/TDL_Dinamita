package logger

import (
	"io"
)

// Common interface for loggers
type Logger interface {
	Println(i ...interface{})
	Debug(i ...interface{})
	Debugf(format string, args ...interface{})
	Info(i ...interface{})
	Infof(format string, args ...interface{})
	Warn(i ...interface{})
	Warnf(format string, args ...interface{})
	Error(i ...interface{})
	Errorf(format string, args ...interface{})
	Fatal(i ...interface{})
	Fatalf(format string, args ...interface{})
	Panic(i ...interface{})
	Panicf(format string, args ...interface{})
	Writer() io.Writer // Allow integration with Echo or any thrid party library.
}

// This log is a global variable that will be used by the whole application, and it "follows" the singleton pattern.
var Log Logger

func init() {
	Log = NewLogrusLoggerAdapter() // If we want to change the logger framework we currently use, we just need to create a new adapter and change this line.
}

func New() Logger {
	return Log
}
