package postgres

import (
	"context"

	"taine-api/domain"
	"taine-api/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type membershipRepository struct{ db *gorm.DB }

func NewMembershipRepository(db *gorm.DB) domain.MembershipRepository {
	return &membershipRepository{db: db}
}

func (r *membershipRepository) UpsertByUserAndOrg(ctx context.Context, userID, orgID uuid.UUID, role string) (*domain.OrganizationMember, error) {
	row := &models.OrganizationMember{
		UserID:         userID,
		OrganizationID: orgID,
		Role:           role,
	}

	if err := r.db.WithContext(ctx).
		Clauses(
			clause.OnConflict{
				Columns: []clause.Column{{Name: "user_id"}, {Name: "organization_id"}},
				DoUpdates: clause.Assignments(map[string]interface{}{
					"role":       role,
					"updated_at": gorm.Expr("now()"),
				}),
			},
			clause.Returning{}, // Postgres: RETURNING *
		).
		Create(row).Error; err != nil {
		return nil, err
	}

	return &domain.OrganizationMember{
		ID:             row.ID,
		UserID:         row.UserID,
		OrganizationID: row.OrganizationID,
		Role:           row.Role,
	}, nil
}

func (r *membershipRepository) SoftDeleteByUserAndOrg(ctx context.Context, userID, orgID uuid.UUID) error {
	// gorm.DeletedAt なので Delete でソフトデリートになる
	res := r.db.WithContext(ctx).Delete(&models.OrganizationMember{}, "user_id = ? AND organization_id = ?", userID, orgID)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
