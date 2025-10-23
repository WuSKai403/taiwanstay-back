package service

import (
	"context"
	"errors"
	"time"

	"github.com/taiwanstay/taiwanstay-back/internal/domain"
	"github.com/taiwanstay/taiwanstay-back/internal/repository"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// ErrEmailAlreadyExists 表示 email 已被註冊
var ErrEmailAlreadyExists = errors.New("email already exists")

// UserService 定義了與使用者相關的業務邏輯介面
type UserService interface {
	RegisterUser(ctx context.Context, name, email, password string) (*domain.User, error)
}

// userService 是 UserService 的實作
type userService struct {
	userRepo repository.UserRepository
}

// NewUserService 建立一個新的 UserService 實例
func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		userRepo: repo,
	}
}

// RegisterUser 處理使用者註冊的業務邏輯
func (s *userService) RegisterUser(ctx context.Context, name, email, password string) (*domain.User, error) {
	// 1. 檢查 Email 是否已存在
	_, err := s.userRepo.GetByEmail(ctx, email)
	if err == nil {
		// 找到了使用者，表示 email 已存在
		return nil, ErrEmailAlreadyExists
	}
	if !errors.Is(err, mongo.ErrNoDocuments) {
		// 如果是除了 "not found" 以外的其他錯誤，就回傳
		return nil, err
	}

	// 2. 密碼加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 3. 建立 User domain 物件
	newUser := &domain.User{
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
		Role:     domain.RoleUser, // 預設角色為 USER
		Profile: domain.Profile{ // 初始化一些必要的 Profile 欄位
			BirthDate:          time.Now(), // 暫時使用當前時間，未來應由前端傳入
			EmergencyContact:   domain.EmergencyContact{Name: "N/A", Relationship: "N/A", Phone: "N/A"},
			PhysicalCondition:  "N/A",
			PreferredWorkHours: 8,
		},
		PrivacySettings: domain.PrivacySettings{ // 設定預設隱私等級
			Email: domain.PrivacyPrivate,
			Phone: domain.PrivacyPrivate,
			// ... 其他隱私設定
		},
	}

	// 3. 呼叫 Repository 將使用者存入資料庫
	userID, err := s.userRepo.Create(ctx, newUser)
	if err != nil {
		return nil, err
	}

	// 4. 設定回傳物件的 ID，並清除密碼，確保密碼不會外洩
	newUser.ID = userID
	newUser.Password = ""

	return newUser, nil
}
