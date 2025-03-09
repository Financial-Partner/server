//go:build wireinject
// +build wireinject

package main

import (
	"net/http"

	"github.com/google/wire"
)

func InitializeServer(cfgFile string) (*http.Server, error) {
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
		ProvideHandler,
		ProvideAuthMiddleware,
		ProvideRouter,
		ProvideHTTPServer,
	)
	return nil, nil
}
