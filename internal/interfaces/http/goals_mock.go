// Code generated by MockGen. DO NOT EDIT.
// Source: goals.go
//
// Generated by this command:
//
//	mockgen -source=goals.go -destination=goals_mock.go -package=handler
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

// MockGoalService is a mock of GoalService interface.
type MockGoalService struct {
	ctrl     *gomock.Controller
	recorder *MockGoalServiceMockRecorder
	isgomock struct{}
}

// MockGoalServiceMockRecorder is the mock recorder for MockGoalService.
type MockGoalServiceMockRecorder struct {
	mock *MockGoalService
}

// NewMockGoalService creates a new mock instance.
func NewMockGoalService(ctrl *gomock.Controller) *MockGoalService {
	mock := &MockGoalService{ctrl: ctrl}
	mock.recorder = &MockGoalServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGoalService) EXPECT() *MockGoalServiceMockRecorder {
	return m.recorder
}

// CreateGoal mocks base method.
func (m *MockGoalService) CreateGoal(ctx context.Context, userID string, req *dto.CreateGoalRequest) (*entities.Goal, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateGoal", ctx, userID, req)
	ret0, _ := ret[0].(*entities.Goal)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateGoal indicates an expected call of CreateGoal.
func (mr *MockGoalServiceMockRecorder) CreateGoal(ctx, userID, req any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateGoal", reflect.TypeOf((*MockGoalService)(nil).CreateGoal), ctx, userID, req)
}

// GetAutoGoalSuggestion mocks base method.
func (m *MockGoalService) GetAutoGoalSuggestion(ctx context.Context, userID string) (*entities.GoalSuggestion, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAutoGoalSuggestion", ctx, userID)
	ret0, _ := ret[0].(*entities.GoalSuggestion)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAutoGoalSuggestion indicates an expected call of GetAutoGoalSuggestion.
func (mr *MockGoalServiceMockRecorder) GetAutoGoalSuggestion(ctx, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAutoGoalSuggestion", reflect.TypeOf((*MockGoalService)(nil).GetAutoGoalSuggestion), ctx, userID)
}

// GetGoal mocks base method.
func (m *MockGoalService) GetGoal(ctx context.Context, userID string) (*entities.Goal, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGoal", ctx, userID)
	ret0, _ := ret[0].(*entities.Goal)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGoal indicates an expected call of GetGoal.
func (mr *MockGoalServiceMockRecorder) GetGoal(ctx, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGoal", reflect.TypeOf((*MockGoalService)(nil).GetGoal), ctx, userID)
}

// GetGoalSuggestion mocks base method.
func (m *MockGoalService) GetGoalSuggestion(ctx context.Context, userID string, req *dto.GoalSuggestionRequest) (*entities.GoalSuggestion, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGoalSuggestion", ctx, userID, req)
	ret0, _ := ret[0].(*entities.GoalSuggestion)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGoalSuggestion indicates an expected call of GetGoalSuggestion.
func (mr *MockGoalServiceMockRecorder) GetGoalSuggestion(ctx, userID, req any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGoalSuggestion", reflect.TypeOf((*MockGoalService)(nil).GetGoalSuggestion), ctx, userID, req)
}
