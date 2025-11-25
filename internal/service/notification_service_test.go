package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/taiwanstay/taiwanstay-back/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Mocks
type MockNotificationRepository struct {
	mock.Mock
}

func (m *MockNotificationRepository) Create(ctx context.Context, notification *domain.Notification) error {
	args := m.Called(ctx, notification)
	return args.Error(0)
}

func (m *MockNotificationRepository) ListByUserID(ctx context.Context, userID string, limit, offset int64) ([]*domain.Notification, int64, error) {
	args := m.Called(ctx, userID, limit, offset)
	return args.Get(0).([]*domain.Notification), args.Get(1).(int64), args.Error(2)
}

func (m *MockNotificationRepository) MarkAsRead(ctx context.Context, id string, userID string) error {
	args := m.Called(ctx, id, userID)
	return args.Error(0)
}

func (m *MockNotificationRepository) MarkAllAsRead(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

type MockEmailSender struct {
	mock.Mock
}

func (m *MockEmailSender) Send(toEmail, toName, subject, htmlBody string) error {
	args := m.Called(toEmail, toName, subject, htmlBody)
	return args.Error(0)
}

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *domain.User) (string, error) {
	args := m.Called(ctx, user)
	return args.String(0), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, id string, payload bson.M) error {
	args := m.Called(ctx, id, payload)
	return args.Error(0)
}

func (m *MockUserRepository) GetAll(ctx context.Context) ([]*domain.User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*domain.User), args.Error(1)
}

func (m *MockUserRepository) Count(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockUserRepository) List(ctx context.Context, filter bson.M, limit, offset int64) ([]*domain.User, int64, error) {
	args := m.Called(ctx, filter, limit, offset)
	return args.Get(0).([]*domain.User), args.Get(1).(int64), args.Error(2)
}

func TestSendNotification(t *testing.T) {
	mockRepo := new(MockNotificationRepository)
	mockUserRepo := new(MockUserRepository)
	mockEmailSender := new(MockEmailSender)
	service := NewNotificationService(mockRepo, mockUserRepo, mockEmailSender)

	userID := primitive.NewObjectID().Hex()
	user := &domain.User{
		ID:    primitive.NewObjectID().Hex(),
		Email: "test@example.com",
		Name:  "Test User",
	}

	// Expectation: Create notification in DB
	mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(n *domain.Notification) bool {
		return n.Title == "Test Title"
	})).Return(nil)

	// Expectation: Get User for email
	mockUserRepo.On("GetByID", mock.Anything, userID).Return(user, nil)

	// Expectation: Send Email
	mockEmailSender.On("Send", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	err := service.SendNotification(context.Background(), userID, domain.NotificationTypeApplicationCreated, "Test Title", "Test Message", nil)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

func TestListNotifications(t *testing.T) {
	mockRepo := new(MockNotificationRepository)
	mockUserRepo := new(MockUserRepository)
	mockEmailSender := new(MockEmailSender)
	service := NewNotificationService(mockRepo, mockUserRepo, mockEmailSender)

	userID := primitive.NewObjectID().Hex()
	expectedNotifs := []*domain.Notification{{Title: "Test"}}

	mockRepo.On("ListByUserID", mock.Anything, userID, int64(10), int64(0)).Return(expectedNotifs, int64(1), nil)

	notifs, total, err := service.ListNotifications(context.Background(), userID, 10, 0)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Equal(t, expectedNotifs, notifs)
	mockRepo.AssertExpectations(t)
}
