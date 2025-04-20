package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Financial-Partner/server/internal/entities"
)

const (
	investmentCacheKey  = "user:%s:investments"
	opportunityCacheKey = "user:%s:opportunities"
	investmentCacheTTL  = time.Hour * 24
)

type InvestmentStore struct {
	cacheClient RedisClient
}

func NewInvestmentStore(cacheClient RedisClient) *InvestmentStore {
	return &InvestmentStore{cacheClient: cacheClient}
}

func (s *InvestmentStore) SetOpportunities(ctx context.Context, userID string, opportunities []entities.Opportunity) error {
	data, err := json.Marshal(opportunities)
	if err != nil {
		return err
	}

	return s.cacheClient.Set(ctx, fmt.Sprintf(opportunityCacheKey, userID), data, investmentCacheTTL)
}

func (s *InvestmentStore) SetInvestments(ctx context.Context, userID string, investments []entities.Investment) error {
	data, err := json.Marshal(investments)
	if err != nil {
		return err
	}

	return s.cacheClient.Set(ctx, fmt.Sprintf(investmentCacheKey, userID), data, investmentCacheTTL)
}

func (s *InvestmentStore) DeleteInvestments(ctx context.Context, userID string) error {
	return s.cacheClient.Delete(ctx, fmt.Sprintf(investmentCacheKey, userID))
}

func (s *InvestmentStore) DeleteOpportunities(ctx context.Context, userID string) error {
	return s.cacheClient.Delete(ctx, fmt.Sprintf(opportunityCacheKey, userID))
}

func (s *InvestmentStore) GetOpportunities(ctx context.Context, userID string) ([]entities.Opportunity, error) {
	var opportunities []entities.Opportunity
	err := s.cacheClient.Get(ctx, fmt.Sprintf(opportunityCacheKey, userID), &opportunities)
	if err != nil {
		return nil, err
	}
	return opportunities, nil
}

func (s *InvestmentStore) GetInvestments(ctx context.Context, userID string) ([]entities.Investment, error) {
	var investments []entities.Investment
	err := s.cacheClient.Get(ctx, fmt.Sprintf(investmentCacheKey, userID), &investments)
	if err != nil {
		return nil, err
	}
	return investments, nil
}
