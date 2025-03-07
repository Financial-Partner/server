package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/Financial-Partner/server/internal/entities"
)

const (
	userCacheKey = "user:%s"
	userCacheTTL = time.Hour * 24
)

type UserStore struct {
	cacheClient RedisClient
}

func NewUserStore(cacheClient RedisClient) *UserStore {
	return &UserStore{cacheClient: cacheClient}
}

func (s *UserStore) Get(ctx context.Context, email string) (*entities.User, error) {
	var user entities.User
	err := s.cacheClient.Get(ctx, fmt.Sprintf(userCacheKey, email), &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserStore) Set(ctx context.Context, user *entities.User) error {
	return s.cacheClient.Set(ctx, fmt.Sprintf(userCacheKey, user.Email), user, userCacheTTL)
}

func (s *UserStore) Delete(ctx context.Context, email string) error {
	return s.cacheClient.Delete(ctx, email)
}
