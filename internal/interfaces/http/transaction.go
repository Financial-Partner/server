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

//go:generate mockgen -source=transaction.go -destination=transaction_mock.go -package=handler

type TransactionService interface {
	CreateTransaction(ctx context.Context, UserID string, transaction *dto.CreateTransactionRequest) (*entities.Transaction, error)
	GetTransactions(ctx context.Context, UserId string) ([]entities.Transaction, error)
}

// @Summary Get transactions
// @Description Get transactions for a user
// @Tags transactions
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer {token}" default "Bearer "
// @Success 200 {object} dto.GetTransactionsResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /transactions [get]
func (h *Handler) GetTransactions(w http.ResponseWriter, r *http.Request) {
	userID, ok := contextutil.GetUserID(r.Context())
	if !ok {
		h.log.Warnf("failed to get user ID from context")
		responde.WithError(w, r, h.log, nil, httperror.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	transactions, err := h.transactionService.GetTransactions(r.Context(), userID)
	if err != nil {
		h.log.Errorf("failed to get transactions")
		responde.WithError(w, r, h.log, err, httperror.ErrFailedToGetTransactions, http.StatusInternalServerError)
		return
	}

	var transactionResponses []dto.TransactionResponse
	for _, transaction := range transactions {
		transactionResponses = append(transactionResponses, dto.TransactionResponse{
			Amount:      transaction.Amount,
			Category:    transaction.Category,
			Type:        transaction.Type,
			Date:        transaction.Date.Format(time.DateOnly),
			Description: transaction.Description,
			CreatedAt:   transaction.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   transaction.UpdatedAt.Format(time.RFC3339),
		})
	}

	resp := dto.GetTransactionsResponse{
		Transactions: transactionResponses,
	}

	responde.WithJSON(w, r, resp, http.StatusOK)
}

// @Summary Create a transaction
// @Description Create a transaction for user
// @Tags transactions
// @Accept json
// @Produce json
// @Param request body dto.CreateTransactionRequest true "Create transaction request"
// @Param Authorization header string true "Bearer {token}" default "Bearer "
// @Success 200 {object} dto.TransactionResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /transactions [post]
func (h *Handler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	userID, ok := contextutil.GetUserID(r.Context())
	if !ok {
		h.log.Warnf("failed to get user ID from context")
		responde.WithError(w, r, h.log, nil, httperror.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	var req dto.CreateTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.WithError(err).Warnf("failed to decode request body")
		responde.WithError(w, r, h.log, err, httperror.ErrInvalidRequest, http.StatusBadRequest)
		return
	}

	transaction, err := h.transactionService.CreateTransaction(r.Context(), userID, &req)
	if err != nil {
		h.log.Errorf("failed to create transaction")
		responde.WithError(w, r, h.log, err, httperror.ErrFailedToCreateTransaction, http.StatusInternalServerError)
		return
	}

	resp := dto.TransactionResponse{
		Amount:      transaction.Amount,
		Category:    transaction.Category,
		Type:        transaction.Type,
		Date:        transaction.Date.Format(time.DateOnly),
		Description: transaction.Description,
		CreatedAt:   transaction.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   transaction.UpdatedAt.Format(time.RFC3339),
	}

	responde.WithJSON(w, r, resp, http.StatusOK)
}
