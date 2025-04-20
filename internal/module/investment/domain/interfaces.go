package investment_domain

import (
	"context"

	"github.com/Financial-Partner/server/internal/entities"
	"github.com/Financial-Partner/server/internal/interfaces/http/dto"
)

//go:generate mockgen -source=interfaces.go -destination=interfaces_mock.go -package=investment_domain

type InvestmentService interface {
	GetOpportunities(ctx context.Context, userID string) ([]entities.Opportunity, error)
	CreateUserInvestment(ctx context.Context, userID string, req *dto.CreateUserInvestmentRequest) (*entities.Investment, error)
	GetUserInvestments(ctx context.Context, userID string) ([]entities.Investment, error)
}
