package entities

import (
	"time"
)

type Opportunity struct {
	ID          string    `bson:"_id,omitempty" json:"id"`
	Title       string    `bson:"title" json:"title"`
	Description string    `bson:"description" json:"description"`
	Tags        []string  `bson:"tags" json:"tags"`
	IsIncrease  bool      `bson:"is_increase" json:"is_increase"`
	Variation   int64     `bson:"variation" json:"variation"`
	Duration    string    `bson:"duration" json:"duration"`
	MinAmount   int64     `bson:"min_amount" json:"min_amount"`
	CreatedAt   time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at" json:"updated_at"`
}

type Investment struct {
	ID            string    `bson:"_id,omitempty" json:"id"`
	UserID        string    `bson:"user_id" json:"user_id"`
	OpportunityID string    `bson:"opportunity_id" json:"opportunity_id"`
	Amount        int64     `bson:"amount" json:"amount"`
	CreatedAt     time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time `bson:"updated_at" json:"updated_at"`
}
