package contextutil_test

import (
	"context"
	"testing"

	"github.com/Financial-Partner/server/internal/contextutil"
	"github.com/stretchr/testify/assert"
)

func TestGetUserEmail(t *testing.T) {
	t.Run("should return user email", func(t *testing.T) {
		ctx := context.Background()
		email := "test@example.com"
		ctx = context.WithValue(ctx, contextutil.UserEmailKey, email)

		gotEmail, ok := contextutil.GetUserEmail(ctx)
		assert.True(t, ok)
		assert.Equal(t, email, gotEmail)
	})

	t.Run("should return false if user email is not present", func(t *testing.T) {
		ctx := context.Background()
		gotEmail, ok := contextutil.GetUserEmail(ctx)
		assert.False(t, ok)
		assert.Empty(t, gotEmail)
	})
}
