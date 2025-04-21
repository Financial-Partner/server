package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GoalSuggestion struct {
	SuggestedAmount int64  `bson:"suggested_amount" json:"suggested_amount"`
	Period          int    `bson:"period" json:"period"`
	Message         string `bson:"message" json:"message"`
}

type Goal struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID        primitive.ObjectID `bson:"user_id" json:"user_id"`
	TargetAmount  int64              `bson:"target_amount" json:"target_amount"`
	CurrentAmount int64              `bson:"current_amount" json:"current_amount"`
	Period        int                `bson:"period" json:"period"`
	Status        string             `bson:"status" json:"status"` // "active", "completed", "failed"
	CreatedAt     time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time          `bson:"updated_at" json:"updated_at"`
}

type GoalMilestone struct {
	Title         string     `bson:"title" json:"title"`
	TargetPercent int        `bson:"target_percent" json:"target_percent"`
	Reward        string     `bson:"reward" json:"reward"`
	IsCompleted   bool       `bson:"is_completed" json:"is_completed"`
	CompletedAt   *time.Time `bson:"completed_at" json:"completed_at"`
}
