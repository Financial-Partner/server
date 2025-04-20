package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/Financial-Partner/server/internal/entities"
	transaction_repository "github.com/Financial-Partner/server/internal/module/transaction/repository"
)

type MongoTransactionRepository struct {
	collection *mongo.Collection
}

func NewTransactionRepository(db MongoClient) transaction_repository.Repository {
	return &MongoTransactionRepository{
		collection: db.Collection("transactions"),
	}
}

func (r *MongoTransactionRepository) Create(ctx context.Context, entity *entities.Transaction) (*entities.Transaction, error) {
	_, err := r.collection.InsertOne(ctx, entity)
	if err != nil {
		return nil, err
	}
	return entity, nil
}

func (r *MongoTransactionRepository) FindByUserId(ctx context.Context, userID string) ([]entities.Transaction, error) {
	var transactions []entities.Transaction
	cursor, err := r.collection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &transactions); err != nil {
		return nil, err
	}

	return transactions, nil
}
