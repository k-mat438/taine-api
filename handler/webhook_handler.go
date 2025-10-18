// internal/handler/webhook_handler.go
package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"taine-api/usecase"

	"github.com/gin-gonic/gin"
	svix "github.com/svix/svix-webhooks/go"
)

type WebhookHandler struct {
	UserSvc       usecase.UserService     // SyncMe(subID, name, avatar)
	OrgSvc        usecase.OrganizationSvc // UpsertByExternalID / SoftDeleteByExternalID
	MembershipSvc usecase.MembershipSvc   // UpsertByExternalIDs / SoftDeleteByExternalIDs
}

func NewWebhookHandler(u usecase.UserService, o usecase.OrganizationSvc, m usecase.MembershipSvc) *WebhookHandler {
	return &WebhookHandler{UserSvc: u, OrgSvc: o, MembershipSvc: m}
}

func (h *WebhookHandler) Clerk(c *gin.Context) {
	secret := os.Getenv("CLERK_WEBHOOK_SECRET")
	if secret == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "webhook secret not configured"})
		return
	}

	// 署名ヘッダ（Svix）
	id := c.GetHeader("svix-id")
	ts := c.GetHeader("svix-timestamp")
	sig := c.GetHeader("svix-signature")

	// 生ボディ
	body, err := io.ReadAll(io.LimitReader(c.Request.Body, 1<<20)) // 1MB上限（任意）
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "read body failed"})
		return
	}

	wh, err := svix.NewWebhook(secret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "svix init failed"})
		return
	}

	// Svixのヘッダーを構築
	headers := make(http.Header)
	headers.Set("svix-id", id)
	headers.Set("svix-timestamp", ts)
	headers.Set("svix-signature", sig)

	if err := wh.Verify(body, headers); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "signature verification failed"})
		return
	}

	var evt struct {
		Type string          `json:"type"`
		Data json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(body, &evt); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	switch evt.Type {

	// -------- users --------
	case "user.created", "user.updated":
		var d struct {
			ID        string `json:"id"` // "user_..."
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
			ImageURL  string `json:"image_url"`
		}
		if err := json.Unmarshal(evt.Data, &d); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user payload"})
			return
		}
		name := strings.TrimSpace(d.FirstName + " " + d.LastName)
		if _, err := h.UserSvc.SyncMe(c.Request.Context(), &usecase.AuthClaims{
			SubID:     d.ID,
			Name:      name,
			AvatarURL: d.ImageURL,
		}); err != nil {
			// Webhookは再送されるので5xxでOK（ただしログ推奨）
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user upsert failed"})
			return
		}

	case "user.deleted":
		var d struct {
			ID string `json:"id"`
		}
		if err := json.Unmarshal(evt.Data, &d); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user payload"})
			return
		}
		if err := h.UserSvc.SoftDeleteBySubID(c.Request.Context(), d.ID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user soft delete failed"})
			return
		}

	// -------- organizations --------
	case "organization.created":
		var d struct {
			ID        string `json:"id"`         // Clerk org id
			Name      string `json:"name"`       // Org名
			CreatedBy string `json:"created_by"` // 作成者のClerk user id
			// Slug等が必要なら追加
		}
		if err := json.Unmarshal(evt.Data, &d); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid org payload"})
			return
		}
		// 組織作成時は作成者をownerとして追加
		if err := h.OrgSvc.UpsertByExternalIDWithCreator(c.Request.Context(), d.ID, d.Name, d.CreatedBy); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "org create with owner failed"})
			return
		}

	case "organization.updated":
		var d struct {
			ID   string `json:"id"`   // Clerk org id
			Name string `json:"name"` // Org名
			// Slug等が必要なら追加
		}
		if err := json.Unmarshal(evt.Data, &d); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid org payload"})
			return
		}
		// 組織更新時は組織情報のみ更新
		if err := h.OrgSvc.UpsertByExternalID(c.Request.Context(), d.ID, d.Name); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "org update failed"})
			return
		}

	case "organization.deleted":
		var d struct {
			ID string `json:"id"`
		}
		if err := json.Unmarshal(evt.Data, &d); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid org payload"})
			return
		}
		if err := h.OrgSvc.SoftDeleteByExternalID(c.Request.Context(), d.ID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "org delete failed"})
			return
		}

	// -------- memberships --------
	case "organizationMembership.created", "organizationMembership.updated":
		var d struct {
			Role         string `json:"role"` // e.g. "org:admin" / "org:member"
			Organization struct {
				ID string `json:"id"`
			} `json:"organization"`
			PublicUserData struct {
				UserID string `json:"user_id"`
			} `json:"public_user_data"`
		}
		if err := json.Unmarshal(evt.Data, &d); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid membership payload"})
			return
		}
		if err := h.MembershipSvc.UpsertByExternalIDs(
			c.Request.Context(),
			d.PublicUserData.UserID, // clerk user id
			d.Organization.ID,       // clerk org id
			d.Role,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "membership upsert failed"})
			return
		}

	case "organizationMembership.deleted":
		var d struct {
			Organization struct {
				ID string `json:"id"`
			} `json:"organization"`
			PublicUserData struct {
				UserID string `json:"user_id"`
			} `json:"public_user_data"`
		}
		if err := json.Unmarshal(evt.Data, &d); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid membership payload"})
			return
		}
		if err := h.MembershipSvc.SoftDeleteByExternalIDs(
			c.Request.Context(),
			d.PublicUserData.UserID,
			d.Organization.ID,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "membership delete failed"})
			return
		}

	default:
		// 未対応イベントは200でOK（Clerkの再送を止める）
	}

	// なるべく早く2xxを返す。重い処理は別キューに逃がす設計も可。
	c.JSON(http.StatusOK, gin.H{"status": "ok", "at": time.Now().UTC()})
}
