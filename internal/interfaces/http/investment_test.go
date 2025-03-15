package handler_test

import (
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

func TestGetInvestments(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("Unauthorized request", func(t *testing.T) {
		h, _ := newTestHandler(t)

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/investments", nil)

		h.GetInvestments(w, r)

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

		mockServices.InvestmentService.EXPECT().
			GetInvestments(gomock.Any(), "test@example.com").
			Return(nil, errors.New("service error"))

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/investments", nil)
		ctx := context.WithValue(r.Context(), contextutil.UserEmailKey, "test@example.com")
		r = r.WithContext(ctx)

		h.GetInvestments(w, r)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var errorResp dto.ErrorResponse
		err := json.NewDecoder(w.Body).Decode(&errorResp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, errorResp.Code)
		assert.Equal(t, httperror.ErrFailedToGetInvestments, errorResp.Message)
	})

	t.Run("Success", func(t *testing.T) {
		h, mockServices := newTestHandler(t)

		now := time.Now()
		investments := []entities.Investment{
			{
				ID:          "investment_123456",
				UserID:      "test@example.com",
				Title:       "Test investments",
				Description: "Test description",
				Tags:        []string{"test", "investments"},
				IsIncrease:  true,
				Variation:   30,
				Duration:    "a month",
				MinAmount:   1000,
				CreatedAt:   now.AddDate(0, -1, 0),
				UpdatedAt:   now,
			},
		}

		mockServices.InvestmentService.EXPECT().
			GetInvestments(gomock.Any(), "test@example.com").
			Return(investments, nil)

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/investments", nil)
		ctx := context.WithValue(r.Context(), contextutil.UserEmailKey, "test@example.com")
		r = r.WithContext(ctx)

		h.GetInvestments(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var response dto.GetInvestmentsResponse
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		// Compare the response with the expected data
		assert.Len(t, response.Investments, 1)
		assert.Equal(t, investments[0].Title, response.Investments[0].Title)
		assert.Equal(t, investments[0].Description, response.Investments[0].Description)
		assert.Equal(t, investments[0].Tags, response.Investments[0].Tags)
		assert.Equal(t, investments[0].IsIncrease, response.Investments[0].IsIncrease)
		assert.Equal(t, investments[0].Variation, response.Investments[0].Variation)
		assert.Equal(t, investments[0].Duration, response.Investments[0].Duration)
		assert.Equal(t, investments[0].MinAmount, response.Investments[0].MinAmount)
		assert.Equal(t, investments[0].CreatedAt.Format(time.RFC3339), response.Investments[0].CreatedAt)
		assert.Equal(t, investments[0].UpdatedAt.Format(time.RFC3339), response.Investments[0].UpdatedAt)
	})
}
