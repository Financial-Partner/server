package investment_repository

import (
	"context"

	"github.com/Financial-Partner/server/internal/entities"
)

//go:generate mockgen -source=repository.go -destination=repository_mock.go -package=investment_repository

type Repository interface {
	CreateInvestment(ctx context.Context, entity *entities.Investment) (*entities.Investment, error)
	FindOpportunitiesByUserId(ctx context.Context, userID string) ([]entities.Opportunity, error)
	FindInvestmentsByUserId(ctx context.Context, userID string) ([]entities.Investment, error)
}

type InvestmentStore interface {
	GetOpportunities(ctx context.Context, userID string) ([]entities.Opportunity, error)
	GetInvestments(ctx context.Context, userID string) ([]entities.Investment, error)
	SetOpportunities(ctx context.Context, userID string, opportunities []entities.Opportunity) error
	SetInvestments(ctx context.Context, userID string, investments []entities.Investment) error
	DeleteInvestments(ctx context.Context, userID string) error
	DeleteOpportunities(ctx context.Context, userID string) error
}
