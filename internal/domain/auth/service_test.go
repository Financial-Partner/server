package auth_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Financial-Partner/server/internal/config"
	"github.com/Financial-Partner/server/internal/domain/auth"
	"github.com/Financial-Partner/server/internal/domain/auth/mocks"
	"github.com/Financial-Partner/server/internal/entities"
	infraAuth "github.com/Financial-Partner/server/internal/infrastructure/auth"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestRefreshToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockJWTManager := mocks.NewMockJWTManager(ctrl)
	mockTokenStore := mocks.NewMockTokenStore(ctrl)
	mockUserService := mocks.NewMockUserService(ctrl)

	cfg := &config.Config{}
	service := auth.NewService(cfg, nil, mockJWTManager, mockTokenStore, mockUserService)

	t.Run("Success case", func(t *testing.T) {
		claims := &infraAuth.Claims{Email: "test@example.com"}

		mockJWTManager.EXPECT().
			ValidateToken("valid_refresh_token").
			Return(claims, nil)

		mockTokenStore.EXPECT().
			GetRefreshToken(gomock.Any(), "valid_refresh_token").
			Return("test@example.com", nil)

		mockJWTManager.EXPECT().
			GenerateAccessToken("test@example.com").
			Return("new_access_token", time.Now().Add(time.Hour), nil)

		mockJWTManager.EXPECT().
			GenerateRefreshToken("test@example.com").
			Return("new_refresh_token", time.Now().Add(24*time.Hour), nil)

		mockTokenStore.EXPECT().
			DeleteRefreshToken(gomock.Any(), "valid_refresh_token").
			Return(nil)

		mockTokenStore.EXPECT().
			SaveRefreshToken(gomock.Any(), "test@example.com", "new_refresh_token", gomock.Any()).
			Return(nil)

		accessToken, refreshToken, expiresIn, err := service.RefreshToken(context.Background(), "valid_refresh_token")

		assert.NoError(t, err)
		assert.Equal(t, "new_access_token", accessToken)
		assert.Equal(t, "new_refresh_token", refreshToken)
		assert.Greater(t, expiresIn, 0)
	})

	t.Run("Invalid token", func(t *testing.T) {
		mockJWTManager.EXPECT().
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
		claims := &infraAuth.Claims{Email: "test@example.com"}

		mockJWTManager.EXPECT().
			ValidateToken("valid_but_deleted_token").
			Return(claims, nil)

		mockTokenStore.EXPECT().
			GetRefreshToken(gomock.Any(), "valid_but_deleted_token").
			Return("", errors.New("token not found"))

		accessToken, refreshToken, expiresIn, err := service.RefreshToken(context.Background(), "valid_but_deleted_token")

		assert.Error(t, err)
		assert.Equal(t, "", accessToken)
		assert.Equal(t, "", refreshToken)
		assert.Equal(t, 0, expiresIn)
		assert.Contains(t, err.Error(), "refresh token not found")
	})

	t.Run("Email mismatch", func(t *testing.T) {
		claims := &infraAuth.Claims{Email: "test@example.com"}

		mockJWTManager.EXPECT().
			ValidateToken("mismatched_token").
			Return(claims, nil)

		mockTokenStore.EXPECT().
			GetRefreshToken(gomock.Any(), "mismatched_token").
			Return("other@example.com", nil)

		accessToken, refreshToken, expiresIn, err := service.RefreshToken(context.Background(), "mismatched_token")

		assert.Error(t, err)
		assert.Equal(t, "", accessToken)
		assert.Equal(t, "", refreshToken)
		assert.Equal(t, 0, expiresIn)
		assert.Contains(t, err.Error(), "token email mismatch")
	})

	t.Run("Failed to generate access token", func(t *testing.T) {
		claims := &infraAuth.Claims{Email: "test@example.com"}

		mockJWTManager.EXPECT().
			ValidateToken("valid_token").
			Return(claims, nil)

		mockTokenStore.EXPECT().
			GetRefreshToken(gomock.Any(), "valid_token").
			Return("test@example.com", nil)

		mockJWTManager.EXPECT().
			GenerateAccessToken("test@example.com").
			Return("", time.Time{}, errors.New("failed to generate token"))

		accessToken, refreshToken, expiresIn, err := service.RefreshToken(context.Background(), "valid_token")

		assert.Error(t, err)
		assert.Equal(t, "", accessToken)
		assert.Equal(t, "", refreshToken)
		assert.Equal(t, 0, expiresIn)
		assert.Contains(t, err.Error(), "failed to generate access token")
	})

	t.Run("Failed to generate refresh token", func(t *testing.T) {
		claims := &infraAuth.Claims{Email: "test@example.com"}

		mockJWTManager.EXPECT().
			ValidateToken("valid_token").
			Return(claims, nil)

		mockTokenStore.EXPECT().
			GetRefreshToken(gomock.Any(), "valid_token").
			Return("test@example.com", nil)

		mockJWTManager.EXPECT().
			GenerateAccessToken("test@example.com").
			Return("new_access_token", time.Now().Add(time.Hour), nil)

		mockJWTManager.EXPECT().
			GenerateRefreshToken("test@example.com").
			Return("", time.Time{}, errors.New("failed to generate refresh token"))

		accessToken, refreshToken, expiresIn, err := service.RefreshToken(context.Background(), "valid_token")

		assert.Error(t, err)
		assert.Equal(t, "", accessToken)
		assert.Equal(t, "", refreshToken)
		assert.Equal(t, 0, expiresIn)
		assert.Contains(t, err.Error(), "failed to generate refresh token")
	})

	t.Run("Failed to delete old refresh token", func(t *testing.T) {
		claims := &infraAuth.Claims{Email: "test@example.com"}

		mockJWTManager.EXPECT().
			ValidateToken("valid_token").
			Return(claims, nil)

		mockTokenStore.EXPECT().
			GetRefreshToken(gomock.Any(), "valid_token").
			Return("test@example.com", nil)

		mockJWTManager.EXPECT().
			GenerateAccessToken("test@example.com").
			Return("new_access_token", time.Now().Add(time.Hour), nil)

		mockJWTManager.EXPECT().
			GenerateRefreshToken("test@example.com").
			Return("new_refresh_token", time.Now().Add(24*time.Hour), nil)

		mockTokenStore.EXPECT().
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
		claims := &infraAuth.Claims{Email: "test@example.com"}

		mockJWTManager.EXPECT().
			ValidateToken("valid_token").
			Return(claims, nil)

		mockTokenStore.EXPECT().
			GetRefreshToken(gomock.Any(), "valid_token").
			Return("test@example.com", nil)

		mockJWTManager.EXPECT().
			GenerateAccessToken("test@example.com").
			Return("new_access_token", time.Now().Add(time.Hour), nil)

		mockJWTManager.EXPECT().
			GenerateRefreshToken("test@example.com").
			Return("new_refresh_token", time.Now().Add(24*time.Hour), nil)

		mockTokenStore.EXPECT().
			DeleteRefreshToken(gomock.Any(), "valid_token").
			Return(nil)

		mockTokenStore.EXPECT().
			SaveRefreshToken(gomock.Any(), "test@example.com", "new_refresh_token", gomock.Any()).
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
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFirebaseAuth := mocks.NewMockFirebaseAuth(ctrl)
	mockJWTManager := mocks.NewMockJWTManager(ctrl)
	mockTokenStore := mocks.NewMockTokenStore(ctrl)
	mockUserService := mocks.NewMockUserService(ctrl)

	cfg := &config.Config{}
	service := auth.NewService(cfg, mockFirebaseAuth, mockJWTManager, mockTokenStore, mockUserService)

	t.Run("Success case", func(t *testing.T) {
		token := &infraAuth.Token{
			Claims: map[string]interface{}{
				"email": "test@example.com",
				"name":  "Test User",
			},
		}

		mockUser := &entities.User{
			Email: "test@example.com",
			Name:  "Test User",
		}

		mockFirebaseAuth.EXPECT().
			VerifyToken(gomock.Any(), "valid_firebase_token").
			Return(token, nil)

		mockUserService.EXPECT().
			GetOrCreateUser(gomock.Any(), "test@example.com", "Test User").
			Return(mockUser, nil)

		mockJWTManager.EXPECT().
			GenerateAccessToken("test@example.com").
			Return("access_token", time.Now().Add(time.Hour), nil)

		mockJWTManager.EXPECT().
			GenerateRefreshToken("test@example.com").
			Return("refresh_token", time.Now().Add(24*time.Hour), nil)

		mockTokenStore.EXPECT().
			SaveRefreshToken(gomock.Any(), "test@example.com", "refresh_token", gomock.Any()).
			Return(nil)

		accessToken, refreshToken, expiresIn, user, err := service.LoginWithFirebase(context.Background(), "valid_firebase_token")

		assert.NoError(t, err)
		assert.Equal(t, "access_token", accessToken)
		assert.Equal(t, "refresh_token", refreshToken)
		assert.Greater(t, expiresIn, 0)
		assert.Equal(t, mockUser, user)
	})

	t.Run("Invalid firebase token", func(t *testing.T) {
		mockFirebaseAuth.EXPECT().
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
		token := &infraAuth.Token{
			Claims: map[string]interface{}{
				"name": "Test User",
			},
		}

		mockFirebaseAuth.EXPECT().
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
		token := &infraAuth.Token{
			Claims: map[string]interface{}{
				"email": "test@example.com",
				"name":  "Test User",
			},
		}

		mockFirebaseAuth.EXPECT().
			VerifyToken(gomock.Any(), "valid_token").
			Return(token, nil)

		mockUserService.EXPECT().
			GetOrCreateUser(gomock.Any(), "test@example.com", "Test User").
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
		token := &infraAuth.Token{
			Claims: map[string]interface{}{
				"email": "test@example.com",
				"name":  "Test User",
			},
		}

		mockUser := &entities.User{
			Email: "test@example.com",
			Name:  "Test User",
		}

		mockFirebaseAuth.EXPECT().
			VerifyToken(gomock.Any(), "valid_token").
			Return(token, nil)

		mockUserService.EXPECT().
			GetOrCreateUser(gomock.Any(), "test@example.com", "Test User").
			Return(mockUser, nil)

		mockJWTManager.EXPECT().
			GenerateAccessToken("test@example.com").
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
		token := &infraAuth.Token{
			Claims: map[string]interface{}{
				"email": "test@example.com",
				"name":  "Test User",
			},
		}

		mockUser := &entities.User{
			Email: "test@example.com",
			Name:  "Test User",
		}

		mockFirebaseAuth.EXPECT().
			VerifyToken(gomock.Any(), "valid_token").
			Return(token, nil)

		mockUserService.EXPECT().
			GetOrCreateUser(gomock.Any(), "test@example.com", "Test User").
			Return(mockUser, nil)

		mockJWTManager.EXPECT().
			GenerateAccessToken("test@example.com").
			Return("access_token", time.Now().Add(time.Hour), nil)

		mockJWTManager.EXPECT().
			GenerateRefreshToken("test@example.com").
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
		token := &infraAuth.Token{
			Claims: map[string]interface{}{
				"email": "test@example.com",
				"name":  "Test User",
			},
		}

		mockUser := &entities.User{
			Email: "test@example.com",
			Name:  "Test User",
		}

		mockFirebaseAuth.EXPECT().
			VerifyToken(gomock.Any(), "valid_token").
			Return(token, nil)

		mockUserService.EXPECT().
			GetOrCreateUser(gomock.Any(), "test@example.com", "Test User").
			Return(mockUser, nil)

		mockJWTManager.EXPECT().
			GenerateAccessToken("test@example.com").
			Return("access_token", time.Now().Add(time.Hour), nil)

		mockJWTManager.EXPECT().
			GenerateRefreshToken("test@example.com").
			Return("refresh_token", time.Now().Add(24*time.Hour), nil)

		mockTokenStore.EXPECT().
			SaveRefreshToken(gomock.Any(), "test@example.com", "refresh_token", gomock.Any()).
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
		token := &infraAuth.Token{
			Claims: map[string]interface{}{
				"email": "test@example.com",
			},
		}

		mockUser := &entities.User{
			Email: "test@example.com",
			Name:  "test@example.com",
		}

		mockFirebaseAuth.EXPECT().
			VerifyToken(gomock.Any(), "token_without_name").
			Return(token, nil)

		mockUserService.EXPECT().
			GetOrCreateUser(gomock.Any(), "test@example.com", "test@example.com").
			Return(mockUser, nil)

		mockJWTManager.EXPECT().
			GenerateAccessToken("test@example.com").
			Return("access_token", time.Now().Add(time.Hour), nil)

		mockJWTManager.EXPECT().
			GenerateRefreshToken("test@example.com").
			Return("refresh_token", time.Now().Add(24*time.Hour), nil)

		mockTokenStore.EXPECT().
			SaveRefreshToken(gomock.Any(), "test@example.com", "refresh_token", gomock.Any()).
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
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockJWTManager := mocks.NewMockJWTManager(ctrl)
	mockTokenStore := mocks.NewMockTokenStore(ctrl)
	mockUserService := mocks.NewMockUserService(ctrl)

	cfg := &config.Config{}
	service := auth.NewService(cfg, nil, mockJWTManager, mockTokenStore, mockUserService)

	t.Run("Success case", func(t *testing.T) {
		mockTokenStore.EXPECT().
			DeleteRefreshToken(gomock.Any(), "valid_refresh_token").
			Return(nil)

		err := service.Logout(context.Background(), "valid_refresh_token")

		assert.NoError(t, err)
	})

	t.Run("Failed to delete refresh token", func(t *testing.T) {
		mockTokenStore.EXPECT().
			DeleteRefreshToken(gomock.Any(), "valid_refresh_token").
			Return(errors.New("database error"))

		err := service.Logout(context.Background(), "valid_refresh_token")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to delete refresh token")
	})
}
