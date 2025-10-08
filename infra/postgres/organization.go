package postgres

import (
	"context"

	"taine-api/domain"
	"taine-api/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type organizationRepository struct{ db *gorm.DB }

func NewOrganizationRepository(db *gorm.DB) domain.OrganizationRepository {
	return &organizationRepository{db: db}
}

func (r *organizationRepository) UpsertByExternalID(ctx context.Context, externalID, name string) (*domain.Organization, error) {
	row := &models.Organization{
		ExternalID: externalID,
		Name:       name,
	}

	if err := r.db.WithContext(ctx).
		Clauses(
			clause.OnConflict{
				Columns: []clause.Column{{Name: "external_id"}},
				DoUpdates: clause.Assignments(map[string]interface{}{
					"name":       name,
					"updated_at": gorm.Expr("now()"),
					"deleted_at": nil, // ソフトデリートを解除
				}),
			},
			clause.Returning{}, // Postgres: RETURNING *
		).
		Create(row).Error; err != nil {
		return nil, err
	}

	return &domain.Organization{
		ID:         row.ID,
		ExternalID: row.ExternalID,
		Name:       row.Name,
	}, nil
}

func (r *organizationRepository) SoftDeleteByExternalID(ctx context.Context, externalID string) error {
	// gorm.DeletedAt なので Delete でソフトデリートになる
	res := r.db.WithContext(ctx).Delete(&models.Organization{}, "external_id = ?", externalID)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *organizationRepository) FindByExternalID(ctx context.Context, externalID string) (*domain.Organization, error) {
	var row models.Organization
	if err := r.db.WithContext(ctx).First(&row, "external_id = ?", externalID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &domain.Organization{
		ID:         row.ID,
		ExternalID: row.ExternalID,
		Name:       row.Name,
	}, nil
}
