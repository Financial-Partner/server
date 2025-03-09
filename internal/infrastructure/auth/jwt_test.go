package auth_test

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Financial-Partner/server/internal/infrastructure/auth"
)

func TestJWTManager(t *testing.T) {
	secretKey := "test-secret-key"
	accessExpiry := 15 * time.Minute
	refreshExpiry := 24 * time.Hour
	jwtManager := auth.NewJWTManager(secretKey, accessExpiry, refreshExpiry)

	t.Run("GenerateAccessToken", func(t *testing.T) {
		email := "test@example.com"

		token, expiryTime, err := jwtManager.GenerateAccessToken(email)

		require.NoError(t, err)
		require.NotEmpty(t, token)

		expectedExpiry := time.Now().Add(accessExpiry)
		assert.WithinDuration(t, expectedExpiry, expiryTime, 2*time.Second)

		claims, err := jwtManager.ValidateToken(token)
		require.NoError(t, err)
		assert.Equal(t, email, claims.Email)
	})

	t.Run("GenerateRefreshToken", func(t *testing.T) {
		email := "test@example.com"

		token, expiryTime, err := jwtManager.GenerateRefreshToken(email)

		require.NoError(t, err)
		require.NotEmpty(t, token)

		expectedExpiry := time.Now().Add(refreshExpiry)
		assert.WithinDuration(t, expectedExpiry, expiryTime, 2*time.Second)

		claims, err := jwtManager.ValidateToken(token)
		require.NoError(t, err)
		assert.Equal(t, email, claims.Email)
	})

	t.Run("ValidateToken_Valid", func(t *testing.T) {
		email := "test@example.com"
		expiresAt := time.Now().Add(time.Hour)

		claims := &auth.Claims{
			Email: email,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(expiresAt),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte(secretKey))
		require.NoError(t, err)

		parsedClaims, err := jwtManager.ValidateToken(tokenString)

		require.NoError(t, err)
		assert.Equal(t, email, parsedClaims.Email)
	})

	t.Run("ValidateToken_Expired", func(t *testing.T) {
		email := "test@example.com"
		expiresAt := time.Now().Add(-time.Hour)

		claims := &auth.Claims{
			Email: email,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(expiresAt),
				IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte(secretKey))
		require.NoError(t, err)

		parsedClaims, err := jwtManager.ValidateToken(tokenString)

		assert.Error(t, err)
		assert.Nil(t, parsedClaims)
		assert.Contains(t, err.Error(), "token is expired")
	})

	t.Run("ValidateToken_InvalidSignature", func(t *testing.T) {
		email := "test@example.com"
		claims := &auth.Claims{
			Email: email,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte("wrong-secret-key"))
		require.NoError(t, err)

		parsedClaims, err := jwtManager.ValidateToken(tokenString)

		assert.Error(t, err)
		assert.Nil(t, parsedClaims)
		assert.Contains(t, err.Error(), "signature is invalid")
	})

	t.Run("ValidateToken_InvalidAlgorithm", func(t *testing.T) {
		invalidToken := "invalid.token.format"

		parsedClaims, err := jwtManager.ValidateToken(invalidToken)

		assert.Error(t, err)
		assert.Nil(t, parsedClaims)
	})

	t.Run("ValidateToken_MalformedToken", func(t *testing.T) {
		malformedToken := "this.is.not.a.valid.jwt"

		parsedClaims, err := jwtManager.ValidateToken(malformedToken)

		assert.Error(t, err)
		assert.Nil(t, parsedClaims)
	})

	t.Run("ValidateToken_EmptyToken", func(t *testing.T) {
		emptyToken := ""

		parsedClaims, err := jwtManager.ValidateToken(emptyToken)

		assert.Error(t, err)
		assert.Nil(t, parsedClaims)
	})
}
