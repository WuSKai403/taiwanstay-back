package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/taiwanstay/taiwanstay-back/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MockBookmarkRepository
type MockBookmarkRepository struct {
	mock.Mock
}

func (m *MockBookmarkRepository) Create(ctx context.Context, bookmark *domain.Bookmark) error {
	args := m.Called(ctx, bookmark)
	return args.Error(0)
}

func (m *MockBookmarkRepository) Delete(ctx context.Context, userID, opportunityID string) error {
	args := m.Called(ctx, userID, opportunityID)
	return args.Error(0)
}

func (m *MockBookmarkRepository) ListByUserID(ctx context.Context, userID string, limit, offset int64) ([]*domain.Bookmark, int64, error) {
	args := m.Called(ctx, userID, limit, offset)
	return args.Get(0).([]*domain.Bookmark), args.Get(1).(int64), args.Error(2)
}

func (m *MockBookmarkRepository) Exists(ctx context.Context, userID, opportunityID string) (bool, error) {
	args := m.Called(ctx, userID, opportunityID)
	return args.Bool(0), args.Error(1)
}

func TestAddBookmark(t *testing.T) {
	mockBookmarkRepo := new(MockBookmarkRepository)
	mockOppRepo := new(MockOpportunityRepository)
	service := NewBookmarkService(mockBookmarkRepo, mockOppRepo)

	ctx := context.Background()
	userID := "user1"
	oppID := "507f1f77bcf86cd799439011" // Valid ObjectID hex string
	oppObjectID, _ := primitive.ObjectIDFromHex(oppID)

	// Case 1: Success
	mockOppRepo.On("GetByID", ctx, oppID).Return(&domain.Opportunity{ID: oppObjectID}, nil).Once()
	mockBookmarkRepo.On("Exists", ctx, userID, oppID).Return(false, nil).Once()
	mockBookmarkRepo.On("Create", ctx, mock.AnythingOfType("*domain.Bookmark")).Return(nil).Once()

	err := service.AddBookmark(ctx, userID, oppID)
	assert.NoError(t, err)

	// Case 2: Already Exists
	mockOppRepo.On("GetByID", ctx, oppID).Return(&domain.Opportunity{ID: oppObjectID}, nil).Once()
	mockBookmarkRepo.On("Exists", ctx, userID, oppID).Return(true, nil).Once()

	err = service.AddBookmark(ctx, userID, oppID)
	assert.Error(t, err)
	assert.Equal(t, ErrBookmarkAlreadyExists, err)
}

func TestRemoveBookmark(t *testing.T) {
	mockBookmarkRepo := new(MockBookmarkRepository)
	mockOppRepo := new(MockOpportunityRepository)
	service := NewBookmarkService(mockBookmarkRepo, mockOppRepo)

	ctx := context.Background()
	userID := "user1"
	oppID := "opp1"

	mockBookmarkRepo.On("Delete", ctx, userID, oppID).Return(nil).Once()

	err := service.RemoveBookmark(ctx, userID, oppID)
	assert.NoError(t, err)
}

func TestListUserBookmarks(t *testing.T) {
	mockBookmarkRepo := new(MockBookmarkRepository)
	mockOppRepo := new(MockOpportunityRepository)
	service := NewBookmarkService(mockBookmarkRepo, mockOppRepo)

	ctx := context.Background()
	userID := "user1"
	expectedBookmarks := []*domain.Bookmark{{UserID: userID, OpportunityID: "opp1"}}

	mockBookmarkRepo.On("ListByUserID", ctx, userID, int64(10), int64(0)).Return(expectedBookmarks, int64(1), nil).Once()

	bookmarks, total, err := service.ListUserBookmarks(ctx, userID, 10, 0)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Equal(t, expectedBookmarks, bookmarks)
}
