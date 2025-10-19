package handler

import (
	"net/http"
	"taine-api/usecase"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type WishHandler struct {
	wishSvc usecase.WishSvc
	userSvc usecase.UserUsecase
}

func NewWishHandler(wishSvc usecase.WishSvc, userSvc usecase.UserUsecase) *WishHandler {
	return &WishHandler{
		wishSvc: wishSvc,
		userSvc: userSvc,
	}
}

type CreateWishRequest struct {
	Title   string `json:"title" binding:"required"`
	Note    string `json:"note"`
	OrderNo int    `json:"order_no"`
}

type UpdateWishRequest struct {
	Title   string `json:"title" binding:"required"`
	Note    string `json:"note"`
	OrderNo int    `json:"order_no"`
}

type UpdateWishOrderRequest struct {
	OrderNo int `json:"order_no" binding:"required"`
}

type WishResponse struct {
	ID             string  `json:"id"`
	OrganizationID string  `json:"organization_id"`
	Title          string  `json:"title"`
	Note           string  `json:"note"`
	OrderNo        int     `json:"order_no"`
	CreatedAt      string  `json:"created_at"`
	UpdatedAt      string  `json:"updated_at"`
	DeletedAt      *string `json:"deleted_at,omitempty"`
}

// CreateWish - 新しいWishを作成
func (h *WishHandler) CreateWish(c *gin.Context) {
	var req CreateWishRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 組織IDをパース
	orgID := c.GetString("org_external_id")

	// Wishを作成
	wish, err := h.wishSvc.CreateWishByOrganizationExternalID(c.Request.Context(), orgID, req.Title, req.Note, req.OrderNo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// レスポンスを作成
	response := WishResponse{
		ID:             wish.ID.String(),
		OrganizationID: wish.OrganizationID.String(),
		Title:          wish.Title,
		Note:           wish.Note,
		OrderNo:        wish.OrderNo,
		CreatedAt:      wish.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:      wish.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}

	c.JSON(http.StatusCreated, response)
}

// GetWish - 特定のWishを取得
func (h *WishHandler) GetWish(c *gin.Context) {
	wishIDStr := c.Param("id")
	wishID, err := uuid.Parse(wishIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid wish ID"})
		return
	}

	wish, err := h.wishSvc.GetWish(c.Request.Context(), wishID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// レスポンスを作成
	response := WishResponse{
		ID:             wish.ID.String(),
		OrganizationID: wish.OrganizationID.String(),
		Title:          wish.Title,
		Note:           wish.Note,
		OrderNo:        wish.OrderNo,
		CreatedAt:      wish.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:      wish.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
	if wish.DeletedAt != nil {
		deletedAtStr := wish.DeletedAt.Format("2006-01-02T15:04:05Z")
		response.DeletedAt = &deletedAtStr
	}

	c.JSON(http.StatusOK, response)
}

// GetWishesByOrganization - 組織のWish一覧を取得
func (h *WishHandler) GetWishesByOrganization(c *gin.Context) {
	orgIDStr := c.Param("org_id")

	// UUIDとして解析を試みる
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		// UUID形式でない場合は、external_idとして扱う
		// external_idから実際の組織IDを取得する必要がある
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID format - UUID required"})
		return
	}

	wishes, err := h.wishSvc.GetWishesByOrganization(c.Request.Context(), orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// レスポンスを作成
	responses := make([]WishResponse, len(wishes))
	for i, wish := range wishes {
		responses[i] = WishResponse{
			ID:             wish.ID.String(),
			OrganizationID: wish.OrganizationID.String(),
			Title:          wish.Title,
			Note:           wish.Note,
			OrderNo:        wish.OrderNo,
			CreatedAt:      wish.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt:      wish.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		}
		if wish.DeletedAt != nil {
			deletedAtStr := wish.DeletedAt.Format("2006-01-02T15:04:05Z")
			responses[i].DeletedAt = &deletedAtStr
		}
	}

	c.JSON(http.StatusOK, gin.H{"wishes": responses})
}

// GetWishesForCurrentOrg - 現在のユーザーの組織のWish一覧を取得（JWTのorg_idを使用）
func (h *WishHandler) GetWishesForCurrentOrg(c *gin.Context) {
	orgExternalID := c.GetString("org_external_id")

	if orgExternalID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Organization ID not found in token"})
		return
	}

	// external_idから実際の組織を取得してWishを取得
	wishes, err := h.wishSvc.GetWishesByOrganizationExternalID(c.Request.Context(), orgExternalID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// レスポンスを作成
	responses := make([]WishResponse, len(wishes))
	for i, wish := range wishes {
		responses[i] = WishResponse{
			ID:             wish.ID.String(),
			OrganizationID: wish.OrganizationID.String(),
			Title:          wish.Title,
			Note:           wish.Note,
			OrderNo:        wish.OrderNo,
			CreatedAt:      wish.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt:      wish.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		}
		if wish.DeletedAt != nil {
			deletedAtStr := wish.DeletedAt.Format("2006-01-02T15:04:05Z")
			responses[i].DeletedAt = &deletedAtStr
		}
	}

	c.JSON(http.StatusOK, gin.H{"wishes": responses})
}

// UpdateWish - Wishを更新
func (h *WishHandler) UpdateWish(c *gin.Context) {
	wishIDStr := c.Param("id")
	wishID, err := uuid.Parse(wishIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid wish ID"})
		return
	}

	var req UpdateWishRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	wish, err := h.wishSvc.UpdateWish(c.Request.Context(), wishID, req.Title, req.Note, req.OrderNo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// レスポンスを作成
	response := WishResponse{
		ID:             wish.ID.String(),
		OrganizationID: wish.OrganizationID.String(),
		Title:          wish.Title,
		Note:           wish.Note,
		OrderNo:        wish.OrderNo,
		CreatedAt:      wish.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:      wish.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}

	c.JSON(http.StatusOK, response)
}

// DeleteWish - Wishを物理削除
func (h *WishHandler) DeleteWish(c *gin.Context) {
	wishIDStr := c.Param("id")
	wishID, err := uuid.Parse(wishIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid wish ID"})
		return
	}

	err = h.wishSvc.DeleteWish(c.Request.Context(), wishID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Wish deleted successfully"})
}

// SoftDeleteWish - Wishをソフトデリート
func (h *WishHandler) SoftDeleteWish(c *gin.Context) {
	wishIDStr := c.Param("id")
	wishID, err := uuid.Parse(wishIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid wish ID"})
		return
	}

	err = h.wishSvc.SoftDeleteWish(c.Request.Context(), wishID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Wish soft deleted successfully"})
}

// RestoreWish - ソフトデリートされたWishを復元
func (h *WishHandler) RestoreWish(c *gin.Context) {
	wishIDStr := c.Param("id")
	wishID, err := uuid.Parse(wishIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid wish ID"})
		return
	}

	err = h.wishSvc.RestoreWish(c.Request.Context(), wishID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Wish restored successfully"})
}

// UpdateWishOrder - Wishの並び順を更新
func (h *WishHandler) UpdateWishOrder(c *gin.Context) {
	wishIDStr := c.Param("id")
	wishID, err := uuid.Parse(wishIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid wish ID"})
		return
	}

	var req UpdateWishOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.wishSvc.UpdateWishOrder(c.Request.Context(), wishID, req.OrderNo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Wish order updated successfully"})
}
