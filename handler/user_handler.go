package handler

import (
	"net/http"
	"taine-api/domain"
	"taine-api/usecase"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userUsecase usecase.UserUsecase
}

func NewUserHandler(userUsecase usecase.UserUsecase) *UserHandler {
	return &UserHandler{userUsecase: userUsecase}
}

func (h *UserHandler) UpsertUser(c *gin.Context) {
	subID := c.GetString("sub_id")
	user, err := h.userUsecase.UpsertUser(c.Request.Context(), &domain.User{
		SubID: subID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) GetUserBySubID(c *gin.Context) {
	subID := c.GetString("sub_id")
	user, err := h.userUsecase.GetUserBySubID(c.Request.Context(), subID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	subID := c.GetString("sub_id")
	user, err := h.userUsecase.GetUserBySubID(c.Request.Context(), subID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	subID := c.GetString("sub_id")
	user, err := h.userUsecase.GetUserBySubID(c.Request.Context(), subID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	err = h.userUsecase.DeleteUser(c.Request.Context(), user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
