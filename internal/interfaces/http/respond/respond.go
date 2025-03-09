package responde

import (
	"encoding/json"
	"net/http"

	"github.com/Financial-Partner/server/internal/infrastructure/logger"
	"github.com/Financial-Partner/server/internal/interfaces/http/dto"
)

func WithError(w http.ResponseWriter, r *http.Request, log logger.Logger, err error, message string, statusCode int) {
	logField := log.WithError(err)

	if statusCode >= 500 {
		logField.Errorf("Server error: %s", message)
	} else if statusCode >= 400 {
		logField.Warnf("Client error: %s", message)
	}
	// Should not get other status code

	errResp := dto.ErrorResponse{
		Code:    statusCode,
		Message: message,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(errResp)
}

func WithJSON(w http.ResponseWriter, r *http.Request, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(data)
}
