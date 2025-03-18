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

	// Check if the comment exists and get the post
	var post models.Post
	filter := bson.M{"comments._id": objID}
	err = col.FindOne(c.Request.Context(), filter).Decode(&post)
	if err != nil {
		return errors.New("comment not found")
	}

	// Check if the user is the author of the comment or the post
	var commentAuthorID primitive.ObjectID
	for _, comment := range post.Comments {
		if comment.ID == objID {
			commentAuthorID = comment.UserId
			break
		}
	}

	userObjID, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil {
		return errors.New("invalid user ID format")
	}

	if commentAuthorID != userObjID && post.UserId != userObjID {
		return errors.New("you are not authorized to delete this comment")
	}

	// Delete the comment
	update := bson.M{"$pull": bson.M{"comments": bson.M{"_id": objID}}}
	_, err = col.UpdateOne(c.Request.Context(), filter, update)
	if err != nil {
		return err
	}
	return nil
}
