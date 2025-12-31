package usecase

import (
	"context"
	"synthetica/internal/domain"
	"time"
)

type storyUsecase struct {
	storyRepo      domain.StoryRepository
	contextTimeout time.Duration
}

func NewStoryUsecase(storyRepo domain.StoryRepository, timeout time.Duration) domain.StoryUsecase {
	return &storyUsecase{
		storyRepo:      storyRepo,
		contextTimeout: timeout,
	}
}

func (u *storyUsecase) Create(c context.Context, story *domain.Story) error {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()
	return u.storyRepo.Create(ctx, story)
}

func (u *storyUsecase) Fetch(c context.Context, userID uint) ([]domain.Story, error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()
	return u.storyRepo.Fetch(ctx, userID)
}

func (u *storyUsecase) ToggleLike(c context.Context, storyID uint, userID uint) error {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()
	return u.storyRepo.ToggleLike(ctx, storyID, userID)
}
