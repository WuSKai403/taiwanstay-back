package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Bookmark 代表使用者收藏的機會
type Bookmark struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID        string             `bson:"userId" json:"userId"`
	OpportunityID string             `bson:"opportunityId" json:"opportunityId"`
	CreatedAt     time.Time          `bson:"createdAt" json:"createdAt"`
}
