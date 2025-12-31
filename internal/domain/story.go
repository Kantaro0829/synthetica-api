package domain

import (
	"context"
	"time"
)

type Story struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Title     string    `json:"title"`
	Detail    string    `json:"detail"`
	UserID    uint      `json:"user_id"`
	User      User      `json:"user" gorm:"foreignKey:UserID"`
	Likes     []Like    `json:"likes" gorm:"foreignKey:StoryID"`
	Liked     bool      `json:"liked" gorm:"-"`
	Comments  []Comment `json:"comments" gorm:"foreignKey:StoryID"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type StoryRepository interface {
	Create(ctx context.Context, story *Story) error
	Fetch(ctx context.Context, userID uint) ([]Story, error)
	ToggleLike(ctx context.Context, storyID uint, userID uint) error
}

type StoryUsecase interface {
	Create(ctx context.Context, story *Story) error
	Fetch(ctx context.Context, userID uint) ([]Story, error)
	ToggleLike(ctx context.Context, storyID uint, userID uint) error
}
