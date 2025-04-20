package redis_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/Financial-Partner/server/internal/entities"
	"github.com/Financial-Partner/server/internal/infrastructure/persistence/redis"
	goredis "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/mock/gomock"
)

func TestTransactionStore(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("GetByUserIdSuccess", func(t *testing.T) {
		mockRedisClient := redis.NewMockRedisClient(ctrl)
		transactionStore := redis.NewTransactionStore(mockRedisClient)

		// Mock data
		userID := primitive.NewObjectID().Hex()

		mockTransactions := []entities.Transaction{
			{
				ID:          primitive.NewObjectID(),
				UserID:      userID,
				Amount:      100,
				Description: "Groceries",
				Date:        time.Now(),
				Category:    "Food",
				Type:        "expense",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				ID:          primitive.NewObjectID(),
				UserID:      userID,
				Amount:      200,
				Description: "Rent",
				Date:        time.Now(),
				Category:    "Housing",
				Type:        "expense",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		}

		// Serialize mockTransactions to JSON
		mockData, _ := json.Marshal(mockTransactions)

		// Mock the Get method to return the serialized JSON data
		mockRedisClient.EXPECT().Get(gomock.Any(), fmt.Sprintf("user:%s:transactions", userID), gomock.Any()).DoAndReturn(
			func(_ context.Context, _ string, dest interface{}) error {
				return json.Unmarshal(mockData, dest)
			},
		)
		transactions, err := transactionStore.GetByUserId(context.Background(), userID)
		require.NoError(t, err)
		assert.NotNil(t, transactions)
	})

	t.Run("GetByUserIdNotFound", func(t *testing.T) {
		mockRedisClient := redis.NewMockRedisClient(ctrl)
		transactionStore := redis.NewTransactionStore(mockRedisClient)

		userID := primitive.NewObjectID().Hex()
		mockRedisClient.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(goredis.Nil)

		transactions, err := transactionStore.GetByUserId(context.Background(), userID)
		require.Error(t, err)
		assert.Nil(t, transactions)
	})

	t.Run("SetByUserIdSuccess", func(t *testing.T) {
		mockRedisClient := redis.NewMockRedisClient(ctrl)
		transactionStore := redis.NewTransactionStore(mockRedisClient)

		// mock data
		userID := primitive.NewObjectID().Hex()
		mockTransactions := []entities.Transaction{
			{
				ID:          primitive.NewObjectID(),
				UserID:      userID,
				Amount:      100,
				Description: "Groceries",
				Date:        time.Now(),
				Category:    "Food",
				Type:        "expense",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				ID:          primitive.NewObjectID(),
				UserID:      userID,
				Amount:      200,
				Description: "Rent",
				Date:        time.Now(),
				Category:    "Housing",
				Type:        "expense",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		}
		mockRedisClient.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

		err := transactionStore.SetByUserId(context.Background(), userID, mockTransactions)
		require.NoError(t, err)
	})

	t.Run("DeleteTransactionSuccess", func(t *testing.T) {
		mockRedisClient := redis.NewMockRedisClient(ctrl)
		transactionStore := redis.NewTransactionStore(mockRedisClient)

		userID := primitive.NewObjectID().Hex()
		mockRedisClient.EXPECT().Delete(gomock.Any(), fmt.Sprintf("user:%s:transactions", userID)).Return(nil)

		err := transactionStore.DeleteByUserId(context.Background(), userID)
		require.NoError(t, err)
	})
}
