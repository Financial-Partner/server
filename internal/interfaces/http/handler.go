package handler

import (
	"github.com/Financial-Partner/server/internal/infrastructure/logger"
)

type Handler struct {
	userService        UserService
	authService        AuthService
	goalService        GoalService
	investmentService  InvestmentService
	transactionService TransactionService
	gachaService       GachaService
	log                logger.Logger
}

func NewHandler(us UserService, as AuthService, gs GoalService, is InvestmentService, ts TransactionService, gcs GachaService, log logger.Logger) *Handler {
	return &Handler{
		userService:        us,
		authService:        as,
		goalService:        gs,
		investmentService:  is,
		transactionService: ts,
		gachaService:       gcs,
		log:                log,
	}
}
