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
	postsCol := config.DB.Collection("posts")
	commentsCol := config.DB.Collection("comments")
	savedCol := config.DB.Collection("bookmarked")

	objID, err := primitive.ObjectIDFromHex(c.Query("id"))
	if err != nil {
		return err
	}

	// Delete all comments for the post
	_, err = commentsCol.DeleteMany(context.TODO(), bson.M{"postId": objID})
	if err != nil {
		return err
	}

	// Delete all bookmarks for the post
	_, err = savedCol.DeleteMany(context.TODO(), bson.M{"postId": objID})
	if err != nil {
		return err
	}

	// Delete the post
	_, err = postsCol.DeleteOne(context.TODO(), bson.M{"_id": objID})
	return err
}
