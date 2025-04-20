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

func TestGetOpportunities(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("Unauthorized request", func(t *testing.T) {
		h, _ := newTestHandler(t)

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/investments", nil)

		h.GetOpportunities(w, r)

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

		mockServices.InvestmentService.EXPECT().
			GetOpportunities(gomock.Any(), userID).
			Return(nil, errors.New("service error"))

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/investments", nil)
		ctx := newContext(userID, userEmail)
		r = r.WithContext(ctx)

		h.GetOpportunities(w, r)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var errorResp dto.ErrorResponse
		err := json.NewDecoder(w.Body).Decode(&errorResp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, errorResp.Code)
		assert.Equal(t, httperror.ErrFailedToGetOpportunities, errorResp.Message)
	})

	t.Run("Success", func(t *testing.T) {
		h, mockServices := newTestHandler(t)

		now := time.Now()
		userID := primitive.NewObjectID().Hex()
		userEmail := "test@example.com"
		opportunities := []entities.Opportunity{
			{
				ID:          primitive.NewObjectID(),
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
			GetOpportunities(gomock.Any(), userID).
			Return(opportunities, nil)

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/investments", nil)
		ctx := newContext(userID, userEmail)
		r = r.WithContext(ctx)

		h.GetOpportunities(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var response dto.GetOpportunitiesResponse
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		// Compare the response with the expected data
		assert.Len(t, response.Opportunities, 1)
		assert.Equal(t, opportunities[0].Title, response.Opportunities[0].Title)
		assert.Equal(t, opportunities[0].Description, response.Opportunities[0].Description)
		assert.Equal(t, opportunities[0].Tags, response.Opportunities[0].Tags)
		assert.Equal(t, opportunities[0].IsIncrease, response.Opportunities[0].IsIncrease)
		assert.Equal(t, opportunities[0].Variation, response.Opportunities[0].Variation)
		assert.Equal(t, opportunities[0].Duration, response.Opportunities[0].Duration)
		assert.Equal(t, opportunities[0].MinAmount, response.Opportunities[0].MinAmount)
		assert.Equal(t, opportunities[0].CreatedAt.Format(time.RFC3339), response.Opportunities[0].CreatedAt)
		assert.Equal(t, opportunities[0].UpdatedAt.Format(time.RFC3339), response.Opportunities[0].UpdatedAt)
	})
}

func TestCreateUserInvestment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("Invalid request format", func(t *testing.T) {
		h, _ := newTestHandler(t)

		invalidBody := bytes.NewBufferString(`{invalid json`)
		userID := primitive.NewObjectID().Hex()
		ctx := context.WithValue(context.Background(), contextutil.UserEmailKey, "test@example.com")
		ctx = context.WithValue(ctx, contextutil.UserIDKey, userID)
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/me/investments", invalidBody)
		r = r.WithContext(ctx)

		h.CreateUserInvestment(w, r)

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

		req := dto.CreateUserInvestmentRequest{
			OpportunityID: primitive.NewObjectID().Hex(),
			Amount:        1000,
		}
		body, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/me/investments", bytes.NewBuffer(body))

		h.CreateUserInvestment(w, r)

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

		mockServices.InvestmentService.EXPECT().
			CreateUserInvestment(gomock.Any(), userID, gomock.Any()).
			Return(nil, errors.New("service error"))

		req := dto.CreateUserInvestmentRequest{
			OpportunityID: primitive.NewObjectID().Hex(),
			Amount:        1000,
		}
		body, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/me/investments", bytes.NewBuffer(body))
		ctx := newContext(userID, userEmail)
		r = r.WithContext(ctx)

		h.CreateUserInvestment(w, r)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var errorResp dto.ErrorResponse
		err := json.NewDecoder(w.Body).Decode(&errorResp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, errorResp.Code)
		assert.Equal(t, httperror.ErrFailedToCreateUserInvestment, errorResp.Message)
	})

	t.Run("Success", func(t *testing.T) {
		h, mockServices := newTestHandler(t)

		userID := primitive.NewObjectID().Hex()
		userEmail := "test@example.com"

		now := time.Now()
		investment := &entities.Investment{
			ID:            primitive.NewObjectID(),
			UserID:        userID,
			OpportunityID: primitive.NewObjectID().Hex(),
			Amount:        1000,
			CreatedAt:     now,
			UpdatedAt:     now,
		}

		mockServices.InvestmentService.EXPECT().
			CreateUserInvestment(gomock.Any(), userID, gomock.Any()).
			Return(investment, nil)

		req := dto.CreateUserInvestmentRequest{
			OpportunityID: primitive.NewObjectID().Hex(),
			Amount:        1000,
		}

		body, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/me/investments", bytes.NewBuffer(body))
		ctx := newContext(userID, userEmail)
		r = r.WithContext(ctx)

		h.CreateUserInvestment(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var response dto.InvestmentResponse
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, investment.OpportunityID, investment.OpportunityID)
		assert.Equal(t, investment.Amount, investment.Amount)
		assert.Equal(t, investment.CreatedAt.Format(time.RFC3339), response.CreatedAt)
		assert.Equal(t, investment.UpdatedAt.Format(time.RFC3339), response.UpdatedAt)
	})
}

func TestGetUserInvestments(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("Unauthorized request", func(t *testing.T) {
		h, _ := newTestHandler(t)

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/me/investments", nil)

		h.GetUserInvestments(w, r)

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

		mockServices.InvestmentService.EXPECT().
			GetUserInvestments(gomock.Any(), userID).
			Return(nil, errors.New("service error"))

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/me/investments", nil)
		ctx := newContext(userID, userEmail)
		r = r.WithContext(ctx)

		h.GetUserInvestments(w, r)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var errorResp dto.ErrorResponse
		err := json.NewDecoder(w.Body).Decode(&errorResp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, errorResp.Code)
		assert.Equal(t, httperror.ErrFailedToGetUserInvestments, errorResp.Message)
	})

	t.Run("Success", func(t *testing.T) {
		h, mockServices := newTestHandler(t)

		now := time.Now()
		userID := primitive.NewObjectID().Hex()
		userEmail := "test@example.com"
		investments := []entities.Investment{
			{
				ID:            primitive.NewObjectID(),
				OpportunityID: primitive.NewObjectID().Hex(),
				Amount:        1000,
				CreatedAt:     now.AddDate(0, -1, 0),
				UpdatedAt:     now,
			},
		}

		mockServices.InvestmentService.EXPECT().
			GetUserInvestments(gomock.Any(), userID).
			Return(investments, nil)

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/me/investments", nil)
		ctx := newContext(userID, userEmail)
		r = r.WithContext(ctx)

		h.GetUserInvestments(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var response dto.GetUserInvestmentsResponse
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		// Compare the response with the expected data
		assert.Len(t, response.Investments, 1)
		assert.Equal(t, investments[0].OpportunityID, response.Investments[0].OpportunityID)
		assert.Equal(t, investments[0].Amount, response.Investments[0].Amount)
		assert.Equal(t, investments[0].CreatedAt.Format(time.RFC3339), response.Investments[0].CreatedAt)
		assert.Equal(t, investments[0].UpdatedAt.Format(time.RFC3339), response.Investments[0].UpdatedAt)
	})
}
