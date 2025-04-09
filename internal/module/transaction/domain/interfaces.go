package transaction_domain

import (
	"context"

	"github.com/Financial-Partner/server/internal/entities"
	"github.com/Financial-Partner/server/internal/interfaces/http/dto"
)

//go:generate mockgen -source=interfaces.go -destination=interfaces_mock.go -package=transaction_domain

type TransactionService interface {
	CreateTransaction(ctx context.Context, UserID string, transaction *dto.CreateTransactionRequest) (*entities.Transaction, error)
	GetTransactions(ctx context.Context, UserId string) ([]entities.Transaction, error)
}
