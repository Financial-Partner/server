package investment_usecase

import (
	"context"

	"github.com/Financial-Partner/server/internal/entities"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) GetInvestments(ctx context.Context, userID string) ([]entities.Investment, error) {
	return nil, nil
}
