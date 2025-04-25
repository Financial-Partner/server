package main

import (
	"net/http"
	"net/url"

	"github.com/gorilla/mux"

	handler "github.com/Financial-Partner/server/internal/interfaces/http"
	"github.com/Financial-Partner/server/internal/interfaces/http/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

func SetupRoutes(
	router *mux.Router,
	handlers *handler.Handler,
	authMiddleware *middleware.AuthMiddleware,
	loggerMiddleware *middleware.LoggerMiddleware,
	apiBaseURL url.URL,
) {
	router.Use(loggerMiddleware.LogRequest)

	router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL(apiBaseURL.String()+"/swagger/doc.json"),
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	))

	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	api := router.PathPrefix("/api").Subrouter()

	setupPublicRoutes(api, handlers)

	protectedRoutes := api.NewRoute().Subrouter()
	protectedRoutes.Use(authMiddleware.Authenticate)
	setupProtectedRoutes(protectedRoutes, handlers)
}

func setupPublicRoutes(router *mux.Router, handlers *handler.Handler) {
	authRoutes := router.PathPrefix("/auth").Subrouter()
	authRoutes.HandleFunc("/login", handlers.Login).Methods(http.MethodPost)
	authRoutes.HandleFunc("/refresh", handlers.RefreshToken).Methods(http.MethodPost)
	authRoutes.HandleFunc("/logout", handlers.Logout).Methods(http.MethodPost)
}

func setupProtectedRoutes(router *mux.Router, handlers *handler.Handler) {
	userRoutes := router.PathPrefix("/users").Subrouter()
	userRoutes.HandleFunc("/me", handlers.GetUser).Methods(http.MethodGet)
	userRoutes.HandleFunc("/me", handlers.UpdateUser).Methods(http.MethodPut)
	userRoutes.HandleFunc("/me/character", handlers.UpdateUserCharacter).Methods(http.MethodPut)

	goalRoutes := router.PathPrefix("/goals").Subrouter()
	goalRoutes.HandleFunc("", handlers.CreateGoal).Methods(http.MethodPost)
	goalRoutes.HandleFunc("", handlers.GetGoal).Methods(http.MethodGet)
	goalRoutes.HandleFunc("/suggestion", handlers.GetGoalSuggestion).Methods(http.MethodPost)
	goalRoutes.HandleFunc("/suggestion/me", handlers.GetAutoGoalSuggestion).Methods(http.MethodGet)

	investmentRoutes := router.PathPrefix("/investments").Subrouter()
	investmentRoutes.HandleFunc("", handlers.GetOpportunities).Methods(http.MethodGet)
	investmentRoutes.HandleFunc("", handlers.CreateOpportunity).Methods(http.MethodPost)

	userInvestmentRoutes := router.PathPrefix("/users/me/investment").Subrouter()
	userInvestmentRoutes.HandleFunc("/", handlers.CreateUserInvestment).Methods(http.MethodPost)
	userInvestmentRoutes.HandleFunc("/", handlers.GetUserInvestments).Methods(http.MethodGet)

	transactionRoutes := router.PathPrefix("/transactions").Subrouter()
	transactionRoutes.HandleFunc("", handlers.CreateTransaction).Methods(http.MethodPost)
	transactionRoutes.HandleFunc("", handlers.GetTransactions).Methods(http.MethodGet)

	gachaRoutes := router.PathPrefix("/gacha").Subrouter()
	gachaRoutes.HandleFunc("/draw", handlers.DrawGacha).Methods(http.MethodPost)
	gachaRoutes.HandleFunc("/preview", handlers.PreviewGachas).Methods(http.MethodGet)

	reportRoutes := router.PathPrefix("/reports").Subrouter()
	reportRoutes.HandleFunc("/finance", handlers.GetReport).Methods(http.MethodGet)
	reportRoutes.HandleFunc("/analysis", handlers.GetReportSummary).Methods(http.MethodGet)
}
