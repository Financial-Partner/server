package handler_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/Financial-Partner/server/internal/contextutil"
	"github.com/Financial-Partner/server/internal/domain/user"
	"github.com/Financial-Partner/server/internal/infrastructure/logger"
	handler "github.com/Financial-Partner/server/internal/interfaces/http"
	"github.com/Financial-Partner/server/internal/interfaces/http/mocks"
)

func TestGetUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("User not found in context", func(t *testing.T) {
		mockUserService := mocks.NewMockUserService(ctrl)
		mockLogger := logger.NewNopLogger()

		h := handler.NewHandler(mockUserService, mockLogger)

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/user", nil)

		h.GetUser(w, r)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("GetUser error", func(t *testing.T) {
		mockUserService := mocks.NewMockUserService(ctrl)
		mockLogger := logger.NewNopLogger()

		mockUserService.EXPECT().
			GetOrCreateUser(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil, errors.New("create user error"))

		h := handler.NewHandler(mockUserService, mockLogger)
		ctx := context.Background()
		ctx = context.WithValue(ctx, contextutil.UserEmailKey, "test@example.com")

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/user", nil).WithContext(ctx)

		h.GetUser(w, r)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("User found", func(t *testing.T) {
		mockUserService := mocks.NewMockUserService(ctrl)
		mockLogger := logger.NewNopLogger()

		testUser := &user.UserEntity{
			Email: "test@example.com",
		}

		mockUserService.EXPECT().
			GetOrCreateUser(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(testUser, nil)

		h := handler.NewHandler(mockUserService, mockLogger)
		ctx := context.Background()
		ctx = context.WithValue(ctx, contextutil.UserEmailKey, testUser.Email)

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/user", nil).WithContext(ctx)

		h.GetUser(w, r)

		assert.Equal(t, http.StatusOK, w.Code)

		var actualUser user.UserEntity
		err := json.NewDecoder(w.Body).Decode(&actualUser)
		assert.NoError(t, err)
		assert.Equal(t, testUser, &actualUser)
	})
}
