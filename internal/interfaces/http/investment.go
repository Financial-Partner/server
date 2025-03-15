package handler

import (
	"context"
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
	GetInvestments(ctx context.Context, userID string) ([]entities.Investment, error)
}

// @Summary Get investments
// @Description Get investments for a user
// @Tags investments
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer {token}" default "Bearer "
// @Success 200 {object} dto.GetInvestmentsResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /investments [get]
func (h *Handler) GetInvestments(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(contextutil.UserEmailKey).(string)
	if !ok {
		h.log.Warnf("failed to get user ID from context")
		responde.WithError(w, r, h.log, nil, httperror.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	investments, err := h.investmentService.GetInvestments(r.Context(), userID)
	if err != nil {
		h.log.WithError(err).Warnf("failed to get investments")
		responde.WithError(w, r, h.log, err, httperror.ErrFailedToGetInvestments, http.StatusInternalServerError)
		return
	}

	var investmentResponses []dto.InvestmentResponse
	for _, investment := range investments {
		investmentResponses = append(investmentResponses, dto.InvestmentResponse{
			Title:       investment.Title,
			Description: investment.Description,
			Tags:        investment.Tags,
			IsIncrease:  investment.IsIncrease,
			Variation:   investment.Variation,
			Duration:    investment.Duration,
			MinAmount:   investment.MinAmount,
			CreatedAt:   investment.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   investment.UpdatedAt.Format(time.RFC3339),
		})
	}

	resp := dto.GetInvestmentsResponse{
		Investments: investmentResponses,
	}

	responde.WithJSON(w, r, resp, http.StatusOK)
}
