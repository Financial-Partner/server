package main

import (
	"net/http"

	"github.com/Financial-Partner/server/internal/config"
	"github.com/Financial-Partner/server/internal/infrastructure/logger"
)

type Server struct {
	httpServer *http.Server
	cfg        *config.Config
	logger     logger.Logger
}

func NewServer(server *http.Server, cfg *config.Config, logger logger.Logger) *Server {
	return &Server{
		httpServer: server,
		cfg:        cfg,
		logger:     logger,
	}
}
