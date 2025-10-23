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
	// 在未來，這個連線字串應該來自設定檔或環境變數
	mongoURI := "mongodb://localhost:27017"
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Failed to create mongo client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatalf("Failed to connect to mongo: %v", err)
	}

	// 檢查連線
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to ping mongo: %v", err)
	}

	log.Println("Successfully connected to MongoDB!")

	// 取得 user collection
	userCollection := client.Database("taiwanstay").Collection("users")

	// ===================================================================
	// 依賴注入 (Dependency Injection)
	// ===================================================================

	// 1. 建立 Repository 層
	userRepository := repository.NewUserRepository(userCollection)

	// 2. 建立 Service 層，並注入 Repository
	userService := service.NewUserService(userRepository)

	// 3. 建立 Handler 層，並注入 Service
	userHandler := api.NewUserHandler(userService)

	// ===================================================================
	// 伺服器設定
	// ===================================================================

	// 初始化 Gin 引擎
	router := gin.Default()

	// 設定所有 API 路由，並傳入 Handler
	api.SetupRoutes(router, userHandler)

	// 啟動 HTTP 伺服器，監聽在 8080 port
	router.Run(":8080")
}
