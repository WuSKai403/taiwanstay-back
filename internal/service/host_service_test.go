package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/taiwanstay/taiwanstay-back/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MockHostRepository is a mock implementation of HostRepository
type MockHostRepository struct {
	mock.Mock
}

func (m *MockHostRepository) Create(ctx context.Context, host *domain.Host) error {
	args := m.Called(ctx, host)
	return args.Error(0)
}

func (m *MockHostRepository) GetByID(ctx context.Context, id string) (*domain.Host, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Host), args.Error(1)
}

func (m *MockHostRepository) GetByUserID(ctx context.Context, userID string) (*domain.Host, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Host), args.Error(1)
}

func (m *MockHostRepository) Update(ctx context.Context, id string, host *domain.Host) error {
	args := m.Called(ctx, id, host)
	return args.Error(0)
}

func TestCreateHost(t *testing.T) {
	mockRepo := new(MockHostRepository)
	service := NewHostService(mockRepo)

	ctx := context.Background()
	userID := primitive.NewObjectID()
	host := &domain.Host{
		UserID: userID,
		Name:   "Test Farm",
	}

	// Expect Create to be called
	mockRepo.On("Create", ctx, host).Return(nil)

	createdHost, err := service.CreateHost(ctx, host)

	assert.NoError(t, err)
	assert.NotNil(t, createdHost)
	assert.Equal(t, "Test Farm", createdHost.Name)
	assert.NotEmpty(t, createdHost.Slug) // Slug should be generated
	assert.Contains(t, createdHost.Slug, "test-farm")
	assert.Equal(t, domain.HostStatusPending, createdHost.Status) // Default status

	mockRepo.AssertExpectations(t)
}

func TestGetHostByUserID(t *testing.T) {
	mockRepo := new(MockHostRepository)
	service := NewHostService(mockRepo)

	ctx := context.Background()
	userID := primitive.NewObjectID().Hex()
	expectedHost := &domain.Host{
		Name: "Existing Host",
	}

	mockRepo.On("GetByUserID", ctx, userID).Return(expectedHost, nil)

	host, err := service.GetHostByUserID(ctx, userID)

	assert.NoError(t, err)
	assert.Equal(t, expectedHost, host)
	mockRepo.AssertExpectations(t)
}
