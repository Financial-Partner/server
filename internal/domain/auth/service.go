package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/Financial-Partner/server/internal/config"
	"github.com/Financial-Partner/server/internal/entities"
	"github.com/Financial-Partner/server/internal/infrastructure/auth"
)

type TokenStore interface {
	SaveRefreshToken(ctx context.Context, email, refreshToken string, expiry time.Time) error
	GetRefreshToken(ctx context.Context, refreshToken string) (string, error)
	DeleteRefreshToken(ctx context.Context, refreshToken string) error
}

type UserService interface {
	GetUser(ctx context.Context, email string) (*entities.User, error)
	GetOrCreateUser(ctx context.Context, email, name string) (*entities.User, error)
}

type FirebaseAuth interface {
	VerifyToken(ctx context.Context, idToken string) (*auth.Token, error)
}

type JWTManager interface {
	GenerateAccessToken(email string) (string, time.Time, error)
	GenerateRefreshToken(email string) (string, time.Time, error)
	ValidateToken(tokenString string) (*auth.Claims, error)
}

type Service struct {
	cfg          *config.Config
	firebaseAuth FirebaseAuth
	jwtManager   JWTManager
	tokenStore   TokenStore
	userService  UserService
}

func NewService(
	cfg *config.Config,
	firebaseAuth FirebaseAuth,
	jwtManager JWTManager,
	tokenStore TokenStore,
	userService UserService,
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

	accessToken, expiryTime, err := s.jwtManager.GenerateAccessToken(email)
	if err != nil {
		return "", "", 0, nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, refreshExpiryTime, err := s.jwtManager.GenerateRefreshToken(email)
	if err != nil {
		return "", "", 0, nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	err = s.tokenStore.SaveRefreshToken(ctx, email, refreshToken, refreshExpiryTime)
	if err != nil {
		return "", "", 0, nil, fmt.Errorf("failed to save refresh token: %w", err)
	}

	expiresIn = int(time.Until(expiryTime).Seconds())

	return accessToken, refreshToken, expiresIn, user, nil
}

func (s *Service) RefreshToken(ctx context.Context, refreshToken string) (newAccessToken, newRefreshToken string, expiresIn int, err error) {
	claims, err := s.jwtManager.ValidateToken(refreshToken)
	if err != nil {
		return "", "", 0, fmt.Errorf("invalid refresh token: %w", err)
	}

	email, err := s.tokenStore.GetRefreshToken(ctx, refreshToken)
	if err != nil {
		return "", "", 0, fmt.Errorf("refresh token not found: %w", err)
	}

	if email != claims.Email {
		return "", "", 0, fmt.Errorf("token email mismatch")
	}

	accessToken, expiryTime, err := s.jwtManager.GenerateAccessToken(email)
	if err != nil {
		return "", "", 0, fmt.Errorf("failed to generate access token: %w", err)
	}

	newRefreshToken, refreshExpiryTime, err := s.jwtManager.GenerateRefreshToken(email)
	if err != nil {
		return "", "", 0, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	if err := s.tokenStore.DeleteRefreshToken(ctx, refreshToken); err != nil {
		return "", "", 0, fmt.Errorf("failed to delete old refresh token: %w", err)
	}

	if err := s.tokenStore.SaveRefreshToken(ctx, email, newRefreshToken, refreshExpiryTime); err != nil {
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
