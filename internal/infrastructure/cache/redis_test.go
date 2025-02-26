package cache_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	redismock "github.com/go-redis/redismock/v9"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/Financial-Partner/server/internal/config"
	"github.com/Financial-Partner/server/internal/infrastructure/cache"
)

func TestNewClient(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("success", func(t *testing.T) {
		db, mock := redismock.NewClientMock()
		mock.ExpectPing().SetVal("PONG")

		newRedisClientFunc := func(opts *redis.Options) redis.UniversalClient {
			return db
		}

		client, err := cache.NewClient(&config.Config{}, newRedisClientFunc)
		require.NoError(t, err)
		require.NotNil(t, client)
	})

	t.Run("error", func(t *testing.T) {
		db, mock := redismock.NewClientMock()
		mock.ExpectPing().SetErr(errors.New("ping error"))

		newRedisClientFunc := func(opts *redis.Options) redis.UniversalClient {
			return db
		}

		client, err := cache.NewClient(&config.Config{}, newRedisClientFunc)
		require.Error(t, err)
		require.Nil(t, client)
	})
}

func TestSet(t *testing.T) {
	db, mock := redismock.NewClientMock()
	client := cache.NewWithRedisClient(db)

	t.Run("Marshal error", func(t *testing.T) {
		ctx := context.Background()
		err := client.Set(ctx, "test-key", func() {}, time.Minute)
		assert.Error(t, err)
	})

	t.Run("Set success", func(t *testing.T) {
		ctx := context.Background()
		testValue := map[string]interface{}{"foo": "bar"}
		data, err := json.Marshal(testValue)
		assert.NoError(t, err)

		mock.ExpectSet("test-key", data, time.Minute).SetVal("OK")

		err = client.Set(ctx, "test-key", testValue, time.Minute)
		assert.NoError(t, err)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unmet expectations: %s", err)
		}
	})
}

func TestGet(t *testing.T) {
	db, mock := redismock.NewClientMock()
	client := cache.NewWithRedisClient(db)

	t.Run("Not found", func(t *testing.T) {
		ctx := context.Background()
		mock.ExpectGet("test-key").SetErr(redis.Nil)

		var result map[string]interface{}
		err := client.Get(ctx, "test-key", &result)
		assert.Error(t, err)
	})

	t.Run("Unmarshal error", func(t *testing.T) {
		ctx := context.Background()
		err := client.Get(ctx, "test-key", func() {})
		assert.Error(t, err)
	})

	t.Run("Get success", func(t *testing.T) {
		ctx := context.Background()
		expectedValue := map[string]interface{}{"foo": "bar"}
		data, err := json.Marshal(expectedValue)
		assert.NoError(t, err)

		mock.ExpectGet("test-key").SetVal(string(data))

		var result map[string]interface{}
		err = client.Get(ctx, "test-key", &result)
		assert.NoError(t, err)
		assert.Equal(t, expectedValue, result)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unmet expectations: %s", err)
		}
	})
}

func TestDelete(t *testing.T) {
	db, mock := redismock.NewClientMock()
	client := cache.NewWithRedisClient(db)

	t.Run("Delete failed", func(t *testing.T) {
		ctx := context.Background()
		mock.ExpectDel("test-key").SetErr(redis.Nil)

		err := client.Delete(ctx, "test-key")
		assert.Error(t, err)
	})

	t.Run("Delete success", func(t *testing.T) {
		ctx := context.Background()

		mock.ExpectDel("test-key").SetVal(1)

		err := client.Delete(ctx, "test-key")
		assert.NoError(t, err)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unmet expectations: %s", err)
		}
	})
}
