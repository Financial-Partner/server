package logger

import (
	"io"

	"github.com/sirupsen/logrus"
)

// LogrusLogger is a logger implementation using logrus
type LogrusLogger struct {
	logger *logrus.Logger
}

// Ensure LogrusLogger implements Logger interface
var _ Logger = (*LogrusLogger)(nil)

// Singleton instance
var defaultLogger = NewLogrusLogger()

// NewLogrusLogger creates a new logrus logger
func NewLogrusLogger() *LogrusLogger {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetLevel(logrus.DebugLevel)
	return &LogrusLogger{logger: log}
}

// GetLogger returns the default logger
func GetLogger() Logger {
	return defaultLogger
}

// Debugf implements Logger interface
func (l *LogrusLogger) Debugf(format string, args ...interface{}) {
	l.logger.Debugf(format, args...)
}

// Infof implements Logger interface
func (l *LogrusLogger) Infof(format string, args ...interface{}) {
	l.logger.Infof(format, args...)
}

// Warnf implements Logger interface
func (l *LogrusLogger) Warnf(format string, args ...interface{}) {
	l.logger.Warnf(format, args...)
}

// Errorf implements Logger interface
func (l *LogrusLogger) Errorf(format string, args ...interface{}) {
	l.logger.Errorf(format, args...)
}

// Fatalf implements Logger interface
func (l *LogrusLogger) Fatalf(format string, args ...interface{}) {
	l.logger.Fatalf(format, args...)
}

// WithField implements Logger interface
func (l *LogrusLogger) WithField(key string, value interface{}) Logger {
	return &LogrusLoggerEntry{
		entry: l.logger.WithField(key, value),
	}
}

// WithFields implements Logger interface
func (l *LogrusLogger) WithFields(fields map[string]interface{}) Logger {
	return &LogrusLoggerEntry{
		entry: l.logger.WithFields(logrus.Fields(fields)),
	}
}

// WithError implements Logger interface
func (l *LogrusLogger) WithError(err error) Logger {
	return &LogrusLoggerEntry{
		entry: l.logger.WithError(err),
	}
}

// SetOutput sets the output destination for the logger - used in tests
func (l *LogrusLogger) SetOutput(w io.Writer) {
	l.logger.Out = w
}

// SetExitFunc sets the exit function for the logger - used in tests
func (l *LogrusLogger) SetExitFunc(exitFunc func(int)) {
	l.logger.ExitFunc = exitFunc
}
