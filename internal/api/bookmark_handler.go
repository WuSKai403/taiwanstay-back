package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/taiwanstay/taiwanstay-back/internal/service"
)

type BookmarkHandler struct {
	bookmarkService service.BookmarkService
}

func NewBookmarkHandler(bookmarkService service.BookmarkService) *BookmarkHandler {
	return &BookmarkHandler{bookmarkService: bookmarkService}
}

func (h *BookmarkHandler) AddBookmark(c *gin.Context) {
	opportunityID := c.Param("id")
	userID := c.GetString("userID") // Assuming AuthMiddleware sets this

	err := h.bookmarkService.AddBookmark(c.Request.Context(), userID, opportunityID)
	if err != nil {
		if err == service.ErrBookmarkAlreadyExists {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "bookmark added"})
}

func (h *BookmarkHandler) RemoveBookmark(c *gin.Context) {
	opportunityID := c.Param("id")
	userID := c.GetString("userID")

	err := h.bookmarkService.RemoveBookmark(c.Request.Context(), userID, opportunityID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "bookmark removed"})
}

func (h *BookmarkHandler) ListBookmarks(c *gin.Context) {
	userID := c.GetString("userID")
	limitStr := c.DefaultQuery("limit", "20")
	offsetStr := c.DefaultQuery("offset", "0")
	limit, _ := strconv.ParseInt(limitStr, 10, 64)
	offset, _ := strconv.ParseInt(offsetStr, 10, 64)

	bookmarks, total, err := h.bookmarkService.ListUserBookmarks(c.Request.Context(), userID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list bookmarks"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  bookmarks,
		"total": total,
	})
}
