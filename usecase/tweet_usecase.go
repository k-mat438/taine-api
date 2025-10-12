package usecase

import (
	"context"
	"taine-api/domain"

	"github.com/google/uuid"
)

type TweetUsecase interface {
	CreateTweet(ctx context.Context, userID uuid.UUID, content string) (*domain.Tweet, error)
	GetTweetsByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Tweet, error)
	GetAllTweets(ctx context.Context) ([]*domain.Tweet, error)
	GetTweetByID(ctx context.Context, id uuid.UUID) (*domain.Tweet, error)
	UpdateTweet(ctx context.Context, id uuid.UUID, userID uuid.UUID, content string) (*domain.Tweet, error)
	DeleteTweet(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
	GetAllTweetsWithUsers(ctx context.Context) ([]*domain.TweetResponse, error)
	GetTweetsByUserIDWithUser(ctx context.Context, userID uuid.UUID) ([]*domain.TweetResponse, error)
	GetTweetByIDWithUser(ctx context.Context, id uuid.UUID) (*domain.TweetResponse, error)
}

type tweetUsecase struct {
	tweetRepository domain.TweetRepository
	userRepository  domain.UserRepository
}

func NewTweetUsecase(tweetRepository domain.TweetRepository, userRepository domain.UserRepository) TweetUsecase {
	return &tweetUsecase{tweetRepository: tweetRepository, userRepository: userRepository}
}

func (u *tweetUsecase) CreateTweet(ctx context.Context, userID uuid.UUID, content string) (*domain.Tweet, error) {
	tweet := &domain.Tweet{
		UserID:  userID,
		Content: content,
	}
	return u.tweetRepository.CreateTweet(ctx, tweet)
}

func (u *tweetUsecase) GetTweetsByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Tweet, error) {
	return u.tweetRepository.GetTweetsByUserID(ctx, userID)
}

func (u *tweetUsecase) GetAllTweets(ctx context.Context) ([]*domain.Tweet, error) {
	return u.tweetRepository.GetAllTweets(ctx)
}

func (u *tweetUsecase) GetTweetByID(ctx context.Context, id uuid.UUID) (*domain.Tweet, error) {
	return u.tweetRepository.GetTweetByID(ctx, id)
}

func (u *tweetUsecase) UpdateTweet(ctx context.Context, id uuid.UUID, userID uuid.UUID, content string) (*domain.Tweet, error) {
	// まずtweetが存在するかチェック
	tweet, err := u.tweetRepository.GetTweetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 所有者チェック
	if tweet.UserID != userID {
		return nil, domain.ErrTweetNotFound // セキュリティ上、詳細なエラーは返さない
	}

	// 更新
	tweet.Content = content
	return u.tweetRepository.UpdateTweet(ctx, tweet)
}

func (u *tweetUsecase) DeleteTweet(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	// まずtweetが存在するかチェック
	tweet, err := u.tweetRepository.GetTweetByID(ctx, id)
	if err != nil {
		return err
	}

	// 所有者チェック
	if tweet.UserID != userID {
		return domain.ErrTweetNotFound // セキュリティ上、詳細なエラーは返さない
	}

	return u.tweetRepository.DeleteTweet(ctx, id)
}

// GetAllTweetsWithUsers - 全てのtweetをuser情報と一緒に取得
func (u *tweetUsecase) GetAllTweetsWithUsers(ctx context.Context) ([]*domain.TweetResponse, error) {
	tweets, err := u.tweetRepository.GetAllTweets(ctx)
	if err != nil {
		return nil, err
	}

	tweetResponses := make([]*domain.TweetResponse, len(tweets))
	for i, tweet := range tweets {
		user, err := u.userRepository.GetUserByID(ctx, tweet.UserID)
		if err != nil {
			return nil, err
		}
		tweetResponses[i] = &domain.TweetResponse{
			ID:        tweet.ID,
			SubID:     user.SubID,
			Content:   tweet.Content,
			CreatedAt: tweet.CreatedAt,
			UpdatedAt: tweet.UpdatedAt,
		}
	}
	return tweetResponses, nil
}

// GetTweetsByUserIDWithUser - 特定のユーザーのtweetをuser情報と一緒に取得
func (u *tweetUsecase) GetTweetsByUserIDWithUser(ctx context.Context, userID uuid.UUID) ([]*domain.TweetResponse, error) {
	tweets, err := u.tweetRepository.GetTweetsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	user, err := u.userRepository.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	tweetResponses := make([]*domain.TweetResponse, len(tweets))
	for i, tweet := range tweets {
		tweetResponses[i] = &domain.TweetResponse{
			ID:        tweet.ID,
			SubID:     user.SubID,
			Content:   tweet.Content,
			CreatedAt: tweet.CreatedAt,
			UpdatedAt: tweet.UpdatedAt,
		}
	}
	return tweetResponses, nil
}

// GetTweetByIDWithUser - 特定のtweetをuser情報と一緒に取得
func (u *tweetUsecase) GetTweetByIDWithUser(ctx context.Context, id uuid.UUID) (*domain.TweetResponse, error) {
	tweet, err := u.tweetRepository.GetTweetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	user, err := u.userRepository.GetUserByID(ctx, tweet.UserID)
	if err != nil {
		return nil, err
	}

	return &domain.TweetResponse{
		ID:        tweet.ID,
		SubID:     user.SubID,
		Content:   tweet.Content,
		CreatedAt: tweet.CreatedAt,
		UpdatedAt: tweet.UpdatedAt,
	}, nil
}
