package repository

import (
	"context"
	"time"

	"github.com/taiwanstay/taiwanstay-back/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ImageRepository interface {
	Create(ctx context.Context, image *domain.Image) error
	GetByID(ctx context.Context, id string) (*domain.Image, error)
	UpdateStatus(ctx context.Context, id string, status domain.ImageStatus) error
}

type mongoImageRepository struct {
	collection *mongo.Collection
}

func NewImageRepository(collection *mongo.Collection) ImageRepository {
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
	return err
}
