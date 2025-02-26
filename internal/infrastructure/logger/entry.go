package logger

import "github.com/sirupsen/logrus"

// LogrusLoggerEntry is a wrapper for logrus.Entry
type LogrusLoggerEntry struct {
	entry *logrus.Entry
}

// Ensure LogrusLoggerEntry implements Logger interface
var _ Logger = (*LogrusLoggerEntry)(nil)

// Debugf implements Logger interface
func (e *LogrusLoggerEntry) Debugf(format string, args ...interface{}) {
	e.entry.Debugf(format, args...)
}

// Infof implements Logger interface
func (e *LogrusLoggerEntry) Infof(format string, args ...interface{}) {
	e.entry.Infof(format, args...)
}

// Warnf implements Logger interface
func (e *LogrusLoggerEntry) Warnf(format string, args ...interface{}) {
	e.entry.Warnf(format, args...)
}

// Errorf implements Logger interface
func (e *LogrusLoggerEntry) Errorf(format string, args ...interface{}) {
	e.entry.Errorf(format, args...)
}

// Fatalf implements Logger interface
func (e *LogrusLoggerEntry) Fatalf(format string, args ...interface{}) {
	e.entry.Fatalf(format, args...)
}

// WithField implements Logger interface
func (e *LogrusLoggerEntry) WithField(key string, value interface{}) Logger {
	return &LogrusLoggerEntry{
		entry: e.entry.WithField(key, value),
	}
}

// WithFields implements Logger interface
func (e *LogrusLoggerEntry) WithFields(fields map[string]interface{}) Logger {
	return &LogrusLoggerEntry{
		entry: e.entry.WithFields(logrus.Fields(fields)),
	}
}

// WithError implements Logger interface
func (e *LogrusLoggerEntry) WithError(err error) Logger {
	return &LogrusLoggerEntry{
		entry: e.entry.WithError(err),
	}
}
