package logger

import (
	"os"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

var Log Logger
var once sync.Once

type Logger struct {
	zerolog.Logger
}

type LogVals struct {
	Endpoint        string
	ResponseCodeErr *error
}

func (l *Logger) LogDebug(msg string) {
	l.Debug().Msg(msg)
}

func (l *Logger) LogDebugf(format string, v ...interface{}) {
	l.Debug().Msgf(format, v...)
}

func (l *Logger) LogInfo(msg string) {
	l.Info().Msg(msg)
}

func (l *Logger) LogInfof(format string, v ...interface{}) {
	l.Info().Msgf(format, v...)
}

func (l *Logger) LogWarn(msg string) {
	l.Warn().Msg(msg)
}

func (l *Logger) LogWarnf(format string, v ...interface{}) {
	l.Warn().Msgf(format, v...)
}

func (l *Logger) LogError(msg string) {
	l.Error().Msg(msg)
}

func (l *Logger) LogErrorf(format string, v ...interface{}) {
	l.Error().Msgf(format, v...)
}

func (l *Logger) LogFatal(msg string) {
	l.Fatal().Msg(msg)
}

func (l *Logger) LogFatalf(format string, v ...interface{}) {
	l.Fatal().Msgf(format, v...)
}

func getLog(serviceName string, isPretty bool, logLevel string) zerolog.Logger {
	logger := zerolog.New(os.Stdout).With().Str("module", serviceName).Timestamp()

	if isPretty {
		logger = zerolog.New(os.Stdout).Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}).With().Str("module", serviceName).Timestamp()
	}

	if logLevel != "" {
		level, err := zerolog.ParseLevel(logLevel)
		if err != nil {
			// default log level is debug
			return logger.Logger()
		}
		zerolog.SetGlobalLevel(level)
	}

	return logger.Logger()
}

func New(serviceName string, isPretty bool, logLevel string) *Logger {
	once.Do(func() {
		Log = Logger{getLog(serviceName, isPretty, logLevel)}
	})

	return &Log
}
