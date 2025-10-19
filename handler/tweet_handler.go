package handler

import (
	"net/http"
	"taine-api/usecase"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TweetHandler struct {
	tweetUsecase usecase.TweetUsecase
	userUsecase  usecase.UserUsecase
}

func NewTweetHandler(tweetUsecase usecase.TweetUsecase, userUsecase usecase.UserUsecase) *TweetHandler {
	return &TweetHandler{tweetUsecase: tweetUsecase, userUsecase: userUsecase}
}

type CreateTweetRequest struct {
	Content string `json:"content" binding:"required"`
}

type UpdateTweetRequest struct {
	Content string `json:"content" binding:"required"`
}

// CreateTweet - 新しいtweetを作成
func (h *TweetHandler) CreateTweet(c *gin.Context) {
	subID := c.GetString("sub_id")
	user, err := h.userUsecase.GetUserBySubID(c.Request.Context(), subID)
	if err != nil || user == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		return
	}

	var req CreateTweetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tweet, err := h.tweetUsecase.CreateTweet(c.Request.Context(), user.ID, req.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 作成されたtweetをsub_idを含むレスポンス形式で返す
	tweetResponse, err := h.tweetUsecase.GetTweetByIDWithUser(c.Request.Context(), tweet.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, tweetResponse)
}

// GetTweets - 全てのtweetを取得（認証されたユーザーが他のユーザーのtweetも見れるかテスト用）
func (h *TweetHandler) GetTweets(c *gin.Context) {
	tweets, err := h.tweetUsecase.GetAllTweetsWithUsers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"tweets": tweets})
}

// GetMyTweets - 自分のtweetのみを取得
func (h *TweetHandler) GetMyTweets(c *gin.Context) {
	subID := c.GetString("sub_id")
	user, err := h.userUsecase.GetUserBySubID(c.Request.Context(), subID)
	if err != nil || user == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		return
	}

	tweets, err := h.tweetUsecase.GetTweetsByUserIDWithUser(c.Request.Context(), user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"tweets": tweets})
}

// GetTweetByID - 特定のtweetを取得
func (h *TweetHandler) GetTweetByID(c *gin.Context) {
	tweetIDStr := c.Param("id")
	tweetID, err := uuid.Parse(tweetIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tweet ID"})
		return
	}

	tweet, err := h.tweetUsecase.GetTweetByIDWithUser(c.Request.Context(), tweetID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tweet)
}

// UpdateTweet - tweetを更新
func (h *TweetHandler) UpdateTweet(c *gin.Context) {
	subID := c.GetString("sub_id")
	user, err := h.userUsecase.GetUserBySubID(c.Request.Context(), subID)
	if err != nil || user == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		return
	}

	tweetIDStr := c.Param("id")
	tweetID, err := uuid.Parse(tweetIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tweet ID"})
		return
	}

	var req UpdateTweetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tweet, err := h.tweetUsecase.UpdateTweet(c.Request.Context(), tweetID, user.ID, req.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 更新されたtweetをsub_idを含むレスポンス形式で返す
	tweetResponse, err := h.tweetUsecase.GetTweetByIDWithUser(c.Request.Context(), tweet.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tweetResponse)
}

// DeleteTweet - tweetを削除
func (h *TweetHandler) DeleteTweet(c *gin.Context) {
	subID := c.GetString("sub_id")
	user, err := h.userUsecase.GetUserBySubID(c.Request.Context(), subID)
	if err != nil || user == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		return
	}

	tweetIDStr := c.Param("id")
	tweetID, err := uuid.Parse(tweetIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tweet ID"})
		return
	}

	err = h.tweetUsecase.DeleteTweet(c.Request.Context(), tweetID, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tweet deleted successfully"})
}

// GetTweetsTest - テスト用: 認証なしで全てのtweetを取得
func (h *TweetHandler) GetTweetsTest(c *gin.Context) {
	tweets, err := h.tweetUsecase.GetAllTweetsWithUsers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"tweets": tweets})
}
