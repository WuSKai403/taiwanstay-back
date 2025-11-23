你現在是一位資深的 Go 後端工程師。接下來，我們將要一起合作完成一個名為 `taiwanstay-back` 的專案。

**[第一部分：專案核心指南]**

以下是這個專案的 **核心開發指南**，你必須嚴格遵守其中的所有架構、規範和程式碼風格。請先閱讀並理解這份指南。

```markdown
# TaiwanStay Go 後端開發指南

本文檔是 `taiwanstay-back` 專案的核心開發指南，整合了專案架構、API 遷移計畫、開發與測試流程、以及 AI 協作規範。所有開發者都應以此文件為準。

---

## 1. 專案架構 (Architecture)

### 1.1. 核心技術棧

| 類別 | 推薦函式庫 | 說明 |
| :--- | :--- | :--- |
| **Web 框架** | [Gin](https://github.com/gin-gonic/gin) | 一個高效能、API 精簡的 Web 框架。社群龐大，文件豐富，對新手友好。 |
| **資料庫驅動** | [Official MongoDB Driver](https://github.com/mongodb/mongo-go-driver) | 官方支援的 MongoDB 驅動，穩定可靠。 |
| **環境變數管理** | [Viper](https://github.com/spf13/viper) | 強大的設定檔管理工具，能輕鬆讀取 `.env` 檔案和 YAML。 |
| **請求驗證** | [Validator v10](https://github.com/go-playground/validator) | 基於 struct tag 的驗證函式庫，可以與 Gin 完美整合。 |
| **相依性管理** | Go Modules (內建) | Go 官方的相依性管理工具。 |

### 1.2. 專案結構 (Project Layout)

我們採用一個清晰的、分層的專案結構，以實現關注點分離。

/taiwanstay-back
├── cmd/
│   └── server/
│       └── main.go            # 程式進入點，初始化與啟動伺服器
├── internal/
│   ├── api/                   # HTTP Handlers (控制器)
│   ├── service/               # 業務邏輯層
│   ├── repository/            # 資料存取層
│   └── domain/                # 核心領域模型 (資料結構)
├── pkg/
│   ├── config/                # 設定檔讀取 (Viper)
│   └── database/              # 資料庫連線
├── .env.example               # 環境變數範例
├── go.mod                     # Go Modules 檔案
└── go.sum

### 1.3. 系統分層架構圖 (Layered Architecture)

此圖展示了系統的核心分層結構，以及各層之間的單向依賴關係。

graph TD
    subgraph "外部請求 (External Requests)"
        Client[用戶端 / Client]
    end

    subgraph "Go 後端應用 (Go Backend Application)"
        direction LR
        subgraph "接入層 (API Layer)"
            Router["路由 (Gin Router)"]
            Handlers[API Handlers]
            Middleware[中介軟體]
        end

        subgraph "業務邏輯層 (Service Layer)"
            Services[業務邏輯服務]
        end

        subgraph "資料存取層 (Repository Layer)"
            Repositories[資料倉儲]
        end

        subgraph "領域模型 (Domain Layer)"
            Models["資料模型 (Structs)"]
        end
    end

    subgraph "外部服務 (External Services)"
        MongoDB[(MongoDB)]
        GoogleVision[Google Vision API]
        Cloudinary[Cloudinary API]
    end

    Client --> Router
    Router --> Middleware
    Middleware --> Handlers
    Handlers --> Services
    Services --> Repositories
    Repositories --> Models
    Repositories --> MongoDB
    Services --> GoogleVision
    Services --> Cloudinary

---

## 2. API 遷移計畫

... (以下內容省略以保持簡潔，實際使用時應包含完整內容) ...

---

## 3. 開發與測試流程

### 3.1. 開發流程閉環

對於每一個新的 API 端點，都必須遵循以下**開發與測試的完整閉環**：

1.  **Domain 層**: 在 `internal/domain/` 中定義或確認資料模型。
2.  **Repository 層**: 在 `internal/repository/` 中定義介面並實作資料庫操作。
3.  **Service 層**: 在 `internal/service/` 中定義介面並實作業務邏輯。
4.  **Handler 層**: 在 `internal/api/` 中建立 Handler 方法處理 HTTP 請求。
5.  **Router 層**: 在 `internal/api/router.go` 中註冊路由。

... (以下內容省略以保持簡潔，實際使用時應包含完整內容) ...

---

## 4. 日誌與錯誤處理

... (以下內容省略以保持簡潔，實際使用時應包含完整內容) ...

---

## 5. AI 協作指南 (Continuation Prompt)

### 5.1. 核心原則

你 (AI) 必須嚴格遵守本文檔建立的架構、規範和開發模式。

### 5.2. 程式碼風格與規範

-   **顯式優於隱式**: 嚴格遵守手動依賴注入。
-   **介面導向**: 所有 `Repository` 和 `Service` 都必須先定義介面。
-   **錯誤處理**: 嚴格處理所有函式可能回傳的 `error`。
-   **安全性**: 密碼等敏感資訊必須加密處理。
-   **持續整合 (CI)**: 確保所有程式碼變更都能通過 CI 檢查。
```

**[第二部分：專案當前程式碼狀態]**

以下是專案目前的核心檔案內容。你生成的任何新程式碼都必須與這些現有的結構和介面相容。

```go
// FILE: cmd/server/main.go
package main

import (
	"context"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/taiwanstay/taiwanstay-back/internal/api"
	"github.com/taiwanstay/taiwanstay-back/internal/repository"
	"github.com/taiwanstay/taiwanstay-back/internal/service"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// ===================================================================
	// 資料庫連線
	// ===================================================================
	mongoURI := "mongodb://localhost:27017"
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Failed to connect to mongo: %v", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to ping mongo: %v", err)
	}

	log.Println("Successfully connected to MongoDB!")

	userCollection := client.Database("taiwanstay").Collection("users")

	// ===================================================================
	// 依賴注入 (Dependency Injection)
	// ===================================================================

	userRepository := repository.NewUserRepository(userCollection)
	userService := service.NewUserService(userRepository)
	userHandler := api.NewUserHandler(userService)

	// ===================================================================
	// 伺服器設定
	// ===================================================================

	router := gin.Default()
	api.SetupRoutes(router, userHandler)

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
```

```go
// FILE: internal/domain/user.go
package domain

import (
	"time"
)

type UserRole string
type PrivacyLevel string

const (
	RoleUser  UserRole = "USER"
	RoleHost  UserRole = "HOST"
	RoleAdmin UserRole = "ADMIN"

	PrivacyPublic     PrivacyLevel = "PUBLIC"
	PrivacyRegistered PrivacyLevel = "REGISTERED"
	PrivacyPrivate    PrivacyLevel = "PRIVATE"
)

type User struct {
	ID              string          `json:"id" bson:"_id,omitempty"`
	Name            string          `json:"name" bson:"name"`
	Email           string          `json:"email" bson:"email"`
	Image           string          `json:"image,omitempty" bson:"image,omitempty"`
	EmailVerified   *time.Time      `json:"emailVerified,omitempty" bson:"emailVerified,omitempty"`
	Password        string          `json:"-" bson:"password,omitempty"`
	Role            UserRole        `json:"role" bson:"role"`
	Profile         Profile         `json:"profile" bson:"profile"`
	HostID          string          `json:"hostId,omitempty" bson:"hostId,omitempty"`
	OrganizationID  string          `json:"organizationId,omitempty" bson:"organizationId,omitempty"`
	PrivacySettings PrivacySettings `json:"privacySettings" bson:"privacySettings"`
	CreatedAt       time.Time       `json:"createdAt" bson:"createdAt"`
	UpdatedAt       time.Time       `json:"updatedAt" bson:"updatedAt"`
}

type Profile struct {
	Avatar                  string                   `json:"avatar,omitempty" bson:"avatar,omitempty"`
	Bio                     string                   `json:"bio,omitempty" bson:"bio,omitempty"`
	Skills                  []string                 `json:"skills,omitempty" bson:"skills,omitempty"`
	Languages               []string                 `json:"languages,omitempty" bson:"languages,omitempty"`
	Location                *Location                `json:"location,omitempty" bson:"location,omitempty"`
	SocialMedia             *SocialMedia             `json:"socialMedia,omitempty" bson:"socialMedia,omitempty"`
	PersonalInfo            *PersonalInfo            `json:"personalInfo,omitempty" bson:"personalInfo,omitempty"`
	WorkExchangePreferences *WorkExchangePreferences `json:"workExchangePreferences,omitempty" bson:"workExchangePreferences,omitempty"`
	BirthDate               time.Time                `json:"birthDate" bson:"birthDate"`
	EmergencyContact        EmergencyContact         `json:"emergencyContact" bson:"emergencyContact"`
	WorkExperience          []WorkExperience         `json:"workExperience" bson:"workExperience"`
	PhysicalCondition       string                   `json:"physicalCondition" bson:"physicalCondition"`
	AccommodationNeeds      string                   `json:"accommodationNeeds,omitempty" bson:"accommodationNeeds,omitempty"`
	CulturalInterests       []string                 `json:"culturalInterests,omitempty" bson:"culturalInterests,omitempty"`
	LearningGoals           []string                 `json:"learningGoals,omitempty" bson:"learningGoals,omitempty"`
	PhoneNumber             string                   `json:"phoneNumber,omitempty" bson:"phoneNumber,omitempty"`
	Address                 string                   `json:"address,omitempty" bson:"address,omitempty"`
	IsPhoneVerified         bool                     `json:"isPhoneVerified" bson:"isPhoneVerified"`
	PreferredWorkHours      int                      `json:"preferredWorkHours" bson:"preferredWorkHours"`
}

// ... (其他 sub-structs 省略，實際使用時應包含) ...
```

```go
// FILE: internal/repository/user_repo.go
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
	GetAll(ctx context.Context) ([]*domain.User, error)
	GetByID(ctx context.Context, id string) (*domain.User, error)
	Update(ctx context.Context, id string, payload bson.M) error
}

// mongoUserRepository 是 UserRepository 的 MongoDB 實作
type mongoUserRepository struct {
	collection *mongo.Collection
}

// NewUserRepository 建立一個新的 UserRepository 實例
func NewUserRepository(collection *mongo.Collection) UserRepository {
	return &mongoUserRepository{collection: collection}
}

// ... (實作方法省略，實際使用時應包含) ...
```

```go
// FILE: internal/service/user_service.go
package service

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/taiwanstay/taiwanstay-back/internal/domain"
	"github.com/taiwanstay/taiwanstay-back/internal/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid email or password")
)

// UserService 定義了與使用者相關的業務邏輯介面
type UserService interface {
	RegisterUser(ctx context.Context, name, email, password string) (*domain.User, error)
	LoginUser(ctx context.Context, email, password string) (*domain.User, string, error)
	LogoutUser(ctx context.Context) error
	GetAllUsers(ctx context.Context) ([]*domain.User, error)
	GetUserByID(ctx context.Context, id string) (*domain.User, error)
	UpdateUser(ctx context.Context, id string, payload bson.M) (*domain.User, error)
}

// userService 是 UserService 的實作
type userService struct {
	userRepo      repository.UserRepository
	jwtSecret     string
	tokenDuration time.Duration
}

// NewUserService 建立一個新的 UserService 實例
func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		userRepo:      repo,
		jwtSecret:     "your-super-secret-key", // TODO: Move to config
		tokenDuration: time.Hour * 24,
	}
}

// ... (實作方法省略，實際使用時應包含) ...
```

```go
// FILE: internal/api/user_handler.go
package api

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/taiwanstay/taiwanstay-back/internal/domain"
	"github.com/taiwanstay/taiwanstay-back/internal/service"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// UserHandler 負責處理與使用者相關的 HTTP 請求
type UserHandler struct {
	userService service.UserService
}

// NewUserHandler 建立一個新的 UserHandler 實例
func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// ... (DTOs 和 Handler 方法省略，實際使用時應包含) ...
```

```go
// FILE: internal/api/router.go
package api

import "github.com/gin-gonic/gin"

// SetupRoutes 負責設定所有 API 路由
func SetupRoutes(router *gin.Engine, userHandler *UserHandler) {
	v1 := router.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/register", userHandler.Register)
			auth.POST("/login", userHandler.Login)
			auth.POST("/logout", userHandler.Logout)
		}

		users := v1.Group("/users")
		users.Use(AuthMiddleware())
		users.Use(AdminAuthMiddleware())
		{
			users.GET("/", userHandler.GetAllUsers)
			users.GET("/:id", userHandler.GetUserByID)
		}

		user := v1.Group("/user")
		user.Use(AuthMiddleware())
		{
			user.GET("/me", userHandler.GetMe)
			user.PUT("/me", userHandler.UpdateMe)
		}
	}
}
```

**[第三部分：你的任務]**

現在，請基於以上的指南和程式碼，開始我們的下一個任務。

(在這裡貼上你的「功能切片」請求，例如：「現在我們要來實作『工作機會』的功能...」)
