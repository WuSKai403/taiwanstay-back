package service

import (
	"context"
	"time"

	"github.com/taiwanstay/taiwanstay-back/internal/domain"
	"github.com/taiwanstay/taiwanstay-back/internal/repository"
	"github.com/taiwanstay/taiwanstay-back/pkg/email"
	"github.com/taiwanstay/taiwanstay-back/pkg/logger"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NotificationService interface {
	SendNotification(ctx context.Context, userID string, notifType domain.NotificationType, title, message string, data map[string]string) error
	ListNotifications(ctx context.Context, userID string, limit, offset int64) ([]*domain.Notification, int64, error)
	MarkAsRead(ctx context.Context, id string, userID string) error
	MarkAllAsRead(ctx context.Context, userID string) error
}

type notificationService struct {
	repo        repository.NotificationRepository
	userRepo    repository.UserRepository
	emailSender email.EmailSender
}

func NewNotificationService(repo repository.NotificationRepository, userRepo repository.UserRepository, emailSender email.EmailSender) NotificationService {
	return &notificationService{
		repo:        repo,
		userRepo:    userRepo,
		emailSender: emailSender,
	}
}

func (s *notificationService) SendNotification(ctx context.Context, userID string, notifType domain.NotificationType, title, message string, data map[string]string) error {
	// 1. Save In-App Notification
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	notification := &domain.Notification{
		UserID:    userObjID,
		Type:      notifType,
		Title:     title,
		Message:   message,
		IsRead:    false,
		Data:      data,
		CreatedAt: time.Now(),
	}

	if err := s.repo.Create(ctx, notification); err != nil {
		logger.Error("Failed to create in-app notification", "error", err)
		// We continue to try sending email even if DB fails?
		// Ideally yes, but for consistency let's log and proceed.
	}

	// 2. Send Email (Async)
	// We need to fetch user email first
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		logger.Error("Failed to get user for email notification", "userId", userID, "error", err)
		return nil // Don't fail the whole operation if user not found for email
	}

	go func() {
		// Use a background context for async email sending
		// In production, use a proper worker queue
		err := s.emailSender.Send(user.Email, user.Name, title, message)
		if err != nil {
			logger.Error("Failed to send email notification", "to", user.Email, "error", err)
		}
	}()

	return nil
}

func (s *notificationService) ListNotifications(ctx context.Context, userID string, limit, offset int64) ([]*domain.Notification, int64, error) {
	return s.repo.ListByUserID(ctx, userID, limit, offset)
}

func (s *notificationService) MarkAsRead(ctx context.Context, id string, userID string) error {
	return s.repo.MarkAsRead(ctx, id, userID)
}

func (s *notificationService) MarkAllAsRead(ctx context.Context, userID string) error {
	return s.repo.MarkAllAsRead(ctx, userID)
}
