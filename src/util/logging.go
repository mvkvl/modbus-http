package util

import (
	"errors"
	"fmt"
	"github.com/rs/zerolog"
	"io"
	"os"
	"strings"
	"time"
)

type Logger interface {
	Trace(message string, args ...interface{})
	Debug(message string, args ...interface{})
	Info(message string, args ...interface{})
	Warning(message string, args ...interface{})
	Error(message string, args ...interface{})
	Fatal(message string, args ...interface{})
	Panic(message string, args ...interface{})
	IsTraceEnabled() bool
	IsDebugEnabled() bool
	IsInfoEnabled() bool
	IsWarningEnabled() bool
	IsErrorEnabled() bool
	IsFatalEnabled() bool
	IsPanicEnabled() bool
}

type loggerImpl struct {
	lg *zerolog.Logger
}

func (l *loggerImpl) Trace(message string, args ...interface{}) {
	if l.lg.GetLevel() == zerolog.TraceLevel {
		l.lg.Trace().Msg(fmt.Sprintf(message, args...))
	}
}
func (l *loggerImpl) Debug(message string, args ...interface{}) {
	if l.lg.GetLevel() <= zerolog.DebugLevel {
		l.lg.Debug().Msg(fmt.Sprintf(message, args...))
	}
}
func (l *loggerImpl) Info(message string, args ...interface{}) {
	if l.lg.GetLevel() <= zerolog.InfoLevel {
		l.lg.Info().Msg(fmt.Sprintf(message, args...))
	}
}
func (l *loggerImpl) Warning(message string, args ...interface{}) {
	if l.lg.GetLevel() <= zerolog.WarnLevel {
		l.lg.Warn().Msg(fmt.Sprintf(message, args...))
	}
}
func (l *loggerImpl) Error(message string, args ...interface{}) {
	if l.lg.GetLevel() <= zerolog.ErrorLevel {
		l.lg.Error().Msg(fmt.Sprintf(message, args...))
	}
}
func (l *loggerImpl) Fatal(message string, args ...interface{}) {
	if l.lg.GetLevel() <= zerolog.FatalLevel {
		l.lg.Fatal().Msg(fmt.Sprintf(message, args...))
	}
}
func (l *loggerImpl) Panic(message string, args ...interface{}) {
	if l.lg.GetLevel() <= zerolog.PanicLevel {
		l.lg.Panic().Msg(fmt.Sprintf(message, args...))
	}
}

func (l *loggerImpl) IsTraceEnabled() bool {
	return l.lg.GetLevel() == zerolog.TraceLevel
}
func (l *loggerImpl) IsDebugEnabled() bool {
	return l.lg.GetLevel() <= zerolog.DebugLevel
}
func (l *loggerImpl) IsInfoEnabled() bool {
	return l.lg.GetLevel() <= zerolog.InfoLevel
}
func (l *loggerImpl) IsWarningEnabled() bool {
	return l.lg.GetLevel() <= zerolog.WarnLevel
}
func (l *loggerImpl) IsErrorEnabled() bool {
	return l.lg.GetLevel() <= zerolog.ErrorLevel
}
func (l *loggerImpl) IsFatalEnabled() bool {
	return l.lg.GetLevel() <= zerolog.FatalLevel
}
func (l *loggerImpl) IsPanicEnabled() bool {
	return l.lg.GetLevel() <= zerolog.PanicLevel
}

func init() {
	loggerFactory = loggerFactoryImpl{
		loggers: make(map[string]Logger),
	}
}

var loggerFactory loggerFactoryImpl

type loggerFactoryImpl struct {
	loggers map[string]Logger
}

func GetLogger(id string) Logger {
	v, ok := loggerFactory.loggers[id]
	if ok {
		return v
	}
	l := newLogger(id)
	loggerFactory.loggers[id] = l
	return l
}

func newLogger(id string) Logger {
	var w io.Writer
	if os.Getenv("GO_ENV") != "dev" {
		w = os.Stdout
	} else {
		w = configureConsoleWriter(id)
	}
	logger := zerolog.
		New(w).
		Level(getLoggingLevel(id)).
		With().Str("logger", id).Timestamp(). //Caller().
		Logger()
	result := &loggerImpl{
		lg: &logger,
	}
	return result
}

func configureConsoleWriter(id string) io.Writer {
	return zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.RFC3339,
		FormatMessage: func(i interface{}) string {
			return fmt.Sprintf("[%s] %s", id, i)
		},
		FieldsExclude: []string{
			"logger",
		},
		PartsExclude: []string{
			zerolog.CallerFieldName,
			"logger",
		},
	}
}
func getLoggingLevel(id string) zerolog.Level {
	/*
		LevelOffValue = "off" // special value to turn logs off
		LevelTraceValue = "trace"
		LevelDebugValue = "debug"
		LevelInfoValue = "info"
		LevelWarnValue = "warn"
		LevelErrorValue = "error"
		LevelFatalValue = "fatal"
		LevelPanicValue = "panic"
	*/
	key := getConfiguredLevel(id)
	lvl, err := stringToLevel(key)
	if err != nil {
		lvl = zerolog.InfoLevel
	}
	return lvl
}
func stringToLevel(key string) (zerolog.Level, error) {
	if strings.ToUpper(key) == "OFF" {
		return zerolog.NoLevel, nil
	}
	if key == "" {
		return zerolog.NoLevel, errors.New("empty config")
	}
	return zerolog.ParseLevel(key)
}
func getConfiguredLevel(id string) string {
	key := os.Getenv("LOGGING_LEVEL_" + strings.ToUpper(strings.ReplaceAll(id, "-", "_")))
	if strings.TrimSpace(key) == "" {
		key = os.Getenv("LOGGING_LEVEL_ROOT")
	}
	return key
}
