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
		userID := primitive.NewObjectID()

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
		mockRedisClient.EXPECT().Get(gomock.Any(), fmt.Sprintf("user:%s:transactions", userID.Hex()), gomock.Any()).DoAndReturn(
			func(_ context.Context, _ string, dest interface{}) error {
				return json.Unmarshal(mockData, dest)
			},
		)
		transactions, err := transactionStore.GetByUserId(context.Background(), userID.Hex())
		require.NoError(t, err)
		assert.NotNil(t, transactions)
	})

	t.Run("GetByUserIdNotFound", func(t *testing.T) {
		mockRedisClient := redis.NewMockRedisClient(ctrl)
		transactionStore := redis.NewTransactionStore(mockRedisClient)

		userID := primitive.NewObjectID()
		mockRedisClient.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(goredis.Nil)

		transactions, err := transactionStore.GetByUserId(context.Background(), userID.Hex())
		require.Error(t, err)
		assert.Nil(t, transactions)
	})

	t.Run("SetMultipleByUserIdWithExistingTransactions", func(t *testing.T) {
		mockRedisClient := redis.NewMockRedisClient(ctrl)
		transactionStore := redis.NewTransactionStore(mockRedisClient)

		// Mock data
		userID := primitive.NewObjectID()
		existingTransactions := []entities.Transaction{
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
		}
		newTransactions := []entities.Transaction{
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
			{
				ID:          primitive.NewObjectID(),
				UserID:      userID,
				Amount:      50,
				Description: "Utilities",
				Date:        time.Now(),
				Category:    "Bills",
				Type:        "expense",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		}

		// Serialize existing transactions to JSON
		mockData, _ := json.Marshal(existingTransactions)

		// Mock the Get method to return the existing transactions
		mockRedisClient.EXPECT().Get(
			gomock.Any(),
			fmt.Sprintf("user:%s:transactions", userID.Hex()),
			gomock.Any(),
		).DoAndReturn(
			func(_ context.Context, _ string, dest interface{}) error {
				return json.Unmarshal(mockData, dest)
			},
		)

		// Mock the Set method to save the updated transactions
		mockRedisClient.EXPECT().Set(
			gomock.Any(),
			fmt.Sprintf("user:%s:transactions", userID.Hex()),
			gomock.Any(),
			gomock.Any(),
		).DoAndReturn(
			func(_ context.Context, _ string, value interface{}, _ time.Duration) error {
				return nil
			},
		)

		// Call SetMultipleByUserId
		err := transactionStore.SetMultipleByUserId(context.Background(), userID.Hex(), newTransactions)
		require.NoError(t, err)
	})

	t.Run("DeleteTransactionSuccess", func(t *testing.T) {
		mockRedisClient := redis.NewMockRedisClient(ctrl)
		transactionStore := redis.NewTransactionStore(mockRedisClient)

		userID := primitive.NewObjectID()
		mockRedisClient.EXPECT().Delete(gomock.Any(), fmt.Sprintf("user:%s:transactions", userID.Hex())).Return(nil)

		err := transactionStore.DeleteByUserId(context.Background(), userID.Hex())
		require.NoError(t, err)
	})
}
