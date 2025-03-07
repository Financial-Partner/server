package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/mock/gomock"

	"github.com/Financial-Partner/server/internal/contextutil"
	"github.com/Financial-Partner/server/internal/entities"
	"github.com/Financial-Partner/server/internal/infrastructure/logger"
	handler "github.com/Financial-Partner/server/internal/interfaces/http"
	"github.com/Financial-Partner/server/internal/interfaces/http/dto"
	"github.com/Financial-Partner/server/internal/interfaces/http/mocks"
)

func TestCreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("Invalid request format", func(t *testing.T) {
		mockUserService := mocks.NewMockUserService(ctrl)
		mockAuthService := mocks.NewMockAuthService(ctrl)
		mockLogger := logger.NewNopLogger()

		h := handler.NewHandler(mockUserService, mockAuthService, mockLogger)

		invalidBody := bytes.NewBufferString(`{invalid json`)
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/users", invalidBody)

		h.CreateUser(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var errorResp dto.ErrorResponse
		err := json.NewDecoder(w.Body).Decode(&errorResp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, errorResp.Code)
		assert.Equal(t, handler.ErrInvalidRequest, errorResp.Message)
	})

	t.Run("Email mismatch", func(t *testing.T) {
		mockUserService := mocks.NewMockUserService(ctrl)
		mockAuthService := mocks.NewMockAuthService(ctrl)
		mockLogger := logger.NewNopLogger()

		h := handler.NewHandler(mockUserService, mockAuthService, mockLogger)

		createReq := dto.CreateUserRequest{
			Email: "user@example.com",
			Name:  "Test User",
		}
		body, _ := json.Marshal(createReq)

		ctx := context.Background()
		ctx = context.WithValue(ctx, contextutil.UserEmailKey, "different@example.com")

		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/users", bytes.NewBuffer(body)).WithContext(ctx)

		h.CreateUser(w, r)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var errorResp dto.ErrorResponse
		err := json.NewDecoder(w.Body).Decode(&errorResp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, errorResp.Code)
		assert.Equal(t, handler.ErrEmailMismatch, errorResp.Message)
	})

	t.Run("Create user failed", func(t *testing.T) {
		mockUserService := mocks.NewMockUserService(ctrl)
		mockAuthService := mocks.NewMockAuthService(ctrl)
		mockLogger := logger.NewNopLogger()

		mockUserService.EXPECT().
			GetOrCreateUser(gomock.Any(), "user@example.com", "Test User").
			Return(nil, errors.New("failed to create user"))

		h := handler.NewHandler(mockUserService, mockAuthService, mockLogger)

		createReq := dto.CreateUserRequest{
			Email: "user@example.com",
			Name:  "Test User",
		}
		body, _ := json.Marshal(createReq)

		ctx := context.Background()
		ctx = context.WithValue(ctx, contextutil.UserEmailKey, "user@example.com")

		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/users", bytes.NewBuffer(body)).WithContext(ctx)

		h.CreateUser(w, r)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var errorResp dto.ErrorResponse
		err := json.NewDecoder(w.Body).Decode(&errorResp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, errorResp.Code)
		assert.Equal(t, handler.ErrFailedToCreateUser, errorResp.Message)
		assert.Contains(t, errorResp.Error, "failed to create user")
	})

	t.Run("Create user successful", func(t *testing.T) {
		mockUserService := mocks.NewMockUserService(ctrl)
		mockAuthService := mocks.NewMockAuthService(ctrl)
		mockLogger := logger.NewNopLogger()

		objectID := primitive.NewObjectID()
		testUser := &entities.User{
			ID:    objectID,
			Email: "user@example.com",
			Name:  "Test User",
			Wallet: entities.Wallet{
				Diamonds: 100,
				Savings:  0,
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mockUserService.EXPECT().
			GetOrCreateUser(gomock.Any(), "user@example.com", "Test User").
			Return(testUser, nil)

		h := handler.NewHandler(mockUserService, mockAuthService, mockLogger)

		createReq := dto.CreateUserRequest{
			Email: "user@example.com",
			Name:  "Test User",
		}
		body, _ := json.Marshal(createReq)

		ctx := context.Background()
		ctx = context.WithValue(ctx, contextutil.UserEmailKey, "user@example.com")

		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/users", bytes.NewBuffer(body)).WithContext(ctx)

		h.CreateUser(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var response dto.CreateUserResponse
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, testUser.ID.Hex(), response.ID)
		assert.Equal(t, testUser.Email, response.Email)
		assert.Equal(t, testUser.Name, response.Name)
		assert.Equal(t, testUser.Wallet.Diamonds, response.Diamonds)
		assert.Equal(t, testUser.Wallet.Savings, response.Savings)
	})
}

func TestUpdateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("Invalid request format", func(t *testing.T) {
		mockUserService := mocks.NewMockUserService(ctrl)
		mockAuthService := mocks.NewMockAuthService(ctrl)
		mockLogger := logger.NewNopLogger()

		h := handler.NewHandler(mockUserService, mockAuthService, mockLogger)

		invalidBody := bytes.NewBufferString(`{invalid json`)
		w := httptest.NewRecorder()
		r := httptest.NewRequest("PUT", "/users/me", invalidBody)

		h.UpdateUser(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var errorResp dto.ErrorResponse
		err := json.NewDecoder(w.Body).Decode(&errorResp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, errorResp.Code)
		assert.Equal(t, handler.ErrInvalidRequest, errorResp.Message)
	})

	t.Run("Email not in context", func(t *testing.T) {
		mockUserService := mocks.NewMockUserService(ctrl)
		mockAuthService := mocks.NewMockAuthService(ctrl)
		mockLogger := logger.NewNopLogger()

		h := handler.NewHandler(mockUserService, mockAuthService, mockLogger)

		updateReq := dto.UpdateUserRequest{
			Name: "New Name",
		}
		body, _ := json.Marshal(updateReq)

		w := httptest.NewRecorder()
		r := httptest.NewRequest("PUT", "/users/me", bytes.NewBuffer(body))

		h.UpdateUser(w, r)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var errorResp dto.ErrorResponse
		err := json.NewDecoder(w.Body).Decode(&errorResp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, errorResp.Code)
		assert.Equal(t, handler.ErrEmailNotFound, errorResp.Message)
	})

	t.Run("Update user failed", func(t *testing.T) {
		mockUserService := mocks.NewMockUserService(ctrl)
		mockAuthService := mocks.NewMockAuthService(ctrl)
		mockLogger := logger.NewNopLogger()

		mockUserService.EXPECT().
			UpdateUserName(gomock.Any(), "user@example.com", "New Name").
			Return(nil, errors.New("update failed"))

		h := handler.NewHandler(mockUserService, mockAuthService, mockLogger)

		updateReq := dto.UpdateUserRequest{
			Name: "New Name",
		}
		body, _ := json.Marshal(updateReq)

		ctx := context.Background()
		ctx = context.WithValue(ctx, contextutil.UserEmailKey, "user@example.com")

		w := httptest.NewRecorder()
		r := httptest.NewRequest("PUT", "/users/me", bytes.NewBuffer(body)).WithContext(ctx)

		h.UpdateUser(w, r)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var errorResp dto.ErrorResponse
		err := json.NewDecoder(w.Body).Decode(&errorResp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, errorResp.Code)
		assert.Equal(t, handler.ErrFailedToUpdateUser, errorResp.Message)
		assert.Contains(t, errorResp.Error, "update failed")
	})

	t.Run("Update user successful", func(t *testing.T) {
		mockUserService := mocks.NewMockUserService(ctrl)
		mockAuthService := mocks.NewMockAuthService(ctrl)
		mockLogger := logger.NewNopLogger()

		objectID := primitive.NewObjectID()
		testUser := &entities.User{
			ID:    objectID,
			Email: "user@example.com",
			Name:  "New Name",
			Wallet: entities.Wallet{
				Diamonds: 100,
				Savings:  5000,
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mockUserService.EXPECT().
			UpdateUserName(gomock.Any(), "user@example.com", "New Name").
			Return(testUser, nil)

		h := handler.NewHandler(mockUserService, mockAuthService, mockLogger)

		updateReq := dto.UpdateUserRequest{
			Name: "New Name",
		}
		body, _ := json.Marshal(updateReq)

		ctx := context.Background()
		ctx = context.WithValue(ctx, contextutil.UserEmailKey, "user@example.com")

		w := httptest.NewRecorder()
		r := httptest.NewRequest("PUT", "/users/me", bytes.NewBuffer(body)).WithContext(ctx)

		h.UpdateUser(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var response dto.UpdateUserResponse
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, testUser.ID.Hex(), response.ID)
		assert.Equal(t, testUser.Email, response.Email)
		assert.Equal(t, testUser.Name, response.Name)
		assert.Equal(t, testUser.Wallet.Diamonds, response.Diamonds)
		assert.Equal(t, testUser.Wallet.Savings, response.Savings)
	})
}

func TestGetUserWithScope(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 準備測試用户數據
	objectID := primitive.NewObjectID()
	testUser := &entities.User{
		ID:    objectID,
		Email: "user@example.com",
		Name:  "Test User",
		Wallet: entities.Wallet{
			Diamonds: 100,
			Savings:  5000,
		},
		Character: entities.Character{
			ID:       "char_001",
			Name:     "理財顧問 AI",
			ImageURL: "https://example.com/characters/advisor.png",
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	t.Run("User not found in context", func(t *testing.T) {
		mockUserService := mocks.NewMockUserService(ctrl)
		mockAuthService := mocks.NewMockAuthService(ctrl)
		mockLogger := logger.NewNopLogger()

		h := handler.NewHandler(mockUserService, mockAuthService, mockLogger)

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/users/me", nil)

		h.GetUser(w, r)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var errorResp dto.ErrorResponse
		err := json.NewDecoder(w.Body).Decode(&errorResp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, errorResp.Code)
		assert.Equal(t, handler.ErrEmailNotFound, errorResp.Message)
	})

	t.Run("Get user failed", func(t *testing.T) {
		mockUserService := mocks.NewMockUserService(ctrl)
		mockAuthService := mocks.NewMockAuthService(ctrl)
		mockLogger := logger.NewNopLogger()

		mockUserService.EXPECT().
			GetUser(gomock.Any(), "user@example.com").
			Return(nil, errors.New("user not found"))

		h := handler.NewHandler(mockUserService, mockAuthService, mockLogger)

		ctx := context.Background()
		ctx = context.WithValue(ctx, contextutil.UserEmailKey, "user@example.com")

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/users/me", nil).WithContext(ctx)

		h.GetUser(w, r)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var errorResp dto.ErrorResponse
		err := json.NewDecoder(w.Body).Decode(&errorResp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, errorResp.Code)
		assert.Equal(t, handler.ErrFailedToGetUser, errorResp.Message)
		assert.Contains(t, errorResp.Error, "user not found")
	})

	t.Run("Get user with no scope (all info)", func(t *testing.T) {
		mockUserService := mocks.NewMockUserService(ctrl)
		mockAuthService := mocks.NewMockAuthService(ctrl)
		mockLogger := logger.NewNopLogger()

		mockUserService.EXPECT().
			GetUser(gomock.Any(), "user@example.com").
			Return(testUser, nil)

		h := handler.NewHandler(mockUserService, mockAuthService, mockLogger)

		ctx := context.Background()
		ctx = context.WithValue(ctx, contextutil.UserEmailKey, "user@example.com")

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/users/me", nil).WithContext(ctx)

		h.GetUser(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var response dto.GetUserResponse
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, testUser.ID.Hex(), response.ID)
		assert.NotNil(t, response.Email)
		assert.Equal(t, testUser.Email, *response.Email)
		assert.NotNil(t, response.Name)
		assert.Equal(t, testUser.Name, *response.Name)
		assert.NotNil(t, response.Wallet)
		assert.Equal(t, testUser.Wallet.Diamonds, response.Wallet.Diamonds)
		assert.Equal(t, testUser.Wallet.Savings, response.Wallet.Savings)
		assert.NotNil(t, response.Character)
		assert.Equal(t, testUser.Character.ID, response.Character.ID)
		assert.Equal(t, testUser.Character.Name, response.Character.Name)
		assert.Equal(t, testUser.Character.ImageURL, response.Character.ImageURL)
	})

	t.Run("Get user with profile scope only", func(t *testing.T) {
		mockUserService := mocks.NewMockUserService(ctrl)
		mockAuthService := mocks.NewMockAuthService(ctrl)
		mockLogger := logger.NewNopLogger()

		mockUserService.EXPECT().
			GetUser(gomock.Any(), "user@example.com").
			Return(testUser, nil)

		h := handler.NewHandler(mockUserService, mockAuthService, mockLogger)

		ctx := context.Background()
		ctx = context.WithValue(ctx, contextutil.UserEmailKey, "user@example.com")

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/users/me?scope=profile", nil).WithContext(ctx)

		h.GetUser(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var response dto.GetUserResponse
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, testUser.ID.Hex(), response.ID)
		assert.NotNil(t, response.Email)
		assert.Equal(t, testUser.Email, *response.Email)
		assert.NotNil(t, response.Name)
		assert.Equal(t, testUser.Name, *response.Name)
		assert.Nil(t, response.Wallet)
		assert.Nil(t, response.Character)
	})

	t.Run("Get user with wallet scope only", func(t *testing.T) {
		mockUserService := mocks.NewMockUserService(ctrl)
		mockAuthService := mocks.NewMockAuthService(ctrl)
		mockLogger := logger.NewNopLogger()

		mockUserService.EXPECT().
			GetUser(gomock.Any(), "user@example.com").
			Return(testUser, nil)

		h := handler.NewHandler(mockUserService, mockAuthService, mockLogger)

		ctx := context.Background()
		ctx = context.WithValue(ctx, contextutil.UserEmailKey, "user@example.com")

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/users/me?scope=wallet", nil).WithContext(ctx)

		h.GetUser(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var response dto.GetUserResponse
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, testUser.ID.Hex(), response.ID)
		assert.Nil(t, response.Email)
		assert.Nil(t, response.Name)
		assert.NotNil(t, response.Wallet)
		assert.Equal(t, testUser.Wallet.Diamonds, response.Wallet.Diamonds)
		assert.Equal(t, testUser.Wallet.Savings, response.Wallet.Savings)
		assert.Nil(t, response.Character)
	})

	t.Run("Get user with character scope only", func(t *testing.T) {
		mockUserService := mocks.NewMockUserService(ctrl)
		mockAuthService := mocks.NewMockAuthService(ctrl)
		mockLogger := logger.NewNopLogger()

		mockUserService.EXPECT().
			GetUser(gomock.Any(), "user@example.com").
			Return(testUser, nil)

		h := handler.NewHandler(mockUserService, mockAuthService, mockLogger)

		ctx := context.Background()
		ctx = context.WithValue(ctx, contextutil.UserEmailKey, "user@example.com")

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/users/me?scope=character", nil).WithContext(ctx)

		h.GetUser(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var response dto.GetUserResponse
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, testUser.ID.Hex(), response.ID)
		assert.Nil(t, response.Email)
		assert.Nil(t, response.Name)
		assert.Nil(t, response.Wallet)
		assert.NotNil(t, response.Character)
		assert.Equal(t, testUser.Character.ID, response.Character.ID)
		assert.Equal(t, testUser.Character.Name, response.Character.Name)
		assert.Equal(t, testUser.Character.ImageURL, response.Character.ImageURL)
	})

	t.Run("Get user with multiple scopes", func(t *testing.T) {
		mockUserService := mocks.NewMockUserService(ctrl)
		mockAuthService := mocks.NewMockAuthService(ctrl)
		mockLogger := logger.NewNopLogger()

		mockUserService.EXPECT().
			GetUser(gomock.Any(), "user@example.com").
			Return(testUser, nil)

		h := handler.NewHandler(mockUserService, mockAuthService, mockLogger)

		ctx := context.Background()
		ctx = context.WithValue(ctx, contextutil.UserEmailKey, "user@example.com")

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/users/me?scope=profile&scope=wallet", nil).WithContext(ctx)

		h.GetUser(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var response dto.GetUserResponse
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, testUser.ID.Hex(), response.ID)
		assert.NotNil(t, response.Email)
		assert.Equal(t, testUser.Email, *response.Email)
		assert.NotNil(t, response.Name)
		assert.Equal(t, testUser.Name, *response.Name)
		assert.NotNil(t, response.Wallet)
		assert.Equal(t, testUser.Wallet.Diamonds, response.Wallet.Diamonds)
		assert.Equal(t, testUser.Wallet.Savings, response.Wallet.Savings)
		assert.Nil(t, response.Character)
	})

	t.Run("Get user with all scopes explicitly", func(t *testing.T) {
		mockUserService := mocks.NewMockUserService(ctrl)
		mockAuthService := mocks.NewMockAuthService(ctrl)
		mockLogger := logger.NewNopLogger()

		mockUserService.EXPECT().
			GetUser(gomock.Any(), "user@example.com").
			Return(testUser, nil)

		h := handler.NewHandler(mockUserService, mockAuthService, mockLogger)

		ctx := context.Background()
		ctx = context.WithValue(ctx, contextutil.UserEmailKey, "user@example.com")

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/users/me?scope=profile&scope=wallet&scope=character", nil).WithContext(ctx)

		h.GetUser(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var response dto.GetUserResponse
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, testUser.ID.Hex(), response.ID)
		assert.NotNil(t, response.Email)
		assert.Equal(t, testUser.Email, *response.Email)
		assert.NotNil(t, response.Name)
		assert.Equal(t, testUser.Name, *response.Name)
		assert.NotNil(t, response.Wallet)
		assert.Equal(t, testUser.Wallet.Diamonds, response.Wallet.Diamonds)
		assert.Equal(t, testUser.Wallet.Savings, response.Wallet.Savings)
		assert.NotNil(t, response.Character)
		assert.Equal(t, testUser.Character.ID, response.Character.ID)
		assert.Equal(t, testUser.Character.Name, response.Character.Name)
		assert.Equal(t, testUser.Character.ImageURL, response.Character.ImageURL)
	})

	t.Run("Get user with unknown scope", func(t *testing.T) {
		mockUserService := mocks.NewMockUserService(ctrl)
		mockAuthService := mocks.NewMockAuthService(ctrl)
		mockLogger := logger.NewNopLogger()

		mockUserService.EXPECT().
			GetUser(gomock.Any(), "user@example.com").
			Return(testUser, nil)

		h := handler.NewHandler(mockUserService, mockAuthService, mockLogger)

		ctx := context.Background()
		ctx = context.WithValue(ctx, contextutil.UserEmailKey, "user@example.com")

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/users/me?scope=unknown", nil).WithContext(ctx)

		h.GetUser(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var response dto.GetUserResponse
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, testUser.ID.Hex(), response.ID)
		assert.Nil(t, response.Email)
		assert.Nil(t, response.Name)
		assert.Nil(t, response.Wallet)
		assert.Nil(t, response.Character)
		assert.NotEmpty(t, response.CreatedAt)
		assert.NotEmpty(t, response.UpdatedAt)
	})
}
