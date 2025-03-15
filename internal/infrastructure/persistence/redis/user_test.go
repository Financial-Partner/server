package redis_test

import (
	"context"
	"testing"

	"github.com/Financial-Partner/server/internal/entities"
	"github.com/Financial-Partner/server/internal/infrastructure/persistence/redis"
	goredis "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestUserStore(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("GetUserSuccess", func(t *testing.T) {
		mockRedisClient := redis.NewMockRedisClient(ctrl)
		userStore := redis.NewUserStore(mockRedisClient)

		mockRedisClient.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

		user, err := userStore.Get(context.Background(), "test@example.com")
		require.NoError(t, err)
		assert.NotNil(t, user)
	})

	t.Run("GetUserNotFound", func(t *testing.T) {
		mockRedisClient := redis.NewMockRedisClient(ctrl)
		userStore := redis.NewUserStore(mockRedisClient)

		mockRedisClient.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(goredis.Nil)

		user, err := userStore.Get(context.Background(), "test@example.com")
		require.Error(t, err)
		assert.Nil(t, user)
	})

	t.Run("SetUserSuccess", func(t *testing.T) {
		mockRedisClient := redis.NewMockRedisClient(ctrl)
		userStore := redis.NewUserStore(mockRedisClient)

		mockRedisClient.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

		err := userStore.Set(context.Background(), &entities.User{
			Email: "test@example.com",
			Name:  "Test User",
		})
		require.NoError(t, err)
	})

	t.Run("DeleteUserSuccess", func(t *testing.T) {
		mockRedisClient := redis.NewMockRedisClient(ctrl)
		userStore := redis.NewUserStore(mockRedisClient)

		mockRedisClient.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil)

		err := userStore.Delete(context.Background(), "test@example.com")
		require.NoError(t, err)
	})
}
