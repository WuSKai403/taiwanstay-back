package repository

import (
	"context"
	"errors"
	"time"

	"github.com/taiwanstay/taiwanstay-back/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserRepository 定義了與使用者資料庫操作相關的介面
type UserRepository interface {
	Create(ctx context.Context, user *domain.User) (string, error)
	// GetByEmail(ctx context.Context, email string) (*domain.User, error)
	// GetByID(ctx context.Context, id string) (*domain.User, error)
}

// mongoUserRepository 是 UserRepository 的 MongoDB 實作
type mongoUserRepository struct {
	// TODO: 在這裡加入 MongoDB collection 的參考
	// collection *mongo.Collection
}

// NewUserRepository 建立一個新的 UserRepository 實例
func NewUserRepository() UserRepository {
	// TODO: 在這裡傳入 MongoDB collection
	return &mongoUserRepository{}
}

// Create 在資料庫中建立一個新使用者
func (r *mongoUserRepository) Create(ctx context.Context, user *domain.User) (string, error) {
	// --- 這是模擬的資料庫操作 ---
	// 在未來的實作中，這裡將會是實際的 MongoDB 插入邏輯

	// 1. 檢查 Email 是否重複 (模擬)
	if user.Email == "exists@example.com" {
		return "", errors.New("email already exists")
	}

	// 2. 設定時間戳和 ID (模擬)
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	// 在真實情境中，MongoDB 會自動生成 _id
	generatedID := primitive.NewObjectID().Hex()
	user.ID = generatedID

	// 3. 模擬插入資料庫
	// collection.InsertOne(ctx, user)
	// log.Printf("User created with ID: %s", generatedID)

	// --- 模擬結束 ---

	// 返回生成的 ID
	return generatedID, nil
}
