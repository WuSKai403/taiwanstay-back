package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/taiwanstay/taiwanstay-back/internal/domain"
	"github.com/taiwanstay/taiwanstay-back/pkg/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// mockUserRepository 是一個用於測試的 UserRepository mock
type mockUserRepository struct {
	mock.Mock
}

func (m *mockUserRepository) Create(ctx context.Context, user *domain.User) (string, error) {
	args := m.Called(ctx, user)
	return args.String(0), args.Error(1)
}

func (m *mockUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *mockUserRepository) GetAll(ctx context.Context) ([]*domain.User, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.User), args.Error(1)
}

func (m *mockUserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *mockUserRepository) Update(ctx context.Context, id string, payload bson.M) error {
	args := m.Called(ctx, id, payload)
	return args.Error(0)
}

func (m *mockUserRepository) Count(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockUserRepository) List(ctx context.Context, filter bson.M, limit, offset int64) ([]*domain.User, int64, error) {
	args := m.Called(ctx, filter, limit, offset)
	return args.Get(0).([]*domain.User), args.Get(1).(int64), args.Error(2)
}

func (m *mockUserRepository) UpdateStatus(ctx context.Context, id string, status domain.UserStatus) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func TestLoginUser(t *testing.T) {
	// 準備加密後的密碼
	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	// 準備一個測試用的使用者
	testUser := &domain.User{
		ID:       "some-user-id",
		Name:     "testuser",
		Email:    "test@example.com",
		Password: string(hashedPassword),
		Role:     domain.RoleUser,
	}

	t.Run("Successful Login", func(t *testing.T) {
		// 準備 mock
		mockRepo := new(mockUserRepository)
		mockRepo.On("GetByEmail", mock.Anything, testUser.Email).Return(testUser, nil)

		// 建立 service
		cfg := &config.Config{Server: config.ServerConfig{JWTSecret: "test-secret"}}
		userService := NewUserService(mockRepo, cfg)

		// 執行登入
		user, token, err := userService.LoginUser(context.Background(), testUser.Email, password)

		// 斷言
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.NotEmpty(t, token)
		assert.Equal(t, testUser.Email, user.Email)
		assert.Empty(t, user.Password) // 確保密碼已被清除
		mockRepo.AssertExpectations(t)
	})

	t.Run("Failed Login - Wrong Password", func(t *testing.T) {
		// 準備 mock
		mockRepo := new(mockUserRepository)
		mockRepo.On("GetByEmail", mock.Anything, testUser.Email).Return(testUser, nil)

		// 建立 service
		cfg := &config.Config{Server: config.ServerConfig{JWTSecret: "test-secret"}}
		userService := NewUserService(mockRepo, cfg)

		// 執行登入
		_, _, err := userService.LoginUser(context.Background(), testUser.Email, "wrongpassword")

		// 斷言
		assert.Error(t, err)
		assert.Equal(t, ErrInvalidCredentials, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Failed Login - User Not Found", func(t *testing.T) {
		// 準備 mock
		mockRepo := new(mockUserRepository)
		mockRepo.On("GetByEmail", mock.Anything, "notfound@example.com").Return(nil, mongo.ErrNoDocuments)

		// 建立 service
		cfg := &config.Config{Server: config.ServerConfig{JWTSecret: "test-secret"}}
		userService := NewUserService(mockRepo, cfg)

		// 執行登入
		_, _, err := userService.LoginUser(context.Background(), "notfound@example.com", "password123")

		// 斷言
		assert.Error(t, err)
		assert.Equal(t, ErrInvalidCredentials, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestRegisterUser(t *testing.T) {
	t.Run("Successful Registration", func(t *testing.T) {
		mockRepo := new(mockUserRepository)
		mockRepo.On("GetByEmail", mock.Anything, "new@example.com").Return(nil, mongo.ErrNoDocuments)
		mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.User")).Return("new-user-id", nil)

		cfg := &config.Config{Server: config.ServerConfig{JWTSecret: "test-secret"}}
		userService := NewUserService(mockRepo, cfg)
		user, err := userService.RegisterUser(context.Background(), "newuser", "new@example.com", "password123")

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, "new-user-id", user.ID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Failed Registration - Email Exists", func(t *testing.T) {
		existingUser := &domain.User{Email: "exists@example.com"}
		mockRepo := new(mockUserRepository)
		mockRepo.On("GetByEmail", mock.Anything, "exists@example.com").Return(existingUser, nil)

		cfg := &config.Config{Server: config.ServerConfig{JWTSecret: "test-secret"}}
		userService := NewUserService(mockRepo, cfg)
		_, err := userService.RegisterUser(context.Background(), "anotheruser", "exists@example.com", "password123")

		assert.Error(t, err)
		assert.Equal(t, ErrEmailAlreadyExists, err)
		mockRepo.AssertExpectations(t)
	})
}
