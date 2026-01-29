package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// BlacklistedToken represents an invalidated access token
type BlacklistedToken struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Token     string             `bson:"token" json:"token"`
	ExpiresAt time.Time          `bson:"expires_at" json:"expiresAt"`
	CreatedAt time.Time          `bson:"created_at" json:"createdAt"`
}
