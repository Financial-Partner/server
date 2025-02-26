package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/Financial-Partner/server/internal/contextutil"
	"github.com/Financial-Partner/server/internal/infrastructure/auth"
	"github.com/Financial-Partner/server/internal/infrastructure/logger"
)

//go:generate mockgen -source=auth.go -destination=mocks/auth_mock.go -package=mocks

type AuthClient interface {
	VerifyToken(ctx context.Context, idToken string) (*auth.Token, error)
}

type AuthMiddleware struct {
	firebaseAuth AuthClient
	log          logger.Logger
}

func NewAuthMiddleware(fa AuthClient, log logger.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		firebaseAuth: fa,
		log:          log,
	}
}

func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			m.log.Warnf("Authentication failed: no token provided")
			http.Error(w, "Authorization token required", http.StatusUnauthorized)
			return
		}

		idToken := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := m.firebaseAuth.VerifyToken(r.Context(), idToken)
		if err != nil {
			m.log.WithError(err).Warnf("Token verification failed")
			http.Error(w, "Invalid authorization token", http.StatusUnauthorized)
			return
		}

		email, ok := token.Claims["email"].(string)
		if !ok {
			m.log.Warnf("No email found in token claims")
			http.Error(w, "Unable to retrieve user email", http.StatusUnauthorized)
			return
		}

		m.log.WithField("email", email).Infof("User authenticated successfully")
		ctx := context.WithValue(r.Context(), contextutil.UserEmailKey, email)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
