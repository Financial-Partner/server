package handler

import (
	"github.com/Financial-Partner/server/internal/infrastructure/logger"
)

type Handler struct {
	userService UserService
	authService AuthService
	log         logger.Logger
}

func NewHandler(us UserService, as AuthService, log logger.Logger) *Handler {
	return &Handler{
		userService: us,
		authService: as,
		log:         log,
	}
}
