package api

import (
	"errors"
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

// LoginRequest 定義了登入請求的資料結構，以支援多種登入方式
type LoginRequest struct {
	LoginType string `json:"loginType" binding:"required,oneof=password google apple"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password"` // 密碼登入時為必需
	Token     string `json:"token"`    // OAuth 登入時為必需
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
		if errors.Is(err, service.ErrEmailAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"error": "email already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to register user"})
		return
	}

	// 3. 回傳成功的響應
	c.JSON(http.StatusCreated, user)
}

// Login 處理使用者登入請求，支援多種登入方式
func (h *UserHandler) Login(c *gin.Context) {
	var req LoginRequest
	// 1. 綁定並驗證請求資料
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 2. 根據登入類型執行不同邏輯
	switch req.LoginType {
	case "password":
		// 驗證密碼是否存在
		if req.Password == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "password is required for this login type"})
			return
		}
		h.handlePasswordLogin(c, req)
	case "google":
		// TODO: 實作 Google OAuth 登入
		c.JSON(http.StatusNotImplemented, gin.H{"error": "Google login not implemented yet"})
	case "apple":
		// TODO: 實作 Apple OAuth 登入
		c.JSON(http.StatusNotImplemented, gin.H{"error": "Apple login not implemented yet"})
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid login type"})
	}
}

// handlePasswordLogin 處理傳統的密碼登入
func (h *UserHandler) handlePasswordLogin(c *gin.Context, req LoginRequest) {
	user, token, err := h.userService.LoginUser(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to login user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user":  user,
		"token": token,
	})
}
