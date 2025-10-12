package postgres

import (
	"context"
	"time"

	"taine-api/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// GORM用のテーブル行
type Tweet struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;column:user_id"`
	Content   string    `gorm:"type:text;not null"`
	CreatedAt time.Time `gorm:"type:timestamptz;not null;default:now();column:created_at"`
	UpdatedAt time.Time `gorm:"type:timestamptz;not null;default:now();column:updated_at"`
}

func (Tweet) TableName() string { return "tweets" }

type tweetRepository struct{ db *gorm.DB }

func NewTweetRepository(db *gorm.DB) domain.TweetRepository { return &tweetRepository{db: db} }

func (r *tweetRepository) CreateTweet(ctx context.Context, t *domain.Tweet) (*domain.Tweet, error) {
	row := &Tweet{
		UserID:  t.UserID,
		Content: t.Content,
	}

	if err := r.db.WithContext(ctx).Create(row).Error; err != nil {
		return nil, err
	}

	return &domain.Tweet{
		ID:        row.ID,
		UserID:    row.UserID,
		Content:   row.Content,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}, nil
}

func (r *tweetRepository) GetTweetsByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Tweet, error) {
	var rows []Tweet
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at DESC").Find(&rows).Error; err != nil {
		return nil, err
	}

	tweets := make([]*domain.Tweet, len(rows))
	for i, row := range rows {
		tweets[i] = &domain.Tweet{
			ID:        row.ID,
			UserID:    row.UserID,
			Content:   row.Content,
			CreatedAt: row.CreatedAt,
			UpdatedAt: row.UpdatedAt,
		}
	}
	return tweets, nil
}

func (r *tweetRepository) GetAllTweets(ctx context.Context) ([]*domain.Tweet, error) {
	var rows []Tweet
	if err := r.db.WithContext(ctx).Order("created_at DESC").Find(&rows).Error; err != nil {
		return nil, err
	}

	tweets := make([]*domain.Tweet, len(rows))
	for i, row := range rows {
		tweets[i] = &domain.Tweet{
			ID:        row.ID,
			UserID:    row.UserID,
			Content:   row.Content,
			CreatedAt: row.CreatedAt,
			UpdatedAt: row.UpdatedAt,
		}
	}
	return tweets, nil
}

func (r *tweetRepository) GetTweetByID(ctx context.Context, id uuid.UUID) (*domain.Tweet, error) {
	var row Tweet
	if err := r.db.WithContext(ctx).First(&row, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrTweetNotFound
		}
		return nil, err
	}
	return &domain.Tweet{
		ID:        row.ID,
		UserID:    row.UserID,
		Content:   row.Content,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}, nil
}

func (r *tweetRepository) UpdateTweet(ctx context.Context, t *domain.Tweet) (*domain.Tweet, error) {
	row := &Tweet{
		ID:        t.ID,
		UserID:    t.UserID,
		Content:   t.Content,
		UpdatedAt: time.Now(),
	}

	if err := r.db.WithContext(ctx).Save(row).Error; err != nil {
		return nil, err
	}

	return &domain.Tweet{
		ID:        row.ID,
		UserID:    row.UserID,
		Content:   row.Content,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}, nil
}

func (r *tweetRepository) DeleteTweet(ctx context.Context, id uuid.UUID) error {
	res := r.db.WithContext(ctx).Delete(&Tweet{}, "id = ?", id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return domain.ErrTweetNotFound
	}
	return nil
}
