package comments

import (
	"errors"

	"github.com/Rajkumar-coderm/go-blog-backend/config"
	"github.com/Rajkumar-coderm/go-blog-backend/internal/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func DeleteComment(c *gin.Context) error {
	col := config.DB.Collection("posts")
	commentsCol := config.DB.Collection("comments")

	commentID := c.Query("id")
	userID, ok := c.Get("userID")
	if !ok {
		return errors.New("invalid user ID")
	}

	if commentID == "" {
		return errors.New("comment ID is required")
	}

	objID, err := primitive.ObjectIDFromHex(commentID)
	if err != nil {
		return errors.New("invalid comment ID format")
	}

	// Find the comment
	var comment models.Comment
	err = commentsCol.FindOne(c.Request.Context(), bson.M{"_id": objID}).Decode(&comment)
	if err != nil {
		return errors.New("comment not found")
	}

	// Check if the user is the author of the comment or the post
	var post models.Post
	err = col.FindOne(c.Request.Context(), bson.M{"_id": comment.PostId}).Decode(&post)
	if err != nil {
		return errors.New("post not found")
	}

	userObjID, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil {
		return errors.New("invalid user ID format")
	}

	if comment.UserId != userObjID && post.UserId != userObjID {
		return errors.New("you are not authorized to delete this comment")
	}

	// Delete the comment
	_, err = commentsCol.DeleteOne(c.Request.Context(), bson.M{"_id": objID})
	if err != nil {
		return err
	}

	// Decrement commentsCount in posts
	_, err = col.UpdateOne(c.Request.Context(), bson.M{"_id": comment.PostId}, bson.M{
		"$inc": bson.M{"commentsCount": -1},
	})
	if err != nil {
		return err
	}

	return nil
}
