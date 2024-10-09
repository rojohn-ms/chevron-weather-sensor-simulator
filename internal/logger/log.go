// Package logger contains simple loggers.
package logger

import (
	"log"
	"strings"
)

const (
	// InfoLogLevel prints everything
	InfoLogLevel = "INFO"
	// DebugLogLevel prints only debug
	DebugLogLevel = "DEBUG"
	// NoneLogLevel prints nothing
	NoneLogLevel = "NONE"
)

type (
	// Log is an interface for logging
	Log interface {
		Printf(format string, v ...any)

		Println(v ...any)
	}

	// LogLevel is the log level to print
	LogLevel string

	internalLog struct {
		lvl LogLevel
	}
)

// New returns a new Log.
func New(logLevel string) Log {
	logLvl := NoneLogLevel
	logLvlLower := strings.ToLower(logLevel)
	switch logLvlLower {
	case "info":
		logLvl = InfoLogLevel

	case "debug":
		logLvl = DebugLogLevel
	}

	return &internalLog{
		lvl: LogLevel(logLvl),
	}
}

// Printf prints to the console.
func (l *internalLog) Printf(format string, v ...any) {
	if l.lvl == InfoLogLevel {
		log.Printf(format, v...)
	}
}

// Println prints to the console.
func (l *internalLog) Println(v ...any) {
	if l.lvl == InfoLogLevel {
		log.Println(v...)
	}
}
