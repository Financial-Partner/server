package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/Financial-Partner/server/swagger"
)

// @title Financial Partner API
// @version 1.0
// @description API for the Financial Partner application
// @BasePath /api
func main() {
	cfgFile := flag.String("c", "config.yaml", "config file")
	flag.Parse()

	srv, err := InitializeServer(*cfgFile)
	if err != nil {
		ProvideLogger().WithError(err).Fatalf("Failed to initialize server")
	}

	srv.logger.Infof("Server is starting on port %s", srv.cfg.Server.Port)
	if err := srv.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		srv.logger.WithError(err).Fatalf("Server failed to start")
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	srv.logger.Infof("Server is shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.httpServer.Shutdown(ctx); err != nil {
		srv.logger.Fatalf("Server forced to shutdown: %v", err)
	}

	srv.logger.Infof("Server exited properly")
}
