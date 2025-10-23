package domain

import (
	"time"
)

// 為了與前端對應，我們先定義 Enum 的類型
type UserRole string
type PrivacyLevel string

const (
	RoleUser  UserRole = "USER"
	RoleHost  UserRole = "HOST"
	RoleAdmin UserRole = "ADMIN"

	PrivacyPublic     PrivacyLevel = "PUBLIC"
	PrivacyRegistered PrivacyLevel = "REGISTERED"
	PrivacyPrivate    PrivacyLevel = "PRIVATE"
)

// User 定義了與前端 User.ts 對應的完整使用者模型
type User struct {
	ID              string          `json:"id" bson:"_id,omitempty"`
	Name            string          `json:"name" bson:"name"`
	Email           string          `json:"email" bson:"email"`
	Image           string          `json:"image,omitempty" bson:"image,omitempty"`
	EmailVerified   *time.Time      `json:"emailVerified,omitempty" bson:"emailVerified,omitempty"`
	Password        string          `json:"-" bson:"password,omitempty"`
	Role            UserRole        `json:"role" bson:"role"`
	Profile         Profile         `json:"profile" bson:"profile"`
	HostID          string          `json:"hostId,omitempty" bson:"hostId,omitempty"`
	OrganizationID  string          `json:"organizationId,omitempty" bson:"organizationId,omitempty"`
	PrivacySettings PrivacySettings `json:"privacySettings" bson:"privacySettings"`
	CreatedAt       time.Time       `json:"createdAt" bson:"createdAt"`
	UpdatedAt       time.Time       `json:"updatedAt" bson:"updatedAt"`
}

// Profile 對應前端的 profile 物件
type Profile struct {
	Avatar                  string                   `json:"avatar,omitempty" bson:"avatar,omitempty"`
	Bio                     string                   `json:"bio,omitempty" bson:"bio,omitempty"`
	Skills                  []string                 `json:"skills,omitempty" bson:"skills,omitempty"`
	Languages               []string                 `json:"languages,omitempty" bson:"languages,omitempty"`
	Location                *Location                `json:"location,omitempty" bson:"location,omitempty"`
	SocialMedia             *SocialMedia             `json:"socialMedia,omitempty" bson:"socialMedia,omitempty"`
	PersonalInfo            *PersonalInfo            `json:"personalInfo,omitempty" bson:"personalInfo,omitempty"`
	WorkExchangePreferences *WorkExchangePreferences `json:"workExchangePreferences,omitempty" bson:"workExchangePreferences,omitempty"`
	BirthDate               time.Time                `json:"birthDate" bson:"birthDate"`
	EmergencyContact        EmergencyContact         `json:"emergencyContact" bson:"emergencyContact"`
	WorkExperience          []WorkExperience         `json:"workExperience" bson:"workExperience"`
	PhysicalCondition       string                   `json:"physicalCondition" bson:"physicalCondition"`
	AccommodationNeeds      string                   `json:"accommodationNeeds,omitempty" bson:"accommodationNeeds,omitempty"`
	CulturalInterests       []string                 `json:"culturalInterests,omitempty" bson:"culturalInterests,omitempty"`
	LearningGoals           []string                 `json:"learningGoals,omitempty" bson:"learningGoals,omitempty"`
	PhoneNumber             string                   `json:"phoneNumber,omitempty" bson:"phoneNumber,omitempty"`
	IsPhoneVerified         bool                     `json:"isPhoneVerified" bson:"isPhoneVerified"`
	PreferredWorkHours      int                      `json:"preferredWorkHours" bson:"preferredWorkHours"`
}

// Location 地理位置
type Location struct {
	Type        string    `json:"type" bson:"type"` // "Point"
	Coordinates []float64 `json:"coordinates" bson:"coordinates"`
}

// SocialMedia 社群媒體
type SocialMedia struct {
	Instagram string  `json:"instagram,omitempty" bson:"instagram,omitempty"`
	Facebook  string  `json:"facebook,omitempty" bson:"facebook,omitempty"`
	Threads   string  `json:"threads,omitempty" bson:"threads,omitempty"`
	Linkedin  string  `json:"linkedin,omitempty" bson:"linkedin,omitempty"`
	Twitter   string  `json:"twitter,omitempty" bson:"twitter,omitempty"`
	Youtube   string  `json:"youtube,omitempty" bson:"youtube,omitempty"`
	Tiktok    string  `json:"tiktok,omitempty" bson:"tiktok,omitempty"`
	Website   string  `json:"website,omitempty" bson:"website,omitempty"`
	Other     []Other `json:"other,omitempty" bson:"other,omitempty"`
}

type Other struct {
	Name string `json:"name" bson:"name"`
	URL  string `json:"url" bson:"url"`
}

// PersonalInfo 個人資訊
type PersonalInfo struct {
	Birthdate       *time.Time `json:"birthdate,omitempty" bson:"birthdate,omitempty"`
	Gender          string     `json:"gender,omitempty" bson:"gender,omitempty"`
	Nationality     string     `json:"nationality,omitempty" bson:"nationality,omitempty"`
	CurrentLocation string     `json:"currentLocation,omitempty" bson:"currentLocation,omitempty"`
	Occupation      string     `json:"occupation,omitempty" bson:"occupation,omitempty"`
	Education       string     `json:"education,omitempty" bson:"education,omitempty"`
}

// WorkExchangePreferences 工作換宿偏好
type WorkExchangePreferences struct {
	PreferredWorkTypes  []string   `json:"preferredWorkTypes,omitempty" bson:"preferredWorkTypes,omitempty"`
	PreferredLocations  []string   `json:"preferredLocations,omitempty" bson:"preferredLocations,omitempty"`
	AvailableFrom       *time.Time `json:"availableFrom,omitempty" bson:"availableFrom,omitempty"`
	AvailableTo         *time.Time `json:"availableTo,omitempty" bson:"availableTo,omitempty"`
	MinDuration         int        `json:"minDuration,omitempty" bson:"minDuration,omitempty"`
	MaxDuration         int        `json:"maxDuration,omitempty" bson:"maxDuration,omitempty"`
	HasDriverLicense    bool       `json:"hasDriverLicense,omitempty" bson:"hasDriverLicense,omitempty"`
	DietaryRestrictions []string   `json:"dietaryRestrictions,omitempty" bson:"dietaryRestrictions,omitempty"`
	SpecialNeeds        string     `json:"specialNeeds,omitempty" bson:"specialNeeds,omitempty"`
	Notes               string     `json:"notes,omitempty" bson:"notes,omitempty"`
}

// EmergencyContact 緊急聯絡人
type EmergencyContact struct {
	Name         string `json:"name" bson:"name"`
	Relationship string `json:"relationship" bson:"relationship"`
	Phone        string `json:"phone" bson:"phone"`
	Email        string `json:"email,omitempty" bson:"email,omitempty"`
}

// WorkExperience 工作經驗
type WorkExperience struct {
	Title       string `json:"title" bson:"title"`
	Description string `json:"description,omitempty" bson:"description,omitempty"`
	Duration    string `json:"duration,omitempty" bson:"duration,omitempty"`
}

// PrivacySettings 隱私設定
type PrivacySettings struct {
	Email                   PrivacyLevel `json:"email" bson:"email"`
	Phone                   PrivacyLevel `json:"phone" bson:"phone"`
	PersonalInfo            PrivacyLevel `json:"personalInfo" bson:"personalInfo"`
	SocialMedia             PrivacyLevel `json:"socialMedia" bson:"socialMedia"`
	WorkExchangePreferences PrivacyLevel `json:"workExchangePreferences" bson:"workExchangePreferences"`
	Skills                  PrivacyLevel `json:"skills" bson:"skills"`
	Languages               PrivacyLevel `json:"languages" bson:"languages"`
	Bio                     PrivacyLevel `json:"bio" bson:"bio"`
}
