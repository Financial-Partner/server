package auth_test

import (
	"context"
	"testing"

	fbAuth "firebase.google.com/go/v4/auth"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/Financial-Partner/server/internal/config"
	"github.com/Financial-Partner/server/internal/infrastructure/auth"
	"github.com/Financial-Partner/server/internal/infrastructure/auth/mocks"
)

func TestNewClient(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// unit test should	not test correct behavior

	_, err := auth.NewClient(context.Background(), &config.Config{
		Firebase: config.Firebase{
			ProjectID:      "test-project-id",
			CredentialFile: "credentials.json",
		},
	})
	assert.Error(t, err)
}

func TestVerifyToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuth := mocks.NewMockFirebaseAuth(ctrl)
	client := auth.NewWithAuth(mockAuth)

	mockAuth.EXPECT().VerifyIDToken(gomock.Any(), gomock.Any()).Return(&fbAuth.Token{
		Claims: map[string]interface{}{
			"email": "test@example.com",
		},
	}, nil)

	token, err := client.VerifyToken(context.Background(), "test-token")
	assert.NoError(t, err)
	assert.NotNil(t, token)
}
