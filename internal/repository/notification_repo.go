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

type NotificationRepository interface {
	Create(ctx context.Context, notification *domain.Notification) error
	ListByUserID(ctx context.Context, userID string, limit, offset int64) ([]*domain.Notification, int64, error)
	MarkAsRead(ctx context.Context, id string, userID string) error
	MarkAllAsRead(ctx context.Context, userID string) error
}

type mongoNotificationRepository struct {
	collection *mongo.Collection
}

func NewNotificationRepository(collection *mongo.Collection) NotificationRepository {
	// Create Indexes
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Indexes for filtering
	collection.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "userId", Value: 1}}},
		{Keys: bson.D{{Key: "userId", Value: 1}, {Key: "isRead", Value: 1}}},
	})

	return &mongoNotificationRepository{collection: collection}
}

func (r *mongoNotificationRepository) Create(ctx context.Context, notification *domain.Notification) error {
	if notification.ID.IsZero() {
		notification.ID = primitive.NewObjectID()
	}
	_, err := r.collection.InsertOne(ctx, notification)
	return err
}

func (r *mongoNotificationRepository) ListByUserID(ctx context.Context, userID string, limit, offset int64) ([]*domain.Notification, int64, error) {
	objID, _ := primitive.ObjectIDFromHex(userID)
	filter := bson.M{"userId": objID}

	// Count total
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	// Find with pagination and sort by CreatedAt desc
	opts := options.Find().
		SetSort(bson.D{{Key: "createdAt", Value: -1}}).
		SetSkip(offset).
		SetLimit(limit)

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var notifications []*domain.Notification
	if err = cursor.All(ctx, &notifications); err != nil {
		return nil, 0, err
	}

	return notifications, total, nil
}

func (r *mongoNotificationRepository) MarkAsRead(ctx context.Context, id string, userID string) error {
	notifID, _ := primitive.ObjectIDFromHex(id)
	userObjID, _ := primitive.ObjectIDFromHex(userID)

	filter := bson.M{"_id": notifID, "userId": userObjID}
	update := bson.M{"$set": bson.M{"isRead": true}}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

func (r *mongoNotificationRepository) MarkAllAsRead(ctx context.Context, userID string) error {
	userObjID, _ := primitive.ObjectIDFromHex(userID)
	filter := bson.M{"userId": userObjID, "isRead": false}
	update := bson.M{"$set": bson.M{"isRead": true}}

	_, err := r.collection.UpdateMany(ctx, filter, update)
	return err
}
