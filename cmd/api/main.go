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

	err := database.DB.AutoMigrate(&domain.User{}, &domain.Questionnaire{})
	if err != nil {
		logger.Log.Fatal("failed to migrate database: " + err.Error())
	}

	// 3. Setup Layers
	userRepo := repository.NewUserRepository(database.DB)
	userUsecase := usecase.NewUserUsecase(userRepo, 2*time.Second)

	questionnaireRepo := repository.NewQuestionnaireRepository(database.DB)
	questionnaireUsecase := usecase.NewQuestionnaireUsecase(questionnaireRepo, userRepo, 2*time.Second)

	// 4. Setup Router
	r := gin.Default()
	// CORS Middleware to allow requests from localhost:3000
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000") // Temporarily allow 3000 as requested, or * but Credentials requires specific origin.
		// Actually, "Access-Control-Allow-Credentials" is true, so "Access-Control-Allow-Origin" CANNOT be "*".
		// It must be the specific origin.
		// I will use the Origin header from the request to allow "any" origin dynamically if I want to simulate "*".
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		} else {
			// Fallback or default
			c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		}
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	http.NewUserHandler(r, userUsecase)
	http.NewAuthHandler(r, userUsecase)
	http.NewQuestionnaireHandler(r, questionnaireUsecase)

	// 5. Run
	if err := r.Run(":8080"); err != nil {
		logger.Log.Fatal("failed to start server")
	}
}
