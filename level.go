package log

import (
	"io"

	"github.com/rs/zerolog"
)

// Level defines log levels.
type Level = zerolog.Level

const (
	LevelDebug    = zerolog.DebugLevel
	LevelInfo     = zerolog.InfoLevel
	LevelWarn     = zerolog.WarnLevel
	LevelError    = zerolog.ErrorLevel
	LevelFatal    = zerolog.FatalLevel
	LevelPanic    = zerolog.PanicLevel
	LevelNone     = zerolog.NoLevel
	LevelDisabled = zerolog.Disabled
	LevelTrace    = zerolog.TraceLevel
)

// LevelWriter ...
type LevelWriter struct {
	io.Writer
	Level Level
}

// WriteLevel ...
func (w *LevelWriter) WriteLevel(level Level, p []byte) (n int, err error) {
	if level >= w.Level {
		return w.Writer.Write(p)
	}
	return len(p), nil
}
