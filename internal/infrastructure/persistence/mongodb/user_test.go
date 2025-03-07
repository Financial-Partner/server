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

func TestMongoUserRepository(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	testUserID := primitive.NewObjectID()
	testUser := &entities.User{
		ID:        testUserID,
		Email:     "test@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	testUserBSON, err := bson.Marshal(testUser)
	require.NoError(t, err)
	var testUserDoc bson.D
	err = bson.Unmarshal(testUserBSON, &testUserDoc)
	require.NoError(t, err)
	t.Run("FindByEmail", func(t *testing.T) {
		mt.Run("success", func(mt *mtest.T) {
			mt.AddMockResponses(mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, testUserDoc))
			repo := mongodb.NewUserRepository(mt.Client.Database("testdb"))
			result, err := repo.FindByEmail(context.Background(), "test@example.com")
			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, testUser.ID, result.ID)
			assert.Equal(t, testUser.Email, result.Email)
		})
		mt.Run("not found", func(mt *mtest.T) {
			mt.AddMockResponses(mtest.CreateCursorResponse(0, "foo.bar", mtest.FirstBatch))
			repo := mongodb.NewUserRepository(mt.DB)
			result, err := repo.FindByEmail(context.Background(), "nonexistent@example.com")
			assert.Error(t, err)
			assert.Nil(t, result)
		})
		mt.Run("database error", func(mt *mtest.T) {
			mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
				Code:    11000,
				Message: "database error",
			}))
			repo := mongodb.NewUserRepository(mt.DB)
			result, err := repo.FindByEmail(context.Background(), "test@example.com")
			assert.Error(t, err)
			assert.Nil(t, result)
		})
	})
	t.Run("Create", func(t *testing.T) {
		mt.Run("success", func(mt *mtest.T) {
			mt.AddMockResponses(mtest.CreateSuccessResponse())
			repo := mongodb.NewUserRepository(mt.DB)
			result, err := repo.Create(context.Background(), testUser)
			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, testUser.ID, result.ID)
			assert.Equal(t, testUser.Email, result.Email)
		})
		mt.Run("error", func(mt *mtest.T) {
			mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
				Code:    11000,
				Message: "duplicate key error",
			}))
			repo := mongodb.NewUserRepository(mt.DB)
			result, err := repo.Create(context.Background(), testUser)
			assert.Error(t, err)
			assert.Nil(t, result)
		})
	})
	t.Run("Update", func(t *testing.T) {
		mt.Run("success", func(mt *mtest.T) {
			mt.AddMockResponses(mtest.CreateSuccessResponse())
			repo := mongodb.NewUserRepository(mt.DB)
			err := repo.Update(context.Background(), testUser)
			assert.NoError(t, err)
		})
		mt.Run("error", func(mt *mtest.T) {
			mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
				Code:    11000,
				Message: "update error",
			}))
			repo := mongodb.NewUserRepository(mt.DB)
			err := repo.Update(context.Background(), testUser)
			assert.Error(t, err)
		})
	})
}
