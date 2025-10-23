package api

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/taiwanstay/taiwanstay-back/internal/domain"
	"github.com/taiwanstay/taiwanstay-back/internal/repository"
	"github.com/taiwanstay/taiwanstay-back/internal/service"
)

// setupTestMongoDB 啟動一個用於測試的 MongoDB Docker 容器
func setupTestMongoDB(ctx context.Context) (*mongo.Collection, func()) {
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

	// 取得 collection
	collection := client.Database("testdb").Collection("users")

	// 定義清理函式
	cleanup := func() {
		if err := mongodbContainer.Terminate(ctx); err != nil {
			log.Fatalf("failed to terminate container: %s", err)
		}
	}

	return collection, cleanup
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

func TestRegister_Success(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	collection, cleanup := setupTestMongoDB(ctx)
	defer cleanup()

	router := setupTestRouter(collection)

	// 準備請求 Body
	registerData := gin.H{
		"name":     "testuser",
		"email":    "test@example.com",
		"password": "password123",
	}
	body, _ := json.Marshal(registerData)

	// 建立請求
	req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// 執行請求
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 斷言
	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.NotEmpty(t, response["id"])
	assert.Equal(t, "testuser", response["name"])
	assert.Equal(t, "test@example.com", response["email"])
	assert.Nil(t, response["password"])
}

func TestRegister_EmailAlreadyExists(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	collection, cleanup := setupTestMongoDB(ctx)
	defer cleanup()

	router := setupTestRouter(collection)

	// 1. 先建立一個使用者
	existingUser := domain.User{Name: "existing", Email: "exist@example.com", Password: "password"}
	_, err := collection.InsertOne(ctx, &existingUser)
	assert.NoError(t, err)

	// 2. 嘗試用相同的 email 註冊
	registerData := gin.H{
		"name":     "anotheruser",
		"email":    "exist@example.com",
		"password": "anotherpassword",
	}
	body, _ := json.Marshal(registerData)

	req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 斷言
	assert.Equal(t, http.StatusConflict, w.Code)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "email already exists", response["error"])
}

func TestRegister_MissingFields(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	collection, cleanup := setupTestMongoDB(ctx)
	defer cleanup()

	router := setupTestRouter(collection)

	// 準備一個缺少 password 的請求 Body
	registerData := gin.H{
		"name":  "testuser",
		"email": "test@example.com",
	}
	body, _ := json.Marshal(registerData)

	req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 斷言
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
