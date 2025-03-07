package user_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Financial-Partner/server/internal/domain/user"
	"github.com/Financial-Partner/server/internal/domain/user/mocks"
	"github.com/Financial-Partner/server/internal/entities"
	"github.com/Financial-Partner/server/internal/infrastructure/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestService(t *testing.T) {
	t.Run("GetUserFromStore", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockRepository(ctrl)
		mockStore := mocks.NewMockUserStore(ctrl)
		mockLogger := logger.NewNopLogger()

		svc := user.NewService(mockRepo, mockStore, mockLogger)
		ctx := context.Background()
		email := "test@example.com"

		expectedUser := &entities.User{
			Email: email,
			Name:  "Test User",
		}

		mockStore.EXPECT().Get(ctx, email).Return(expectedUser, nil)

		result, err := svc.GetUser(ctx, email)
		require.NoError(t, err)
		assert.Equal(t, expectedUser, result)
	})

	t.Run("GetUserFromRepo", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockRepository(ctrl)
		mockStore := mocks.NewMockUserStore(ctrl)
		mockLogger := logger.NewNopLogger()

		svc := user.NewService(mockRepo, mockStore, mockLogger)
		ctx := context.Background()
		email := "test@example.com"

		storeErr := errors.New("not found in store")
		expectedUser := &entities.User{
			Email: email,
			Name:  "Test User",
		}

		mockStore.EXPECT().Get(ctx, email).Return(nil, storeErr)
		mockRepo.EXPECT().FindByEmail(ctx, email).Return(expectedUser, nil)
		mockStore.EXPECT().Set(ctx, expectedUser).Return(nil)

		result, err := svc.GetUser(ctx, email)
		require.NoError(t, err)
		assert.Equal(t, expectedUser, result)
	})

	t.Run("GetUserFromRepoSuccessButSetUserToStoreFailure", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockRepository(ctrl)
		mockStore := mocks.NewMockUserStore(ctrl)
		mockLogger := logger.NewNopLogger()

		svc := user.NewService(mockRepo, mockStore, mockLogger)
		ctx := context.Background()
		email := "test@example.com"

		storeErr := errors.New("not found in store")
		expectedUser := &entities.User{
			Email: email,
			Name:  "Test User",
		}

		mockStore.EXPECT().Get(ctx, email).Return(nil, storeErr)
		mockRepo.EXPECT().FindByEmail(ctx, email).Return(expectedUser, nil)
		mockStore.EXPECT().Set(ctx, expectedUser).Return(errors.New("failed to set user to store"))

		result, err := svc.GetUser(ctx, email)
		require.NoError(t, err)
		assert.Equal(t, expectedUser, result)
	})

	t.Run("GetUserErrorFromRepo", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockRepository(ctrl)
		mockStore := mocks.NewMockUserStore(ctrl)
		mockLogger := logger.NewNopLogger()

		svc := user.NewService(mockRepo, mockStore, mockLogger)
		ctx := context.Background()
		email := "test@example.com"

		storeErr := errors.New("store: not found")
		repoErr := errors.New("repo: not found")
		mockStore.EXPECT().Get(ctx, email).Return(nil, storeErr)
		mockRepo.EXPECT().FindByEmail(ctx, email).Return(nil, repoErr)

		result, err := svc.GetUser(ctx, email)
		require.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, repoErr, err)
	})

	t.Run("GetOrCreateUserExistingUser", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockRepository(ctrl)
		mockStore := mocks.NewMockUserStore(ctrl)
		mockLogger := logger.NewNopLogger()

		svc := user.NewService(mockRepo, mockStore, mockLogger)
		ctx := context.Background()
		email := "test@example.com"
		name := "Existing User"

		existingUser := &entities.User{
			Email: email,
			Name:  name,
		}

		mockStore.EXPECT().Get(ctx, email).Return(existingUser, nil)

		result, err := svc.GetOrCreateUser(ctx, email, name)
		require.NoError(t, err)
		assert.Equal(t, existingUser, result)
	})

	t.Run("GetOrCreateUserNewUserSuccess", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockRepository(ctrl)
		mockStore := mocks.NewMockUserStore(ctrl)
		mockLogger := logger.NewNopLogger()

		svc := user.NewService(mockRepo, mockStore, mockLogger)
		ctx := context.Background()
		email := "new@example.com"
		name := "New User"

		storeErr := errors.New("store: not found")
		repoErr := errors.New("repo: not found")
		mockStore.EXPECT().Get(ctx, email).Return(nil, storeErr)
		mockRepo.EXPECT().FindByEmail(ctx, email).Return(nil, repoErr)

		createdUser := &entities.User{
			Email: email,
			Name:  name,
		}
		mockRepo.EXPECT().Create(ctx, gomock.Any()).DoAndReturn(
			func(ctx context.Context, entity *entities.User) (*entities.User, error) {
				assert.Equal(t, email, entity.Email)
				assert.Equal(t, name, entity.Name)
				return createdUser, nil
			},
		)
		mockStore.EXPECT().Set(ctx, createdUser).Return(nil)

		result, err := svc.GetOrCreateUser(ctx, email, name)
		require.NoError(t, err)
		assert.Equal(t, createdUser, result)
	})

	t.Run("GetOrCreateUserNewUserSuccessButSetUserToStoreFailure", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockRepository(ctrl)
		mockStore := mocks.NewMockUserStore(ctrl)
		mockLogger := logger.NewNopLogger()

		svc := user.NewService(mockRepo, mockStore, mockLogger)
		ctx := context.Background()
		email := "new@example.com"
		name := "New User"

		storeErr := errors.New("store: not found")
		repoErr := errors.New("repo: not found")
		mockStore.EXPECT().Get(ctx, email).Return(nil, storeErr)
		mockRepo.EXPECT().FindByEmail(ctx, email).Return(nil, repoErr)

		createdUser := &entities.User{
			Email: email,
			Name:  name,
		}
		mockRepo.EXPECT().Create(ctx, gomock.Any()).DoAndReturn(
			func(ctx context.Context, entity *entities.User) (*entities.User, error) {
				assert.Equal(t, email, entity.Email)
				assert.Equal(t, name, entity.Name)
				return createdUser, nil
			},
		)
		mockStore.EXPECT().Set(ctx, createdUser).Return(errors.New("failed to set user to store"))

		result, err := svc.GetOrCreateUser(ctx, email, name)
		require.NoError(t, err)
		assert.Equal(t, createdUser, result)
	})

	t.Run("GetOrCreateUserNewUserFailure", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockRepository(ctrl)
		mockStore := mocks.NewMockUserStore(ctrl)
		mockLogger := logger.NewNopLogger()

		svc := user.NewService(mockRepo, mockStore, mockLogger)
		ctx := context.Background()
		email := "fail@example.com"
		name := "Fail User"

		storeErr := errors.New("store: not found")
		repoErr := errors.New("repo: not found")
		mockStore.EXPECT().Get(ctx, email).Return(nil, storeErr)
		mockRepo.EXPECT().FindByEmail(ctx, email).Return(nil, repoErr)

		creationErr := errors.New("creation failed")
		mockRepo.EXPECT().Create(ctx, gomock.Any()).Return(nil, creationErr)

		result, err := svc.GetOrCreateUser(ctx, email, name)
		require.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, creationErr, err)
	})
}
