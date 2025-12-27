package main

import (
	"log"
	"synthetica/internal/domain"
	"synthetica/pkg/database"
	"synthetica/pkg/logger"

	"github.com/joho/godotenv"
)

func main() {
	logger.InitLogger()
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or failed to load")
	}
	database.InitDB()

	user := domain.User{
		Name:     "Test User 2",
		Email:    "test2@example.com",
		GoogleID: "999999999999999999999",
	}

	// Check if exists
	var existing domain.User
	if err := database.DB.Where("google_id = ?", user.GoogleID).First(&existing).Error; err == nil {
		log.Println("User already exists, skipping")
		return
	}

	if err := database.DB.Create(&user).Error; err != nil {
		log.Fatal("Failed to create user:", err)
	}

	log.Println("User created successfully")
}
