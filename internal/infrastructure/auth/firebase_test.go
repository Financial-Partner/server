package auth_test

import (
	"context"
	"errors"
	"testing"

	fbAuth "firebase.google.com/go/v4/auth"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/Financial-Partner/server/internal/config"
	"github.com/Financial-Partner/server/internal/infrastructure/auth"
)

func TestNewClient(t *testing.T) {
	t.Run("Error creating client with invalid credentials", func(t *testing.T) {
		_, err := auth.NewClient(context.Background(), &config.Config{
			Firebase: config.Firebase{
				ProjectID:      "test-project-id",
				CredentialFile: "credentials.json",
			},
		})
		assert.Error(t, err)
	})

	t.Run("Success with mock auth", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAuth := auth.NewMockFirebaseAuth(ctrl)
		client := auth.NewWithAuth(mockAuth)
		assert.NotNil(t, client)
	})
}

func TestVerifyToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("Success case", func(t *testing.T) {
		mockAuth := auth.NewMockFirebaseAuth(ctrl)
		client := auth.NewWithAuth(mockAuth)

		expectedToken := &fbAuth.Token{
			Claims: map[string]interface{}{
				"email": "test@example.com",
				"name":  "Test User",
			},
		}

		mockAuth.EXPECT().
			VerifyIDToken(gomock.Any(), "valid-token").
			Return(expectedToken, nil)

		token, err := client.VerifyToken(context.Background(), "valid-token")

		assert.NoError(t, err)
		assert.Equal(t, expectedToken, token)
		assert.Equal(t, "test@example.com", token.Claims["email"])
		assert.Equal(t, "Test User", token.Claims["name"])
	})

	t.Run("Error case - invalid token", func(t *testing.T) {
		mockAuth := auth.NewMockFirebaseAuth(ctrl)
		client := auth.NewWithAuth(mockAuth)

		mockAuth.EXPECT().
			VerifyIDToken(gomock.Any(), "invalid-token").
			Return(nil, errors.New("invalid token"))

		token, err := client.VerifyToken(context.Background(), "invalid-token")

		assert.Error(t, err)
		assert.Nil(t, token)
		assert.Contains(t, err.Error(), "invalid token")
	})

	t.Run("Error case - expired token", func(t *testing.T) {
		mockAuth := auth.NewMockFirebaseAuth(ctrl)
		client := auth.NewWithAuth(mockAuth)

		mockAuth.EXPECT().
			VerifyIDToken(gomock.Any(), "expired-token").
			Return(nil, errors.New("token has expired"))

		token, err := client.VerifyToken(context.Background(), "expired-token")

		assert.Error(t, err)
		assert.Nil(t, token)
		assert.Contains(t, err.Error(), "token has expired")
	})
}
