package comments

import (
	"context"
	"errors"
	"time"

	"github.com/Rajkumar-coderm/go-blog-backend/config"
	"github.com/Rajkumar-coderm/go-blog-backend/internal/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CommentPost(c *gin.Context) error {
	col := config.DB.Collection("posts")
	commentsCol := config.DB.Collection("comments")

	userID, ok := c.Get("userID")
	if !ok {
		return errors.New("invalid user ID")
	}

	userCol := config.DB.Collection("users")

	// Convert userID to ObjectID
	_id, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil {
		return errors.New("invalid user ID format")
	}

	// Get user details from users collection
	var user models.User
	err = userCol.FindOne(context.TODO(), bson.M{"_id": _id}).Decode(&user)
	if err != nil {
		return errors.New("user not found")
	}

	// Parse request body
	var request struct {
		PostID  string `json:"postId" binding:"required"`
		Content string `json:"content" binding:"required"`
		ReplyTo string `json:"replyTo,omitempty"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		return errors.New("invalid request payload" + err.Error())
	}

	if request.Content == "" {
		return errors.New("content is required")
	}

	if request.PostID == "" {
		return errors.New("post ID is required")
	}

	// Convert IDs to ObjectID
	userObjID, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil {
		return errors.New("invalid user ID format")
	}

	postObjID, err := primitive.ObjectIDFromHex(request.PostID)
	if err != nil {
		return errors.New("invalid post ID format")
	}

	comment := models.Comment{
		ID:        primitive.NewObjectID(),
		UserId:    userObjID,
		PostId:    postObjID,
		Content:   request.Content,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Replies:   []*models.Comment{},
		Auther: models.Auther{
			UserId:   userObjID,
			Name:     user.FirstName + " " + user.LastName,
			Email:    user.Email,
			UserName: user.Username,
			Bio:      user.Bio,
		},
	}

	if request.ReplyTo != "" {
		replyToID, err := primitive.ObjectIDFromHex(request.ReplyTo)
		if err != nil {
			return errors.New("invalid reply-to comment ID format")
		}
		comment.ParentId = &replyToID
	}

	// Insert comment into comments collection
	_, err = commentsCol.InsertOne(context.TODO(), comment)
	if err != nil {
		return err
	}

	// Increment commentsCount in posts collection
	_, err = col.UpdateOne(context.TODO(), bson.M{"_id": postObjID}, bson.M{
		"$inc": bson.M{"commentsCount": 1},
	})
	if err != nil {
		return err
	}

	return nil
}
