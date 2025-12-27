package repository

import (
	"context"
	"synthetica/internal/domain"

	"gorm.io/gorm"
)

type questionnaireRepository struct {
	Conn *gorm.DB
}

func NewQuestionnaireRepository(conn *gorm.DB) domain.QuestionnaireRepository {
	return &questionnaireRepository{Conn: conn}
}

func (r *questionnaireRepository) Store(ctx context.Context, q *domain.Questionnaire) error {
	return r.Conn.Create(q).Error
}

func (r *questionnaireRepository) GetByUserID(ctx context.Context, userID uint) (*domain.Questionnaire, error) {
	var q domain.Questionnaire
	// Find the FIRST answer. If multiple, we just take one.
	if err := r.Conn.Where("user_id = ?", userID).First(&q).Error; err != nil {
		return nil, err
	}
	return &q, nil
}
