package api

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"

	"github.com/taiwanstay/taiwanstay-back/internal/domain"
	"github.com/taiwanstay/taiwanstay-back/internal/repository"
	"github.com/taiwanstay/taiwanstay-back/internal/service"
	"github.com/taiwanstay/taiwanstay-back/pkg/config"
)

var (
	testCollection *mongo.Collection
	testRouter     *gin.Engine
	testConfig     *config.Config
)

// TestMain 是測試的主進入點，用於設定和清理測試環境
func TestMain(m *testing.M) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	// Setup test config
	testConfig = &config.Config{
		Server: config.ServerConfig{
			JWTSecret: "test-secret",
			Mode:      "test",
		},
	}

	// 啟動 MongoDB 容器
	mongodbContainer, err := mongodb.Run(ctx, "mongo:6")
	if err != nil {
		log.Fatalf("failed to start container: %s", err)
	}

	// 取得連線 URI
	uri, err := mongodbContainer.ConnectionString(ctx)
	if err != nil {
		log.Fatalf("failed to get connection string: %s", err)
	}

	// 連線到 MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("failed to connect to mongo: %s", err)
	}

	// 設定全域變數
	testCollection = client.Database("testdb").Collection("users")
	testRouter = setupTestRouter(testCollection)

	// 執行所有測試
	exitCode := m.Run()

	// 清理容器
	if err := mongodbContainer.Terminate(ctx); err != nil {
		log.Fatalf("failed to terminate container: %s", err)
	}

	os.Exit(exitCode)
}

// setupTestRouter 設定一個用於測試的 Gin 路由器
func setupTestRouter(collection *mongo.Collection) *gin.Engine {
	gin.SetMode(gin.TestMode)

	userRepo := repository.NewUserRepository(collection)
	userService := service.NewUserService(userRepo, testConfig)
	userHandler := NewUserHandler(userService)

	router := gin.Default()
	// Pass nil for ImageHandler, HostHandler, OppHandler, AppHandler as we are not testing them here yet
	SetupRoutes(router, userHandler, nil, nil, nil, nil, testConfig)
	return router
}

// cleanupCollection 在每個測試前清理集合
func cleanupCollection(ctx context.Context) {
	if err := testCollection.Drop(ctx); err != nil {
		log.Fatalf("failed to drop collection: %s", err)
	}
}

func TestRegister_Success(t *testing.T) {
	ctx := context.Background()
	cleanupCollection(ctx)

	// 準備請求 Body
	registerData := gin.H{
		"name":     "testuser",
		"email":    "test@example.com",
		"password": "password123",
	}
	body, _ := json.Marshal(registerData)

	// 建立請求
	req, _ := http.NewRequestWithContext(ctx, "POST", "/api/v1/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// 執行請求
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	// 斷言
	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.NotEmpty(t, response["id"])
	assert.Equal(t, "testuser", response["name"])
	assert.Equal(t, "test@example.com", response["email"])
	assert.Nil(t, response["password"])
}

func TestRegister_EmailAlreadyExists(t *testing.T) {
	ctx := context.Background()
	cleanupCollection(ctx)

	// 1. 先建立一個使用者
	existingUser := domain.User{Name: "existing", Email: "exist@example.com", Password: "password"}
	_, err := testCollection.InsertOne(ctx, &existingUser)
	assert.NoError(t, err)

	// 2. 嘗試用相同的 email 註冊
	registerData := gin.H{
		"name":     "anotheruser",
		"email":    "exist@example.com",
		"password": "anotherpassword",
	}
	body, _ := json.Marshal(registerData)

	req, _ := http.NewRequestWithContext(ctx, "POST", "/api/v1/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	// 斷言
	assert.Equal(t, http.StatusConflict, w.Code)

	var response map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "email already exists", response["error"])
}

func TestRegister_MissingFields(t *testing.T) {
	ctx := context.Background()
	cleanupCollection(ctx)

	// 準備一個缺少 password 的請求 Body
	registerData := gin.H{
		"name":  "testuser",
		"email": "test@example.com",
	}
	body, _ := json.Marshal(registerData)

	req, _ := http.NewRequestWithContext(ctx, "POST", "/api/v1/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	// 斷言
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestLogin_Success(t *testing.T) {
	ctx := context.Background()
	cleanupCollection(ctx)

	// 1. 先註冊一個使用者
	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	testUser := domain.User{
		Name:     "testuser",
		Email:    "test@example.com",
		Password: string(hashedPassword),
	}
	_, err := testCollection.InsertOne(ctx, &testUser)
	assert.NoError(t, err)

	// 2. 準備登入請求
	loginData := gin.H{
		"loginType": "password",
		"email":     "test@example.com",
		"password":  "password123",
	}
	body, _ := json.Marshal(loginData)

	req, _ := http.NewRequestWithContext(ctx, "POST", "/api/v1/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	// 3. 斷言
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response["token"])
	assert.NotNil(t, response["user"])
}

func TestLogin_WrongPassword(t *testing.T) {
	ctx := context.Background()
	cleanupCollection(ctx)

	// 1. 先註冊一個使用者
	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	testUser := domain.User{
		Name:     "testuser",
		Email:    "test@example.com",
		Password: string(hashedPassword),
	}
	_, err := testCollection.InsertOne(ctx, &testUser)
	assert.NoError(t, err)

	// 2. 準備登入請求 (密碼錯誤)
	loginData := gin.H{
		"loginType": "password",
		"email":     "test@example.com",
		"password":  "wrongpassword",
	}
	body, _ := json.Marshal(loginData)

	req, _ := http.NewRequestWithContext(ctx, "POST", "/api/v1/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	// 3. 斷言
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestLogin_UserNotFound(t *testing.T) {
	ctx := context.Background()
	cleanupCollection(ctx)

	// 準備登入請求 (使用者不存在)
	loginData := gin.H{
		"loginType": "password",
		"email":     "notfound@example.com",
		"password":  "password123",
	}
	body, _ := json.Marshal(loginData)

	req, _ := http.NewRequestWithContext(ctx, "POST", "/api/v1/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	// 斷言
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// generateTestToken 根據使用者資訊生成用於測試的 JWT
func generateTestToken(t *testing.T, user *domain.User) string {
	jwtSecret := testConfig.Server.JWTSecret // 與 middleware 中保持一致
	tokenDuration := time.Hour * 24

	claims := jwt.MapClaims{
		"sub":  user.ID,
		"name": user.Name,
		"role": user.Role,
		"exp":  time.Now().Add(tokenDuration).Unix(),
		"iat":  time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(jwtSecret))
	assert.NoError(t, err)
	return signedToken
}

// createAndLoginUser 是一個輔助函式，用於在資料庫中建立指定角色的使用者，並返回 userID 和 token
func createAndLoginUser(t *testing.T, ctx context.Context, name, email, password string, role domain.UserRole) (*domain.User, string) {
	// 1. 直接在資料庫中建立使用者
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := &domain.User{
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
		Role:     role,
	}
	res, err := testCollection.InsertOne(ctx, user)
	assert.NoError(t, err)
	userID := res.InsertedID.(primitive.ObjectID)
	user.ID = userID.Hex()

	// 2. 產生 token
	token := generateTestToken(t, user)

	return user, token
}

func TestGetMe_Success(t *testing.T) {
	ctx := context.Background()
	cleanupCollection(ctx)

	// 1. 註冊並登入使用者
	_, token := createAndLoginUser(t, ctx, "me_user", "me@example.com", "password123", domain.RoleUser)

	// 2. 建立 GetMe 請求
	req, _ := http.NewRequestWithContext(ctx, "GET", "/api/v1/user/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	// 3. 斷言
	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "me@example.com", response["email"])
}

func TestUpdateMe_Success(t *testing.T) {
	ctx := context.Background()
	cleanupCollection(ctx)

	// 1. 註冊並登入使用者
	user, token := createAndLoginUser(t, ctx, "update_user", "update@example.com", "password123", domain.RoleUser)
	userID := user.ID

	// 2. 準備更新資料
	updateData := gin.H{
		"name":              "Updated Name",
		"physicalCondition": "Excellent",
		"emergencyContact": gin.H{
			"name":         "Jane Doe",
			"relationship": "Spouse",
			"phone":        "123456789",
		},
	}
	body, _ := json.Marshal(updateData)

	// 3. 建立 UpdateMe 請求
	req, _ := http.NewRequestWithContext(ctx, "PUT", "/api/v1/user/me", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	// 4. 斷言 HTTP 回應
	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Name", response["name"])

	// 5. 直接從資料庫驗證資料是否真的被更新
	var updatedUser domain.User
	objID, _ := primitive.ObjectIDFromHex(userID)
	err = testCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&updatedUser)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Name", updatedUser.Name)
	assert.Equal(t, "Excellent", updatedUser.Profile.PhysicalCondition)
	assert.Equal(t, "Jane Doe", updatedUser.Profile.EmergencyContact.Name)
}

func TestLogout_Success(t *testing.T) {
	ctx := context.Background()
	cleanupCollection(ctx)

	// 1. 註冊並登入使用者
	_, token := createAndLoginUser(t, ctx, "logout_user", "logout@example.com", "password123", domain.RoleUser)

	// 2. 建立 Logout 請求
	req, _ := http.NewRequestWithContext(ctx, "POST", "/api/v1/auth/logout", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	// 3. 斷言
	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "logout successful", response["message"])
}

func TestAdminActions_SuccessAsAdmin(t *testing.T) {
	ctx := context.Background()
	cleanupCollection(ctx)

	// 1. 建立一個管理員和一個普通使用者
	adminUser, adminToken := createAndLoginUser(t, ctx, "admin", "admin@example.com", "password123", domain.RoleAdmin)
	normalUser, _ := createAndLoginUser(t, ctx, "user", "user@example.com", "password123", domain.RoleUser)
	assert.NotEqual(t, adminUser.ID, normalUser.ID)

	// 2. 作為管理員，嘗試獲取所有使用者
	reqGetAll, _ := http.NewRequestWithContext(ctx, "GET", "/api/v1/users/", nil)
	reqGetAll.Header.Set("Authorization", "Bearer "+adminToken)
	wGetAll := httptest.NewRecorder()
	testRouter.ServeHTTP(wGetAll, reqGetAll)

	// 斷言 GetAllUsers
	assert.Equal(t, http.StatusOK, wGetAll.Code)
	var users []map[string]interface{}
	err := json.Unmarshal(wGetAll.Body.Bytes(), &users)
	assert.NoError(t, err)
	assert.Len(t, users, 2) // 應該包含 admin 和 user

	// 3. 作為管理員，嘗試透過 ID 獲取特定使用者
	reqGetByID, _ := http.NewRequestWithContext(ctx, "GET", "/api/v1/users/"+normalUser.ID, nil)
	reqGetByID.Header.Set("Authorization", "Bearer "+adminToken)
	wGetByID := httptest.NewRecorder()
	testRouter.ServeHTTP(wGetByID, reqGetByID)

	// 斷言 GetUserByID
	assert.Equal(t, http.StatusOK, wGetByID.Code)
	var fetchedUser map[string]interface{}
	err = json.Unmarshal(wGetByID.Body.Bytes(), &fetchedUser)
	assert.NoError(t, err)
	assert.Equal(t, normalUser.Email, fetchedUser["email"])
}

func TestAdminActions_ForbiddenAsUser(t *testing.T) {
	ctx := context.Background()
	cleanupCollection(ctx)

	// 1. 建立兩個普通使用者
	_, user1Token := createAndLoginUser(t, ctx, "user1", "user1@example.com", "password123", domain.RoleUser)
	user2, _ := createAndLoginUser(t, ctx, "user2", "user2@example.com", "password123", domain.RoleUser)

	// 2. 作為普通使用者，嘗試獲取所有使用者
	reqGetAll, _ := http.NewRequestWithContext(ctx, "GET", "/api/v1/users/", nil)
	reqGetAll.Header.Set("Authorization", "Bearer "+user1Token)
	wGetAll := httptest.NewRecorder()
	testRouter.ServeHTTP(wGetAll, reqGetAll)

	// 斷言 GetAllUsers 失敗
	assert.Equal(t, http.StatusForbidden, wGetAll.Code)

	// 3. 作為普通使用者，嘗試透過 ID 獲取另一個使用者
	reqGetByID, _ := http.NewRequestWithContext(ctx, "GET", "/api/v1/users/"+user2.ID, nil)
	reqGetByID.Header.Set("Authorization", "Bearer "+user1Token)
	wGetByID := httptest.NewRecorder()
	testRouter.ServeHTTP(wGetByID, reqGetByID)

	// 斷言 GetUserByID 失敗
	assert.Equal(t, http.StatusForbidden, wGetByID.Code)
}
