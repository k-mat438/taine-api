package domain

import (
	"context"

	"github.com/google/uuid"
)

type OrganizationMember struct {
	ID             uuid.UUID
	UserID         uuid.UUID
	OrganizationID uuid.UUID
	Role           string
}

type MembershipRepository interface {
	UpsertByUserAndOrg(ctx context.Context, userID, orgID uuid.UUID, role string) (*OrganizationMember, error)
	SoftDeleteByUserAndOrg(ctx context.Context, userID, orgID uuid.UUID) error
}
