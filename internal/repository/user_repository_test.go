package repository

import (
	"context"
	"testing"

	"synthetica/internal/domain"

	"github.com/stretchr/testify/assert"
)

func TestUserRepository_Create(t *testing.T) {
	cleanDB(t)
	repo := NewUserRepository(testDB)
	ctx := context.Background()

	user := &domain.User{
		Name:     "Test User",
		Email:    "test@example.com",
		GoogleID: "google123",
	}

	err := repo.Create(ctx, user)
	assert.NoError(t, err)
	assert.NotZero(t, user.ID)
}

func TestUserRepository_GetByEmail(t *testing.T) {
	cleanDB(t)
	repo := NewUserRepository(testDB)
	ctx := context.Background()

	user := &domain.User{Name: "Finder", Email: "find@example.com"}
	testDB.Create(user)

	found, err := repo.GetByEmail(ctx, "find@example.com")
	assert.NoError(t, err)
	assert.Equal(t, user.ID, found.ID)

	_, err = repo.GetByEmail(ctx, "missing@example.com")
	assert.Error(t, err)
}
