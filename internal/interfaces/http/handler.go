package handler

import (
	"github.com/Financial-Partner/server/internal/infrastructure/logger"
)

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
