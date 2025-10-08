package models

import (
	"time"

	"github.com/google/uuid"
)

// OrganizationMember represents a user's membership in an organization
type OrganizationMember struct {
	ID             uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID         uuid.UUID `gorm:"type:uuid;not null;index:idx_uom_user_id"`
	OrganizationID uuid.UUID `gorm:"type:uuid;not null;index:idx_uom_org_id"`
	Role           string    `gorm:"type:text;not null;index:idx_uom_role"` // 'owner'|'admin'|'member' etc.
	CreatedAt      time.Time `gorm:"type:timestamptz;not null;default:now()"`
	UpdatedAt      time.Time `gorm:"type:timestamptz;not null;default:now()"`
}

// TableName returns the table name for the OrganizationMember model
func (OrganizationMember) TableName() string {
	return "user_organization_memberships"
}

// Role constants
const (
	RoleOwner  = "owner"
	RoleAdmin  = "admin"
	RoleMember = "member"
)
