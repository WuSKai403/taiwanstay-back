package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/taiwanstay/taiwanstay-back/internal/service"
)

type NotificationHandler struct {
	notifService service.NotificationService
}

func NewNotificationHandler(notifService service.NotificationService) *NotificationHandler {
	return &NotificationHandler{notifService: notifService}
}

func (h *NotificationHandler) List(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "20")
	offsetStr := c.DefaultQuery("offset", "0")
	limit, _ := strconv.ParseInt(limitStr, 10, 64)
	offset, _ := strconv.ParseInt(offsetStr, 10, 64)

	claims, _ := c.Get("userClaims")
	mapClaims := claims.(jwt.MapClaims)
	userID := mapClaims["sub"].(string)

	notifs, total, err := h.notifService.ListNotifications(c.Request.Context(), userID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list notifications"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  notifs,
		"total": total,
	})
}

func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	id := c.Param("id")
	claims, _ := c.Get("userClaims")
	mapClaims := claims.(jwt.MapClaims)
	userID := mapClaims["sub"].(string)

	err := h.notifService.MarkAsRead(c.Request.Context(), id, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "marked as read"})
}

func (h *NotificationHandler) MarkAllAsRead(c *gin.Context) {
	claims, _ := c.Get("userClaims")
	mapClaims := claims.(jwt.MapClaims)
	userID := mapClaims["sub"].(string)

	err := h.notifService.MarkAllAsRead(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "all marked as read"})
}
