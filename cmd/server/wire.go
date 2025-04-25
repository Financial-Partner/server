//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
)

func InitializeServer(cfgFile string) (*Server, error) {
	wire.Build(
		ProvideConfig,
		ProvideLogger,
		ProvideDBClient,
		ProvideCacheClient,
		ProvideAuthClient,
		ProvideUserRepository,
		ProvideUserStore,
		ProvideUserService,
		ProvideJWTManager,
		ProvideLoggerMiddleware,
		ProvideTokenStore,
		ProvideAuthService,
		ProvideGoalService,
		ProvideInvestmentService,
		ProvideTransactionRepository,
		ProvideTransactionStore,
		ProvideTransactionService,
		ProvideGachaService,
		ProvideReportService,
		ProvideHandler,
		ProvideAuthMiddleware,
		ProvideRouter,
		ProvideServer,
	)
	return nil, nil
}
