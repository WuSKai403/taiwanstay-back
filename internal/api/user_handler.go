package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/taiwanstay/taiwanstay-back/internal/service"
)

// UserHandler 負責處理與使用者相關的 HTTP 請求
type UserHandler struct {
	userService service.UserService
}

// NewUserHandler 建立一個新的 UserHandler 實例
func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// RegisterRequest 定義了註冊請求的資料結構
type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// Register 處理使用者註冊請求
func (h *UserHandler) Register(c *gin.Context) {
	var req RegisterRequest
	// 1. 綁定並驗證請求資料
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 2. 呼叫 Service 層執行業務邏輯
	user, err := h.userService.RegisterUser(c.Request.Context(), req.Name, req.Email, req.Password)
	if err != nil {
		// TODO: 根據錯誤類型回傳不同的 HTTP 狀態碼
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 3. 回傳成功的響應
	c.JSON(http.StatusCreated, user)
}
