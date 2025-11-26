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
	Delete(ctx context.Context, id string) error
	Search(ctx context.Context, filter OpportunityFilter) ([]*domain.Opportunity, int64, error)
}

type OpportunityFilter struct {
	Query     string
	Type      string
	City      string
	Country   string
	StartDate string // YYYY-MM-DD
	EndDate   string // YYYY-MM-DD
	Lat       float64
	Lng       float64
	Distance  float64 // in meters
	Limit     int64
	Offset    int64
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

	// Text index for search
	collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{Key: "title", Value: "text"},
			{Key: "description", Value: "text"},
			{Key: "shortDescription", Value: "text"},
			{Key: "location.city", Value: "text"},
			{Key: "location.country", Value: "text"},
		},
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

func (r *mongoOpportunityRepository) Delete(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	update := bson.M{
		"$set": bson.M{
			"status":    domain.OpportunityStatusDeleted,
			"updatedAt": time.Now(),
		},
	}
	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	return err
}

func (r *mongoOpportunityRepository) Search(ctx context.Context, filter OpportunityFilter) ([]*domain.Opportunity, int64, error) {
	query := bson.M{"status": domain.OpportunityStatusActive}

	// Text Search
	if filter.Query != "" {
		query["$text"] = bson.M{"$search": filter.Query}
	}

	// Filters
	if filter.Type != "" {
		query["type"] = filter.Type
	}
	if filter.City != "" {
		query["location.city"] = filter.City
	}
	if filter.Country != "" {
		query["location.country"] = filter.Country
	}

	// Geospatial Search
	if filter.Lat != 0 && filter.Lng != 0 {
		maxDist := filter.Distance
		if maxDist == 0 {
			maxDist = 50000 // Default 50km
		}
		query["location.coordinates"] = bson.M{
			"$near": bson.M{
				"$geometry": bson.M{
					"type":        "Point",
					"coordinates": []float64{filter.Lng, filter.Lat},
				},
				"$maxDistance": maxDist,
			},
		}
	}

	// Time Slot Logic (Range Overlap)
	if filter.StartDate != "" && filter.EndDate != "" {
		// Find opportunities that have at least one time slot overlapping with the requested range
		// AND that time slot is OPEN
		query["timeSlots"] = bson.M{
			"$elemMatch": bson.M{
				"status":    domain.TimeSlotStatusOpen,
				"startDate": bson.M{"$lte": filter.EndDate},
				"endDate":   bson.M{"$gte": filter.StartDate},
			},
		}
	}

	// Count total
	total, err := r.collection.CountDocuments(ctx, query)
	if err != nil {
		return nil, 0, err
	}

	// Find
	opts := options.Find().SetLimit(filter.Limit).SetSkip(filter.Offset)

	// If text search, sort by score
	if filter.Query != "" {
		opts.SetProjection(bson.M{"score": bson.M{"$meta": "textScore"}})
		opts.SetSort(bson.M{"score": bson.M{"$meta": "textScore"}})
	} else {
		opts.SetSort(bson.D{{Key: "createdAt", Value: -1}})
	}

	cursor, err := r.collection.Find(ctx, query, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var opps []*domain.Opportunity
	if err := cursor.All(ctx, &opps); err != nil {
		return nil, 0, err
	}

	return opps, total, nil
}
