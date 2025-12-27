package domain

import (
	"context"
	"time"
)

type Questionnaire struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id"`
	User      User      `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Answer    int       `json:"answer"`
	CreatedAt time.Time `json:"created_at"`
}

type QuestionnaireRepository interface {
	Store(ctx context.Context, q *Questionnaire) error
	GetByUserID(ctx context.Context, userID uint) (*Questionnaire, error)
}

type QuestionnaireUsecase interface {
	Store(ctx context.Context, googleID string, answer int) error
	GetStatus(ctx context.Context, googleID string) (*Questionnaire, error)
}
