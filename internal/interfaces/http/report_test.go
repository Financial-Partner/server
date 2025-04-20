package handler_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Financial-Partner/server/internal/entities"
	"github.com/Financial-Partner/server/internal/interfaces/http/dto"
	httperror "github.com/Financial-Partner/server/internal/interfaces/http/error"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestGetReport(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("Unauthorized request", func(t *testing.T) {
		h, _ := newTestHandler(t)

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/reports/finance", nil)

		h.GetReport(w, r)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var errorResp dto.ErrorResponse
		err := json.NewDecoder(w.Body).Decode(&errorResp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, errorResp.Code)
		assert.Equal(t, httperror.ErrUnauthorized, errorResp.Message)
	})

	t.Run("Invalid parameter format", func(t *testing.T) {
		h, _ := newTestHandler(t)

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/reports/finance?type=summary&start=invalid-date&end=2025-01-31", nil)

		h.GetReport(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		// Decode the error response
		var errorResp dto.ErrorResponse
		err := json.NewDecoder(w.Body).Decode(&errorResp)
		assert.NoError(t, err)

		// Assert the error response details
		assert.Equal(t, http.StatusBadRequest, errorResp.Code)
		assert.Equal(t, httperror.ErrInvalidParameter, errorResp.Message)
	})

	t.Run("Service error", func(t *testing.T) {
		h, mockService := newTestHandler(t)

		userID := "testUserID"
		userEmail := "test@example.com"

		mockService.ReportService.EXPECT().
			GetReport(gomock.Any(), userID, gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil, errors.New("service error"))

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/reports/finance?type=summary&start=2025-01-01&end=2025-01-31", nil)
		ctx := newContext(userID, userEmail)
		r = r.WithContext(ctx)

		h.GetReport(w, r)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var errorResp dto.ErrorResponse
		err := json.NewDecoder(w.Body).Decode(&errorResp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, errorResp.Code)
		assert.Equal(t, httperror.ErrFailedToGetReport, errorResp.Message)
	})

	t.Run("Success", func(t *testing.T) {
		h, mockService := newTestHandler(t)

		userID := "testUserID"
		userEmail := "test@example.com"

		report := &entities.Report{
			Revenue:     1000,
			Expenses:    500,
			NetProfit:   500,
			Categorys:   []string{"Food", "Transport"},
			Amounts:     []int64{200, 300},
			Percentages: []float64{0.4, 0.6},
		}

		mockService.ReportService.EXPECT().
			GetReport(gomock.Any(), userID, gomock.Any(), gomock.Any(), gomock.Any()).
			Return(report, nil)

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/reports/finance?type=summary&start=2025-01-01&end=2025-01-31", nil)
		ctx := newContext(userID, userEmail)
		r = r.WithContext(ctx)

		h.GetReport(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var response dto.ReportResponse
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, report.Revenue, response.Revenue)
		assert.Equal(t, report.Expenses, response.Expenses)
		assert.Equal(t, report.NetProfit, response.NetProfit)
		assert.Equal(t, report.Categorys, response.Categorys)
		assert.Equal(t, report.Amounts, response.Amounts)
		assert.Equal(t, report.Percentages, response.Percentages)
	})
}

func TestGetReportSummary(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("Unauthorized request", func(t *testing.T) {
		h, _ := newTestHandler(t)

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/reports/analysis", nil)

		h.GetReportSummary(w, r)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var errorResp dto.ErrorResponse
		err := json.NewDecoder(w.Body).Decode(&errorResp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, errorResp.Code)
		assert.Equal(t, httperror.ErrUnauthorized, errorResp.Message)
	})

	t.Run("Service error", func(t *testing.T) {
		h, mockService := newTestHandler(t)

		userID := "testUserID"
		userEmail := "test@example.com"

		mockService.ReportService.EXPECT().
			GetReportSummary(gomock.Any(), userID).
			Return(nil, errors.New("service error"))

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/reports/analysis", nil)
		ctx := newContext(userID, userEmail)
		r = r.WithContext(ctx)

		h.GetReportSummary(w, r)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var errorResp dto.ErrorResponse
		err := json.NewDecoder(w.Body).Decode(&errorResp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, errorResp.Code)
		assert.Equal(t, httperror.ErrFailedToGetReportSummary, errorResp.Message)
	})

	t.Run("Success", func(t *testing.T) {
		h, mockService := newTestHandler(t)

		userID := "testUserID"
		userEmail := "test@example.com"

		reportSummary := &entities.ReportSummary{
			Summary: "This is a summary",
		}

		mockService.ReportService.EXPECT().
			GetReportSummary(gomock.Any(), userID).
			Return(reportSummary, nil)

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/reports/analysis", nil)
		ctx := newContext(userID, userEmail)
		r = r.WithContext(ctx)

		h.GetReportSummary(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var response dto.ReportSummaryResponse
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, reportSummary.Summary, response.Summary)
	})
}
