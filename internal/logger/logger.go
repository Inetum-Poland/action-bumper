// Copyright (c) 2024 Inetum Poland.

package logger

import (
	"io"
	"log"
	"os"
)

// Logger defines the logging interface
type Logger interface {
	Printf(format string, v ...interface{})
	Println(v ...interface{})
}

// StandardLogger wraps standard log package
type StandardLogger struct {
	logger *log.Logger
}

// NewStandardLogger creates a new standard logger
func NewStandardLogger(w io.Writer, prefix string, flag int) *StandardLogger {
	return &StandardLogger{
		logger: log.New(w, prefix, flag),
	}
}

// NewDefaultLogger creates a default logger for production
func NewDefaultLogger() *StandardLogger {
	return &StandardLogger{
		logger: log.New(os.Stdout, "", log.LstdFlags),
	}
}

// Printf logs formatted message
func (l *StandardLogger) Printf(format string, v ...interface{}) {
	l.logger.Printf(format, v...)
}

// Println logs message with newline
func (l *StandardLogger) Println(v ...interface{}) {
	l.logger.Println(v...)
}

// NoopLogger is a logger that does nothing (useful for tests)
type NoopLogger struct{}

// NewNoopLogger creates a no-op logger
func NewNoopLogger() *NoopLogger {
	return &NoopLogger{}
}

// Printf does nothing
func (l *NoopLogger) Printf(format string, v ...interface{}) {}

// Println does nothing
func (l *NoopLogger) Println(v ...interface{}) {}
