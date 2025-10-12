package domain

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrTweetNotFound = errors.New("tweet not found")
)

type Tweet struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// TweetResponse - APIレスポンス用のtweet構造体（sub_idを含む）
type TweetResponse struct {
	ID        uuid.UUID `json:"id"`
	SubID     string    `json:"sub_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type TweetRepository interface {
	CreateTweet(ctx context.Context, tweet *Tweet) (*Tweet, error)
	GetTweetsByUserID(ctx context.Context, userID uuid.UUID) ([]*Tweet, error)
	GetAllTweets(ctx context.Context) ([]*Tweet, error)
	GetTweetByID(ctx context.Context, id uuid.UUID) (*Tweet, error)
	UpdateTweet(ctx context.Context, tweet *Tweet) (*Tweet, error)
	DeleteTweet(ctx context.Context, id uuid.UUID) error
}
