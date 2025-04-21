package investment_usecase

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

func (s *Service) GetOpportunities(ctx context.Context, userID string) ([]entities.Opportunity, error) {
	return nil, nil
}

func (s *Service) CreateUserInvestment(ctx context.Context, userID string, req *dto.CreateUserInvestmentRequest) (*entities.Investment, error) {
	return nil, nil
}

func (s *Service) GetUserInvestments(ctx context.Context, userID string) ([]entities.Investment, error) {
	return nil, nil
}

func (s *Service) CreateOpportunity(ctx context.Context, userID string, req *dto.CreateOpportunityRequest) (*entities.Opportunity, error) {
	return nil, nil
}
