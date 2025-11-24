package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/taiwanstay/taiwanstay-back/internal/domain"
	"github.com/taiwanstay/taiwanstay-back/internal/repository"
)

type HostService interface {
	CreateHost(ctx context.Context, host *domain.Host) (*domain.Host, error)
	GetHostByUserID(ctx context.Context, userID string) (*domain.Host, error)
	GetHostByID(ctx context.Context, id string) (*domain.Host, error)
	UpdateHost(ctx context.Context, id string, host *domain.Host) error
}

type hostService struct {
	repo repository.HostRepository
}

func NewHostService(repo repository.HostRepository) HostService {
	return &hostService{repo: repo}
}

func (s *hostService) CreateHost(ctx context.Context, host *domain.Host) (*domain.Host, error) {
	// Generate Slug
	if host.Slug == "" {
		host.Slug = generateSlug(host.Name)
	}

	// Set default status if not provided
	if host.Status == "" {
		host.Status = domain.HostStatusPending
	}

	err := s.repo.Create(ctx, host)
	if err != nil {
		return nil, err
	}
	return host, nil
}

func (s *hostService) GetHostByUserID(ctx context.Context, userID string) (*domain.Host, error) {
	return s.repo.GetByUserID(ctx, userID)
}

func (s *hostService) GetHostByID(ctx context.Context, id string) (*domain.Host, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *hostService) UpdateHost(ctx context.Context, id string, host *domain.Host) error {
	return s.repo.Update(ctx, id, host)
}

// Helper to generate simple slug
func generateSlug(name string) string {
	slug := strings.ToLower(name)
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = fmt.Sprintf("%s-%s", slug, uuid.New().String()[:8])
	return slug
}
