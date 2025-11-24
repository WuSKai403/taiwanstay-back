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

type OpportunityRepository interface {
	Create(ctx context.Context, opp *domain.Opportunity) error
	GetByID(ctx context.Context, id string) (*domain.Opportunity, error)
	List(ctx context.Context, filter bson.M, limit, offset int64) ([]*domain.Opportunity, error)
	Update(ctx context.Context, id string, opp *domain.Opportunity) error
}

type mongoOpportunityRepository struct {
	collection *mongo.Collection
}

func NewOpportunityRepository(collection *mongo.Collection) OpportunityRepository {
	// Create Indexes
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 2dsphere index for location
	collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "location.coordinates", Value: "2dsphere"}},
	})

	// Unique index for slug and publicId
	collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "slug", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "publicId", Value: 1}},
		Options: options.Index().SetUnique(true),
	})

	// Other indexes for filtering
	collection.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "hostId", Value: 1}}},
		{Keys: bson.D{{Key: "status", Value: 1}}},
		{Keys: bson.D{{Key: "type", Value: 1}}},
		{Keys: bson.D{{Key: "location.city", Value: 1}}},
		{Keys: bson.D{{Key: "location.country", Value: 1}}},
	})

	return &mongoOpportunityRepository{collection: collection}
}

func (r *mongoOpportunityRepository) Create(ctx context.Context, opp *domain.Opportunity) error {
	opp.CreatedAt = time.Now()
	opp.UpdatedAt = time.Now()
	res, err := r.collection.InsertOne(ctx, opp)
	if err != nil {
		return err
	}
	opp.ID = res.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *mongoOpportunityRepository) GetByID(ctx context.Context, id string) (*domain.Opportunity, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var opp domain.Opportunity
	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&opp)
	if err != nil {
		return nil, err
	}
	return &opp, nil
}

func (r *mongoOpportunityRepository) List(ctx context.Context, filter bson.M, limit, offset int64) ([]*domain.Opportunity, error) {
	opts := options.Find().SetLimit(limit).SetSkip(offset).SetSort(bson.D{{Key: "createdAt", Value: -1}})
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var opps []*domain.Opportunity
	if err := cursor.All(ctx, &opps); err != nil {
		return nil, err
	}
	return opps, nil
}

func (r *mongoOpportunityRepository) Update(ctx context.Context, id string, opp *domain.Opportunity) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	opp.UpdatedAt = time.Now()
	_, err = r.collection.ReplaceOne(ctx, bson.M{"_id": objID}, opp)
	return err
}
