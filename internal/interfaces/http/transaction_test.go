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
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/mock/gomock"
)

func TestCreateTransaction(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("Invalid request format", func(t *testing.T) {
		h, _ := newTestHandler(t)

		invalidBody := bytes.NewBufferString(`{invalid json`)
		userID := primitive.NewObjectID().Hex()
		ctx := context.WithValue(context.Background(), contextutil.UserEmailKey, "test@example.com")
		ctx = context.WithValue(ctx, contextutil.UserIDKey, userID)
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/transactions", invalidBody)
		r = r.WithContext(ctx)

		h.CreateTransaction(w, r)

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

		req := dto.CreateTransactionRequest{
			Amount:      1000,
			Category:    "Food",
			Type:        "expense",
			Date:        "2023-01-01",
			Description: "Lunch",
		}
		body, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/transactions", bytes.NewBuffer(body))

		h.CreateTransaction(w, r)

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

		userID := primitive.NewObjectID().Hex()
		userEmail := "test@example.com"

		mockServices.TransactionService.EXPECT().
			CreateTransaction(gomock.Any(), userID, gomock.Any()).
			Return(nil, errors.New("service error"))

		req := dto.CreateTransactionRequest{
			Amount:      1000,
			Category:    "Food",
			Type:        "expense",
			Date:        "2023-01-01",
			Description: "Lunch",
		}
		body, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/transactions", bytes.NewBuffer(body))
		ctx := newContext(userID, userEmail)
		r = r.WithContext(ctx)

		h.CreateTransaction(w, r)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var errorResp dto.ErrorResponse
		err := json.NewDecoder(w.Body).Decode(&errorResp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, errorResp.Code)
		assert.Equal(t, httperror.ErrFailedToCreateTransaction, errorResp.Message)
	})

	t.Run("Success", func(t *testing.T) {
		h, mockServices := newTestHandler(t)

		userID := primitive.NewObjectID()
		userEmail := "test@example.com"

		now := time.Now()
		objectID := primitive.NewObjectID()
		transaction := &entities.Transaction{
			ID:          objectID,
			Amount:      1000,
			Description: "Lunch",
			Date:        time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC),
			Category:    "Food",
			Type:        "expense",
			UserID:      primitive.NewObjectID(),
			CreatedAt:   now,
			UpdatedAt:   now,
		}

		mockServices.TransactionService.EXPECT().
			CreateTransaction(gomock.Any(), userID.Hex(), gomock.Any()).
			Return(transaction, nil)

		req := dto.CreateTransactionRequest{
			Amount:      1000,
			Category:    "Food",
			Type:        "expense",
			Date:        "2023-01-01",
			Description: "Lunch",
		}
		body, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/transactions", bytes.NewBuffer(body))
		ctx := newContext(userID.Hex(), userEmail)
		r = r.WithContext(ctx)

		h.CreateTransaction(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var response dto.TransactionResponse
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, transaction.Amount, response.Amount)
		assert.Equal(t, transaction.Description, response.Description)
		assert.Equal(t, transaction.Date.Format(time.DateOnly), response.Date)
		assert.Equal(t, transaction.Category, response.Category)
		assert.Equal(t, transaction.Type, response.Type)
		assert.Equal(t, transaction.CreatedAt.Format(time.RFC3339), response.CreatedAt)
		assert.Equal(t, transaction.UpdatedAt.Format(time.RFC3339), response.UpdatedAt)
	})
}

func TestGetTransactions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("Unauthorized request", func(t *testing.T) {
		h, _ := newTestHandler(t)

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/transactions", nil)

		h.GetTransactions(w, r)

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

		userID := primitive.NewObjectID()
		userEmail := "test@example.com"

		mockServices.TransactionService.EXPECT().
			GetTransactions(gomock.Any(), userID.Hex()).
			Return(nil, errors.New("service error"))

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/transactions", nil)
		ctx := newContext(userID.Hex(), userEmail)
		r = r.WithContext(ctx)

		h.GetTransactions(w, r)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var errorResp dto.ErrorResponse
		err := json.NewDecoder(w.Body).Decode(&errorResp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, errorResp.Code)
		assert.Equal(t, httperror.ErrFailedToGetTransactions, errorResp.Message)
	})

	t.Run("Success", func(t *testing.T) {
		h, mockServices := newTestHandler(t)

		userID := primitive.NewObjectID()
		userEmail := "test@example.com"

		now := time.Now()
		objectID := primitive.NewObjectID()
		transactions := []entities.Transaction{
			{
				ID:          objectID,
				Amount:      1000,
				Description: "Lunch",
				Date:        time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC),
				Category:    "Food",
				Type:        "expense",
				UserID:      primitive.NewObjectID(),
				CreatedAt:   now,
				UpdatedAt:   now,
			},
		}

		mockServices.TransactionService.EXPECT().
			GetTransactions(gomock.Any(), userID.Hex()).
			Return(transactions, nil)

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/transactions", nil)
		ctx := newContext(userID.Hex(), userEmail)
		r = r.WithContext(ctx)

		h.GetTransactions(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var response dto.GetTransactionsResponse
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, len(transactions), len(response.Transactions))
		assert.Equal(t, transactions[0].Amount, response.Transactions[0].Amount)
		assert.Equal(t, transactions[0].Description, response.Transactions[0].Description)
		assert.Equal(t, transactions[0].Date.Format(time.DateOnly), response.Transactions[0].Date)
		assert.Equal(t, transactions[0].Category, response.Transactions[0].Category)
		assert.Equal(t, transactions[0].Type, response.Transactions[0].Type)
		assert.Equal(t, transactions[0].CreatedAt.Format(time.RFC3339), response.Transactions[0].CreatedAt)
		assert.Equal(t, transactions[0].UpdatedAt.Format(time.RFC3339), response.Transactions[0].UpdatedAt)
	})
}
