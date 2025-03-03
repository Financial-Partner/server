package main

import (
	"context"
	"flag"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"

	"github.com/Financial-Partner/server/internal/config"
	userDomain "github.com/Financial-Partner/server/internal/domain/user"
	authInfra "github.com/Financial-Partner/server/internal/infrastructure/auth"
	cacheInfra "github.com/Financial-Partner/server/internal/infrastructure/cache"
	dbInfra "github.com/Financial-Partner/server/internal/infrastructure/database"
	loggerInfra "github.com/Financial-Partner/server/internal/infrastructure/logger"
	mongodb "github.com/Financial-Partner/server/internal/infrastructure/persistence/mongodb"
	redis "github.com/Financial-Partner/server/internal/infrastructure/persistence/redis"
	handler "github.com/Financial-Partner/server/internal/interfaces/http"
	"github.com/Financial-Partner/server/internal/interfaces/http/middleware"
	_ "github.com/Financial-Partner/server/swagger"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Financial Partner API
// @version 1.0
// @description API for the Financial Partner application
// @BasePath /api
func main() {
	log := loggerInfra.GetLogger()

	cfgFile := flag.String("c", "config.yaml", "config file")
	flag.Parse()

	cfg, err := config.LoadConfig(*cfgFile)
	if err != nil {
		log.WithError(err).Fatalf("Failed to load config")
	}

	db, err := dbInfra.NewClient(cfg)
	if err != nil {
		log.WithError(err).Fatalf("Failed to connect to MongoDB")
	}
	defer db.Close(context.Background())

	cacheClient, err := cacheInfra.NewClient(cfg)
	if err != nil {
		log.WithError(err).Fatalf("Failed to connect to Redis")
	}

	authClient, err := authInfra.NewClient(context.Background(), cfg)
	if err != nil {
		log.WithError(err).Fatalf("Failed to initialize Firebase Auth")
	}

	userRepo := mongodb.NewUserRepository(db)
	userStore := redis.NewUserStore(cacheClient)

	userService := userDomain.NewService(userRepo, userStore, log)
	authMiddleware := middleware.NewAuthMiddleware(authClient, log)
	handlers := handler.NewHandler(userService, log)

	router := mux.NewRouter()

	apiBaseURL := url.URL{
		Scheme: "http",
		Host:   cfg.Server.Host + ":" + cfg.Server.Port,
	}

	router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL(apiBaseURL.String()+"/swagger/doc.json"),
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	))

	api := router.PathPrefix("/api").Subrouter()
	api.Use(authMiddleware.Authenticate)

	api.HandleFunc("/users/me", handlers.GetUser).Methods(http.MethodGet)

	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Infof("Server is starting on port %s", cfg.Server.Port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.WithError(err).Fatalf("Server failed to start")
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Infof("Server is shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Infof("Server exited properly")
}
