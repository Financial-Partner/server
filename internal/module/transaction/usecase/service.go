package transaction_usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/Financial-Partner/server/internal/entities"
	"github.com/Financial-Partner/server/internal/infrastructure/logger"
	"github.com/Financial-Partner/server/internal/interfaces/http/dto"
	transaction_repository "github.com/Financial-Partner/server/internal/module/transaction/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service struct {
	repo  transaction_repository.Repository
	store transaction_repository.TransactionStore
	log   logger.Logger
}

func NewService(repo transaction_repository.Repository, store transaction_repository.TransactionStore, log logger.Logger) *Service {
	return &Service{
		repo:  repo,
		store: store,
		log:   log,
	}
}

func (s *Service) CreateTransaction(ctx context.Context, userID string, req *dto.CreateTransactionRequest) (*entities.Transaction, error) {
	transactionDate, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, fmt.Errorf("invalid date format: %w", err)
	}

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	// Convert DTO to Entity
	transaction := &entities.Transaction{
		UserID:      objectID,
		Amount:      req.Amount,
		Category:    req.Category,
		Type:        req.Type,
		Date:        transactionDate.UTC(),
		Description: req.Description,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	// Call the repository to persist the transaction
	createdTransaction, err := s.repo.Create(ctx, transaction)
	if err != nil {
		cacheErr := s.store.SetByUserId(ctx, userID, []entities.Transaction{*createdTransaction})
		if cacheErr != nil {
			s.log.Warnf("Failed to cache transaction for userID %s: %v", userID, cacheErr)
		}
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	// Return the created transaction
	return createdTransaction, nil
}

func (s *Service) GetTransactions(ctx context.Context, userID string) ([]entities.Transaction, error) {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}
	// Check if transactions are cached
	cachedTransactions, err := s.store.GetByUserId(ctx, userID)
	if err != nil && err.Error() != "redis: nil" {
		return nil, fmt.Errorf("failed to get cached transactions: %w", err)
	}
	if cachedTransactions != nil {
		return cachedTransactions, nil
	}
	// If not cached, fetch from the repository
	transactions, err := s.repo.FindByUserId(ctx, objectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions: %w", err)
	}

	// Cache the fetched transactions
	cacheErr := s.store.SetByUserId(ctx, userID, transactions)
	if cacheErr != nil {
		s.log.Warnf("Failed to cache transaction for userID %s: %v", userID, cacheErr)
	}
	return transactions, nil
}
