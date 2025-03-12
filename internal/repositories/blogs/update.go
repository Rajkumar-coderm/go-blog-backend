package blogs

import (
	"context"
	"time"

	"github.com/Rajkumar-coderm/go-blog-backend/config"
	"github.com/Rajkumar-coderm/go-blog-backend/internal/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

// Update a post
func UpdatePost(c *gin.Context, post *models.Post) error {
	objID := post.ID
	col := config.DB.Collection("posts")

	post.UpdatedAt = time.Now()
	_, err := col.UpdateOne(context.TODO(), bson.M{"_id": objID}, bson.M{"$set": post})
	return err
}
