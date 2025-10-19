package postgres

import (
	"context"
	"time"

	"taine-api/domain"
	"taine-api/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type wishRepository struct {
	db *gorm.DB
}

func NewWishRepository(db *gorm.DB) domain.WishRepository {
	return &wishRepository{db: db}
}

func (r *wishRepository) Create(ctx context.Context, wish *domain.Wish) (*domain.Wish, error) {
	row := &models.Wish{
		OrganizationID: wish.OrganizationID,
		Title:          wish.Title,
		Note:           wish.Note,
		OrderNo:        wish.OrderNo,
	}

	if err := r.db.WithContext(ctx).Create(row).Error; err != nil {
		return nil, err
	}

	return r.toDomain(row), nil
}

func (r *wishRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Wish, error) {
	var row models.Wish
	if err := r.db.WithContext(ctx).First(&row, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return r.toDomain(&row), nil
}

func (r *wishRepository) FindByOrganizationID(ctx context.Context, organizationID uuid.UUID) ([]*domain.Wish, error) {
	var rows []models.Wish
	if err := r.db.WithContext(ctx).
		Where("organization_id = ?", organizationID).
		Where("deleted_at IS NULL").
		Order("order_no DESC, created_at DESC").
		Find(&rows).Error; err != nil {
		return nil, err
	}

	wishes := make([]*domain.Wish, len(rows))
	for i, row := range rows {
		wishes[i] = r.toDomain(&row)
	}
	return wishes, nil
}

func (r *wishRepository) Update(ctx context.Context, wish *domain.Wish) (*domain.Wish, error) {
	updates := map[string]interface{}{
		"title":      wish.Title,
		"note":       wish.Note,
		"order_no":   wish.OrderNo,
		"updated_at": time.Now(),
	}

	var row models.Wish
	if err := r.db.WithContext(ctx).
		Model(&row).
		Where("id = ?", wish.ID).
		Updates(updates).
		First(&row).Error; err != nil {
		return nil, err
	}

	return r.toDomain(&row), nil
}

func (r *wishRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&models.Wish{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *wishRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	now := time.Now()
	result := r.db.WithContext(ctx).
		Model(&models.Wish{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"deleted_at": now,
			"updated_at": now,
		})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *wishRepository) Restore(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).
		Model(&models.Wish{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"deleted_at": nil,
			"updated_at": time.Now(),
		})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *wishRepository) UpdateOrder(ctx context.Context, id uuid.UUID, orderNo int) error {
	result := r.db.WithContext(ctx).
		Model(&models.Wish{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"order_no":   orderNo,
			"updated_at": time.Now(),
		})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *wishRepository) toDomain(row *models.Wish) *domain.Wish {
	return &domain.Wish{
		ID:             row.ID,
		OrganizationID: row.OrganizationID,
		Title:          row.Title,
		Note:           row.Note,
		OrderNo:        row.OrderNo,
		CreatedAt:      row.CreatedAt,
		UpdatedAt:      row.UpdatedAt,
		DeletedAt:      row.DeletedAt,
	}
}
