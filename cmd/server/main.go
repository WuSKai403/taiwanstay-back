package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/taiwanstay/taiwanstay-back/internal/api"
	"github.com/taiwanstay/taiwanstay-back/internal/repository"
	"github.com/taiwanstay/taiwanstay-back/internal/service"
	"github.com/taiwanstay/taiwanstay-back/pkg/config"
	"github.com/taiwanstay/taiwanstay-back/pkg/database"
	"github.com/taiwanstay/taiwanstay-back/pkg/email"
	"github.com/taiwanstay/taiwanstay-back/pkg/gcp"
	"github.com/taiwanstay/taiwanstay-back/pkg/logger"
)

func main() {
	// 1. Load Config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 2. Init Logger
	logger.InitLogger(cfg.Server.Mode)
	logger.Info("Starting TaiwanStay Backend", "port", cfg.Server.Port, "env", cfg.Server.Mode)

	// 3. Connect to Database
	mongoClient, err := database.Connect(cfg.Database.URI)
	if err != nil {
		logger.Error("Failed to connect to MongoDB", "error", err)
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer database.Close(mongoClient)
	logger.Info("Connected to MongoDB")

	// 4. Init GCP Clients
	ctx := context.Background()
	storageClient, err := gcp.NewStorageClient(ctx)
	if err != nil {
		logger.Warn("Failed to init GCP Storage Client (Check credentials)", "error", err)
	} else {
		defer storageClient.Close()
		logger.Info("GCP Storage Client initialized")
	}

	visionClient, err := gcp.NewVisionClient(ctx)
	if err != nil {
		logger.Warn("Failed to init GCP Vision Client (Check credentials)", "error", err)
	} else {
		defer visionClient.Close()
		logger.Info("GCP Vision Client initialized")
	}

	// 5. Dependency Injection
	db := mongoClient.Database(cfg.Database.Database)
	// Email Sender
	primarySender := email.NewBrevoSender(cfg)
	secondarySender := email.NewMailerLiteSender(cfg)
	emailSender := email.NewFallbackSender(primarySender, secondarySender)

	// Repositories
	userRepo := repository.NewUserRepository(db.Collection("users"))
	hostRepo := repository.NewHostRepository(db.Collection("hosts"))
	oppRepo := repository.NewOpportunityRepository(db.Collection("opportunities"))
	appRepo := repository.NewApplicationRepository(db.Collection("applications"))
	notifRepo := repository.NewNotificationRepository(db.Collection("notifications"))
	bookmarkRepo := repository.NewBookmarkRepository(db.Collection("bookmarks"))

	// Services
	userService := service.NewUserService(userRepo, cfg)

	imageCollection := db.Collection("images")
	imageRepo := repository.NewImageRepository(imageCollection)
	imageService := service.NewImageService(imageRepo, storageClient, visionClient, cfg)

	hostService := service.NewHostService(hostRepo)
	oppService := service.NewOpportunityService(oppRepo)
	notifService := service.NewNotificationService(notifRepo, userRepo, emailSender)
	appService := service.NewApplicationService(appRepo, oppRepo, hostRepo, notifService)
	adminService := service.NewAdminService(userRepo, imageRepo, appRepo, imageService)
	bookmarkService := service.NewBookmarkService(bookmarkRepo, oppRepo)

	// Handlers
	userHandler := api.NewUserHandler(userService)
	imageHandler := api.NewImageHandler(imageService)
	hostHandler := api.NewHostHandler(hostService)
	oppHandler := api.NewOpportunityHandler(oppService, hostService)
	appHandler := api.NewApplicationHandler(appService)
	notifHandler := api.NewNotificationHandler(notifService)
	adminHandler := api.NewAdminHandler(adminService, oppService)
	bookmarkHandler := api.NewBookmarkHandler(bookmarkService)

	// 6. Setup Server
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()

	// Setup Routes
	api.SetupRoutes(router, userHandler, imageHandler, hostHandler, oppHandler, appHandler, notifHandler, adminHandler, bookmarkHandler, cfg)

	// 7. Run Server
	addr := ":" + cfg.Server.Port
	logger.Info("Server listening on " + addr)
	if err := router.Run(addr); err != nil {
		logger.Error("Failed to run server", "error", err)
	}
}
