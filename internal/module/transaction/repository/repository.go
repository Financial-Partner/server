package transaction_repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Financial-Partner/server/internal/entities"
)

//go:generate mockgen -source=repository.go -destination=repository_mock.go -package=transaction_repository

type Repository interface {
	Create(ctx context.Context, transaction *entities.Transaction) (*entities.Transaction, error)
	FindByUserId(ctx context.Context, userID primitive.ObjectID) ([]entities.Transaction, error)
}

type TransactionStore interface {
	GetByUserId(ctx context.Context, userID string) ([]entities.Transaction, error)
	SetByUserId(ctx context.Context, userID string, transactions []entities.Transaction) error
	DeleteByUserId(ctx context.Context, userID string) error
}
