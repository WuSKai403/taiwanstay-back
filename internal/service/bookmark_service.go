package service

import (
	"context"
	"errors"

	"github.com/taiwanstay/taiwanstay-back/internal/domain"
	"github.com/taiwanstay/taiwanstay-back/internal/repository"
)

var (
	ErrBookmarkAlreadyExists = errors.New("bookmark already exists")
	ErrBookmarkNotFound      = errors.New("bookmark not found")
)

type BookmarkService interface {
	AddBookmark(ctx context.Context, userID, opportunityID string) error
	RemoveBookmark(ctx context.Context, userID, opportunityID string) error
	ListUserBookmarks(ctx context.Context, userID string, limit, offset int64) ([]*domain.Bookmark, int64, error)
}

type bookmarkService struct {
	bookmarkRepo repository.BookmarkRepository
	oppRepo      repository.OpportunityRepository
}

func NewBookmarkService(bookmarkRepo repository.BookmarkRepository, oppRepo repository.OpportunityRepository) BookmarkService {
	return &bookmarkService{
		bookmarkRepo: bookmarkRepo,
		oppRepo:      oppRepo,
	}
}

func (s *bookmarkService) AddBookmark(ctx context.Context, userID, opportunityID string) error {
	// Check if opportunity exists
	_, err := s.oppRepo.GetByID(ctx, opportunityID)
	if err != nil {
		return err // Could be ErrOpportunityNotFound if that's what repo returns
	}

	// Check if already bookmarked
	exists, err := s.bookmarkRepo.Exists(ctx, userID, opportunityID)
	if err != nil {
		return err
	}
	if exists {
		return ErrBookmarkAlreadyExists
	}

	bookmark := &domain.Bookmark{
		UserID:        userID,
		OpportunityID: opportunityID,
	}

	return s.bookmarkRepo.Create(ctx, bookmark)
}

func (s *bookmarkService) RemoveBookmark(ctx context.Context, userID, opportunityID string) error {
	return s.bookmarkRepo.Delete(ctx, userID, opportunityID)
}

func (s *bookmarkService) ListUserBookmarks(ctx context.Context, userID string, limit, offset int64) ([]*domain.Bookmark, int64, error) {
	return s.bookmarkRepo.ListByUserID(ctx, userID, limit, offset)
}
