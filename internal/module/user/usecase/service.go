package user_usecase

import (
	"context"
	"time"

	"github.com/Financial-Partner/server/internal/entities"
	"github.com/Financial-Partner/server/internal/infrastructure/logger"
	user_repository "github.com/Financial-Partner/server/internal/module/user/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service struct {
	repo  user_repository.Repository
	store user_repository.UserStore
	log   logger.Logger
}

func NewService(repo user_repository.Repository, store user_repository.UserStore, log logger.Logger) *Service {
	return &Service{
		repo:  repo,
		store: store,
		log:   log,
	}
}

func (s *Service) GetUser(ctx context.Context, email string) (*entities.User, error) {
	entity, err := s.store.Get(ctx, email)
	if err == nil {
		return entity, nil
	}

	entity, err = s.repo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	s.setUserToStore(ctx, entity)

	return entity, nil
}

func (s *Service) GetOrCreateUser(ctx context.Context, email, name string) (*entities.User, error) {
	logger := s.log.WithField("email", email)

	entity, err := s.GetUser(ctx, email)
	if err == nil {
		return entity, nil
	}

	logger.Infof("Creating new user")
	newEntity := &entities.User{
		ID:    primitive.NewObjectID(),
		Email: email,
		Name:  name,
		Wallet: entities.Wallet{
			Diamonds: 0,
			Savings:  0,
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	entity, err = s.repo.Create(ctx, newEntity)
	if err != nil {
		logger.WithError(err).Errorf("Failed to create new user")
		return nil, err
	}

	logger.Infof("New user created successfully")

	s.setUserToStore(ctx, entity)

	return entity, nil
}

func (s *Service) UpdateUserName(ctx context.Context, email, name string) (*entities.User, error) {
	return nil, nil
}

func (s *Service) setUserToStore(ctx context.Context, entity *entities.User) {
	err := s.store.Set(ctx, entity)
	if err != nil {
		s.log.WithError(err).Errorf("Failed to create user in store")
	}
}
