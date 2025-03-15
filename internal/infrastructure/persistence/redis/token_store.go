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

func (s *TokenStore) SaveRefreshToken(ctx context.Context, id, refreshToken string, expiry time.Time) error {
	key := "refresh_token:" + refreshToken

	ttl := time.Until(expiry)

	return s.client.Set(ctx, key, id, ttl)
}

func (s *TokenStore) GetRefreshToken(ctx context.Context, refreshToken string) (string, error) {
	key := "refresh_token:" + refreshToken

	var id string
	err := s.client.Get(ctx, key, &id)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (s *TokenStore) DeleteRefreshToken(ctx context.Context, refreshToken string) error {
	key := "refresh_token:" + refreshToken

	return s.client.Delete(ctx, key)
}
