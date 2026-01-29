package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserSession represents a user session in MongoDB
type UserSession struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID       primitive.ObjectID `bson:"user_id" json:"userId"`
	RefreshToken string             `bson:"refresh_token" json:"refreshToken"`
	DeviceInfo   string             `bson:"device_info,omitempty" json:"deviceInfo,omitempty"`
	IPAddress    string             `bson:"ip_address,omitempty" json:"ipAddress,omitempty"`
	ExpiresAt    time.Time          `bson:"expires_at" json:"expiresAt"`
	CreatedAt    time.Time          `bson:"created_at,omitempty" json:"createdAt,omitempty"`
}
