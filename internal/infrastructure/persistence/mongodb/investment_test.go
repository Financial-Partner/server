package mongodb_test

import (
	"context"
	"testing"
	"time"

	"github.com/Financial-Partner/server/internal/entities"
	"github.com/Financial-Partner/server/internal/infrastructure/persistence/mongodb"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestMongoInvestmentRepository(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	testUserID := primitive.NewObjectID().Hex()

	testInvestment := &entities.Investment{
		ID:            primitive.NewObjectID(),
		UserID:        testUserID,
		OpportunityID: primitive.NewObjectID().Hex(),
		Amount:        1000,
		CreatedAt:     time.Date(2023, time.January, 31, 0, 0, 0, 0, time.UTC),
		UpdatedAt:     time.Date(2023, time.January, 31, 0, 0, 0, 0, time.UTC),
	}
	testInvestments := []entities.Investment{
		*testInvestment,
	}

	testOpportunity := &entities.Opportunity{
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
	}
	testOpportunities := []entities.Opportunity{
		*testOpportunity,
	}

	var testInvestmentDocs []bson.D
	for _, investment := range testInvestments {
		investmentBSON, err := bson.Marshal(investment)
		assert.NoError(t, err)
		var investmentDoc bson.D
		err = bson.Unmarshal(investmentBSON, &investmentDoc)
		assert.NoError(t, err)
		testInvestmentDocs = append(testInvestmentDocs, investmentDoc)
	}

	var testOpportunityDocs []bson.D
	for _, opportunity := range testOpportunities {
		opportunityBSON, err := bson.Marshal(opportunity)
		assert.NoError(t, err)
		var opportunityDoc bson.D
		err = bson.Unmarshal(opportunityBSON, &opportunityDoc)
		assert.NoError(t, err)
		testOpportunityDocs = append(testOpportunityDocs, opportunityDoc)
	}

	t.Run("CreateInvestment", func(t *testing.T) {
		mt.Run("error", func(mt *mtest.T) {
			mt.AddMockResponses(
				mtest.CreateCommandErrorResponse(mtest.CommandError{
					Code:    11000,
					Message: "Duplicate key error",
				}),
			)

			repo := mongodb.NewInvestmentRepository(mt.DB)
			result, err := repo.CreateInvestment(context.Background(), testInvestment)
			assert.Error(t, err)
			assert.Nil(t, result)
		})

		mt.Run("success", func(mt *mtest.T) {
			mt.AddMockResponses(
				mtest.CreateSuccessResponse(),
			)

			repo := mongodb.NewInvestmentRepository(mt.DB)
			result, err := repo.CreateInvestment(context.Background(), testInvestment)
			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, testInvestment, result)
		})
	})

	t.Run("FindOpportunitiesByUserId", func(t *testing.T) {
		mt.Run("database error", func(mt *mtest.T) {
			mt.AddMockResponses(
				mtest.CreateCommandErrorResponse(mtest.CommandError{
					Code:    11000,
					Message: "Database error",
				}),
			)

			repo := mongodb.NewInvestmentRepository(mt.DB)
			result, err := repo.FindOpportunitiesByUserId(context.Background(), testUserID)
			assert.Error(t, err)
			assert.Nil(t, result)
		})
		mt.Run("not found", func(mt *mtest.T) {
			mt.AddMockResponses(mtest.CreateCursorResponse(0, "foo.bar", mtest.FirstBatch))
			repo := mongodb.NewInvestmentRepository(mt.DB)
			result, err := repo.FindOpportunitiesByUserId(context.Background(), testUserID)
			assert.NoError(t, err)
			assert.Nil(t, result)
		})
		mt.Run("success", func(mt *mtest.T) {
			mt.AddMockResponses(
				mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, testOpportunityDocs...),
				mtest.CreateCursorResponse(0, "foo.bar", mtest.NextBatch),
			)
			repo := mongodb.NewInvestmentRepository(mt.Client.Database("testdb"))
			result, err := repo.FindOpportunitiesByUserId(context.Background(), testUserID)
			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.Len(t, result, len(testOpportunityDocs))
			// Validate each investment
			for i, opportunity := range result {
				assert.Equal(t, testOpportunities[i], opportunity)
			}
		})
	})

	t.Run("FindInvestmentsByUserId", func(t *testing.T) {
		mt.Run("database error", func(mt *mtest.T) {
			mt.AddMockResponses(
				mtest.CreateCommandErrorResponse(mtest.CommandError{
					Code:    11000,
					Message: "Database error",
				}),
			)

			repo := mongodb.NewInvestmentRepository(mt.DB)
			result, err := repo.FindInvestmentsByUserId(context.Background(), testUserID)
			assert.Error(t, err)
			assert.Nil(t, result)
		})
		mt.Run("not found", func(mt *mtest.T) {
			mt.AddMockResponses(mtest.CreateCursorResponse(0, "foo.bar", mtest.FirstBatch))
			repo := mongodb.NewInvestmentRepository(mt.DB)
			result, err := repo.FindInvestmentsByUserId(context.Background(), testUserID)
			assert.NoError(t, err)
			assert.Nil(t, result)
		})
		mt.Run("success", func(mt *mtest.T) {
			mt.AddMockResponses(
				mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, testInvestmentDocs...),
				mtest.CreateCursorResponse(0, "foo.bar", mtest.NextBatch),
			)
			repo := mongodb.NewInvestmentRepository(mt.Client.Database("testdb"))
			result, err := repo.FindInvestmentsByUserId(context.Background(), testUserID)
			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.Len(t, result, len(testInvestmentDocs))
			// Validate each investment
			for i, investment := range result {
				assert.Equal(t, testInvestments[i], investment)
			}
		})
	})
}
