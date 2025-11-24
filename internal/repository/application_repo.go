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

type ApplicationRepository interface {
	Create(ctx context.Context, app *domain.Application) error
	GetByID(ctx context.Context, id string) (*domain.Application, error)
	List(ctx context.Context, filter bson.M, limit, offset int64) ([]*domain.Application, int64, error)
	Update(ctx context.Context, id string, app *domain.Application) error
	Delete(ctx context.Context, id string) error
}

type mongoApplicationRepository struct {
	collection *mongo.Collection
}

func NewApplicationRepository(collection *mongo.Collection) ApplicationRepository {
	// Create Indexes
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Compound index for user and opportunity to prevent duplicate applications
	collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "userId", Value: 1}, {Key: "opportunityId", Value: 1}},
		Options: options.Index().SetUnique(true),
	})

	// Indexes for filtering
	collection.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "hostId", Value: 1}}},
		{Keys: bson.D{{Key: "status", Value: 1}}},
		{Keys: bson.D{{Key: "userId", Value: 1}}},
	})

	return &mongoApplicationRepository{collection: collection}
}

func (r *mongoApplicationRepository) Create(ctx context.Context, app *domain.Application) error {
	app.CreatedAt = time.Now()
	app.UpdatedAt = time.Now()
	res, err := r.collection.InsertOne(ctx, app)
	if err != nil {
		return err
	}
	app.ID = res.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *mongoApplicationRepository) GetByID(ctx context.Context, id string) (*domain.Application, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var app domain.Application
	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&app)
	if err != nil {
		return nil, err
	}
	return &app, nil
}

func (r *mongoApplicationRepository) List(ctx context.Context, filter bson.M, limit, offset int64) ([]*domain.Application, int64, error) {
	// Count total
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

	var apps []*domain.Application
	if err := cursor.All(ctx, &apps); err != nil {
		return nil, 0, err
	}
	return apps, total, nil
}

func (r *mongoApplicationRepository) Update(ctx context.Context, id string, app *domain.Application) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	app.UpdatedAt = time.Now()
	_, err = r.collection.ReplaceOne(ctx, bson.M{"_id": objID}, app)
	return err
}

func (r *mongoApplicationRepository) Delete(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objID})
	return err
}
