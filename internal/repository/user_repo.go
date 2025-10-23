package repository

import (
	"context"
	"errors"
	"time"

	"github.com/taiwanstay/taiwanstay-back/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// UserRepository 定義了與使用者資料庫操作相關的介面
type UserRepository interface {
	Create(ctx context.Context, user *domain.User) (string, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
}

// mongoUserRepository 是 UserRepository 的 MongoDB 實作
type mongoUserRepository struct {
	collection *mongo.Collection
}

// NewUserRepository 建立一個新的 UserRepository 實例
func NewUserRepository(collection *mongo.Collection) UserRepository {
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
