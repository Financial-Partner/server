package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Financial-Partner/server/internal/contextutil"
	"github.com/Financial-Partner/server/internal/entities"
	"github.com/Financial-Partner/server/internal/interfaces/http/dto"
	httperror "github.com/Financial-Partner/server/internal/interfaces/http/error"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/mock/gomock"
)

func TestDrawGacha(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("Invalid request format", func(t *testing.T) {
		h, _ := newTestHandler(t)

		invalidBody := bytes.NewBufferString(`{invalid json`)
		userID := primitive.NewObjectID().Hex()
		ctx := context.WithValue(context.Background(), contextutil.UserEmailKey, "test@example.com")
		ctx = context.WithValue(ctx, contextutil.UserIDKey, userID)
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/gacha/draw", invalidBody)
		r = r.WithContext(ctx)

		h.DrawGacha(w, r)

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

		req := dto.DrawGachaRequest{
			Amount: 100,
		}
		body, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/gacha/draw", bytes.NewBuffer(body))

		h.DrawGacha(w, r)

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

		mockServices.GachaService.EXPECT().
			DrawGacha(gomock.Any(), userID, gomock.Any()).
			Return(nil, errors.New("service error"))

		req := dto.DrawGachaRequest{
			Amount: 100,
		}
		body, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/gacha/draw", bytes.NewBuffer(body))
		ctx := newContext(userID, userEmail)
		r = r.WithContext(ctx)

		h.DrawGacha(w, r)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var errorResp dto.ErrorResponse
		err := json.NewDecoder(w.Body).Decode(&errorResp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, errorResp.Code)
		assert.Equal(t, httperror.ErrFailedToDrawGacha, errorResp.Message)
	})

	t.Run("Success", func(t *testing.T) {
		h, mockServices := newTestHandler(t)

		userID := primitive.NewObjectID().Hex()
		userEmail := "test@example.com"

		objectID := primitive.NewObjectID()
		gacha := &entities.Gacha{
			ID:     objectID,
			ImgSrc: "https://example.com/image.png",
		}

		mockServices.GachaService.EXPECT().
			DrawGacha(gomock.Any(), userID, gomock.Any()).
			Return(gacha, nil)

		req := dto.DrawGachaRequest{
			Amount: 100,
		}
		body, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/gacha/draw", bytes.NewBuffer(body))
		ctx := newContext(userID, userEmail)
		r = r.WithContext(ctx)

		h.DrawGacha(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var response dto.GachaResponse
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, gacha.ID.Hex(), response.ID)
		assert.Equal(t, gacha.ImgSrc, response.ImgSrc)
	})
}

func TestPreviewGacha(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("Unauthorized request", func(t *testing.T) {
		h, _ := newTestHandler(t)

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/gacha/preview", nil)

		h.PreviewGachas(w, r)

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

		mockServices.GachaService.EXPECT().
			PreviewGachas(gomock.Any(), userID).
			Return(nil, errors.New("service error"))

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/gacha/preview", nil)
		ctx := newContext(userID, userEmail)
		r = r.WithContext(ctx)

		h.PreviewGachas(w, r)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var errorResp dto.ErrorResponse
		err := json.NewDecoder(w.Body).Decode(&errorResp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, errorResp.Code)
		assert.Equal(t, httperror.ErrFailedToPreviewGachas, errorResp.Message)
	})

	t.Run("Success", func(t *testing.T) {
		h, mockServices := newTestHandler(t)

		userID := primitive.NewObjectID().Hex()
		userEmail := "test@example.com"

		gachas := make([]entities.Gacha, 9) // Create a slice for 9 Gacha objects
		for i := 0; i < 9; i++ {
			objectID := primitive.NewObjectID()
			gachas[i] = entities.Gacha{
				ID:     objectID,
				ImgSrc: fmt.Sprintf("https://example.com/image%d.png", i+1),
			}
		}

		mockServices.GachaService.EXPECT().
			PreviewGachas(gomock.Any(), userID).
			Return(gachas, nil)

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/gacha/preview", nil)
		ctx := newContext(userID, userEmail)
		r = r.WithContext(ctx)

		h.PreviewGachas(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var response dto.PreviewGachasResponse
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Len(t, response.Gachas, 9) // Check if the length of the Gachas slice is 9
		for i, gacha := range gachas {
			assert.Equal(t, gacha.ID.Hex(), response.Gachas[i].ID)
			assert.Equal(t, gacha.ImgSrc, response.Gachas[i].ImgSrc)
		}
	})
}
