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

type HostRepository interface {
	Create(ctx context.Context, host *domain.Host) error
	GetByID(ctx context.Context, id string) (*domain.Host, error)
	GetByUserID(ctx context.Context, userID string) (*domain.Host, error)
	Update(ctx context.Context, id string, host *domain.Host) error
}

type mongoHostRepository struct {
	collection *mongo.Collection
}

func NewHostRepository(collection *mongo.Collection) HostRepository {
	// Create Indexes
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 2dsphere index for location
	collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "location.coordinates", Value: "2dsphere"}},
	})

	// Text index for name and description
	collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "name", Value: "text"}, {Key: "description", Value: "text"}},
	})

	// Unique index for slug
	collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "slug", Value: 1}},
		Options: options.Index().SetUnique(true),
	})

	return &mongoHostRepository{collection: collection}
}

func (r *mongoHostRepository) Create(ctx context.Context, host *domain.Host) error {
	host.CreatedAt = time.Now()
	host.UpdatedAt = time.Now()
	res, err := r.collection.InsertOne(ctx, host)
	if err != nil {
		return err
	}
	host.ID = res.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *mongoHostRepository) GetByID(ctx context.Context, id string) (*domain.Host, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var host domain.Host
	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&host)
	if err != nil {
		return nil, err
	}
	return &host, nil
}

func (r *mongoHostRepository) GetByUserID(ctx context.Context, userID string) (*domain.Host, error) {
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}
	var host domain.Host
	err = r.collection.FindOne(ctx, bson.M{"userId": objID}).Decode(&host)
	if err != nil {
		return nil, err
	}
	return &host, nil
}

func (r *mongoHostRepository) Update(ctx context.Context, id string, host *domain.Host) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	host.UpdatedAt = time.Now()

	// Ensure the ID in the struct matches the ID we are updating
	host.ID = objID

	_, err = r.collection.ReplaceOne(ctx, bson.M{"_id": objID}, host)
	return err
}
