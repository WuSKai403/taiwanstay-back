package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HostStatus string

const (
	HostStatusPending   HostStatus = "PENDING"
	HostStatusActive    HostStatus = "ACTIVE"
	HostStatusInactive  HostStatus = "INACTIVE"
	HostStatusRejected  HostStatus = "REJECTED"
	HostStatusSuspended HostStatus = "SUSPENDED"
	HostStatusEditing   HostStatus = "EDITING"
)

type HostType string

const (
	HostTypeFarm            HostType = "FARM"
	HostTypeHostel          HostType = "HOSTEL"
	HostTypeHomestay        HostType = "HOMESTAY"
	HostTypeEcoVillage      HostType = "ECO_VILLAGE"
	HostTypeRetreatCenter   HostType = "RETREAT_CENTER"
	HostTypeCommunity       HostType = "COMMUNITY"
	HostTypeNGO             HostType = "NGO"
	HostTypeSchool          HostType = "SCHOOL"
	HostTypeCafe            HostType = "CAFE"
	HostTypeRestaurant      HostType = "RESTAURANT"
	HostTypeArtCenter       HostType = "ART_CENTER"
	HostTypeAnimalShelter   HostType = "ANIMAL_SHELTER"
	HostTypeOutdoorActivity HostType = "OUTDOOR_ACTIVITY"
	HostTypeOther           HostType = "OTHER"
	HostTypeCoworkingSpace  HostType = "COWORKING_SPACE"
	HostTypeCulturalVenue   HostType = "CULTURAL_VENUE"
	HostTypeCommunityCenter HostType = "COMMUNITY_CENTER"
)

type HostStatusHistory struct {
	Status     HostStatus         `bson:"status" json:"status"`
	StatusNote string             `bson:"statusNote,omitempty" json:"statusNote,omitempty"`
	UpdatedBy  primitive.ObjectID `bson:"updatedBy,omitempty" json:"updatedBy,omitempty"`
	UpdatedAt  time.Time          `bson:"updatedAt" json:"updatedAt"`
}

type HostSocialMedia struct {
	Facebook  string `bson:"facebook,omitempty" json:"facebook,omitempty"`
	Instagram string `bson:"instagram,omitempty" json:"instagram,omitempty"`
	Line      string `bson:"line,omitempty" json:"line,omitempty"`
	Threads   string `bson:"threads,omitempty" json:"threads,omitempty"`
	Linkedin  string `bson:"linkedin,omitempty" json:"linkedin,omitempty"`
	Twitter   string `bson:"twitter,omitempty" json:"twitter,omitempty"`
	Youtube   string `bson:"youtube,omitempty" json:"youtube,omitempty"`
	Tiktok    string `bson:"tiktok,omitempty" json:"tiktok,omitempty"`
	Other     []struct {
		Name string `bson:"name" json:"name"`
		URL  string `bson:"url" json:"url"`
	} `bson:"other,omitempty" json:"other,omitempty"`
}

type ContactInfo struct {
	ContactEmail  string           `bson:"contactEmail" json:"contactEmail"`
	Phone         string           `bson:"phone,omitempty" json:"phone,omitempty"`
	ContactMobile string           `bson:"contactMobile" json:"contactMobile"`
	Website       string           `bson:"website,omitempty" json:"website,omitempty"`
	ContactHours  string           `bson:"contactHours,omitempty" json:"contactHours,omitempty"`
	Notes         string           `bson:"notes,omitempty" json:"notes,omitempty"`
	SocialMedia   *HostSocialMedia `bson:"socialMedia,omitempty" json:"socialMedia,omitempty"`
}

type GeoJSON struct {
	Type        string    `bson:"type" json:"type" binding:"required,eq=Point"`
	Coordinates []float64 `bson:"coordinates" json:"coordinates" binding:"required,min=2,max=2"`
}

type HostLocation struct {
	Address           string   `bson:"address" json:"address"`
	City              string   `bson:"city" json:"city"`
	District          string   `bson:"district,omitempty" json:"district,omitempty"`
	ZipCode           string   `bson:"zipCode,omitempty" json:"zipCode,omitempty"`
	Country           string   `bson:"country" json:"country"`
	Coordinates       *GeoJSON `bson:"coordinates" json:"coordinates" binding:"required"`
	ShowExactLocation bool     `bson:"showExactLocation" json:"showExactLocation"`
}

type HostPhoto struct {
	PublicID     string `bson:"publicId" json:"publicId"`
	SecureURL    string `bson:"secureUrl" json:"secureUrl"`
	ThumbnailURL string `bson:"thumbnailUrl,omitempty" json:"thumbnailUrl,omitempty"`
	PreviewURL   string `bson:"previewUrl,omitempty" json:"previewUrl,omitempty"`
	OriginalURL  string `bson:"originalUrl,omitempty" json:"originalUrl,omitempty"`
}

type VideoIntroduction struct {
	URL         string `bson:"url" json:"url"`
	Description string `bson:"description,omitempty" json:"description,omitempty"`
}

type AdditionalMedia struct {
	VirtualTour  string     `bson:"virtualTour,omitempty" json:"virtualTour,omitempty"`
	Presentation *HostPhoto `bson:"presentation,omitempty" json:"presentation,omitempty"`
}

type Amenities struct {
	Basics                  map[string]bool `bson:"basics,omitempty" json:"basics,omitempty"`
	Accommodation           map[string]bool `bson:"accommodation,omitempty" json:"accommodation,omitempty"`
	WorkExchange            map[string]bool `bson:"workExchange,omitempty" json:"workExchange,omitempty"`
	Lifestyle               map[string]bool `bson:"lifestyle,omitempty" json:"lifestyle,omitempty"`
	Activities              map[string]bool `bson:"activities,omitempty" json:"activities,omitempty"`
	CustomAmenities         []string        `bson:"customAmenities,omitempty" json:"customAmenities,omitempty"`
	AmenitiesNotes          string          `bson:"amenitiesNotes,omitempty" json:"amenitiesNotes,omitempty"`
	WorkExchangeDescription string          `bson:"workExchangeDescription,omitempty" json:"workExchangeDescription,omitempty"`
}

type HostDetails struct {
	FoundingYear          int      `bson:"foundingYear,omitempty" json:"foundingYear,omitempty"`
	TeamSize              int      `bson:"teamSize,omitempty" json:"teamSize,omitempty"`
	Languages             []string `bson:"languages,omitempty" json:"languages,omitempty"`
	AcceptsChildren       bool     `bson:"acceptsChildren" json:"acceptsChildren"`
	AcceptsPets           bool     `bson:"acceptsPets" json:"acceptsPets"`
	AcceptsCouples        bool     `bson:"acceptsCouples" json:"acceptsCouples"`
	MinStayDuration       int      `bson:"minStayDuration,omitempty" json:"minStayDuration,omitempty"`
	MaxStayDuration       int      `bson:"maxStayDuration,omitempty" json:"maxStayDuration,omitempty"`
	WorkHoursPerWeek      int      `bson:"workHoursPerWeek,omitempty" json:"workHoursPerWeek,omitempty"`
	WorkDaysPerWeek       int      `bson:"workDaysPerWeek,omitempty" json:"workDaysPerWeek,omitempty"`
	ProvidesAccommodation bool     `bson:"providesAccommodation" json:"providesAccommodation"`
	ProvidesMeals         bool     `bson:"providesMeals" json:"providesMeals"`
	DietaryOptions        []string `bson:"dietaryOptions,omitempty" json:"dietaryOptions,omitempty"`
	Rules                 []string `bson:"rules,omitempty" json:"rules,omitempty"`
	Expectations          []string `bson:"expectations,omitempty" json:"expectations,omitempty"`
}

type HostFeatures struct {
	Features    []string `bson:"features,omitempty" json:"features,omitempty"`
	Story       string   `bson:"story,omitempty" json:"story,omitempty"`
	Experience  string   `bson:"experience,omitempty" json:"experience,omitempty"`
	Environment struct {
		Surroundings      string   `bson:"surroundings,omitempty" json:"surroundings,omitempty"`
		Accessibility     string   `bson:"accessibility,omitempty" json:"accessibility,omitempty"`
		NearbyAttractions []string `bson:"nearbyAttractions,omitempty" json:"nearbyAttractions,omitempty"`
	} `bson:"environment,omitempty" json:"environment,omitempty"`
}

type HostRatings struct {
	Overall               float64 `bson:"overall" json:"overall"`
	WorkEnvironment       float64 `bson:"workEnvironment" json:"workEnvironment"`
	Accommodation         float64 `bson:"accommodation" json:"accommodation"`
	Food                  float64 `bson:"food" json:"food"`
	HostHospitality       float64 `bson:"hostHospitality" json:"hostHospitality"`
	LearningOpportunities float64 `bson:"learningOpportunities" json:"learningOpportunities"`
	ReviewCount           int     `bson:"reviewCount" json:"reviewCount"`
}

type Host struct {
	ID                primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	UserID            primitive.ObjectID  `bson:"userId" json:"userId"`
	Name              string              `bson:"name" json:"name"`
	Slug              string              `bson:"slug" json:"slug"`
	Description       string              `bson:"description" json:"description"`
	Status            HostStatus          `bson:"status" json:"status"`
	StatusNote        string              `bson:"statusNote,omitempty" json:"statusNote,omitempty"`
	StatusHistory     []HostStatusHistory `bson:"statusHistory,omitempty" json:"statusHistory,omitempty"`
	Type              HostType            `bson:"type" json:"type"`
	Category          string              `bson:"category" json:"category"`
	Verified          bool                `bson:"verified" json:"verified"`
	VerifiedAt        *time.Time          `bson:"verifiedAt,omitempty" json:"verifiedAt,omitempty"`
	Email             string              `bson:"email" json:"email"`
	Mobile            string              `bson:"mobile" json:"mobile"`
	ContactInfo       ContactInfo         `bson:"contactInfo" json:"contactInfo"`
	Location          HostLocation        `bson:"location" json:"location"`
	Photos            []HostPhoto         `bson:"photos,omitempty" json:"photos,omitempty"`
	PhotoDescriptions []string            `bson:"photoDescriptions,omitempty" json:"photoDescriptions,omitempty"`
	VideoIntroduction *VideoIntroduction  `bson:"videoIntroduction,omitempty" json:"videoIntroduction,omitempty"`
	AdditionalMedia   *AdditionalMedia    `bson:"additionalMedia,omitempty" json:"additionalMedia,omitempty"`
	Amenities         Amenities           `bson:"amenities" json:"amenities"`
	Details           HostDetails         `bson:"details" json:"details"`
	Features          *HostFeatures       `bson:"features,omitempty" json:"features,omitempty"`
	Ratings           HostRatings         `bson:"ratings" json:"ratings"`
	OrganizationID    *primitive.ObjectID `bson:"organizationId,omitempty" json:"organizationId,omitempty"`
	CreatedAt         time.Time           `bson:"createdAt" json:"createdAt"`
	UpdatedAt         time.Time           `bson:"updatedAt" json:"updatedAt"`
}
