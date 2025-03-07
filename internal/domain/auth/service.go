package auth

import (
	"context"

	"github.com/Financial-Partner/server/internal/entities"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) LoginWithFirebase(ctx context.Context, firebaseToken string) (
	accessToken, refreshToken string, expiresIn int, userInfo *entities.User, err error,
) {
	return "", "", 0, nil, nil
}

func (s *Service) RefreshToken(ctx context.Context, refreshToken string) (newAccessToken string, expiresIn int, err error) {
	return "", 0, nil
}

func (s *Service) Logout(ctx context.Context, refreshToken string) error {
	return nil
}
