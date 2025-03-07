package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Financial-Partner/server/internal/entities"
	"github.com/Financial-Partner/server/internal/infrastructure/logger"
	handler "github.com/Financial-Partner/server/internal/interfaces/http"
	"github.com/Financial-Partner/server/internal/interfaces/http/dto"
	"github.com/Financial-Partner/server/internal/interfaces/http/mocks"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/mock/gomock"
)

func TestLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("Invalid request format", func(t *testing.T) {
		mockUserService := mocks.NewMockUserService(ctrl)
		mockAuthService := mocks.NewMockAuthService(ctrl)
		mockLogger := logger.NewNopLogger()

		h := handler.NewHandler(mockUserService, mockAuthService, mockLogger)

		invalidBody := bytes.NewBufferString(`{invalid json`)
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/auth/login", invalidBody)

		h.Login(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var errorResp dto.ErrorResponse
		err := json.NewDecoder(w.Body).Decode(&errorResp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, errorResp.Code)
		assert.Equal(t, handler.ErrInvalidRequest, errorResp.Message)
	})

	t.Run("Authentication failed", func(t *testing.T) {
		mockUserService := mocks.NewMockUserService(ctrl)
		mockAuthService := mocks.NewMockAuthService(ctrl)
		mockLogger := logger.NewNopLogger()

		mockAuthService.EXPECT().
			LoginWithFirebase(gomock.Any(), gomock.Any()).
			Return("", "", 0, nil, errors.New("authentication failed"))

		h := handler.NewHandler(mockUserService, mockAuthService, mockLogger)

		loginReq := dto.LoginRequest{
			FirebaseToken: "valid_firebase_token",
		}
		body, _ := json.Marshal(loginReq)
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/auth/login", bytes.NewBuffer(body))

		h.Login(w, r)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var errorResp dto.ErrorResponse
		err := json.NewDecoder(w.Body).Decode(&errorResp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, errorResp.Code)
		assert.Equal(t, handler.ErrUnauthorized, errorResp.Message)
		assert.Contains(t, errorResp.Error, "authentication failed")
	})

	t.Run("Login successful", func(t *testing.T) {
		mockUserService := mocks.NewMockUserService(ctrl)
		mockAuthService := mocks.NewMockAuthService(ctrl)
		mockLogger := logger.NewNopLogger()

		objectID := primitive.NewObjectID()
		testUser := &entities.User{
			ID:    objectID,
			Email: "test@example.com",
			Name:  "Test User",
			Wallet: entities.Wallet{
				Diamonds: 100,
				Savings:  5000,
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mockAuthService.EXPECT().
			LoginWithFirebase(gomock.Any(), "valid_firebase_token").
			Return("test_access_token", "test_refresh_token", 3600, testUser, nil)

		h := handler.NewHandler(mockUserService, mockAuthService, mockLogger)

		loginReq := dto.LoginRequest{
			FirebaseToken: "valid_firebase_token",
		}
		body, _ := json.Marshal(loginReq)
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/auth/login", bytes.NewBuffer(body))

		h.Login(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var response dto.LoginResponse
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "test_access_token", response.AccessToken)
		assert.Equal(t, "test_refresh_token", response.RefreshToken)
		assert.Equal(t, 3600, response.ExpiresIn)
		assert.Equal(t, "Bearer", response.TokenType)
		assert.Equal(t, testUser.ID.Hex(), response.User.ID)
		assert.Equal(t, testUser.Email, response.User.Email)
		assert.Equal(t, testUser.Name, response.User.Name)
		assert.Equal(t, testUser.Wallet.Diamonds, response.User.Diamonds)
	})
}

func TestRefreshToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("Invalid request format", func(t *testing.T) {
		mockUserService := mocks.NewMockUserService(ctrl)
		mockAuthService := mocks.NewMockAuthService(ctrl)
		mockLogger := logger.NewNopLogger()

		h := handler.NewHandler(mockUserService, mockAuthService, mockLogger)

		invalidBody := bytes.NewBufferString(`{invalid json`)
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/auth/refresh", invalidBody)

		h.RefreshToken(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var errorResp dto.ErrorResponse
		err := json.NewDecoder(w.Body).Decode(&errorResp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, errorResp.Code)
		assert.Equal(t, handler.ErrInvalidRequest, errorResp.Message)
	})

	t.Run("Token refresh failed", func(t *testing.T) {
		mockUserService := mocks.NewMockUserService(ctrl)
		mockAuthService := mocks.NewMockAuthService(ctrl)
		mockLogger := logger.NewNopLogger()

		mockAuthService.EXPECT().
			RefreshToken(gomock.Any(), gomock.Any()).
			Return("", 0, errors.New("invalid refresh token"))

		h := handler.NewHandler(mockUserService, mockAuthService, mockLogger)

		refreshReq := dto.RefreshTokenRequest{
			RefreshToken: "invalid_refresh_token",
		}
		body, _ := json.Marshal(refreshReq)
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/auth/refresh", bytes.NewBuffer(body))

		h.RefreshToken(w, r)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var errorResp dto.ErrorResponse
		err := json.NewDecoder(w.Body).Decode(&errorResp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, errorResp.Code)
		assert.Equal(t, handler.ErrInvalidRefreshToken, errorResp.Message)
		assert.Contains(t, errorResp.Error, "invalid refresh token")
	})

	t.Run("Token refresh successful", func(t *testing.T) {
		mockUserService := mocks.NewMockUserService(ctrl)
		mockAuthService := mocks.NewMockAuthService(ctrl)
		mockLogger := logger.NewNopLogger()

		mockAuthService.EXPECT().
			RefreshToken(gomock.Any(), "valid_refresh_token").
			Return("new_access_token", 3600, nil)

		h := handler.NewHandler(mockUserService, mockAuthService, mockLogger)

		refreshReq := dto.RefreshTokenRequest{
			RefreshToken: "valid_refresh_token",
		}
		body, _ := json.Marshal(refreshReq)
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/auth/refresh", bytes.NewBuffer(body))

		h.RefreshToken(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var response dto.RefreshTokenResponse
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "new_access_token", response.AccessToken)
		assert.Equal(t, 3600, response.ExpiresIn)
		assert.Equal(t, "Bearer", response.TokenType)
	})
}

func TestLogout(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("Invalid request format", func(t *testing.T) {
		mockUserService := mocks.NewMockUserService(ctrl)
		mockAuthService := mocks.NewMockAuthService(ctrl)
		mockLogger := logger.NewNopLogger()

		h := handler.NewHandler(mockUserService, mockAuthService, mockLogger)

		invalidBody := bytes.NewBufferString(`{invalid json`)
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/auth/logout", invalidBody)

		h.Logout(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var errorResp dto.ErrorResponse
		err := json.NewDecoder(w.Body).Decode(&errorResp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, errorResp.Code)
		assert.Equal(t, handler.ErrInvalidRequest, errorResp.Message)
	})

	t.Run("Logout failed", func(t *testing.T) {
		mockUserService := mocks.NewMockUserService(ctrl)
		mockAuthService := mocks.NewMockAuthService(ctrl)
		mockLogger := logger.NewNopLogger()

		mockAuthService.EXPECT().
			Logout(gomock.Any(), gomock.Any()).
			Return(errors.New("logout failed"))

		h := handler.NewHandler(mockUserService, mockAuthService, mockLogger)

		logoutReq := dto.LogoutRequest{
			RefreshToken: "refresh_token",
		}
		body, _ := json.Marshal(logoutReq)
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/auth/logout", bytes.NewBuffer(body))

		h.Logout(w, r)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var errorResp dto.ErrorResponse
		err := json.NewDecoder(w.Body).Decode(&errorResp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, errorResp.Code)
		assert.Equal(t, handler.ErrFailedToLogout, errorResp.Message)
		assert.Contains(t, errorResp.Error, "logout failed")
	})

	t.Run("Logout successful", func(t *testing.T) {
		mockUserService := mocks.NewMockUserService(ctrl)
		mockAuthService := mocks.NewMockAuthService(ctrl)
		mockLogger := logger.NewNopLogger()

		mockAuthService.EXPECT().
			Logout(gomock.Any(), "refresh_token").
			Return(nil)

		h := handler.NewHandler(mockUserService, mockAuthService, mockLogger)

		logoutReq := dto.LogoutRequest{
			RefreshToken: "refresh_token",
		}
		body, _ := json.Marshal(logoutReq)
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/auth/logout", bytes.NewBuffer(body))

		h.Logout(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var response dto.LogoutResponse
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.True(t, response.Success)
		assert.Equal(t, "Logout successfully", response.Message)
	})
}
