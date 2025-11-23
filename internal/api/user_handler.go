package api

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/taiwanstay/taiwanstay-back/internal/domain"
	"github.com/taiwanstay/taiwanstay-back/internal/service"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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

// UpdateUserRequest 定義了更新使用者資訊請求的資料結構。
// 使用指標類型 (*string, *int) 來區分「未提供」和「提供空值」的情況。
// 對於 slice 和 struct，如果請求中未包含該 key，它們的值會是 nil。
type UpdateUserRequest struct {
	Name                    *string                         `json:"name,omitempty"`
	Avatar                  *string                         `json:"avatar,omitempty"`
	Bio                     *string                         `json:"bio,omitempty"`
	Skills                  []string                        `json:"skills"`
	Languages               []string                        `json:"languages"`
	Location                *domain.Location                `json:"location,omitempty"`
	SocialMedia             *domain.SocialMedia             `json:"socialMedia,omitempty"`
	PersonalInfo            *domain.PersonalInfo            `json:"personalInfo,omitempty"`
	WorkExchangePreferences *domain.WorkExchangePreferences `json:"workExchangePreferences,omitempty"`
	BirthDate               *time.Time                      `json:"birthDate,omitempty"`
	EmergencyContact        *domain.EmergencyContact        `json:"emergencyContact,omitempty"`
	WorkExperience          []domain.WorkExperience         `json:"workExperience"`
	PhysicalCondition       *string                         `json:"physicalCondition,omitempty"`
	AccommodationNeeds      *string                         `json:"accommodationNeeds,omitempty"`
	CulturalInterests       []string                        `json:"culturalInterests"`
	LearningGoals           []string                        `json:"learningGoals"`
	PhoneNumber             *string                         `json:"phoneNumber,omitempty"`
	Address                 *string                         `json:"address,omitempty"`
	PreferredWorkHours      *int                            `json:"preferredWorkHours,omitempty"`
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

// Logout 處理使用者登出請求
func (h *UserHandler) Logout(c *gin.Context) {
	// 呼叫 Service 層執行登出邏輯
	err := h.userService.LogoutUser(c.Request.Context())
	if err != nil {
		// 理論上目前版本的 LogoutUser 不會回傳錯誤，但保留以備未來擴充
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to logout user"})
		return
	}

	// 回傳成功的響應
	c.JSON(http.StatusOK, gin.H{"message": "logout successful"})
}

// GetAllUsers 處理取得所有使用者的請求 (僅限管理員)
func (h *UserHandler) GetAllUsers(c *gin.Context) {
	users, err := h.userService.GetAllUsers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get users"})
		return
	}

	c.JSON(http.StatusOK, users)
}

// GetUserByID 處理根據 ID 取得單一使用者的請求 (僅限管理員)
func (h *UserHandler) GetUserByID(c *gin.Context) {
	userID := c.Param("id")

	user, err := h.userService.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateMe 處理更新當前登入使用者資訊的請求
func (h *UserHandler) UpdateMe(c *gin.Context) {
	// 1. 從 context 取得使用者 ID
	claims, _ := c.Get("userClaims")
	mapClaims := claims.(jwt.MapClaims)
	userID := mapClaims["sub"].(string)

	// 2. 綁定請求資料
	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 3. 動態建立更新 payload
	payload := bson.M{}

	if req.Name != nil {
		payload["name"] = *req.Name
	}
	// 注意：對於 Profile 內嵌結構，我們使用點表示法來更新特定欄位
	if req.Avatar != nil {
		payload["profile.avatar"] = *req.Avatar
	}
	if req.Bio != nil {
		payload["profile.bio"] = *req.Bio
	}
	if req.Skills != nil {
		payload["profile.skills"] = req.Skills
	}
	if req.Languages != nil {
		payload["profile.languages"] = req.Languages
	}
	if req.Location != nil {
		payload["profile.location"] = req.Location
	}
	if req.SocialMedia != nil {
		payload["profile.socialMedia"] = req.SocialMedia
	}
	if req.PersonalInfo != nil {
		payload["profile.personalInfo"] = req.PersonalInfo
	}
	if req.WorkExchangePreferences != nil {
		payload["profile.workExchangePreferences"] = req.WorkExchangePreferences
	}
	if req.BirthDate != nil {
		payload["profile.birthDate"] = *req.BirthDate
	}
	if req.EmergencyContact != nil {
		payload["profile.emergencyContact"] = *req.EmergencyContact
	}
	if req.WorkExperience != nil {
		payload["profile.workExperience"] = req.WorkExperience
	}
	if req.PhysicalCondition != nil {
		payload["profile.physicalCondition"] = *req.PhysicalCondition
	}
	if req.AccommodationNeeds != nil {
		payload["profile.accommodationNeeds"] = *req.AccommodationNeeds
	}
	if req.CulturalInterests != nil {
		payload["profile.culturalInterests"] = req.CulturalInterests
	}
	if req.LearningGoals != nil {
		payload["profile.learningGoals"] = req.LearningGoals
	}
	if req.PhoneNumber != nil {
		payload["profile.phoneNumber"] = *req.PhoneNumber
	}
	if req.Address != nil {
		payload["profile.address"] = *req.Address
	}
	if req.PreferredWorkHours != nil {
		payload["profile.preferredWorkHours"] = *req.PreferredWorkHours
	}

	// 如果沒有任何欄位需要更新，直接回傳成功
	if len(payload) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "no fields to update"})
		return
	}

	// 4. 呼叫 Service 層執行更新
	updatedUser, err := h.userService.UpdateUser(c.Request.Context(), userID, payload)
	if err != nil {
		// 這邊可以根據 service 回傳的 error 類型做更細緻的處理
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user"})
		return
	}

	// 5. 回傳更新後的使用者資訊
	c.JSON(http.StatusOK, updatedUser)
}

// GetMe 處理取得當前登入使用者資訊的請求
func (h *UserHandler) GetMe(c *gin.Context) {
	// 從 context 中取得 AuthMiddleware 注入的 userClaims
	claims, exists := c.Get("userClaims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user claims not found"})
		return
	}

	mapClaims, ok := claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid claims format"})
		return
	}

	// 從 claims 中取得使用者 ID (sub)
	userID, ok := mapClaims["sub"].(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user id in token"})
		return
	}

	// 複用 GetUserByID 服務
	user, err := h.userService.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user information"})
		return
	}

	c.JSON(http.StatusOK, user)
}
