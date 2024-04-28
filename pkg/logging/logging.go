package logging

import (
	"github.com/rs/zerolog/pkgerrors"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

/*
panic (zerolog.PanicLevel, 5)
fatal (zerolog.FatalLevel, 4)
error (zerolog.ErrorLevel, 3)
warn (zerolog.WarnLevel, 2)
info (zerolog.InfoLevel, 1)
debug (zerolog.DebugLevel, 0)
trace (zerolog.TraceLevel, -1)
*/
type ILogger interface {
	LogError(error, string)
	LogInfo(string)
	LogTrace(string)
	HadError() bool
}

type Logger struct {
	hadError bool
}

// The main logging function for Error
func (h *Logger) LogError(err error, extraMsg string) {
	log.Error().Stack().Err(err).Msg(extraMsg)
}

// The main logging function for Info
func (h *Logger) LogInfo(msg string) {
	log.Info().Msg(msg)
}

// The main logging function for Trace
func (h *Logger) LogTrace(msg string) {
	log.Trace().Msg(msg)
}

// Error testing
func (h *Logger) HadError() bool {
	return h.hadError
}

func NewLogger() ILogger {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	return &Logger{
		hadError: false,
	}
}
