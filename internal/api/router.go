package api

import (
	"github.com/gin-gonic/gin"
	"github.com/taiwanstay/taiwanstay-back/pkg/config"
)

// SetupRoutes 負責設定所有 API 路由
func SetupRoutes(router *gin.Engine, userHandler *UserHandler, imageHandler *ImageHandler, hostHandler *HostHandler, oppHandler *OpportunityHandler, appHandler *ApplicationHandler, notifHandler *NotificationHandler, adminHandler *AdminHandler, cfg *config.Config) {
	// Global Middleware
	router.Use(gin.Recovery())
	router.Use(Logger())

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
			opps.GET("/search", oppHandler.Search)
			opps.GET("/:id", oppHandler.GetByID)

			// 需要認證
			authOpps := opps.Group("")
			authOpps.Use(AuthMiddleware(cfg))
			{
				authOpps.POST("", oppHandler.Create)
			}

		}

		// Applications
		applications := v1.Group("/applications")
		applications.Use(AuthMiddleware(cfg)) // Assuming authMiddleware refers to AuthMiddleware(cfg)
		{
			applications.POST("", appHandler.Create)
			applications.GET("", appHandler.List)
			applications.GET("/:id", appHandler.GetByID)
			applications.PUT("/:id", appHandler.UpdateStatus)
			applications.DELETE("/:id", appHandler.Delete)
		}

		// Notifications
		notifications := v1.Group("/users/me/notifications")
		notifications.Use(AuthMiddleware(cfg)) // Assuming authMiddleware refers to AuthMiddleware(cfg)
		{
			notifications.GET("", notifHandler.List)
			notifications.PUT("/:id/read", notifHandler.MarkAsRead)
			notifications.PUT("/read-all", notifHandler.MarkAllAsRead)
		}

		// Admin
		admin := v1.Group("/admin")
		admin.Use(AuthMiddleware(cfg))   // First check if authenticated
		admin.Use(AdminAuthMiddleware()) // Then check if admin
		{
			admin.GET("/stats", adminHandler.GetStats)
			admin.GET("/images/pending", adminHandler.ListPendingImages)
			admin.PUT("/images/:id/review", adminHandler.ReviewImage)
			admin.GET("/users", adminHandler.ListUsers)
			admin.PUT("/users/:id/status", adminHandler.UpdateUserStatus)
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
