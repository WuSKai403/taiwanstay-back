package api

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/taiwanstay/taiwanstay-back/internal/domain"
	"github.com/taiwanstay/taiwanstay-back/internal/service"
)

type ImageHandler struct {
	imageService service.ImageService
}

func NewImageHandler(imageService service.ImageService) *ImageHandler {
	return &ImageHandler{imageService: imageService}
}

func (h *ImageHandler) Upload(c *gin.Context) {
	// 1. Get User ID from context
	claims, exists := c.Get("userClaims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	mapClaims := claims.(jwt.MapClaims)
	userID := mapClaims["sub"].(string)

	// 2. Get File
	file, header, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "image file is required"})
		return
	}
	defer file.Close()

	// 3. Call Service
	image, err := h.imageService.UploadImage(c.Request.Context(), file, header, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upload image"})
		return
	}

	c.JSON(http.StatusCreated, image)
}

func (h *ImageHandler) GetPrivateImage(c *gin.Context) {
	id := c.Param("id")
	// TODO: Check permission (Owner or Admin)

	content, err := h.imageService.GetImageContent(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "image not found"})
		return
	}
	defer content.Close()

	_, err = io.Copy(c.Writer, content)
	if err != nil {
		// Log error
	}
}

func (h *ImageHandler) UpdateStatus(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Status string `json:"status" binding:"required,oneof=PENDING APPROVED REJECTED"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.imageService.UpdateImageStatus(c.Request.Context(), id, domain.ImageStatus(req.Status))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "status updated"})
}
