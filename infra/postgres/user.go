package postgres

import (
	"context"
	"time"

	"taine-api/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// GORM用のテーブル行
type User struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	SubID     string         `gorm:"type:text;not null;uniqueIndex"`
	Name      string         `gorm:"type:text;not null;default:''"`
	AvatarURL string         `gorm:"type:text;not null;default:'';column:avatar_url"`
	CreatedAt time.Time      `gorm:"type:timestamptz;not null;default:now();column:created_at"`
	UpdatedAt time.Time      `gorm:"type:timestamptz;not null;default:now();column:updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"type:timestamptz;index;column:deleted_at"`
}

func (User) TableName() string { return "users" }

type userRepository struct{ db *gorm.DB }

func NewUserRepository(db *gorm.DB) domain.UserRepository { return &userRepository{db: db} }

func (r *userRepository) UpsertUser(ctx context.Context, u *domain.User) (*domain.User, error) {
	row := &User{
		SubID:     u.SubID,
		Name:      u.Name,
		AvatarURL: u.AvatarURL,
	}

	if err := r.db.WithContext(ctx).
		Clauses(
			clause.OnConflict{
				Columns: []clause.Column{{Name: "sub_id"}},
				DoUpdates: clause.Assignments(map[string]interface{}{
					"name":       u.Name,
					"avatar_url": u.AvatarURL,
					"updated_at": gorm.Expr("now()"),
				}),
			},
			clause.Returning{}, // Postgres: RETURNING *
		).
		Create(row).Error; err != nil {
		return nil, err
	}

	return &domain.User{
		ID:        row.ID,
		SubID:     row.SubID,
		Name:      row.Name,
		AvatarURL: row.AvatarURL,
	}, nil
}

func (r *userRepository) GetUserBySubID(ctx context.Context, subID string) (*domain.User, error) {
	var row User
	if err := r.db.WithContext(ctx).First(&row, "sub_id = ?", subID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &domain.User{ID: row.ID, SubID: row.SubID, Name: row.Name, AvatarURL: row.AvatarURL}, nil
}

func (r *userRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	var row User
	if err := r.db.WithContext(ctx).First(&row, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &domain.User{ID: row.ID, SubID: row.SubID, Name: row.Name, AvatarURL: row.AvatarURL}, nil
}

func (r *userRepository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	// gorm.DeletedAt なので Delete でソフトデリートになる
	res := r.db.WithContext(ctx).Delete(&User{}, "id = ?", id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *userRepository) SoftDeleteBySubID(ctx context.Context, subID string) error {
	// gorm.DeletedAt なので Delete でソフトデリートになる
	res := r.db.WithContext(ctx).Delete(&User{}, "sub_id = ?", subID)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
