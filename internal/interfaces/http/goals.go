package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Financial-Partner/server/internal/contextutil"
	"github.com/Financial-Partner/server/internal/entities"
	"github.com/Financial-Partner/server/internal/interfaces/http/dto"
)

//go:generate mockgen -source=goals.go -destination=mocks/goals_mock.go -package=mocks

type GoalService interface {
	GetGoalSuggestion(ctx context.Context, userID string, req *dto.GoalSuggestionRequest) (*entities.GoalSuggestion, error)
	GetAutoGoalSuggestion(ctx context.Context, userID string) (*entities.GoalSuggestion, error)
	CreateGoal(ctx context.Context, userID string, req *dto.CreateGoalRequest) (*entities.Goal, error)
	GetGoal(ctx context.Context, userID string) (*entities.Goal, error)
}

// @Summary Calculate and return suggested saving goals based on user's input expense data
// @Description Calculate and return suggested saving goals based on user's input expense data
// @Tags goals
// @Accept json
// @Produce json
// @Param request body dto.GoalSuggestionRequest true "Goal suggestion request"
// @Param Authorization header string true "Bearer {token}" default "Bearer "
// @Success 200 {object} dto.GoalSuggestionResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /goals/suggestion [post]
func (h *Handler) GetGoalSuggestion(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(contextutil.UserEmailKey).(string)
	if !ok {
		h.log.Warnf("failed to get user ID from context")
		h.RespondWithError(w, r, h.log, nil, ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	var req dto.GoalSuggestionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.WithError(err).Warnf("failed to decode request body")
		h.RespondWithError(w, r, h.log, err, ErrInvalidRequest, http.StatusBadRequest)
		return
	}

	suggestion, err := h.goalService.GetGoalSuggestion(r.Context(), userID, &req)
	if err != nil {
		h.log.WithError(err).Warnf("failed to get goal suggestion")
		h.RespondWithError(w, r, h.log, err, ErrFailedToGetGoalSuggestion, http.StatusInternalServerError)
		return
	}

	resp := dto.GoalSuggestionResponse{
		SuggestedAmount: suggestion.SuggestedAmount,
		Period:          suggestion.Period,
		Message:         suggestion.Message,
	}

	h.RespondWithJSON(w, r, resp, http.StatusOK)
}

// @Summary Calculate and return suggested saving goals based on user's expense data
// @Description Calculate and return suggested saving goals based on user's expense data
// @Tags goals
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer {token}" default "Bearer "
// @Success 200 {object} dto.GoalSuggestionResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /goals/suggestion/me [get]
func (h *Handler) GetAutoGoalSuggestion(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(contextutil.UserEmailKey).(string)
	if !ok {
		h.log.Warnf("failed to get user ID from context")
		h.RespondWithError(w, r, h.log, nil, ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	suggestion, err := h.goalService.GetAutoGoalSuggestion(r.Context(), userID)
	if err != nil {
		h.log.WithError(err).Warnf("failed to get auto goal suggestion")
		h.RespondWithError(w, r, h.log, err, ErrFailedToGetGoalSuggestion, http.StatusInternalServerError)
		return
	}

	resp := dto.GoalSuggestionResponse{
		SuggestedAmount: suggestion.SuggestedAmount,
		Period:          suggestion.Period,
		Message:         suggestion.Message,
	}

	h.RespondWithJSON(w, r, resp, http.StatusOK)
}

// @Summary Create user's saving goal
// @Description Set user's saving goal amount and period
// @Tags goals
// @Accept json
// @Produce json
// @Param request body dto.CreateGoalRequest true "Create goal request"
// @Param Authorization header string true "Bearer {token}" default "Bearer "
// @Success 200 {object} dto.GoalResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /goals [post]
func (h *Handler) CreateGoal(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(contextutil.UserEmailKey).(string)
	if !ok {
		h.log.Warnf("failed to get user ID from context")
		h.RespondWithError(w, r, h.log, nil, ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	var req dto.CreateGoalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.WithError(err).Warnf("failed to decode request body")
		h.RespondWithError(w, r, h.log, err, ErrInvalidRequest, http.StatusBadRequest)
		return
	}

	goal, err := h.goalService.CreateGoal(r.Context(), userID, &req)
	if err != nil {
		h.log.WithError(err).Warnf("failed to create or update goal")
		h.RespondWithError(w, r, h.log, err, ErrFailedToCreateGoal, http.StatusInternalServerError)
		return
	}

	resp := dto.GoalResponse{
		TargetAmount:  goal.TargetAmount,
		CurrentAmount: goal.CurrentAmount,
		Period:        goal.Period,
		Status:        goal.Status,
		CreatedAt:     goal.CreatedAt.Format(time.RFC3339),
		UpdatedAt:     goal.UpdatedAt.Format(time.RFC3339),
	}

	h.RespondWithJSON(w, r, resp, http.StatusOK)
}

// @Summary Get current saving goal
// @Description Get user's current saving goal and status
// @Tags goals
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer {token}" default "Bearer "
// @Success 200 {object} dto.GetGoalResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /goals [get]
func (h *Handler) GetGoal(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(contextutil.UserEmailKey).(string)
	if !ok {
		h.log.Warnf("failed to get user ID from context")
		h.RespondWithError(w, r, h.log, nil, ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	goal, err := h.goalService.GetGoal(r.Context(), userID)
	if err != nil {
		h.log.WithError(err).Warnf("failed to get goal")
		h.RespondWithError(w, r, h.log, err, ErrFailedToGetGoal, http.StatusInternalServerError)
		return
	}

	resp := dto.GetGoalResponse{
		Goal: dto.GoalResponse{
			TargetAmount:  goal.TargetAmount,
			CurrentAmount: goal.CurrentAmount,
			Period:        goal.Period,
			Status:        goal.Status,
			CreatedAt:     goal.CreatedAt.Format(time.RFC3339),
			UpdatedAt:     goal.UpdatedAt.Format(time.RFC3339),
		},
	}

	h.RespondWithJSON(w, r, resp, http.StatusOK)
}
