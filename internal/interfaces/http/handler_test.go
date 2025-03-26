package handler_test

import (
	"context"
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/Financial-Partner/server/internal/contextutil"
	"github.com/Financial-Partner/server/internal/infrastructure/logger"
	handler "github.com/Financial-Partner/server/internal/interfaces/http"
)

type MockServices struct {
	UserService        *handler.MockUserService
	AuthService        *handler.MockAuthService
	GoalService        *handler.MockGoalService
	InvestmentService  *handler.MockInvestmentService
	TransactionService *handler.MockTransactionService
}

func newTestHandler(t *testing.T) (*handler.Handler, *MockServices) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ms := &MockServices{
		UserService:        handler.NewMockUserService(ctrl),
		AuthService:        handler.NewMockAuthService(ctrl),
		GoalService:        handler.NewMockGoalService(ctrl),
		InvestmentService:  handler.NewMockInvestmentService(ctrl),
		TransactionService: handler.NewMockTransactionService(ctrl),
	}
	h := handler.NewHandler(ms.UserService, ms.AuthService, ms.GoalService, ms.InvestmentService, ms.TransactionService, logger.NewNopLogger())

	return h, ms
}

func newContext(userID, userEmail string) context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, contextutil.UserIDKey, userID)
	ctx = context.WithValue(ctx, contextutil.UserEmailKey, userEmail)

	return ctx
}
