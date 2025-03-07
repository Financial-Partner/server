package goal

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

func (s *Service) GetGoalSuggestion(ctx context.Context, userID string, req *dto.GoalSuggestionRequest) (*entities.GoalSuggestion, error) {
	return nil, nil
}

func (s *Service) GetAutoGoalSuggestion(ctx context.Context, userID string) (*entities.GoalSuggestion, error) {
	return nil, nil
}

func (s *Service) CreateGoal(ctx context.Context, userID string, req *dto.CreateGoalRequest) (*entities.Goal, error) {
	return nil, nil
}

func (s *Service) GetGoal(ctx context.Context, userID string) (*entities.Goal, error) {
	return nil, nil
}
