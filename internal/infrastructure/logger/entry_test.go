package logger_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"testing"

	"github.com/Financial-Partner/server/internal/infrastructure/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLogrusLoggerEntry_LogMethods(t *testing.T) {
	log := logger.NewLogrusLogger()
	buffer := &bytes.Buffer{}
	log.SetOutput(buffer)

	entry := log.WithField("test_key", "test_value")

	testCases := []struct {
		name     string
		logFunc  func(string, ...interface{})
		level    string
		message  string
		args     []interface{}
		expected string
	}{
		{
			name:    "Debugf",
			logFunc: entry.Debugf,
			level:   "debug",
			message: "test debug message %s",
			args:    []interface{}{"arg"},
		},
		{
			name:    "Infof",
			logFunc: entry.Infof,
			level:   "info",
			message: "test info message %s",
			args:    []interface{}{"arg"},
		},
		{
			name:    "Warnf",
			logFunc: entry.Warnf,
			level:   "warning",
			message: "test warning message %s",
			args:    []interface{}{"arg"},
		},
		{
			name:    "Errorf",
			logFunc: entry.Errorf,
			level:   "error",
			message: "test error message %s",
			args:    []interface{}{"arg"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			buffer.Reset()
			tc.logFunc(tc.message, tc.args...)

			var logEntry map[string]interface{}
			err := json.Unmarshal(buffer.Bytes(), &logEntry)
			require.NoError(t, err)

			assert.Equal(t, tc.level, logEntry["level"])
			expectedMsg := tc.message
			if len(tc.args) > 0 {
				expectedMsg = "test " + tc.level + " message arg"
			}
			assert.Equal(t, expectedMsg, logEntry["msg"])
			assert.Equal(t, "test_value", logEntry["test_key"])
		})
	}
}

func TestLogrusLoggerEntryWithField(t *testing.T) {
	log := logger.NewLogrusLogger()
	buffer := &bytes.Buffer{}
	log.SetOutput(buffer)

	entry := log.WithField("test_key", "test_value")

	fieldLogger := entry.WithField("key", "value")
	assert.NotNil(t, fieldLogger)

	fieldLogger.Infof("test message")

	var logEntry map[string]interface{}
	err := json.Unmarshal(buffer.Bytes(), &logEntry)
	require.NoError(t, err)

	assert.Equal(t, "info", logEntry["level"])
	assert.Equal(t, "test message", logEntry["msg"])
	assert.Equal(t, "test_value", logEntry["test_key"])
	assert.Equal(t, "value", logEntry["key"])
}

func TestLogrusLoggerEntryWithFields(t *testing.T) {
	log := logger.NewLogrusLogger()
	buffer := &bytes.Buffer{}
	log.SetOutput(buffer)

	entry := log.WithField("test_key", "test_value")

	fields := map[string]interface{}{
		"key1": "value1",
		"key2": 42,
	}
	fieldsLogger := entry.WithFields(fields)
	assert.NotNil(t, fieldsLogger)

	fieldsLogger.Infof("test message")

	var logEntry map[string]interface{}
	err := json.Unmarshal(buffer.Bytes(), &logEntry)
	require.NoError(t, err)

	assert.Equal(t, "info", logEntry["level"])
	assert.Equal(t, "test message", logEntry["msg"])
	assert.Equal(t, "test_value", logEntry["test_key"])
	assert.Equal(t, "value1", logEntry["key1"])
	assert.Equal(t, float64(42), logEntry["key2"])
}

func TestLogrusLoggerEntryWithError(t *testing.T) {
	log := logger.NewLogrusLogger()
	buffer := &bytes.Buffer{}
	log.SetOutput(buffer)

	entry := log.WithField("test_key", "test_value")

	testErr := errors.New("test error")
	errorLogger := entry.WithError(testErr)
	assert.NotNil(t, errorLogger)

	errorLogger.Infof("test message")

	var logEntry map[string]interface{}
	err := json.Unmarshal(buffer.Bytes(), &logEntry)
	require.NoError(t, err)

	assert.Equal(t, "info", logEntry["level"])
	assert.Equal(t, "test message", logEntry["msg"])
	assert.Equal(t, "test_value", logEntry["test_key"])
	assert.Equal(t, "test error", logEntry["error"])
}

func TestLogrusLoggerEntryFatalf(t *testing.T) {
	log := logger.NewLogrusLogger()
	buffer := &bytes.Buffer{}
	log.SetOutput(buffer)

	exitCalled := false
	log.SetExitFunc(func(int) { exitCalled = true })

	entry := log.WithField("test_key", "test_value")

	entry.Fatalf("fatal error %s", "test")

	var logEntry map[string]interface{}
	err := json.Unmarshal(buffer.Bytes(), &logEntry)
	require.NoError(t, err)

	assert.Equal(t, "fatal", logEntry["level"])
	assert.Equal(t, "fatal error test", logEntry["msg"])
	assert.Equal(t, "test_value", logEntry["test_key"])
	assert.True(t, exitCalled, "Exit function should have been called")
}
