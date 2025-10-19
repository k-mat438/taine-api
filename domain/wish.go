package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Wish represents a wish domain model
type Wish struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID
	Title          string
	Note           string
	OrderNo        int
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      *time.Time
}

// WishRepository defines the interface for wish data operations
type WishRepository interface {
	Create(ctx context.Context, wish *Wish) (*Wish, error)
	FindByID(ctx context.Context, id uuid.UUID) (*Wish, error)
	FindByOrganizationID(ctx context.Context, organizationID uuid.UUID) ([]*Wish, error)
	Update(ctx context.Context, wish *Wish) (*Wish, error)
	Delete(ctx context.Context, id uuid.UUID) error
	SoftDelete(ctx context.Context, id uuid.UUID) error
	Restore(ctx context.Context, id uuid.UUID) error
	UpdateOrder(ctx context.Context, id uuid.UUID, orderNo int) error
}
