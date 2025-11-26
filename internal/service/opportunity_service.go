package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/taiwanstay/taiwanstay-back/internal/domain"
	"github.com/taiwanstay/taiwanstay-back/internal/repository"
	"go.mongodb.org/mongo-driver/bson"
)

type OpportunityService interface {
	CreateOpportunity(ctx context.Context, opp *domain.Opportunity) (*domain.Opportunity, error)
	GetOpportunityByID(ctx context.Context, id string) (*domain.Opportunity, error)
	ListOpportunities(ctx context.Context, filter bson.M, limit, offset int64) ([]*domain.Opportunity, error)
	UpdateOpportunity(ctx context.Context, id string, opp *domain.Opportunity) error
	DeleteOpportunity(ctx context.Context, id string) error
	SearchOpportunities(ctx context.Context, filter repository.OpportunityFilter) ([]*domain.Opportunity, int64, error)
}

type opportunityService struct {
	repo repository.OpportunityRepository
}

func NewOpportunityService(repo repository.OpportunityRepository) OpportunityService {
	return &opportunityService{repo: repo}
}

func (s *opportunityService) CreateOpportunity(ctx context.Context, opp *domain.Opportunity) (*domain.Opportunity, error) {
	// Generate Slug
	if opp.Slug == "" {
		opp.Slug = generateSlug(opp.Title)
	}

	// Generate PublicID
	if opp.PublicID == "" {
		opp.PublicID = uuid.New().String()
	}

	// Set default status
	if opp.Status == "" {
		opp.Status = domain.OpportunityStatusDraft
	}

	err := s.repo.Create(ctx, opp)
	if err != nil {
		return nil, err
	}
	return opp, nil
}

func (s *opportunityService) GetOpportunityByID(ctx context.Context, id string) (*domain.Opportunity, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *opportunityService) ListOpportunities(ctx context.Context, filter bson.M, limit, offset int64) ([]*domain.Opportunity, error) {
	return s.repo.List(ctx, filter, limit, offset)
}

func (s *opportunityService) UpdateOpportunity(ctx context.Context, id string, opp *domain.Opportunity) error {
	return s.repo.Update(ctx, id, opp)
}

func (s *opportunityService) DeleteOpportunity(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *opportunityService) SearchOpportunities(ctx context.Context, filter repository.OpportunityFilter) ([]*domain.Opportunity, int64, error) {
	return s.repo.Search(ctx, filter)
}
