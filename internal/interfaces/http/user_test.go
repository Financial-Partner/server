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
	"github.com/Financial-Partner/server/internal/interfaces/http/dto"
	httperror "github.com/Financial-Partner/server/internal/interfaces/http/error"
)

func TestUpdateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("Invalid request format", func(t *testing.T) {
		h, _ := newTestHandler(t)

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
		assert.Equal(t, httperror.ErrInvalidRequest, errorResp.Message)
	})

	t.Run("Email not in context", func(t *testing.T) {
		h, _ := newTestHandler(t)

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
		assert.Equal(t, httperror.ErrEmailNotFound, errorResp.Message)
	})

	t.Run("Update user failed", func(t *testing.T) {
		h, mockServices := newTestHandler(t)

		mockServices.UserService.EXPECT().
			UpdateUserName(gomock.Any(), "user@example.com", "New Name").
			Return(nil, errors.New("update failed"))

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
		assert.Equal(t, httperror.ErrFailedToUpdateUser, errorResp.Message)
	})

	t.Run("Update user successful", func(t *testing.T) {
		h, mockServices := newTestHandler(t)

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

		mockServices.UserService.EXPECT().
			UpdateUserName(gomock.Any(), "user@example.com", "New Name").
			Return(testUser, nil)

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
		h, _ := newTestHandler(t)

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/users/me", nil)

		h.GetUser(w, r)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var errorResp dto.ErrorResponse
		err := json.NewDecoder(w.Body).Decode(&errorResp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, errorResp.Code)
		assert.Equal(t, httperror.ErrEmailNotFound, errorResp.Message)
	})

	t.Run("Get user failed", func(t *testing.T) {
		h, mockServices := newTestHandler(t)

		mockServices.UserService.EXPECT().
			GetUser(gomock.Any(), "user@example.com").
			Return(nil, errors.New("user not found"))

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
		assert.Equal(t, httperror.ErrFailedToGetUser, errorResp.Message)
	})

	t.Run("Get user with no scope (all info)", func(t *testing.T) {
		h, mockServices := newTestHandler(t)

		mockServices.UserService.EXPECT().
			GetUser(gomock.Any(), "user@example.com").
			Return(testUser, nil)

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
		h, mockServices := newTestHandler(t)

		mockServices.UserService.EXPECT().
			GetUser(gomock.Any(), "user@example.com").
			Return(testUser, nil)

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
		h, mockServices := newTestHandler(t)

		mockServices.UserService.EXPECT().
			GetUser(gomock.Any(), "user@example.com").
			Return(testUser, nil)

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
		h, mockServices := newTestHandler(t)

		mockServices.UserService.EXPECT().
			GetUser(gomock.Any(), "user@example.com").
			Return(testUser, nil)

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
		h, mockServices := newTestHandler(t)

		mockServices.UserService.EXPECT().
			GetUser(gomock.Any(), "user@example.com").
			Return(testUser, nil)

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
		h, mockServices := newTestHandler(t)

		mockServices.UserService.EXPECT().
			GetUser(gomock.Any(), "user@example.com").
			Return(testUser, nil)

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
		h, mockServices := newTestHandler(t)

		mockServices.UserService.EXPECT().
			GetUser(gomock.Any(), "user@example.com").
			Return(testUser, nil)

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
