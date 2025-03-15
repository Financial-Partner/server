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
		ProvideTokenStore,
		ProvideAuthService,
		ProvideGoalService,
		ProvideInvestmentService,
		ProvideHandler,
		ProvideAuthMiddleware,
		ProvideRouter,
		ProvideServer,
	)
	return nil, nil
}
