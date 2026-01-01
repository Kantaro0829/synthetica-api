package repository

import (
	"log"
	"os"
	"testing"

	"synthetica/internal/domain"
	"synthetica/pkg/database"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var testDB *gorm.DB

func TestMain(m *testing.M) {
	// Setup
	dsn := "host=localhost user=user password=password dbname=synthetica_test port=5433 sslmode=disable"
	var err error
	testDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to test database: %v", err)
	}

	// AutoMigrate
	err = testDB.AutoMigrate(&domain.User{}, &domain.Questionnaire{}, &domain.Story{}, &domain.Comment{}, &domain.Like{})
	if err != nil {
		log.Fatalf("Failed to migrate test database: %v", err)
	}

	// Assign to global DB if needed, though we prefer dependency injection
	database.DB = testDB

	code := m.Run()

	// Teardown (optional for Docker containers, but good for cleanup if needed)

	os.Exit(code)
}

// Helper to clean database between tests
func cleanDB(t *testing.T) {
	testDB.Exec("TRUNCATE TABLE users, stories, questionnaires, comments, likes RESTART IDENTITY CASCADE")
}
