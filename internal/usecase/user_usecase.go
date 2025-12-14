package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"synthetica/internal/domain"
	"time"

	"golang.org/x/oauth2"
)

type userUsecase struct {
	userRepo       domain.UserRepository
	contextTimeout time.Duration
}

func NewUserUsecase(userRepo domain.UserRepository, timeout time.Duration) domain.UserUsecase {
	return &userUsecase{
		userRepo:       userRepo,
		contextTimeout: timeout,
	}
}

func (u *userUsecase) Store(ctx context.Context, user *domain.User) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()
	return u.userRepo.Create(ctx, user)
}

func (u *userUsecase) GetByID(ctx context.Context, id uint) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()
	return u.userRepo.GetByID(ctx, id)
}

func (u *userUsecase) Fetch(ctx context.Context) ([]domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()
	return u.userRepo.Fetch(ctx)
}

func (u *userUsecase) LoginWithGoogleOAuth(ctx context.Context, token *oauth2.Token) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// Use the token to get user info
	client := http.Client{}
	req, err := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v2/userinfo", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to get user info")
	}

	var payload struct {
		Email string `json:"email"`
		Name  string `json:"name"`
		ID    string `json:"id"` // Google ID from UserInfo endpoint
	}

	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}

	email := payload.Email
	name := payload.Name
	googleID := payload.ID

	// 2. Check if user exists by Google ID
	user, err := u.userRepo.GetByGoogleID(ctx, googleID)
	if err == nil {
		return user, nil
	}

	// 3. Check if user exists by Email (to link)
	user, err = u.userRepo.GetByEmail(ctx, email)
	if err == nil {
		user.GoogleID = googleID
		return user, nil
	}

	// 4. Create new User
	newUser := &domain.User{
		Name:     name,
		Email:    email,
		GoogleID: googleID,
	}
	err = u.userRepo.Create(ctx, newUser)
	if err != nil {
		return nil, err
	}

	return newUser, nil
}
