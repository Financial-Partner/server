package main

import (
	"context"
	"net/http"
	"net/url"
	"time"

	"github.com/Financial-Partner/server/internal/config"
	authDomain "github.com/Financial-Partner/server/internal/domain/auth"
	goalDomain "github.com/Financial-Partner/server/internal/domain/goal"
	userDomain "github.com/Financial-Partner/server/internal/domain/user"
	authInfra "github.com/Financial-Partner/server/internal/infrastructure/auth"
	cacheInfra "github.com/Financial-Partner/server/internal/infrastructure/cache"
	dbInfra "github.com/Financial-Partner/server/internal/infrastructure/database"
	loggerInfra "github.com/Financial-Partner/server/internal/infrastructure/logger"
	perMongo "github.com/Financial-Partner/server/internal/infrastructure/persistence/mongodb"
	perRedis "github.com/Financial-Partner/server/internal/infrastructure/persistence/redis"
	handler "github.com/Financial-Partner/server/internal/interfaces/http"
	"github.com/Financial-Partner/server/internal/interfaces/http/middleware"
	"github.com/gorilla/mux"
)

func ProvideConfig(cfgFile string) (*config.Config, error) {
	return config.LoadConfig(cfgFile)
}

func ProvideLogger() loggerInfra.Logger {
	return loggerInfra.GetLogger()
}

func ProvideDBClient(cfg *config.Config) (*dbInfra.Client, error) {
	return dbInfra.NewClient(cfg)
}

func ProvideCacheClient(cfg *config.Config) (*cacheInfra.Client, error) {
	return cacheInfra.NewClient(cfg)
}

func ProvideAuthClient(cfg *config.Config) (*authInfra.Client, error) {
	return authInfra.NewClient(context.Background(), cfg)
}

func ProvideUserRepository(db *dbInfra.Client) userDomain.Repository {
	return perMongo.NewUserRepository(db)
}

func ProvideUserStore(cache *cacheInfra.Client) *perRedis.UserStore {
	return perRedis.NewUserStore(cache)
}

func ProvideUserService(repo userDomain.Repository, store *perRedis.UserStore, log loggerInfra.Logger) *userDomain.Service {
	return userDomain.NewService(repo, store, log)
}

func ProvideJWTManager(cfg *config.Config) *authInfra.JWTManager {
	return authInfra.NewJWTManager(cfg.JWT.SecretKey, cfg.JWT.AccessExpiry, cfg.JWT.RefreshExpiry)
}

func ProvideTokenStore(cache *cacheInfra.Client) *perRedis.TokenStore {
	return perRedis.NewTokenStore(cache)
}

func ProvideAuthService(
	cfg *config.Config,
	authClient *authInfra.Client,
	jwtManager *authInfra.JWTManager,
	tokenStore *perRedis.TokenStore,
	userService *userDomain.Service,
) *authDomain.Service {
	return authDomain.NewService(cfg, authClient, jwtManager, tokenStore, userService)
}

func ProvideGoalService() *goalDomain.Service {
	return goalDomain.NewService()
}

func ProvideHandler(
	userService *userDomain.Service,
	authService *authDomain.Service,
	goalService *goalDomain.Service,
	log loggerInfra.Logger,
) *handler.Handler {
	return handler.NewHandler(userService, authService, goalService, log)
}

func ProvideAuthMiddleware(jwtManager *authInfra.JWTManager, log loggerInfra.Logger) *middleware.AuthMiddleware {
	return middleware.NewAuthMiddleware(jwtManager, log)
}

func ProvideRouter(
	h *handler.Handler,
	authMiddleware *middleware.AuthMiddleware,
	cfg *config.Config,
) *mux.Router {
	router := mux.NewRouter()
	apiBaseURL := url.URL{
		Scheme: "http",
		Host:   cfg.Server.Host + ":" + cfg.Server.Port,
	}

	SetupRoutes(router, h, authMiddleware, apiBaseURL)
	return router
}

func ProvideServer(router *mux.Router, cfg *config.Config, log loggerInfra.Logger) *Server {
	httpServer := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return NewServer(httpServer, cfg, log)
}
