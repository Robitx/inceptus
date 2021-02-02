// Package log provides basic logging mechanism.
//
// You can log into physical file, but log rotation is not implemented,
// see https://12factor.net/logs =>
// Ideally write to /dev/stderr|stdout|null and use something like
// https://github.com/logrotate/logrotate or
// https://github.com/agnivade/funnel for log post processing
// (funnel supports multiple targets, storing into files with rotation and compression)
// so let's not reinvent the wheel here..
// TODO:
// - Enrich logger with methods of std logger for use in third party libs that needs it
// - If you get really bored some day:
// -- implement rolling files with https://github.com/natefinch/lumberjack
// -- allow writing into multiple files with io.MultiWriter
package log

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

// Logger represents wrapper for logging mechanism
type Logger struct {
	zerolog.Logger
	file *os.File
	mask string
}

// OpenFileHelper tries to make or create file (including dirs)
// with proper flags for use in logger
func OpenFileHelper(name string) (*os.File, error) {
	_ = os.MkdirAll(filepath.Dir(name), os.ModePerm)
	return os.OpenFile(name,
		os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
}

// closeFile tries to close file and returns it's name
func (l *Logger) closeFile() string {
	if l.file == nil {
		return ""
	}

	fileName := l.file.Name()
	_ = l.file.Close()
	l.file = nil
	return fileName
}

// Close cleans up after logger
func (l *Logger) Close() {
	l.Info().
		Str("file", l.file.Name()).
		Msg("closing log")
	l.closeFile()
	l = nil
}

// TimestampHook returns hook for logging UnixNano timestamp as float64
func TimestampHook() zerolog.HookFunc {
	return zerolog.HookFunc(
		func(e *zerolog.Event, level zerolog.Level, msg string) {
			e.Float64("timestamp", float64(time.Now().UnixNano())/1e9)
		})
}

// StaticHook returns hook for logging static name, value pair to every probe
// (like hostname, PID, local IP and such..), basic types are cast to avoid reflection
func StaticHook(name string, value interface{}) zerolog.HookFunc {
	switch typedValue := value.(type) {
	case float64:
		return zerolog.HookFunc(
			func(e *zerolog.Event, level zerolog.Level, msg string) {
				e.Float64(name, typedValue)
			})
	case int:
		return zerolog.HookFunc(
			func(e *zerolog.Event, level zerolog.Level, msg string) {
				e.Int(name, typedValue)
			})
	case string:
		return zerolog.HookFunc(
			func(e *zerolog.Event, level zerolog.Level, msg string) {
				e.Str(name, typedValue)
			})
	default:
		return zerolog.HookFunc(
			func(e *zerolog.Event, level zerolog.Level, msg string) {
				e.Interface(name, typedValue)
			})
	}
}

// New creates logger
// if opening log file fails, logger will fallback to stderr
// you might want to use a log router like https://github.com/agnivade/funnel
// to handle app log output
func New(fileName string, mask string, hooks ...zerolog.HookFunc) *Logger {
	// preparing log level mask with fallback to debug
	level, merr := zerolog.ParseLevel(strings.ToLower(mask))
	if merr != nil {
		level = zerolog.DebugLevel
	}

	file, ferr := OpenFileHelper(fileName)
	if ferr != nil {
		file = os.Stderr
	}

	l := &Logger{}

	zerolog.DisableSampling(true)

	zlogger := zerolog.New(file).Level(level).
		// With().Caller().Logger().Hook(timestampHook)
		With().Logger()

	// Register logger hook functions
	for _, hook := range hooks {
		zlogger = zlogger.Hook(hook)
	}

	l = &Logger{zlogger, file, mask}

	// now we can log errors if something happened..
	if merr != nil {
		l.Error().
			Str("mask", mask).
			Msg("failed to parse log mask value, using debug level..")
	}

	if ferr != nil {
		l.Error().
			Str("fileName", fileName).
			Msg("failed to open log file, using stderr..")
	}

	l.Info().Msg("logger ready")
	return l
}
