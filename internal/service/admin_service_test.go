package service

import (
	"context"
	"io"
	"mime/multipart"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/taiwanstay/taiwanstay-back/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
)

// MockImageRepository
type MockImageRepository struct {
	mock.Mock
}

func (m *MockImageRepository) Create(ctx context.Context, image *domain.Image) error {
	args := m.Called(ctx, image)
	return args.Error(0)
}

func (m *MockImageRepository) GetByID(ctx context.Context, id string) (*domain.Image, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Image), args.Error(1)
}

func (m *MockImageRepository) UpdateStatus(ctx context.Context, id string, status domain.ImageStatus) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockImageRepository) CountByStatus(ctx context.Context, status domain.ImageStatus) (int64, error) {
	args := m.Called(ctx, status)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockImageRepository) ListByStatus(ctx context.Context, status domain.ImageStatus, limit, offset int64) ([]*domain.Image, int64, error) {
	args := m.Called(ctx, status, limit, offset)
	return args.Get(0).([]*domain.Image), args.Get(1).(int64), args.Error(2)
}

func (m *MockUserRepository) UpdateStatus(ctx context.Context, id string, status domain.UserStatus) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

// MockImageService (reusing interface from image_service.go)
type MockImageService struct {
	mock.Mock
}

func (m *MockImageService) UploadImage(ctx context.Context, file multipart.File, header *multipart.FileHeader, userID string) (*domain.Image, error) {
	// Simplified signature for mock, as we don't use this in AdminService
	return nil, nil
}

func (m *MockImageService) GetImage(ctx context.Context, id string) (*domain.Image, error) {
	return nil, nil
}

func (m *MockImageService) UpdateImageStatus(ctx context.Context, id string, status domain.ImageStatus) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockImageService) GetImageContent(ctx context.Context, id string) (io.ReadCloser, error) {
	return nil, nil
}

func TestGetSystemStats(t *testing.T) {
	mockUserRepo := new(MockUserRepository) // Reusing from notification_service_test.go if in same package, but we are in same package 'service' so it should be available?
	// Wait, MockUserRepository is defined in notification_service_test.go which is in package service_test or service?
	// It's in package service. So it is available.
	mockImageRepo := new(MockImageRepository)
	mockAppRepo := new(MockApplicationRepository)
	mockImageService := new(MockImageService)

	adminService := NewAdminService(mockUserRepo, mockImageRepo, mockAppRepo, mockImageService)

	ctx := context.Background()

	mockUserRepo.On("Count", ctx).Return(int64(10), nil)
	mockImageRepo.On("CountByStatus", ctx, domain.ImageStatusPending).Return(int64(5), nil)
	mockAppRepo.On("CountByDate", ctx, mock.Anything).Return(int64(3), nil)

	stats, err := adminService.GetSystemStats(ctx)

	assert.NoError(t, err)
	assert.Equal(t, int64(10), stats["totalUsers"])
	assert.Equal(t, int64(5), stats["pendingImages"])
	assert.Equal(t, int64(3), stats["todayApplications"])
}

func TestReviewImage(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockImageRepo := new(MockImageRepository)
	mockAppRepo := new(MockApplicationRepository)
	mockImageService := new(MockImageService)

	adminService := NewAdminService(mockUserRepo, mockImageRepo, mockAppRepo, mockImageService)

	ctx := context.Background()
	imageID := "img123"

	// Approve
	mockImageService.On("UpdateImageStatus", ctx, imageID, domain.ImageStatusApproved).Return(nil)
	err := adminService.ReviewImage(ctx, imageID, true)
	assert.NoError(t, err)

	// Reject
	mockImageService.On("UpdateImageStatus", ctx, imageID, domain.ImageStatusRejected).Return(nil)
	err = adminService.ReviewImage(ctx, imageID, false)
	assert.NoError(t, err)
}

func TestListUsers(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockImageRepo := new(MockImageRepository)
	mockAppRepo := new(MockApplicationRepository)
	mockImageService := new(MockImageService)

	adminService := NewAdminService(mockUserRepo, mockImageRepo, mockAppRepo, mockImageService)

	ctx := context.Background()
	expectedUsers := []*domain.User{{Name: "Test"}}

	// Filter by Role
	mockUserRepo.On("List", ctx, bson.M{"role": domain.RoleHost}, int64(10), int64(0)).Return(expectedUsers, int64(1), nil)
	users, total, err := adminService.ListUsers(ctx, domain.RoleHost, 10, 0)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Equal(t, expectedUsers, users)

	// No Filter
	mockUserRepo.On("List", ctx, bson.M{}, int64(10), int64(0)).Return(expectedUsers, int64(1), nil)
	users, total, err = adminService.ListUsers(ctx, "", 10, 0)
	assert.NoError(t, err)
}

func TestUpdateUserStatus(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockImageRepo := new(MockImageRepository)
	mockAppRepo := new(MockApplicationRepository)
	mockImageService := new(MockImageService)

	adminService := NewAdminService(mockUserRepo, mockImageRepo, mockAppRepo, mockImageService)

	ctx := context.Background()
	userID := "user123"
	status := domain.UserStatusSuspended

	mockUserRepo.On("UpdateStatus", ctx, userID, status).Return(nil)

	err := adminService.UpdateUserStatus(ctx, userID, status)
	assert.NoError(t, err)
	mockUserRepo.AssertExpectations(t)
}
