package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ImageStatus string

const (
	ImageStatusPending  ImageStatus = "PENDING"
	ImageStatusApproved ImageStatus = "APPROVED"
	ImageStatusRejected ImageStatus = "REJECTED"
)

// VisionAIRawData 儲存 Vision AI 回傳的原始判斷 (Lickelihood: VERY_UNLIKELY ~ VERY_LIKELY)
type VisionAIRawData struct {
	Adult    string `json:"adult" bson:"adult"`
	Racy     string `json:"racy" bson:"racy"`
	Violence string `json:"violence" bson:"violence"`
	Medical  string `json:"medical" bson:"medical"`
	Spoof    string `json:"spoof" bson:"spoof"`
}

type Image struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID     primitive.ObjectID `bson:"userId" json:"userId"`
	GCSPath    string             `bson:"gcsPath" json:"gcsPath"`     // GCS 中的檔案路徑 (object name)
	PublicURL  string             `bson:"publicUrl" json:"publicUrl"` // 透過 ImageKit 或 GCS 公開的 URL
	Status     ImageStatus        `bson:"status" json:"status"`
	VisionData VisionAIRawData    `bson:"visionData" json:"visionData"`
	CreatedAt  time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt  time.Time          `bson:"updatedAt" json:"updatedAt"`
}
