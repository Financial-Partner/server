package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Financial-Partner/server/internal/infrastructure/logger"
	"github.com/Financial-Partner/server/internal/interfaces/http/dto"
)

const (
	ErrInvalidRequest      = "Invalid request format"
	ErrUnauthorized        = "Unauthorized"
	ErrEmailNotFound       = "Email not found"
	ErrEmailMismatch       = "Email mismatch"
	ErrInternalServer      = "Internal server error"
	ErrUserNotFound        = "User not found"
	ErrInvalidRefreshToken = "Invalid refresh token"
	ErrFailedToCreateUser  = "Failed to create user"
	ErrFailedToUpdateUser  = "Failed to update user"
	ErrFailedToGetUser     = "Failed to get user"
	ErrFailedToLogout      = "Failed to logout"
)

func (h *Handler) RespondWithError(w http.ResponseWriter, r *http.Request, log logger.Logger, err error, message string, statusCode int) {
	logField := log.WithError(err)

	if statusCode >= 500 {
		logField.Errorf("Server error: %s", message)
	} else if statusCode >= 400 {
		logField.Warnf("Client error: %s", message)
	} else {
		logField.Infof("Info: %s", message)
	}

	errResp := dto.ErrorResponse{
		Code:    statusCode,
		Message: message,
	}

	if h.IsDevelopment() && err != nil {
		errResp.Error = err.Error()
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(errResp)
}

func (h *Handler) RespondWithJSON(w http.ResponseWriter, r *http.Request, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(data)
}

func (h *Handler) IsDevelopment() bool {
	// TODO: parse from config
	return true
}
