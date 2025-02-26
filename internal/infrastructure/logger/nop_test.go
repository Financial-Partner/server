package logger_test

import (
	"testing"

	"github.com/Financial-Partner/server/internal/infrastructure/logger"
	"github.com/stretchr/testify/assert"
)

func TestNopLogger(t *testing.T) {
	log := logger.NewNopLogger()
	assert.NotNil(t, log)

	log.Debugf("debug message")
	log.Infof("info message")
	log.Warnf("warn message")
	log.Errorf("error message")
	log.Fatalf("fatal message")

	withFieldLogger := log.WithField("key", "value")
	assert.NotNil(t, withFieldLogger)
	withFieldLogger.Infof("with field message")

	withFieldsLogger := log.WithFields(map[string]interface{}{"key": "value"})
	assert.NotNil(t, withFieldsLogger)
	withFieldsLogger.Infof("with fields message")

	withErrorLogger := log.WithError(assert.AnError)
	assert.NotNil(t, withErrorLogger)
	withErrorLogger.Infof("with error message")
}
