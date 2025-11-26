package repository

import (
	"context"
	"time"

	"github.com/taiwanstay/taiwanstay-back/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BookmarkRepository interface {
	Create(ctx context.Context, bookmark *domain.Bookmark) error
	Delete(ctx context.Context, userID, opportunityID string) error
	ListByUserID(ctx context.Context, userID string, limit, offset int64) ([]*domain.Bookmark, int64, error)
	Exists(ctx context.Context, userID, opportunityID string) (bool, error)
}

type mongoBookmarkRepository struct {
	collection *mongo.Collection
}

func NewBookmarkRepository(collection *mongo.Collection) BookmarkRepository {
	// Create Indexes
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Compound unique index for user and opportunity
	collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "userId", Value: 1}, {Key: "opportunityId", Value: 1}},
		Options: options.Index().SetUnique(true),
	})

	// Index for userId
	collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "userId", Value: 1}},
	})

	return &mongoBookmarkRepository{collection: collection}
}

func (r *mongoBookmarkRepository) Create(ctx context.Context, bookmark *domain.Bookmark) error {
	bookmark.CreatedAt = time.Now()
	_, err := r.collection.InsertOne(ctx, bookmark)
	return err
}

func (r *mongoBookmarkRepository) Delete(ctx context.Context, userID, opportunityID string) error {
	filter := bson.M{
		"userId":        userID,
		"opportunityId": opportunityID,
	}
	_, err := r.collection.DeleteOne(ctx, filter)
	return err
}

func (r *mongoBookmarkRepository) ListByUserID(ctx context.Context, userID string, limit, offset int64) ([]*domain.Bookmark, int64, error) {
	filter := bson.M{"userId": userID}

	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	opts := options.Find().SetLimit(limit).SetSkip(offset).SetSort(bson.D{{Key: "createdAt", Value: -1}})
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var bookmarks []*domain.Bookmark
	if err := cursor.All(ctx, &bookmarks); err != nil {
		return nil, 0, err
	}

	return bookmarks, total, nil
}

func (r *mongoBookmarkRepository) Exists(ctx context.Context, userID, opportunityID string) (bool, error) {
	filter := bson.M{
		"userId":        userID,
		"opportunityId": opportunityID,
	}
	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
