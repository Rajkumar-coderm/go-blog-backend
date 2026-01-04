package models

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SavedPost struct {
	ID        primitive.ObjectID `bson:"_id"`
	UserID    primitive.ObjectID `bson:"userId"`
	PostID    primitive.ObjectID `bson:"postId"`
	CreatedAt time.Time          `bson:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt"`
}
