package user_domain

import (
	"context"

	"github.com/Financial-Partner/server/internal/entities"
)

//go:generate mockgen -source=interfaces.go -destination=interfaces_mock.go -package=user_domain

type UserService interface {
	GetUser(ctx context.Context, email string) (*entities.User, error)
	GetOrCreateUser(ctx context.Context, email, name string) (*entities.User, error)
	UpdateUserName(ctx context.Context, email, name string) (*entities.User, error)
}
