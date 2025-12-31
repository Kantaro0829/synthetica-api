package repository

import (
	"context"
	"synthetica/internal/domain"

	"gorm.io/gorm"
)

type storyRepository struct {
	db *gorm.DB
}

func NewStoryRepository(db *gorm.DB) domain.StoryRepository {
	return &storyRepository{db}
}

func (r *storyRepository) Create(ctx context.Context, story *domain.Story) error {
	return getDB(ctx, r.db).WithContext(ctx).Create(story).Error
}

func (r *storyRepository) Fetch(ctx context.Context, userID uint) ([]domain.Story, error) {
	var stories []domain.Story
	err := getDB(ctx, r.db).WithContext(ctx).Preload("User").Preload("Likes").Preload("Comments").Preload("Comments.User").Order("created_at desc").Find(&stories).Error
	if err != nil {
		return nil, err
	}

	// Calculate Liked status
	for i := range stories {
		for _, like := range stories[i].Likes {
			if like.UserID == userID {
				stories[i].Liked = true
				break
			}
		}
	}

	return stories, nil
}

func (r *storyRepository) ToggleLike(ctx context.Context, storyID uint, userID uint) error {
	var like domain.Like
	db := getDB(ctx, r.db).WithContext(ctx)

	// Check if like exists
	err := db.Where("story_id = ? AND from_user_id = ?", storyID, userID).First(&like).Error
	if err == nil {
		// Like exists, user wants to unlike? Or disable/error?
		// Requirement: "if record already recorded disable click" -> Implies we probably just fail or do nothing if user tries to "like" again.
		// However, "use other like icon" usually implies toggle.
		// But User said: "disable click and use other other like icon".
		// This means ON THE FRONTEND it is disabled.
		// ON THE BACKEND, we should probably return an error if they try to like again, or just return success (idempotent).
		// I'll make it idempotent-ish or return error. Let's return error "already liked" so frontend knows?
		// Or simpler: The requirement says "add record... if record already recorded".
		// It DOES NOT say "remove record". It just says "disable click".
		// So I will only implement ADD. If exists, do nothing or error.
		return nil // Doing nothing if already exists is safe/idempotent.
	} else if err != gorm.ErrRecordNotFound {
		return err
	}

	// Create like
	newLike := domain.Like{
		StoryID: storyID,
		UserID:  userID,
	}
	return db.Create(&newLike).Error
}
