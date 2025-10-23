package main

import (
	"github.com/gin-gonic/gin"
	"github.com/taiwanstay/taiwanstay-back/internal/api"
	"github.com/taiwanstay/taiwanstay-back/internal/repository"
	"github.com/taiwanstay/taiwanstay-back/internal/service"
)

func main() {
	// ===================================================================
	// 依賴注入 (Dependency Injection)
	// ===================================================================

	// 1. 建立 Repository 層
	//    (目前是模擬的實作，未來會在這裡初始化資料庫連線並傳入)
	userRepository := repository.NewUserRepository()

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
