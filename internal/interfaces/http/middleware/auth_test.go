package middleware_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/Financial-Partner/server/internal/contextutil"
	"github.com/Financial-Partner/server/internal/infrastructure/auth"
	"github.com/Financial-Partner/server/internal/infrastructure/logger"
	"github.com/Financial-Partner/server/internal/interfaces/http/middleware"
)

type contextKey string

func TestNewAuthMiddleware(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockJWTValidator := middleware.NewMockJWTValidator(ctrl)
	mockLogger := logger.NewNopLogger()

	middleware := middleware.NewAuthMiddleware(mockJWTValidator, mockLogger)

	assert.NotNil(t, middleware)
}

func TestAuthMiddlewareAuthenticate(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockJWTValidator := middleware.NewMockJWTValidator(ctrl)
		mockLogger := logger.NewNopLogger()

		token := &auth.Claims{
			Email: "test@example.com",
		}
		mockJWTValidator.EXPECT().ValidateToken("valid-token").Return(token, nil)

		middleware := middleware.NewAuthMiddleware(mockJWTValidator, mockLogger)

		var capturedEmail string
		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			email, ok := contextutil.GetUserEmail(r.Context())
			assert.True(t, ok)
			capturedEmail = email
			w.WriteHeader(http.StatusOK)
		})

		req := httptest.NewRequest("GET", "/api/test", nil)
		req.Header.Set("Authorization", "Bearer valid-token")
		w := httptest.NewRecorder()

		handler := middleware.Authenticate(nextHandler)
		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "test@example.com", capturedEmail)
	})

	t.Run("NoToken", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockJWTValidator := middleware.NewMockJWTValidator(ctrl)
		mockLogger := logger.NewNopLogger()

		middleware := middleware.NewAuthMiddleware(mockJWTValidator, mockLogger)

		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Fatal("Next handler should not be called")
		})

		req := httptest.NewRequest("GET", "/api/test", nil)
		w := httptest.NewRecorder()

		handler := middleware.Authenticate(nextHandler)
		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Authorization token required")
	})

	t.Run("InvalidToken", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockJWTValidator := middleware.NewMockJWTValidator(ctrl)
		mockLogger := logger.NewNopLogger()

		tokenError := errors.New("invalid token")
		mockJWTValidator.EXPECT().ValidateToken("invalid-token").Return(nil, tokenError)

		middleware := middleware.NewAuthMiddleware(mockJWTValidator, mockLogger)

		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Fatal("Next handler should not be called")
		})

		req := httptest.NewRequest("GET", "/api/test", nil)
		req.Header.Set("Authorization", "Bearer invalid-token")
		w := httptest.NewRecorder()

		handler := middleware.Authenticate(nextHandler)
		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid authorization token")
	})

	t.Run("NoEmailInToken", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockJWTValidator := middleware.NewMockJWTValidator(ctrl)
		mockLogger := logger.NewNopLogger()

		token := &auth.Claims{
			Email: "",
		}
		mockJWTValidator.EXPECT().ValidateToken("no-email-token").Return(token, nil)

		middleware := middleware.NewAuthMiddleware(mockJWTValidator, mockLogger)

		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Fatal("Next handler should not be called")
		})

		req := httptest.NewRequest("GET", "/api/test", nil)
		req.Header.Set("Authorization", "Bearer no-email-token")
		w := httptest.NewRecorder()

		handler := middleware.Authenticate(nextHandler)
		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Unable to retrieve user email")
	})

	t.Run("BearerPrefix", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockJWTValidator := middleware.NewMockJWTValidator(ctrl)
		mockLogger := logger.NewNopLogger()

		token := &auth.Claims{
			Email: "test@example.com",
		}
		mockJWTValidator.EXPECT().ValidateToken("valid-token").Return(token, nil)
		mockLogger.WithField("email", "test@example.com").Infof("User authenticated successfully", gomock.Any())

		middleware := middleware.NewAuthMiddleware(mockJWTValidator, mockLogger)

		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		req := httptest.NewRequest("GET", "/api/test", nil)
		req.Header.Set("Authorization", "Bearer valid-token")
		w := httptest.NewRecorder()

		handler := middleware.Authenticate(nextHandler)
		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("ContextPropagation", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockJWTValidator := middleware.NewMockJWTValidator(ctrl)
		mockLogger := logger.NewNopLogger()

		token := &auth.Claims{
			Email: "test@example.com",
		}
		mockJWTValidator.EXPECT().ValidateToken("valid-token").Return(token, nil)
		mockLogger.WithField("email", "test@example.com").Infof("User authenticated successfully", gomock.Any())

		middleware := middleware.NewAuthMiddleware(mockJWTValidator, mockLogger)

		var capturedCtx context.Context
		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			capturedCtx = r.Context()
			w.WriteHeader(http.StatusOK)
		})

		originalCtx := context.WithValue(context.Background(), contextKey("test-key"), "test-value")
		req := httptest.NewRequest("GET", "/api/test", nil).WithContext(originalCtx)
		req.Header.Set("Authorization", "Bearer valid-token")
		w := httptest.NewRecorder()

		handler := middleware.Authenticate(nextHandler)
		handler.ServeHTTP(w, req)

		require.NotNil(t, capturedCtx)
		assert.Equal(t, "test-value", capturedCtx.Value(contextKey("test-key")))

		email, ok := contextutil.GetUserEmail(capturedCtx)
		assert.True(t, ok)
		assert.Equal(t, "test@example.com", email)
	})
}
