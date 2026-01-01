package repository

import (
	"context"
	"testing"

	"synthetica/internal/domain"

	"github.com/stretchr/testify/assert"
)

func TestStoryRepository_Create(t *testing.T) {
	cleanDB(t)
	repo := NewStoryRepository(testDB)
	ctx := context.Background()

	// Create User first
	user := &domain.User{Name: "Story Author", Email: "author@example.com"}
	testDB.Create(user)

	story := &domain.Story{
		Title:  "Test Story",
		Detail: "Details here",
		UserID: user.ID,
	}

	err := repo.Create(ctx, story)
	assert.NoError(t, err)
	assert.NotZero(t, story.ID)

	// Verify
	var saved domain.Story
	testDB.First(&saved, story.ID)
	assert.Equal(t, "Test Story", saved.Title)
}

func TestStoryRepository_Fetch(t *testing.T) {
	cleanDB(t)
	repo := NewStoryRepository(testDB)
	ctx := context.Background()

	// Setup Data
	user1 := &domain.User{Name: "User 1", Email: "u1@example.com"}
	user2 := &domain.User{Name: "User 2", Email: "u2@example.com"}
	testDB.Create(user1)
	testDB.Create(user2)

	s1 := &domain.Story{Title: "Story 1", Detail: "D1", UserID: user1.ID}
	s2 := &domain.Story{Title: "Story 2", Detail: "D2", UserID: user2.ID}
	testDB.Create(s1)
	testDB.Create(s2)

	// Tests
	stories, err := repo.Fetch(ctx, user1.ID)
	assert.NoError(t, err)
	assert.Len(t, stories, 2)
	// Order is created_at desc (default check might need to be careful with timestamp precision but usually s2 is newer)
	assert.Equal(t, "Story 2", stories[0].Title)
}

func TestStoryRepository_ToggleLike(t *testing.T) {
	cleanDB(t)
	repo := NewStoryRepository(testDB)
	ctx := context.Background()

	user := &domain.User{Name: "Liker", Email: "liker@example.com"}
	testDB.Create(user)
	story := &domain.Story{Title: "Liked Story", Detail: "Details", UserID: user.ID}
	testDB.Create(story)

	// Like
	err := repo.ToggleLike(ctx, story.ID, user.ID)
	assert.NoError(t, err)

	// Verify Like exists
	var count int64
	testDB.Model(&domain.Like{}).Where("story_id = ? AND from_user_id = ?", story.ID, user.ID).Count(&count)
	assert.Equal(t, int64(1), count)

	// Toggle again (should do nothing or error based on implementation - checked implementation: currently returns nil if exists)
	err = repo.ToggleLike(ctx, story.ID, user.ID)
	assert.NoError(t, err)

	// Count still 1
	testDB.Model(&domain.Like{}).Where("story_id = ? AND from_user_id = ?", story.ID, user.ID).Count(&count)
	assert.Equal(t, int64(1), count)
}
