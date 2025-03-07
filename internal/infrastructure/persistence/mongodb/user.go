package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/Financial-Partner/server/internal/domain/user"
	"github.com/Financial-Partner/server/internal/entities"
)

type MongoClient interface {
	Collection(name string, opts ...*options.CollectionOptions) *mongo.Collection
}

type MongoUserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(db MongoClient) user.Repository {
	return &MongoUserRepository{
		collection: db.Collection("users"),
	}
}

func (r *MongoUserRepository) FindByEmail(ctx context.Context, email string) (*entities.User, error) {
	var entity entities.User
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&entity)
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *MongoUserRepository) Create(ctx context.Context, entity *entities.User) (*entities.User, error) {
	_, err := r.collection.InsertOne(ctx, entity)
	if err != nil {
		return nil, err
	}
	return entity, nil
}

func (r *MongoUserRepository) Update(ctx context.Context, entity *entities.User) error {
	entity.UpdatedAt = time.Now()
	_, err := r.collection.ReplaceOne(ctx, bson.M{"_id": entity.ID}, entity)
	return err
}
