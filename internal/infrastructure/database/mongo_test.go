package database_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Financial-Partner/server/internal/config"
	"github.com/Financial-Partner/server/internal/infrastructure/database"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/mock/gomock"
)

func TestNewClient(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{}

	t.Run("success", func(t *testing.T) {
		mockClient := database.NewMockMongoClient(ctrl)
		mockClient.EXPECT().Ping(gomock.Any(), nil).Return(nil)
		mockClient.EXPECT().Database(gomock.Any(), gomock.Any()).Return(nil)

		connectFunc := func(ctx context.Context, opts ...*options.ClientOptions) (database.MongoClient, error) {
			return mockClient, nil
		}

		client, err := database.NewClient(cfg, connectFunc)
		require.NoError(t, err)
		require.NotNil(t, client)
	})

	t.Run("ping error", func(t *testing.T) {
		mockClient := database.NewMockMongoClient(ctrl)
		mockClient.EXPECT().Ping(gomock.Any(), nil).Return(errors.New("ping error"))

		connectFunc := func(ctx context.Context, opts ...*options.ClientOptions) (database.MongoClient, error) {
			return mockClient, nil
		}

		client, err := database.NewClient(cfg, connectFunc)
		require.Error(t, err)
		require.Nil(t, client)
	})

	t.Run("connect error", func(t *testing.T) {
		connectFunc := func(ctx context.Context, opts ...*options.ClientOptions) (database.MongoClient, error) {
			return nil, errors.New("connect error")
		}

		client, err := database.NewClient(cfg, connectFunc)
		require.Error(t, err)
		require.Nil(t, client)
	})

	t.Run("fail to create client", func(t *testing.T) {
		client, err := database.NewClient(cfg)
		require.Error(t, err)
		require.Nil(t, client)
	})
}

func TestClientCollection(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	cfg := &config.Config{}

	mt.Run("success", func(mt *mtest.T) {
		mockClient := database.NewMockMongoClient(ctrl)
		mockClient.EXPECT().Ping(gomock.Any(), nil).Return(nil)
		mockDB := mt.Client.Database(mt.DB.Name(), nil)
		mockClient.EXPECT().Database(gomock.Any(), gomock.Any()).Return(mockDB)

		connectFunc := func(ctx context.Context, opts ...*options.ClientOptions) (database.MongoClient, error) {
			return mockClient, nil
		}

		client, err := database.NewClient(cfg, connectFunc)
		require.NoError(t, err)
		require.NotNil(t, client)

		collection := client.Collection("test")
		require.NotNil(t, collection)
	})
}

func TestClientClose(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	cfg := &config.Config{}

	mt.Run("success", func(mt *mtest.T) {
		mockClient := database.NewMockMongoClient(ctrl)
		mockClient.EXPECT().Ping(gomock.Any(), nil).Return(nil)
		mockDB := mt.Client.Database(mt.DB.Name(), nil)
		mockClient.EXPECT().Database(gomock.Any(), gomock.Any()).Return(mockDB)
		mockClient.EXPECT().Disconnect(gomock.Any()).Return(nil)

		connectFunc := func(ctx context.Context, opts ...*options.ClientOptions) (database.MongoClient, error) {
			return mockClient, nil
		}

		client, err := database.NewClient(cfg, connectFunc)
		require.NoError(t, err)
		require.NotNil(t, client)

		err = client.Close(context.Background())
		require.NoError(t, err)
	})
}
