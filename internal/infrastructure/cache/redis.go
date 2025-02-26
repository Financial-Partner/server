package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/Financial-Partner/server/internal/config"
)

type NewRedisClientFunc func(opts *redis.Options) redis.UniversalClient

type Client struct {
	redisClient redis.UniversalClient
}

var newRedisClient NewRedisClientFunc = func(opts *redis.Options) redis.UniversalClient {
	return redis.NewClient(opts)
}

func NewWithRedisClient(redisClient redis.UniversalClient) *Client {
	return &Client{redisClient: redisClient}
}

func NewClient(cfg *config.Config, newRedisClientFuncs ...NewRedisClientFunc) (*Client, error) {
	newRedisClientFunc := newRedisClient
	if len(newRedisClientFuncs) > 0 {
		newRedisClientFunc = newRedisClientFuncs[0]
	}

	redisClient := newRedisClientFunc(&redis.Options{
		Addr:     cfg.Redis.Host,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := redisClient.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return NewWithRedisClient(redisClient), nil
}

func (c *Client) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.redisClient.Set(ctx, key, data, expiration).Err()
}

func (c *Client) Get(ctx context.Context, key string, dest interface{}) error {
	data, err := c.redisClient.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dest)
}

func (c *Client) Delete(ctx context.Context, key string) error {
	return c.redisClient.Del(ctx, key).Err()
}
