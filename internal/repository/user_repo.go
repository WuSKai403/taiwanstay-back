package repository

import (
	"context"
	"errors"
	"time"

	"github.com/taiwanstay/taiwanstay-back/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// UserRepository 定義了與使用者資料庫操作相關的介面
type UserRepository interface {
	Create(ctx context.Context, user *domain.User) (string, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetAll(ctx context.Context) ([]*domain.User, error)
	GetByID(ctx context.Context, id string) (*domain.User, error)
	Update(ctx context.Context, id string, payload bson.M) error
	Count(ctx context.Context) (int64, error)
	List(ctx context.Context, filter bson.M, limit, offset int64) ([]*domain.User, int64, error)
	UpdateStatus(ctx context.Context, id string, status domain.UserStatus) error
}

// mongoUserRepository 是 UserRepository 的 MongoDB 實作
type mongoUserRepository struct {
	collection *mongo.Collection
}

// NewUserRepository 建立一個新的 UserRepository 實例
func NewUserRepository(collection *mongo.Collection) UserRepository {
	// Create Indexes
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Unique index for email
	collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true),
	})

	return &mongoUserRepository{collection: collection}
}

// Create 在資料庫中建立一個新使用者
func (r *mongoUserRepository) Create(ctx context.Context, user *domain.User) (string, error) {
	// 設定時間戳
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	// 插入資料庫
	res, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		return "", err
	}

	// 返回生成的 ID
	return res.InsertedID.(primitive.ObjectID).Hex(), nil
}

// GetByEmail 透過 Email 尋找使用者
func (r *mongoUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, mongo.ErrNoDocuments
		}
		return nil, err
	}
	return &user, nil
}

// Update 更新資料庫中的使用者資訊
func (r *mongoUserRepository) Update(ctx context.Context, id string, payload bson.M) error {
	// 將 ID 字串轉換為 ObjectID
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid user id format")
	}

	// 確保更新時間戳被設定
	payload["updatedAt"] = time.Now()

	// 建立更新的 filter 和 update document
	filter := bson.M{"_id": objID}
	update := bson.M{"$set": payload}

	// 執行更新操作
	res, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	// 檢查是否有文件被更新
	if res.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

// GetAll 從資料庫中取得所有使用者
func (r *mongoUserRepository) GetAll(ctx context.Context) ([]*domain.User, error) {
	var users []*domain.User

	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
}

// Count 計算使用者總數
func (r *mongoUserRepository) Count(ctx context.Context) (int64, error) {
	return r.collection.CountDocuments(ctx, bson.M{})
}

// List 根據條件列出使用者
func (r *mongoUserRepository) List(ctx context.Context, filter bson.M, limit, offset int64) ([]*domain.User, int64, error) {
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

	var users []*domain.User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// GetByID 透過 ID 尋找使用者
func (r *mongoUserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	var user domain.User
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid user id format")
	}

	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, mongo.ErrNoDocuments
		}
		return nil, err
	}
	return &user, nil
}

// UpdateStatus 更新使用者狀態
func (r *mongoUserRepository) UpdateStatus(ctx context.Context, id string, status domain.UserStatus) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid user id format")
	}

	filter := bson.M{"_id": objID}
	update := bson.M{"$set": bson.M{"status": status, "updatedAt": time.Now()}}

	res, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}
