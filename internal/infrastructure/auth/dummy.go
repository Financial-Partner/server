package auth

import (
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Financial-Partner/server/internal/config"
)

type DummyJWTValidator struct {
	cfg *config.Config
}

func NewDummyJWTValidator(cfg *config.Config) *DummyJWTValidator {
	return &DummyJWTValidator{
		cfg: cfg,
	}
}

func (v *DummyJWTValidator) ValidateToken(tokenString string) (*Claims, error) {
	if tokenString != v.cfg.Firebase.BypassToken {
		return nil, errors.New("invalid token")
	}

	dummyObjectID, err := primitive.ObjectIDFromHex("680b4fc122fc6fd9212d78f9")
	if err != nil {
		return nil, errors.New("failed to create dummy ObjectID")
	}

	return &Claims{
		ID:    dummyObjectID.Hex(),
		Email: "bypass@example.com",
	}, nil
}
