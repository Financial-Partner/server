package redis_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/Financial-Partner/server/internal/entities"
	"github.com/Financial-Partner/server/internal/infrastructure/persistence/redis"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/mock/gomock"
)

func TestInvestmentStore(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("SetOpportunitiesSuccess", func(t *testing.T) {
		mockRedisClient := redis.NewMockRedisClient(ctrl)
		investmentStore := redis.NewInvestmentStore(mockRedisClient)

		// Mock data
		userID := primitive.NewObjectID().Hex()
		opportunities := []entities.Opportunity{
			{
				ID:          primitive.NewObjectID(),
				Title:       "real estate",
				Description: "Invest in real estate",
				Tags:        []string{"high risk", "long term"},
				IsIncrease:  true,
				Variation:   10,
				Duration:    "1 year",
				MinAmount:   1000,
				CreatedAt:   time.Date(2023, time.January, 31, 0, 0, 0, 0, time.UTC),
				UpdatedAt:   time.Date(2023, time.January, 31, 0, 0, 0, 0, time.UTC),
			},
		}

		mockRedisClient.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

		err := investmentStore.SetOpportunities(context.Background(), userID, opportunities)
		require.NoError(t, err)
	})

	t.Run("SetInvestmentsSuccess", func(t *testing.T) {
		mockRedisClient := redis.NewMockRedisClient(ctrl)
		investmentStore := redis.NewInvestmentStore(mockRedisClient)

		// Mock data
		userID := primitive.NewObjectID().Hex()
		investments := []entities.Investment{
			{
				ID:            primitive.NewObjectID(),
				UserID:        primitive.NewObjectID(),
				OpportunityID: primitive.NewObjectID(),
				Amount:        1000,
				CreatedAt:     time.Date(2023, time.January, 31, 0, 0, 0, 0, time.UTC),
				UpdatedAt:     time.Date(2023, time.January, 31, 0, 0, 0, 0, time.UTC),
			},
		}

		mockRedisClient.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

		err := investmentStore.SetInvestments(context.Background(), userID, investments)
		require.NoError(t, err)
	})

	t.Run("DeleteInvestmentsSuccess", func(t *testing.T) {
		mockRedisClient := redis.NewMockRedisClient(ctrl)
		investmentStore := redis.NewInvestmentStore(mockRedisClient)

		userID := primitive.NewObjectID().Hex()
		mockRedisClient.EXPECT().Delete(gomock.Any(), fmt.Sprintf("user:%s:investments", userID)).Return(nil)

		err := investmentStore.DeleteInvestments(context.Background(), userID)
		require.NoError(t, err)
	})

	t.Run("DeleteOpportunitiesSuccess", func(t *testing.T) {
		mockRedisClient := redis.NewMockRedisClient(ctrl)
		investmentStore := redis.NewInvestmentStore(mockRedisClient)

		userID := primitive.NewObjectID().Hex()
		mockRedisClient.EXPECT().Delete(gomock.Any(), fmt.Sprintf("user:%s:opportunities", userID)).Return(nil)

		err := investmentStore.DeleteOpportunities(context.Background(), userID)
		require.NoError(t, err)
	})

	t.Run("GetInvestmentsSuccess", func(t *testing.T) {
		mockRedisClient := redis.NewMockRedisClient(ctrl)
		investmentStore := redis.NewInvestmentStore(mockRedisClient)

		// Mock data
		userID := primitive.NewObjectID().Hex()

		mockInvestments := []entities.Investment{
			{
				ID:            primitive.NewObjectID(),
				UserID:        primitive.NewObjectID(),
				OpportunityID: primitive.NewObjectID(),
				Amount:        1000,
				CreatedAt:     time.Date(2023, time.January, 31, 0, 0, 0, 0, time.UTC),
				UpdatedAt:     time.Date(2023, time.January, 31, 0, 0, 0, 0, time.UTC),
			},
		}

		// Serialize mockInvestments to JSON
		mockData, _ := json.Marshal(mockInvestments)

		// Mock the Get method to return the serialized JSON data
		mockRedisClient.EXPECT().Get(gomock.Any(), fmt.Sprintf("user:%s:investments", userID), gomock.Any()).DoAndReturn(
			func(_ context.Context, _ string, dest interface{}) error {
				return json.Unmarshal(mockData, dest)
			},
		)
		investments, err := investmentStore.GetInvestments(context.Background(), userID)
		require.NoError(t, err)
		assert.NotNil(t, investments)
	})

	t.Run("GetOpportunitiesSuccess", func(t *testing.T) {
		mockRedisClient := redis.NewMockRedisClient(ctrl)
		investmentStore := redis.NewInvestmentStore(mockRedisClient)

		// Mock data
		userID := primitive.NewObjectID().Hex()

		mockOpportunities := []entities.Opportunity{
			{
				ID:          primitive.NewObjectID(),
				Title:       "real estate",
				Description: "Invest in real estate",
				Tags:        []string{"high risk", "long term"},
				IsIncrease:  true,
				Variation:   10,
				Duration:    "1 year",
				MinAmount:   1000,
				CreatedAt:   time.Date(2023, time.January, 31, 0, 0, 0, 0, time.UTC),
				UpdatedAt:   time.Date(2023, time.January, 31, 0, 0, 0, 0, time.UTC),
			},
		}

		// Serialize mockOpportunities to JSON
		mockData, _ := json.Marshal(mockOpportunities)

		// Mock the Get method to return the serialized JSON data
		mockRedisClient.EXPECT().Get(gomock.Any(), fmt.Sprintf("user:%s:opportunities", userID), gomock.Any()).DoAndReturn(
			func(_ context.Context, _ string, dest interface{}) error {
				return json.Unmarshal(mockData, dest)
			},
		)
		opportunities, err := investmentStore.GetOpportunities(context.Background(), userID)
		require.NoError(t, err)
		assert.NotNil(t, opportunities)
	})
}
