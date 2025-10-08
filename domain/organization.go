package domain

import (
	"context"

	"github.com/google/uuid"
)

type Organization struct {
	ID         uuid.UUID
	ExternalID string
	Name       string
}

type OrganizationRepository interface {
	UpsertByExternalID(ctx context.Context, externalID, name string) (*Organization, error)
	SoftDeleteByExternalID(ctx context.Context, externalID string) error
	FindByExternalID(ctx context.Context, externalID string) (*Organization, error)
}
