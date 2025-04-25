package mongodb_test

import (
	"context"
	"testing"
	"time"

	"github.com/Financial-Partner/server/internal/entities"
	"github.com/Financial-Partner/server/internal/infrastructure/persistence/mongodb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestMongoTransactionRepository(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	testUserID := primitive.NewObjectID()
	testTransactions := []entities.Transaction{
		{
			ID:          primitive.NewObjectID(),
			UserID:      primitive.NewObjectID(),
			Amount:      100,
			Description: "Dinner",
			Date:        time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC),
			Category:    "Food",
			Type:        "expense",
			CreatedAt:   time.Date(2023, time.January, 31, 0, 0, 0, 0, time.UTC),
			UpdatedAt:   time.Date(2023, time.January, 31, 0, 0, 0, 0, time.UTC),
		},
		{
			ID:          primitive.NewObjectID(),
			UserID:      primitive.NewObjectID(),
			Amount:      200,
			Description: "Rent",
			Date:        time.Date(2023, time.January, 2, 0, 0, 0, 0, time.UTC),
			Category:    "Housing",
			Type:        "expense",
			CreatedAt:   time.Date(2023, time.January, 30, 0, 0, 0, 0, time.UTC),
			UpdatedAt:   time.Date(2023, time.January, 30, 0, 0, 0, 0, time.UTC),
		},
	}

	// Convert transactions to BSON documents
	var testTransactionDocs []bson.D
	for _, transaction := range testTransactions {
		transactionBSON, err := bson.Marshal(transaction)
		require.NoError(t, err)
		var transactionDoc bson.D
		err = bson.Unmarshal(transactionBSON, &transactionDoc)
		require.NoError(t, err)
		testTransactionDocs = append(testTransactionDocs, transactionDoc)
	}

	t.Run("FindByUserId", func(t *testing.T) {
		mt.Run("success", func(mt *mtest.T) {
			mt.AddMockResponses(
				mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, testTransactionDocs...),
				mtest.CreateCursorResponse(0, "foo.bar", mtest.NextBatch), // Simulate end of cursor
			)
			repo := mongodb.NewTransactionRepository(mt.Client.Database("testdb"))
			result, err := repo.FindByUserId(context.Background(), testUserID)
			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.Len(t, result, len(testTransactions))

			// Validate each transaction
			for i, transaction := range result {
				assert.Equal(t, testTransactions[i], transaction)
			}
		})
		mt.Run("not found", func(mt *mtest.T) {
			mt.AddMockResponses(mtest.CreateCursorResponse(0, "foo.bar", mtest.FirstBatch))
			repo := mongodb.NewTransactionRepository(mt.DB)
			result, err := repo.FindByUserId(context.Background(), testUserID)
			assert.Nil(t, err)
			assert.Nil(t, result)
		})
		mt.Run("database error", func(mt *mtest.T) {
			mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
				Code:    11000,
				Message: "database error",
			}))
			repo := mongodb.NewTransactionRepository(mt.DB)
			result, err := repo.FindByUserId(context.Background(), testUserID)
			assert.Error(t, err)
			assert.Nil(t, result)
		})
	})

	t.Run("Create", func(t *testing.T) {
		mt.Run("success", func(mt *mtest.T) {
			mt.AddMockResponses(mtest.CreateSuccessResponse())
			repo := mongodb.NewTransactionRepository(mt.DB)
			result, err := repo.Create(context.Background(), &(testTransactions[0]))
			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, testTransactions[0], *result)
		})
		mt.Run("error", func(mt *mtest.T) {
			mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
				Code:    11000,
				Message: "duplicate key error",
			}))
			repo := mongodb.NewTransactionRepository(mt.DB)
			result, err := repo.Create(context.Background(), &(testTransactions[0]))
			assert.Error(t, err)
			assert.Nil(t, result)
		})
	})
}
