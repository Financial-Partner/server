package auth_domain

import (
	"context"
	"time"

	"github.com/Financial-Partner/server/internal/infrastructure/auth"
)

//go:generate mockgen -source=interfaces.go -destination=interfaces_mock.go -package=auth_domain

type TokenStore interface {
	SaveRefreshToken(ctx context.Context, id, refreshToken string, expiry time.Time) error
	GetRefreshToken(ctx context.Context, refreshToken string) (string, error)
	DeleteRefreshToken(ctx context.Context, refreshToken string) error
}

type FirebaseAuth interface {
	VerifyToken(ctx context.Context, idToken string) (*auth.Token, error)
}

type JWTManager interface {
	GenerateAccessToken(id, email string) (string, time.Time, error)
	GenerateRefreshToken(id, email string) (string, time.Time, error)
	ValidateToken(tokenString string) (*auth.Claims, error)
}
