package transaction_usecase

import (
	"context"

	"github.com/Financial-Partner/server/internal/entities"
	"github.com/Financial-Partner/server/internal/interfaces/http/dto"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) CreateTransaction(ctx context.Context, userID string, req *dto.CreateTransactionRequest) (*entities.Transaction, error) {
	return nil, nil
}

func (s *Service) GetTransactions(ctx context.Context, userID string) ([]entities.Transaction, error) {
	return nil, nil
}
