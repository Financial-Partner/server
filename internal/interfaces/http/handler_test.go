package handler_test

import (
	"testing"

	"github.com/Financial-Partner/server/internal/infrastructure/logger"
	handler "github.com/Financial-Partner/server/internal/interfaces/http"
	"go.uber.org/mock/gomock"
)

type MockServices struct {
	UserService *handler.MockUserService
	AuthService *handler.MockAuthService
	GoalService *handler.MockGoalService
}

func newTestHandler(t *testing.T) (*handler.Handler, *MockServices) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ms := &MockServices{
		UserService: handler.NewMockUserService(ctrl),
		AuthService: handler.NewMockAuthService(ctrl),
		GoalService: handler.NewMockGoalService(ctrl),
	}
	h := handler.NewHandler(ms.UserService, ms.AuthService, ms.GoalService, logger.NewNopLogger())

	return h, ms
}
