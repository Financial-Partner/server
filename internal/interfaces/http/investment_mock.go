// Code generated by MockGen. DO NOT EDIT.
// Source: investment.go
//
// Generated by this command:
//
//	mockgen -source=investment.go -destination=investment_mock.go -package=handler
//

// Package handler is a generated GoMock package.
package handler

import (
	context "context"
	reflect "reflect"

	entities "github.com/Financial-Partner/server/internal/entities"
	dto "github.com/Financial-Partner/server/internal/interfaces/http/dto"
	gomock "go.uber.org/mock/gomock"
)

// MockInvestmentService is a mock of InvestmentService interface.
type MockInvestmentService struct {
	ctrl     *gomock.Controller
	recorder *MockInvestmentServiceMockRecorder
	isgomock struct{}
}

// MockInvestmentServiceMockRecorder is the mock recorder for MockInvestmentService.
type MockInvestmentServiceMockRecorder struct {
	mock *MockInvestmentService
}

// NewMockInvestmentService creates a new mock instance.
func NewMockInvestmentService(ctrl *gomock.Controller) *MockInvestmentService {
	mock := &MockInvestmentService{ctrl: ctrl}
	mock.recorder = &MockInvestmentServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockInvestmentService) EXPECT() *MockInvestmentServiceMockRecorder {
	return m.recorder
}

// GetOpportunities mocks base method.
func (m *MockInvestmentService) GetOpportunities(ctx context.Context, userID string) ([]entities.Opportunity, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOpportunities", ctx, userID)
	ret0, _ := ret[0].([]entities.Opportunity)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOpportunities indicates an expected call of GetOpportunities.
func (mr *MockInvestmentServiceMockRecorder) GetOpportunities(ctx, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOpportunities", reflect.TypeOf((*MockInvestmentService)(nil).GetOpportunities), ctx, userID)
}

// GetUserInvestments mocks base method.
func (m *MockInvestmentService) GetUserInvestments(ctx context.Context, userID string) ([]entities.Investment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserInvestments", ctx, userID)
	ret0, _ := ret[0].([]entities.Investment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetInvestments indicates an expected call of GetInvestments.
func (mr *MockInvestmentServiceMockRecorder) GetUserInvestments(ctx, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserInvestments", reflect.TypeOf((*MockInvestmentService)(nil).GetUserInvestments), ctx, userID)
}

// CreateUserInvestment mocks base method.
func (m *MockInvestmentService) CreateUserInvestment(ctx context.Context, userID string, req *dto.CreateUserInvestmentRequest) (*entities.Investment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUserInvestment", ctx, userID, req)
	ret0, _ := ret[0].(*entities.Investment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUserInvestment indicates an expected call of CreateUserInvestment.
func (mr *MockInvestmentServiceMockRecorder) CreateUserInvestment(ctx, userID, req any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUserInvestment", reflect.TypeOf((*MockInvestmentService)(nil).CreateUserInvestment), ctx, userID, req)
}
