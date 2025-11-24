package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/taiwanstay/taiwanstay-back/internal/domain"
	"github.com/taiwanstay/taiwanstay-back/internal/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HostHandler struct {
	hostService service.HostService
}

func NewHostHandler(hostService service.HostService) *HostHandler {
	return &HostHandler{hostService: hostService}
}

func (h *HostHandler) Create(c *gin.Context) {
	var host domain.Host
	if err := c.ShouldBindJSON(&host); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get User ID from context
	claims, exists := c.Get("userClaims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	mapClaims := claims.(jwt.MapClaims)
	userID := mapClaims["sub"].(string)
	userObjID, _ := primitive.ObjectIDFromHex(userID)
	host.UserID = userObjID

	createdHost, err := h.hostService.CreateHost(c.Request.Context(), &host)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create host"})
		return
	}

	c.JSON(http.StatusCreated, createdHost)
}

func (h *HostHandler) GetMe(c *gin.Context) {
	// Get User ID from context
	claims, exists := c.Get("userClaims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	mapClaims := claims.(jwt.MapClaims)
	userID := mapClaims["sub"].(string)

	host, err := h.hostService.GetHostByUserID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "host profile not found"})
		return
	}

	c.JSON(http.StatusOK, host)
}

func (h *HostHandler) UpdateMe(c *gin.Context) {
	// Get User ID from context
	claims, exists := c.Get("userClaims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	mapClaims := claims.(jwt.MapClaims)
	userID := mapClaims["sub"].(string)

	// First get existing host to ensure ownership
	existingHost, err := h.hostService.GetHostByUserID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "host profile not found"})
		return
	}

	var updateData domain.Host
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ensure ID matches
	updateData.ID = existingHost.ID
	updateData.UserID = existingHost.UserID

	err = h.hostService.UpdateHost(c.Request.Context(), existingHost.ID.Hex(), &updateData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update host"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "host updated"})
}
