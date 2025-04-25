package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Financial-Partner/server/internal/entities"
	"github.com/redis/go-redis/v9"
)

const (
	transactionCacheKey = "user:%s:transactions"
	transactionCacheTTL = time.Hour * 24
)

type TransactionStore struct {
	cacheClient RedisClient
}

func NewTransactionStore(cacheClient RedisClient) *TransactionStore {
	return &TransactionStore{cacheClient: cacheClient}
}

func (s *TransactionStore) GetByUserId(ctx context.Context, userID string) ([]entities.Transaction, error) {
	var transactions []entities.Transaction
	err := s.cacheClient.Get(ctx, fmt.Sprintf(transactionCacheKey, userID), &transactions)
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

func (s *TransactionStore) AddByUserId(ctx context.Context, userID string, transaction *entities.Transaction) error {
	transactions, err := s.GetByUserId(ctx, userID)
	if err != nil && !errors.Is(err, redis.Nil) {
		return err
	}

	transactions = append(transactions, *transaction)
	err = s.cacheClient.Set(ctx, fmt.Sprintf(transactionCacheKey, userID), transactions, transactionCacheTTL)
	if err != nil {
		return err
	}
	return nil
}

func (s *TransactionStore) SetMultipleByUserId(ctx context.Context, userID string, transactions []entities.Transaction) error {
	existingTransactions, err := s.GetByUserId(ctx, userID)
	if err != nil && !errors.Is(err, redis.Nil) {
		return err
	}

	transactions = append(existingTransactions, transactions...)
	err = s.cacheClient.Set(ctx, fmt.Sprintf(transactionCacheKey, userID), transactions, transactionCacheTTL)
	if err != nil {
		return err
	}
	return nil
}

func (s *TransactionStore) DeleteByUserId(ctx context.Context, userID string) error {
	return s.cacheClient.Delete(ctx, fmt.Sprintf(transactionCacheKey, userID))
}
