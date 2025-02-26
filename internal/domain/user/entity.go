package user

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WalletEntity struct {
	Diamonds int64 `bson:"diamonds" json:"diamonds"`
	Savings  int64 `bson:"savings" json:"savings"`
}

type UserEntity struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Email     string             `bson:"email" json:"email"`
	Name      string             `bson:"name" json:"name"`
	Wallet    WalletEntity       `bson:"wallet" json:"wallet"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}
