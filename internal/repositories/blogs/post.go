package blogs

import (
	"context"
	"errors"
	"time"

	"github.com/Rajkumar-coderm/go-blog-backend/config"
	"github.com/Rajkumar-coderm/go-blog-backend/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Create a new blog post
func CreatePost(post *models.Post, userID string) (*mongo.InsertOneResult, error) {
	col := config.DB.Collection("posts")

	authorID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}
	post.ID = primitive.NewObjectID()
	post.UserId = authorID
	post.CreatedAt = time.Now()
	post.UpdatedAt = time.Now()
	return col.InsertOne(context.TODO(), post)
}
