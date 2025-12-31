package domain

import (
	"time"
)

type Comment struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	StoryID   uint      `json:"story_id"`
	UserID    uint      `json:"user_id" gorm:"column:from_user_id"`
	User      User      `json:"user" gorm:"foreignKey:UserID"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}
