package auth_usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/Financial-Partner/server/internal/config"
	"github.com/Financial-Partner/server/internal/entities"
	auth_domain "github.com/Financial-Partner/server/internal/module/auth/domain"
	user_domain "github.com/Financial-Partner/server/internal/module/user/domain"
)

type Service struct {
	cfg          *config.Config
	firebaseAuth auth_domain.FirebaseAuth
	jwtManager   auth_domain.JWTManager
	tokenStore   auth_domain.TokenStore
	userService  user_domain.UserService
}

func NewService(
	cfg *config.Config,
	firebaseAuth auth_domain.FirebaseAuth,
	jwtManager auth_domain.JWTManager,
	tokenStore auth_domain.TokenStore,
	userService user_domain.UserService,
) *Service {
	return &Service{
		cfg:          cfg,
		firebaseAuth: firebaseAuth,
		jwtManager:   jwtManager,
		tokenStore:   tokenStore,
		userService:  userService,
	}
}

func (s *Service) LoginWithFirebase(ctx context.Context, firebaseToken string) (
	accessToken, refreshToken string, expiresIn int, userInfo *entities.User, err error,
) {
	token, err := s.firebaseAuth.VerifyToken(ctx, firebaseToken)
	if err != nil {
		return "", "", 0, nil, fmt.Errorf("invalid firebase token: %w", err)
	}

	email, ok := token.Claims["email"].(string)
	if !ok || email == "" {
		return "", "", 0, nil, fmt.Errorf("email not found in token claims")
	}

	name, _ := token.Claims["name"].(string)
	if name == "" {
		name = email
	}

	user, err := s.userService.GetOrCreateUser(ctx, email, name)
	if err != nil {
		return "", "", 0, nil, fmt.Errorf("failed to get or create user: %w", err)
	}

	accessToken, expiryTime, err := s.jwtManager.GenerateAccessToken(user.ID.Hex(), email)
	if err != nil {
		return "", "", 0, nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, refreshExpiryTime, err := s.jwtManager.GenerateRefreshToken(user.ID.Hex(), email)
	if err != nil {
		return "", "", 0, nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	err = s.tokenStore.SaveRefreshToken(ctx, user.ID.Hex(), refreshToken, refreshExpiryTime)
	if err != nil {
		return "", "", 0, nil, fmt.Errorf("failed to save refresh token: %w", err)
	}

	expiresIn = int(time.Until(expiryTime).Seconds())

	return accessToken, refreshToken, expiresIn, user, nil
}

func (s *Service) RefreshToken(ctx context.Context, refreshToken string) (newAccessToken, newRefreshToken string, expiresIn int, err error) {
	if s.cfg.Firebase.BypassEnabled && refreshToken == s.cfg.Firebase.BypassRefreshToken {
		return s.cfg.Firebase.BypassToken, s.cfg.Firebase.BypassRefreshToken, 0, nil
	}

	claims, err := s.jwtManager.ValidateToken(refreshToken)
	if err != nil {
		return "", "", 0, fmt.Errorf("invalid refresh token: %w", err)
	}

	id, err := s.tokenStore.GetRefreshToken(ctx, refreshToken)
	if err != nil {
		return "", "", 0, fmt.Errorf("refresh token not found: %w", err)
	}

	if id != claims.ID {
		return "", "", 0, fmt.Errorf("token id mismatch")
	}

	accessToken, expiryTime, err := s.jwtManager.GenerateAccessToken(claims.ID, claims.Email)
	if err != nil {
		return "", "", 0, fmt.Errorf("failed to generate access token: %w", err)
	}

	newRefreshToken, refreshExpiryTime, err := s.jwtManager.GenerateRefreshToken(claims.ID, claims.Email)
	if err != nil {
		return "", "", 0, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	if err := s.tokenStore.DeleteRefreshToken(ctx, refreshToken); err != nil {
		return "", "", 0, fmt.Errorf("failed to delete old refresh token: %w", err)
	}

	if err := s.tokenStore.SaveRefreshToken(ctx, claims.ID, newRefreshToken, refreshExpiryTime); err != nil {
		return "", "", 0, fmt.Errorf("failed to save new refresh token: %w", err)
	}

	expiresIn = int(time.Until(expiryTime).Seconds())

	return accessToken, newRefreshToken, expiresIn, nil
}

func (s *Service) Logout(ctx context.Context, refreshToken string) error {
	err := s.tokenStore.DeleteRefreshToken(ctx, refreshToken)
	if err != nil {
		return fmt.Errorf("failed to delete refresh token: %w", err)
	}

	return nil
}
