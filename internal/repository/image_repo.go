package repository

import (
	"context"
	"time"

	"github.com/taiwanstay/taiwanstay-back/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ImageRepository interface {
	Create(ctx context.Context, image *domain.Image) error
	GetByID(ctx context.Context, id string) (*domain.Image, error)
	UpdateStatus(ctx context.Context, id string, status domain.ImageStatus) error
	CountByStatus(ctx context.Context, status domain.ImageStatus) (int64, error)
	ListByStatus(ctx context.Context, status domain.ImageStatus, limit, offset int64) ([]*domain.Image, int64, error)
}

type mongoImageRepository struct {
	collection *mongo.Collection
}

func NewImageRepository(collection *mongo.Collection) ImageRepository {
	// Create Indexes
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Indexes for filtering
	collection.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "userId", Value: 1}}},
		{Keys: bson.D{{Key: "status", Value: 1}}},
	})

	return &mongoImageRepository{collection: collection}
}

func (r *mongoImageRepository) Create(ctx context.Context, image *domain.Image) error {
	image.CreatedAt = time.Now()
	image.UpdatedAt = time.Now()
	res, err := r.collection.InsertOne(ctx, image)
	if err != nil {
		return err
	}
	image.ID = res.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *mongoImageRepository) GetByID(ctx context.Context, id string) (*domain.Image, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var image domain.Image
	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&image)
	if err != nil {
		return nil, err
	}
	return &image, nil
}

func (r *mongoImageRepository) UpdateStatus(ctx context.Context, id string, status domain.ImageStatus) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	update := bson.M{
		"$set": bson.M{
			"status":    status,
			"updatedAt": time.Now(),
		},
	}
	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	return err
}

func (r *mongoImageRepository) CountByStatus(ctx context.Context, status domain.ImageStatus) (int64, error) {
	return r.collection.CountDocuments(ctx, bson.M{"status": status})
}

func (r *mongoImageRepository) ListByStatus(ctx context.Context, status domain.ImageStatus, limit, offset int64) ([]*domain.Image, int64, error) {
	filter := bson.M{"status": status}
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

	var images []*domain.Image
	if err := cursor.All(ctx, &images); err != nil {
		return nil, 0, err
	}

	return images, total, nil
}
