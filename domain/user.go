package domain

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

var (
	ErrUserNotFound         = errors.New("user not found")
	ErrOrganizationNotFound = errors.New("organization not found")
)

type User struct {
	ID        uuid.UUID
	SubID     string
	Name      string
	AvatarURL string
}

type UserRepository interface {
	UpsertUser(ctx context.Context, user *User) (*User, error)
	GetUserBySubID(ctx context.Context, subID string) (*User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*User, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
	SoftDeleteBySubID(ctx context.Context, subID string) error
}
