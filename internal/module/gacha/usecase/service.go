package gacha_usecase

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

func (s *Service) PreviewGachas(ctx context.Context, userID string) ([]entities.Gacha, error) {
	return nil, nil
}

func (s *Service) DrawGacha(ctx context.Context, userID string, req *dto.DrawGachaRequest) (*entities.Gacha, error) {
	return nil, nil
}
