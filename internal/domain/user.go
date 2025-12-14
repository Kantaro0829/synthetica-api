package domain

import (
	"context"
	"time"

	"golang.org/x/oauth2"
)

type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name"`
	Email     string    `json:"email" gorm:"unique"`
	Password  string    `json:"-"`
	GoogleID  string    `json:"google_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id uint) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByGoogleID(ctx context.Context, googleID string) (*User, error)
	Fetch(ctx context.Context) ([]User, error)
}

type UserUsecase interface {
	Store(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id uint) (*User, error)
	Fetch(ctx context.Context) ([]User, error)
	LoginWithGoogleOAuth(ctx context.Context, token *oauth2.Token) (*User, error)
}
