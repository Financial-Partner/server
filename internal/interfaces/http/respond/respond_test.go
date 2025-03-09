package responde_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Financial-Partner/server/internal/infrastructure/logger"
	"github.com/Financial-Partner/server/internal/interfaces/http/dto"
	respond "github.com/Financial-Partner/server/internal/interfaces/http/respond"
	"github.com/stretchr/testify/assert"
)

func TestWithError(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		message    string
		err        error
	}{
		{
			name:       "client error",
			statusCode: http.StatusBadRequest,
			message:    "invalid request",
			err:        errors.New("validation error"),
		},
		{
			name:       "server error",
			statusCode: http.StatusInternalServerError,
			message:    "server error",
			err:        errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/test", nil)
			mockLoggger := logger.NewNopLogger()

			respond.WithError(w, r, mockLoggger, tt.err, tt.message, tt.statusCode)

			assert.Equal(t, tt.statusCode, w.Code)
			assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

			var response dto.ErrorResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tt.statusCode, response.Code)
			assert.Equal(t, tt.message, response.Message)
		})
	}
}

func TestWithJSON(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		data       interface{}
	}{
		{
			name:       "success response",
			statusCode: http.StatusOK,
			data:       map[string]interface{}{"success": true, "data": "test data"},
		},
		{
			name:       "created response",
			statusCode: http.StatusCreated,
			data:       map[string]interface{}{"id": 123, "name": "test resource"},
		},
		{
			name:       "empty response",
			statusCode: http.StatusNoContent,
			data:       nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/test", nil)

			respond.WithJSON(w, r, tt.data, tt.statusCode)

			assert.Equal(t, tt.statusCode, w.Code)
			assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

			if tt.data != nil {
				expectedJSON, err := json.Marshal(tt.data)
				assert.NoError(t, err)
				expectedJSON = append(expectedJSON, '\n')
				assert.Equal(t, expectedJSON, w.Body.Bytes())
			}
		})
	}
}
