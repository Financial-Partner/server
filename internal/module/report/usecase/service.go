package report_usecase

import (
	"context"
	"time"

	"github.com/Financial-Partner/server/internal/entities"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) GetReport(ctx context.Context, userID string, startTime time.Time, endTime time.Time, reportType string) (*entities.Report, error) {
	return nil, nil
}

func (s *Service) GetReportSummary(ctx context.Context, userID string) (*entities.ReportSummary, error) {
	return nil, nil
}
