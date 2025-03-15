package auth_usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/mock/gomock"

	"github.com/Financial-Partner/server/internal/config"
	"github.com/Financial-Partner/server/internal/entities"
	infraAuth "github.com/Financial-Partner/server/internal/infrastructure/auth"
	auth_domain "github.com/Financial-Partner/server/internal/module/auth/domain"
	auth_usecase "github.com/Financial-Partner/server/internal/module/auth/usecase"
	user_domain "github.com/Financial-Partner/server/internal/module/user/domain"
)

type Mocks struct {
	ctrl             *gomock.Controller
	mockJWTManager   *auth_domain.MockJWTManager
	mockTokenStore   *auth_domain.MockTokenStore
	mockUserService  *user_domain.MockUserService
	mockFirebaseAuth *auth_domain.MockFirebaseAuth
}

func NewMocks(t *testing.T) *Mocks {
	ctrl := gomock.NewController(t)

	return &Mocks{
		ctrl:             ctrl,
		mockJWTManager:   auth_domain.NewMockJWTManager(ctrl),
		mockTokenStore:   auth_domain.NewMockTokenStore(ctrl),
		mockUserService:  user_domain.NewMockUserService(ctrl),
		mockFirebaseAuth: auth_domain.NewMockFirebaseAuth(ctrl),
	}
}

func TestRefreshToken(t *testing.T) {
	cfg := &config.Config{}

	t.Run("Success case", func(t *testing.T) {
		mocks := NewMocks(t)
		defer mocks.ctrl.Finish()
		service := auth_usecase.NewService(cfg, mocks.mockFirebaseAuth, mocks.mockJWTManager, mocks.mockTokenStore, mocks.mockUserService)

		id := primitive.NewObjectID().Hex()
		email := "test@example.com"
		claims := &infraAuth.Claims{
			ID:    id,
			Email: email,
		}

		mocks.mockJWTManager.EXPECT().
			ValidateToken("valid_refresh_token").
			Return(claims, nil)

		mocks.mockTokenStore.EXPECT().
			GetRefreshToken(gomock.Any(), "valid_refresh_token").
			Return(id, nil)

		mocks.mockJWTManager.EXPECT().
			GenerateAccessToken(id, email).
			Return("new_access_token", time.Now().Add(time.Hour), nil)

		mocks.mockJWTManager.EXPECT().
			GenerateRefreshToken(id, email).
			Return("new_refresh_token", time.Now().Add(24*time.Hour), nil)

		mocks.mockTokenStore.EXPECT().
			DeleteRefreshToken(gomock.Any(), "valid_refresh_token").
			Return(nil)

		mocks.mockTokenStore.EXPECT().
			SaveRefreshToken(gomock.Any(), id, "new_refresh_token", gomock.Any()).
			Return(nil)

		accessToken, refreshToken, expiresIn, err := service.RefreshToken(context.Background(), "valid_refresh_token")

		assert.NoError(t, err)
		assert.Equal(t, "new_access_token", accessToken)
		assert.Equal(t, "new_refresh_token", refreshToken)
		assert.Greater(t, expiresIn, 0)
	})

	t.Run("Invalid token", func(t *testing.T) {
		mocks := NewMocks(t)
		defer mocks.ctrl.Finish()
		service := auth_usecase.NewService(cfg, mocks.mockFirebaseAuth, mocks.mockJWTManager, mocks.mockTokenStore, mocks.mockUserService)

		mocks.mockJWTManager.EXPECT().
			ValidateToken("invalid_token").
			Return(nil, errors.New("invalid token"))

		accessToken, refreshToken, expiresIn, err := service.RefreshToken(context.Background(), "invalid_token")

		assert.Error(t, err)
		assert.Equal(t, "", accessToken)
		assert.Equal(t, "", refreshToken)
		assert.Equal(t, 0, expiresIn)
		assert.Contains(t, err.Error(), "invalid refresh token")
	})

	t.Run("Token not found in store", func(t *testing.T) {
		mocks := NewMocks(t)
		defer mocks.ctrl.Finish()
		service := auth_usecase.NewService(cfg, mocks.mockFirebaseAuth, mocks.mockJWTManager, mocks.mockTokenStore, mocks.mockUserService)

		id := primitive.NewObjectID().Hex()
		email := "test@example.com"
		claims := &infraAuth.Claims{
			ID:    id,
			Email: email,
		}

		mocks.mockJWTManager.EXPECT().
			ValidateToken("valid_but_deleted_token").
			Return(claims, nil)

		mocks.mockTokenStore.EXPECT().
			GetRefreshToken(gomock.Any(), "valid_but_deleted_token").
			Return("", errors.New("token not found"))

		accessToken, refreshToken, expiresIn, err := service.RefreshToken(context.Background(), "valid_but_deleted_token")

		assert.Error(t, err)
		assert.Equal(t, "", accessToken)
		assert.Equal(t, "", refreshToken)
		assert.Equal(t, 0, expiresIn)
		assert.Contains(t, err.Error(), "refresh token not found")
	})

	t.Run("ID mismatch", func(t *testing.T) {
		mocks := NewMocks(t)
		defer mocks.ctrl.Finish()
		service := auth_usecase.NewService(cfg, mocks.mockFirebaseAuth, mocks.mockJWTManager, mocks.mockTokenStore, mocks.mockUserService)

		id := primitive.NewObjectID().Hex()
		email := "test@example.com"
		claims := &infraAuth.Claims{
			ID:    id,
			Email: email,
		}

		mocks.mockJWTManager.EXPECT().
			ValidateToken("mismatched_token").
			Return(claims, nil)

		mocks.mockTokenStore.EXPECT().
			GetRefreshToken(gomock.Any(), "mismatched_token").
			Return("mismatched_id", nil)

		accessToken, refreshToken, expiresIn, err := service.RefreshToken(context.Background(), "mismatched_token")

		assert.Error(t, err)
		assert.Equal(t, "", accessToken)
		assert.Equal(t, "", refreshToken)
		assert.Equal(t, 0, expiresIn)
		assert.Contains(t, err.Error(), "token id mismatch")
	})

	t.Run("Failed to generate access token", func(t *testing.T) {
		mocks := NewMocks(t)
		defer mocks.ctrl.Finish()
		service := auth_usecase.NewService(cfg, mocks.mockFirebaseAuth, mocks.mockJWTManager, mocks.mockTokenStore, mocks.mockUserService)

		id := primitive.NewObjectID().Hex()
		email := "test@example.com"
		claims := &infraAuth.Claims{
			ID:    id,
			Email: email,
		}

		mocks.mockJWTManager.EXPECT().
			ValidateToken("valid_token").
			Return(claims, nil)

		mocks.mockTokenStore.EXPECT().
			GetRefreshToken(gomock.Any(), "valid_token").
			Return(id, nil)

		mocks.mockJWTManager.EXPECT().
			GenerateAccessToken(id, email).
			Return("", time.Time{}, errors.New("failed to generate token"))

		accessToken, refreshToken, expiresIn, err := service.RefreshToken(context.Background(), "valid_token")

		assert.Error(t, err)
		assert.Equal(t, "", accessToken)
		assert.Equal(t, "", refreshToken)
		assert.Equal(t, 0, expiresIn)
		assert.Contains(t, err.Error(), "failed to generate access token")
	})

	t.Run("Failed to generate refresh token", func(t *testing.T) {
		mocks := NewMocks(t)
		defer mocks.ctrl.Finish()
		service := auth_usecase.NewService(cfg, mocks.mockFirebaseAuth, mocks.mockJWTManager, mocks.mockTokenStore, mocks.mockUserService)

		id := primitive.NewObjectID().Hex()
		email := "test@example.com"
		claims := &infraAuth.Claims{
			ID:    id,
			Email: email,
		}

		mocks.mockJWTManager.EXPECT().
			ValidateToken("valid_token").
			Return(claims, nil)

		mocks.mockTokenStore.EXPECT().
			GetRefreshToken(gomock.Any(), "valid_token").
			Return(id, nil)

		mocks.mockJWTManager.EXPECT().
			GenerateAccessToken(id, email).
			Return("new_access_token", time.Now().Add(time.Hour), nil)

		mocks.mockJWTManager.EXPECT().
			GenerateRefreshToken(id, email).
			Return("", time.Time{}, errors.New("failed to generate refresh token"))

		accessToken, refreshToken, expiresIn, err := service.RefreshToken(context.Background(), "valid_token")

		assert.Error(t, err)
		assert.Equal(t, "", accessToken)
		assert.Equal(t, "", refreshToken)
		assert.Equal(t, 0, expiresIn)
		assert.Contains(t, err.Error(), "failed to generate refresh token")
	})

	t.Run("Failed to delete old refresh token", func(t *testing.T) {
		mocks := NewMocks(t)
		defer mocks.ctrl.Finish()
		service := auth_usecase.NewService(cfg, mocks.mockFirebaseAuth, mocks.mockJWTManager, mocks.mockTokenStore, mocks.mockUserService)

		id := primitive.NewObjectID().Hex()
		email := "test@example.com"
		claims := &infraAuth.Claims{
			ID:    id,
			Email: email,
		}

		mocks.mockJWTManager.EXPECT().
			ValidateToken("valid_token").
			Return(claims, nil)

		mocks.mockTokenStore.EXPECT().
			GetRefreshToken(gomock.Any(), "valid_token").
			Return(id, nil)

		mocks.mockJWTManager.EXPECT().
			GenerateAccessToken(id, email).
			Return("new_access_token", time.Now().Add(time.Hour), nil)

		mocks.mockJWTManager.EXPECT().
			GenerateRefreshToken(id, email).
			Return("new_refresh_token", time.Now().Add(24*time.Hour), nil)

		mocks.mockTokenStore.EXPECT().
			DeleteRefreshToken(gomock.Any(), "valid_token").
			Return(errors.New("failed to delete token"))

		accessToken, refreshToken, expiresIn, err := service.RefreshToken(context.Background(), "valid_token")

		assert.Error(t, err)
		assert.Equal(t, "", accessToken)
		assert.Equal(t, "", refreshToken)
		assert.Equal(t, 0, expiresIn)
		assert.Contains(t, err.Error(), "failed to delete old refresh token")
	})

	t.Run("Failed to save new refresh token", func(t *testing.T) {
		mocks := NewMocks(t)
		defer mocks.ctrl.Finish()
		service := auth_usecase.NewService(cfg, mocks.mockFirebaseAuth, mocks.mockJWTManager, mocks.mockTokenStore, mocks.mockUserService)

		id := primitive.NewObjectID().Hex()
		email := "test@example.com"
		claims := &infraAuth.Claims{
			ID:    id,
			Email: email,
		}

		mocks.mockJWTManager.EXPECT().
			ValidateToken("valid_token").
			Return(claims, nil)

		mocks.mockTokenStore.EXPECT().
			GetRefreshToken(gomock.Any(), "valid_token").
			Return(id, nil)

		mocks.mockJWTManager.EXPECT().
			GenerateAccessToken(id, email).
			Return("new_access_token", time.Now().Add(time.Hour), nil)

		mocks.mockJWTManager.EXPECT().
			GenerateRefreshToken(id, email).
			Return("new_refresh_token", time.Now().Add(24*time.Hour), nil)

		mocks.mockTokenStore.EXPECT().
			DeleteRefreshToken(gomock.Any(), "valid_token").
			Return(nil)

		mocks.mockTokenStore.EXPECT().
			SaveRefreshToken(gomock.Any(), id, "new_refresh_token", gomock.Any()).
			Return(errors.New("failed to save token"))

		accessToken, refreshToken, expiresIn, err := service.RefreshToken(context.Background(), "valid_token")

		assert.Error(t, err)
		assert.Equal(t, "", accessToken)
		assert.Equal(t, "", refreshToken)
		assert.Equal(t, 0, expiresIn)
		assert.Contains(t, err.Error(), "failed to save new refresh token")
	})
}

func TestLoginWithFirebase(t *testing.T) {
	cfg := &config.Config{}

	t.Run("Success case", func(t *testing.T) {
		mocks := NewMocks(t)
		defer mocks.ctrl.Finish()
		service := auth_usecase.NewService(cfg, mocks.mockFirebaseAuth, mocks.mockJWTManager, mocks.mockTokenStore, mocks.mockUserService)

		id := primitive.NewObjectID()
		email := "test@example.com"
		name := "Test User"
		token := &infraAuth.Token{
			Claims: map[string]interface{}{
				"email": email,
				"name":  name,
			},
		}

		mockUser := &entities.User{
			ID:    id,
			Email: email,
			Name:  name,
		}

		mocks.mockFirebaseAuth.EXPECT().
			VerifyToken(gomock.Any(), "valid_firebase_token").
			Return(token, nil)

		mocks.mockUserService.EXPECT().
			GetOrCreateUser(gomock.Any(), email, name).
			Return(mockUser, nil)

		mocks.mockJWTManager.EXPECT().
			GenerateAccessToken(id.Hex(), email).
			Return("access_token", time.Now().Add(time.Hour), nil)

		mocks.mockJWTManager.EXPECT().
			GenerateRefreshToken(id.Hex(), email).
			Return("refresh_token", time.Now().Add(24*time.Hour), nil)

		mocks.mockTokenStore.EXPECT().
			SaveRefreshToken(gomock.Any(), id.Hex(), "refresh_token", gomock.Any()).
			Return(nil)

		accessToken, refreshToken, expiresIn, user, err := service.LoginWithFirebase(context.Background(), "valid_firebase_token")

		assert.NoError(t, err)
		assert.Equal(t, "access_token", accessToken)
		assert.Equal(t, "refresh_token", refreshToken)
		assert.Greater(t, expiresIn, 0)
		assert.Equal(t, mockUser, user)
	})

	t.Run("Invalid firebase token", func(t *testing.T) {
		mocks := NewMocks(t)
		defer mocks.ctrl.Finish()
		service := auth_usecase.NewService(cfg, mocks.mockFirebaseAuth, mocks.mockJWTManager, mocks.mockTokenStore, mocks.mockUserService)

		mocks.mockFirebaseAuth.EXPECT().
			VerifyToken(gomock.Any(), "invalid_token").
			Return(nil, errors.New("invalid token"))

		accessToken, refreshToken, expiresIn, user, err := service.LoginWithFirebase(context.Background(), "invalid_token")

		assert.Error(t, err)
		assert.Equal(t, "", accessToken)
		assert.Equal(t, "", refreshToken)
		assert.Equal(t, 0, expiresIn)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "invalid firebase token")
	})

	t.Run("Missing email in token claims", func(t *testing.T) {
		mocks := NewMocks(t)
		defer mocks.ctrl.Finish()
		service := auth_usecase.NewService(cfg, mocks.mockFirebaseAuth, mocks.mockJWTManager, mocks.mockTokenStore, mocks.mockUserService)

		token := &infraAuth.Token{
			Claims: map[string]interface{}{
				"name": "Test User",
			},
		}

		mocks.mockFirebaseAuth.EXPECT().
			VerifyToken(gomock.Any(), "token_without_email").
			Return(token, nil)

		accessToken, refreshToken, expiresIn, user, err := service.LoginWithFirebase(context.Background(), "token_without_email")

		assert.Error(t, err)
		assert.Equal(t, "", accessToken)
		assert.Equal(t, "", refreshToken)
		assert.Equal(t, 0, expiresIn)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "email not found in token claims")
	})

	t.Run("Failed to get or create user", func(t *testing.T) {
		mocks := NewMocks(t)
		defer mocks.ctrl.Finish()
		service := auth_usecase.NewService(cfg, mocks.mockFirebaseAuth, mocks.mockJWTManager, mocks.mockTokenStore, mocks.mockUserService)

		email := "test@example.com"
		name := "Test User"
		token := &infraAuth.Token{
			Claims: map[string]interface{}{
				"email": email,
				"name":  name,
			},
		}

		mocks.mockFirebaseAuth.EXPECT().
			VerifyToken(gomock.Any(), "valid_token").
			Return(token, nil)

		mocks.mockUserService.EXPECT().
			GetOrCreateUser(gomock.Any(), email, name).
			Return(nil, errors.New("database error"))

		accessToken, refreshToken, expiresIn, user, err := service.LoginWithFirebase(context.Background(), "valid_token")

		assert.Error(t, err)
		assert.Equal(t, "", accessToken)
		assert.Equal(t, "", refreshToken)
		assert.Equal(t, 0, expiresIn)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "failed to get or create user")
	})

	t.Run("Failed to generate access token", func(t *testing.T) {
		mocks := NewMocks(t)
		defer mocks.ctrl.Finish()
		service := auth_usecase.NewService(cfg, mocks.mockFirebaseAuth, mocks.mockJWTManager, mocks.mockTokenStore, mocks.mockUserService)

		email := "test@example.com"
		name := "Test User"
		token := &infraAuth.Token{
			Claims: map[string]interface{}{
				"email": email,
				"name":  name,
			},
		}

		id := primitive.NewObjectID()
		mockUser := &entities.User{
			ID:    id,
			Email: email,
			Name:  name,
		}

		mocks.mockFirebaseAuth.EXPECT().
			VerifyToken(gomock.Any(), "valid_token").
			Return(token, nil)

		mocks.mockUserService.EXPECT().
			GetOrCreateUser(gomock.Any(), email, name).
			Return(mockUser, nil)

		mocks.mockJWTManager.EXPECT().
			GenerateAccessToken(id.Hex(), email).
			Return("", time.Time{}, errors.New("failed to generate token"))

		accessToken, refreshToken, expiresIn, user, err := service.LoginWithFirebase(context.Background(), "valid_token")

		assert.Error(t, err)
		assert.Equal(t, "", accessToken)
		assert.Equal(t, "", refreshToken)
		assert.Equal(t, 0, expiresIn)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "failed to generate access token")
	})

	t.Run("Failed to generate refresh token", func(t *testing.T) {
		mocks := NewMocks(t)
		defer mocks.ctrl.Finish()
		service := auth_usecase.NewService(cfg, mocks.mockFirebaseAuth, mocks.mockJWTManager, mocks.mockTokenStore, mocks.mockUserService)

		email := "test@example.com"
		name := "Test User"
		token := &infraAuth.Token{
			Claims: map[string]interface{}{
				"email": email,
				"name":  name,
			},
		}

		id := primitive.NewObjectID()
		mockUser := &entities.User{
			ID:    id,
			Email: email,
			Name:  name,
		}

		mocks.mockFirebaseAuth.EXPECT().
			VerifyToken(gomock.Any(), "valid_token").
			Return(token, nil)

		mocks.mockUserService.EXPECT().
			GetOrCreateUser(gomock.Any(), email, name).
			Return(mockUser, nil)

		mocks.mockJWTManager.EXPECT().
			GenerateAccessToken(id.Hex(), email).
			Return("access_token", time.Now().Add(time.Hour), nil)

		mocks.mockJWTManager.EXPECT().
			GenerateRefreshToken(id.Hex(), email).
			Return("", time.Time{}, errors.New("failed to generate token"))

		accessToken, refreshToken, expiresIn, user, err := service.LoginWithFirebase(context.Background(), "valid_token")

		assert.Error(t, err)
		assert.Equal(t, "", accessToken)
		assert.Equal(t, "", refreshToken)
		assert.Equal(t, 0, expiresIn)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "failed to generate refresh token")
	})

	t.Run("Failed to save refresh token", func(t *testing.T) {
		mocks := NewMocks(t)
		defer mocks.ctrl.Finish()
		service := auth_usecase.NewService(cfg, mocks.mockFirebaseAuth, mocks.mockJWTManager, mocks.mockTokenStore, mocks.mockUserService)

		email := "test@example.com"
		name := "Test User"
		token := &infraAuth.Token{
			Claims: map[string]interface{}{
				"email": email,
				"name":  name,
			},
		}

		id := primitive.NewObjectID()
		mockUser := &entities.User{
			ID:    id,
			Email: email,
			Name:  name,
		}

		mocks.mockFirebaseAuth.EXPECT().
			VerifyToken(gomock.Any(), "valid_token").
			Return(token, nil)

		mocks.mockUserService.EXPECT().
			GetOrCreateUser(gomock.Any(), email, name).
			Return(mockUser, nil)

		mocks.mockJWTManager.EXPECT().
			GenerateAccessToken(id.Hex(), email).
			Return("access_token", time.Now().Add(time.Hour), nil)

		mocks.mockJWTManager.EXPECT().
			GenerateRefreshToken(id.Hex(), email).
			Return("refresh_token", time.Now().Add(24*time.Hour), nil)

		mocks.mockTokenStore.EXPECT().
			SaveRefreshToken(gomock.Any(), id.Hex(), "refresh_token", gomock.Any()).
			Return(errors.New("database error"))

		accessToken, refreshToken, expiresIn, user, err := service.LoginWithFirebase(context.Background(), "valid_token")

		assert.Error(t, err)
		assert.Equal(t, "", accessToken)
		assert.Equal(t, "", refreshToken)
		assert.Equal(t, 0, expiresIn)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "failed to save refresh token")
	})

	t.Run("Empty name uses email as name", func(t *testing.T) {
		mocks := NewMocks(t)
		defer mocks.ctrl.Finish()
		service := auth_usecase.NewService(cfg, mocks.mockFirebaseAuth, mocks.mockJWTManager, mocks.mockTokenStore, mocks.mockUserService)

		email := "test@example.com"
		token := &infraAuth.Token{
			Claims: map[string]interface{}{
				"email": email,
			},
		}

		id := primitive.NewObjectID()
		mockUser := &entities.User{
			ID:    id,
			Email: email,
			Name:  email,
		}

		mocks.mockFirebaseAuth.EXPECT().
			VerifyToken(gomock.Any(), "token_without_name").
			Return(token, nil)

		mocks.mockUserService.EXPECT().
			GetOrCreateUser(gomock.Any(), email, email).
			Return(mockUser, nil)

		mocks.mockJWTManager.EXPECT().
			GenerateAccessToken(id.Hex(), email).
			Return("access_token", time.Now().Add(time.Hour), nil)

		mocks.mockJWTManager.EXPECT().
			GenerateRefreshToken(id.Hex(), email).
			Return("refresh_token", time.Now().Add(24*time.Hour), nil)

		mocks.mockTokenStore.EXPECT().
			SaveRefreshToken(gomock.Any(), id.Hex(), "refresh_token", gomock.Any()).
			Return(nil)

		accessToken, refreshToken, expiresIn, user, err := service.LoginWithFirebase(context.Background(), "token_without_name")

		assert.NoError(t, err)
		assert.Equal(t, "access_token", accessToken)
		assert.Equal(t, "refresh_token", refreshToken)
		assert.Greater(t, expiresIn, 0)
		assert.Equal(t, mockUser, user)
	})
}

func TestLogout(t *testing.T) {
	cfg := &config.Config{}

	t.Run("Success case", func(t *testing.T) {
		mocks := NewMocks(t)
		defer mocks.ctrl.Finish()
		service := auth_usecase.NewService(cfg, mocks.mockFirebaseAuth, mocks.mockJWTManager, mocks.mockTokenStore, mocks.mockUserService)

		mocks.mockTokenStore.EXPECT().
			DeleteRefreshToken(gomock.Any(), "valid_refresh_token").
			Return(nil)

		err := service.Logout(context.Background(), "valid_refresh_token")

		assert.NoError(t, err)
	})

	t.Run("Failed to delete refresh token", func(t *testing.T) {
		mocks := NewMocks(t)
		defer mocks.ctrl.Finish()
		service := auth_usecase.NewService(cfg, mocks.mockFirebaseAuth, mocks.mockJWTManager, mocks.mockTokenStore, mocks.mockUserService)

		mocks.mockTokenStore.EXPECT().
			DeleteRefreshToken(gomock.Any(), "valid_refresh_token").
			Return(errors.New("database error"))

		err := service.Logout(context.Background(), "valid_refresh_token")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to delete refresh token")
	})
}
