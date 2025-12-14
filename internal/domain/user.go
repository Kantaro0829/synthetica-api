package domain

import (
	"context"
	"time"
)

type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name"`
	Email     string    `json:"email" gorm:"unique"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id uint) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	Fetch(ctx context.Context) ([]User, error)
}

type UserUsecase interface {
	Store(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id uint) (*User, error)
	Fetch(ctx context.Context) ([]User, error)
}
