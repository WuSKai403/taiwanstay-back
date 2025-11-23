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
			auth.POST("/login", userHandler.Login)
			auth.POST("/logout", userHandler.Logout)
		}

		// 用戶相關路由 (需要管理員權限)
		users := v1.Group("/users")
		users.Use(AuthMiddleware())
		users.Use(AdminAuthMiddleware())
		{
			users.GET("/", userHandler.GetAllUsers)
			users.GET("/:id", userHandler.GetUserByID)
		}

		// ... 其他資源的路由設定

		// 當前登入者相關路由
		user := v1.Group("/user")
		user.Use(AuthMiddleware())
		{
			user.GET("/me", userHandler.GetMe)
			user.PUT("/me", userHandler.UpdateMe)
		}
	}
}
