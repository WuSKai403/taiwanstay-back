package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// NotificationType 定義通知類型
type NotificationType string

const (
	NotificationTypeApplicationCreated       NotificationType = "APPLICATION_CREATED"
	NotificationTypeApplicationStatusChanged NotificationType = "APPLICATION_STATUS_CHANGED"
)

// Notification 代表一則系統通知
type Notification struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    primitive.ObjectID `bson:"userId" json:"userId"` // 接收者 ID
	Type      NotificationType   `bson:"type" json:"type"`
	Title     string             `bson:"title" json:"title"`
	Message   string             `bson:"message" json:"message"`
	IsRead    bool               `bson:"isRead" json:"isRead"`
	Data      map[string]string  `bson:"data,omitempty" json:"data,omitempty"` // 額外資訊 (e.g., {"applicationId": "..."})
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
}
