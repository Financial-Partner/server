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

func TestLogrusLogger(t *testing.T) {
	t.Run("NewLogrusLogger", func(t *testing.T) {
		log := logger.NewLogrusLogger()
		assert.NotNil(t, log)

		buffer := &bytes.Buffer{}
		log.SetOutput(buffer)

		log.Debugf("test debug")

		var logEntry map[string]interface{}
		err := json.Unmarshal(buffer.Bytes(), &logEntry)
		require.NoError(t, err)

		assert.Equal(t, "debug", logEntry["level"])
	})

	t.Run("GetLogger", func(t *testing.T) {
		log := logger.GetLogger()
		assert.NotNil(t, log)

		buffer := &bytes.Buffer{}
		log.(*logger.LogrusLogger).SetOutput(buffer)

		log.Infof("test info")

		var logEntry map[string]interface{}
		err := json.Unmarshal(buffer.Bytes(), &logEntry)
		require.NoError(t, err)

		assert.Equal(t, "info", logEntry["level"])
	})
}

func TestLogrusLoggerLogMethods(t *testing.T) {
	log := logger.NewLogrusLogger()
	buffer := &bytes.Buffer{}
	log.SetOutput(buffer)

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
			logFunc: log.Debugf,
			level:   "debug",
			message: "test debug message %s",
			args:    []interface{}{"arg"},
		},
		{
			name:    "Infof",
			logFunc: log.Infof,
			level:   "info",
			message: "test info message %s",
			args:    []interface{}{"arg"},
		},
		{
			name:    "Warnf",
			logFunc: log.Warnf,
			level:   "warning",
			message: "test warning message %s",
			args:    []interface{}{"arg"},
		},
		{
			name:    "Errorf",
			logFunc: log.Errorf,
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
		})
	}
}

func TestLogrusLoggerWithField(t *testing.T) {
	log := logger.NewLogrusLogger()
	buffer := &bytes.Buffer{}
	log.SetOutput(buffer)

	fieldLogger := log.WithField("key", "value")
	assert.NotNil(t, fieldLogger)

	fieldLogger.Infof("test message")

	var logEntry map[string]interface{}
	err := json.Unmarshal(buffer.Bytes(), &logEntry)
	require.NoError(t, err)

	assert.Equal(t, "info", logEntry["level"])
	assert.Equal(t, "test message", logEntry["msg"])
	assert.Equal(t, "value", logEntry["key"])
}

func TestLogrusLoggerWithFields(t *testing.T) {
	log := logger.NewLogrusLogger()
	buffer := &bytes.Buffer{}
	log.SetOutput(buffer)

	fields := map[string]interface{}{
		"key1": "value1",
		"key2": 42,
	}
	fieldsLogger := log.WithFields(fields)
	assert.NotNil(t, fieldsLogger)

	fieldsLogger.Infof("test message")

	var logEntry map[string]interface{}
	err := json.Unmarshal(buffer.Bytes(), &logEntry)
	require.NoError(t, err)

	assert.Equal(t, "info", logEntry["level"])
	assert.Equal(t, "test message", logEntry["msg"])
	assert.Equal(t, "value1", logEntry["key1"])
	assert.Equal(t, float64(42), logEntry["key2"])
}

func TestLogrusLoggerWithError(t *testing.T) {
	log := logger.NewLogrusLogger()
	buffer := &bytes.Buffer{}
	log.SetOutput(buffer)

	testErr := errors.New("test error")
	errorLogger := log.WithError(testErr)
	assert.NotNil(t, errorLogger)

	errorLogger.Infof("test message")

	var logEntry map[string]interface{}
	err := json.Unmarshal(buffer.Bytes(), &logEntry)
	require.NoError(t, err)

	assert.Equal(t, "info", logEntry["level"])
	assert.Equal(t, "test message", logEntry["msg"])
	assert.Equal(t, "test error", logEntry["error"])
}

func TestLogrusLoggerFatalf(t *testing.T) {
	log := logger.NewLogrusLogger()
	buffer := &bytes.Buffer{}
	log.SetOutput(buffer)

	exitCalled := false
	log.SetExitFunc(func(int) { exitCalled = true })

	log.Fatalf("fatal error %s", "test")

	var logEntry map[string]interface{}
	err := json.Unmarshal(buffer.Bytes(), &logEntry)
	require.NoError(t, err)

	assert.Equal(t, "fatal", logEntry["level"])
	assert.Equal(t, "fatal error test", logEntry["msg"])
	assert.True(t, exitCalled, "Exit function should have been called")
}
