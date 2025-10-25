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
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"

	"github.com/taiwanstay/taiwanstay-back/internal/domain"
	"github.com/taiwanstay/taiwanstay-back/internal/repository"
	"github.com/taiwanstay/taiwanstay-back/internal/service"
)

var (
	testCollection *mongo.Collection
	testRouter     *gin.Engine
)

// TestMain 是測試的主進入點，用於設定和清理測試環境
func TestMain(m *testing.M) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

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
	userService := service.NewUserService(userRepo)
	userHandler := NewUserHandler(userService)

	router := gin.Default()
	SetupRoutes(router, userHandler)
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
