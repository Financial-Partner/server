package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/Financial-Partner/server/internal/config"
)

//go:generate mockgen -source=mongo.go -destination=mongo_mock.go -package=database

type MongoClient interface {
	Ping(ctx context.Context, rp *readpref.ReadPref) error
	Database(name string, opts ...*options.DatabaseOptions) *mongo.Database
	Disconnect(ctx context.Context) error
}

type ConnectFunc func(ctx context.Context, opts ...*options.ClientOptions) (MongoClient, error)

var mongoConnect ConnectFunc = func(ctx context.Context, opts ...*options.ClientOptions) (MongoClient, error) {
	return mongo.Connect(ctx, opts...)
}

type Client struct {
	client   MongoClient
	database *mongo.Database
}

func NewClient(cfg *config.Config, connectFuncs ...ConnectFunc) (*Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	connectFunc := mongoConnect
	if len(connectFuncs) > 0 {
		connectFunc = connectFuncs[0]
	}

	client, err := connectFunc(ctx, options.Client().ApplyURI(cfg.MongoDB.URI))
	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	database := client.Database(cfg.MongoDB.Database)
	return &Client{
		client:   client,
		database: database,
	}, nil
}

func (m *Client) Collection(name string, opts ...*options.CollectionOptions) *mongo.Collection {
	return m.database.Collection(name)
}

func (m *Client) Close(ctx context.Context) error {
	return m.client.Disconnect(ctx)
}
