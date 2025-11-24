package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ApplicationStatus string

const (
	ApplicationStatusDraft     ApplicationStatus = "DRAFT"
	ApplicationStatusPending   ApplicationStatus = "PENDING"
	ApplicationStatusAccepted  ApplicationStatus = "ACCEPTED"
	ApplicationStatusRejected  ApplicationStatus = "REJECTED"
	ApplicationStatusCancelled ApplicationStatus = "CANCELLED"
)

type ApplicationDetails struct {
	Message       string `bson:"message" json:"message"`
	StartDate     string `bson:"startDate" json:"startDate"` // YYYY-MM-DD
	EndDate       string `bson:"endDate" json:"endDate"`     // YYYY-MM-DD
	Duration      int    `bson:"duration" json:"duration"`   // Days
	TravelingWith struct {
		Partner  bool `bson:"partner" json:"partner"`
		Children bool `bson:"children" json:"children"`
		Pets     bool `bson:"pets" json:"pets"`
	} `bson:"travelingWith" json:"travelingWith"`
	Languages          []string `bson:"languages" json:"languages"`
	RelevantExperience string   `bson:"relevantExperience" json:"relevantExperience"`
}

type ReviewDetails struct {
	ReviewedBy primitive.ObjectID `bson:"reviewedBy,omitempty" json:"reviewedBy,omitempty"`
	ReviewedAt time.Time          `bson:"reviewedAt,omitempty" json:"reviewedAt,omitempty"`
	Notes      string             `bson:"notes,omitempty" json:"notes,omitempty"`
	Rating     int                `bson:"rating,omitempty" json:"rating,omitempty"`
}

type Application struct {
	ID                 primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID             primitive.ObjectID `bson:"userId" json:"userId"`
	OpportunityID      primitive.ObjectID `bson:"opportunityId" json:"opportunityId"`
	HostID             primitive.ObjectID `bson:"hostId" json:"hostId"`
	Status             ApplicationStatus  `bson:"status" json:"status"`
	StatusNote         string             `bson:"statusNote,omitempty" json:"statusNote,omitempty"`
	ApplicationDetails ApplicationDetails `bson:"applicationDetails" json:"applicationDetails"`
	ReviewDetails      ReviewDetails      `bson:"reviewDetails,omitempty" json:"reviewDetails,omitempty"`
	CreatedAt          time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt          time.Time          `bson:"updatedAt" json:"updatedAt"`
}
