package log

import (
	"fmt"
	"io"

	kitlog "github.com/go-kit/kit/log"
	kitlevel "github.com/go-kit/kit/log/level"
	"github.com/go-kit/kit/log/term"
)

const (
	msgKey    = "_msg" // "_" prefixed to avoid collisions
	moduleKey = "module"
)

type tmLogger struct {
	srcLogger kitlog.Logger
}

// Interface assertions
var _ Logger = (*tmLogger)(nil)

// NewTMTermLogger returns a logger that encodes msg and keyvals to the Writer
// using go-kit's log as an underlying logger and our custom formatter. Note
// that underlying logger could be swapped with something else.
func NewTMLogger(w io.Writer) Logger {
	// Color by level value
	colorFn := func(keyvals ...interface{}) term.FgBgColor {
		if keyvals[0] != kitlevel.Key() {
			panic(fmt.Sprintf("expected level key to be first, got %v", keyvals[0]))
		}
		switch keyvals[1].(kitlevel.Value).String() {
		case "debug":
			return term.FgBgColor{Fg: term.DarkGray}
		case "error":
			return term.FgBgColor{Fg: term.Red}
		default:
			return term.FgBgColor{}
		}
	}

	return &tmLogger{term.NewLogger(w, NewTMFmtLogger, colorFn)}
}

// NewTMLoggerWithColorFn allow you to provide your own color function. See
// NewTMLogger for documentation.
func NewTMLoggerWithColorFn(w io.Writer, colorFn func(keyvals ...interface{}) term.FgBgColor) Logger {
	return &tmLogger{term.NewLogger(w, NewTMFmtLogger, colorFn)}
}

// Info logs a message at level Info.
func (l *tmLogger) Info(msg string, keyvals ...interface{}) error {
	lWithLevel := kitlevel.Info(l.srcLogger)
	return kitlog.With(lWithLevel, msgKey, msg).Log(keyvals...)
}

// Debug logs a message at level Debug.
func (l *tmLogger) Debug(msg string, keyvals ...interface{}) error {
	lWithLevel := kitlevel.Debug(l.srcLogger)
	return kitlog.With(lWithLevel, msgKey, msg).Log(keyvals...)
}

// Error logs a message at level Error.
func (l *tmLogger) Error(msg string, keyvals ...interface{}) error {
	lWithLevel := kitlevel.Error(l.srcLogger)
	return kitlog.With(lWithLevel, msgKey, msg).Log(keyvals...)
}

// With returns a new contextual logger with keyvals prepended to those passed
// to calls to Info, Debug or Error.
func (l *tmLogger) With(keyvals ...interface{}) Logger {
	return &tmLogger{kitlog.With(l.srcLogger, keyvals...)}
}
