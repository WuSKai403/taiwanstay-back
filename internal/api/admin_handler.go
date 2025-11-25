package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/taiwanstay/taiwanstay-back/internal/domain"
	"github.com/taiwanstay/taiwanstay-back/internal/service"
)

type AdminHandler struct {
	adminService service.AdminService
}

func NewAdminHandler(adminService service.AdminService) *AdminHandler {
	return &AdminHandler{adminService: adminService}
}

func (h *AdminHandler) GetStats(c *gin.Context) {
	stats, err := h.adminService.GetSystemStats(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get stats"})
		return
	}
	c.JSON(http.StatusOK, stats)
}

func (h *AdminHandler) ListPendingImages(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "20")
	offsetStr := c.DefaultQuery("offset", "0")
	limit, _ := strconv.ParseInt(limitStr, 10, 64)
	offset, _ := strconv.ParseInt(offsetStr, 10, 64)

	images, total, err := h.adminService.ListPendingImages(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list pending images"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  images,
		"total": total,
	})
}

func (h *AdminHandler) ReviewImage(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Approved bool `json:"approved"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.adminService.ReviewImage(c.Request.Context(), id, req.Approved)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "image reviewed"})
}

func (h *AdminHandler) ListUsers(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "20")
	offsetStr := c.DefaultQuery("offset", "0")
	roleStr := c.Query("role")
	limit, _ := strconv.ParseInt(limitStr, 10, 64)
	offset, _ := strconv.ParseInt(offsetStr, 10, 64)

	users, total, err := h.adminService.ListUsers(c.Request.Context(), domain.UserRole(roleStr), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list users"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  users,
		"total": total,
	})
}

func (h *AdminHandler) UpdateUserStatus(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Status domain.UserStatus `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.adminService.UpdateUserStatus(c.Request.Context(), id, req.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user status updated"})
}
