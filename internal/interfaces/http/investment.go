package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Financial-Partner/server/internal/contextutil"
	"github.com/Financial-Partner/server/internal/entities"
	"github.com/Financial-Partner/server/internal/interfaces/http/dto"
	httperror "github.com/Financial-Partner/server/internal/interfaces/http/error"
	responde "github.com/Financial-Partner/server/internal/interfaces/http/respond"
)

//go:generate mockgen -source=investment.go -destination=investment_mock.go -package=handler

type InvestmentService interface {
	GetOpportunities(ctx context.Context, userID string) ([]entities.Opportunity, error)
	CreateUserInvestment(ctx context.Context, userID string, req *dto.CreateUserInvestmentRequest) (*entities.Investment, error)
	GetUserInvestments(ctx context.Context, userID string) ([]entities.Investment, error)
}

// @Summary Get investment opportunities
// @Description Get investment opportunities for a user
// @Tags investments
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer {token}" default "Bearer "
// @Success 200 {object} dto.GetOpportunitiesResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /investments [get]
func (h *Handler) GetOpportunities(w http.ResponseWriter, r *http.Request) {
	userID, ok := contextutil.GetUserID(r.Context())
	if !ok {
		h.log.Warnf("failed to get user ID from context")
		responde.WithError(w, r, h.log, nil, httperror.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	opportunities, err := h.investmentService.GetOpportunities(r.Context(), userID)
	if err != nil {
		h.log.WithError(err).Warnf("failed to get investment opportunities")
		responde.WithError(w, r, h.log, err, httperror.ErrFailedToGetOpportunities, http.StatusInternalServerError)
		return
	}

	var opportunitiesResponses []dto.OpportunityResponse
	for _, opportunity := range opportunities {
		opportunitiesResponses = append(opportunitiesResponses, dto.OpportunityResponse{
			Title:       opportunity.Title,
			Description: opportunity.Description,
			Tags:        opportunity.Tags,
			IsIncrease:  opportunity.IsIncrease,
			Variation:   opportunity.Variation,
			Duration:    opportunity.Duration,
			MinAmount:   opportunity.MinAmount,
			CreatedAt:   opportunity.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   opportunity.UpdatedAt.Format(time.RFC3339),
		})
	}

	resp := dto.GetOpportunitiesResponse{
		Opportunities: opportunitiesResponses,
	}

	responde.WithJSON(w, r, resp, http.StatusOK)
}

// @Summary Create user investment
// @Description Create investment for a user
// @Tags investments
// @Accept json
// @Produce json
// @Param request body dto.CreateUserInvestmentRequest true "Create user investment request"
// @Param Authorization header string true "Bearer {token}" default "Bearer "
// @Success 201 {object} dto.CreateUserInvestmentResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users/me/investments [post]
func (h *Handler) CreateUserInvestment(w http.ResponseWriter, r *http.Request) {
	userID, ok := contextutil.GetUserID(r.Context())
	if !ok {
		h.log.Warnf("failed to get user ID from context")
		responde.WithError(w, r, h.log, nil, httperror.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	var req dto.CreateUserInvestmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.WithError(err).Warnf("failed to decode request body")
		responde.WithError(w, r, h.log, err, httperror.ErrInvalidRequest, http.StatusBadRequest)
		return
	}

	investment, err := h.investmentService.CreateUserInvestment(r.Context(), userID, &req)
	if err != nil {
		h.log.WithError(err).Warnf("failed to create an user investment")
		responde.WithError(w, r, h.log, err, httperror.ErrFailedToCreateUserInvestment, http.StatusInternalServerError)
		return
	}

	resp := dto.InvestmentResponse{
		OpportunityID: investment.OpportunityID,
		UserID:        investment.UserID,
		Amount:        investment.Amount,
		CreatedAt:     investment.CreatedAt.Format(time.RFC3339),
		UpdatedAt:     investment.UpdatedAt.Format(time.RFC3339),
	}

	responde.WithJSON(w, r, resp, http.StatusOK)
}

// @Summary Get user investments
// @Description Get user investments
// @Tags investments
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer {token}" default "Bearer "
// @Success 200 {object} dto.GetUserInvestmentsResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users/me/investments [get]
func (h *Handler) GetUserInvestments(w http.ResponseWriter, r *http.Request) {
	userID, ok := contextutil.GetUserID(r.Context())
	if !ok {
		h.log.Warnf("failed to get user ID from context")
		responde.WithError(w, r, h.log, nil, httperror.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	investments, err := h.investmentService.GetUserInvestments(r.Context(), userID)
	if err != nil {
		h.log.WithError(err).Warnf("failed to get user investments")
		responde.WithError(w, r, h.log, err, httperror.ErrFailedToGetUserInvestments, http.StatusInternalServerError)
		return
	}

	var investmentsResponse []dto.InvestmentResponse
	for _, investment := range investments {
		investmentsResponse = append(investmentsResponse, dto.InvestmentResponse{
			OpportunityID: investment.OpportunityID,
			Amount:        investment.Amount,
			CreatedAt:     investment.CreatedAt.Format(time.RFC3339),
			UpdatedAt:     investment.UpdatedAt.Format(time.RFC3339),
		})
	}

	resp := dto.GetUserInvestmentsResponse{
		Investments: investmentsResponse,
	}

	responde.WithJSON(w, r, resp, http.StatusOK)
}
