package subscribe

import (
	"os"

	"github.com/kyokomi/emoji"
	"github.com/tendermint/tendermint/libs/log"
)

// Logger is used for logging when running the websocket
type Logger struct {
	logger log.Logger
}

// NewLogger is the constructor of the Logger struct
func NewLogger(level log.Option) *Logger {
	return &Logger{logger: log.NewFilter(log.NewTMLogger(os.Stdout), level)}
}

// Debug is logger debug option
func (l *Logger) Debug(format string, args ...interface{}) {
	l.logger.Debug(emoji.Sprintf(format, args...))
}

// Info is logger information option
func (l *Logger) Info(format string, args ...interface{}) {
	l.logger.Info(emoji.Sprintf(format, args...))
}

// Error is logger error option
func (l *Logger) Error(format string, args ...interface{}) {
	l.logger.Error(emoji.Sprintf(format, args...))
}

// With does something
func (l *Logger) With(keyvals ...interface{}) *Logger {
	return &Logger{logger: l.logger.With(keyvals...)}
}
