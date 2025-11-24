package api

import (
	"github.com/gin-gonic/gin"
	"github.com/taiwanstay/taiwanstay-back/pkg/config"
)

// SetupRoutes 負責設定所有 API 路由
func SetupRoutes(router *gin.Engine, userHandler *UserHandler, imageHandler *ImageHandler, hostHandler *HostHandler, oppHandler *OpportunityHandler, cfg *config.Config) {
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
		users.Use(AuthMiddleware(cfg))
		users.Use(AdminAuthMiddleware())
		{
			users.GET("", userHandler.GetAllUsers)
			users.GET("/:id", userHandler.GetUserByID)
		}

		// 圖片相關路由
		images := v1.Group("/images")
		images.Use(AuthMiddleware(cfg))
		{
			images.POST("/upload", imageHandler.Upload)
			images.GET("/private/:id", imageHandler.GetPrivateImage)

			// Admin only
			adminImages := images.Group("/")
			adminImages.Use(AdminAuthMiddleware())
			{
				adminImages.PUT("/:id/status", imageHandler.UpdateStatus)
			}
		}

		// 接待主 (Host) 相關路由
		hosts := v1.Group("/hosts")
		hosts.Use(AuthMiddleware(cfg))
		{
			hosts.POST("", hostHandler.Create)
			hosts.GET("/me", hostHandler.GetMe)
			hosts.PUT("/me", hostHandler.UpdateMe)
		}

		// 機會 (Opportunity) 相關路由
		opps := v1.Group("/opportunities")
		{
			opps.GET("", oppHandler.List)
			opps.GET("/:id", oppHandler.GetByID)

			// 需要認證
			authOpps := opps.Group("")
			authOpps.Use(AuthMiddleware(cfg))
			{
				authOpps.POST("", oppHandler.Create)
			}
		}

		// ... 其他資源的路由設定

		// 當前登入者相關路由
		user := v1.Group("/user")
		user.Use(AuthMiddleware(cfg))
		{
			user.GET("/me", userHandler.GetMe)
			user.PUT("/me", userHandler.UpdateMe)
		}
	}
}
