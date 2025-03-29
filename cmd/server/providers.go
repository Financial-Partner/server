package main

import (
	"context"
	"net/http"
	"net/url"
	"time"

	"github.com/Financial-Partner/server/internal/config"
	authInfra "github.com/Financial-Partner/server/internal/infrastructure/auth"
	cacheInfra "github.com/Financial-Partner/server/internal/infrastructure/cache"
	dbInfra "github.com/Financial-Partner/server/internal/infrastructure/database"
	loggerInfra "github.com/Financial-Partner/server/internal/infrastructure/logger"
	perMongo "github.com/Financial-Partner/server/internal/infrastructure/persistence/mongodb"
	perRedis "github.com/Financial-Partner/server/internal/infrastructure/persistence/redis"
	handler "github.com/Financial-Partner/server/internal/interfaces/http"
	"github.com/Financial-Partner/server/internal/interfaces/http/middleware"
	auth_usecase "github.com/Financial-Partner/server/internal/module/auth/usecase"
	gacha_usecase "github.com/Financial-Partner/server/internal/module/gacha/usecase"
	goal_usecase "github.com/Financial-Partner/server/internal/module/goal/usecase"
	investment_usecase "github.com/Financial-Partner/server/internal/module/investment/usecase"
	transaction_usecase "github.com/Financial-Partner/server/internal/module/transaction/usecase"
	user_repository "github.com/Financial-Partner/server/internal/module/user/repository"
	user_usecase "github.com/Financial-Partner/server/internal/module/user/usecase"
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

func ProvideUserRepository(db *dbInfra.Client) user_repository.Repository {
	return perMongo.NewUserRepository(db)
}

func ProvideUserStore(cache *cacheInfra.Client) *perRedis.UserStore {
	return perRedis.NewUserStore(cache)
}

func ProvideUserService(repo user_repository.Repository, store *perRedis.UserStore, log loggerInfra.Logger) *user_usecase.Service {
	return user_usecase.NewService(repo, store, log)
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
	userService *user_usecase.Service,
) *auth_usecase.Service {
	return auth_usecase.NewService(cfg, authClient, jwtManager, tokenStore, userService)
}

func ProvideGoalService() *goal_usecase.Service {
	return goal_usecase.NewService()
}

func ProvideInvestmentService() *investment_usecase.Service {
	return investment_usecase.NewService()
}

func ProvideTransactionService() *transaction_usecase.Service {
	return transaction_usecase.NewService()
}

func ProvideGachaService() *gacha_usecase.Service {
	return gacha_usecase.NewService()
}

func ProvideHandler(
	userService *user_usecase.Service,
	authService *auth_usecase.Service,
	goalService *goal_usecase.Service,
	investmentService *investment_usecase.Service,
	transactionService *transaction_usecase.Service,
	gachaService *gacha_usecase.Service,
	log loggerInfra.Logger,
) *handler.Handler {
	return handler.NewHandler(userService, authService, goalService, investmentService, transactionService, gachaService, log)
}

func ProvideAuthMiddleware(jwtManager *authInfra.JWTManager, cfg *config.Config, log loggerInfra.Logger) *middleware.AuthMiddleware {
	if cfg.Firebase.BypassEnabled {
		return middleware.NewAuthMiddleware(authInfra.NewDummyJWTValidator(cfg), log)
	}
	return middleware.NewAuthMiddleware(jwtManager, log)
}

func ProvideLoggerMiddleware(log loggerInfra.Logger) *middleware.LoggerMiddleware {
	return middleware.NewLoggerMiddleware(log)
}

func ProvideRouter(
	h *handler.Handler,
	authMiddleware *middleware.AuthMiddleware,
	loggerMiddleware *middleware.LoggerMiddleware,
	cfg *config.Config,
) *mux.Router {
	router := mux.NewRouter()
	apiBaseURL := url.URL{
		Scheme: "http",
		Host:   cfg.Server.Host + ":" + cfg.Server.Port,
	}

	SetupRoutes(router, h, authMiddleware, loggerMiddleware, apiBaseURL)
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
