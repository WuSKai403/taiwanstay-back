package service

import (
	"context"
	"time"

	"github.com/taiwanstay/taiwanstay-back/internal/domain"
	"github.com/taiwanstay/taiwanstay-back/internal/repository"
	"go.mongodb.org/mongo-driver/bson"
)

type AdminService interface {
	GetSystemStats(ctx context.Context) (map[string]int64, error)
	ListPendingImages(ctx context.Context, limit, offset int64) ([]*domain.Image, int64, error)
	ReviewImage(ctx context.Context, imageID string, approved bool) error
	ListUsers(ctx context.Context, role domain.UserRole, limit, offset int64) ([]*domain.User, int64, error)
	UpdateUserStatus(ctx context.Context, userID string, status domain.UserStatus) error
}

type adminService struct {
	userRepo     repository.UserRepository
	imageRepo    repository.ImageRepository
	appRepo      repository.ApplicationRepository
	imageService ImageService
}

func NewAdminService(userRepo repository.UserRepository, imageRepo repository.ImageRepository, appRepo repository.ApplicationRepository, imageService ImageService) AdminService {
	return &adminService{
		userRepo:     userRepo,
		imageRepo:    imageRepo,
		appRepo:      appRepo,
		imageService: imageService,
	}
}

func (s *adminService) GetSystemStats(ctx context.Context) (map[string]int64, error) {
	stats := make(map[string]int64)

	// 1. Total Users
	userCount, err := s.userRepo.Count(ctx)
	if err != nil {
		return nil, err
	}
	stats["totalUsers"] = userCount

	// 2. Pending Images
	pendingImageCount, err := s.imageRepo.CountByStatus(ctx, domain.ImageStatusPending)
	if err != nil {
		return nil, err
	}
	stats["pendingImages"] = pendingImageCount

	// 3. Today's Applications
	appCount, err := s.appRepo.CountByDate(ctx, time.Now())
	if err != nil {
		return nil, err
	}
	stats["todayApplications"] = appCount

	return stats, nil
}

func (s *adminService) ListPendingImages(ctx context.Context, limit, offset int64) ([]*domain.Image, int64, error) {
	return s.imageRepo.ListByStatus(ctx, domain.ImageStatusPending, limit, offset)
}

func (s *adminService) ReviewImage(ctx context.Context, imageID string, approved bool) error {
	status := domain.ImageStatusRejected
	if approved {
		status = domain.ImageStatusApproved
	}
	return s.imageService.UpdateImageStatus(ctx, imageID, status)
}

func (s *adminService) ListUsers(ctx context.Context, role domain.UserRole, limit, offset int64) ([]*domain.User, int64, error) {
	filter := bson.M{}
	if role != "" {
		filter["role"] = role
	}
	return s.userRepo.List(ctx, filter, limit, offset)
}

func (s *adminService) UpdateUserStatus(ctx context.Context, userID string, status domain.UserStatus) error {
	return s.userRepo.UpdateStatus(ctx, userID, status)
}
