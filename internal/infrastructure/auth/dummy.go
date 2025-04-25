package auth

import (
	"errors"

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

	return &Claims{
		ID:    "680bc51354f0a75ad544ff79",
		Email: "bypass@example.com",
	}, nil
}
