package usecase

import (
	"context"
	"synthetica/internal/domain"
	"time"
)

type questionnaireUsecase struct {
	questionnaireRepo domain.QuestionnaireRepository
	userRepo          domain.UserRepository // Need this to lookup UserID from GoogleID
	contextTimeout    time.Duration
}

func NewQuestionnaireUsecase(qRepo domain.QuestionnaireRepository, uRepo domain.UserRepository, timeout time.Duration) domain.QuestionnaireUsecase {
	return &questionnaireUsecase{
		questionnaireRepo: qRepo,
		userRepo:          uRepo,
		contextTimeout:    timeout,
	}
}

func (u *questionnaireUsecase) Store(c context.Context, googleID string, answer int) error {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	// 1. Get User by GoogleID
	user, err := u.userRepo.GetByGoogleID(ctx, googleID)
	if err != nil {
		return err // User not found or DB error
	}

	// 2. Create Questionnaire
	q := &domain.Questionnaire{
		UserID: user.ID,
		Answer: answer,
	}

	return u.questionnaireRepo.Store(ctx, q)
}

func (u *questionnaireUsecase) GetStatus(c context.Context, googleID string) (*domain.Questionnaire, error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	user, err := u.userRepo.GetByGoogleID(ctx, googleID)
	if err != nil {
		return nil, err
	}

	q, err := u.questionnaireRepo.GetByUserID(ctx, user.ID)
	if err == nil {
		return q, nil
	}

	return nil, nil // Not found is not an error here, just nil result
}
