package service

import (
	"context"
	"errors"

	"github.com/taiwanstay/taiwanstay-back/internal/domain"
	"github.com/taiwanstay/taiwanstay-back/internal/repository"
	"go.mongodb.org/mongo-driver/bson"
)

type ApplicationService interface {
	CreateApplication(ctx context.Context, app *domain.Application) (*domain.Application, error)
	GetApplicationByID(ctx context.Context, id string) (*domain.Application, error)
	ListApplications(ctx context.Context, filter bson.M, limit, offset int64) ([]*domain.Application, int64, error)
	UpdateApplication(ctx context.Context, id string, app *domain.Application) error
	UpdateApplicationStatus(ctx context.Context, id string, status domain.ApplicationStatus, note string, userID string) error
	DeleteApplication(ctx context.Context, id string, userID string) error
}

type applicationService struct {
	repo    repository.ApplicationRepository
	oppRepo repository.OpportunityRepository
}

func NewApplicationService(repo repository.ApplicationRepository, oppRepo repository.OpportunityRepository) ApplicationService {
	return &applicationService{
		repo:    repo,
		oppRepo: oppRepo,
	}
}

func (s *applicationService) CreateApplication(ctx context.Context, app *domain.Application) (*domain.Application, error) {
	// 1. Check Opportunity existence
	opp, err := s.oppRepo.GetByID(ctx, app.OpportunityID.Hex())
	if err != nil {
		return nil, errors.New("opportunity not found")
	}

	// 2. Validate TimeSlot (if applicable)
	if opp.HasTimeSlots {
		valid := false
		reqStart := app.ApplicationDetails.StartDate
		reqEnd := app.ApplicationDetails.EndDate

		for _, slot := range opp.TimeSlots {
			if slot.Status == domain.TimeSlotStatusOpen &&
				slot.StartDate <= reqStart &&
				slot.EndDate >= reqEnd {
				valid = true
				break
			}
		}
		if !valid {
			return nil, errors.New("selected dates are not available in any open time slot")
		}
	}

	// 3. Set HostID from Opportunity
	app.HostID = opp.HostID
	app.Status = domain.ApplicationStatusPending // Default to Pending

	err = s.repo.Create(ctx, app)
	if err != nil {
		return nil, err
	}
	return app, nil
}

func (s *applicationService) GetApplicationByID(ctx context.Context, id string) (*domain.Application, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *applicationService) ListApplications(ctx context.Context, filter bson.M, limit, offset int64) ([]*domain.Application, int64, error) {
	return s.repo.List(ctx, filter, limit, offset)
}

func (s *applicationService) UpdateApplication(ctx context.Context, id string, app *domain.Application) error {
	return s.repo.Update(ctx, id, app)
}

func (s *applicationService) UpdateApplicationStatus(ctx context.Context, id string, status domain.ApplicationStatus, note string, userID string) error {
	app, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Verify ownership (Host) - This logic might be better in Handler or Middleware, but service check is safe
	// Here we just update the status
	app.Status = status
	app.StatusNote = note

	// If Accepted/Rejected, update ReviewDetails? Or just StatusHistory?
	// For MVP, just update Status.

	return s.repo.Update(ctx, id, app)
}

func (s *applicationService) DeleteApplication(ctx context.Context, id string, userID string) error {
	app, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Only allow deleting if status is DRAFT or PENDING
	if app.Status != domain.ApplicationStatusDraft && app.Status != domain.ApplicationStatusPending {
		return errors.New("cannot delete application that is not draft or pending")
	}

	// Verify ownership
	if app.UserID.Hex() != userID {
		return errors.New("unauthorized to delete this application")
	}

	return s.repo.Delete(ctx, id)
}
