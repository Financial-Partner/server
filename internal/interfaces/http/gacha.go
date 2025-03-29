package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Financial-Partner/server/internal/contextutil"
	"github.com/Financial-Partner/server/internal/entities"
	"github.com/Financial-Partner/server/internal/interfaces/http/dto"
	httperror "github.com/Financial-Partner/server/internal/interfaces/http/error"
	responde "github.com/Financial-Partner/server/internal/interfaces/http/respond"
)

//go:generate mockgen -source=gacha.go -destination=gacha_mock.go -package=handler

type GachaService interface {
	DrawGacha(ctx context.Context, userID string, req *dto.DrawGachaRequest) (*entities.Gacha, error)
	PreviewGachas(ctx context.Context, userID string) ([]entities.Gacha, error)
}

// @Summary Decrease user's gacha amount and return gacha result
// @Description Decrease user's gacha amount and return gacha result
// @Tags gacha
// @Accept json
// @Produce json
// @Param request body dto.DrawGachaRequest true "Draw gacha request"
// @Param Authorization header string true "Bearer {token}" default "Bearer "
// @Success 200 {object} dto.GachaResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /gacha/draw [post]
func (h *Handler) DrawGacha(w http.ResponseWriter, r *http.Request) {
	userID, ok := contextutil.GetUserID(r.Context())
	if !ok {
		h.log.Warnf("failed to get user ID from context")
		responde.WithError(w, r, h.log, nil, httperror.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	var req dto.DrawGachaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.WithError(err).Warnf("failed to decode request body")
		responde.WithError(w, r, h.log, err, httperror.ErrInvalidRequest, http.StatusBadRequest)
		return
	}

	gacha, err := h.gachaService.DrawGacha(r.Context(), userID, &req)
	if err != nil {
		h.log.WithError(err).Warnf("failed to draw a gacha")
		responde.WithError(w, r, h.log, err, httperror.ErrFailedToDrawGacha, http.StatusInternalServerError)
		return
	}

	resp := dto.GachaResponse{
		ID:     gacha.ID.Hex(),
		ImgSrc: gacha.ImgSrc,
	}

	responde.WithJSON(w, r, resp, http.StatusOK)
}

// @Summary Get 9 gacha images for preview
// @Description Get 9 gacha images for preview
// @Tags gacha
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer {token}" default "Bearer "
// @Success 200 {object} dto.PreviewGachasResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /gacha/preview [get]
func (h *Handler) PreviewGachas(w http.ResponseWriter, r *http.Request) {
	userID, ok := contextutil.GetUserID(r.Context())
	if !ok {
		h.log.Warnf("failed to get user ID from context")
		responde.WithError(w, r, h.log, nil, httperror.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	gachas, err := h.gachaService.PreviewGachas(r.Context(), userID)
	if err != nil {
		h.log.WithError(err).Warnf("failed to get preview gacha images")
		responde.WithError(w, r, h.log, err, httperror.ErrFailedToPreviewGachas, http.StatusInternalServerError)
		return
	}

	resp := dto.PreviewGachasResponse{
		Gachas: make([]dto.GachaResponse, 0, 9),
	}

	for _, gacha := range gachas {
		resp.Gachas = append(resp.Gachas, dto.GachaResponse{
			ID:     gacha.ID.Hex(),
			ImgSrc: gacha.ImgSrc,
		})
	}

	responde.WithJSON(w, r, resp, http.StatusOK)
}
