package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/taiwanstay/taiwanstay-back/internal/domain"
	"github.com/taiwanstay/taiwanstay-back/internal/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Mocks
type MockApplicationRepository struct {
	mock.Mock
}

func (m *MockApplicationRepository) Create(ctx context.Context, app *domain.Application) error {
	args := m.Called(ctx, app)
	return args.Error(0)
}

func (m *MockApplicationRepository) GetByID(ctx context.Context, id string) (*domain.Application, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Application), args.Error(1)
}

func (m *MockApplicationRepository) List(ctx context.Context, filter bson.M, limit, offset int64) ([]*domain.Application, int64, error) {
	args := m.Called(ctx, filter, limit, offset)
	return args.Get(0).([]*domain.Application), args.Get(1).(int64), args.Error(2)
}

func (m *MockApplicationRepository) Update(ctx context.Context, id string, app *domain.Application) error {
	args := m.Called(ctx, id, app)
	return args.Error(0)
}

func (m *MockApplicationRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type MockOpportunityRepository struct {
	mock.Mock
}

func (m *MockOpportunityRepository) Create(ctx context.Context, opp *domain.Opportunity) error {
	args := m.Called(ctx, opp)
	return args.Error(0)
}

func (m *MockOpportunityRepository) GetByID(ctx context.Context, id string) (*domain.Opportunity, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Opportunity), args.Error(1)
}

func (m *MockOpportunityRepository) List(ctx context.Context, filter bson.M, limit, offset int64) ([]*domain.Opportunity, error) {
	args := m.Called(ctx, filter, limit, offset)
	return args.Get(0).([]*domain.Opportunity), args.Error(1)
}

func (m *MockOpportunityRepository) Update(ctx context.Context, id string, opp *domain.Opportunity) error {
	args := m.Called(ctx, id, opp)
	return args.Error(0)
}

func (m *MockOpportunityRepository) Search(ctx context.Context, filter repository.OpportunityFilter) ([]*domain.Opportunity, int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]*domain.Opportunity), args.Get(1).(int64), args.Error(2)
}

// Tests
func TestCreateApplication_Success(t *testing.T) {
	mockAppRepo := new(MockApplicationRepository)
	mockOppRepo := new(MockOpportunityRepository)
	service := NewApplicationService(mockAppRepo, mockOppRepo)

	ctx := context.Background()
	oppID := primitive.NewObjectID()
	hostID := primitive.NewObjectID()

	// Opportunity with valid time slot
	opp := &domain.Opportunity{
		ID:           oppID,
		HostID:       hostID,
		HasTimeSlots: true,
		TimeSlots: []domain.TimeSlot{
			{
				StartDate: "2023-01-01",
				EndDate:   "2023-01-31",
				Status:    domain.TimeSlotStatusOpen,
			},
		},
	}

	app := &domain.Application{
		OpportunityID: oppID,
		ApplicationDetails: domain.ApplicationDetails{
			StartDate: "2023-01-05",
			EndDate:   "2023-01-10",
		},
	}

	mockOppRepo.On("GetByID", ctx, oppID.Hex()).Return(opp, nil)
	mockAppRepo.On("Create", ctx, app).Return(nil)

	createdApp, err := service.CreateApplication(ctx, app)

	assert.NoError(t, err)
	assert.NotNil(t, createdApp)
	assert.Equal(t, domain.ApplicationStatusPending, createdApp.Status)
	assert.Equal(t, hostID, createdApp.HostID)
	mockOppRepo.AssertExpectations(t)
	mockAppRepo.AssertExpectations(t)
}

func TestCreateApplication_InvalidDates(t *testing.T) {
	mockAppRepo := new(MockApplicationRepository)
	mockOppRepo := new(MockOpportunityRepository)
	service := NewApplicationService(mockAppRepo, mockOppRepo)

	ctx := context.Background()
	oppID := primitive.NewObjectID()

	// Opportunity with time slot NOT covering application dates
	opp := &domain.Opportunity{
		ID:           oppID,
		HasTimeSlots: true,
		TimeSlots: []domain.TimeSlot{
			{
				StartDate: "2023-01-01",
				EndDate:   "2023-01-31",
				Status:    domain.TimeSlotStatusOpen,
			},
		},
	}

	app := &domain.Application{
		OpportunityID: oppID,
		ApplicationDetails: domain.ApplicationDetails{
			StartDate: "2023-02-01", // Outside range
			EndDate:   "2023-02-05",
		},
	}

	mockOppRepo.On("GetByID", ctx, oppID.Hex()).Return(opp, nil)

	createdApp, err := service.CreateApplication(ctx, app)

	assert.Error(t, err)
	assert.Nil(t, createdApp)
	assert.Contains(t, err.Error(), "selected dates are not available")
	mockOppRepo.AssertExpectations(t)
}
