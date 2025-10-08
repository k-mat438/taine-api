package models

import (
	"time"

	"github.com/google/uuid"
)

// Organization represents an organization in the system
type Organization struct {
	ID         uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ExternalID string     `gorm:"type:text;not null;uniqueIndex"`
	Name       string     `gorm:"type:text;not null"`
	CreatedAt  time.Time  `gorm:"type:timestamptz;not null;default:now()"`
	UpdatedAt  time.Time  `gorm:"type:timestamptz;not null;default:now()"`
	DeletedAt  *time.Time `gorm:"type:timestamptz;index"`
}

// TableName returns the table name for the Organization model
func (Organization) TableName() string {
	return "organizations"
}
