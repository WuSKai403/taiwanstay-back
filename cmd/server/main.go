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
	userCollection := db.Collection("users")

	userRepo := repository.NewUserRepository(userCollection)
	userService := service.NewUserService(userRepo, cfg)
	userHandler := api.NewUserHandler(userService)

	imageCollection := db.Collection("images")
	imageRepo := repository.NewImageRepository(imageCollection)
	imageService := service.NewImageService(imageRepo, storageClient, visionClient, cfg)
	imageHandler := api.NewImageHandler(imageService)

	hostCollection := db.Collection("hosts")
	hostRepo := repository.NewHostRepository(hostCollection)
	hostService := service.NewHostService(hostRepo)
	hostHandler := api.NewHostHandler(hostService)

	oppCollection := db.Collection("opportunities")
	oppRepo := repository.NewOpportunityRepository(oppCollection)
	oppService := service.NewOpportunityService(oppRepo)
	oppHandler := api.NewOpportunityHandler(oppService, hostService)

	// 6. Setup Server
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()

	// Setup Routes
	api.SetupRoutes(router, userHandler, imageHandler, hostHandler, oppHandler, cfg)

	// 7. Run Server
	addr := ":" + cfg.Server.Port
	logger.Info("Server listening on " + addr)
	if err := router.Run(addr); err != nil {
		logger.Error("Failed to run server", "error", err)
	}
}
