package redis

import (
	"context"
	"time"
)

//go:generate mockgen -source=redis.go -destination=mocks/redis_mock.go -package=mocks

type RedisClient interface {
	Get(ctx context.Context, key string, dest interface{}) error
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Delete(ctx context.Context, key string) error
}
