package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/Financial-Partner/server/internal/entities"
)

type MongoInvestmentResporitory struct {
	collection *mongo.Collection
}

func NewInvestmentRepository(db MongoClient) *MongoInvestmentResporitory {
	return &MongoInvestmentResporitory{
		collection: db.Collection("investments"),
	}
}

func (r *MongoInvestmentResporitory) CreateInvestment(ctx context.Context, entity *entities.Investment) (*entities.Investment, error) {
	_, err := r.collection.InsertOne(ctx, entity)
	if err != nil {
		return nil, err
	}
	return entity, nil
}

func (r *MongoInvestmentResporitory) CreateOpportunity(ctx context.Context, entity *entities.Opportunity) (*entities.Opportunity, error) {
	_, err := r.collection.InsertOne(ctx, entity)
	if err != nil {
		return nil, err
	}
	return entity, nil
}

func (r *MongoInvestmentResporitory) FindOpportunitiesByUserId(ctx context.Context, userID string) ([]entities.Opportunity, error) {
	var opportunities []entities.Opportunity
	cursor, err := r.collection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &opportunities); err != nil {
		return nil, err
	}

	return opportunities, nil
}

func (r *MongoInvestmentResporitory) FindInvestmentsByUserId(ctx context.Context, userID string) ([]entities.Investment, error) {
	var investments []entities.Investment
	cursor, err := r.collection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &investments); err != nil {
		return nil, err
	}

	return investments, nil
}
