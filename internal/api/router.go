package api

import "github.com/gin-gonic/gin"

// SetupRoutes 負責設定所有 API 路由
func SetupRoutes(router *gin.Engine, userHandler *UserHandler) {
	// 建立 API 版本分組
	v1 := router.Group("/api/v1")
	{
		// 認證相關路由
		auth := v1.Group("/auth")
		{
			auth.POST("/register", userHandler.Register)
			// auth.POST("/login", userHandler.Login) // 未來實作
		}

		// 用戶相關路由
		// users := v1.Group("/users")
		// {
		// 	users.GET("/", GetUsersHandler)
		// }

		// ... 其他資源的路由設定
	}
}
