package redis

import (
	"context"
	"time"
)

type TokenStore struct {
	client RedisClient
}

func NewTokenStore(client RedisClient) *TokenStore {
	return &TokenStore{client: client}
}

func (s *TokenStore) SaveRefreshToken(ctx context.Context, email, refreshToken string, expiry time.Time) error {
	key := "refresh_token:" + refreshToken

	ttl := time.Until(expiry)

	return s.client.Set(ctx, key, email, ttl)
}

func (s *TokenStore) GetRefreshToken(ctx context.Context, refreshToken string) (string, error) {
	key := "refresh_token:" + refreshToken

	var email string
	err := s.client.Get(ctx, key, &email)
	if err != nil {
		return "", err
	}

	return email, nil
}

func (s *TokenStore) DeleteRefreshToken(ctx context.Context, refreshToken string) error {
	key := "refresh_token:" + refreshToken

	return s.client.Delete(ctx, key)
}
