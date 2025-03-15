package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/Financial-Partner/server/internal/contextutil"
	"github.com/Financial-Partner/server/internal/infrastructure/auth"
	"github.com/Financial-Partner/server/internal/infrastructure/logger"
)

//go:generate mockgen -source=auth.go -destination=auth_mock.go -package=middleware

type JWTValidator interface {
	ValidateToken(tokenString string) (*auth.Claims, error)
}

type AuthMiddleware struct {
	jwtValidator JWTValidator
	log          logger.Logger
}

func NewAuthMiddleware(jwt JWTValidator, log logger.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		jwtValidator: jwt,
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

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := m.jwtValidator.ValidateToken(tokenString)
		if err != nil {
			m.log.WithError(err).Warnf("JWT token validation failed")
			http.Error(w, "Invalid authorization token", http.StatusUnauthorized)
			return
		}

		if claims.ID == "" {
			m.log.Warnf("No user ID found in token claims")
			http.Error(w, "Unable to retrieve user ID", http.StatusUnauthorized)
			return
		}

		if claims.Email == "" {
			m.log.Warnf("No email found in token claims")
			http.Error(w, "Unable to retrieve user email", http.StatusUnauthorized)
			return
		}

		m.log.WithFields(map[string]any{
			"email": claims.Email,
			"id":    claims.ID,
		}).Infof("User authenticated successfully")
		ctx := context.WithValue(r.Context(), contextutil.UserIDKey, claims.ID)
		ctx = context.WithValue(ctx, contextutil.UserEmailKey, claims.Email)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
