package user

import (
	"context"
)

//go:generate mockgen -source=repository.go -destination=mocks/user_repository_mock.go -package=mocks

type Repository interface {
	FindByEmail(ctx context.Context, email string) (*UserEntity, error)
	Create(ctx context.Context, entity *UserEntity) (*UserEntity, error)
	Update(ctx context.Context, entity *UserEntity) error
}

type UserStore interface {
	Get(ctx context.Context, email string) (*UserEntity, error)
	Set(ctx context.Context, entity *UserEntity) error
	Delete(ctx context.Context, email string) error
}
