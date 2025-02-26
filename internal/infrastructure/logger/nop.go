package logger

// NopLogger is a no-operation logger that discards all log messages
// Useful for testing when you don't want to see any logs
type NopLogger struct{}

// Ensure NopLogger implements Logger interface
var _ Logger = (*NopLogger)(nil)

// NewNopLogger creates a new no-operation logger
func NewNopLogger() Logger {
	return &NopLogger{}
}

// Debugf implements Logger interface - does nothing
func (l *NopLogger) Debugf(format string, args ...interface{}) {}

// Infof implements Logger interface - does nothing
func (l *NopLogger) Infof(format string, args ...interface{}) {}

// Warnf implements Logger interface - does nothing
func (l *NopLogger) Warnf(format string, args ...interface{}) {}

// Errorf implements Logger interface - does nothing
func (l *NopLogger) Errorf(format string, args ...interface{}) {}

// Fatalf implements Logger interface - does nothing
func (l *NopLogger) Fatalf(format string, args ...interface{}) {}

// WithField implements Logger interface - returns the same nop logger
func (l *NopLogger) WithField(key string, value interface{}) Logger {
	return l
}

// WithFields implements Logger interface - returns the same nop logger
func (l *NopLogger) WithFields(fields map[string]interface{}) Logger {
	return l
}

// WithError implements Logger interface - returns the same nop logger
func (l *NopLogger) WithError(err error) Logger {
	return l
}
