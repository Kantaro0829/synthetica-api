package repository

import (
	"context"
	"testing"

	"synthetica/internal/domain"

	"github.com/stretchr/testify/assert"
)

func TestQuestionnaireRepository_GetByUserID(t *testing.T) {
	cleanDB(t)
	repo := NewQuestionnaireRepository(testDB)
	ctx := context.Background()

	// Create User
	user := &domain.User{Name: "Q User", Email: "q@example.com"}
	testDB.Create(user)

	// Tests empty
	found, err := repo.GetByUserID(ctx, user.ID)
	assert.Error(t, err) // Or nil depending on impl, assuming record not found error or similar
	assert.Nil(t, found)

	// Store Q (using direct DB for setup if Store isn't available or just use repo Store if implemented)
	// Assuming Store exists in repo interface
	q := &domain.Questionnaire{UserID: user.ID, Answer: 5}
	err = repo.Store(ctx, q)
	assert.NoError(t, err)

	// Check
	found, err = repo.GetByUserID(ctx, user.ID)
	assert.NoError(t, err)
	assert.Equal(t, 5, found.Answer)
}
