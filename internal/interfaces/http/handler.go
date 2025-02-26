package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Financial-Partner/server/internal/contextutil"
	"github.com/Financial-Partner/server/internal/domain/user"
	"github.com/Financial-Partner/server/internal/infrastructure/logger"
)

//go:generate mockgen -source=handler.go -destination=mocks/handler_mock.go -package=mocks

type UserService interface {
	GetUser(ctx context.Context, email string) (*user.UserEntity, error)
	GetOrCreateUser(ctx context.Context, email, name string) (*user.UserEntity, error)
}

type Handler struct {
	userService UserService
	log         logger.Logger
}

func NewHandler(us UserService, log logger.Logger) *Handler {
	return &Handler{
		userService: us,
		log:         log,
	}
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	email := r.Context().Value(contextutil.UserEmailKey)
	if email == nil {
		h.log.Errorf("User email not found in context")
		http.Error(w, "User email not found in context", http.StatusInternalServerError)
		return
	}

	logger := h.log.WithField("email", email)

	userEntity, err := h.userService.GetOrCreateUser(r.Context(), email.(string), "")
	if err != nil {
		logger.WithError(err).Errorf("Failed to get or create user")
		http.Error(w, "Failed to retrieve user information", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(userEntity)
}
