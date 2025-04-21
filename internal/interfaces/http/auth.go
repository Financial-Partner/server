package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Financial-Partner/server/internal/entities"
	"github.com/Financial-Partner/server/internal/interfaces/http/dto"
	httperror "github.com/Financial-Partner/server/internal/interfaces/http/error"
	respond "github.com/Financial-Partner/server/internal/interfaces/http/respond"
)

//go:generate mockgen -source=auth.go -destination=auth_mock.go -package=handler

type AuthService interface {
	LoginWithFirebase(ctx context.Context, firebaseToken string) (accessToken, refreshToken string, expiresIn int, userInfo *entities.User, err error)
	RefreshToken(ctx context.Context, refreshToken string) (newAccessToken, newRefreshToken string, expiresIn int, err error)
	Logout(ctx context.Context, refreshToken string) error
}

// Login Login with Firebase
// @Summary Login with Firebase
// @Description Login with Firebase, get Access Token and Refresh Token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login request"
// @Success 200 {object} dto.LoginResponse "Login successfully"
// @Failure 400 {object} dto.ErrorResponse "Invalid request format"
// @Failure 401 {object} dto.ErrorResponse "Authentication failed"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /auth/login [post]
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.WithError(err).Warnf("Invalid request format")
		respond.WithError(w, r, h.log, err, httperror.ErrInvalidRequest, http.StatusBadRequest)
		return
	}

	accessToken, refreshToken, expiresIn, userInfo, err := h.authService.LoginWithFirebase(r.Context(), req.FirebaseToken)
	if err != nil {
		h.log.WithError(err).Errorf("Login failed")
		respond.WithError(w, r, h.log, err, httperror.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	response := dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
		TokenType:    "Bearer",
		User: dto.UserResponse{
			ID:        userInfo.ID.Hex(),
			Email:     userInfo.Email,
			Name:      userInfo.Name,
			Diamonds:  userInfo.Wallet.Diamonds,
			CreatedAt: userInfo.CreatedAt.Format(time.RFC3339),
		},
	}

	respond.WithJSON(w, r, response, http.StatusOK)
}

// RefreshToken Refresh Access Token
// @Summary Refresh Access Token
// @Description Use Refresh Token to get a new Access Token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RefreshTokenRequest true "Token refresh request"
// @Success 200 {object} dto.RefreshTokenResponse "Token refresh successfully"
// @Failure 400 {object} dto.ErrorResponse "Invalid request format"
// @Failure 401 {object} dto.ErrorResponse "Invalid refresh token"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /auth/refresh [post]
func (h *Handler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req dto.RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.WithError(err).Warnf("Invalid request format")
		respond.WithError(w, r, h.log, err, httperror.ErrInvalidRequest, http.StatusBadRequest)
		return
	}

	newAccessToken, newRefreshToken, expiresIn, err := h.authService.RefreshToken(r.Context(), req.RefreshToken)
	if err != nil {
		h.log.WithError(err).Errorf("Token refresh failed")
		respond.WithError(w, r, h.log, err, httperror.ErrInvalidRefreshToken, http.StatusUnauthorized)
		return
	}

	response := dto.RefreshTokenResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    expiresIn,
		TokenType:    "Bearer",
	}

	respond.WithJSON(w, r, response, http.StatusOK)
}

// Logout Logout
// @Summary User logout
// @Description Invalidate the current Refresh Token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LogoutRequest true "Logout request"
// @Success 200 {object} dto.LogoutResponse "Logout successfully"
// @Failure 400 {object} dto.ErrorResponse "Invalid request format"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /auth/logout [post]
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	var req dto.LogoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.WithError(err).Warnf("Invalid request format")
		respond.WithError(w, r, h.log, err, httperror.ErrInvalidRequest, http.StatusBadRequest)
		return
	}

	err := h.authService.Logout(r.Context(), req.RefreshToken)
	if err != nil {
		h.log.WithError(err).Errorf("Logout failed")
		respond.WithError(w, r, h.log, err, httperror.ErrFailedToLogout, http.StatusInternalServerError)
		return
	}

	response := dto.LogoutResponse{
		Success: true,
		Message: "Logout successfully",
	}

	respond.WithJSON(w, r, response, http.StatusOK)
}
