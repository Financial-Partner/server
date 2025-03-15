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
	apiBaseURL url.URL,
) {
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
}

func setupProtectedRoutes(router *mux.Router, handlers *handler.Handler) {
	userRoutes := router.PathPrefix("/users").Subrouter()
	// WHY: create user will happen after exchange token by firebase token, so it belongs to protected route
	userRoutes.HandleFunc("", handlers.CreateUser).Methods(http.MethodPost)
	userRoutes.HandleFunc("/me", handlers.GetUser).Methods(http.MethodGet)
	userRoutes.HandleFunc("/me", handlers.UpdateUser).Methods(http.MethodPut)

	goalRoutes := router.PathPrefix("/goals").Subrouter()
	goalRoutes.HandleFunc("", handlers.CreateGoal).Methods(http.MethodPost)
	goalRoutes.HandleFunc("", handlers.GetGoal).Methods(http.MethodGet)
	goalRoutes.HandleFunc("/suggestion", handlers.GetGoalSuggestion).Methods(http.MethodPost)
	goalRoutes.HandleFunc("/suggestion/me", handlers.GetAutoGoalSuggestion).Methods(http.MethodGet)
}
