package mongodb

import (
	"context"
	"fmt"

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

	if len(transactions) == 0 {
		return nil, fmt.Errorf("no transactions found for user ID: %s", userID)
	}

	for cursor.Next(ctx) {
		var transaction entities.Transaction
		if err := cursor.Decode(&transaction); err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}
