package domain

import (
	"time"
)

type Like struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	StoryID   uint      `json:"story_id"`
	UserID    uint      `json:"user_id" gorm:"column:from_user_id"`
	User      User      `json:"-" gorm:"foreignKey:UserID"`
	CreatedAt time.Time `json:"created_at"`
}
