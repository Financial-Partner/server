package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/Financial-Partner/server/internal/entities"
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

func (s *TransactionStore) SetByUserId(ctx context.Context, userID string, transactions []entities.Transaction) error {
	return s.cacheClient.Set(ctx, fmt.Sprintf(transactionCacheKey, userID), transactions, transactionCacheTTL)
}

func (s *TransactionStore) DeleteByUserId(ctx context.Context, userID string) error {
	return s.cacheClient.Delete(ctx, fmt.Sprintf(transactionCacheKey, userID))
}
