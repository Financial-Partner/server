package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Financial-Partner/server/internal/contextutil"
	"github.com/Financial-Partner/server/internal/entities"
	"github.com/Financial-Partner/server/internal/interfaces/http/dto"
	httperror "github.com/Financial-Partner/server/internal/interfaces/http/error"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestGetGoalSuggestion(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("Invalid request format", func(t *testing.T) {
		h, _ := newTestHandler(t)

		invalidBody := bytes.NewBufferString(`{invalid json`)
		ctx := context.WithValue(context.Background(), contextutil.UserEmailKey, "test@example.com")
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/goals/suggestion", invalidBody)
		r = r.WithContext(ctx)

		h.GetGoalSuggestion(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var errorResp dto.ErrorResponse
		err := json.NewDecoder(w.Body).Decode(&errorResp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, errorResp.Code)
		assert.Equal(t, httperror.ErrInvalidRequest, errorResp.Message)
	})

	t.Run("Unauthorized request", func(t *testing.T) {
		h, _ := newTestHandler(t)

		req := dto.GoalSuggestionRequest{
			DailyExpenses:   1000,
			DailyIncome:     2000,
			WeeklyExpenses:  7000,
			WeeklyIncome:    14000,
			MonthlyExpenses: 30000,
			MonthlyIncome:   50000,
		}
		body, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/goals/suggestion", bytes.NewBuffer(body))

		h.GetGoalSuggestion(w, r)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var errorResp dto.ErrorResponse
		err := json.NewDecoder(w.Body).Decode(&errorResp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, errorResp.Code)
		assert.Equal(t, httperror.ErrUnauthorized, errorResp.Message)
	})

	t.Run("Service error", func(t *testing.T) {
		h, mockServices := newTestHandler(t)

		mockServices.GoalService.EXPECT().
			GetGoalSuggestion(gomock.Any(), "test@example.com", gomock.Any()).
			Return(nil, errors.New("service error"))

		req := dto.GoalSuggestionRequest{
			DailyExpenses:   1000,
			DailyIncome:     2000,
			WeeklyExpenses:  7000,
			WeeklyIncome:    14000,
			MonthlyExpenses: 30000,
			MonthlyIncome:   50000,
		}
		body, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/goals/suggestion", bytes.NewBuffer(body))
		ctx := context.WithValue(r.Context(), contextutil.UserEmailKey, "test@example.com")
		r = r.WithContext(ctx)

		h.GetGoalSuggestion(w, r)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var errorResp dto.ErrorResponse
		err := json.NewDecoder(w.Body).Decode(&errorResp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, errorResp.Code)
		assert.Equal(t, httperror.ErrFailedToGetGoalSuggestion, errorResp.Message)
	})

	t.Run("Success", func(t *testing.T) {
		h, mockServices := newTestHandler(t)

		suggestion := &entities.GoalSuggestion{
			SuggestedAmount: 15000,
			Period:          30,
			Message:         "Based on your income and expense analysis, we recommend that you can save 15,000 yuan per month.",
		}

		mockServices.GoalService.EXPECT().
			GetGoalSuggestion(gomock.Any(), "test@example.com", gomock.Any()).
			Return(suggestion, nil)

		req := dto.GoalSuggestionRequest{
			DailyExpenses:   1000,
			DailyIncome:     2000,
			WeeklyExpenses:  7000,
			WeeklyIncome:    14000,
			MonthlyExpenses: 30000,
			MonthlyIncome:   50000,
		}
		body, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/goals/suggestion", bytes.NewBuffer(body))
		ctx := context.WithValue(r.Context(), contextutil.UserEmailKey, "test@example.com")
		r = r.WithContext(ctx)

		h.GetGoalSuggestion(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var response dto.GoalSuggestionResponse
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, suggestion.SuggestedAmount, response.SuggestedAmount)
		assert.Equal(t, suggestion.Period, response.Period)
		assert.Equal(t, suggestion.Message, response.Message)
	})
}

func TestGetAutoGoalSuggestion(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("Unauthorized request", func(t *testing.T) {
		h, _ := newTestHandler(t)

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/goals/suggestion/me", nil)

		h.GetAutoGoalSuggestion(w, r)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var errorResp dto.ErrorResponse
		err := json.NewDecoder(w.Body).Decode(&errorResp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, errorResp.Code)
		assert.Equal(t, httperror.ErrUnauthorized, errorResp.Message)
	})

	t.Run("Service error", func(t *testing.T) {
		h, mockServices := newTestHandler(t)

		mockServices.GoalService.EXPECT().
			GetAutoGoalSuggestion(gomock.Any(), "test@example.com").
			Return(nil, errors.New("service error"))

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/goals/suggestion/me", nil)
		ctx := context.WithValue(r.Context(), contextutil.UserEmailKey, "test@example.com")
		r = r.WithContext(ctx)

		h.GetAutoGoalSuggestion(w, r)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var errorResp dto.ErrorResponse
		err := json.NewDecoder(w.Body).Decode(&errorResp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, errorResp.Code)
		assert.Equal(t, httperror.ErrFailedToGetGoalSuggestion, errorResp.Message)
	})

	t.Run("Success", func(t *testing.T) {
		h, mockServices := newTestHandler(t)

		suggestion := &entities.GoalSuggestion{
			SuggestedAmount: 15000,
			Period:          30,
			Message:         "Based on your past financial data analysis, we recommend that you can save 15,000 yuan per month.",
		}

		mockServices.GoalService.EXPECT().
			GetAutoGoalSuggestion(gomock.Any(), "test@example.com").
			Return(suggestion, nil)

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/goals/suggestion/me", nil)
		ctx := context.WithValue(r.Context(), contextutil.UserEmailKey, "test@example.com")
		r = r.WithContext(ctx)

		h.GetAutoGoalSuggestion(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var response dto.GoalSuggestionResponse
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, suggestion.SuggestedAmount, response.SuggestedAmount)
		assert.Equal(t, suggestion.Period, response.Period)
		assert.Equal(t, suggestion.Message, response.Message)
	})
}

func TestCreateGoal(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("Invalid request format", func(t *testing.T) {
		h, _ := newTestHandler(t)

		invalidBody := bytes.NewBufferString(`{invalid json`)
		ctx := context.WithValue(context.Background(), contextutil.UserEmailKey, "test@example.com")
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/goals", invalidBody)
		r = r.WithContext(ctx)

		h.CreateGoal(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var errorResp dto.ErrorResponse
		err := json.NewDecoder(w.Body).Decode(&errorResp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, errorResp.Code)
		assert.Equal(t, httperror.ErrInvalidRequest, errorResp.Message)
	})

	t.Run("Unauthorized request", func(t *testing.T) {
		h, _ := newTestHandler(t)

		req := dto.CreateGoalRequest{
			TargetAmount: 10000,
			Period:       30,
		}
		body, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/goals", bytes.NewBuffer(body))

		h.CreateGoal(w, r)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var errorResp dto.ErrorResponse
		err := json.NewDecoder(w.Body).Decode(&errorResp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, errorResp.Code)
		assert.Equal(t, httperror.ErrUnauthorized, errorResp.Message)
	})

	t.Run("Service error", func(t *testing.T) {
		h, mockServices := newTestHandler(t)

		mockServices.GoalService.EXPECT().
			CreateGoal(gomock.Any(), "test@example.com", gomock.Any()).
			Return(nil, errors.New("service error"))

		req := dto.CreateGoalRequest{
			TargetAmount: 10000,
			Period:       30,
		}
		body, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/goals", bytes.NewBuffer(body))
		ctx := context.WithValue(r.Context(), contextutil.UserEmailKey, "test@example.com")
		r = r.WithContext(ctx)

		h.CreateGoal(w, r)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var errorResp dto.ErrorResponse
		err := json.NewDecoder(w.Body).Decode(&errorResp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, errorResp.Code)
		assert.Equal(t, httperror.ErrFailedToCreateGoal, errorResp.Message)
	})

	t.Run("Success", func(t *testing.T) {
		h, mockServices := newTestHandler(t)

		now := time.Now()
		goal := &entities.Goal{
			ID:            "goal_123456",
			UserID:        "test@example.com",
			TargetAmount:  10000,
			CurrentAmount: 0,
			Period:        30,
			Status:        "",
			CreatedAt:     now,
			UpdatedAt:     now,
		}

		mockServices.GoalService.EXPECT().
			CreateGoal(gomock.Any(), "test@example.com", gomock.Any()).
			Return(goal, nil)

		req := dto.CreateGoalRequest{
			TargetAmount: 10000,
			Period:       30,
		}
		body, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/goals", bytes.NewBuffer(body))
		ctx := context.WithValue(r.Context(), contextutil.UserEmailKey, "test@example.com")
		r = r.WithContext(ctx)

		h.CreateGoal(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var response dto.GoalResponse
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, goal.TargetAmount, response.TargetAmount)
		assert.Equal(t, goal.CurrentAmount, response.CurrentAmount)
		assert.Equal(t, goal.Period, response.Period)
		assert.Equal(t, goal.Status, response.Status)
		assert.Equal(t, goal.CreatedAt.Format(time.RFC3339), response.CreatedAt)
		assert.Equal(t, goal.UpdatedAt.Format(time.RFC3339), response.UpdatedAt)
	})
}

func TestGetGoal(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("Unauthorized request", func(t *testing.T) {
		h, _ := newTestHandler(t)

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/goals", nil)

		h.GetGoal(w, r)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var errorResp dto.ErrorResponse
		err := json.NewDecoder(w.Body).Decode(&errorResp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, errorResp.Code)
		assert.Equal(t, httperror.ErrUnauthorized, errorResp.Message)
	})

	t.Run("Service error", func(t *testing.T) {
		h, mockServices := newTestHandler(t)

		mockServices.GoalService.EXPECT().
			GetGoal(gomock.Any(), "test@example.com").
			Return(nil, errors.New("service error"))

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/goals", nil)
		ctx := context.WithValue(r.Context(), contextutil.UserEmailKey, "test@example.com")
		r = r.WithContext(ctx)

		h.GetGoal(w, r)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var errorResp dto.ErrorResponse
		err := json.NewDecoder(w.Body).Decode(&errorResp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, errorResp.Code)
		assert.Equal(t, httperror.ErrFailedToGetGoal, errorResp.Message)
	})

	t.Run("Success", func(t *testing.T) {
		h, mockServices := newTestHandler(t)

		now := time.Now()
		goal := &entities.Goal{
			ID:            "goal_123456",
			UserID:        "test@example.com",
			TargetAmount:  10000,
			CurrentAmount: 5000,
			Period:        30,
			Status:        "Good progress",
			CreatedAt:     now.AddDate(0, -1, 0),
			UpdatedAt:     now,
		}

		mockServices.GoalService.EXPECT().
			GetGoal(gomock.Any(), "test@example.com").
			Return(goal, nil)

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/goals", nil)
		ctx := context.WithValue(r.Context(), contextutil.UserEmailKey, "test@example.com")
		r = r.WithContext(ctx)

		h.GetGoal(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var response dto.GetGoalResponse
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, goal.TargetAmount, response.Goal.TargetAmount)
		assert.Equal(t, goal.CurrentAmount, response.Goal.CurrentAmount)
		assert.Equal(t, goal.Period, response.Goal.Period)
		assert.Equal(t, goal.Status, response.Goal.Status)
		assert.Equal(t, goal.CreatedAt.Format(time.RFC3339), response.Goal.CreatedAt)
		assert.Equal(t, goal.UpdatedAt.Format(time.RFC3339), response.Goal.UpdatedAt)
	})
}
