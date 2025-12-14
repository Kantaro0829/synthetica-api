package main

import (
	"synthetica/internal/config"
	"synthetica/internal/delivery/http"
	"synthetica/internal/domain"
	"synthetica/internal/repository"
	"synthetica/internal/usecase"
	"synthetica/pkg/database"
	"synthetica/pkg/logger"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// 1. Init Logger
	logger.InitLogger()

	// Load .env
	if err := godotenv.Load(); err != nil {
		logger.Log.Info("No .env file found or failed to load")
	}

	// 2. Init DB
	database.InitDB()
	config.InitOauth()
	// Auto Migrate
	err := database.DB.AutoMigrate(&domain.User{})
	if err != nil {
		logger.Log.Fatal("failed to migrate database")
	}

	// 3. Setup Layers
	userRepo := repository.NewUserRepository(database.DB)
	userUsecase := usecase.NewUserUsecase(userRepo, 2*time.Second)

	// 4. Setup Router
	r := gin.Default()
	http.NewUserHandler(r, userUsecase)
	http.NewAuthHandler(r, userUsecase)

	// 5. Run
	if err := r.Run(":8080"); err != nil {
		logger.Log.Fatal("failed to start server")
	}
}
