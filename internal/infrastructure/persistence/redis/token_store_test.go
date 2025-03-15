package redis_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/Financial-Partner/server/internal/infrastructure/persistence/redis"
)

func TestTokenStore(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRedisClient := redis.NewMockRedisClient(ctrl)
	tokenStore := redis.NewTokenStore(mockRedisClient)

	testEmail := "test@example.com"
	testRefreshToken := "refresh-token-123"
	testExpiry := time.Now().Add(24 * time.Hour)
	testKey := "refresh_token:" + testRefreshToken

	t.Run("SaveRefreshToken Success", func(t *testing.T) {
		mockRedisClient.EXPECT().
			Set(gomock.Any(), testKey, testEmail, gomock.Any()).
			DoAndReturn(func(_ context.Context, _ string, _ interface{}, ttl time.Duration) error {
				expectedTTL := time.Until(testExpiry)
				assert.InDelta(t, expectedTTL.Seconds(), ttl.Seconds(), 1.0)
				return nil
			})

		err := tokenStore.SaveRefreshToken(context.Background(), testEmail, testRefreshToken, testExpiry)
		assert.NoError(t, err)
	})

	t.Run("SaveRefreshToken Error", func(t *testing.T) {
		mockRedisClient.EXPECT().
			Set(gomock.Any(), testKey, testEmail, gomock.Any()).
			Return(errors.New("redis error"))

		err := tokenStore.SaveRefreshToken(context.Background(), testEmail, testRefreshToken, testExpiry)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "redis error")
	})

	t.Run("GetRefreshToken Success", func(t *testing.T) {
		mockRedisClient.EXPECT().
			Get(gomock.Any(), testKey, gomock.Any()).
			DoAndReturn(func(_ context.Context, _ string, dest interface{}) error {
				emailPtr, ok := dest.(*string)
				if !ok {
					return errors.New("dest is not a *string")
				}
				*emailPtr = testEmail
				return nil
			})

		email, err := tokenStore.GetRefreshToken(context.Background(), testRefreshToken)
		assert.NoError(t, err)
		assert.Equal(t, testEmail, email)
	})

	t.Run("GetRefreshToken NotFound", func(t *testing.T) {
		mockRedisClient.EXPECT().
			Get(gomock.Any(), testKey, gomock.Any()).
			Return(errors.New("key not found"))

		email, err := tokenStore.GetRefreshToken(context.Background(), testRefreshToken)
		assert.Error(t, err)
		assert.Equal(t, "", email)
		assert.Contains(t, err.Error(), "key not found")
	})

	t.Run("DeleteRefreshToken Success", func(t *testing.T) {
		mockRedisClient.EXPECT().
			Delete(gomock.Any(), testKey).
			Return(nil)

		err := tokenStore.DeleteRefreshToken(context.Background(), testRefreshToken)
		assert.NoError(t, err)
	})

	t.Run("DeleteRefreshToken Error", func(t *testing.T) {
		mockRedisClient.EXPECT().
			Delete(gomock.Any(), testKey).
			Return(errors.New("redis error"))

		err := tokenStore.DeleteRefreshToken(context.Background(), testRefreshToken)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "redis error")
	})

	t.Run("Empty RefreshToken", func(t *testing.T) {
		emptyToken := ""
		emptyKey := "refresh_token:" + emptyToken

		mockRedisClient.EXPECT().
			Set(gomock.Any(), emptyKey, testEmail, gomock.Any()).
			Return(nil)

		err := tokenStore.SaveRefreshToken(context.Background(), testEmail, emptyToken, testExpiry)
		assert.NoError(t, err)

		mockRedisClient.EXPECT().
			Get(gomock.Any(), emptyKey, gomock.Any()).
			DoAndReturn(func(_ context.Context, _ string, dest interface{}) error {
				emailPtr, ok := dest.(*string)
				if !ok {
					return errors.New("dest is not a *string")
				}
				*emailPtr = testEmail
				return nil
			})

		email, err := tokenStore.GetRefreshToken(context.Background(), emptyToken)
		assert.NoError(t, err)
		assert.Equal(t, testEmail, email)

		mockRedisClient.EXPECT().
			Delete(gomock.Any(), emptyKey).
			Return(nil)

		err = tokenStore.DeleteRefreshToken(context.Background(), emptyToken)
		assert.NoError(t, err)
	})

	t.Run("Empty Email", func(t *testing.T) {
		emptyEmail := ""

		mockRedisClient.EXPECT().
			Set(gomock.Any(), testKey, emptyEmail, gomock.Any()).
			Return(nil)

		err := tokenStore.SaveRefreshToken(context.Background(), emptyEmail, testRefreshToken, testExpiry)
		assert.NoError(t, err)

		mockRedisClient.EXPECT().
			Get(gomock.Any(), testKey, gomock.Any()).
			DoAndReturn(func(_ context.Context, _ string, dest interface{}) error {
				emailPtr, ok := dest.(*string)
				if !ok {
					return errors.New("dest is not a *string")
				}
				*emailPtr = emptyEmail
				return nil
			})

		email, err := tokenStore.GetRefreshToken(context.Background(), testRefreshToken)
		assert.NoError(t, err)
		assert.Equal(t, emptyEmail, email)
	})

	t.Run("Context Canceled", func(t *testing.T) {
		canceledCtx, cancel := context.WithCancel(context.Background())
		cancel()

		mockRedisClient.EXPECT().
			Set(gomock.Any(), testKey, testEmail, gomock.Any()).
			Return(context.Canceled)

		err := tokenStore.SaveRefreshToken(canceledCtx, testEmail, testRefreshToken, testExpiry)
		assert.Error(t, err)
		assert.Equal(t, context.Canceled, err)

		mockRedisClient.EXPECT().
			Get(gomock.Any(), testKey, gomock.Any()).
			Return(context.Canceled)

		email, err := tokenStore.GetRefreshToken(canceledCtx, testRefreshToken)
		assert.Error(t, err)
		assert.Equal(t, "", email)
		assert.Equal(t, context.Canceled, err)

		mockRedisClient.EXPECT().
			Delete(gomock.Any(), testKey).
			Return(context.Canceled)

		err = tokenStore.DeleteRefreshToken(canceledCtx, testRefreshToken)
		assert.Error(t, err)
		assert.Equal(t, context.Canceled, err)
	})
}
