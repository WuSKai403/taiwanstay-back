package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OpportunityStatus string

const (
	OpportunityStatusDraft       OpportunityStatus = "DRAFT"
	OpportunityStatusPending     OpportunityStatus = "PENDING"
	OpportunityStatusActive      OpportunityStatus = "ACTIVE"
	OpportunityStatusPaused      OpportunityStatus = "PAUSED"
	OpportunityStatusExpired     OpportunityStatus = "EXPIRED"
	OpportunityStatusFilled      OpportunityStatus = "FILLED"
	OpportunityStatusRejected    OpportunityStatus = "REJECTED"
	OpportunityStatusAdminPaused OpportunityStatus = "ADMIN_PAUSED"
)

type OpportunityType string

const (
	OpportunityTypeFarming          OpportunityType = "FARMING"
	OpportunityTypeGardening        OpportunityType = "GARDENING"
	OpportunityTypeAnimalCare       OpportunityType = "ANIMAL_CARE"
	OpportunityTypeConstruction     OpportunityType = "CONSTRUCTION"
	OpportunityTypeHospitality      OpportunityType = "HOSPITALITY"
	OpportunityTypeCooking          OpportunityType = "COOKING"
	OpportunityTypeCleaning         OpportunityType = "CLEANING"
	OpportunityTypeChildcare        OpportunityType = "CHILDCARE"
	OpportunityTypeElderlyCare      OpportunityType = "ELDERLY_CARE"
	OpportunityTypeTeaching         OpportunityType = "TEACHING"
	OpportunityTypeLanguageExchange OpportunityType = "LANGUAGE_EXCHANGE"
	OpportunityTypeCreative         OpportunityType = "CREATIVE"
	OpportunityTypeDigitalNomad     OpportunityType = "DIGITAL_NOMAD"
	OpportunityTypeAdministration   OpportunityType = "ADMINISTRATION"
	OpportunityTypeMaintenance      OpportunityType = "MAINTENANCE"
	OpportunityTypeTourism          OpportunityType = "TOURISM"
	OpportunityTypeConservation     OpportunityType = "CONSERVATION"
	OpportunityTypeCommunity        OpportunityType = "COMMUNITY"
	OpportunityTypeEvent            OpportunityType = "EVENT"
	OpportunityTypeOther            OpportunityType = "OTHER"
)

type TimeSlotStatus string

const (
	TimeSlotStatusOpen   TimeSlotStatus = "OPEN"
	TimeSlotStatusFilled TimeSlotStatus = "FILLED"
	TimeSlotStatusClosed TimeSlotStatus = "CLOSED"
)

type OpportunityStatusHistory struct {
	Status    OpportunityStatus  `bson:"status" json:"status"`
	Reason    string             `bson:"reason,omitempty" json:"reason,omitempty"`
	ChangedBy primitive.ObjectID `bson:"changedBy,omitempty" json:"changedBy,omitempty"`
	ChangedAt time.Time          `bson:"changedAt" json:"changedAt"`
}

type WorkDetails struct {
	Tasks                 []string `bson:"tasks" json:"tasks"`
	Skills                []string `bson:"skills,omitempty" json:"skills,omitempty"`
	LearningOpportunities []string `bson:"learningOpportunities,omitempty" json:"learningOpportunities,omitempty"`
	PhysicalDemand        string   `bson:"physicalDemand" json:"physicalDemand"` // low, medium, high
	Languages             []string `bson:"languages" json:"languages"`
	AvailableMonths       []int    `bson:"availableMonths,omitempty" json:"availableMonths,omitempty"`
}

type Benefits struct {
	Accommodation struct {
		Provided    bool   `bson:"provided" json:"provided"`
		Type        string `bson:"type,omitempty" json:"type,omitempty"` // private_room, shared_room, etc.
		Description string `bson:"description,omitempty" json:"description,omitempty"`
	} `bson:"accommodation" json:"accommodation"`
	Meals struct {
		Provided    bool   `bson:"provided" json:"provided"`
		Count       int    `bson:"count,omitempty" json:"count,omitempty"`
		Description string `bson:"description,omitempty" json:"description,omitempty"`
	} `bson:"meals" json:"meals"`
	Stipend struct {
		Provided  bool    `bson:"provided" json:"provided"`
		Amount    float64 `bson:"amount,omitempty" json:"amount,omitempty"`
		Currency  string  `bson:"currency,omitempty" json:"currency,omitempty"`
		Frequency string  `bson:"frequency,omitempty" json:"frequency,omitempty"` // daily, weekly, monthly
	} `bson:"stipend" json:"stipend"`
	OtherBenefits []string `bson:"otherBenefits,omitempty" json:"otherBenefits,omitempty"`
}

type Requirements struct {
	MinAge          int    `bson:"minAge,omitempty" json:"minAge,omitempty"`
	MaxAge          int    `bson:"maxAge,omitempty" json:"maxAge,omitempty"`
	Gender          string `bson:"gender" json:"gender"` // any, male, female
	AcceptsCouples  bool   `bson:"acceptsCouples" json:"acceptsCouples"`
	AcceptsFamilies bool   `bson:"acceptsFamilies" json:"acceptsFamilies"`
	AcceptsPets     bool   `bson:"acceptsPets" json:"acceptsPets"`
	DrivingLicense  struct {
		CarRequired        bool   `bson:"carRequired" json:"carRequired"`
		MotorcycleRequired bool   `bson:"motorcycleRequired" json:"motorcycleRequired"`
		OtherRequired      bool   `bson:"otherRequired" json:"otherRequired"`
		OtherDescription   string `bson:"otherDescription,omitempty" json:"otherDescription,omitempty"`
	} `bson:"drivingLicense" json:"drivingLicense"`
	SpecificNationalities []string `bson:"specificNationalities,omitempty" json:"specificNationalities,omitempty"`
	SpecificSkills        []string `bson:"specificSkills,omitempty" json:"specificSkills,omitempty"`
	OtherRequirements     []string `bson:"otherRequirements,omitempty" json:"otherRequirements,omitempty"`
}

type OpportunityMedia struct {
	CoverImage       *HostPhoto  `bson:"coverImage,omitempty" json:"coverImage,omitempty"`
	Images           []HostPhoto `bson:"images,omitempty" json:"images,omitempty"`
	Descriptions     []string    `bson:"descriptions,omitempty" json:"descriptions,omitempty"`
	VideoURL         string      `bson:"videoUrl,omitempty" json:"videoUrl,omitempty"`
	VideoDescription string      `bson:"videoDescription,omitempty" json:"videoDescription,omitempty"`
	VirtualTour      string      `bson:"virtualTour,omitempty" json:"virtualTour,omitempty"`
}

type OpportunityLocation struct {
	Address     string   `bson:"address,omitempty" json:"address,omitempty"`
	City        string   `bson:"city" json:"city"`
	District    string   `bson:"district,omitempty" json:"district,omitempty"`
	Country     string   `bson:"country" json:"country"`
	Coordinates *GeoJSON `bson:"coordinates,omitempty" json:"coordinates,omitempty"`
}

type ApplicationProcess struct {
	Instructions        string    `bson:"instructions,omitempty" json:"instructions,omitempty"`
	Questions           []string  `bson:"questions,omitempty" json:"questions,omitempty"`
	Deadline            time.Time `bson:"deadline,omitempty" json:"deadline,omitempty"`
	MaxApplications     int       `bson:"maxApplications,omitempty" json:"maxApplications,omitempty"`
	CurrentApplications int       `bson:"currentApplications" json:"currentApplications"`
}

type Impact struct {
	EnvironmentalContribution   string   `bson:"environmentalContribution,omitempty" json:"environmentalContribution,omitempty"`
	SocialContribution          string   `bson:"socialContribution,omitempty" json:"socialContribution,omitempty"`
	CulturalExchange            string   `bson:"culturalExchange,omitempty" json:"culturalExchange,omitempty"`
	SustainableDevelopmentGoals []string `bson:"sustainableDevelopmentGoals,omitempty" json:"sustainableDevelopmentGoals,omitempty"`
}

type OpportunityStats struct {
	Views        int `bson:"views" json:"views"`
	Applications int `bson:"applications" json:"applications"`
	Bookmarks    int `bson:"bookmarks" json:"bookmarks"`
	Shares       int `bson:"shares" json:"shares"`
}

type CapacityOverride struct {
	StartDate time.Time `bson:"startDate" json:"startDate"`
	EndDate   time.Time `bson:"endDate" json:"endDate"`
	Capacity  int       `bson:"capacity" json:"capacity"`
}

type MonthlyCapacity struct {
	Month       string `bson:"month" json:"month"` // YYYY-MM
	Capacity    int    `bson:"capacity" json:"capacity"`
	BookedCount int    `bson:"bookedCount" json:"bookedCount"`
}

type TimeSlot struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	StartDate         string             `bson:"startDate" json:"startDate"` // YYYY-MM-DD
	EndDate           string             `bson:"endDate" json:"endDate"`     // YYYY-MM-DD
	DefaultCapacity   int                `bson:"defaultCapacity" json:"defaultCapacity"`
	MinimumStay       int                `bson:"minimumStay" json:"minimumStay"`
	WorkDaysPerWeek   int                `bson:"workDaysPerWeek" json:"workDaysPerWeek"`
	WorkHoursPerDay   int                `bson:"workHoursPerDay" json:"workHoursPerDay"`
	AppliedCount      int                `bson:"appliedCount" json:"appliedCount"`
	ConfirmedCount    int                `bson:"confirmedCount" json:"confirmedCount"`
	Status            TimeSlotStatus     `bson:"status" json:"status"`
	Description       string             `bson:"description,omitempty" json:"description,omitempty"`
	CapacityOverrides []CapacityOverride `bson:"capacityOverrides,omitempty" json:"capacityOverrides,omitempty"`
	MonthlyCapacities []MonthlyCapacity  `bson:"monthlyCapacities,omitempty" json:"monthlyCapacities,omitempty"`
}

type Opportunity struct {
	ID                 primitive.ObjectID         `bson:"_id,omitempty" json:"id"`
	HostID             primitive.ObjectID         `bson:"hostId" json:"hostId"`
	Title              string                     `bson:"title" json:"title"`
	Slug               string                     `bson:"slug" json:"slug"`
	PublicID           string                     `bson:"publicId" json:"publicId"`
	Description        string                     `bson:"description" json:"description"`
	ShortDescription   string                     `bson:"shortDescription" json:"shortDescription"`
	Status             OpportunityStatus          `bson:"status" json:"status"`
	StatusNote         string                     `bson:"statusNote,omitempty" json:"statusNote,omitempty"`
	Type               OpportunityType            `bson:"type" json:"type"`
	StatusHistory      []OpportunityStatusHistory `bson:"statusHistory,omitempty" json:"statusHistory,omitempty"`
	WorkDetails        WorkDetails                `bson:"workDetails" json:"workDetails"`
	Benefits           Benefits                   `bson:"benefits" json:"benefits"`
	Requirements       Requirements               `bson:"requirements" json:"requirements"`
	Media              OpportunityMedia           `bson:"media" json:"media"`
	Location           OpportunityLocation        `bson:"location" json:"location"`
	ApplicationProcess ApplicationProcess         `bson:"applicationProcess" json:"applicationProcess"`
	Impact             Impact                     `bson:"impact" json:"impact"`
	Ratings            HostRatings                `bson:"ratings" json:"ratings"` // Reusing HostRatings as structure is same
	Stats              OpportunityStats           `bson:"stats" json:"stats"`
	TimeSlots          []TimeSlot                 `bson:"timeSlots,omitempty" json:"timeSlots,omitempty"`
	HasTimeSlots       bool                       `bson:"hasTimeSlots" json:"hasTimeSlots"`
	CreatedAt          time.Time                  `bson:"createdAt" json:"createdAt"`
	UpdatedAt          time.Time                  `bson:"updatedAt" json:"updatedAt"`
}
