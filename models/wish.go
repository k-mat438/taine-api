package models

import (
	"time"

	"github.com/google/uuid"
)

// Wish represents a wish in the system
type Wish struct {
	ID             uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	OrganizationID uuid.UUID  `gorm:"type:uuid;not null"`
	Title          string     `gorm:"type:text;not null"`
	Note           string     `gorm:"type:text;not null;default:''"`
	OrderNo        int        `gorm:"type:int;not null;default:0"`
	CreatedAt      time.Time  `gorm:"type:timestamptz;not null;default:now()"`
	UpdatedAt      time.Time  `gorm:"type:timestamptz;not null;default:now()"`
	DeletedAt      *time.Time `gorm:"type:timestamptz;index"`

	// Relations
	Organization Organization `gorm:"foreignKey:OrganizationID;constraint:OnDelete:CASCADE"`
}

// TableName returns the table name for the Wish model
func (Wish) TableName() string {
	return "wishes"
}
