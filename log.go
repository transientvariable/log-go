package log

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/transientvariable/config-go"

	"github.com/mattn/go-colorable"
	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"

	zlog "github.com/rs/zerolog/log"
	stdLog "log"
)

const (
	callerSkipFrames        = 3
	defaultSize             = 10
	defaultRetentionAge     = 14
	defaultRetentionBackups = 10
)

var (
	defaultLogger atomic.Value
	once          sync.Once
)

// Logger defines the type for the logger.
type Logger = zerolog.Logger

// Init initializes the logging system and sets the default logger. If the default logger has already been set
// (e.g. SetDefault), no further action is taken.
func Init() {
	once.Do(func() {
		if _, ok := defaultLogger.Load().(*Logger); ok {
			return
		}

		// call to New() should never be nil, but panic nonetheless.
		if err := SetDefault(New()); err != nil {
			panic(err)
		}
	})
}

// New creates a new logger using the provided Option properties.
func New(options ...func(*Option)) *Logger {
	opts := &Option{}
	for _, opt := range options {
		opt(opts)
	}

	logger := zlog.Output(zerolog.ConsoleWriter{Out: colorable.NewColorableStderr()}).With().Timestamp().Caller().Logger()
	switch opts.level {
	case "debug":
		logger = logger.Level(LevelDebug)
		break
	case "error":
		logger = logger.Level(LevelError)
		break
	case "fatal":
		logger = logger.Level(LevelFatal)
		break
	case "panic":
		logger = logger.Level(LevelPanic)
		break
	case "trace":
		logger = logger.Level(LevelTrace)
		break
	case "warn":
		logger = logger.Level(LevelWarn)
		break
	default:
		logger = logger.Level(LevelInfo)
	}
	return &logger
}

// SetDefault sets the default logger used by all logging functions within the log package.
func SetDefault(logger *Logger) error {
	if logger == nil {
		return errors.New("log: default logger cannot be nil")
	}
	defaultLogger.Store(logger)
	return nil
}

// Default returns the default logger. If the default logger has not been initialized, Init is called before returning.
func Default() *Logger {
	Init()
	return defaultLogger.Load().(*Logger)
}

// Debug records debug log event with the provided msg and arguments.
func Debug(msg string, args ...func(*Record)) {
	Log(LevelDebug, msg, args...)
}

// Error records and error log event on the default logger with the provided msg and arguments.
func Error(msg string, args ...func(*Record)) {
	Log(LevelError, msg, args...)
}

// Fatal records a fatal log event with the provided msg and arguments, then immediately terminates program execution
// by calling os.Exit(1).
func Fatal(msg string, args ...func(*Record)) {
	Log(LevelFatal, msg, args...)
}

// Info calls info on the default logger.
func Info(msg string, args ...func(*Record)) {
	Log(LevelInfo, msg, args...)
}

// Panic calls panic on the default logger.
func Panic(msg string, args ...func(*Record)) {
	Log(LevelPanic, msg, args...)
}

// Trace calls trace on the default logger.
func Trace(msg string, args ...func(*Record)) {
	Log(LevelTrace, msg, args...)
}

// Warn calls warn on the default logger.
func Warn(msg string, args ...func(*Record)) {
	Log(LevelWarn, msg, args...)
}

// Log records a log event using the provided Level.
func Log(level Level, msg string, args ...func(*Record)) {
	var event *zerolog.Event
	switch level {
	case LevelDebug:
		lvl := Default().GetLevel()
		if lvl == LevelDebug || lvl == LevelTrace {
			event = Default().Debug()
		}
		break
	case LevelError:
		event = Default().Error()
		break
	case LevelFatal:
		event = Default().Fatal()
		break
	case LevelInfo:
		event = Default().Info()
		break
	case LevelPanic:
		event = Default().Panic()
		break
	case LevelTrace:
		if Default().GetLevel() == LevelTrace {
			event = Default().Trace()
		}
		break
	case LevelWarn:
		event = Default().Warn()
		break
	default:
		// no-op
	}

	if event != nil {
		handleEvent(event, msg, args...)
	}
}

func handleEvent(e *zerolog.Event, msg string, args ...func(*Record)) {
	r := acquireRecord()
	defer releaseRecord(r)

	for _, arg := range args {
		arg(r)
	}

	e.CallerSkipFrame(callerSkipFrames)
	e.Ctx(r.ctx)
	e.Err(r.err)

	for k, attr := range r.attrs {
		switch attr.kind {
		case kindAny:
			e.Any(k, attr.value)
		case kindBool:
			e.Bool(k, attr.value.(bool))
		case kindDuration:
			e.Dur(k, attr.value.(time.Duration))
		case kindFloat32:
			e.Float32(k, attr.value.(float32))
		case kindFloat64:
			e.Float64(k, attr.value.(float64))
		case kindInt64:
			e.Int64(k, attr.value.(int64))
		case kindString:
			e.Str(k, attr.value.(string))
		case kindTime:
			e.Time(k, attr.value.(time.Time))
		case kindUint64:
			e.Uint64(k, attr.value.(uint64))
		default:
			// no-op
		}
	}
	e.Msg(msg)
}

func prepareFileWriters(file string, levels ...zerolog.Level) []*LevelWriter {
	if len(levels) == 0 {
		return nil
	}

	lj := &lumberjack.Logger{
		Compress:   false,
		Filename:   file,
		LocalTime:  false,
		MaxAge:     intValue(FileRetentionAge, defaultRetentionAge),
		MaxBackups: intValue(FileRetentionBackups, defaultRetentionBackups),
		MaxSize:    intValue(FileSize, defaultSize),
	}

	var w []*LevelWriter
	for _, l := range levels {
		w = append(w, &LevelWriter{Writer: lj, Level: l})
	}
	return w
}

func prepareDir(path string) string {
	dir := strings.TrimSpace(path)
	if err := statDir(dir); err != nil {
		stdLog.Println(fmt.Errorf("log: could not stat directory: %s: %w", dir, err))
		dir = os.TempDir()
		stdLog.Printf("log: using default directory: %s", dir)
	}
	return dir
}

func statDir(path string) error {
	fi, err := os.Stat(path)
	if err != nil {
		return err
	}

	if !fi.Mode().IsDir() {
		return fmt.Errorf("log: path is not a directory: %s", path)
	}

	if _, err := os.OpenFile(path, os.O_WRONLY, fi.Mode()); err != nil {
		return err
	}
	return nil
}

func intValue(path string, defaultValue int) int {
	if s, _ := config.Int(path); s >= 0 {
		return s
	}
	return defaultValue
}
