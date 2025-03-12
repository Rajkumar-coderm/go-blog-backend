package blogs

import (
	"context"

	"github.com/Rajkumar-coderm/go-blog-backend/config"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Delete a post
func DeletePost(c *gin.Context) error {
	col := config.DB.Collection("posts")
	objID, err := primitive.ObjectIDFromHex(c.Query("id"))
	if err != nil {
		return err
	}
	_, err = col.DeleteOne(context.TODO(), bson.M{"_id": objID})
	return err
}
