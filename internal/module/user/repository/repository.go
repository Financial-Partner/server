package user_repository

import (
	"context"

	"github.com/Financial-Partner/server/internal/entities"
)

//go:generate mockgen -source=repository.go -destination=user_repository_mock.go -package=user_repository

type Repository interface {
	FindByEmail(ctx context.Context, email string) (*entities.User, error)
	Create(ctx context.Context, entity *entities.User) (*entities.User, error)
	Update(ctx context.Context, entity *entities.User) error
}

type UserStore interface {
	Get(ctx context.Context, email string) (*entities.User, error)
	Set(ctx context.Context, entity *entities.User) error
	Delete(ctx context.Context, email string) error
}
