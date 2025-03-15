package handler_test

import (
	"testing"

	"github.com/Financial-Partner/server/internal/infrastructure/logger"
	handler "github.com/Financial-Partner/server/internal/interfaces/http"
	"github.com/Financial-Partner/server/internal/interfaces/http/mocks"
	"go.uber.org/mock/gomock"
)

type MockServices struct {
	UserService       *mocks.MockUserService
	AuthService       *mocks.MockAuthService
	GoalService       *mocks.MockGoalService
	InvestmentService *mocks.MockInvestmentService
}

func newTestHandler(t *testing.T) (*handler.Handler, *MockServices) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ms := &MockServices{
		UserService:       mocks.NewMockUserService(ctrl),
		AuthService:       mocks.NewMockAuthService(ctrl),
		GoalService:       mocks.NewMockGoalService(ctrl),
		InvestmentService: mocks.NewMockInvestmentService(ctrl),
	}
	h := handler.NewHandler(ms.UserService, ms.AuthService, ms.GoalService, ms.InvestmentService, logger.NewNopLogger())

	return h, ms
}
